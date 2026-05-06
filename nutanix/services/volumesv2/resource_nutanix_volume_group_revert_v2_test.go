package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVolumeGroupRevert = "nutanix_volume_group_revert_v2.test"

func TestAccV2NutanixVolumeGroupRevertResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-vg-revert-%d", r)
	desc := "terraform test volume group revert"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupRevertResourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVolumeGroupRevert, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameVolumeGroupRevert, "volume_group_recovery_point_ext_id"),
				),
			},
		},
	})
}

func testAccVolumeGroupRevertResourceConfig(name, desc string) string {
	return testAccVolumeGroupResourceConfig(name, desc) + fmt.Sprintf(`
		resource "nutanix_volume_group_revert_v2" "test" {
			ext_id                             = resource.nutanix_volume_group_v2.test.id
			volume_group_recovery_point_ext_id = "00000000-0000-0000-0000-000000000000"
			depends_on = [resource.nutanix_volume_group_v2.test]
		}
	`)
}
