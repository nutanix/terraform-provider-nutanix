package tenancyv2_test

// SP Admin Tenant Setup acceptance tests, organised as two cases:
//
//  1. Add domain -> create tenant -> domain allocation without networking -> Update Quota -->
//     destroy the tenant and the allocation.
//  2. Add domain -> create tenant -> domain allocation without networking ->
//     update that allocation with external network config (NAT, and no-NAT
//     when those inputs are set). This case has no delete step.
//
// Gating:
//   - Tests are guarded by acc.TestAccPreCheck (standard PC/NC admin run).
//   - Tests do NOT declare their own provider block. They rely on the default
//     provider (acc.TestAccProviders), which is configured from environment
//     variables.
//
// Infrastructure inputs:
//   - Because the SP admin is also an NC admin with access to every resource,
//     the workflow self-provisions all of its infrastructure instead of
//     relying on user-supplied identifiers:
//       * the domain ext_id is discovered from the nutanix_domains_v2 data
//         source (registered domains),
//       * the cluster ext_id is discovered from nutanix_clusters_v2 (the first
//         non-Prism-Central cluster),
//       * two storage containers are created by the test on that cluster and
//         the pre-existing self-service container is fetched via data source;
//         the allocation uses these across its create/update/delete lifecycle,
//         and
//       * Case 2's networking (NAT and no-NAT transit VPCs plus their overlay
//         external subnets) is created by the test as well.
//   - There are therefore no required environment variables beyond the
//     provider configuration itself.

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Resource addresses for the NC Admin tenant setup workflow.
const (
	ncTenantRes              = "nutanix_tenant_v2.test"
	ncAddDomainRes           = "nutanix_add_domain_v2.test"
	ncDomainAllocationRes    = "nutanix_tenant_domain_allocation_v2.test"
	ncStorageContainerRes1   = "nutanix_storage_containers_v2.test1"
	ncStorageContainerRes2   = "nutanix_storage_containers_v2.test2"
	ncStorageContainerRes3   = "nutanix_storage_containers_v2.test3"
	ncSelfServiceContainerDS = "data.nutanix_storage_containers_v2.self_service"
	ncNatTransitVPCRes       = "nutanix_vpc_v2.nat_transit"
	ncNatOverlaySubnetRes    = "nutanix_subnet_v2.nat_overlay_external"
	ncNoNatTransitVPCRes     = "nutanix_vpc_v2.no_nat_transit"
	ncNoNatOverlaySubnetRes  = "nutanix_subnet_v2.no_nat_overlay_external"
	ncDomainsDS              = "data.nutanix_domains_v2.test"
	ncTenantDS               = "data.nutanix_tenant_v2.test"
	ncTenantsDS              = "data.nutanix_tenants_v2.test"
	ncDomainAllocationDS     = "data.nutanix_tenant_domain_allocation_v2.test"
	ncDomainAllocationsDS    = "data.nutanix_tenant_domain_allocations_v2.test"

	// ncSelfServiceContainerName is the name of the pre-existing self-service
	// storage container the allocation reuses (fetched via data source).
	ncSelfServiceContainerName = "SelfServiceContainer"

	// ncExternallyRoutablePrefix is advertised by the no-NAT allocation config.
	ncExternallyRoutablePrefix = "10.52.0.0/24"
)

// ===========================================================================
// Case 1: tenant + domain allocation without networking, then delete.
//
//	add_domain -> tenant -> domain allocation (projects, resource groups,
//	placement targets; no external_network_config) -> data sources -> destroy
// ===========================================================================

