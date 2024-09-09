package iamv2_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameUsers = "nutanix_users_v2.test"

// create local Active user, and test update the username and display name
func TestAccNutanixUsersV4Resource_LocalActiveUser(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLocalActiveUserResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUsers, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUsers, "username", testVars.Iam.Users.Username),
					resource.TestCheckResourceAttr(resourceNameUsers, "display_name", testVars.Iam.Users.DisplayName),
					resource.TestCheckResourceAttr(resourceNameUsers, "user_type", "LOCAL"),
					resource.TestCheckResourceAttr(resourceNameUsers, "first_name", testVars.Iam.Users.FirstName),
					resource.TestCheckResourceAttr(resourceNameUsers, "middle_initial", testVars.Iam.Users.MiddleInitial),
					resource.TestCheckResourceAttr(resourceNameUsers, "last_name", testVars.Iam.Users.LastName),
					resource.TestCheckResourceAttr(resourceNameUsers, "email_id", testVars.Iam.Users.EmailId),
					resource.TestCheckResourceAttr(resourceNameUsers, "status", "ACTIVE"),
				),
			},
			// test update
			{
				Config: testLocalActiveUserResourceUpdateConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUsers, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUsers, "username", testVars.Iam.Users.Username),
					resource.TestCheckResourceAttr(resourceNameUsers, "display_name", fmt.Sprintf("%s_updated", testVars.Iam.Users.DisplayName)),
					resource.TestCheckResourceAttr(resourceNameUsers, "user_type", "LOCAL"),
					resource.TestCheckResourceAttr(resourceNameUsers, "first_name", fmt.Sprintf("%s_updated", testVars.Iam.Users.FirstName)),
					resource.TestCheckResourceAttr(resourceNameUsers, "middle_initial", fmt.Sprintf("%s_updated", testVars.Iam.Users.MiddleInitial)),
					resource.TestCheckResourceAttr(resourceNameUsers, "last_name", fmt.Sprintf("%s_updated", testVars.Iam.Users.LastName)),
					resource.TestCheckResourceAttr(resourceNameUsers, "email_id", fmt.Sprintf("updated_%s", testVars.Iam.Users.EmailId)),
				),
			},
		},
	})
}

// test duplicate user creation
func TestAccNutanixUsersV4Resource_AlreadyExistsUser(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLocalActiveUserResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUsers, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUsers, "username", testVars.Iam.Users.Username),
					resource.TestCheckResourceAttr(resourceNameUsers, "display_name", testVars.Iam.Users.DisplayName),
					resource.TestCheckResourceAttr(resourceNameUsers, "user_type", "LOCAL"),
					resource.TestCheckResourceAttr(resourceNameUsers, "first_name", testVars.Iam.Users.FirstName),
					resource.TestCheckResourceAttr(resourceNameUsers, "last_name", testVars.Iam.Users.LastName),
				),
			},
			{
				Config:      testLocalUserAlreadyExistsResourceConfig(filepath),
				ExpectError: regexp.MustCompile("already existing User with given username"),
			},
		},
	})
}

// create local Inactive user
func TestAccNutanixUsersV4Resource_LocalInactiveUser(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLocalInactiveUserResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUsers, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUsers, "username", testVars.Iam.Users.Username),
					resource.TestCheckResourceAttr(resourceNameUsers, "display_name", testVars.Iam.Users.DisplayName),
					resource.TestCheckResourceAttr(resourceNameUsers, "user_type", "LOCAL"),
					resource.TestCheckResourceAttr(resourceNameUsers, "first_name", testVars.Iam.Users.FirstName),
					resource.TestCheckResourceAttr(resourceNameUsers, "last_name", testVars.Iam.Users.LastName),
					resource.TestCheckResourceAttr(resourceNameUsers, "status", "INACTIVE"),
				),
			},
		},
	})
}

// create SAML user
func TestAccNutanixUsersV4Resource_SAMLUser(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSAMLUserResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUsers, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUsers, "username", testVars.Iam.Users.Username),
					resource.TestCheckResourceAttr(resourceNameUsers, "user_type", "SAML"),
					resource.TestCheckResourceAttr(resourceNameUsers, "idp_id", testVars.Iam.Users.IdpId),
				),
			},
		},
	})
}

// create LDAP user
func TestAccNutanixUsersV4Resource_LDAPUser(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLDAPUserWithMinimalConfigResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUsers, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUsers, "username", testVars.Iam.Users.DirectoryServiceUsername),
					resource.TestCheckResourceAttr(resourceNameUsers, "user_type", "LDAP"),
					resource.TestCheckResourceAttr(resourceNameUsers, "idp_id", testVars.Iam.Users.DirectoryServiceId),
				),
			},
		},
	})
}

