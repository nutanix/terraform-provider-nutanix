package prismv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameCatg = "data.nutanix_category_v2.test"

func TestAccNutanixCategoryDataSourceV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameCatg, "description"),
					resource.TestCheckResourceAttrSet(datasourceNameCatg, "key"),
					resource.TestCheckResourceAttrSet(datasourceNameCatg, "value"),
					resource.TestCheckResourceAttrSet(datasourceNameCatg, "type"),
					resource.TestCheckResourceAttrSet(datasourceNameCatg, "associations.#"),
					resource.TestCheckResourceAttrSet(datasourceNameCatg, "detailed_associations.#"),
				),
			},
		},
	})
}

func testAccCategoryDataSourceConfig() string {
	return (`
		data "nutanix_categories_v2" "dtest" { }

		data "nutanix_category_v2" "test" {
			ext_id = data.nutanix_categories_v2.dtest.categories.0.ext_id

			depends_on = [
				data.nutanix_categories_v2.dtest
			]
		}
	`)
}