func TestAccNCAdminTenantSetupWithoutNetworking(t *testing.T) {
	r := acctest.RandInt()
	tenantName := fmt.Sprintf("tf-test-tenant-%d", r)
	tenantDesc := "Tenant created by NC admin workflow test"
	quota := ncQuota{Vcpu: 30, Memory: 25, Storage: 40, FloatingIPCount: 25}
	updatedQuota := ncQuota{Vcpu: 40, Memory: 35, Storage: 60, FloatingIPCount: 25}

	storageContainerUUIDs := "projects.0.resource_group.0.placement_targets.0.storage_container_uuids"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckNutanixTenantV2Destroy,
			testAccCheckNutanixTenantDomainAllocationV2Destroy,
		),
		Steps: []resource.TestStep{
			{
				Config: ncTenantSetupConfig(tenantName, tenantDesc, quota, false, 1, ncStageTenant),
				Check:  ncTenantSetupTenantChecks(tenantName, tenantDesc),
			},
			// Create the allocation using the self-service + first created storage containers.
			{
				Config: ncTenantSetupConfig(tenantName, tenantDesc, quota, false, 1, ncStageDataSources),
				Check: resource.ComposeTestCheckFunc(
					ncDomainAllocationChecks(quota),
					ncStorageContainerChecks(storageContainerUUIDs, 1),
					ncExternalNetworkNotConfiguredCheck(),
					ncTenantSetupDataSourceChecks(),
				),
			},
			// Update the quota (same storage containers).
			{
				Config: ncTenantSetupConfig(tenantName, tenantDesc, updatedQuota, false, 1, ncStageDataSources),
				Check: resource.ComposeTestCheckFunc(
					ncDomainAllocationChecks(updatedQuota),
					ncStorageContainerChecks(storageContainerUUIDs, 1),
					ncExternalNetworkNotConfiguredCheck(),
					ncTenantSetupDataSourceChecks(),
				),
			},
			// Add the second created storage container to the allocation.
			//
			// NOTE: there is no "remove storage container" step. The API rejects
			// detaching an already-allocated storage container with
			// TEN-33214 / RESOURCE_REMOVAL_NOT_ALLOWED, so a container can only be
			// added to an existing allocation, never removed.
			{
				Config: ncTenantSetupConfig(tenantName, tenantDesc, updatedQuota, false, 2, ncStageDataSources),
				Check: resource.ComposeTestCheckFunc(
					ncDomainAllocationChecks(updatedQuota),
					ncStorageContainerChecks(storageContainerUUIDs, 2),
					ncExternalNetworkNotConfiguredCheck(),
					ncTenantSetupDataSourceChecks(),
				),
			},
			// Add a third created storage container to the allocation.
			{
				Config: ncTenantSetupConfig(tenantName, tenantDesc, updatedQuota, false, 3, ncStageDataSources),
				Check: resource.ComposeTestCheckFunc(
					ncDomainAllocationChecks(updatedQuota),
					ncStorageContainerChecks(storageContainerUUIDs, 3),
					ncExternalNetworkNotConfiguredCheck(),
					ncTenantSetupDataSourceChecks(),
				),
			},
		},
	})
}

// ncExternalNetworkNotConfiguredCheck asserts the allocation carries an external
// network config block in the NOT_CONFIGURED state. The API always returns the
// block; it stays unconfigured (state NOT_CONFIGURED) until networking is set,
// so the block count is 1 rather than 0.
func ncExternalNetworkNotConfiguredCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "external_network_config.#", "1"),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "external_network_config.0.state_info.0.state", "NOT_CONFIGURED"),
	)
}

// ncStorageContainerChecks verifies the allocation's storage_container_uuids
// set. The self-service container and the first created container are always
// present; createdCount selects how many of the created containers (test1,
// test2, test3) are attached, so the total set size is createdCount + 1 (the
// self-service container). Containers can only be added to an allocation, never
// removed (the API rejects removal with TEN-33214), so createdCount only grows.
func ncStorageContainerChecks(path string, createdCount int) resource.TestCheckFunc {
	createdRes := []string{ncStorageContainerRes1, ncStorageContainerRes2, ncStorageContainerRes3}
	checks := []resource.TestCheckFunc{
		resource.TestCheckTypeSetElemAttrPair(ncDomainAllocationRes, path+".*", ncSelfServiceContainerDS, "storage_containers.0.container_ext_id"),
	}
	for i := 0; i < createdCount; i++ {
		checks = append(checks, resource.TestCheckTypeSetElemAttrPair(ncDomainAllocationRes, path+".*", createdRes[i], "container_ext_id"))
	}
	checks = append(checks, resource.TestCheckResourceAttr(ncDomainAllocationRes, path+".#", fmt.Sprintf("%d", createdCount+1)))
	return resource.ComposeTestCheckFunc(checks...)
}

