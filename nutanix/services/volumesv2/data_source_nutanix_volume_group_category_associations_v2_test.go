package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupCategoryAssociations = "data.nutanix_volume_group_category_associations_v2.test"

func TestAccV2NutanixVolumeGroupCategoryAssociationsDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-vg-catassoc-%d", r)
	desc := "terraform test volume group category associations"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupCategoryAssociationsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupCategoryAssociations, "volume_group_ext_id"),
				),
			},
		},
	})
}

func testAccVolumeGroupCategoryAssociationsDataSourceConfig(name, desc string) string {
	return testAccVolumeGroupResourceConfig(name, desc) + `
		data "nutanix_volume_group_category_associations_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			depends_on = [resource.nutanix_volume_group_v2.test]
		}
	`
}
