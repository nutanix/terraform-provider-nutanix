package storagecontainersv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameStorageContainers = "nutanix_storage_containers_v2.test"

// Hardcoded storage-container fixtures. Shared across the
// storagecontainersv2 resource/data-source tests.
const (
	storageContainerName                     = "terraform_storage_container_test"
	storageContainerAdvertisedCapacityBytes  = "1073741824000"
	storageContainerReservedCapacityBytes    = "20"
	storageContainerReplicationFactor        = "1"
	storageContainerNfsWhitelistIPv4         = "192.168.14.0"
	storageContainerNfsWhitelistIPv4Updated  = "192.168.15.0"
	storageContainerNfsWhitelistPrefixLength = "32"
)

func TestAccV2NutanixStorageContainersResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-storage-container-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainersResourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStorageContainers, "container_ext_id"),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "name", name),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "logical_advertised_capacity_bytes", storageContainerAdvertisedCapacityBytes),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "logical_explicit_reserved_capacity_bytes", storageContainerReservedCapacityBytes),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "replication_factor", storageContainerReplicationFactor),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "nfs_whitelist_addresses.0.ipv4.0.value", storageContainerNfsWhitelistIPv4),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "nfs_whitelist_addresses.0.ipv4.0.prefix_length", storageContainerNfsWhitelistPrefixLength),
				),
			},
			// test update
			{
				Config: testStorageContainersResourceUpdateConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStorageContainers, "container_ext_id"),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "name", fmt.Sprintf("%s_updated", name)),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "logical_advertised_capacity_bytes", storageContainerAdvertisedCapacityBytes),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "logical_explicit_reserved_capacity_bytes", storageContainerReservedCapacityBytes),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "replication_factor", storageContainerReplicationFactor),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "nfs_whitelist_addresses.0.ipv4.0.value", storageContainerNfsWhitelistIPv4Updated),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "nfs_whitelist_addresses.0.ipv4.0.prefix_length", storageContainerNfsWhitelistPrefixLength),
				),
			},
		},
	})
}

func TestAccV2NutanixStorageContainersResource_WithNoClusterExtId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStorageContainersResourceWithoutClusterExtIDConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixStorageContainersResource_WithNoName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStorageContainersResourceWithoutNameConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testStorageContainersResourceConfig(name string) string {
	return fmt.Sprintf(`

		data "nutanix_clusters_v2" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_storage_containers_v2" "test" {
			name = "%[1]s"
			cluster_ext_id = local.cluster
			logical_advertised_capacity_bytes = 1073741824000
			logical_explicit_reserved_capacity_bytes = 20
			replication_factor = 1
			nfs_whitelist_addresses {
				ipv4  {
					value = "192.168.14.0"
					prefix_length = 32
				}
			}
			erasure_code = "OFF"
			is_inline_ec_enabled = false
			has_higher_ec_fault_domain_preference = false
			cache_deduplication = "OFF"
			on_disk_dedup = "OFF"
			is_compression_enabled = true
			is_internal = false
			is_software_encryption_enabled = false
		}`, name)
}

func testStorageContainersResourceUpdateConfig(name string) string {
	return fmt.Sprintf(`

		data "nutanix_clusters_v2" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_storage_containers_v2" "test" {
			name = "%[1]s_updated"
			cluster_ext_id = local.cluster
			logical_advertised_capacity_bytes = 1073741824000
			logical_explicit_reserved_capacity_bytes = 20
			replication_factor = 1
			nfs_whitelist_addresses {
				ipv4  {
					value = "192.168.15.0"
					prefix_length = 32
				}
			}
			erasure_code = "OFF"
			is_inline_ec_enabled = false
			has_higher_ec_fault_domain_preference = false
			cache_deduplication = "OFF"
			on_disk_dedup = "OFF"
			is_compression_enabled = true
			is_internal = false
			is_software_encryption_enabled = false
		}`, name)
}

func testStorageContainersResourceWithoutNameConfig() string {
	return `

		data "nutanix_clusters_v2" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_storage_containers_v2" "test" {
			cluster_ext_id = local.cluster
			logical_advertised_capacity_bytes = 1073741824000
			logical_explicit_reserved_capacity_bytes = 20
			replication_factor = 1
			nfs_whitelist_addresses {
				ipv4  {
					value = "192.168.14.0"
					prefix_length = 32
				}
			}
			erasure_code = "OFF"
			is_inline_ec_enabled = false
			has_higher_ec_fault_domain_preference = false
			cache_deduplication = "OFF"
			on_disk_dedup = "OFF"
			is_compression_enabled = true
			is_internal = false
			is_software_encryption_enabled = false
		}`
}

func testStorageContainersResourceWithoutClusterExtIDConfig() string {
	return fmt.Sprintf(`

		data "nutanix_clusters_v2" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_storage_containers_v2" "test" {
			name = "%[1]s"
			logical_advertised_capacity_bytes = 1073741824000
			logical_explicit_reserved_capacity_bytes = 20
			replication_factor = 1
			nfs_whitelist_addresses {
				ipv4  {
					value = "192.168.14.0"
					prefix_length = 32
				}
			}
			erasure_code = "OFF"
			is_inline_ec_enabled = false
			has_higher_ec_fault_domain_preference = false
			cache_deduplication = "OFF"
			on_disk_dedup = "OFF"
			is_compression_enabled = true
			is_internal = false
			is_software_encryption_enabled = false
		}`, storageContainerName)
}
