package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRole = "data.nutanix_role_v2.test"

func TestAccV2NutanixRolesDatasource_Basic_Role(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRoleDatasourceV2Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRole, "display_name"),
					resource.TestCheckResourceAttr(datasourceNameRole, "display_name", testVars.Iam.Roles.DisplayName),
					resource.TestCheckResourceAttr(datasourceNameRole, "description", testVars.Iam.Roles.Description),
				),
			},
		},
	})
}

func testRoleDatasourceV2Config(filepath string) string {
	return fmt.Sprintf(`

		locals{
			config = (jsondecode(file("%s")))
			roles = local.config.iam.roles
		}

		data "nutanix_operations_v2" "test" {
			filter = "startswith(displayName, 'Create_')"
		}

		resource "nutanix_roles_v2" "test" {
			display_name = local.roles.display_name
			description  = local.roles.description
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
	`, filepath)
}
