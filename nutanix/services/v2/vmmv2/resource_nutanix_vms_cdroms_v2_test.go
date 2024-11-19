package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmCdrom = "nutanix_vm_cdroms_v4.test"

func TestAccNutanixVmsCdromV4_Basic(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCdromV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "disk_address.#"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.#"),
				),
			},
		},
	})
}

func TestAccNutanixVmsCdromV4_SATA(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCdromV4SATAConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "disk_address.#"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.#"),
				),
			},
		},
	})
}
func TestAccNutanixVmsCdromV4_WithBackingInfo(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCdromV4WithBackingInfoConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "disk_address.#"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.#"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "backing_info.0.disk_size_bytes", "1073741824"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.0.storage_container.#"),
				),
			},
		},
	})
}

func TestAccNutanixVmsCdromV4_WithDataSource(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCdromV4WithDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "disk_address.#"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.#"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "backing_info.0.disk_size_bytes", "21474836480"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.0.storage_container.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.0.data_source.#"),
				),
			},
		},
	})
}

func TestAccNutanixVmsCdromV4_WithVmReference(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCdromV4WithVmReferenceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "disk_address.#"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "disk_address.0.index", "1"),
					resource.TestCheckResourceAttr(resourceNameVmCdrom, "iso_type", "OTHER"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.0.disk_size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.0.storage_container.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmCdrom, "backing_info.0.data_source.#"),
				),
			},
		},
	})
}

func testVmsCdromV4Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
	
		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}

		resource "nutanix_vm_cdroms_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 1
			}
		}
`, name, desc)
}

func testVmsCdromV4SATAConfig(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
	
		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}

		resource "nutanix_vm_cdroms_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "SATA"
			  index= 1
			}
		}
`, name, desc)
}

func testVmsCdromV4WithBackingInfoConfig(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
	
		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}

		resource "nutanix_vm_cdroms_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 1
			}
			backing_info{
				disk_size_bytes = 1073741824
				storage_container {
				ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
				}
			}
		}
`, name, desc)
}

func testVmsCdromV4WithDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
	
		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}

		resource "nutanix_vm_cdroms_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 1
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
`, name, desc)
}

func testVmsCdromV4WithVmReferenceConfig(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
	
		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
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
							ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
						}
					}
				}
			}
		}

		resource "nutanix_vm_cdroms_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 1
			}
			backing_info{
				storage_container {
					ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
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
			lifecycle{
				ignore_changes = [
					backing_info.0.data_source
				]
			}
		}
`, name, desc)
}
