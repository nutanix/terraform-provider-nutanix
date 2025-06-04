package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNutanixUserRevokeKeyV2Create = "nutanix_user_key_revoke_v2.revoke_key"
const dataSourceNutanixRevokeKeyV2 = "data.nutanix_user_key_v2.get_revoke_key"

func TestAccV2NutanixUsersRevokeKey(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-revoke-api-%d", r)
	keyName := fmt.Sprintf("tf-revoke-api-key-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAPIKeyRevokeResourceConfig(name, keyName, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNutanixUserRevokeKeyV2Create, "message", "Revoke user key successful."),
					resource.TestCheckResourceAttr(dataSourceNutanixRevokeKeyV2, "status", "REVOKED"),
				),
			},
		},
	})
}

func testAPIKeyRevokeResourceConfig(name string, keyName string, expirationTimeFormatted string) string {
	return fmt.Sprintf(`
	resource "nutanix_users_v2" "service_account" {
		username = "%[2]s"
		description = "test service account tf"
		email_id = "terraform_plugin@domain.com"
		user_type = "SERVICE_ACCOUNT"
	}

	resource "nutanix_user_key_v2" "create_key" {
		user_ext_id = nutanix_users_v2.service_account.ext_id
		name = "%[3]s"
		key_type = "API_KEY"
		expiry_time = "%[4]s"
  }

	resource "nutanix_user_key_revoke_v2" "revoke_key" {
		user_ext_id = nutanix_users_v2.service_account.ext_id
		ext_id = nutanix_user_key_v2.create_key.ext_id
  }

	data "nutanix_user_key_v2" "get_revoke_key" {
		user_ext_id = nutanix_users_v2.service_account.ext_id
		ext_id = nutanix_user_key_v2.create_key.ext_id
		depends_on = [nutanix_user_key_revoke_v2.revoke_key]
	}
	`, filepath, name, keyName, expirationTimeFormatted)
}
