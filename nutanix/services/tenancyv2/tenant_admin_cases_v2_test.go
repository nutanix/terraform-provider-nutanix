package tenancyv2_test

// Tenant Admin acceptance tests, organised as end-to-end WORKFLOWS.
//
// Instead of one isolated test per resource (each rebuilding its own
// prerequisites), resources are grouped into workflows that build on shared,
// chained prerequisites. Each workflow is a single resource.Test whose steps
// grow cumulatively, so step N contains everything from step N-1 plus the next
// resource. This mirrors how a tenant admin actually stands up an environment:
//
//	VM workflow:      image -> transit VPC + overlay subnet -> VM(from image+subnet) ->
//	                  install NGT -> attach categories -> template(from VM) -> deploy template
//	VM actions:       VM -> guest customization update -> assign IP -> clone -> shutdown
//	Network workflow: transit VPC + overlay external subnet -> REGULAR VPC ->
//	                  address group -> entity group -> service group ->
//	                  network security policy
//	IAM workflow:     service account -> user API key (create, update, delete)
//	Membership:       create AD (nutanix_directory_services_v2) ->
//	                  Tenant Admin role (nutanix_roles_v2) ->
//	                  directory search (user + user group) -> two tenant memberships
//
// Gating:
//   - Every test is guarded by acc.TestAccTenantAdminPreCheck, which SKIPS unless
//     NUTANIX_TENANT_NAME is set. A standard PC acceptance run does not set that
//     variable, so these cases never execute alongside the PC suite (and vice
//     versa).
//   - Tests do NOT declare their own provider block. They rely on the default
//     provider (acc.TestAccProviders), which the terraform-plugin-sdk harness
//     configures entirely from environment variables. The provider schema reads
//     NUTANIX_CENTRAL_ENDPOINT, NUTANIX_DOMAIN_NAME, NUTANIX_TENANT_NAME and the
//     credentials (NUTANIX_API_KEY, or NUTANIX_USERNAME + NUTANIX_PASSWORD) via
//     DefaultFunc, so both basic-auth and API-key tenant admins are supported
//     without any provider HCL here. Set NUTANIX_PORT=443 / NUTANIX_INSECURE for
//     the NC endpoint.
//
// Project association ("only one change"):
//   - The project ext_id is NOT hard-coded. It is fetched from the
//     `nutanix_projects_v2` data source and the FIRST entity
//     (projects[0].ext_id) is reused across every resource that supports a
//     project association (project_ext_id, nested project { ext_id }, or the v3
//     project_reference map).

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// projectDataSource is the data source whose first entity supplies the project
// ext_id reused across all tenant admin resources.
const projectDataSource = "data.nutanix_projects_v2.tenant_projects"

// Stable resource addresses used across the workflow steps.
const (
	taImageResource      = "nutanix_images_v2.image"
	taTransitVPCRes      = "nutanix_vpc_v2.transit"
	taSubnetResource     = "nutanix_subnet_v2.subnet"
	taVMResource         = "nutanix_virtual_machine_v2.vm"
	taNGTResource        = "nutanix_ngt_installation_v2.ngt"
	taTemplateRes        = "nutanix_template_v2.template"
	taDeployResource     = "nutanix_deploy_templates_v2.deploy"
	taExtSubnetRes       = "nutanix_subnet_v2.ext"
	taVPCResource        = "nutanix_vpc_v2.vpc"
	taAddrGroupRes       = "nutanix_address_groups_v2.ag"
	taEntityGroupRes     = "nutanix_entity_group_v2.eg"
	taServiceGroupRes    = "nutanix_service_groups_v2.sg"
	taNSPResource        = "nutanix_network_security_policy_v2.nsp"
	taNSRResource        = "nutanix_network_security_rule.nsr"
	taAntiAffinityRes    = "nutanix_vm_anti_affinity_policy_v2.aap"
	taServiceAcctRes     = "nutanix_users_v2.service_account"
	taUserKeyRes         = "nutanix_user_key_v2.key"
	taDirectorySvcRes    = "nutanix_directory_services_v2.ad"
	taUserMembershipRes  = "nutanix_tenant_membership_v2.user"
	taGroupMembershipRes = "nutanix_tenant_membership_v2.group"
	taVolumeGroupRes     = "nutanix_volume_group_v2.vg"
)

// Data source addresses (single + plural) used to verify that every resource is
// also readable through its corresponding v2 data sources. Each resource in a
// workflow gets both a singular data source (looked up by ext_id) and a plural
// data source (filtered by name), and the steps assert them field-by-field
// against the created resource.
const (
	// Workflow 1: VM / compute
	taImageDS          = "data.nutanix_image_v2.image"
	taImagesDS         = "data.nutanix_images_v2.images"
	taSubnetDS         = "data.nutanix_subnet_v2.subnet"
	taSubnetsDS        = "data.nutanix_subnets_v2.subnets"
	taVMDS             = "data.nutanix_virtual_machine_v2.vm"
	taVMsDS            = "data.nutanix_virtual_machines_v2.vms"
	taNGTConfigDS      = "data.nutanix_ngt_configuration_v2.ngt"
	taCategoryDS       = "data.nutanix_category_v2.cat"
	taCategoriesDS     = "data.nutanix_categories_v2.cat"
	taTemplateDS       = "data.nutanix_template_v2.template"
	taTemplatesDS      = "data.nutanix_templates_v2.templates"
	taActVMDS          = "data.nutanix_virtual_machine_v2.act_vm"
	taActVMsDS         = "data.nutanix_virtual_machines_v2.act_vms"
	taNGTUpgradeDS     = "data.nutanix_ngt_configuration_v2.upgrade"
	taAntiAffinityDS   = "data.nutanix_vm_anti_affinity_policy_v2.aap"
	taAntiAffinitiesDS = "data.nutanix_vm_anti_affinity_policies_v2.aap"

	// Workflow 3: network
	taExtSubnetDS     = "data.nutanix_subnet_v2.ext"
	taExtSubnetsDS    = "data.nutanix_subnets_v2.ext"
	taVPCDS           = "data.nutanix_vpc_v2.vpc"
	taVPCsDS          = "data.nutanix_vpcs_v2.vpcs"
	taAddrGroupDS     = "data.nutanix_address_group_v2.ag"
	taAddrGroupsDS    = "data.nutanix_address_groups_v2.ag"
	taEntityGroupDS   = "data.nutanix_entity_group_v2.eg"
	taEntityGroupsDS  = "data.nutanix_entity_groups_v2.eg"
	taServiceGroupDS  = "data.nutanix_service_group_v2.sg"
	taServiceGroupsDS = "data.nutanix_service_groups_v2.sg"
	taNSPDS           = "data.nutanix_network_security_policy_v2.nsp"
	taNSPsDS          = "data.nutanix_network_security_policies_v2.nsp"
	taNSRDS           = "data.nutanix_network_security_rule.nsr"

	// Workflow 4: IAM
	taServiceAcctDS  = "data.nutanix_user_v2.sa"
	taServiceAcctsDS = "data.nutanix_users_v2.sa"
	taUserKeyDS      = "data.nutanix_user_key_v2.key"
	taUserKeysDS     = "data.nutanix_user_keys_v2.keys"

	// Workflow 5: Membership
	taDirectorySvcDS    = "data.nutanix_directory_service_v2.ad"
	taDirectorySvcsDS   = "data.nutanix_directory_services_v2.ad"
	taRolesDS           = "data.nutanix_roles_v2.roles"
	taDomainsDS         = "data.nutanix_registered_domains_v2.domains"
	taUserSearchDS      = "data.nutanix_directory_service_users_search_v2.user"
	taGroupSearchDS     = "data.nutanix_directory_service_users_search_v2.group"
	taUserMembershipDS  = "data.nutanix_tenant_membership_v2.user"
	taGroupMembershipDS = "data.nutanix_tenant_membership_v2.group"
	taMembershipsDS     = "data.nutanix_tenant_memberships_v2.memberships"

	// Workflow 6: Volume Group
	taVolumeGroupDS  = "data.nutanix_volume_group_v2.vg"
	taVolumeGroupsDS = "data.nutanix_volume_groups_v2.vg"
)

// ---------------------------------------------------------------------------
// Shared helpers: base infrastructure data sources
// ---------------------------------------------------------------------------

// tenantAdminBaseConfig returns the shared data sources and locals reused by
// every workflow. There is intentionally NO provider block: the default
// provider (acc.TestAccProviders) is configured from environment variables
// (NUTANIX_CENTRAL_ENDPOINT / NUTANIX_DOMAIN_NAME / NUTANIX_TENANT_NAME plus
// NUTANIX_API_KEY or NUTANIX_USERNAME+PASSWORD).
//   - local.project_ext_id  -> FIRST project entity's ext_id (the "one change")
//   - local.cluster_ext_id  -> cluster from the project's first resource-group
//     placement target (not nutanix_clusters_v2)
//   - local.storage_container_ext_id -> first storage container on that
//     placement target (not nutanix_storage_containers_v2)
//   - local.config/local.vmm -> fixtures from test_config_v2.json
func tenantAdminBaseConfig() string {
	return fmt.Sprintf(`
# Fetch the tenant's projects; the FIRST entity's ext_id is reused across every
# resource that supports a project association.
data "nutanix_projects_v2" "tenant_projects" {}

locals {
  project_ext_id = data.nutanix_projects_v2.tenant_projects.projects[0].ext_id
  config         = jsondecode(file("%s"))
  vmm            = local.config.vmm
}

data "nutanix_resource_groups_v2" "resource_groups" {
  filter = "projectExtId eq '${local.project_ext_id}'"
}

locals {
  cluster_ext_id           = data.nutanix_resource_groups_v2.resource_groups.resource_groups[0].placement_targets[0].cluster_ext_id
  storage_container_ext_id = data.nutanix_resource_groups_v2.resource_groups.resource_groups[0].placement_targets[0].storage_containers[0].ext_id
}

`, filepath)
}

// ===========================================================================
// Workflow 1: VM / Compute
//
//	image -> transit VPC + overlay subnet -> VM(from image+subnet) -> install NGT -> attach categories
//	-> template(from VM) -> deploy template
//
// The steps grow cumulatively, so each step re-declares everything created in
// the previous step. This is the workflow described in the request.
// ===========================================================================

