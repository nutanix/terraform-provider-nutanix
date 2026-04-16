package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameAuthorizationPolicies = "data.nutanix_authorization_policies_v2.test"

const authPolicy = `
data "nutanix_operations_v2" "test" {
	filter = "startswith(displayName, 'Create_')"
}

resource "nutanix_roles_v2" "test" {
	display_name = local.roles.display_name
	description  = local.roles.description
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
	display_name = local.auth_policies.display_name
	description  = local.auth_policies.description
	authorization_policy_type = local.auth_policies.authorization_policy_type
	identities {
		reserved = local.auth_policies.identities[0]
	}
	entities {
		reserved = local.auth_policies.entities[0]
	}
	entities {
		reserved = local.auth_policies.entities[1]
	}
	depends_on = [nutanix_roles_v2.test]
}
  `

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
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPoliciesDatasourceV4WithFilterConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameAuthorizationPolicies, "auth_policies.#"),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicies, "auth_policies.0.display_name", testVars.Iam.AuthPolicies.DisplayName),
				),
			},
		},
	})
}

func TestAccV2NutanixAuthorizationPoliciesDatasource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPoliciesDatasourceV4WithLimitConfig(filepath),
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

func testAuthorizationPoliciesDatasourceV4WithFilterConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		auth_policies = local.config.iam.auth_policies
		roles = local.config.iam.roles
	}

	%s

	data "nutanix_authorization_policies_v2" "test" {
		filter = "displayName eq '${local.auth_policies.display_name}'"
		depends_on = [resource.nutanix_authorization_policy_v2.auth_policy_test]
	}


	`, filepath, authPolicy)
}

func testAuthorizationPoliciesDatasourceV4WithLimitConfig(filepath string) string {
	return fmt.Sprintf(`
		locals{
			config = (jsondecode(file("%s")))
			auth_policies = local.config.iam.auth_policies
			roles = local.config.iam.roles
		}

		%s

		data "nutanix_authorization_policies_v2" "test" {
			limit     = 1
			depends_on = [resource.nutanix_authorization_policy_v2.auth_policy_test]
		}
	`, filepath, authPolicy)
}

func testAuthorizationPoliciesDatasourceV4WithInvalidFilterConfig() string {
	return `
	data "nutanix_authorization_policies_v2" "test" {
		filter = "displayName eq 'invalid_filter'"
	}
	`
}
