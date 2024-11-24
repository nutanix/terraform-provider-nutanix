package prismv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameCategories = "data.nutanix_categories_v2.test"

func TestAccNutanixCategoriesV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoriesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.#"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.key"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.value"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.type"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.associations.#"),
				),
			},
		},
	})
}

func TestAccNutanixCategoriesV2DataSource_WithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoriesDataSourceConfigWithFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.#"),
					resource.TestCheckResourceAttr(datasourceNameCategories, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.key"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.value"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.type"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.associations.#"),
				),
			},
		},
	})
}

func TestAccNutanixCategoriesV2DataSource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoriesDataSourceConfigWithLimit(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.#"),
					resource.TestCheckResourceAttr(datasourceNameCategories, "categories.#", "2"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.key"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.value"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.type"),
					resource.TestCheckResourceAttrSet(datasourceNameCategories, "categories.0.associations.#"),
				),
			},
		},
	})
}

func testAccCategoriesDataSourceConfig() string {
	return (`
		data "nutanix_categories_v2" "test" { }
	`)
}

func testAccCategoriesDataSourceConfigWithFilter() string {
	return (`

		data "nutanix_categories_v2" "dtest" { }

		locals{
			kk = data.nutanix_categories_v2.dtest.categories.0.key
		}
		data "nutanix_categories_v2" "test" {
			filter = "key eq '${local.kk}'"
			depends_on = [
				data.nutanix_categories_v2.dtest
			]
		}
	`)
}

func testAccCategoriesDataSourceConfigWithLimit() string {
	return (`
		data "nutanix_categories_v2" "test" {
			limit = 2
		}
	`)
}
