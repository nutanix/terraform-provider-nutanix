package storagecontainersv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameStorageContainersV4 = "data.nutanix_storage_containers_v2.test"

func TestAccNutanixStorageContainersV2Datasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainersV4DatasourceV4Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainersV4, "storage_containers.#"),
				),
			},
		},
	})
}

func TestAccNutanixStorageContainersV2Datasource_WithFilter(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainersV4DatasourceV4WithFilterConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainersV4, "storage_containers.#"),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainersV4, "storage_containers.0.container_ext_id"),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.name", testVars.StorageContainer.Name),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.logical_advertised_capacity_bytes", strconv.Itoa(testVars.StorageContainer.LogicalAdvertisedCapacityBytes)),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.logical_explicit_reserved_capacity_bytes", strconv.Itoa(testVars.StorageContainer.LogicalExplicitReservedCapacityBytes)),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.replication_factor", strconv.Itoa(testVars.StorageContainer.ReplicationFactor)),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.nfs_whitelist_addresses.0.ipv4.0.value", testVars.StorageContainer.NfsWhitelistAddresses.Ipv4.Value),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.nfs_whitelist_addresses.0.ipv4.0.prefix_length", strconv.Itoa(testVars.StorageContainer.NfsWhitelistAddresses.Ipv4.PrefixLength)),
				),
			},
		},
	})
}

func TestAccNutanixStorageContainersV2Datasource_WithLimit(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainersV4DatasourceV4WithLimitConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainersV4, "storage_containers.#"),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.#", "1"),
				),
			},
		},
	})
}

func testStorageContainersV4DatasourceV4Config() string {
	return `
	data "nutanix_storage_containers_v2" "test"{}
	`
}

func testStorageContainersV4DatasourceV4WithFilterConfig(filepath string) string {
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
		}

		data "nutanix_storage_containers_v2" "test" {
			filter = "name eq '${local.storage_container.name}'"
			depends_on = [nutanix_storage_containers_v2.test]
		}

	`, filepath)
}

func testStorageContainersV4DatasourceV4WithLimitConfig(filepath string) string {
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
		}

		data "nutanix_storage_containers_v2" "test" {
			limit     = 1
			depends_on = [nutanix_storage_containers_v2.test]
		}
	`, filepath)
}
