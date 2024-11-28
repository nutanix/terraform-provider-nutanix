package iamv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameUserGroups = "nutanix_user_groups_v2.test"

func TestAccNutanixUserGroupsV2Resource_LDAPUserGroup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLDAPUserGroupsResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameUserGroups, "name", testVars.Iam.UserGroups.Name),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "idp_id", testVars.Iam.Users.DirectoryServiceId),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "group_type", "LDAP"),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "distinguished_name", testVars.Iam.UserGroups.DistinguishedName),
				),
			},
			{
				Config:      testLDAPUserGroupsResourceAlreadyExistsConfig(filepath),
				ExpectError: regexp.MustCompile("Failed to create the user group as an user group already exists with same DN"),
			},
		},
	})
}

func TestAccNutanixUserGroupsV2Resource_SAMLUserGroup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSAMLUserGroupsResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameUserGroups, "name", testVars.Iam.UserGroups.SAMLName),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "idp_id", testVars.Iam.Users.IdpId),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "group_type", "SAML"),
				),
			},
			{
				Config:      testSAMLUserGroupsResourceConfig(filepath) + testSAMLAlreadyExistsUserGroupsResourceConfig(),
				ExpectError: regexp.MustCompile("Failed to create the user group as an user group already exists with same DN"),
			},
		},
	})
}

func TestAccNutanixUserGroupsV2Resource_WithNoGroupType(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testUserGroupsResourceWithoutGroupTypeConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixUserGroupsV2Resource_WithNoIdpId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testUserGroupsResourceWithoutIdpIdConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testLDAPUserGroupsResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
		user_groups = local.config.iam.user_groups
	}

	resource "nutanix_user_groups_v2" "test" {
		group_type = "LDAP"
		idp_id = local.users.directory_service_id
		name = local.user_groups.name
		distinguished_name = local.user_groups.distinguished_name
	  }`, filepath)
}

func testLDAPUserGroupsResourceAlreadyExistsConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
		user_groups = local.config.iam.user_groups
	}

	resource "nutanix_user_groups_v2" "test_2" {
		group_type = "LDAP"
		idp_id = local.users.directory_service_id
		name = local.user_groups.name
		distinguished_name = local.user_groups.distinguished_name
	  }`, filepath)
}

func testSAMLUserGroupsResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
		user_groups = local.config.iam.user_groups
	}

	resource "nutanix_user_groups_v2" "test" {
		group_type = "SAML"
		idp_id = local.users.idp_id
		name = local.user_groups.saml_name
	  }`, filepath)
}

func testSAMLAlreadyExistsUserGroupsResourceConfig() string {
	return `
	resource "nutanix_user_groups_v2" "test_1" {
		group_type = "SAML"
		idp_id = local.users.idp_id
		name = local.user_groups.saml_name
	  }`
}

func testUserGroupsResourceWithoutGroupTypeConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
		user_groups = local.config.iam.user_groups
	}

	resource "nutanix_user_groups_v2" "test" {
		idp_id = local.user_groups.idp_id
		name = local.user_groups.name
		distinguished_name = local.user_groups.distinguished_name
	  }`, filepath)
}

func testUserGroupsResourceWithoutIdpIdConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
		user_groups = local.config.iam.user_groups
	}

	resource "nutanix_user_groups_v2" "test" {
		group_type = "LDAP"
		name = local.user_groups.name
		distinguished_name = local.user_groups.distinguished_name
	  }`, filepath)
}
