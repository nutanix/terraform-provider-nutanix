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