// ===========================================================================
// Case 2: create the allocation, then update it with networking.
//
//	add_domain -> tenant -> domain allocation without networking ->
//	update the same allocation with external_network_config
//
// There is no delete step. NAT VPC and overlay subnet ext_ids are required;
// no-NAT is included when those ext_ids are set.
// ===========================================================================

func TestAccNCAdminTenantSetupWithNetworking(t *testing.T) {
	r := acctest.RandInt()
	tenantName := fmt.Sprintf("tf-test-tenant-net-%d", r)
	tenantDesc := "Tenant created by NC admin networking workflow test"
	quota := ncQuota{Vcpu: 30, Memory: 25, Storage: 40, FloatingIPCount: 25}

	networkingChecks := []resource.TestCheckFunc{ncDomainAllocationChecks(quota)}
	networkingChecks = append(networkingChecks, ncNetworkingChecks(quota)...)
	networkingChecks = append(networkingChecks, ncTenantSetupDataSourceChecks())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ncTenantSetupConfig(tenantName, tenantDesc, quota, false, 1, ncStageTenant),
				Check:  ncTenantSetupTenantChecks(tenantName, tenantDesc),
			},
			{
				Config: ncTenantSetupConfig(tenantName, tenantDesc, quota, false, 1, ncStageAllocation),
				Check: resource.ComposeTestCheckFunc(
					ncDomainAllocationChecks(quota),
					ncExternalNetworkNotConfiguredCheck(),
				),
			},
			{
				Config: ncTenantSetupConfig(tenantName, tenantDesc, quota, true, 1, ncStageDataSources),
				Check:  resource.ComposeTestCheckFunc(networkingChecks...),
			},
			// Final step: forget every resource from Terraform state without
			// deleting it on the backend, so the framework's automatic post-test
			// destroy is skipped and this test performs no delete API calls.
			{
				Config: ncTenantSetupForgetConfig(),
			},
		},
	})
}

func ncTenantSetupTenantChecks(tenantName, tenantDesc string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		// NOTE: add_domain checks are commented out while the resource is
		// disabled (the domain is pre-shortlisted). The domain is still
		// validated via the nutanix_domains_v2 data source below.
		// resource.TestCheckResourceAttrSet(ncAddDomainRes, "id"),
		// resource.TestCheckResourceAttrSet(ncAddDomainRes, "domain_ext_id"),
		resource.TestCheckResourceAttrSet(ncDomainsDS, "domains.#"),
		resource.TestCheckResourceAttrSet(ncTenantRes, "ext_id"),
		resource.TestCheckResourceAttr(ncTenantRes, "name", tenantName),
		resource.TestCheckResourceAttr(ncTenantRes, "description", tenantDesc),
		// NOTE: admin_info is write-only (not returned on read) and the resource
		// ignores changes to it, so it is not present in state. Its validation is
		// intentionally omitted here.
		// resource.TestCheckResourceAttr(ncTenantRes, "admin_info.#", "1"),
		// resource.TestCheckResourceAttr(ncTenantRes, "admin_info.0.username", "tenant-admin"),
		// resource.TestCheckResourceAttr(ncTenantRes, "admin_info.0.first_name", "First"),
		// resource.TestCheckResourceAttr(ncTenantRes, "admin_info.0.last_name", "Last"),
		// resource.TestCheckResourceAttr(ncTenantRes, "admin_info.0.email", "tenant-admin@example.com"),
	)
}

