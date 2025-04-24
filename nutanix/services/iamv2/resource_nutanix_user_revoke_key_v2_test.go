package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNutanixUserRevokeKeyV2Create = "nutanix_user_revoke_key_v2.revoke_key"

func TestAccV2NutanixUsersRevokeKey(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-revoke-api-%d", r)
	key_name := fmt.Sprintf("tf-revoke-api-key-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testApiKeyRevokeResourceConfig(filepath, name, key_name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNutanixUserRevokeKeyV2Create, "message", "Revoke user key successful."),
				),
			},
		},
	})
}

func testApiKeyRevokeResourceConfig(filepath, name string, key_name string) string {
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
	 expiry_time = "2026-01-01T00:00:00Z"
  }

	resource "nutanix_user_revoke_key_v2" "revoke_key" {
	 user_ext_id = nutanix_users_v2.service_account.ext_id
	 ext_id = nutanix_user_key_v2.create_key.ext_id
  }
	`, filepath, name, key_name)
}
