package clusters_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixHostDataSource_basic(t *testing.T) {
	dataSourceName := "data.nutanix_host.test"
	vmResourceName := "nutanix_virtual_machine.vm"

	imgName := fmt.Sprintf("test-acc-dou-image-%s", acctest.RandString(3))
	vmName := fmt.Sprintf("test-acc-dou-vm-%s", acctest.RandString(3))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHostDataSourceConfig(imgName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vmResourceName, "name"),
					resource.TestCheckResourceAttrSet(vmResourceName, "cluster_uuid"),
					resource.TestCheckResourceAttrSet(vmResourceName, "num_vcpus_per_socket"),
					resource.TestCheckResourceAttrSet(vmResourceName, "num_sockets"),
					resource.TestCheckResourceAttrSet(vmResourceName, "memory_size_mib"),
					// resource.TestCheckResourceAttrSet(vmResourceName, "serial_port_list.#"),
					resource.TestCheckResourceAttrSet(vmResourceName, "disk_list.#"),

					resource.TestCheckResourceAttr(vmResourceName, "name", vmName),
					resource.TestCheckResourceAttr(vmResourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(vmResourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(vmResourceName, "memory_size_mib", "186"),
					// This check is commented out because the serial port index not returned by the API Response.
					// resource.TestCheckResourceAttr(vmResourceName, "serial_port_list.0.index", "1"),
					// resource.TestCheckResourceAttr(vmResourceName, "serial_port_list.0.is_connected", "true"),
					resource.TestCheckResourceAttr(vmResourceName, "disk_list.#", "4"),

					resource.TestCheckResourceAttrSet(dataSourceName, "host_id"),
				),
			},
		},
	})
}

func testAccHostDataSourceConfig(imgName, vmName string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_image" "cirros-034-disk" {
			name        = "%s"
			source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
			description = "heres a tiny linux image, not an iso, but a real disk!"
		}

		resource "nutanix_virtual_machine" "vm" {
			name         = "%s"
			cluster_uuid = "${local.cluster1}"

			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			#serial_port_list {
			#	index = 1
			#	is_connected = true
			#}

			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = "${nutanix_image.cirros-034-disk.id}"
				}

				device_properties {
					disk_address = {
						device_index = 0,
						adapter_type = "IDE"
					}
					device_type = "CDROM"
				}
			}
			disk_list {
				disk_size_mib = 100
			}
			disk_list {
				disk_size_mib = 200
			}
			disk_list {
				disk_size_mib = 300
			}
		}

		data "nutanix_host" "test" {
			host_id = nutanix_virtual_machine.vm.host_reference.uuid
		}
	`, imgName, vmName)
}
