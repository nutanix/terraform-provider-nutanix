package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixAccessControlPolicyDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	description := "Description of my access control policy"
	uuidRole := "760dac6c-be97-4b24-adb0-e3c3026dc8d5"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessControlPolicyDataSourceConfig(uuidRole, name, description),
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

func testAccAccessControlPolicyDataSourceConfig(uuidRole, name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference{
		kind = "role"
		uuid = "%[3]s"
	}
}

data "nutanix_access_control_policy" "test" {
	access_control_policy_id = nutanix_access_control_policy.test.id
}
`, name, description, uuidRole)
}
