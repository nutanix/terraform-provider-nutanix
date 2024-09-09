package iamv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameUser = "data.nutanix_user_v2.test"

func TestAccNutanixUserV4Datasource_Basic(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserDatasourceV4Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameUser, "username", testVars.Iam.Users.Username),
					resource.TestCheckResourceAttr(datasourceNameUser, "first_name", testVars.Iam.Users.FirstName),
					resource.TestCheckResourceAttr(datasourceNameUser, "middle_initial", testVars.Iam.Users.MiddleInitial),
					resource.TestCheckResourceAttr(datasourceNameUser, "last_name", testVars.Iam.Users.LastName),
					resource.TestCheckResourceAttr(datasourceNameUser, "email_id", testVars.Iam.Users.EmailId),
					resource.TestCheckResourceAttr(datasourceNameUser, "locale", testVars.Iam.Users.Locale),
					resource.TestCheckResourceAttr(datasourceNameUser, "region", testVars.Iam.Users.Region),
					resource.TestCheckResourceAttr(datasourceNameUser, "display_name", testVars.Iam.Users.DisplayName),
					resource.TestCheckResourceAttr(datasourceNameUser, "user_type", "LOCAL"),
					resource.TestCheckResourceAttr(datasourceNameUser, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testUserDatasourceV4Config(filepath string) string {
	return fmt.Sprintf(`

		locals{
			config = (jsondecode(file("%s")))
			users = local.config.iam.users
		}
		
		resource "nutanix_users_v2" "test" {
			username = local.users.username
			first_name = local.users.first_name
			middle_initial = local.users.middle_initial
			last_name = local.users.last_name
			email_id = local.users.email_id
			locale = local.users.locale
			region = local.users.region
			display_name = local.users.display_name
			password = local.users.password
			user_type = "LOCAL"
			status = "ACTIVE"  
			force_reset_password = local.users.force_reset_password   
		}
		
		data "nutanix_user_v2" "test" {
			ext_id = nutanix_users_v2.test.id
			depends_on = [nutanix_users_v2.test]
		}			

		
	`, filepath)
}
