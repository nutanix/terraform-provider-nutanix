package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameAuthorizationPolicies = "data.nutanix_authorization_policies_v2.test"

func authPolicyConfig(acpDisplayName, roleDisplayName, roleDescription string) string {
	return fmt.Sprintf(`
data "nutanix_operations_v2" "test" {
	filter = "startswith(displayName, 'Create_')"
}

resource "nutanix_roles_v2" "test" {
	display_name = "%[1]s"
	description  = "%[2]s"
	operations = [
		data.nutanix_operations_v2.test.operations[0].ext_id,
		data.nutanix_operations_v2.test.operations[1].ext_id,
		data.nutanix_operations_v2.test.operations[2].ext_id,
		data.nutanix_operations_v2.test.operations[3].ext_id
	]
	depends_on = [data.nutanix_operations_v2.test]
}
resource "nutanix_authorization_policy_v2" "auth_policy_test" {
	role         = nutanix_roles_v2.test.id
	display_name = "%[3]s"
	description  = "%[4]s"
	authorization_policy_type = "%[5]s"
	%[6]s
	depends_on = [nutanix_roles_v2.test]
}
  `, roleDisplayName, roleDescription, acpDisplayName, authPolicyDescription, authPolicyType, authPolicyIdentitiesEntitiesHCL())
}

func TestAccV2NutanixAuthorizationPoliciesDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPoliciesDatasourceV4Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameAuthorizationPolicies, "auth_policies.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixAuthorizationPoliciesDatasource_WithFilter(t *testing.T) {
	acpDisplayName := fmt.Sprintf("tf-test-acp-%d", acctest.RandInt())
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPoliciesDatasourceV4WithFilterConfig(acpDisplayName, roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameAuthorizationPolicies, "auth_policies.#"),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicies, "auth_policies.0.display_name", acpDisplayName),
				),
			},
		},
	})
}

func TestAccV2NutanixAuthorizationPoliciesDatasource_WithLimit(t *testing.T) {
	acpDisplayName := fmt.Sprintf("tf-test-acp-%d", acctest.RandInt())
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPoliciesDatasourceV4WithLimitConfig(acpDisplayName, roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameAuthorizationPolicies, "auth_policies.#"),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicies, "auth_policies.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixAuthorizationPoliciesDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPoliciesDatasourceV4WithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicies, "auth_policies.#", "0"),
				),
			},
		},
	})
}

func testAuthorizationPoliciesDatasourceV4Config() string {
	return `
	data "nutanix_authorization_policies_v2" "test"{}
	`
}

func testAuthorizationPoliciesDatasourceV4WithFilterConfig(acpDisplayName, roleDisplayName, roleDescription string) string {
	return fmt.Sprintf(`
	%[1]s

	data "nutanix_authorization_policies_v2" "test" {
		filter = "displayName eq '%[2]s'"
		depends_on = [resource.nutanix_authorization_policy_v2.auth_policy_test]
	}
	`, authPolicyConfig(acpDisplayName, roleDisplayName, roleDescription), acpDisplayName)
}

func testAuthorizationPoliciesDatasourceV4WithLimitConfig(acpDisplayName, roleDisplayName, roleDescription string) string {
	return fmt.Sprintf(`
		%s

		data "nutanix_authorization_policies_v2" "test" {
			limit     = 1
			depends_on = [resource.nutanix_authorization_policy_v2.auth_policy_test]
		}
	`, authPolicyConfig(acpDisplayName, roleDisplayName, roleDescription))
}

func testAuthorizationPoliciesDatasourceV4WithInvalidFilterConfig() string {
	return `
	data "nutanix_authorization_policies_v2" "test" {
		filter = "displayName eq 'invalid_filter'"
	}
	`
}
