package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixRoleDataSourceByID_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-role")
	description := "Description of my role"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfigByID(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_role.test", "name", name),
					resource.TestCheckResourceAttr(
						"data.nutanix_role.test", "description", description),
					resource.TestCheckResourceAttrSet("data.nutanix_role.test", "name"),
				),
			},
		},
	})
}

func TestAccNutanixRoleDataSourceByName_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-role")
	description := "Description of my role"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfigByName(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_role.test", "name", name),
					resource.TestCheckResourceAttr(
						"data.nutanix_role.test", "description", description),
					resource.TestCheckResourceAttrSet("data.nutanix_role.test", "name"),
				),
			},
		},
	})
}

func testAccRoleDataSourceConfigByID(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	permission_reference_list {
		kind = "permission"
		uuid = "%[3]s"
	}
}

data "nutanix_role" "test" {
	role_id = nutanix_role.test.id
}
`, name, description, testVars.Permissions[0].UUID)
}

func testAccRoleDataSourceConfigByName(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	permission_reference_list {
		kind = "permission"
		uuid = "%[3]s"
	}
}

data "nutanix_role" "test" {
	role_name = nutanix_role.test.name
}
`, name, description, testVars.Permissions[0].UUID)
}