func TestAccTenantAdminVMWorkflow(t *testing.T) {
	r := acctest.RandInt()
	imageName := fmt.Sprintf("tf-ta-image-%d", r)
	subnetName := fmt.Sprintf("tf-ta-subnet-%d", r)
	vmName := fmt.Sprintf("tf-ta-vm-%d", r)
	tplName := fmt.Sprintf("tf-ta-template-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccTenantAdminPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: create image + subnet
				Config: taVMWorkflowConfig(r, imageName, subnetName, vmName, tplName, vmwStageInfra),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taImageResource, "name", imageName),
					resource.TestCheckResourceAttr(taTransitVPCRes, "vpc_type", "TRANSIT"),
					resource.TestCheckResourceAttr(taExtSubnetRes, "subnet_type", "OVERLAY"),
					resource.TestCheckResourceAttr(taExtSubnetRes, "is_external", "true"),
					resource.TestCheckResourceAttr(taExtSubnetRes, "is_nat_enabled", "true"),
					resource.TestCheckResourceAttr(taVPCResource, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttr(taSubnetResource, "name", subnetName),
					resource.TestCheckResourceAttr(taSubnetResource, "subnet_type", "OVERLAY"),
					resource.TestCheckResourceAttr(taSubnetResource, "is_external", "false"),
					resource.TestCheckResourceAttrPair(taImageResource, "project_ext_id", projectDataSource, "projects.0.ext_id"),

					// Image: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taImageDS, "ext_id", taImageResource, "id"),
					resource.TestCheckResourceAttrPair(taImageDS, "name", taImageResource, "name"),
					resource.TestCheckResourceAttrPair(taImageDS, "description", taImageResource, "description"),
					resource.TestCheckResourceAttrPair(taImageDS, "type", taImageResource, "type"),
					resource.TestCheckResourceAttrPair(taImageDS, "project_ext_id", taImageResource, "project_ext_id"),
					resource.TestCheckResourceAttrSet(taImageDS, "create_time"),
					// Image: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taImagesDS, "images.#", "1"),
					resource.TestCheckResourceAttrPair(taImagesDS, "images.0.ext_id", taImageResource, "id"),
					resource.TestCheckResourceAttrPair(taImagesDS, "images.0.name", taImageResource, "name"),
					resource.TestCheckResourceAttrPair(taImagesDS, "images.0.type", taImageResource, "type"),
					resource.TestCheckResourceAttrPair(taImagesDS, "images.0.project_ext_id", taImageResource, "project_ext_id"),

					// Subnet: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taSubnetDS, "ext_id", taSubnetResource, "id"),
					resource.TestCheckResourceAttrPair(taSubnetDS, "name", taSubnetResource, "name"),
					resource.TestCheckResourceAttrPair(taSubnetDS, "description", taSubnetResource, "description"),
					resource.TestCheckResourceAttrPair(taSubnetDS, "subnet_type", taSubnetResource, "subnet_type"),
					resource.TestCheckResourceAttrPair(taSubnetDS, "vpc_reference", taSubnetResource, "vpc_reference"),
					resource.TestCheckResourceAttrPair(taSubnetDS, "is_external", taSubnetResource, "is_external"),
					// Subnet: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taSubnetsDS, "subnets.#", "1"),
					resource.TestCheckResourceAttrPair(taSubnetsDS, "subnets.0.ext_id", taSubnetResource, "id"),
					resource.TestCheckResourceAttrPair(taSubnetsDS, "subnets.0.name", taSubnetResource, "name"),
					resource.TestCheckResourceAttrPair(taSubnetsDS, "subnets.0.subnet_type", taSubnetResource, "subnet_type"),
					resource.TestCheckResourceAttrPair(taSubnetsDS, "subnets.0.vpc_reference", taSubnetResource, "vpc_reference"),
				),
			},
			{
				// Step 2: create the VM from the image + subnet
				Config: taVMWorkflowConfig(r, imageName, subnetName, vmName, tplName, vmwStageVM),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taVMResource, "name", vmName),
					resource.TestCheckResourceAttrPair(taVMResource, "project.0.ext_id", projectDataSource, "projects.0.ext_id"),

					// VM: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taVMDS, "ext_id", taVMResource, "id"),
					resource.TestCheckResourceAttrPair(taVMDS, "name", taVMResource, "name"),
					resource.TestCheckResourceAttrPair(taVMDS, "description", taVMResource, "description"),
					resource.TestCheckResourceAttrPair(taVMDS, "num_sockets", taVMResource, "num_sockets"),
					resource.TestCheckResourceAttrPair(taVMDS, "num_cores_per_socket", taVMResource, "num_cores_per_socket"),
					resource.TestCheckResourceAttrPair(taVMDS, "memory_size_bytes", taVMResource, "memory_size_bytes"),
					resource.TestCheckResourceAttrPair(taVMDS, "cluster.0.ext_id", taVMResource, "cluster.0.ext_id"),
					// VM: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taVMsDS, "vms.#", "1"),
					resource.TestCheckResourceAttrPair(taVMsDS, "vms.0.ext_id", taVMResource, "id"),
					resource.TestCheckResourceAttrPair(taVMsDS, "vms.0.name", taVMResource, "name"),
					resource.TestCheckResourceAttrPair(taVMsDS, "vms.0.num_sockets", taVMResource, "num_sockets"),
					resource.TestCheckResourceAttrPair(taVMsDS, "vms.0.memory_size_bytes", taVMResource, "memory_size_bytes"),
				),
			},
			{
				// Step 3: install NGT on the VM
				Config: taVMWorkflowConfig(r, imageName, subnetName, vmName, tplName, vmwStageNGT),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(taNGTResource, "ext_id"),

					// NGT: single configuration data source vs installation resource.
					resource.TestCheckResourceAttrPair(taNGTConfigDS, "ext_id", taNGTResource, "ext_id"),
					resource.TestCheckResourceAttr(taNGTConfigDS, "is_installed", "true"),
					resource.TestCheckResourceAttr(taNGTConfigDS, "is_enabled", "true"),
					resource.TestCheckResourceAttrPair(taNGTConfigDS, "capablities.#", taNGTResource, "capablities.#"),
				),
			},
			{
				// Step 4: attach categories to the VM (in-place update)
				Config: taVMWorkflowConfig(r, imageName, subnetName, vmName, tplName, vmwStageCategories),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taVMResource, "categories.#", "2"),
					resource.TestCheckResourceAttrPair("nutanix_category_v2.cat.0", "project_ext_id", projectDataSource, "projects.0.ext_id"),

					// Category: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taCategoryDS, "key", "nutanix_category_v2.cat.0", "key"),
					resource.TestCheckResourceAttrPair(taCategoryDS, "value", "nutanix_category_v2.cat.0", "value"),
					resource.TestCheckResourceAttrPair(taCategoryDS, "type", "nutanix_category_v2.cat.0", "type"),
					resource.TestCheckResourceAttrPair(taCategoryDS, "project_ext_id", "nutanix_category_v2.cat.0", "project_ext_id"),
					// Category: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taCategoriesDS, "categories.#", "1"),
					resource.TestCheckResourceAttrPair(taCategoriesDS, "categories.0.ext_id", "nutanix_category_v2.cat.0", "id"),
					resource.TestCheckResourceAttrPair(taCategoriesDS, "categories.0.key", "nutanix_category_v2.cat.0", "key"),
					resource.TestCheckResourceAttrPair(taCategoriesDS, "categories.0.value", "nutanix_category_v2.cat.0", "value"),
					resource.TestCheckResourceAttrPair(taCategoriesDS, "categories.0.project_ext_id", "nutanix_category_v2.cat.0", "project_ext_id"),
				),
			},
			{
				// Step 5: create a template from the VM
				Config: taVMWorkflowConfig(r, imageName, subnetName, vmName, tplName, vmwStageTemplate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taTemplateRes, "template_name", tplName),
					resource.TestCheckResourceAttrPair(taTemplateRes, "project_ext_id", projectDataSource, "projects.0.ext_id"),

					// Template: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taTemplateDS, "ext_id", taTemplateRes, "id"),
					resource.TestCheckResourceAttrPair(taTemplateDS, "template_name", taTemplateRes, "template_name"),
					resource.TestCheckResourceAttrPair(taTemplateDS, "template_description", taTemplateRes, "template_description"),
					resource.TestCheckResourceAttrPair(taTemplateDS, "project_ext_id", taTemplateRes, "project_ext_id"),
					// Template: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taTemplatesDS, "templates.#", "1"),
					resource.TestCheckResourceAttrPair(taTemplatesDS, "templates.0.ext_id", taTemplateRes, "id"),
					resource.TestCheckResourceAttrPair(taTemplatesDS, "templates.0.template_name", taTemplateRes, "template_name"),
					resource.TestCheckResourceAttrPair(taTemplatesDS, "templates.0.project_ext_id", taTemplateRes, "project_ext_id"),
				),
			},
			{
				// Step 6: deploy the template
				Config: taVMWorkflowConfig(r, imageName, subnetName, vmName, tplName, vmwStageDeploy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taDeployResource, "number_of_vms", "1"),
					resource.TestCheckResourceAttrSet(taDeployResource, "ext_id"),
				),
			},
		},
	})
}

// vmwStage selects how far along the cumulative VM workflow a config should be
// rendered.
type vmwStage int

const (
	vmwStageInfra      vmwStage = iota // image + subnet
	vmwStageVM                         // + VM
	vmwStageNGT                        // + NGT install
	vmwStageCategories                 // + categories attached to VM
	vmwStageTemplate                   // + template
	vmwStageDeploy                     // + deploy template
)

// taVMWorkflowConfig builds the cumulative VM workflow config up to the given
// stage. Categories are attached to the VM from the "categories" stage onward.
func taVMWorkflowConfig(r int, imageName, subnetName, vmName, tplName string, stage vmwStage) string {
	cfg := tenantAdminBaseConfig()

	// Categories are declared from the categories stage onward so they can be
	// attached to the VM.
	withCategories := stage >= vmwStageCategories
	if withCategories {
		cfg += taWorkflowCategories(r)
	}

	// Image + transit VPC + overlay subnet always present. Tenant admins
	// cannot create VLAN subnets; overlay requires a VPC first.
	cfg += taWorkflowImage(imageName)
	cfg += taWorkflowSubnet(subnetName, r%200+10)

	if stage >= vmwStageVM {
		cfg += taWorkflowVM(vmName, withCategories)
	}
	if stage >= vmwStageNGT {
		cfg += taWorkflowNGT()
	}
	if stage >= vmwStageTemplate {
		cfg += taWorkflowTemplate(tplName)
	}
	if stage >= vmwStageDeploy {
		cfg += taWorkflowDeploy(vmName)
	}
	return cfg
}

