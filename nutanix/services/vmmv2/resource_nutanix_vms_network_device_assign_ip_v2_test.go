package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVMNetworkDeviceAssignIP = "nutanix_vm_network_device_assign_ip_v2.test"

func TestAccV2NutanixVmsNetworkDeviceAssignIpResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("test-vm-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMPreEnvConfig(r) + testVMWithNicAndDiskConfig(vmName) + testVmsNetworkDeviceAssignIPV4Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVMNetworkDeviceAssignIP, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVMNetworkDeviceAssignIP, "ip_address.0.value", testVars.VMM.AssignedIP),
				),
			},
		},
	})
}

func testVmsNetworkDeviceAssignIPV4Config() string {
	return `
		resource "nutanix_vm_network_device_assign_ip_v2" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test-vm.id
			ext_id    = resource.nutanix_virtual_machine_v2.test-vm.nics.0.ext_id
			ip_address {
			  value = local.vmm.assigned_ip
			}
		}
`
}

func testVMPreEnvConfig(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}
		
		locals {
		  clusterUUID = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		  config = (jsondecode(file("%[1]s")))
		  vmm    = local.config.vmm
 		  gs = base64encode("#cloud-config\nusers:\n  - name: ubuntu\n    ssh-authorized-keys:\n      - ssh-rsa DUMMYSSH mypass\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']")
		}

	
		
		data "nutanix_storage_containers_v2" "sc" {
		  filter = "clusterExtId eq '${local.clusterUUID}'"
		  limit = 1
		}		

		resource "nutanix_subnet_v2" "subnet" {
			name = "tf-test-subnet-%[2]d"
			description = "terraform test subnet to assign ip"
			cluster_reference = local.clusterUUID
			subnet_type = "VLAN"
			network_id = local.vmm.subnet.network_id
			is_external = false
			ip_config {
				ipv4 {
					ip_subnet {
						ip {
							value = local.vmm.subnet.ip
						}
						prefix_length = local.vmm.subnet.prefix
					}
					default_gateway_ip {
						value = local.vmm.subnet.gateway_ip
					}
					pool_list{
						start_ip {
							value = local.vmm.subnet.start_ip
						}
						end_ip {
							value = local.vmm.subnet.end_ip
						}
					}
				}
			}
		depends_on = [data.nutanix_clusters_v2.clusters]
		}

			
`, filepath, r)
}

func testVMWithNicAndDiskConfig(vmName string) string {
	return fmt.Sprintf(`
		resource "nutanix_virtual_machine_v2" "test-vm" {
		  name                 = "%[1]s"
		  description          = "vm for testing "
		  num_cores_per_socket = 1
		  num_sockets          = 1
		  memory_size_bytes    = 4 * 1024 * 1024 * 1024
		  cluster {
			ext_id = local.clusterUUID
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
				  ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
				}
			  }
			}
		  }
       
          nics {
			network_info {
			  nic_type = "DIRECT_NIC"
			  subnet {
				ext_id = nutanix_subnet_v2.subnet.ext_id
			  }
			  vlan_mode = "ACCESS"
			}
		  }	  
				  
		  power_state = "OFF"

		  lifecycle {
			ignore_changes = [guest_tools, nics]
		  }
		
		  depends_on = [data.nutanix_clusters_v2.clusters, data.nutanix_storage_containers_v2.sc]
		}
`, vmName)
}
