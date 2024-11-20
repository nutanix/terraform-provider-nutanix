package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmCdromInsertEject = "nutanix_vm_cdrom_insert_eject_v2.test"

func TestAccNutanixVmsCdromInsertEjectV2Resource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCdromInsertEjectV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmCdromInsertEject, "backing_info.#"),
				),
			},
		},
	})
}

func testVmsCdromInsertEjectV2Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[1].metadata.uuid
		}

		data "nutanix_images_v2" "images" {
			filter = "name eq 'Nutanix-VirtIO-1.1.4.iso'"
		}
	
		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			cd_roms{
				disk_address{
					bus_type = "IDE"
					index= 0
				}
			}
			power_state = "ON"
		}

		resource "nutanix_vm_cdrom_insert_eject_v2" "test" {
			vm_ext_id= resource.nutanix_virtual_machine_v2.test.id
			ext_id = resource.nutanix_virtual_machine_v2.test.cd_roms.0.ext_id
			backing_info{
			  data_source{
				reference{
				  image_reference{
					image_ext_id = data.nutanix_images_v2.images.images[0].ext_id
				  }
				}
			  }
			}
		}
`, name, desc)
}
