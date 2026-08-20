package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRole = "data.nutanix_role_v2.test"

func TestAccV2NutanixRolesDatasource_Basic_Role(t *testing.T) {
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRoleDatasourceV2Config(roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRole, "display_name"),
					resource.TestCheckResourceAttr(datasourceNameRole, "display_name", roleDisplayName),
					resource.TestCheckResourceAttr(datasourceNameRole, "description", roleDescription),
				),
			},
		},
	})
}

func testRoleDatasourceV2Config(displayName, description string) string {
	return fmt.Sprintf(`

		data "nutanix_operations_v2" "test" {
			filter = "startswith(displayName, 'Create_')"
		}

		resource "nutanix_roles_v2" "test" {
			display_name = "%[1]s"
			description  = "%[2]s"
			operations = [
				data.nutanix_operations_v2.test.operations[0].ext_id,
				data.nutanix_operations_v2.test.operations[1].ext_id,
				data.nutanix_operations_v2.test.operations[2].ext_id,
				data.nutanix_operations_v2.test.operations[3].ext_id
			]
			depends_on = [data.nutanix_operations_v2.test]
		}
		
		data "nutanix_role_v2" "test" {
			ext_id = resource.nutanix_roles_v2.test.id  
		}
	`, displayName, description)
}
