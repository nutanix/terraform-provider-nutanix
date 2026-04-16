package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroup = "data.nutanix_volume_group_v2.test"

func TestAccV2NutanixVolumeGroupDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupDataSourceConfig(filepath, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "name", name),
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "description", desc),
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "should_load_balance_vm_attachments", "false"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "sharing_status", "SHARED"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "created_by", "admin"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "iscsi_features.0.enabled_authentications", "CHAP"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "storage_features.0.flash_mode.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "usage_type", "USER"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroup, "is_hidden", "false"),
				),
			},
		},
	})
}

func testAccVolumeGroupDataSourceConfig(filepath, name, desc string) string {
	return testAccVolumeGroupResourceConfig(name, desc) + `
		data "nutanix_volume_group_v2" "test" {
			ext_id = resource.nutanix_volume_group_v2.test.id
			depends_on = [resource.nutanix_volume_group_v2.test]
		}
	`
}
