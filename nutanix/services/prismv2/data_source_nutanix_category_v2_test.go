package prismv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameCategory = "data.nutanix_category_v2.test"

func TestAccV2NutanixCategoryDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameCategory, "description"),
					resource.TestCheckResourceAttrSet(datasourceNameCategory, "key"),
					resource.TestCheckResourceAttrSet(datasourceNameCategory, "value"),
					resource.TestCheckResourceAttrSet(datasourceNameCategory, "type"),
					resource.TestCheckResourceAttrSet(datasourceNameCategory, "associations.#"),
					resource.TestCheckResourceAttrSet(datasourceNameCategory, "detailed_associations.#"),
				),
			},
		},
	})
}

func testAccCategoryDataSourceConfig() string {
	return `
		data "nutanix_categories_v2" "dtest" { }

		data "nutanix_category_v2" "test" {
			ext_id = data.nutanix_categories_v2.dtest.categories.0.ext_id

			depends_on = [
				data.nutanix_categories_v2.dtest
			]
		}
	`
}
