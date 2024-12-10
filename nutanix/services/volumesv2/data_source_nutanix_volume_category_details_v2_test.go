package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeCategoryDetails = "data.nutanix_volume_group_category_details_v2.test"

func TestAccV2NutanixVolumeCategoryDetailsDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccCategoryDetailsV2Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceVolumeCategoryDetails, "category_details.0.entity_type", "CATEGORY"),
					resource.TestCheckResourceAttr(dataSourceVolumeCategoryDetails, "category_details.0.name", ""),
					resource.TestCheckResourceAttrSet(dataSourceVolumeCategoryDetails, "category_details.#"),
				),
			},
		},
	})
}

func testAccCategoryDetailsV2Config(filepath string) string {
	return fmt.Sprintf(`
		locals {
			config = (jsondecode(file("%s")))
			volumes = local.config.volumes
		}

		data "nutanix_volume_group_category_details_v2" "test" {
			ext_id = local.volumes.vg_ext_id_with_category
		}
`, filepath)
}
