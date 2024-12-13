package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixPermissionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionDataSourceConfig(testVars.Permissions[1].UUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "name", testVars.Permissions[1].Name),
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "operation", "delete"),
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "fields.0.field_mode", "NONE"),
				),
			},
		},
	})
}

func testAccPermissionDataSourceConfig(uuid string) string {
	return fmt.Sprintf(`
data "nutanix_permission" "test" {
	permission_id = "%s"
}
`, uuid)
}

func TestAccNutanixPermissionDataSource_basicByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionDataSourceConfigByName(testVars.Permissions[1].Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "name", testVars.Permissions[1].Name),
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "operation", "delete"),
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "fields.0.field_mode", "NONE"),
				),
			},
		},
	})
}

func testAccPermissionDataSourceConfigByName(name string) string {
	return fmt.Sprintf(`
data "nutanix_permission" "test" {
	permission_name = "%s"
}
`, name)
}
