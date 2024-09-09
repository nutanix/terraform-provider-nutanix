package iamv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRole = "data.nutanix_role_v2.test"

func TestAccNutanixRolesV4Datasource_Basic_Role(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRoleDatasourceV4Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRole, "display_name"),
					resource.TestCheckResourceAttr(datasourceNameRole, "display_name", testVars.Iam.Roles.DisplayName),
					resource.TestCheckResourceAttr(datasourceNameRole, "description", testVars.Iam.Roles.Description),
				),
			},
		},
	})
}

func testRoleDatasourceV4Config(filepath string) string {
	return fmt.Sprintf(`

		locals{
			config = (jsondecode(file("%s")))
			roles = local.config.iam.roles
		}

		data "nutanix_operations_v2" "test" {
			limit = 3
		}

		resource "nutanix_roles_v2" "test" {
			display_name = local.roles.display_name
			description  = local.roles.description
			operations = [
				data.nutanix_operations_v2.test.permissions[0].ext_id,
				data.nutanix_operations_v2.test.permissions[1].ext_id,
				data.nutanix_operations_v2.test.permissions[2].ext_id,
			]
			depends_on = [data.nutanix_operations_v2.test]
		}
		
		data "nutanix_role_v2" "test" {
			ext_id = resource.nutanix_roles_v2.test.id  
		}
	`, filepath)
}
