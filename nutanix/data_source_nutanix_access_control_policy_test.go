package nutanix

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixAccessControlPolicyDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	description := "Description of my access control policy"
	uuidRole := ""

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
					resource.TestCheckResourceAttrSet("data.nutanix_access_control_policy.accest-access-policy", "name"),
				),
			},
		},
	})
}

func testAccAccessControlPolicyDataSourceConfig(uuidRole, name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_access_control_policy" "accest-access-policy" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference{
		kind = "role"
		uuid = "%[3]s"
	}
}

data "nutanix_access_control_policy" "test" {
	access_control_policy_id = nutanix_access_control_policy.accest-access-policy.id
}
`, name, description, uuidRole)
}
