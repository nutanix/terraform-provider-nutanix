package storagecontainersv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameStorageContainer = "data.nutanix_storage_container_v2.test"

func TestAccV2NutanixStorageContainerDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-container-%d", r)
	path, _ := os.Getwd()
	filepath := path + "/../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainerV4Config(filepath, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainer, "container_ext_id"),
					resource.TestCheckResourceAttr(datasourceNameStorageContainer, "name", name),
					resource.TestCheckResourceAttr(datasourceNameStorageContainer, "logical_advertised_capacity_bytes", strconv.Itoa(testVars.StorageContainer.LogicalAdvertisedCapacityBytes)),
					resource.TestCheckResourceAttr(datasourceNameStorageContainer, "logical_explicit_reserved_capacity_bytes", strconv.Itoa(testVars.StorageContainer.LogicalExplicitReservedCapacityBytes)),
					resource.TestCheckResourceAttr(datasourceNameStorageContainer, "replication_factor", strconv.Itoa(testVars.StorageContainer.ReplicationFactor)),
					resource.TestCheckResourceAttr(datasourceNameStorageContainer, "nfs_whitelist_addresses.0.ipv4.0.value", testVars.StorageContainer.NfsWhitelistAddresses.Ipv4.Value),
					resource.TestCheckResourceAttr(datasourceNameStorageContainer, "nfs_whitelist_addresses.0.ipv4.0.prefix_length", strconv.Itoa(testVars.StorageContainer.NfsWhitelistAddresses.Ipv4.PrefixLength)),
				),
			},
		},
	})
}

func testStorageContainerV4Config(filepath, name string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = (jsondecode(file("%[1]s")))
			storage_container = local.config.storage_container			
		}

		resource "nutanix_storage_containers_v2" "test" {
			name = "%[2]s"
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
			
		data "nutanix_storage_container_v2" "test" {
			ext_id = resource.nutanix_storage_containers_v2.test.id
		}

		
	`, filepath, name)
}
