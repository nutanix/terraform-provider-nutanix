package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixAccessControlPolicyDataSourceByID_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	description := "Description of my access control policy"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessControlPolicyDataSourceByIDConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_access_control_policy.test", "name", name),
					resource.TestCheckResourceAttr(
						"data.nutanix_access_control_policy.test", "description", description),
					resource.TestCheckResourceAttrSet("data.nutanix_access_control_policy.test", "name"),
				),
			},
		},
	})
}

func TestAccNutanixAccessControlPolicyDataSourceByName_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	description := "Description of my access control policy"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessControlPolicyDataSourceByNameConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_access_control_policy.test", "name", name),
					resource.TestCheckResourceAttr(
						"data.nutanix_access_control_policy.test", "description", description),
					resource.TestCheckResourceAttrSet("data.nutanix_access_control_policy.test", "name"),
				),
			},
		},
	})
}

func testAccAccessControlPolicyDataSourceByIDConfig(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "test role"
	description = "description role"
	permission_reference_list {
		kind = "permission"
		uuid = "2e9988df-47ae-44ae-9114-ada346657b90"
	}
}
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference{
		kind = "role"
		uuid = nutanix_role.test.id
	}
}

data "nutanix_access_control_policy" "test" {
	access_control_policy_id = nutanix_access_control_policy.test.id
}
`, name, description)
}

func testAccAccessControlPolicyDataSourceByNameConfig(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "test role 2"
	description = "description role"
	permission_reference_list {
		kind = "permission"
		uuid = "2e9988df-47ae-44ae-9114-ada346657b90"
	}
}
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference{
		kind = "role"
		uuid = nutanix_role.test.id
	}
}

data "nutanix_access_control_policy" "test" {
	access_control_policy_name = nutanix_access_control_policy.test.name
}
`, name, description)
}