func ncDomainAllocationChecks(quota ncQuota) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(ncDomainAllocationRes, "ext_id"),
		resource.TestCheckResourceAttrPair(ncDomainAllocationRes, "tenant_ext_id", ncTenantRes, "ext_id"),
		resource.TestCheckResourceAttrPair(ncDomainAllocationRes, "domain_ext_id", ncDomainsDS, "domains.0.domain_ext_id"),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "projects.#", "1"),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "projects.0.quota.#", "1"),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "projects.0.quota.0.vcpu", fmt.Sprintf("%d", quota.Vcpu)),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "projects.0.quota.0.memory", fmt.Sprintf("%d", quota.Memory)),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "projects.0.quota.0.storage", fmt.Sprintf("%d", quota.Storage)),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "projects.0.resource_group.#", "1"),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "projects.0.resource_group.0.placement_targets.#", "1"),
		resource.TestCheckResourceAttrPair(ncDomainAllocationRes, "projects.0.resource_group.0.placement_targets.0.cluster_uuid", ncStorageContainerRes1, "cluster_ext_id"),
	)
}

func ncNetworkingChecks(quota ncQuota) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "external_network_config.#", "1"),
		resource.TestCheckResourceAttrPair(ncDomainAllocationRes, "external_network_config.0.config.0.nat.0.transit_vpc_ext_id", ncNatTransitVPCRes, "ext_id"),
		resource.TestCheckResourceAttrPair(ncDomainAllocationRes, "external_network_config.0.config.0.nat.0.overlay_subnet_ext_id", ncNatOverlaySubnetRes, "ext_id"),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "external_network_config.0.config.0.nat.0.floating_ip_count", fmt.Sprintf("%d", quota.FloatingIPCount)),
		resource.TestCheckResourceAttrPair(ncDomainAllocationRes, "external_network_config.0.config.0.no_nat.0.transit_vpc_ext_id", ncNoNatTransitVPCRes, "ext_id"),
		resource.TestCheckResourceAttrPair(ncDomainAllocationRes, "external_network_config.0.config.0.no_nat.0.overlay_subnet_ext_id", ncNoNatOverlaySubnetRes, "ext_id"),
		resource.TestCheckResourceAttr(ncDomainAllocationRes, "external_network_config.0.config.0.no_nat.0.externally_routable_prefixes.0", ncExternallyRoutablePrefix),
	}
}

func ncTenantSetupDataSourceChecks() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(ncTenantDS, "ext_id", ncTenantRes, "ext_id"),
		resource.TestCheckResourceAttrPair(ncTenantDS, "name", ncTenantRes, "name"),
		resource.TestCheckResourceAttrPair(ncTenantDS, "description", ncTenantRes, "description"),
		checkAttributeLength(ncTenantsDS, "tenants", 1),
		resource.TestCheckResourceAttrPair(ncDomainAllocationDS, "ext_id", ncDomainAllocationRes, "ext_id"),
		resource.TestCheckResourceAttrPair(ncDomainAllocationDS, "domain_ext_id", ncDomainAllocationRes, "domain_ext_id"),
		resource.TestCheckResourceAttrPair(ncDomainAllocationDS, "tenant_ext_id", ncDomainAllocationRes, "tenant_ext_id"),
		checkAttributeLength(ncDomainAllocationsDS, "domain_allocations", 1),
	)
}

// ncQuota holds the quota values for domain allocation projects.
type ncQuota struct {
	Vcpu            int
	Memory          int
	Storage         int
	FloatingIPCount int
}

// ncStage selects how far along the NC Admin tenant setup a config should be
// rendered. Networking is controlled separately so the same allocation can be
// created without it and later updated with it.
type ncStage int

const (
	ncStageTenant      ncStage = iota + 1 // add domain + tenant
	ncStageAllocation                     // + domain allocation
	ncStageDataSources                    // + data sources
)

