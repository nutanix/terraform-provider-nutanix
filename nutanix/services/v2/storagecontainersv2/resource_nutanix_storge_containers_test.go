package storagecontainersv2_test

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameStorageContainers = "nutanix_storage_containers_v2.test"

func TestAccNutanixStorageContainersV2Resource_Basic(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainersResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStorageContainers, "container_ext_id"),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "name", testVars.StorageContainer.Name),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "logical_advertised_capacity_bytes", strconv.Itoa(testVars.StorageContainer.LogicalAdvertisedCapacityBytes)),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "logical_explicit_reserved_capacity_bytes", strconv.Itoa(testVars.StorageContainer.LogicalExplicitReservedCapacityBytes)),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "replication_factor", strconv.Itoa(testVars.StorageContainer.ReplicationFactor)),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "nfs_whitelist_addresses.0.ipv4.0.value", testVars.StorageContainer.NfsWhitelistAddresses.Ipv4.Value),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "nfs_whitelist_addresses.0.ipv4.0.prefix_length", strconv.Itoa(testVars.StorageContainer.NfsWhitelistAddresses.Ipv4.PrefixLength)),
				),
			},
			// test update
			{
				Config: testStorageContainersResourceUpdateConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStorageContainers, "container_ext_id"),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "name", fmt.Sprintf("%s_updated", testVars.StorageContainer.Name)),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "logical_advertised_capacity_bytes", strconv.Itoa(testVars.StorageContainer.LogicalAdvertisedCapacityBytes)),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "logical_explicit_reserved_capacity_bytes", strconv.Itoa(testVars.StorageContainer.LogicalExplicitReservedCapacityBytes)),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "replication_factor", strconv.Itoa(testVars.StorageContainer.ReplicationFactor)),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "nfs_whitelist_addresses.0.ipv4.0.value", "192.168.15.0"),
					resource.TestCheckResourceAttr(resourceNameStorageContainers, "nfs_whitelist_addresses.0.ipv4.0.prefix_length", strconv.Itoa(testVars.StorageContainer.NfsWhitelistAddresses.Ipv4.PrefixLength)),
				),
			},
		},
	})
}

func TestAccNutanixStorageContainersV2Resource_WithNoClusterExtId(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStorageContainersResourceWithoutClusterExtIdConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}
func TestAccNutanixStorageContainersV2Resource_WithNoName(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStorageContainersResourceWithoutNameConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testStorageContainersResourceConfig(filepath string) string {
	return fmt.Sprintf(`

		data "nutanix_clusters" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
				][0]
			config = (jsondecode(file("%s")))
			storage_container = local.config.storage_container			
		}

		resource "nutanix_storage_containers_v2" "test" {
			name = local.storage_container.name
			cluster_ext_id = local.cluster
			logical_advertised_capacity_bytes = local.storage_container.logical_advertised_capacity_bytes
			logical_explicit_reserved_capacity_bytes = local.storage_container.logical_explicit_reserved_capacity_bytes
			replication_factor = local.storage_container.replication_factor
			nfs_whitelist_addresses {
				ipv4  {
					value = local.storage_container.nfs_whitelist_addresses.ipv4.value
					prefix_length = local.storage_container.nfs_whitelist_addresses.ipv4.prefix_length
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
		}`, filepath)
}

func testStorageContainersResourceUpdateConfig(filepath string) string {
	return fmt.Sprintf(`

		data "nutanix_clusters" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
				][0]
			config = (jsondecode(file("%s")))
			storage_container = local.config.storage_container			
		}

		resource "nutanix_storage_containers_v2" "test" {
			name = "${local.storage_container.name}_updated"
			cluster_ext_id = local.cluster
			logical_advertised_capacity_bytes = local.storage_container.logical_advertised_capacity_bytes
			logical_explicit_reserved_capacity_bytes = local.storage_container.logical_explicit_reserved_capacity_bytes
			replication_factor = local.storage_container.replication_factor
			nfs_whitelist_addresses {
				ipv4  {
					value = "192.168.15.0"
					prefix_length = local.storage_container.nfs_whitelist_addresses.ipv4.prefix_length
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
		}`, filepath)
}

func testStorageContainersResourceWithoutNameConfig(filepath string) string {
	return fmt.Sprintf(`

		data "nutanix_clusters" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
				][0]
			config = (jsondecode(file("%s")))
			storage_container = local.config.storage_container			
		}

		resource "nutanix_storage_containers_v2" "test" {
			cluster_ext_id = local.cluster
			logical_advertised_capacity_bytes = local.storage_container.logical_advertised_capacity_bytes
			logical_explicit_reserved_capacity_bytes = local.storage_container.logical_explicit_reserved_capacity_bytes
			replication_factor = local.storage_container.replication_factor
			nfs_whitelist_addresses {
				ipv4  {
					value = local.storage_container.nfs_whitelist_addresses.ipv4.value
					prefix_length = local.storage_container.nfs_whitelist_addresses.ipv4.prefix_length
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
		}`, filepath)
}

func testStorageContainersResourceWithoutClusterExtIdConfig(filepath string) string {
	return fmt.Sprintf(`

		data "nutanix_clusters" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
				][0]
			config = (jsondecode(file("%s")))
			storage_container = local.config.storage_container			
		}

		resource "nutanix_storage_containers_v2" "test" {
			name = local.storage_container.name
			logical_advertised_capacity_bytes = local.storage_container.logical_advertised_capacity_bytes
			logical_explicit_reserved_capacity_bytes = local.storage_container.logical_explicit_reserved_capacity_bytes
			replication_factor = local.storage_container.replication_factor
			nfs_whitelist_addresses {
				ipv4  {
					value = local.storage_container.nfs_whitelist_addresses.ipv4.value
					prefix_length = local.storage_container.nfs_whitelist_addresses.ipv4.prefix_length
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
		}`, filepath)
}
