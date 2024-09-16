package volumesv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupsVmsAttachments = "data.nutanix_volume_group_vms_v2.test"

func TestAccNutanixVolumeGroupVmsAttachmentsV2DataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v4.json"
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
	return testAccVolumeGroupResourceConfig(filepath, name, desc) + `		
		resource "nutanix_volume_group_vm_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			vm_ext_id           = local.volumes.vm_ext_id
			depends_on          = [resource.nutanix_volume_group_v2.test]
		}		
		
		data "nutanix_volume_group_vms_v2" "test" {
			ext_id = resource.nutanix_volume_group_v2.test.id
			depends_on = [ resource.nutanix_volume_group_vm_v2.test ]			
		}
	`
}
