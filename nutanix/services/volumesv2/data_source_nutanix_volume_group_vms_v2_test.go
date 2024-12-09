package volumesv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupsVmsAttachments = "data.nutanix_volume_group_vms_v2.test"

func TestAccNutanixVolumeGroupVmsAttachmentsV2DataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"
	path, _ := os.Getwd()
	filepath := path + "/../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupVmsAttachmentsDataSourceConfig(filepath, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupsVmsAttachments, "vms_attachments.#"),
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupsVmsAttachments, "vms_attachments.0.ext_id"),
				),
			},
		},
	})
}

func testAccVolumeGroupVmsAttachmentsDataSourceConfig(filepath, name, desc string) string {
	return testAccVolumeGroupResourceConfig(filepath, name, desc) + fmt.Sprintf(`
		resource "nutanix_virtual_machine_v2" "test"{
			name= "tf-test-vg-vm-%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster1
			}
			lifecycle{
				ignore_changes = [
					disks
				]
			}
		}

		resource "nutanix_volume_group_vm_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			vm_ext_id           = resource.nutanix_virtual_machine_v2.test.id
			depends_on          = [resource.nutanix_volume_group_v2.test]
		}		
		
		data "nutanix_volume_group_vms_v2" "test" {
			ext_id = resource.nutanix_volume_group_v2.test.id
			depends_on = [ resource.nutanix_volume_group_vm_v2.test ]			
		}
	`, name, desc)
}
