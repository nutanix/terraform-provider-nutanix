package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVMs = "data.nutanix_virtual_machine_v2.test"

func TestAccV2NutanixVmsDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMs, "name", name),
					resource.TestCheckResourceAttr(datasourceNameVMs, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "description", desc),
					resource.TestCheckResourceAttr(datasourceNameVMs, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "machine_type", "PC"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsDatasource_WithConfig(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4WithNic(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMs, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(datasourceNameVMs, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "description", desc),
					resource.TestCheckResourceAttr(datasourceNameVMs, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.#"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "nics.0.network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.0.network_info.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsDatasource_WithCdromConfig(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4WithCdrom(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMs, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(datasourceNameVMs, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "description", desc),
					resource.TestCheckResourceAttr(datasourceNameVMs, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.#"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "nics.0.network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.0.network_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "cd_roms.#"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "cd_roms.0.disk_address.0.bus_type", "SATA"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "cd_roms.0.disk_address.0.index", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsDatasource_WithCdromBackingInfo(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4WithCdromBackingInfo(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMs, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(datasourceNameVMs, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "description", desc),
					resource.TestCheckResourceAttr(datasourceNameVMs, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.#"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "nics.0.network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "nics.0.network_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "cd_roms.#"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "cd_roms.0.disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(datasourceNameVMs, "cd_roms.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "cd_roms.0.backing_info.0.data_source.#"),
				),
			},
		},
	})
}

func testAccVMDataSourceConfigV4(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
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
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
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

func testAccVMDataSourceConfigV4WithCdrom(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
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
			boot_config{
				legacy_boot{
				  boot_order = ["CDROM", "DISK","NETWORK"]
				}
			}
			cd_roms{
				disk_address{
					bus_type = "SATA"
					index= 0
				}
			}
		}

		data "nutanix_virtual_machine_v2" "test" {
			ext_id = nutanix_virtual_machine_v2.vm1.id
		}
`, name, desc, filepath)
}

func testAccVMDataSourceConfigV4WithCdromBackingInfo(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = (jsondecode(file("%[3]s")))
		  	vmm    = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "testWithDisk"{
			name= "%[1]s-second"
			description =  "%[2]s-second"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			power_state = "OFF"
		}
	
		resource "nutanix_virtual_machine_v2" "vm1"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			boot_config{
				legacy_boot{
				  boot_order = ["CDROM", "DISK","NETWORK"]
				}
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
			cd_roms{
				disk_address{
					bus_type = "IDE"
					index= 0
				}
				backing_info{
					storage_container {
						ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
					}
					data_source {
						reference{
							vm_disk_reference {
								disk_address{
									bus_type="SCSI"
									index = 0
								}
								vm_reference{
									ext_id= resource.nutanix_virtual_machine_v2.testWithDisk.id
								}
							}
						}
					}
				}
			}
			power_state = "ON"
			lifecycle{
				ignore_changes = [
					cd_roms.0.backing_info.0.data_source
				]
			}
		}

		data "nutanix_virtual_machine_v2" "test" {
			ext_id = nutanix_virtual_machine_v2.vm1.id
		}
`, name, desc, filepath)
}

func TestAccV2NutanixVmsDatasource_WithMemorySizeGib(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-ds-gib-%d", r)
	desc := "test vm datasource with memory_size_gib"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4WithMemorySizeGib(name, desc, 4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMs, "name", name),
					resource.TestCheckResourceAttr(datasourceNameVMs, "memory_size_gib", "4"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "memory_size_bytes"),
				),
			},
		},
	})
}

func testAccVMDataSourceConfigV4WithMemorySizeGib(name, desc string, memorySizeGib int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_virtual_machine_v2" "test_res" {
			name                 = "%[1]s"
			description          = "%[2]s"
			num_cores_per_socket = 1
			num_sockets          = 1
			memory_size_gib      = %[3]d
			cluster {
				ext_id = local.cluster0
			}
		}

		data "nutanix_virtual_machine_v2" "test" {
			ext_id = nutanix_virtual_machine_v2.test_res.id
		}
`, name, desc, memorySizeGib)
}