func taWorkflowCategories(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "cat" {
  count          = 2
  key            = "tf-ta-vm-cat-%[1]d"
  value          = "tf-ta-vm-cat-value-${count.index}"
  project_ext_id = local.project_ext_id
}

# Single + plural data sources for the first category, verified against the resource.
data "nutanix_category_v2" "cat" {
  ext_id = nutanix_category_v2.cat[0].id
}

data "nutanix_categories_v2" "cat" {
  filter     = "key eq '${nutanix_category_v2.cat[0].key}' and value eq '${nutanix_category_v2.cat[0].value}'"
  depends_on = [nutanix_category_v2.cat]
}
`, r)
}

func taWorkflowImage(imageName string) string {
	return fmt.Sprintf(`
resource "nutanix_images_v2" "image" {
  name                     = "%[1]s"
  description              = "tenant admin workflow image"
  type                     = "DISK_IMAGE"
  project_ext_id           = local.project_ext_id
  cluster_location_ext_ids = [local.cluster_ext_id]
  source {
    url_source {
      url = local.config.images.ngt_image_url
    }
  }
}

# Single + plural data sources for the image, verified against the resource.
data "nutanix_image_v2" "image" {
  ext_id = nutanix_images_v2.image.id
}

data "nutanix_images_v2" "images" {
  filter     = "name eq '${nutanix_images_v2.image.name}'"
  depends_on = [nutanix_images_v2.image]
}
`, imageName)
}

func taWorkflowSubnet(subnetName string, thirdOctet int) string {
	return fmt.Sprintf(`
resource "nutanix_vpc_v2" "transit" {
  name           = "%[1]s-transit"
  description    = "tenant admin transit vpc for overlay external subnet"
  vpc_type       = "TRANSIT"
  project_ext_id = local.project_ext_id
}

# Overlay external NAT lives on the TRANSIT VPC. VMs must not attach here.
resource "nutanix_subnet_v2" "ext" {
  name           = "%[1]s-ext"
  description    = "tenant admin overlay external NAT network"
  vpc_reference  = nutanix_vpc_v2.transit.id
  subnet_type    = "OVERLAY"
  is_external    = true
  is_nat_enabled = true
  project_ext_id = local.project_ext_id
  ip_config {
    ipv4 {
      ip_subnet {
        ip { value = "10.57.%[2]d.0" }
        prefix_length = 24
      }
      default_gateway_ip { value = "10.57.%[2]d.1" }
      pool_list {
        start_ip { value = "10.57.%[2]d.2" }
        end_ip   { value = "10.57.%[2]d.20" }
      }
    }
  }
  depends_on = [nutanix_vpc_v2.transit]
}

resource "nutanix_vpc_v2" "vpc" {
  name           = "%[1]s-vpc"
  description    = "tenant admin regular vpc for vm overlay subnet"
  vpc_type       = "REGULAR"
  project_ext_id = local.project_ext_id
  external_subnets {
    subnet_reference = nutanix_subnet_v2.ext.id
  }
  depends_on = [nutanix_subnet_v2.ext]
}

# Internal overlay on the REGULAR VPC — this is the VM NIC attachment target.
resource "nutanix_subnet_v2" "subnet" {
  name           = "%[1]s"
  description    = "tenant admin overlay subnet for vm nics"
  vpc_reference  = nutanix_vpc_v2.vpc.id
  subnet_type    = "OVERLAY"
  project_ext_id = local.project_ext_id
  ip_config {
    ipv4 {
      ip_subnet {
        ip { value = "192.168.%[2]d.0" }
        prefix_length = 24
      }
      default_gateway_ip { value = "192.168.%[2]d.1" }
      pool_list {
        start_ip { value = "192.168.%[2]d.10" }
        end_ip   { value = "192.168.%[2]d.20" }
      }
    }
  }
  depends_on = [nutanix_vpc_v2.vpc]
}

# Single + plural data sources for the (internal) subnet, verified against the resource.
data "nutanix_subnet_v2" "subnet" {
  ext_id = nutanix_subnet_v2.subnet.id
}

data "nutanix_subnets_v2" "subnets" {
  filter     = "name eq '${nutanix_subnet_v2.subnet.name}'"
  depends_on = [nutanix_subnet_v2.subnet]
}
`, subnetName, thirdOctet)
}

func taWorkflowVM(vmName string, withCategories bool) string {
	categoriesBlock := ""
	if withCategories {
		categoriesBlock = `
  dynamic "categories" {
    for_each = nutanix_category_v2.cat
    content {
      ext_id = categories.value.id
    }
  }`
	}
	return fmt.Sprintf(`
resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "%[1]s"
  description          = "tenant admin workflow vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  memory_size_bytes    = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }
  project {
    ext_id = local.project_ext_id
  }
  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        data_source {
          reference {
            image_reference {
              image_ext_id = nutanix_images_v2.image.id
            }
          }
        }
        storage_container {
          ext_id = local.storage_container_ext_id
        }
      }
    }
  }
  nics {
    nic_backing_info {
      virtual_ethernet_nic {
        is_connected = true
      }
    }
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        ipv4_config {
          should_assign_ip = true
        }
        subnet {
          ext_id = nutanix_subnet_v2.subnet.id
        }
        vlan_mode = "ACCESS"
      }
    }
  }
	cd_roms {
    disk_address {
      bus_type = "IDE"
      index    = 1
    }
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  power_state = "ON"
  %[2]s
  lifecycle {
    ignore_changes = [
      guest_tools,
      nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config.0.should_assign_ip,
    ]
  }
  depends_on = [nutanix_images_v2.image, nutanix_subnet_v2.subnet]
}

# Single + plural data sources for the VM, verified against the resource.
data "nutanix_virtual_machine_v2" "vm" {
  ext_id = nutanix_virtual_machine_v2.vm.id
}

data "nutanix_virtual_machines_v2" "vms" {
  filter     = "name eq '${nutanix_virtual_machine_v2.vm.name}'"
  depends_on = [nutanix_virtual_machine_v2.vm]
}
`, vmName, categoriesBlock)
}

func taWorkflowNGT() string {
	return `
# insert the NGT ISO
resource "nutanix_ngt_insert_iso_v2" "insert-iso" {
  ext_id         = nutanix_virtual_machine_v2.vm.id
  capablities    = ["SELF_SERVICE_RESTORE","VSS_SNAPSHOT"]
  is_config_only = false
	lifecycle {
    ignore_changes = [capablities]
  }
}

resource "nutanix_ngt_installation_v2" "ngt" {
  ext_id = nutanix_virtual_machine_v2.vm.id
  credential {
    username = local.vmm.ngt.credential.username
		password = local.vmm.ngt.credential.password
  }
  reboot_preference {
    schedule_type = "IMMEDIATE"
  }
  capablities = ["SELF_SERVICE_RESTORE", "VSS_SNAPSHOT"]
	depends_on = [nutanix_ngt_insert_iso_v2.insert-iso]
	lifecycle {
    ignore_changes = [capablities]
  }
}

# NGT has no plural data source; the single configuration data source is keyed
# by the VM ext_id and verified against the installation resource.
data "nutanix_ngt_configuration_v2" "ngt" {
  ext_id     = nutanix_virtual_machine_v2.vm.id
  depends_on = [nutanix_ngt_installation_v2.ngt]
}
`
}

func taWorkflowTemplate(tplName string) string {
	return fmt.Sprintf(`
resource "nutanix_template_v2" "template" {
  template_name        = "%[1]s"
  template_description = "tenant admin workflow template"
  project_ext_id       = local.project_ext_id
  template_version_spec {
    version_source {
      template_vm_reference {
        ext_id = nutanix_virtual_machine_v2.vm.id
      }
    }
  }
  depends_on = [nutanix_virtual_machine_v2.vm]
}

# Single + plural data sources for the template, verified against the resource.
data "nutanix_template_v2" "template" {
  ext_id = nutanix_template_v2.template.id
}

data "nutanix_templates_v2" "templates" {
  filter     = "templateName eq '${nutanix_template_v2.template.template_name}'"
  depends_on = [nutanix_template_v2.template]
}
`, tplName)
}

func taWorkflowDeploy(vmName string) string {
	return fmt.Sprintf(`
resource "nutanix_deploy_templates_v2" "deploy" {
  ext_id            = nutanix_template_v2.template.id
  number_of_vms     = 1
  cluster_reference = local.cluster_ext_id
  override_vm_config_map {
    name                 = "%[1]s-deploy"
    memory_size_bytes    = 4294967296
    num_sockets          = 2
    num_cores_per_socket = 1
    num_threads_per_core = 1
    nics {
      nic_backing_info {
        virtual_ethernet_nic {
          is_connected = true
          model        = "VIRTIO"
        }
      }
      nic_network_info {
        virtual_ethernet_nic_network_info {
          nic_type = "NORMAL_NIC"
          ipv4_config {
            should_assign_ip = true
          }
          subnet {
            ext_id = nutanix_subnet_v2.subnet.id
          }
          vlan_mode = "ACCESS"
        }
      }
    }
  }
  depends_on = [nutanix_template_v2.template]
}
`, vmName)
}

// ===========================================================================
// Workflow 2: VM lifecycle actions
//
//	VM (Off) -> guest customization update -> clone -> power on
// ===========================================================================

func TestAccTenantAdminVMActionsWorkflow(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-ta-act-vm-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccTenantAdminPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: VM Off + guest customization (GC is rejected if the VM is On)
				Config: taVMActionsConfig(vmName, actStageGCUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taVMResource, "name", vmName),
					resource.TestCheckResourceAttrSet("nutanix_vm_gc_update_v2.gc", "ext_id"),

					// VM: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taActVMDS, "ext_id", taVMResource, "id"),
					resource.TestCheckResourceAttrPair(taActVMDS, "name", taVMResource, "name"),
					resource.TestCheckResourceAttrPair(taActVMDS, "num_sockets", taVMResource, "num_sockets"),
					resource.TestCheckResourceAttrPair(taActVMDS, "num_cores_per_socket", taVMResource, "num_cores_per_socket"),
					resource.TestCheckResourceAttrPair(taActVMDS, "memory_size_bytes", taVMResource, "memory_size_bytes"),
					resource.TestCheckResourceAttrPair(taActVMDS, "cluster.0.ext_id", taVMResource, "cluster.0.ext_id"),
					// VM: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taActVMsDS, "vms.#", "1"),
					resource.TestCheckResourceAttrPair(taActVMsDS, "vms.0.ext_id", taVMResource, "id"),
					resource.TestCheckResourceAttrPair(taActVMsDS, "vms.0.name", taVMResource, "name"),
				),
			},
			{
				// Step 2: + clone the VM
				Config: taVMActionsConfig(vmName, actStageClone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_vm_clone_v2.clone", "name", vmName+"-clone"),
				),
			},
			{
				// Step 3: power the VM On.
				Config: taVMActionsConfig(vmName, actStagePowerOn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taVMResource, "power_state", "ON"),
				),
			},
		},
	})
}

type actStage int

const (
	actStageGCUpdate actStage = iota
	actStageClone
	actStagePowerOn
	actStageShutdown
)

func taVMActionsConfig(vmName string, stage actStage) string {
	powerState := "OFF"
	ignorePower := false
	if stage >= actStagePowerOn {
		powerState = "ON"
	}
	if stage >= actStageShutdown {
		ignorePower = true
	}
	cfg := tenantAdminBaseConfig() + taActionsVM(vmName, powerState, ignorePower)

	// gc update present from the first stage
	cfg += `