func ncTenantSetupConfig(tenantName, tenantDesc string, quota ncQuota, withNetworking bool, createdStorageContainers int, stage ncStage) string {
	cfg := ncTenantSetupDiscovery()
	cfg += ncTenantSetupAddDomain()
	cfg += ncTenantSetupTenant(tenantName, tenantDesc)
	if stage >= ncStageAllocation {
		cfg += ncTenantSetupStorageContainers(tenantName)
		if withNetworking {
			cfg += ncTenantSetupNetworking(tenantName)
		}
		cfg += ncTenantSetupDomainAllocation(quota, withNetworking, createdStorageContainers)
	}
	if stage >= ncStageDataSources {
		cfg += ncTenantSetupDataSources()
	}
	return cfg
}

// ncTenantSetupDiscovery discovers the infrastructure the workflow needs so the
// test does not depend on user-supplied identifiers:
//   - the first non-Prism-Central cluster, and
//   - the first registered domain.
func ncTenantSetupDiscovery() string {
	return `
data "nutanix_clusters_v2" "clusters" {}

data "nutanix_domains_v2" "available" {}

locals {
  nc_cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  nc_domain_ext_id = data.nutanix_domains_v2.available.domains[0].domain_ext_id
}
`
}

func ncTenantSetupAddDomain() string {
	// NOTE: nutanix_add_domain_v2 is temporarily commented out. In the target
	// environment the domain is already shortlisted, and the AddDomain API is
	// not idempotent (it returns DOMAIN_ALREADY_EXISTS / HTTP 409 for an
	// existing domain). The domain ext_id is still discovered via the
	// nutanix_domains_v2 data source (see local.nc_domain_ext_id).
	return `
# resource "nutanix_add_domain_v2" "test" {
#   domain_ext_id = local.nc_domain_ext_id
# }

data "nutanix_domains_v2" "test" {}
`
}

// ncTenantSetupStorageContainers creates the two storage containers the
// allocation uses on the discovered cluster, and fetches the pre-existing
// self-service container. The allocation references the self-service container
// and the first created one by default; the second created one is
// added/removed during the test. The spec mirrors the nc_pc_setup module's
// storage containers.
func ncTenantSetupStorageContainers(tenantName string) string {
	return fmt.Sprintf(`
data "nutanix_storage_containers_v2" "self_service" {
  filter = "name eq '%[2]s'"
}

locals {
  nc_self_service_container_ext_id = data.nutanix_storage_containers_v2.self_service.storage_containers[0].container_ext_id
}

resource "nutanix_storage_containers_v2" "test1" {
  name                                     = "%[1]s-sc1"
  cluster_ext_id                           = local.nc_cluster_ext_id
  logical_advertised_capacity_bytes        = 1073741824000
  logical_explicit_reserved_capacity_bytes = 32
  replication_factor                       = 1
  nfs_whitelist_addresses {
    ipv4 {
      value         = "192.168.15.0"
      prefix_length = 32
    }
  }
  erasure_code                          = "OFF"
  is_inline_ec_enabled                  = false
  has_higher_ec_fault_domain_preference = false
  cache_deduplication                   = "OFF"
  on_disk_dedup                         = "OFF"
  is_compression_enabled                = true
  is_internal                           = false
  is_software_encryption_enabled        = false
}

resource "nutanix_storage_containers_v2" "test2" {
  name                                     = "%[1]s-sc2"
  cluster_ext_id                           = local.nc_cluster_ext_id
  logical_advertised_capacity_bytes        = 1073741824000
  logical_explicit_reserved_capacity_bytes = 32
  replication_factor                       = 1
  nfs_whitelist_addresses {
    ipv4 {
      value         = "192.168.15.0"
      prefix_length = 32
    }
  }
  erasure_code                          = "OFF"
  is_inline_ec_enabled                  = false
  has_higher_ec_fault_domain_preference = false
  cache_deduplication                   = "OFF"
  on_disk_dedup                         = "OFF"
  is_compression_enabled                = true
  is_internal                           = false
  is_software_encryption_enabled        = false
}

resource "nutanix_storage_containers_v2" "test3" {
  name                                     = "%[1]s-sc3"
  cluster_ext_id                           = local.nc_cluster_ext_id
  logical_advertised_capacity_bytes        = 1073741824000
  logical_explicit_reserved_capacity_bytes = 32
  replication_factor                       = 1
  nfs_whitelist_addresses {
    ipv4 {
      value         = "192.168.15.0"
      prefix_length = 32
    }
  }
  erasure_code                          = "OFF"
  is_inline_ec_enabled                  = false
  has_higher_ec_fault_domain_preference = false
  cache_deduplication                   = "OFF"
  on_disk_dedup                         = "OFF"
  is_compression_enabled                = true
  is_internal                           = false
  is_software_encryption_enabled        = false
}
`, tenantName, ncSelfServiceContainerName)
}

