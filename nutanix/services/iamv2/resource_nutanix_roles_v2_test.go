package iamv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRoles = "nutanix_roles_v2.test"

func TestAccV2NutanixRolesResource_Basic(t *testing.T) {
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRoleResourceConfig(roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRoles, "client_name"),
					resource.TestCheckResourceAttr(resourceNameRoles, "display_name", roleDisplayName),
					resource.TestCheckResourceAttr(resourceNameRoles, "description", roleDescription),
					resource.TestCheckResourceAttrSet(resourceNameRoles, "ext_id"),
				),
			},
			// update role
			{
				Config: testRoleResourceUpdateConfig(roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRoles, "client_name"),
					resource.TestCheckResourceAttr(resourceNameRoles, "display_name", fmt.Sprintf("%s_updated", roleDisplayName)),
					resource.TestCheckResourceAttr(resourceNameRoles, "description", roleDescription),
				),
			},
		},
	})
}

func TestAccV2NutanixRolesResource_IsGlobal(t *testing.T) {
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRoleResourceIsGlobalConfig(roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRoles, "client_name"),
					resource.TestCheckResourceAttr(resourceNameRoles, "display_name", roleDisplayName),
					resource.TestCheckResourceAttr(resourceNameRoles, "description", roleDescription),
					resource.TestCheckResourceAttr(resourceNameRoles, "is_global", "true"),
				),
			},
		},
	})
}

func TestAccV2NutanixRolesResource_DuplicateRole(t *testing.T) {
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testRoleResourceDuplicateRoleConfig(roleDisplayName, roleDescription),
				ExpectError: regexp.MustCompile("Failed to create role as already exists"),
			},
		},
	})
}

func TestAccV2NutanixRolesResource_WithNoDisplayName(t *testing.T) {
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testRoleResourceWithoutDisplayNameConfig(roleDescription),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixRolesResource_WithNoOperations(t *testing.T) {
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testRoleResourceWithoutOperationsConfig(roleDisplayName, roleDescription),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testRoleResourceConfig(displayName, description string) string {
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
	}`, displayName, description)
}

func testRoleResourceUpdateConfig(displayName, description string) string {
	return fmt.Sprintf(`

	data "nutanix_operations_v2" "test" {
	  //filter = "startswith(displayName, 'Create_')"
	  filter = "startswith(displayName, 'Create_')"
	}

	resource "nutanix_roles_v2" "test" {
		display_name = "%[1]s_updated"
		description  = "%[2]s"
		operations = [
			data.nutanix_operations_v2.test.operations[0].ext_id,
			data.nutanix_operations_v2.test.operations[1].ext_id,
			data.nutanix_operations_v2.test.operations[2].ext_id,
			data.nutanix_operations_v2.test.operations[3].ext_id
	  	]
		depends_on = [data.nutanix_operations_v2.test]
	}`, displayName, description)
}

func testRoleResourceDuplicateRoleConfig(displayName, description string) string {
	return fmt.Sprintf(`

	data "nutanix_operations_v2" "test" {
	  filter = "startswith(displayName, 'Create_')"
	}

	resource "nutanix_roles_v2" "test_1" {
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

	resource "nutanix_roles_v2" "test_2" {
		display_name = "%[1]s"
		description  = "%[2]s"
		operations = [
			data.nutanix_operations_v2.test.operations[0].ext_id,
			data.nutanix_operations_v2.test.operations[1].ext_id,
			data.nutanix_operations_v2.test.operations[2].ext_id,
			data.nutanix_operations_v2.test.operations[3].ext_id
	  	]
		depends_on = [data.nutanix_operations_v2.test, resource.nutanix_roles_v2.test_1]
	}

	`, displayName, description)
}

func testRoleResourceWithoutDisplayNameConfig(description string) string {
	return fmt.Sprintf(`

	data "nutanix_operations_v2" "test" {
	  filter = "startswith(displayName, 'Create_')"
	}

	resource "nutanix_roles_v2" "test" {
		description  = "%[1]s"
		operations = [
			data.nutanix_operations_v2.test.operations[0].ext_id,
			data.nutanix_operations_v2.test.operations[1].ext_id,
			data.nutanix_operations_v2.test.operations[2].ext_id,
			data.nutanix_operations_v2.test.operations[3].ext_id
	  	]
		depends_on = [data.nutanix_operations_v2.test]
	}`, description)
}

func TestAccV2NutanixRolesResource_ProjectAssociation(t *testing.T) {
	projectName := fmt.Sprintf("tf-role-projassoc-%d", acctest.RandInt())
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRoleProjectAssociationConfig(projectName, "", roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameRoles, "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_role_v2.test", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttr("data.nutanix_roles_v2.test", "roles.#", "1"),
					resource.TestCheckResourceAttrPair("data.nutanix_roles_v2.test", "roles.0.ext_id", resourceNameRoles, "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_roles_v2.test", "roles.0.project_ext_id", "nutanix_project_v2.test", "ext_id"),
				),
			},
			{
				Config:      testRoleProjectAssociationConfig(projectName, "00000000-0000-0000-0000-000000000000", roleDisplayName, roleDescription),
				ExpectError: regexp.MustCompile("Update of project_ext_id is not supported"),
			},
		},
	})
}

func roleProjectExtIDLine(override string) string {
	if override == "" {
		return `project_ext_id = nutanix_project_v2.test.ext_id`
	}
	return fmt.Sprintf(`project_ext_id = "%s"`, override)
}

func testRoleProjectAssociationConfig(projectName, projectExtIDOverride, displayName, description string) string {
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
		display_name = "%[2]s"
		description  = "%[3]s"
		%[4]s
		operations = [
			data.nutanix_operations_v2.test.operations[0].ext_id,
			data.nutanix_operations_v2.test.operations[1].ext_id,
			data.nutanix_operations_v2.test.operations[2].ext_id,
			data.nutanix_operations_v2.test.operations[3].ext_id
		]
		depends_on = [data.nutanix_operations_v2.test, nutanix_project_v2.test]
	}

	data "nutanix_role_v2" "test" {
		ext_id = nutanix_roles_v2.test.id
	}

	data "nutanix_roles_v2" "test" {
		filter     = "displayName eq '${nutanix_roles_v2.test.display_name}'"
		depends_on = [nutanix_roles_v2.test]
	}
	`, projectName, displayName, description, roleProjectExtIDLine(projectExtIDOverride))
}

func testRoleResourceWithoutOperationsConfig(displayName, description string) string {
	return fmt.Sprintf(`

	resource "nutanix_roles_v2" "test" {
		display_name = "%[1]s"
		description  = "%[2]s"
	}`, displayName, description)
}

func testRoleResourceIsGlobalConfig(displayName, description string) string {
	return fmt.Sprintf(`

	data "nutanix_operations_v2" "test" {
	  filter = "startswith(displayName, 'Create_')"
	}

	resource "nutanix_roles_v2" "test" {
		display_name = "%[1]s"
		description  = "%[2]s"
		is_global = true
		operations = [
			data.nutanix_operations_v2.test.operations[0].ext_id,
			data.nutanix_operations_v2.test.operations[1].ext_id,
			data.nutanix_operations_v2.test.operations[2].ext_id,
			data.nutanix_operations_v2.test.operations[3].ext_id
	  	]
		depends_on = [data.nutanix_operations_v2.test]
	}`, displayName, description)
}