// create local Active user, and test update the username and display name
func TestAccNutanixUsersV4Resource_DeactivateLocalUser(t *testing.T) {
	t.Skip("these test were commented since they are using different APIs")
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLocalActiveUserResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUsers, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUsers, "username", testVars.Iam.Users.Username),
					resource.TestCheckResourceAttr(resourceNameUsers, "display_name", testVars.Iam.Users.DisplayName),
					resource.TestCheckResourceAttr(resourceNameUsers, "user_type", "LOCAL"),
					resource.TestCheckResourceAttr(resourceNameUsers, "first_name", testVars.Iam.Users.FirstName),
					resource.TestCheckResourceAttr(resourceNameUsers, "last_name", testVars.Iam.Users.LastName),
					resource.TestCheckResourceAttr(resourceNameUsers, "status", "ACTIVE"),
				),
			},
			// test Deactivate User
			{
				Config: testDeactivateLocalUserResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUsers, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUsers, "username", testVars.Iam.Users.Username),
					resource.TestCheckResourceAttr(resourceNameUsers, "display_name", testVars.Iam.Users.DisplayName),
					resource.TestCheckResourceAttr(resourceNameUsers, "user_type", "LOCAL"),
					resource.TestCheckResourceAttr(resourceNameUsers, "first_name", testVars.Iam.Users.FirstName),
					resource.TestCheckResourceAttr(resourceNameUsers, "last_name", testVars.Iam.Users.LastName),
					resource.TestCheckResourceAttr(resourceNameUsers, "status", "INACTIVE"),
				),
			},
		},
	})
}

// Test missing username
func TestAccNutanixUsersV4Resource_WithNoUserName(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testUsersResourceWithoutUserNameConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

// Test missing user type
func TestAccNutanixUsersV4Resource_WithNoUserType(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testUsersResourceWithoutUserTypeConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testLocalActiveUserResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		username = local.users.username
		first_name = local.users.first_name
		middle_initial = local.users.middle_initial
		last_name = local.users.last_name
		email_id = local.users.email_id
		locale = local.users.locale
		region = local.users.region
		display_name = local.users.display_name
		password = local.users.password
		user_type = "LOCAL"
		status = "ACTIVE"  
		force_reset_password = local.users.force_reset_password  
	}`, filepath)
}

func testLocalActiveUserResourceUpdateConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		username = local.users.username
		first_name = "${local.users.first_name}_updated"
		middle_initial = "${local.users.middle_initial}_updated"
		last_name = "${local.users.last_name}_updated"
		email_id = "updated_${local.users.email_id}"
		locale = local.users.locale
		region = local.users.region
		display_name = "${local.users.display_name}_updated"
		password = "${local.users.password}_updated"
		user_type = "LOCAL"
		status = "ACTIVE"  
		force_reset_password = local.users.force_reset_password
		
	}`, filepath)
}

func testLocalUserAlreadyExistsResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test2" {
		username = local.users.username
		first_name = local.users.first_name
		middle_initial = local.users.middle_initial
		last_name = local.users.last_name
		email_id = local.users.email_id
		locale = local.users.locale
		region = local.users.region
		display_name = local.users.display_name
		password = local.users.password
		user_type = "LOCAL"
		status = "ACTIVE"  
		force_reset_password = local.users.force_reset_password
	}
		
	`, filepath)
}

func testLocalInactiveUserResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		username = local.users.username
		first_name = local.users.first_name
		middle_initial = local.users.middle_initial
		last_name = local.users.last_name
		email_id = local.users.email_id
		locale = local.users.locale
		region = local.users.region
		display_name = local.users.display_name
		password = local.users.password
		user_type = "LOCAL"
		status = "INACTIVE"  
		force_reset_password = local.users.force_reset_password

	}`, filepath)
}

func testSAMLUserResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		username = local.users.username
		user_type = "SAML"
		idp_id = local.users.idp_id		
	}`, filepath)
}

func testLDAPUserWithMinimalConfigResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		username = local.users.directory_service_username
		user_type = "LDAP"
		idp_id = local.users.directory_service_id
		
	}`, filepath)
}

func testDeactivateLocalUserResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		username = local.users.username
		user_type = "LOCAL"
		idp_id = local.users.idp_id
		display_name = local.users.display_name
		locale = local.users.locale
		region = local.users.region
		password = local.users.password
		force_reset_password = local.users.force_reset_password
		status = INACTIVE  
	}`, filepath)
}

func testUsersResourceWithoutUserNameConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		first_name = local.users.first_name
		middle_initial = local.users.middle_initial
		last_name = local.users.last_name
		email_id = local.users.email_id
		locale = local.users.locale
		region = local.users.region
		display_name = local.users.display_name
		password = local.users.password
		user_type = "LOCAL"
		status = "ACTIVE"  
		force_reset_password = local.users.force_reset_password  
		
	}`, filepath)
}

func testUsersResourceWithoutUserTypeConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		username = local.users.username
		first_name = local.users.first_name
		middle_initial = local.users.middle_initial
		last_name = local.users.last_name
		email_id = local.users.email_id
		locale = local.users.locale
		region = local.users.region
		display_name = local.users.display_name
		password = local.users.password
		status = "ACTIVE"  
		force_reset_password = local.users.force_reset_password  

	}`, filepath)
}