locals {
  gs = base64encode("#cloud-config\nusers:\n  - name: ubuntu\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']")
}

resource "nutanix_vm_gc_update_v2" "gc" {
  ext_id = nutanix_virtual_machine_v2.vm.id
  config {
    cloud_init {
      cloud_init_script {
        user_data {
          value = local.gs
        }
      }
    }
  }
  depends_on = [nutanix_virtual_machine_v2.vm]
}
`
	if stage >= actStageClone {
		cfg += fmt.Sprintf(`
resource "nutanix_vm_clone_v2" "clone" {
  vm_ext_id            = nutanix_virtual_machine_v2.vm.id
  name                 = "%[1]s-clone"
  num_sockets          = 2
  num_cores_per_socket = 2
  num_threads_per_core = 2
}
`, vmName)
	}
	if stage >= actStageShutdown {
		cfg += `
resource "nutanix_vm_shutdown_action_v2" "shutdown" {
  ext_id = nutanix_virtual_machine_v2.vm.id
  action = "guest_shutdown"
  depends_on = [nutanix_virtual_machine_v2.vm]
}
`
	}
	return cfg
}

// taActionsVM builds the actions VM. Guest customization requires the VM Off;
// guest_shutdown requires it On. powerState is OFF through GC/clone and ON
// from the power-on step. ignorePowerState is set on shutdown so the action's
// transition to OFF does not race the VM resource.
func taActionsVM(vmName, powerState string, ignorePowerState bool) string {
	ignoreChanges := "[guest_customization, guest_tools, cd_roms]"
	if ignorePowerState {
		ignoreChanges = "[guest_customization, guest_tools, cd_roms, power_state]"
	}
	return fmt.Sprintf(`
resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "%[1]s"
  description          = "tenant admin vm actions"
  num_cores_per_socket = 1
  num_sockets          = 1
  memory_size_bytes    = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }
  project {
    ext_id = local.project_ext_id
  }
  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        disk_size_bytes = "1073741824"
        storage_container {
          ext_id = local.storage_container_ext_id
        }
      }
    }
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  power_state = "%[2]s"
  lifecycle {
    ignore_changes = %[3]s
  }
}

# Single + plural data sources for the actions VM, verified against the resource.
# The action resources (gc update, assign ip, clone, shutdown) are one-shot
# operations and have no corresponding data sources.
data "nutanix_virtual_machine_v2" "act_vm" {
  ext_id = nutanix_virtual_machine_v2.vm.id
}

data "nutanix_virtual_machines_v2" "act_vms" {
  filter     = "name eq '${nutanix_virtual_machine_v2.vm.name}'"
  depends_on = [nutanix_virtual_machine_v2.vm]
}
`, vmName, powerState, ignoreChanges)
}

// ===========================================================================
// NGT upgrade (standalone) - references a pre-existing NGT-enabled VM discovered
// by name from the test config, so it is kept separate from the creation
// workflow.
// ===========================================================================

func TestAccTenantAdminNGTUpgrade(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccTenantAdminPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tenantAdminBaseConfig() + `
data "nutanix_virtual_machines_v2" "ngt-vm-upgrade" {
  filter = "name eq '${local.vmm.ngt.ngt_upgrade_vm_name}'"
}

resource "nutanix_ngt_upgrade_v2" "upgrade" {
  ext_id = data.nutanix_virtual_machines_v2.ngt-vm-upgrade.vms[0].ext_id
  reboot_preference {
    schedule_type = "IMMEDIATE"
  }
}

# NGT has no plural data source; verify the single configuration data source
# against the upgraded VM.
data "nutanix_ngt_configuration_v2" "upgrade" {
  ext_id     = nutanix_ngt_upgrade_v2.upgrade.ext_id
  depends_on = [nutanix_ngt_upgrade_v2.upgrade]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_ngt_upgrade_v2.upgrade", "ext_id"),

					// NGT: single configuration data source vs upgrade resource.
					resource.TestCheckResourceAttrPair(taNGTUpgradeDS, "ext_id", "nutanix_ngt_upgrade_v2.upgrade", "ext_id"),
					resource.TestCheckResourceAttr(taNGTUpgradeDS, "is_installed", "true"),
					resource.TestCheckResourceAttr(taNGTUpgradeDS, "is_enabled", "true"),
				),
			},
		},
	})
}

// ===========================================================================
// VM anti-affinity policy (standalone) - creates its own categories scoped to
// the tenant project. Tenant admins cannot create unscoped (no project_ext_id)
// categories; that returns IAM PLAT-10007. The policy resource itself has no
// project association.
// ===========================================================================

func TestAccTenantAdminVMAntiAffinityPolicy(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-ta-aap-%d", r)
	updatedName := fmt.Sprintf("tf-ta-aap-updated-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccTenantAdminPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: taAntiAffinityConfig(r, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taAntiAffinityRes, "name", name),

					// Anti-affinity: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taAntiAffinityDS, "ext_id", taAntiAffinityRes, "id"),
					resource.TestCheckResourceAttrPair(taAntiAffinityDS, "name", taAntiAffinityRes, "name"),
					resource.TestCheckResourceAttrPair(taAntiAffinityDS, "description", taAntiAffinityRes, "description"),
					resource.TestCheckResourceAttrPair(taAntiAffinityDS, "categories.#", taAntiAffinityRes, "categories.#"),
					// Anti-affinity: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taAntiAffinitiesDS, "policies.#", "1"),
					resource.TestCheckResourceAttrPair(taAntiAffinitiesDS, "policies.0.ext_id", taAntiAffinityRes, "id"),
					resource.TestCheckResourceAttrPair(taAntiAffinitiesDS, "policies.0.name", taAntiAffinityRes, "name"),
				),
			},
			{
				Config: taAntiAffinityConfig(r, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taAntiAffinityRes, "name", updatedName),

					// Anti-affinity: data sources reflect the updated name.
					resource.TestCheckResourceAttrPair(taAntiAffinityDS, "name", taAntiAffinityRes, "name"),
					resource.TestCheckResourceAttrPair(taAntiAffinitiesDS, "policies.0.name", taAntiAffinityRes, "name"),
				),
			},
		},
	})
}

func taAntiAffinityConfig(r int, name string) string {
	return tenantAdminBaseConfig() + fmt.Sprintf(`
resource "nutanix_category_v2" "aap_cat" {
  count          = 2
  key            = "tf-ta-aap-cat-%[1]d"
  value          = "tf-ta-aap-cat-value-${count.index}"
  project_ext_id = local.project_ext_id
}

resource "nutanix_vm_anti_affinity_policy_v2" "aap" {
  name        = "%[2]s"
  description = "tenant admin anti affinity policy"
  categories  = nutanix_category_v2.aap_cat[*].id
	project_ext_id = local.project_ext_id
}

data "nutanix_vm_anti_affinity_policy_v2" "aap" {
  ext_id     = nutanix_vm_anti_affinity_policy_v2.aap.id
  depends_on = [nutanix_vm_anti_affinity_policy_v2.aap]
}

data "nutanix_vm_anti_affinity_policies_v2" "aap" {
  filter     = "name eq '${nutanix_vm_anti_affinity_policy_v2.aap.name}'"
  depends_on = [nutanix_vm_anti_affinity_policy_v2.aap]
}
`, r, name)
}

// ===========================================================================
// Workflow 3: Network
//
//	transit VPC + overlay external subnet -> REGULAR VPC -> address group ->
//	entity group -> service group -> network security policy
//
// Steps grow cumulatively. Every resource that supports it is associated to the
// FIRST project entity. Floating IPs are omitted: tenant admins use overlay
// networking and microseg groups instead.
// ===========================================================================

