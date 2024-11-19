package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmGpu = "data.nutanix_vm_gpu_v4.test"

func TestAccNutanixVmsGpuDataSourceV4_List(t *testing.T) {
	t.Skip("Skipping test as it requires a VM with GPU attached")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceVmGpuV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmGpu, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpu, "mode"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpu, "vendor"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpu, "name"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpu, "pci_address.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpu, "frame_buffer_size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpu, "num_virtual_display_heads"),
				),
			},
		},
	})
}

func testVmsDatasourceVmGpuV4Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}

		data "nutanix_subnets_v2" "subnets" { }
	
		resource "nutanix_virtual_machine_v2" "rtest"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
					}	
					vlan_mode = "ACCESS"
				}
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					disk_size_bytes = 21474836480
					data_source {
						reference{
							image_reference{
								image_ext_id = "5867f64e-7d0a-4b04-a72e-e26a4dbbaea2"
							}
						}
					}
				}
			}
		}
		resource "nutanix_vm_gpus_v2" "rtest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			mode = "PASSTHROUGH_COMPUTE"
			vendor= "NVIDIA"
		}

		data "nutanix_vm_gpu_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			ext_id = resource.nutanix_vm_gpus_v2.rtest.id
		}
`, name, desc)
}
