package microsegv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	datasourceNameNSPRules = "data.nutanix_network_security_policy_rules_v2.test"
	policyResourceName     = "nutanix_network_security_policy_v2.test"
)

func TestAccV2NutanixNetworkSecurityPolicyRulesDataSource_Basic(t *testing.T) {
	r := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-test-nsp-rules-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityPolicyRulesDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPRules, "policy_ext_id"),
					resource.TestCheckResourceAttrPair(datasourceNameNSPRules, "policy_ext_id", policyResourceName, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPRules, "network_security_policy_rules.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.0.type", "TWO_ENV_ISOLATION"),
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.0.spec.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.0.spec.0.two_env_isolation_rule_spec.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.0.spec.0.two_env_isolation_rule_spec.0.first_isolation_group.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.0.spec.0.two_env_isolation_rule_spec.0.second_isolation_group.#", "1"),
					resource.TestCheckResourceAttrPair(
						datasourceNameNSPRules, "network_security_policy_rules.0.spec.0.two_env_isolation_rule_spec.0.first_isolation_group.0",
						policyResourceName, "rules.0.spec.0.two_env_isolation_rule_spec.0.first_isolation_group.0",
					),
					resource.TestCheckResourceAttrPair(
						datasourceNameNSPRules, "network_security_policy_rules.0.spec.0.two_env_isolation_rule_spec.0.second_isolation_group.0",
						policyResourceName, "rules.0.spec.0.two_env_isolation_rule_spec.0.second_isolation_group.0",
					),
				),
			},
		},
	})
}

func TestAccV2NutanixNetworkSecurityPolicyRulesDataSource_WithPagination(t *testing.T) {
	r := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-test-nsp-rules-pag-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityPolicyRulesDataSourceConfigWithPagination(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPRules, "policy_ext_id"),
					resource.TestCheckResourceAttrPair(datasourceNameNSPRules, "policy_ext_id", policyResourceName, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPRules, "network_security_policy_rules.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.0.type", "TWO_ENV_ISOLATION"),
					resource.TestCheckResourceAttr(datasourceNameNSPRules, "network_security_policy_rules.0.spec.0.two_env_isolation_rule_spec.#", "1"),
					resource.TestCheckResourceAttrPair(
						datasourceNameNSPRules, "network_security_policy_rules.0.spec.0.two_env_isolation_rule_spec.0.first_isolation_group.0",
						policyResourceName, "rules.0.spec.0.two_env_isolation_rule_spec.0.first_isolation_group.0",
					),
					resource.TestCheckResourceAttrPair(
						datasourceNameNSPRules, "network_security_policy_rules.0.spec.0.two_env_isolation_rule_spec.0.second_isolation_group.0",
						policyResourceName, "rules.0.spec.0.two_env_isolation_rule_spec.0.second_isolation_group.0",
					),
				),
			},
		},
	})
}

// TestAccV2NutanixNetworkSecurityPolicyRulesDataSource_InvalidPolicyExtID tests error handling
// when a non-existent policy ext_id is provided to the datasource.
func TestAccV2NutanixNetworkSecurityPolicyRulesDataSource_InvalidPolicyExtID(t *testing.T) {
	invalidPolicyExtID := "00000000-0000-0000-0000-000000000000"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkSecurityPolicyRulesDataSourceInvalidPolicyExtIDConfig(invalidPolicyExtID),
				ExpectError: regexp.MustCompile(`error listing network security policy rules`),
			},
		},
	})
}

// TestAccV2NutanixNetworkSecurityPolicyRulesDataSource_InvalidPolicyExtIDFormat tests error handling
// when an invalid policy_ext_id format is provided (e.g. not a valid UUID).
func TestAccV2NutanixNetworkSecurityPolicyRulesDataSource_InvalidPolicyExtIDFormat(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkSecurityPolicyRulesDataSourceInvalidPolicyExtIDConfig("invalid-ext-id"),
				ExpectError: regexp.MustCompile(`(?s)(error listing network security policy rules|SchemaValidationError|ECMA 262 regex|regular expression|regex).*invalid`),
			},
		},
	})
}

func testAccNetworkSecurityPolicyRulesDataSourceInvalidPolicyExtIDConfig(policyExtID string) string {
	return fmt.Sprintf(`
data "nutanix_network_security_policy_rules_v2" "test" {
  policy_ext_id = "%s"
}
`, policyExtID)
}

func testAccNetworkSecurityPolicyRulesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "nutanix_categories_v2" "test" {}

resource "nutanix_network_security_policy_v2" "test" {
  name        = "%[1]s"
  description = "test nsp for rules datasource"
  state       = "SAVE"
  type        = "ISOLATION"
  rules {
    type = "TWO_ENV_ISOLATION"
    spec {
      two_env_isolation_rule_spec {
        first_isolation_group = [
          data.nutanix_categories_v2.test.categories[0].ext_id,
        ]
        second_isolation_group = [
          data.nutanix_categories_v2.test.categories[1].ext_id,
        ]
      }
    }
  }
  is_hitlog_enabled = true
  depends_on        = [data.nutanix_categories_v2.test]
}

data "nutanix_network_security_policy_rules_v2" "test" {
  policy_ext_id = nutanix_network_security_policy_v2.test.ext_id
}
`, name)
}

func testAccNetworkSecurityPolicyRulesDataSourceConfigWithPagination(name string) string {
	return fmt.Sprintf(`
data "nutanix_categories_v2" "test" {}

resource "nutanix_network_security_policy_v2" "test" {
  name        = "%[1]s"
  description = "test nsp for rules datasource with pagination"
  state       = "SAVE"
  type        = "ISOLATION"
  rules {
    type = "TWO_ENV_ISOLATION"
    spec {
      two_env_isolation_rule_spec {
        first_isolation_group = [
          data.nutanix_categories_v2.test.categories[0].ext_id,
        ]
        second_isolation_group = [
          data.nutanix_categories_v2.test.categories[1].ext_id,
        ]
      }
    }
  }
  is_hitlog_enabled = true
  depends_on        = [data.nutanix_categories_v2.test]
}

data "nutanix_network_security_policy_rules_v2" "test" {
  policy_ext_id = nutanix_network_security_policy_v2.test.ext_id
  page         = 0
  limit        = 50
}
`, name)
}