func TestAccTenantAdminNetworkWorkflow(t *testing.T) {
	r := acctest.RandInt()
	subnetName := fmt.Sprintf("tf-ta-ext-subnet-%d", r)
	vpcName := fmt.Sprintf("tf-ta-vpc-%d", r)
	agName := fmt.Sprintf("tf-ta-ag-%d", r)
	egName := fmt.Sprintf("tf-ta-eg-%d", r)
	sgName := fmt.Sprintf("tf-ta-sg-%d", r)
	nspName := fmt.Sprintf("tf-ta-nsp-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccTenantAdminPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: create a transit VPC + overlay external NAT subnet
				Config: taNetworkWorkflowConfig(r, subnetName, vpcName, agName, egName, sgName, nspName, netStageSubnet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taExtSubnetRes, "name", subnetName),
					resource.TestCheckResourceAttr(taExtSubnetRes, "subnet_type", "OVERLAY"),
					resource.TestCheckResourceAttr(taExtSubnetRes, "is_external", "true"),
					resource.TestCheckResourceAttr(taExtSubnetRes, "is_nat_enabled", "true"),

					// External subnet: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taExtSubnetDS, "ext_id", taExtSubnetRes, "id"),
					resource.TestCheckResourceAttrPair(taExtSubnetDS, "name", taExtSubnetRes, "name"),
					resource.TestCheckResourceAttrPair(taExtSubnetDS, "subnet_type", taExtSubnetRes, "subnet_type"),
					resource.TestCheckResourceAttrPair(taExtSubnetDS, "is_external", taExtSubnetRes, "is_external"),
					resource.TestCheckResourceAttrPair(taExtSubnetDS, "is_nat_enabled", taExtSubnetRes, "is_nat_enabled"),
					resource.TestCheckResourceAttrPair(taExtSubnetDS, "vpc_reference", taExtSubnetRes, "vpc_reference"),
					// External subnet: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taExtSubnetsDS, "subnets.#", "1"),
					resource.TestCheckResourceAttrPair(taExtSubnetsDS, "subnets.0.ext_id", taExtSubnetRes, "id"),
					resource.TestCheckResourceAttrPair(taExtSubnetsDS, "subnets.0.name", taExtSubnetRes, "name"),
					resource.TestCheckResourceAttrPair(taExtSubnetsDS, "subnets.0.is_external", taExtSubnetRes, "is_external"),
					resource.TestCheckResourceAttrPair(taExtSubnetsDS, "subnets.0.subnet_type", taExtSubnetRes, "subnet_type"),
				),
			},
			{
				// Step 2: create a VPC on the external subnet
				Config: taNetworkWorkflowConfig(r, subnetName, vpcName, agName, egName, sgName, nspName, netStageVPC),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taVPCResource, "name", vpcName),
					resource.TestCheckResourceAttrPair(taVPCResource, "project_ext_id", projectDataSource, "projects.0.ext_id"),

					// VPC: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taVPCDS, "ext_id", taVPCResource, "id"),
					resource.TestCheckResourceAttrPair(taVPCDS, "name", taVPCResource, "name"),
					resource.TestCheckResourceAttrPair(taVPCDS, "description", taVPCResource, "description"),
					resource.TestCheckResourceAttrPair(taVPCDS, "vpc_type", taVPCResource, "vpc_type"),
					resource.TestCheckResourceAttrPair(taVPCDS, "project_ext_id", taVPCResource, "project_ext_id"),
					// VPC: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taVPCsDS, "vpcs.#", "1"),
					resource.TestCheckResourceAttrPair(taVPCsDS, "vpcs.0.ext_id", taVPCResource, "id"),
					resource.TestCheckResourceAttrPair(taVPCsDS, "vpcs.0.name", taVPCResource, "name"),
					resource.TestCheckResourceAttrPair(taVPCsDS, "vpcs.0.project_ext_id", taVPCResource, "project_ext_id"),
				),
			},
			{
				// Step 3: create an address group
				Config: taNetworkWorkflowConfig(r, subnetName, vpcName, agName, egName, sgName, nspName, netStageAddrGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taAddrGroupRes, "name", agName),
					resource.TestCheckResourceAttrPair(taAddrGroupRes, "project_ext_id", projectDataSource, "projects.0.ext_id"),

					// Address group: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taAddrGroupDS, "ext_id", taAddrGroupRes, "id"),
					resource.TestCheckResourceAttrPair(taAddrGroupDS, "name", taAddrGroupRes, "name"),
					resource.TestCheckResourceAttrPair(taAddrGroupDS, "description", taAddrGroupRes, "description"),
					resource.TestCheckResourceAttrPair(taAddrGroupDS, "ipv4_addresses.#", taAddrGroupRes, "ipv4_addresses.#"),
					resource.TestCheckResourceAttrPair(taAddrGroupDS, "ipv4_addresses.0.value", taAddrGroupRes, "ipv4_addresses.0.value"),
					resource.TestCheckResourceAttrPair(taAddrGroupDS, "ipv4_addresses.0.prefix_length", taAddrGroupRes, "ipv4_addresses.0.prefix_length"),
					// Address group: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taAddrGroupsDS, "address_groups.#", "1"),
					resource.TestCheckResourceAttrPair(taAddrGroupsDS, "address_groups.0.ext_id", taAddrGroupRes, "id"),
					resource.TestCheckResourceAttrPair(taAddrGroupsDS, "address_groups.0.name", taAddrGroupRes, "name"),
				),
			},
			{
				// Step 4: create an entity group that selects the address group
				Config: taNetworkWorkflowConfig(r, subnetName, vpcName, agName, egName, sgName, nspName, netStageEntityGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taEntityGroupRes, "name", egName),
					resource.TestCheckResourceAttrPair(taEntityGroupRes, "project_ext_id", projectDataSource, "projects.0.ext_id"),
					resource.TestCheckResourceAttr(taEntityGroupRes, "allowed_config.0.entities.0.type", "ADDRESS_GROUP"),
					resource.TestCheckResourceAttr(taEntityGroupRes, "allowed_config.0.entities.0.selected_by", "EXT_ID"),

					// Entity group: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taEntityGroupDS, "ext_id", taEntityGroupRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taEntityGroupDS, "name", taEntityGroupRes, "name"),
					resource.TestCheckResourceAttrPair(taEntityGroupDS, "description", taEntityGroupRes, "description"),
					resource.TestCheckResourceAttrPair(taEntityGroupDS, "project_ext_id", taEntityGroupRes, "project_ext_id"),
					// Entity group: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taEntityGroupsDS, "entity_groups.#", "1"),
					resource.TestCheckResourceAttrPair(taEntityGroupsDS, "entity_groups.0.ext_id", taEntityGroupRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taEntityGroupsDS, "entity_groups.0.name", taEntityGroupRes, "name"),
					resource.TestCheckResourceAttrPair(taEntityGroupsDS, "entity_groups.0.project_ext_id", taEntityGroupRes, "project_ext_id"),
				),
			},
			{
				// Step 5: create a service group
				Config: taNetworkWorkflowConfig(r, subnetName, vpcName, agName, egName, sgName, nspName, netStageServiceGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taServiceGroupRes, "name", sgName),
					resource.TestCheckResourceAttrPair(taServiceGroupRes, "project_ext_id", projectDataSource, "projects.0.ext_id"),
					resource.TestCheckResourceAttr(taServiceGroupRes, "tcp_services.#", "1"),
					resource.TestCheckResourceAttr(taServiceGroupRes, "udp_services.#", "1"),

					// Service group: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taServiceGroupDS, "ext_id", taServiceGroupRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taServiceGroupDS, "name", taServiceGroupRes, "name"),
					resource.TestCheckResourceAttrPair(taServiceGroupDS, "description", taServiceGroupRes, "description"),
					resource.TestCheckResourceAttrPair(taServiceGroupDS, "project_ext_id", taServiceGroupRes, "project_ext_id"),
					// Service group: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taServiceGroupsDS, "service_groups.#", "1"),
					resource.TestCheckResourceAttrPair(taServiceGroupsDS, "service_groups.0.ext_id", taServiceGroupRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taServiceGroupsDS, "service_groups.0.name", taServiceGroupRes, "name"),
					resource.TestCheckResourceAttrPair(taServiceGroupsDS, "service_groups.0.project_ext_id", taServiceGroupRes, "project_ext_id"),
				),
			},
			{
				// Step 6: create a network security policy
				Config: taNetworkWorkflowConfig(r, subnetName, vpcName, agName, egName, sgName, nspName, netStageNSP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taNSPResource, "name", nspName),
					resource.TestCheckResourceAttrPair(taNSPResource, "project_ext_id", projectDataSource, "projects.0.ext_id"),

					// NSP: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taNSPDS, "ext_id", taNSPResource, "id"),
					resource.TestCheckResourceAttrPair(taNSPDS, "name", taNSPResource, "name"),
					resource.TestCheckResourceAttrPair(taNSPDS, "description", taNSPResource, "description"),
					resource.TestCheckResourceAttrPair(taNSPDS, "type", taNSPResource, "type"),
					resource.TestCheckResourceAttrPair(taNSPDS, "state", taNSPResource, "state"),
					resource.TestCheckResourceAttrPair(taNSPDS, "is_hitlog_enabled", taNSPResource, "is_hitlog_enabled"),
					// NSP: plural data source vs resource (deep). List attr is network_policies.
					resource.TestCheckResourceAttr(taNSPsDS, "network_policies.#", "1"),
					resource.TestCheckResourceAttrPair(taNSPsDS, "network_policies.0.ext_id", taNSPResource, "id"),
					resource.TestCheckResourceAttrPair(taNSPsDS, "network_policies.0.name", taNSPResource, "name"),
					resource.TestCheckResourceAttrPair(taNSPsDS, "network_policies.0.type", taNSPResource, "type"),
				),
			},
		},
	})
}

type netStage int

const (
	netStageSubnet netStage = iota
	netStageVPC
	netStageAddrGroup
	netStageEntityGroup
	netStageServiceGroup
	netStageNSP
)

func taNetworkWorkflowConfig(r int, subnetName, vpcName, agName, egName, sgName, nspName string, stage netStage) string {
	cfg := tenantAdminBaseConfig() + taNetExternalSubnet(subnetName, r%200+10)

	if stage >= netStageVPC {
		cfg += fmt.Sprintf(`
resource "nutanix_vpc_v2" "vpc" {
  name           = "%[1]s"
  description    = "tenant admin workflow vpc"
  vpc_type       = "REGULAR"
  project_ext_id = local.project_ext_id
  external_subnets {
    subnet_reference = nutanix_subnet_v2.ext.id
  }
  depends_on = [nutanix_subnet_v2.ext]
}

# Single + plural data sources for the VPC, verified against the resource.
data "nutanix_vpc_v2" "vpc" {
  ext_id     = nutanix_vpc_v2.vpc.id
  depends_on = [nutanix_vpc_v2.vpc]
}

data "nutanix_vpcs_v2" "vpcs" {
  filter     = "name eq '${nutanix_vpc_v2.vpc.name}'"
  depends_on = [nutanix_vpc_v2.vpc]
}
`, vpcName)
	}
	if stage >= netStageAddrGroup {
		cfg += fmt.Sprintf(`
resource "nutanix_address_groups_v2" "ag" {
  name           = "%[1]s"
  description    = "tenant admin workflow address group"
  project_ext_id = local.project_ext_id
  ipv4_addresses {
    value         = "10.0.0.0"
    prefix_length = 24
  }
}

# Single + plural data sources for the address group, verified against the resource.
data "nutanix_address_group_v2" "ag" {
  ext_id     = nutanix_address_groups_v2.ag.id
  depends_on = [nutanix_address_groups_v2.ag]
}

data "nutanix_address_groups_v2" "ag" {
  filter     = "name eq '${nutanix_address_groups_v2.ag.name}'"
  depends_on = [nutanix_address_groups_v2.ag]
}
`, agName)
	}
	if stage >= netStageEntityGroup {
		cfg += fmt.Sprintf(`
resource "nutanix_entity_group_v2" "eg" {
  name           = "%[1]s"
  description    = "tenant admin workflow entity group"
  project_ext_id = local.project_ext_id
  allowed_config {
    entities {
      type        = "ADDRESS_GROUP"
      selected_by = "EXT_ID"
      reference_ext_ids = [
        nutanix_address_groups_v2.ag.id,
      ]
    }
  }
  depends_on = [nutanix_address_groups_v2.ag]
}

# Single + plural data sources for the entity group, verified against the resource.
data "nutanix_entity_group_v2" "eg" {
  ext_id     = nutanix_entity_group_v2.eg.ext_id
  depends_on = [nutanix_entity_group_v2.eg]
}

data "nutanix_entity_groups_v2" "eg" {
  filter     = "name eq '${nutanix_entity_group_v2.eg.name}'"
  depends_on = [nutanix_entity_group_v2.eg]
}
`, egName)
	}
	if stage >= netStageServiceGroup {
		cfg += fmt.Sprintf(`
resource "nutanix_service_groups_v2" "sg" {
  name           = "%[1]s"
  description    = "tenant admin workflow service group"
  project_ext_id = local.project_ext_id
  tcp_services {
    start_port = "232"
    end_port   = "232"
  }
  udp_services {
    start_port = "232"
    end_port   = "232"
  }
}

# Single + plural data sources for the service group, verified against the resource.
data "nutanix_service_group_v2" "sg" {
  ext_id     = nutanix_service_groups_v2.sg.ext_id
  depends_on = [nutanix_service_groups_v2.sg]
}

data "nutanix_service_groups_v2" "sg" {
  filter     = "name eq '${nutanix_service_groups_v2.sg.name}'"
  depends_on = [nutanix_service_groups_v2.sg]
}
`, sgName)
	}
	if stage >= netStageNSP {
		cfg += fmt.Sprintf(`
data "nutanix_categories_v2" "nsp" {}

resource "nutanix_network_security_policy_v2" "nsp" {
  name           = "%[1]s"
  description    = "tenant admin workflow network security policy"
  state          = "SAVE"
  type           = "ISOLATION"
  project_ext_id = local.project_ext_id
  rules {
    type = "TWO_ENV_ISOLATION"
    spec {
      two_env_isolation_rule_spec {
        first_isolation_group = [
          data.nutanix_categories_v2.nsp.categories[0].ext_id,
        ]
        second_isolation_group = [
          data.nutanix_categories_v2.nsp.categories[1].ext_id,
        ]
      }
    }
  }
}

# Single + plural data sources for the network security policy, verified against
# the resource. Note the plural data source exposes the list as network_policies.
data "nutanix_network_security_policy_v2" "nsp" {
  ext_id     = nutanix_network_security_policy_v2.nsp.id
  depends_on = [nutanix_network_security_policy_v2.nsp]
}

data "nutanix_network_security_policies_v2" "nsp" {
  filter     = "name eq '${nutanix_network_security_policy_v2.nsp.name}'"
  depends_on = [nutanix_network_security_policy_v2.nsp]
}
`, nspName)
	}
	return cfg
}

