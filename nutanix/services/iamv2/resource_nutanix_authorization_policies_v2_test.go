package iamv2_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameAuthorizationPolicy = "nutanix_authorization_policy_v2.test"

func TestAccV2NutanixAuthorizationPolicyResource_CreateACP(t *testing.T) {
	acpDisplayName := fmt.Sprintf("tf-test-acp-%d", acctest.RandInt())
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPolicyResourceConfig(acpDisplayName, roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameAuthorizationPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "display_name", acpDisplayName),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "description", authPolicyDescription),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "authorization_policy_type", authPolicyType),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "identities.#", strconv.Itoa(len(authPolicyIdentities))),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "identities.0.reserved", authPolicyIdentities[0]),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.#", strconv.Itoa(len(authPolicyEntities))),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.0.reserved", authPolicyEntities[0]),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.1.reserved", authPolicyEntities[1]),
				),
			},
			// test update ac
			{
				Config: testAuthorizationPolicyResourceUpdateConfig(acpDisplayName, roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameAuthorizationPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "description", authPolicyDescription+"_updated"),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "display_name", acpDisplayName),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "authorization_policy_type", authPolicyType),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "identities.#", strconv.Itoa(len(authPolicyIdentities))),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "identities.0.reserved", authPolicyIdentities[0]),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.#", strconv.Itoa(len(authPolicyEntities))),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.0.reserved", authPolicyEntities[0]),
					resource.TestCheckResourceAttr(resourceNameAuthorizationPolicy, "entities.1.reserved", authPolicyEntities[1]),
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

func testAuthorizationPolicyResourceConfig(acpDisplayName, roleDisplayName, roleDescription string) string {
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

	resource "nutanix_authorization_policy_v2" "test" {
		role         = nutanix_roles_v2.test.id
		display_name = "%[3]s"
		description  = "%[4]s"
		authorization_policy_type = "%[5]s"
		%[6]s
		depends_on = [nutanix_roles_v2.test]

	}`, roleDisplayName, roleDescription, acpDisplayName, authPolicyDescription, authPolicyType, authPolicyIdentitiesEntitiesHCL())
}

func testAuthorizationPolicyResourceUpdateConfig(acpDisplayName, roleDisplayName, roleDescription string) string {
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

	resource "nutanix_authorization_policy_v2" "test" {
		role         =  nutanix_roles_v2.test.id
		display_name = "%[3]s"
		description  = "%[4]s_updated"
		authorization_policy_type = "%[5]s"
		%[6]s
		depends_on = [nutanix_roles_v2.test]

	}`, roleDisplayName, roleDescription, acpDisplayName, authPolicyDescription, authPolicyType, authPolicyIdentitiesEntitiesHCL())
}

func testAuthorizationPolicyResourceWithoutDisplayNameConfig() string {
	return fmt.Sprintf(`
	resource "nutanix_authorization_policy_v2" "test" {
		role         = "00000000-0000-0000-0000-000000000000"
		description  = "%[1]s"
		authorization_policy_type = "%[2]s"
		identities {
			reserved = %[3]q
		}
		entities {
			reserved = %[4]q
		}
		entities {
			reserved = %[5]q
		}

	}`, authPolicyDescription, authPolicyType, authPolicyIdentities[0], authPolicyEntities[0], authPolicyEntities[1])
}

func testAuthorizationPolicyResourceWithoutIdentitiesConfig() string {
	return fmt.Sprintf(`
	resource "nutanix_authorization_policy_v2" "test" {
		role         = "00000000-0000-0000-0000-000000000000"
		display_name = "tf-test-acp-no-identities"
		description  = "%[1]s"
		authorization_policy_type = "%[2]s"

		entities {
			reserved = %[3]q
		}
		entities {
			reserved = %[4]q
		}
	}`, authPolicyDescription, authPolicyType, authPolicyEntities[0], authPolicyEntities[1])
}

func testAuthorizationPolicyResourceWithoutEntitiesConfig() string {
	return fmt.Sprintf(`
	resource "nutanix_authorization_policy_v2" "test" {
		role         = "00000000-0000-0000-0000-000000000000"
		display_name = "tf-test-acp-no-entities"
		description  = "%[1]s"
		authorization_policy_type = "%[2]s"
		identities {
			reserved = %[3]q
		}

	}`, authPolicyDescription, authPolicyType, authPolicyIdentities[0])
}

func TestAccV2NutanixAuthorizationPolicyResource_ProjectAssociation(t *testing.T) {
	projectName := fmt.Sprintf("tf-acp-projassoc-%d", acctest.RandInt())
	acpDisplayName := fmt.Sprintf("tf-test-acp-%d", acctest.RandInt())
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPolicyProjectAssociationConfig(projectName, "", acpDisplayName, roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAuthorizationPolicy, "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_authorization_policy_v2.test", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttr("data.nutanix_authorization_policies_v2.test", "auth_policies.#", "1"),
					resource.TestCheckResourceAttrPair("data.nutanix_authorization_policies_v2.test", "auth_policies.0.ext_id", resourceNameAuthorizationPolicy, "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_authorization_policies_v2.test", "auth_policies.0.project_ext_id", "nutanix_project_v2.test", "ext_id"),
				),
			},
			{
				Config:      testAuthorizationPolicyProjectAssociationConfig(projectName, "00000000-0000-0000-0000-000000000000", acpDisplayName, roleDisplayName, roleDescription),
				ExpectError: regexp.MustCompile("Update of project_ext_id is not supported"),
			},
		},
	})
}

func acpProjectExtIDLine(override string) string {
	if override == "" {
		return `project_ext_id = nutanix_project_v2.test.ext_id`
	}
	return fmt.Sprintf(`project_ext_id = "%s"`, override)
}

func testAuthorizationPolicyProjectAssociationConfig(projectName, projectExtIDOverride, acpDisplayName, roleDisplayName, roleDescription string) string {
	return fmt.Sprintf(`
	resource "nutanix_project_v2" "test" {
		name        = "%[1]s"
		project_id  = "%[1]s"
		description = "project association test"
	}

	data "nutanix_operations_v2" "test" {
		filter = "startswith(displayName, 'Create_')"
	}

	resource "nutanix_roles_v2" "test" {
		display_name = "%[3]s"
		description  = "%[4]s"
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
		display_name = "%[5]s"
		description  = "%[6]s"
		authorization_policy_type = "%[7]s"
		%[2]s
		%[8]s
		depends_on = [nutanix_roles_v2.test, nutanix_project_v2.test]
	}

	data "nutanix_authorization_policy_v2" "test" {
		ext_id     = nutanix_authorization_policy_v2.test.id
		depends_on = [nutanix_authorization_policy_v2.test]
	}

	data "nutanix_authorization_policies_v2" "test" {
		filter     = "displayName eq '%[5]s'"
		depends_on = [nutanix_authorization_policy_v2.test]
	}
	`, projectName, acpProjectExtIDLine(projectExtIDOverride), roleDisplayName, roleDescription, acpDisplayName, authPolicyDescription, authPolicyType, authPolicyIdentitiesEntitiesHCL())
}

func testAuthorizationPolicyResourceWithoutRoleConfig() string {
	return fmt.Sprintf(`
	resource "nutanix_authorization_policy_v2" "test" {
		display_name = "tf-test-acp-no-role"
		description  = "%[1]s"
		authorization_policy_type = "%[2]s"
		identities {
			reserved = %[3]q
		}
		entities {
			reserved = %[4]q
		}

	}`, authPolicyDescription, authPolicyType, authPolicyIdentities[0], authPolicyEntities[0])
}
