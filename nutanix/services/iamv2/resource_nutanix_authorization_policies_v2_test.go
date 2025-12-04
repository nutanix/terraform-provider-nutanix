package iamv2_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameAuthorizationPolicy = "nutanix_authorization_policy_v2.test"

func TestAccV2NutanixAuthorizationPolicyResource_CreateACP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPolicyResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameAuthorizationPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "display_name", testVars.Iam.AuthPolicies.DisplayName),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "description", testVars.Iam.AuthPolicies.Description),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "authorization_policy_type", testVars.Iam.AuthPolicies.AuthPolicyType),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "identities.#", strconv.Itoa(len(testVars.Iam.AuthPolicies.Identities))),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "identities.0.reserved", testVars.Iam.AuthPolicies.Identities[0]),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.#", strconv.Itoa(len(testVars.Iam.AuthPolicies.Entities))),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.0.reserved", testVars.Iam.AuthPolicies.Entities[0]),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.1.reserved", testVars.Iam.AuthPolicies.Entities[1]),
				),
			},
			// test update ac
			{
				Config: testAuthorizationPolicyResourceUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameAuthorizationPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "description", testVars.Iam.AuthPolicies.Description+"_updated"),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "display_name", testVars.Iam.AuthPolicies.DisplayName),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "authorization_policy_type", testVars.Iam.AuthPolicies.AuthPolicyType),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "identities.#", strconv.Itoa(len(testVars.Iam.AuthPolicies.Identities))),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "identities.0.reserved", testVars.Iam.AuthPolicies.Identities[0]),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.#", strconv.Itoa(len(testVars.Iam.AuthPolicies.Entities))),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.0.reserved", testVars.Iam.AuthPolicies.Entities[0]),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.1.reserved", testVars.Iam.AuthPolicies.Entities[1]),
				),
			},
		},
	})
}

func TestAccV2NutanixAuthorizationPolicyResource_WithNoDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAuthorizationPolicyResourceWithoutDisplayNameConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixAuthorizationPolicyResource_WithNoIdentities(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAuthorizationPolicyResourceWithoutIdentitiesConfig(),
				ExpectError: regexp.MustCompile("Insufficient identities blocks"),
			},
		},
	})
}

func TestAccV2NutanixAuthorizationPolicyResource_WithNoEntities(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAuthorizationPolicyResourceWithoutEntitiesConfig(),
				ExpectError: regexp.MustCompile("Insufficient entities blocks"),
			},
		},
	})
}

func TestAccV2NutanixAuthorizationPolicyResource_WithNoRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAuthorizationPolicyResourceWithoutRoleConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testAuthorizationPolicyResourceConfig() string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		auth_policies = local.config.iam.auth_policies
		roles = local.config.iam.roles
	}

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

	resource "nutanix_authorization_policy_v2" "test" {
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

	}`, filepath)
}

func testAuthorizationPolicyResourceUpdateConfig() string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		auth_policies = local.config.iam.auth_policies
		roles = local.config.iam.roles
	}

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

	resource "nutanix_authorization_policy_v2" "test" {
		role         =  nutanix_roles_v2.test.id
		display_name = local.auth_policies.display_name
		description  = "${local.auth_policies.description}_updated"
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

	}`, filepath)
}

func testAuthorizationPolicyResourceWithoutDisplayNameConfig() string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		auth_policies = local.config.iam.auth_policies
	}

	resource "nutanix_authorization_policy_v2" "test" {
		role         = local.auth_policies.role
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

	}`, filepath)
}

func testAuthorizationPolicyResourceWithoutIdentitiesConfig() string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		auth_policies = local.config.iam.auth_policies
	}

	resource "nutanix_authorization_policy_v2" "test" {
		role         = local.auth_policies.role
		display_name = local.auth_policies.display_name
		description  = local.auth_policies.description
		authorization_policy_type = local.auth_policies.authorization_policy_type

		entities {
			reserved = local.auth_policies.entities[0]
		}
		entities {
			reserved = local.auth_policies.entities[1]
		}
	}`, filepath)
}

func testAuthorizationPolicyResourceWithoutEntitiesConfig() string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		auth_policies = local.config.iam.auth_policies
	}

	resource "nutanix_authorization_policy_v2" "test" {
		role         = local.auth_policies.role
		display_name = local.auth_policies.display_name
		description  = local.auth_policies.description
		authorization_policy_type = local.auth_policies.authorization_policy_type
		identities {
			reserved = local.auth_policies.identities[0]
		}

	}`, filepath)
}

func testAuthorizationPolicyResourceWithoutRoleConfig() string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		auth_policies = local.config.iam.auth_policies
	}

	resource "nutanix_authorization_policy_v2" "test" {
		display_name = local.auth_policies.display_name
		description  = local.auth_policies.description
		authorization_policy_type = local.auth_policies.authorization_policy_type
		identities {
			reserved = local.auth_policies.identities[0]
		}
		entities {
			reserved = local.auth_policies.entities[0]
		}

	}`, filepath)
}