func taNetExternalSubnet(subnetName string, thirdOctet int) string {
	return fmt.Sprintf(`
resource "nutanix_vpc_v2" "transit" {
  name           = "%[1]s-transit"
  description    = "tenant admin transit vpc for overlay external subnet"
  vpc_type       = "TRANSIT"
  project_ext_id = local.project_ext_id
}

resource "nutanix_subnet_v2" "ext" {
  name           = "%[1]s"
  description    = "tenant admin overlay external NAT network"
  vpc_reference  = nutanix_vpc_v2.transit.id
  subnet_type    = "OVERLAY"
  is_external    = true
  is_nat_enabled = true
  project_ext_id = local.project_ext_id
  ip_config {
    ipv4 {
      ip_subnet {
        ip { value = "10.58.%[2]d.0" }
        prefix_length = 24
      }
      default_gateway_ip { value = "10.58.%[2]d.1" }
      pool_list {
        start_ip { value = "10.58.%[2]d.2" }
        end_ip   { value = "10.58.%[2]d.20" }
      }
    }
  }
  depends_on = [nutanix_vpc_v2.transit]
}

# Single + plural data sources for the overlay external subnet, verified against the resource.
data "nutanix_subnet_v2" "ext" {
  ext_id = nutanix_subnet_v2.ext.id
}

data "nutanix_subnets_v2" "ext" {
  filter     = "name eq '${nutanix_subnet_v2.ext.name}'"
  depends_on = [nutanix_subnet_v2.ext]
}
`, subnetName, thirdOctet)
}

// ===========================================================================
// Workflow 4: IAM
//
//	service account -> user API key
//
// Create, update, and delete. The API key is created against the service
// account. The service account description is updated in-place.
// nutanix_user_key_v2 does not support in-place updates (name, description,
// key_type, expiry_time, assigned_to are ForceNew), so the key "update" step
// replaces the key. Delete steps drop resources from the config so Terraform
// destroys them before the test case ends.
// ===========================================================================

func TestAccTenantAdminIAMWorkflow(t *testing.T) {
	r := acctest.RandInt()
	saName := fmt.Sprintf("tf-ta-sa-%d", r)
	saDesc := "tenant admin service account"
	saDescUpdated := saDesc + " update"
	keyName := fmt.Sprintf("tf-ta-sa-key-%d", r)
	expiry := time.Now().Add(14 * 24 * time.Hour).UTC().Format(time.RFC3339)
	assignedTo := "tenant-admin-workflow"
	assignedToUpdated := assignedTo + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccTenantAdminPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: create the service account
				Config: taIAMWorkflowConfig(saName, saDesc, keyName, expiry, assignedTo, iamStageSA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(taServiceAcctRes, "ext_id"),
					resource.TestCheckResourceAttr(taServiceAcctRes, "username", saName),
					resource.TestCheckResourceAttr(taServiceAcctRes, "description", saDesc),
					resource.TestCheckResourceAttr(taServiceAcctRes, "user_type", "SERVICE_ACCOUNT"),
					resource.TestCheckResourceAttr(taServiceAcctRes, "email_id", "terraform_plugin@domain.com"),

					// Service account: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taServiceAcctDS, "ext_id", taServiceAcctRes, "id"),
					resource.TestCheckResourceAttrPair(taServiceAcctDS, "username", taServiceAcctRes, "username"),
					resource.TestCheckResourceAttrPair(taServiceAcctDS, "description", taServiceAcctRes, "description"),
					resource.TestCheckResourceAttrPair(taServiceAcctDS, "user_type", taServiceAcctRes, "user_type"),
					resource.TestCheckResourceAttrPair(taServiceAcctDS, "email_id", taServiceAcctRes, "email_id"),
					// Service account: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taServiceAcctsDS, "users.#", "1"),
					resource.TestCheckResourceAttrPair(taServiceAcctsDS, "users.0.ext_id", taServiceAcctRes, "id"),
					resource.TestCheckResourceAttrPair(taServiceAcctsDS, "users.0.username", taServiceAcctRes, "username"),
					resource.TestCheckResourceAttrPair(taServiceAcctsDS, "users.0.user_type", taServiceAcctRes, "user_type"),
				),
			},
			{
				// Step 2: update the service account description in-place
				Config: taIAMWorkflowConfig(saName, saDescUpdated, keyName, expiry, assignedTo, iamStageSA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(taServiceAcctRes, "ext_id"),
					resource.TestCheckResourceAttr(taServiceAcctRes, "username", saName),
					resource.TestCheckResourceAttr(taServiceAcctRes, "description", saDescUpdated),
					resource.TestCheckResourceAttr(taServiceAcctRes, "user_type", "SERVICE_ACCOUNT"),
					resource.TestCheckResourceAttr(taServiceAcctRes, "email_id", "terraform_plugin@domain.com"),

					resource.TestCheckResourceAttrPair(taServiceAcctDS, "description", taServiceAcctRes, "description"),
					resource.TestCheckResourceAttrPair(taServiceAcctsDS, "users.0.username", taServiceAcctRes, "username"),
				),
			},
			{
				// Step 3: create a user API key on the service account
				Config: taIAMWorkflowConfig(saName, saDescUpdated, keyName, expiry, assignedTo, iamStageKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(taUserKeyRes, "user_ext_id", taServiceAcctRes, "ext_id"),
					resource.TestCheckResourceAttrSet(taUserKeyRes, "id"),
					resource.TestCheckResourceAttrSet(taUserKeyRes, "ext_id"),
					resource.TestCheckResourceAttr(taUserKeyRes, "name", keyName),
					resource.TestCheckResourceAttr(taUserKeyRes, "key_type", "API_KEY"),
					resource.TestCheckResourceAttr(taUserKeyRes, "expiry_time", expiry),
					resource.TestCheckResourceAttr(taUserKeyRes, "assigned_to", assignedTo),
					resource.TestCheckResourceAttr(taUserKeyRes, "status", "VALID"),
					resource.TestCheckResourceAttrSet(taUserKeyRes, "key_details.0.api_key_details.0.api_key"),

					// User key: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taUserKeyDS, "ext_id", taUserKeyRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taUserKeyDS, "name", taUserKeyRes, "name"),
					resource.TestCheckResourceAttrPair(taUserKeyDS, "key_type", taUserKeyRes, "key_type"),
					resource.TestCheckResourceAttrPair(taUserKeyDS, "status", taUserKeyRes, "status"),
					resource.TestCheckResourceAttrPair(taUserKeyDS, "assigned_to", taUserKeyRes, "assigned_to"),
					resource.TestCheckResourceAttrPair(taUserKeyDS, "expiry_time", taUserKeyRes, "expiry_time"),
					// User key: plural data source vs resource (deep).
					resource.TestCheckResourceAttr(taUserKeysDS, "keys.#", "1"),
					resource.TestCheckResourceAttrPair(taUserKeysDS, "keys.0.ext_id", taUserKeyRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taUserKeysDS, "keys.0.name", taUserKeyRes, "name"),
					resource.TestCheckResourceAttrPair(taUserKeysDS, "keys.0.key_type", taUserKeyRes, "key_type"),
					resource.TestCheckResourceAttrPair(taUserKeysDS, "keys.0.status", taUserKeyRes, "status"),
				),
			},
			{
				// Step 4: update the user API key. name/assigned_to are ForceNew,
				// so Terraform replaces the key (no in-place update API).
				Config: taIAMWorkflowConfig(saName, saDescUpdated, keyName, expiry, assignedToUpdated, iamStageKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taUserKeyRes, "name", keyName),
					resource.TestCheckResourceAttr(taUserKeyRes, "key_type", "API_KEY"),
					resource.TestCheckResourceAttr(taUserKeyRes, "assigned_to", assignedToUpdated),
					resource.TestCheckResourceAttr(taUserKeyRes, "status", "VALID"),

					resource.TestCheckResourceAttrPair(taUserKeyDS, "assigned_to", taUserKeyRes, "assigned_to"),
					resource.TestCheckResourceAttrPair(taUserKeysDS, "keys.0.assigned_to", taUserKeyRes, "assigned_to"),
				),
			},
			{
				// Step 5: delete the user API key (drop it from the config;
				// the service account remains).
				Config: taIAMWorkflowConfig(saName, saDescUpdated, keyName, expiry, assignedToUpdated, iamStageSA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taServiceAcctRes, "username", saName),
					resource.TestCheckResourceAttr(taServiceAcctRes, "user_type", "SERVICE_ACCOUNT"),
				),
			},
			{
				// Step 6: delete the service account (drop it from the config).
				Config: taIAMWorkflowConfig(saName, saDescUpdated, keyName, expiry, assignedToUpdated, iamStageNone),
			},
		},
	})
}

