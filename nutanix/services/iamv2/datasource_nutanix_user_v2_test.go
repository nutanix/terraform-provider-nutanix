package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameUser = "data.nutanix_user_v2.test"
const datasourceNameServiceAccount = "data.nutanix_user_v2.get_service_account"

func TestAccV2NutanixUserDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-user-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		// using V3 API to delete user
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUserDatasourceV4Config(filepath, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameUser, "username", name),
					resource.TestCheckResourceAttr(datasourceNameUser, "first_name", "first-name-"+name),
					resource.TestCheckResourceAttr(datasourceNameUser, "middle_initial", "middle-initial-"+name),
					resource.TestCheckResourceAttr(datasourceNameUser, "last_name", "last-name-"+name),
					resource.TestCheckResourceAttr(datasourceNameUser, "email_id", testVars.Iam.Users.EmailID),
					resource.TestCheckResourceAttr(datasourceNameUser, "locale", testVars.Iam.Users.Locale),
					resource.TestCheckResourceAttr(datasourceNameUser, "region", testVars.Iam.Users.Region),
					resource.TestCheckResourceAttr(datasourceNameUser, "display_name", "display-name-"+name),
					resource.TestCheckResourceAttr(datasourceNameUser, "user_type", "LOCAL"),
					resource.TestCheckResourceAttr(datasourceNameUser, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccV2NutanixUserDatasourceServiceAccount(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-service-account-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceAccountDataSourceV4Config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameServiceAccount, "username", name),
					resource.TestCheckResourceAttr(datasourceNameServiceAccount, "description", "test service account tf"),
					resource.TestCheckResourceAttr(datasourceNameServiceAccount, "email_id", "terraform_plugin@domain.com"),
					resource.TestCheckResourceAttr(datasourceNameServiceAccount, "user_type", "SERVICE_ACCOUNT"),
				),
			},
		},
	})
}

func testUserDatasourceV4Config(filepath, name string) string {
	return fmt.Sprintf(`

		locals{
			config = (jsondecode(file("%[1]s")))
			users = local.config.iam.users
		}

		resource "nutanix_users_v2" "test" {
			username = "%[2]s"
			first_name = "first-name-%[2]s"
			middle_initial = "middle-initial-%[2]s"
			last_name = "last-name-%[2]s"
			email_id = local.users.email_id
			locale = local.users.locale
			region = local.users.region
			display_name = "display-name-%[2]s"
			password = local.users.password
			user_type = "LOCAL"
			status = "ACTIVE"
			force_reset_password = local.users.force_reset_password
		}

		data "nutanix_user_v2" "test" {
			ext_id = nutanix_users_v2.test.id
			depends_on = [nutanix_users_v2.test]
		}
	`, filepath, name)
}

func testServiceAccountDataSourceV4Config(name string) string {
	return fmt.Sprintf(`
	resource "nutanix_users_v2" "service_account" {
		username = "%[2]s"
		description = "test service account tf"
		email_id = "terraform_plugin@domain.com"
		user_type = "SERVICE_ACCOUNT"
	}

	data "nutanix_user_v2" "get_service_account" {
		ext_id = nutanix_users_v2.service_account.id
	}
	`, filepath, name)
}
