package iamv2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNutanixUserKeyV2 = "data.nutanix_user_key_v2.get_key"

// Expiry time is two week later
var expirationTime = time.Now().Add(14 * 24 * time.Hour)
var expirationTimeFormatted = expirationTime.UTC().Format(time.RFC3339)

func TestAccV2NutanixUsersDataSourceKey(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-revoke-api-%d", r)
	key_name := fmt.Sprintf("tf-revoke-api-key-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testApiKeyDataSourceConfig(name, key_name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeyV2, "name", key_name),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeyV2, "key_type", "API_KEY"),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeyV2, "expiry_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeyV2, "assigned_to", "user1"),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeyV2, "status", "VALID"),
				),
			},
		},
	})
}

func TestAccV2NutanixUsersDataSourceKeyInvalid(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-revoke-api-%d", r)
	key_name := fmt.Sprintf("tf-revoke-api-key-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testApiKeyDataSourceConfigInvalid(name, key_name, expirationTimeFormatted),
				ExpectError: regexp.MustCompile("error while fetching the user key"),
			},
		},
	})
}

func testApiKeyDataSourceConfig(name string, key_name string, expirationTimeFormatted string) string {
	return fmt.Sprintf(`
	// Create service account
	resource "nutanix_users_v2" "service_account" {
		username = "%[2]s"
		description = "test service account tf"
		email_id = "terraform_plugin@domain.com"
		user_type = "SERVICE_ACCOUNT"
	}

	// Create key
	resource "nutanix_user_key_v2" "create_key" {
    user_ext_id = nutanix_users_v2.service_account.ext_id
    name = "%[3]s"
    key_type = "API_KEY"
	  expiry_time = "%[4]s"
	  assigned_to = "user1"
  }
	
	// Get key
	data "nutanix_user_key_v2" "get_key"{
  	user_ext_id = nutanix_users_v2.service_account.ext_id
  	ext_id = nutanix_user_key_v2.create_key.ext_id
	}
	`, filepath, name, key_name, expirationTimeFormatted)
}

func testApiKeyDataSourceConfigInvalid(name string, key_name string, expirationTimeFormatted string) string {
	return fmt.Sprintf(`
	// Create service account
	resource "nutanix_users_v2" "service_account" {
		username = "%[2]s"
		description = "test service account tf"
		email_id = "terraform_plugin@domain.com"
		user_type = "SERVICE_ACCOUNT"
	}

	// Create key
	resource "nutanix_user_key_v2" "create_key" {
   user_ext_id = nutanix_users_v2.service_account.ext_id
   name = "%[3]s"
   key_type = "API_KEY"
	 expiry_time = "%[4]s"
	 assigned_to = "user1"
  }
	
	// Get key
	data "nutanix_user_key_v2" "get_key"{
  	user_ext_id = nutanix_users_v2.service_account.ext_id
  	ext_id = "1234"
	}
	`, filepath, name, key_name, expirationTimeFormatted)
}