// iamStage selects how far along the cumulative IAM workflow a config should
// be rendered. iamStageNone is the delete-everything stage (base data sources
// only).
type iamStage int

const (
	iamStageNone iamStage = iota // no IAM resources
	iamStageSA                   // service account
	iamStageKey                  // + user API key
)

func taIAMWorkflowConfig(saName, saDesc, keyName, expiry, assignedTo string, stage iamStage) string {
	cfg := tenantAdminBaseConfig()
	if stage >= iamStageSA {
		cfg += taWorkflowServiceAccount(saName, saDesc)
	}
	if stage >= iamStageKey {
		cfg += taWorkflowUserKey(keyName, expiry, assignedTo)
	}
	return cfg
}

func taWorkflowServiceAccount(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_users_v2" "service_account" {
  username    = "%[1]s"
  description = "%[2]s"
  email_id    = "terraform_plugin@domain.com"
  user_type   = "SERVICE_ACCOUNT"
}

# Single + plural data sources for the service account, verified against the resource.
data "nutanix_user_v2" "sa" {
  ext_id = nutanix_users_v2.service_account.id
}

data "nutanix_users_v2" "sa" {
  filter     = "username eq '${nutanix_users_v2.service_account.username}'"
  depends_on = [nutanix_users_v2.service_account]
}
`, name, description)
}

func taWorkflowUserKey(keyName, expiry, assignedTo string) string {
	return fmt.Sprintf(`
resource "nutanix_user_key_v2" "key" {
  user_ext_id = nutanix_users_v2.service_account.ext_id
  name        = "%[1]s"
  key_type    = "API_KEY"
  expiry_time = "%[2]s"
  assigned_to = "%[3]s"
}

# Single + plural data sources for the user API key, verified against the resource.
data "nutanix_user_key_v2" "key" {
  user_ext_id = nutanix_users_v2.service_account.ext_id
  ext_id      = nutanix_user_key_v2.key.ext_id
}

data "nutanix_user_keys_v2" "keys" {
  user_ext_id = nutanix_users_v2.service_account.ext_id
  filter      = "name eq '${nutanix_user_key_v2.key.name}'"
  depends_on  = [nutanix_user_key_v2.key]
}
`, keyName, expiry, assignedTo)
}

// ===========================================================================
// Workflow 5: Membership
//
//	create nutanix_directory_services_v2 (AD create) ->
//	Tenant Admin role (nutanix_roles_v2) -> directory search for
//	one user and one user group -> two nutanix_tenant_membership_v2 resources
//
// This is the membership section of terraform/Tenant Admin Workflows/main.tf.
// That example looks up an existing directory service; this workflow still
// creates the AD first, then uses the same membership path: registered
// domain, the first Tenant Admin role, userPrincipalName / distinguishedName
// from directory search, and identity_details on the two memberships.
//
// Fixtures come from iam.directory_services_main.tenant_membership_ad in
// test_config_v2.json. The directory-service name is generated (the IAM API
// rejects dotted names).
// ===========================================================================

func TestAccTenantAdminMembershipWorkflow(t *testing.T) {
	userQuery, groupQuery := taTenantMembershipADQueries(t)
	dsName := fmt.Sprintf("tf-ta-ad-%d", acctest.RandInt())
	ad := testVars.Iam.DirectoryServicesMain.TenantMembershipAD

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccTenantAdminPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixTenantMembershipV2Destroy,
		Steps: []resource.TestStep{
			{
				// Step 1: create the directory service (AD create). Roles and
				// memberships are added in the next step so AD create is
				// asserted on its own.
				Config: taTenantConsumerMembershipConfig(dsName, userQuery, groupQuery, tmStageAD),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(taDirectorySvcRes, "ext_id"),
					resource.TestCheckResourceAttr(taDirectorySvcRes, "name", dsName),
					resource.TestCheckResourceAttr(taDirectorySvcRes, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(taDirectorySvcRes, "domain_name", ad.Name),
					resource.TestCheckResourceAttr(taDirectorySvcRes, "url", ad.URL),
					resource.TestCheckResourceAttr(taDirectorySvcRes, "service_account.0.username", ad.Username),
					resource.TestCheckResourceAttrSet(taDirectorySvcRes, "service_account.0.password"),
					resource.TestCheckResourceAttr(taDirectorySvcRes, "white_listed_groups.0", groupQuery),

					resource.TestCheckResourceAttrPair(taDirectorySvcDS, "ext_id", taDirectorySvcRes, "id"),
					resource.TestCheckResourceAttrPair(taDirectorySvcDS, "name", taDirectorySvcRes, "name"),
					resource.TestCheckResourceAttrPair(taDirectorySvcDS, "directory_type", taDirectorySvcRes, "directory_type"),
					resource.TestCheckResourceAttrPair(taDirectorySvcDS, "domain_name", taDirectorySvcRes, "domain_name"),
					resource.TestCheckResourceAttrPair(taDirectorySvcDS, "url", taDirectorySvcRes, "url"),
					resource.TestCheckResourceAttrPair(taDirectorySvcDS, "service_account.0.username", taDirectorySvcRes, "service_account.0.username"),
					resource.TestCheckResourceAttr(taDirectorySvcsDS, "directory_services.#", "1"),
					resource.TestCheckResourceAttrPair(taDirectorySvcsDS, "directory_services.0.ext_id", taDirectorySvcRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taDirectorySvcsDS, "directory_services.0.name", taDirectorySvcRes, "name"),
					resource.TestCheckResourceAttrPair(taDirectorySvcsDS, "directory_services.0.domain_name", taDirectorySvcRes, "domain_name"),
				),
			},
			{
				// Step 2: search AD for the user and the user group, then
				// create Tenant Admin memberships for both.
				Config: taTenantConsumerMembershipConfig(dsName, userQuery, groupQuery, tmStageMembership),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(taRolesDS, "roles.0.ext_id"),
					resource.TestCheckResourceAttrSet(taRolesDS, "roles.0.display_name"),
					resource.TestCheckResourceAttrSet(taDomainsDS, "registered_domains.0.domain_ext_id"),

					resource.TestCheckResourceAttr(taUserSearchDS, "query", userQuery),
					resource.TestCheckResourceAttrPair(taUserSearchDS, "directory_service_ext_id", taDirectorySvcsDS, "directory_services.0.ext_id"),
					resource.TestCheckResourceAttrSet(taUserSearchDS, "search_results.#"),
					resource.TestCheckResourceAttrSet(taUserSearchDS, "search_results.0.name"),
					resource.TestCheckResourceAttrSet(taUserSearchDS, "search_results.0.entity_type"),
					resource.TestCheckResourceAttrSet(taUserSearchDS, "search_results.0.attributes.#"),

					resource.TestCheckResourceAttr(taGroupSearchDS, "query", groupQuery),
					resource.TestCheckResourceAttrPair(taGroupSearchDS, "directory_service_ext_id", taDirectorySvcsDS, "directory_services.0.ext_id"),
					resource.TestCheckResourceAttrSet(taGroupSearchDS, "search_results.#"),
					resource.TestCheckResourceAttrSet(taGroupSearchDS, "search_results.0.name"),
					resource.TestCheckResourceAttrSet(taGroupSearchDS, "search_results.0.entity_type"),
					resource.TestCheckResourceAttrSet(taGroupSearchDS, "search_results.0.attributes.#"),

					resource.TestCheckResourceAttrSet(taUserMembershipRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taUserMembershipRes, "domain_ext_id", taDomainsDS, "registered_domains.0.domain_ext_id"),
					resource.TestCheckResourceAttrPair(taUserMembershipRes, "role_ext_id", taRolesDS, "roles.0.ext_id"),
					resource.TestCheckResourceAttr(taUserMembershipRes, "identity_details.0.type", "USER"),
					resource.TestCheckResourceAttrSet(taUserMembershipRes, "identity_details.0.principal_value"),
					resource.TestCheckResourceAttrPair(taUserMembershipRes, "identity_details.0.identity_provider_ext_id", taDirectorySvcsDS, "directory_services.0.ext_id"),

					resource.TestCheckResourceAttrSet(taGroupMembershipRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taGroupMembershipRes, "domain_ext_id", taDomainsDS, "registered_domains.0.domain_ext_id"),
					resource.TestCheckResourceAttrPair(taGroupMembershipRes, "role_ext_id", taRolesDS, "roles.0.ext_id"),
					resource.TestCheckResourceAttr(taGroupMembershipRes, "identity_details.0.type", "GROUP"),
					resource.TestCheckResourceAttrSet(taGroupMembershipRes, "identity_details.0.principal_value"),
					resource.TestCheckResourceAttrPair(taGroupMembershipRes, "identity_details.0.identity_provider_ext_id", taDirectorySvcsDS, "directory_services.0.ext_id"),

					resource.TestCheckResourceAttrPair(taUserMembershipDS, "ext_id", taUserMembershipRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taUserMembershipDS, "domain_ext_id", taUserMembershipRes, "domain_ext_id"),
					resource.TestCheckResourceAttrPair(taUserMembershipDS, "role_ext_id", taUserMembershipRes, "role_ext_id"),
					resource.TestCheckResourceAttrPair(taUserMembershipDS, "identity_details.0.principal_value", taUserMembershipRes, "identity_details.0.principal_value"),
					resource.TestCheckResourceAttrPair(taUserMembershipDS, "identity_details.0.type", taUserMembershipRes, "identity_details.0.type"),

					resource.TestCheckResourceAttrPair(taGroupMembershipDS, "ext_id", taGroupMembershipRes, "ext_id"),
					resource.TestCheckResourceAttrPair(taGroupMembershipDS, "domain_ext_id", taGroupMembershipRes, "domain_ext_id"),
					resource.TestCheckResourceAttrPair(taGroupMembershipDS, "role_ext_id", taGroupMembershipRes, "role_ext_id"),
					resource.TestCheckResourceAttrPair(taGroupMembershipDS, "identity_details.0.principal_value", taGroupMembershipRes, "identity_details.0.principal_value"),
					resource.TestCheckResourceAttrPair(taGroupMembershipDS, "identity_details.0.type", taGroupMembershipRes, "identity_details.0.type"),

					checkAttributeLength(taMembershipsDS, "memberships", 2),
				),
			},
		},
	})
}

// tmStage selects how far along the membership workflow a config should be
// rendered.
type tmStage int

const (
	tmStageAD         tmStage = iota + 1 // create directory service
	tmStageMembership                    // + role, searches, user and group memberships
)

func taTenantMembershipADQueries(t *testing.T) (userQuery, groupQuery string) {
	t.Helper()
	ad := testVars.Iam.DirectoryServicesMain.TenantMembershipAD
	if ad.URL == "" || ad.Username == "" {
		t.Skip("skipping membership workflow: iam.directory_services_main.tenant_membership_ad is incomplete")
	}
	userQuery = firstMapKey(ad.UserGroups.Users)
	groupQuery = firstMapKey(ad.UserGroups.UserGroups)
	if userQuery == "" || groupQuery == "" {
		t.Skip("skipping membership workflow: tenant_membership_ad.user_groups is missing a user or user group")
	}
	return userQuery, groupQuery
}

func firstMapKey(m map[string]string) string {
	for k := range m {
		return k
	}
	return ""
}

func taTenantConsumerMembershipConfig(dsName, userQuery, groupQuery string, stage tmStage) string {
	cfg := tenantAdminBaseConfig()
	cfg += taWorkflowDirectoryService(dsName, groupQuery)
	if stage >= tmStageMembership {
		cfg += taWorkflowTenantConsumerMemberships(userQuery, groupQuery)
	}
	return cfg
}

func taWorkflowDirectoryService(dsName, whitelistGroup string) string {
	return fmt.Sprintf(`
