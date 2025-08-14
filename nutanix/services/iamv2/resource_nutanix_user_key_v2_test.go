package iamv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNutanixUserKeyV2Create = "nutanix_user_key_v2.create_key"

func TestAccV2NutanixUsers_CreateKey(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-revoke-api-%d", r)
	keyName := fmt.Sprintf("tf-revoke-api-key-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAPIKeyCreateResourceConfig(name, keyName, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNutanixUserKeyV2Create, "name", keyName),
					resource.TestCheckResourceAttr(resourceNutanixUserKeyV2Create, "key_type", "API_KEY"),
					resource.TestCheckResourceAttr(resourceNutanixUserKeyV2Create, "expiry_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNutanixUserKeyV2Create, "assigned_to", "user1"),
					resource.TestCheckResourceAttr(resourceNutanixUserKeyV2Create, "status", "VALID"),
				),
			},
		},
	})
}

func TestAccV2NutanixUsers_CreateKey_DuplicateName(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-revoke-api-%d", r)
	keyName := fmt.Sprintf("tf-revoke-api-key-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAPIKeyCreateResourceConfigDuplicateName(name, keyName, expirationTimeFormatted),
				ExpectError: regexp.MustCompile("Failed to create key as there is a key with the same name."),
			},
		},
	})
}

func testAPIKeyCreateResourceConfig(name string, keyName string, expirationTimeFormatted string) string {
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
	 assigned_to = "user1"
  }
	`, filepath, name, keyName, expirationTimeFormatted)
}

func testAPIKeyCreateResourceConfigDuplicateName(name string, keyName string, expirationTimeFormatted string) string {
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
	 assigned_to = "user1"
  }

	resource "nutanix_user_key_v2" "create_key_dup_name" {
   user_ext_id = nutanix_users_v2.service_account.ext_id
   name = "%[3]s"
   key_type = "API_KEY"
	 expiry_time = 	"%[4]s"
	 assigned_to = "user1"
  }
	`, filepath, name, keyName, expirationTimeFormatted)
}