// ncTenantSetupNetworking creates the NAT and no-NAT transit VPCs plus their
// overlay external subnets used by Case 2's external_network_config. The spec
// mirrors the nc_pc_setup module's networking resources.
func ncTenantSetupNetworking(tenantName string) string {
	return fmt.Sprintf(`
resource "nutanix_vpc_v2" "nat_transit" {
  name        = "%[1]s-nat-transit-vpc"
  description = "Transit VPC for tenant external connectivity (NAT enabled)"
  vpc_type    = "TRANSIT"
}

resource "nutanix_subnet_v2" "nat_overlay_external" {
  name           = "%[1]s-nat-overlay-external"
  description    = "NAT overlay external network for tenant floating IPs"
  vpc_reference  = nutanix_vpc_v2.nat_transit.id
  subnet_type    = "OVERLAY"
  is_external    = true
  is_nat_enabled = true

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "10.57.183.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "10.57.183.1"
      }
      pool_list {
        start_ip { value = "10.57.183.2" }
        end_ip { value = "10.57.183.100" }
      }
    }
  }
}

resource "nutanix_vpc_v2" "no_nat_transit" {
  name        = "%[1]s-no-nat-transit-vpc"
  description = "Transit VPC for tenant external connectivity (No-NAT / Direct Routing)"
  vpc_type    = "TRANSIT"

  externally_routable_prefixes {
    ipv4 {
      ip {
        value         = "10.52.0.0"
        prefix_length = 32
      }
      prefix_length = 24
    }
  }
}

resource "nutanix_subnet_v2" "no_nat_overlay_external" {
  name           = "%[1]s-no-nat-overlay-external"
  description    = "No-NAT overlay external network for direct routing"
  vpc_reference  = nutanix_vpc_v2.no_nat_transit.id
  subnet_type    = "OVERLAY"
  is_external    = true
  is_nat_enabled = false

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "10.128.42.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "10.128.42.1"
      }
      pool_list {
        start_ip { value = "10.128.42.2" }
        end_ip { value = "10.128.42.100" }
      }
    }
  }
}
`, tenantName)
}

func ncTenantSetupTenant(name, desc string) string {
	return fmt.Sprintf(`
resource "nutanix_tenant_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"
  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = "Nutanix.123"
  }
  lifecycle {
    ignore_changes = [admin_info]
  }
}
`, name, desc)
}