locals {
  tm_ad = local.config.iam.directory_services_main.tenant_membership_ad
}

resource "nutanix_directory_services_v2" "ad" {
  name           = "%[1]s"
  url            = local.tm_ad.url
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = local.tm_ad.name
  service_account {
    username = local.tm_ad.username
    password = local.tm_ad.password
  }
  white_listed_groups = ["%[2]s"]
  lifecycle {
    ignore_changes = [
      service_account.0.password,
    ]
  }
}

data "nutanix_directory_service_v2" "ad" {
  ext_id = nutanix_directory_services_v2.ad.id
}

data "nutanix_directory_services_v2" "ad" {
  filter     = "name eq '${nutanix_directory_services_v2.ad.name}'"
  depends_on = [nutanix_directory_services_v2.ad]
}
`, dsName, whitelistGroup)
}

func taWorkflowTenantConsumerMemberships(userQuery, groupQuery string) string {
	return fmt.Sprintf(`
data "nutanix_registered_domains_v2" "domains" {}

data "nutanix_roles_v2" "roles" {
  filter = "contains(displayName, 'Tenant Admin')"
}

data "nutanix_directory_service_users_search_v2" "user" {
  directory_service_ext_id = data.nutanix_directory_services_v2.ad.directory_services[0].ext_id
  query                    = "%[1]s"
  is_wildcard_search       = false
}

data "nutanix_directory_service_users_search_v2" "group" {
  directory_service_ext_id = data.nutanix_directory_services_v2.ad.directory_services[0].ext_id
  query                    = "%[2]s"
  is_wildcard_search       = false
}

locals {
  domain_ext_id            = data.nutanix_registered_domains_v2.domains.registered_domains[0].domain_ext_id
  tenant_admin_role_ext_id = data.nutanix_roles_v2.roles.roles[0].ext_id

  search_results_user = data.nutanix_directory_service_users_search_v2.user.search_results
  user_principal_name = one(flatten([
    for result in local.search_results_user : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "userPrincipalName"
    ]
  ]))
  user_identity_ext_id = one(flatten([
    for result in local.search_results_user : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "objectGUID"
    ]
  ]))

  search_results_group = data.nutanix_directory_service_users_search_v2.group.search_results
  group_principal_name = one(flatten([
    for result in local.search_results_group : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "distinguishedName"
    ]
  ]))
  group_identity_ext_id = one(flatten([
    for result in local.search_results_group : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "objectGUID"
    ]
  ]))
}

resource "nutanix_tenant_membership_v2" "user" {
  domain_ext_id = local.domain_ext_id
  role_ext_id   = local.tenant_admin_role_ext_id

  identity_details {
    identity_provider_ext_id = data.nutanix_directory_services_v2.ad.directory_services[0].ext_id
    principal_value          = local.user_principal_name
    type                     = "USER"
  }
}

resource "nutanix_tenant_membership_v2" "group" {
  domain_ext_id = local.domain_ext_id
  role_ext_id   = local.tenant_admin_role_ext_id

  identity_details {
    identity_provider_ext_id = data.nutanix_directory_services_v2.ad.directory_services[0].ext_id
    principal_value          = local.group_principal_name
    type                     = "GROUP"
  }
}

data "nutanix_tenant_membership_v2" "user" {
  ext_id = nutanix_tenant_membership_v2.user.ext_id
}

data "nutanix_tenant_membership_v2" "group" {
  ext_id = nutanix_tenant_membership_v2.group.ext_id
}

data "nutanix_tenant_memberships_v2" "memberships" {
  filter = "domainExtId eq '${local.domain_ext_id}'"
  depends_on = [
    nutanix_tenant_membership_v2.user,
    nutanix_tenant_membership_v2.group,
  ]
}
`, userQuery, groupQuery)
}

// ===========================================================================
// Workflow 6: Volume Group
//
//	create nutanix_volume_group_v2 -> update (name + description in-place) ->
//	delete
//
// The volume group is hosted on the cluster resolved from the project's first
// resource-group placement target (local.cluster_ext_id) and is associated to
// the FIRST project entity (local.project_ext_id), consistent with every other
// tenant admin resource. Each step verifies the single + plural data sources
// field-by-field against the resource. The delete step drops the resource from
// the config so Terraform destroys it before the test case ends.
// ===========================================================================

func TestAccTenantAdminVolumeGroupWorkflow(t *testing.T) {
	r := acctest.RandInt()
	vgName := fmt.Sprintf("tf-ta-vg-%d", r)
	vgDesc := "tenant admin volume group"
	vgNameUpdated := fmt.Sprintf("tf-ta-vg-updated-%d", r)
	vgDescUpdated := vgDesc + " update"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccTenantAdminPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: create the volume group
				Config: taVolumeGroupConfig(vgName, vgDesc, vgStageVG),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(taVolumeGroupRes, "ext_id"),
					resource.TestCheckResourceAttr(taVolumeGroupRes, "name", vgName),
					resource.TestCheckResourceAttr(taVolumeGroupRes, "description", vgDesc),
					resource.TestCheckResourceAttr(taVolumeGroupRes, "sharing_status", "SHARED"),
					resource.TestCheckResourceAttr(taVolumeGroupRes, "usage_type", "USER"),
					resource.TestCheckResourceAttrPair(taVolumeGroupRes, "project_ext_id", projectDataSource, "projects.0.ext_id"),

					// Volume group: single data source vs resource (deep).
					resource.TestCheckResourceAttrPair(taVolumeGroupDS, "ext_id", taVolumeGroupRes, "id"),
					resource.TestCheckResourceAttrPair(taVolumeGroupDS, "name", taVolumeGroupRes, "name"),
					resource.TestCheckResourceAttrPair(taVolumeGroupDS, "description", taVolumeGroupRes, "description"),
					resource.TestCheckResourceAttrPair(taVolumeGroupDS, "cluster_reference", taVolumeGroupRes, "cluster_reference"),
					resource.TestCheckResourceAttrPair(taVolumeGroupDS, "sharing_status", taVolumeGroupRes, "sharing_status"),
					resource.TestCheckResourceAttrPair(taVolumeGroupDS, "usage_type", taVolumeGroupRes, "usage_type"),
					// Volume group: plural data source vs resource (deep). List attr is volumes.
					resource.TestCheckResourceAttr(taVolumeGroupsDS, "volumes.#", "1"),
					resource.TestCheckResourceAttrPair(taVolumeGroupsDS, "volumes.0.ext_id", taVolumeGroupRes, "id"),
					resource.TestCheckResourceAttrPair(taVolumeGroupsDS, "volumes.0.name", taVolumeGroupRes, "name"),
					resource.TestCheckResourceAttrPair(taVolumeGroupsDS, "volumes.0.cluster_reference", taVolumeGroupRes, "cluster_reference"),
				),
			},
			{
				// Step 2: update the name + description in-place
				Config: taVolumeGroupConfig(vgNameUpdated, vgDescUpdated, vgStageVG),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(taVolumeGroupRes, "ext_id"),
					resource.TestCheckResourceAttr(taVolumeGroupRes, "name", vgNameUpdated),
					resource.TestCheckResourceAttr(taVolumeGroupRes, "description", vgDescUpdated),

					// Data sources reflect the updated name + description.
					resource.TestCheckResourceAttrPair(taVolumeGroupDS, "name", taVolumeGroupRes, "name"),
					resource.TestCheckResourceAttrPair(taVolumeGroupDS, "description", taVolumeGroupRes, "description"),
					resource.TestCheckResourceAttrPair(taVolumeGroupsDS, "volumes.0.name", taVolumeGroupRes, "name"),
				),
			},
			{
				// Step 3: delete the volume group (drop it from the config).
				Config: taVolumeGroupConfig(vgNameUpdated, vgDescUpdated, vgStageNone),
			},
		},
	})
}

// vgStage selects how far along the volume group workflow a config should be
// rendered. vgStageNone is the delete stage (base data sources only).
type vgStage int

const (
	vgStageNone vgStage = iota // no volume group
	vgStageVG                  // volume group
)

func taVolumeGroupConfig(vgName, vgDesc string, stage vgStage) string {
	cfg := tenantAdminBaseConfig()
	if stage >= vgStageVG {
		cfg += taWorkflowVolumeGroup(vgName, vgDesc)
	}
	return cfg
}

func taWorkflowVolumeGroup(vgName, vgDesc string) string {
	return fmt.Sprintf(`
resource "nutanix_volume_group_v2" "vg" {
  name                               = "%[1]s"
  description                        = "%[2]s"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  cluster_reference                  = local.cluster_ext_id
  project_ext_id                     = local.project_ext_id
  storage_features {
    flash_mode {
      is_enabled = true
    }
  }
  usage_type = "USER"
  is_hidden  = false

  # target_secret is never returned by the API, so ignore it on refresh.
  lifecycle {
    ignore_changes = [
      iscsi_features[0].target_secret,
    ]
  }
}

# Single + plural data sources for the volume group, verified against the resource.
data "nutanix_volume_group_v2" "vg" {
  ext_id = nutanix_volume_group_v2.vg.id
}

data "nutanix_volume_groups_v2" "vg" {
  filter     = "name eq '${nutanix_volume_group_v2.vg.name}'"
  depends_on = [nutanix_volume_group_v2.vg]
}
`, vgName, vgDesc)
}
