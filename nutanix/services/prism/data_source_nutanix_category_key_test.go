package prism_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixCategoryKeyDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryKeyDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_category_key.test_key", "id"),
					resource.TestCheckResourceAttr(
						"data.nutanix_category_key.test_key", "values.#", "0"),
				),
			},
		},
	})
}

func TestAccNutanixCategoryKeyDataSource_withValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryKeyDataSourceConfigWithValues,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_category_key.test_key_value", "id"),
					resource.TestCheckResourceAttr(
						"data.nutanix_category_key.test_key_value", "values.#", "1"),
				),
			},
		},
	})
}

const testAccCategoryKeyDataSourceConfig = `
resource "nutanix_category_key" "test_key"{
    name = "data_source_category_key_test"
    description = "Data Source CategoryKey Test"
}


data "nutanix_category_key" "test_key" {
	name = nutanix_category_key.test_key.name
}`

const testAccCategoryKeyDataSourceConfigWithValues = `
resource "nutanix_category_key" "test_key_value"{
    name = "data_source_category_key_test_values"
    description = "Data Source CategoryKey Test with Values"
}

resource "nutanix_category_value" "test_value"{
	name = nutanix_category_key.test_key_value.name
	value = "test_category_value_data_source"
    description = "Data Source CategoryValue Test with Values"
}


data "nutanix_category_key" "test_key_value" {
	name = nutanix_category_value.test_value.name //creating implicit dependency
}`
