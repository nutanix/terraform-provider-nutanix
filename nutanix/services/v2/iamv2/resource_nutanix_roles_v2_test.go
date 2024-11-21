package iamv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRoles = "nutanix_roles_v2.test"

func TestAccNutanixRolesV2Resource_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRoleResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRoles, "client_name"),
					resource.TestCheckResourceAttr(resourceNameRoles, "display_name", testVars.Iam.Roles.DisplayName),
					resource.TestCheckResourceAttr(resourceNameRoles, "description", testVars.Iam.Roles.Description),
				),
			},
			// update role
			{
				Config: testRoleResourceUpdateConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRoles, "client_name"),
					resource.TestCheckResourceAttr(resourceNameRoles, "display_name", fmt.Sprintf("%s_updated", testVars.Iam.Roles.DisplayName)),
					resource.TestCheckResourceAttr(resourceNameRoles, "description", testVars.Iam.Roles.Description),
				),
			},
		},
	})
}

func TestAccNutanixRolesV2Resource_DuplicateRole(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testRoleResourceDuplicateRoleConfig(filepath),
				ExpectError: regexp.MustCompile("Failed to create role as already exists"),
			},
		},
	})
}

func TestAccNutanixRolesV2Resource_WithNoDisplayName(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testRoleResourceWithoutDisplayNameConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}
func TestAccNutanixRolesV2Resource_WithNoOperations(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testRoleResourceWithoutOperationsConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testRoleResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		roles = local.config.iam.roles
	}
	
	data "nutanix_operations_v2" "test" {
	  limit = 3
	}

	resource "nutanix_roles_v2" "test" {
		display_name = local.roles.display_name
		description  = local.roles.description
		operations = [
			data.nutanix_operations_v2.test.permissions[0].ext_id,
			data.nutanix_operations_v2.test.permissions[1].ext_id,
			data.nutanix_operations_v2.test.permissions[2].ext_id,
	  	]
		depends_on = [data.nutanix_operations_v2.test]
	}`, filepath)
}

func testRoleResourceUpdateConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		roles = local.config.iam.roles
	}
	
	data "nutanix_operations_v2" "test" {
	  limit = 3
	}

	resource "nutanix_roles_v2" "test" {
		display_name = "${local.roles.display_name}_updated"
		description  = local.roles.description
		operations = [
			data.nutanix_operations_v2.test.permissions[0].ext_id,
			data.nutanix_operations_v2.test.permissions[1].ext_id,
			data.nutanix_operations_v2.test.permissions[2].ext_id,
	  	]
		depends_on = [data.nutanix_operations_v2.test]
	}`, filepath)
}

func testRoleResourceDuplicateRoleConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		roles = local.config.iam.roles
	}
	
	data "nutanix_operations_v2" "test" {
	  limit = 3
	}

	resource "nutanix_roles_v2" "test_1" {
		display_name = local.roles.display_name
		description  = local.roles.description
		operations = [
			data.nutanix_operations_v2.test.permissions[0].ext_id,
			data.nutanix_operations_v2.test.permissions[1].ext_id,
			data.nutanix_operations_v2.test.permissions[2].ext_id,
	  	]
		depends_on = [data.nutanix_operations_v2.test]
	}
	
	resource "nutanix_roles_v2" "test_2" {
		display_name = local.roles.display_name
		description  = local.roles.description
		operations = [
			data.nutanix_operations_v2.test.permissions[0].ext_id,
			data.nutanix_operations_v2.test.permissions[1].ext_id,
			data.nutanix_operations_v2.test.permissions[2].ext_id,
	  	]
		depends_on = [data.nutanix_operations_v2.test, resource.nutanix_roles_v2.test_1]
	}
	
	`, filepath)
}

func testRoleResourceWithoutDisplayNameConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		roles = local.config.iam.roles
	}
	
	data "nutanix_operations_v2" "test" {
	  limit = 3
	}

	resource "nutanix_roles_v2" "test" {
		description  = local.roles.description
		operations = [
			data.nutanix_operations_v2.test.permissions[0].ext_id,
			data.nutanix_operations_v2.test.permissions[1].ext_id,
			data.nutanix_operations_v2.test.permissions[2].ext_id,
	  	]
		depends_on = [data.nutanix_operations_v2.test]
	}`, filepath)
}

func testRoleResourceWithoutOperationsConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		roles = local.config.iam.roles
	}

	resource "nutanix_roles_v2" "test" {
		display_name = local.roles.display_name
		description  = local.roles.description
	}`, filepath)
}
