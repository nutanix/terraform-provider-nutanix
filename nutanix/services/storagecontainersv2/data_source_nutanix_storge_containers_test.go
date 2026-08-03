package storagecontainersv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameStorageContainersV4 = "data.nutanix_storage_containers_v2.test"

func TestAccV2NutanixStorageContainersDataSource_Basic(t *testing.T) {
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

func TestAccV2NutanixStorageContainersDataSource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-container-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainersV4DatasourceV4WithFilterConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainersV4, "storage_containers.#"),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainersV4, "storage_containers.0.container_ext_id"),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.logical_advertised_capacity_bytes", storageContainerAdvertisedCapacityBytes),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.logical_explicit_reserved_capacity_bytes", storageContainerReservedCapacityBytes),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.replication_factor", storageContainerReplicationFactor),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.nfs_whitelist_addresses.0.ipv4.0.value", storageContainerNfsWhitelistIPv4),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.0.nfs_whitelist_addresses.0.ipv4.0.prefix_length", storageContainerNfsWhitelistPrefixLength),
				),
			},
		},
	})
}

func TestAccV2NutanixStorageContainersDataSource_WithLimit(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-container-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainersV4DatasourceV4WithLimitConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainersV4, "storage_containers.#"),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixStorageContainersDataSource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainersV4DatasourceV4WithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageContainersV4, "storage_containers.#"),
					resource.TestCheckResourceAttr(datasourceNameStorageContainersV4, "storage_containers.#", "0"),
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

func testStorageContainersV4DatasourceV4WithFilterConfig(name string) string {
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
		}

		data "nutanix_storage_containers_v2" "test" {
			filter = "name eq '${nutanix_storage_containers_v2.test.name}'"
			depends_on = [nutanix_storage_containers_v2.test]
		}

	`, name)
}

func testStorageContainersV4DatasourceV4WithLimitConfig(name string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals{
			cluster =[
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
		}

		data "nutanix_storage_containers_v2" "test" {
			limit     = 1
			depends_on = [nutanix_storage_containers_v2.test]
		}
	`, name)
}

func testStorageContainersV4DatasourceV4WithInvalidFilterConfig() string {
	return `
	data "nutanix_storage_containers_v2" "test" {
		filter = "name eq 'invalid-name'"
	}
	`
}