func ncTenantSetupDomainAllocation(quota ncQuota, withNetworking bool, createdStorageContainers int) string {
	// Reference the storage containers. The allocation always uses the
	// self-service container; createdStorageContainers selects how many of the
	// created containers (test1, test2, test3) are attached. The test adds
	// containers across steps (containers can only be added, never removed).
	createdRefs := []string{
		"nutanix_storage_containers_v2.test1.container_ext_id",
		"nutanix_storage_containers_v2.test2.container_ext_id",
		"nutanix_storage_containers_v2.test3.container_ext_id",
	}
	storageContainers := "[local.nc_self_service_container_ext_id"
	for i := 0; i < createdStorageContainers; i++ {
		storageContainers += ", " + createdRefs[i]
	}
	storageContainers += "]"

	// Build external network config block, referencing the networking created
	// by ncTenantSetupNetworking. Omitted until the update step.
	externalNetworkConfig := ""
	if withNetworking {
		externalNetworkConfig = fmt.Sprintf(`
  external_network_config {
    config {
      nat {
        transit_vpc_ext_id    = nutanix_vpc_v2.nat_transit.ext_id
        overlay_subnet_ext_id = nutanix_subnet_v2.nat_overlay_external.ext_id
        floating_ip_count     = %[1]d
      }
      no_nat {
        transit_vpc_ext_id           = nutanix_vpc_v2.no_nat_transit.ext_id
        overlay_subnet_ext_id        = nutanix_subnet_v2.no_nat_overlay_external.ext_id
        externally_routable_prefixes = ["%[2]s"]
      }
    }
  }`, quota.FloatingIPCount, ncExternallyRoutablePrefix)
	}

	return fmt.Sprintf(`
resource "nutanix_tenant_domain_allocation_v2" "test" {
  tenant_ext_id = nutanix_tenant_v2.test.ext_id
  domain_ext_id = local.nc_domain_ext_id
  projects {
    quota {
      vcpu    = %[2]d
      memory  = %[3]d
      storage = %[4]d
    }
    resource_group {
      placement_targets {
        cluster_uuid            = local.nc_cluster_ext_id
        storage_container_uuids = %[1]s
      }
    }
  }%[5]s
  depends_on = [
    nutanix_tenant_v2.test,
    nutanix_storage_containers_v2.test1,
    nutanix_storage_containers_v2.test2,
    nutanix_storage_containers_v2.test3,
  ]
}
`, storageContainers, quota.Vcpu, quota.Memory, quota.Storage, externalNetworkConfig)
}

func ncTenantSetupDataSources() string {
	return `
data "nutanix_tenant_v2" "test" {
  ext_id = nutanix_tenant_v2.test.ext_id
}

data "nutanix_tenants_v2" "test" {
  depends_on = [nutanix_tenant_v2.test]
}

data "nutanix_tenant_domain_allocation_v2" "test" {
  ext_id        = nutanix_tenant_domain_allocation_v2.test.id
  tenant_ext_id = nutanix_tenant_v2.test.ext_id
}

data "nutanix_tenant_domain_allocations_v2" "test" {
  tenant_ext_id = nutanix_tenant_v2.test.ext_id
  depends_on    = [nutanix_tenant_domain_allocation_v2.test]
}
`
}

// ncTenantSetupForgetConfig returns a config containing only "removed" blocks
// (with destroy = false) for every resource the networking workflow manages.
// Applying it drops those resources from Terraform state WITHOUT deleting them
// on the backend. After this step the state has no managed resources, so the
// acceptance test framework's automatic post-test destroy is skipped and the
// test performs no delete API calls.
//
// This is required because the domain allocation cannot be deleted while its
// external network config exists (TEN-33217 /
// DOMAIN_ALLOCATION_EXTERNAL_NETWORK_CONFIG_EXISTS), so the resources are
// intentionally left on the live system.
func ncTenantSetupForgetConfig() string {
	addresses := []string{
		"nutanix_tenant_domain_allocation_v2.test",
		"nutanix_tenant_v2.test",
		"nutanix_storage_containers_v2.test1",
		"nutanix_storage_containers_v2.test2",
		"nutanix_storage_containers_v2.test3",
		"nutanix_subnet_v2.nat_overlay_external",
		"nutanix_vpc_v2.nat_transit",
		"nutanix_subnet_v2.no_nat_overlay_external",
		"nutanix_vpc_v2.no_nat_transit",
	}
	cfg := ""
	for _, addr := range addresses {
		cfg += fmt.Sprintf(`
removed {
  from = %s
  lifecycle {
    destroy = false
  }
}
`, addr)
	}
	return cfg
}
