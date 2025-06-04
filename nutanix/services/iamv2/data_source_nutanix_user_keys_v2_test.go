package iamv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNutanixUserKeysV2 = "data.nutanix_user_keys_v2.get_keys"
const dataSourceNutanixUserKeysFilterV2 = "data.nutanix_user_keys_v2.get_keys_filter"

func TestAccV2NutanixUsersDataSourceKeys(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-revoke-api-%d", r)
	keyName := fmt.Sprintf("tf-revoke-api-key-%d", r)
	expectedKeys := []map[string]string{
		{
			"name":        keyName,
			"key_type":    "API_KEY",
			"status":      "VALID",
			"expiry_time": expirationTimeFormatted,
			"assigned_to": "user1",
		},
		{
			"name":        keyName + "_another",
			"key_type":    "API_KEY",
			"status":      "VALID",
			"expiry_time": expirationTimeFormatted,
			"assigned_to": "user1_another",
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAPIKeysDataSourceConfig(name, keyName, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNutanixUserKeysV2, "keys.#"),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeysV2, "keys.#", "2"),
					checkNutanixKeys(expectedKeys, dataSourceNutanixUserKeysV2),
					resource.TestCheckResourceAttrSet(dataSourceNutanixUserKeysFilterV2, "keys.#"),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeysFilterV2, "keys.0.name", keyName),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeysFilterV2, "keys.0.key_type", "API_KEY"),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeysFilterV2, "keys.0.status", "VALID"),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeysFilterV2, "keys.0.expiry_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(dataSourceNutanixUserKeysFilterV2, "keys.0.assigned_to", "user1"),
				),
			},
		},
	})
}

func testAPIKeysDataSourceConfig(name string, keyName string, expirationTimeFormatted string) string {
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

	// create another key
	resource "nutanix_user_key_v2" "create_another_key" {
   user_ext_id = nutanix_users_v2.service_account.ext_id
   name = "%[3]s_another"
   key_type = "API_KEY"
	 expiry_time = "%[4]s"
	 assigned_to = "user1_another"
  }

	// Data source to fetch the list of keys
	data "nutanix_user_keys_v2" "get_keys" {
    user_ext_id = nutanix_users_v2.service_account.ext_id
		depends_on = [nutanix_user_key_v2.create_key, nutanix_user_key_v2.create_another_key]
	}
	
	// Data source to fetch the key by name
	data "nutanix_user_keys_v2" "get_keys_filter" {
		user_ext_id = nutanix_users_v2.service_account.ext_id
		filter = "name eq '%[3]s'"
		depends_on = [nutanix_user_key_v2.create_key, nutanix_user_key_v2.create_another_key]
	}
	`, filepath, name, keyName, expirationTimeFormatted)
}

func checkNutanixKeys(expectedKeys []map[string]string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		keys := rs.Primary.Attributes
		keyCount, err := strconv.Atoi(rs.Primary.Attributes["keys.#"])
		if err != nil {
			return fmt.Errorf("error parsing keys count: %s", err)
		}

		if keyCount != len(expectedKeys) {
			return fmt.Errorf("expected %d keys, found %d", len(expectedKeys), keyCount)
		}

		// Match each expected key
		for _, expected := range expectedKeys {
			found := false
			for i := 0; i < keyCount; i++ {
				if keys[fmt.Sprintf("keys.%d.name", i)] == expected["name"] &&
					keys[fmt.Sprintf("keys.%d.key_type", i)] == expected["key_type"] &&
					keys[fmt.Sprintf("keys.%d.status", i)] == expected["status"] &&
					keys[fmt.Sprintf("keys.%d.expiry_time", i)] == expected["expiry_time"] &&
					keys[fmt.Sprintf("keys.%d.assigned_to", i)] == expected["assigned_to"] {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("expected key not found: %+v", expected)
			}
		}
		return nil
	}
}
