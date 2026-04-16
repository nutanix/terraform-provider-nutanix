package datapoliciesv2_test

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// getProviderConfigForAPINonSupportedTests returns a provider block that uses basic auth from test_config_v2.json.
// Use for tests that do not support API key (e.g., protection policy). Prepend to Terraform config.
func getProviderConfigForAPINonSupportedTests() string {
	username := testVars.UsernameForTest
	password := testVars.PasswordForTest
	if username == "" || password == "" {
		username = os.Getenv("NUTANIX_USERNAME")
		password = os.Getenv("NUTANIX_PASSWORD")
	}
	port := os.Getenv("NUTANIX_PORT")
	if port == "" {
		port = "9440"
	}
	insecure := os.Getenv("NUTANIX_INSECURE")
	if insecure == "" {
		insecure = "true"
	}
	endpoint := os.Getenv("NUTANIX_ENDPOINT")
	return fmt.Sprintf(`
provider "nutanix" {
  username  = "%s"
  password  = "%s"
  endpoint  = "%s"
  port      = %s
  insecure  = %s
  api_key   = ""  # Force basic auth; protection policy does not support API key
}
`, username, password, endpoint, port, insecure)
}

func checkAttributeLength(resourceName, attribute string, minLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}

		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}

		if count < minLength {
			return fmt.Errorf("expected %s to be >= %d, got %d", attrKey, minLength, count)
		}

		return nil
	}
}

func checkAttributeLengthEqual(resourceName, attribute string, expectedLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}

		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}

		if count != expectedLength {
			return fmt.Errorf("expected %s to be %d, got %d", attrKey, expectedLength, count)
		}

		return nil
	}
}

func testProtectionPolicyV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	client := conn.DataPoliciesAPI.ProtectionPolicies

	for _, rs := range state.RootModule().Resources {
		if rs.Type == resourceNameProtectionPolicy {
			_, err := client.GetProtectionPolicyById(utils.StringPtr(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("protection policy still exists")
			}
			fmt.Printf("Protection Policy still exists")
			_, err = client.DeleteProtectionPolicyById(utils.StringPtr(rs.Primary.ID))
			if err != nil {
				return fmt.Errorf("error: protection policy still exists : %v", err)
			}
		}
	}

	return nil
}
