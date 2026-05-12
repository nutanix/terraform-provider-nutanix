package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupVmAttachments = "data.nutanix_volume_group_vm_attachments_v2.test"

func TestAccV2NutanixVolumeGroupVmAttachmentsDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-vm-attach-%d", r)
	desc := "terraform test volume group vm attachments description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupVmAttachmentsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupVmAttachments, "volume_group_ext_id"),
				),
			},
		},
	})
}

func testAccVolumeGroupVmAttachmentsDataSourceConfig(name, desc string) string {
	return testAccVolumeGroupResourceConfig(name, desc) + `
		data "nutanix_volume_group_vm_attachments_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			depends_on = [resource.nutanix_volume_group_v2.test]
		}
	`
}
