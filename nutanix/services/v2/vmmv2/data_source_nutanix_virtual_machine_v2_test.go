package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVms = "data.nutanix_virtual_machine_v2.test"

func TestAccNutanixVmsDataSourceV2_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVms, "name", name),
					resource.TestCheckResourceAttr(datasourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(datasourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(datasourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(datasourceNameVms, "machine_type", "PC"),
				),
			},
		},
	})
}

func TestAccNutanixVmsDataSourceV2_WithConfig(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4WithNic(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVms, "name", fmt.Sprintf("test-vm-%d", r)),
					resource.TestCheckResourceAttr(datasourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(datasourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(datasourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(datasourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(datasourceNameVms, "nics.#"),
					resource.TestCheckResourceAttr(datasourceNameVms, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(datasourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
		},
	})
}

func testAccVMDataSourceConfigV4(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters.clusters.entities :
			cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_virtual_machine_v2" "vm1" {
			name = "%[1]s"
			description = "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}

		data "nutanix_virtual_machine_v2" "test" {
			ext_id = nutanix_virtual_machine_v2.vm1.id
		}
`, name, desc)
}

func testAccVMDataSourceConfigV4WithNic(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters.clusters.entities :
			cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
			config = (jsondecode(file("%[3]s")))
		  	vmm    = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		resource "nutanix_virtual_machine_v2" "vm1"{
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
		}

		data "nutanix_virtual_machine_v2" "test" {
			ext_id = nutanix_virtual_machine_v2.vm1.id
		}
`, name, desc, filepath)
}
