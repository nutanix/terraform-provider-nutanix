package iamv2_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameDirectoryServices = "nutanix_directory_services_v2.test"

func TestAccNutanixDirectoryServicesV2Resource_CreateACTIVE_DIRECTORYService(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServices.Url),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "white_listed_groups.0", testVars.Iam.DirectoryServices.WhiteListedGroups[0]),
				),
			},
			{
				Config: testDirectoryServicesUpdateResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServices.Url),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "white_listed_groups.0", testVars.Iam.DirectoryServices.WhiteListedGroups[1]),
				),
			}},
	})
}

func TestAccNutanixDirectoryServicesV2Resource_CreateOpenLDAPService(t *testing.T) {
	t.Skip("Skipping test as OpenLDAP waiting for LDAP configuration")
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryOpenLDAPServicesResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServices.Url),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
				),
			}},
	})
}

func TestAccNutanixDirectoryServicesV2Resource_CreateACTIVE_DIRECTORYAlreadyExists(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServices.Url),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
				),
			},
			{
				Config:      testDirectoryServicesDuplicatedResourceConfig(filepath),
				ExpectError: regexp.MustCompile("Failed to create directory service as directory service with name " + testVars.Iam.DirectoryServices.Name + " already exists"),
			}},
	})
}

func TestAccNutanixDirectoryServicesV2Resource_WithNoName(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutNameConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}
func TestAccNutanixDirectoryServicesV2Resource_WithNoUrl(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutUrlConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixDirectoryServicesV2Resource_WithNoDomainName(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutDomainNameConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixDirectoryServicesV2Resource_WithNoDirectoryType(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutDirectoryTypeConfig(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixDirectoryServicesV2Resource_WithNoServiceAccount(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutServiceAccountConfig(filepath),
				ExpectError: regexp.MustCompile("Insufficient service_account blocks"),
			},
		},
	})
}

func testDirectoryServicesResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		name = local.directory_services.name
		url = local.directory_services.url  
		directory_type = "ACTIVE_DIRECTORY"
		domain_name = local.directory_services.domain_name
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		white_listed_groups = [ local.directory_services.white_listed_groups[0]]
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesUpdateResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		name = local.directory_services.name
		url = local.directory_services.url  
		directory_type = "ACTIVE_DIRECTORY"
		domain_name = local.directory_services.domain_name
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		white_listed_groups = [ local.directory_services.white_listed_groups[1]]
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryOpenLDAPServicesResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		name = local.directory_services.name
		url = local.directory_services.url  
		directory_type = "OPEN_LDAP"
		domain_name = local.directory_services.domain_name
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		open_ldap_configuration {
			user_configuration {
				user_search_base = local.directory_services.open_ldap_configuration.user_configuration.user_search_base
				username_attribute = local.directory_services.open_ldap_configuration.user_configuration.username_attribute
				user_object_class = local.directory_services.open_ldap_configuration.user_configuration.user_object_class
			}
			user_group_configuration {
				group_object_class = local.directory_services.open_ldap_configuration.user_group_configuration.group_object_class
				group_search_base = local.directory_services.open_ldap_configuration.user_group_configuration.group_search_base
				group_member_attribute = local.directory_services.open_ldap_configuration.user_group_configuration.group_member_attribute
				group_member_attribute_value = local.directory_services.open_ldap_configuration.user_group_configuration.group_member_attribute_value
			}
		}
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesDuplicatedResourceConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test_1" {
		name = local.directory_services.name
		url = local.directory_services.url  
		directory_type = "ACTIVE_DIRECTORY"
		domain_name = local.directory_services.domain_name
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesResourceWithoutNameConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		directory_type = local.directory_services.directory_type
		domain_name = local.directory_services.domain_name
		url = local.directory_services.url  
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesResourceWithoutUrlConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		name = local.directory_services.name
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		directory_type = local.directory_services.directory_type
		domain_name = local.directory_services.domain_name
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesResourceWithoutDomainNameConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		name = local.directory_services.name
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		directory_type = local.directory_services.directory_type
		url = local.directory_services.url  
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesResourceWithoutDirectoryTypeConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		name = local.directory_services.name
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		domain_name = local.directory_services.domain_name
		url = local.directory_services.url  
	    lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesResourceWithoutServiceAccountConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		name = local.directory_services.name
		
		directory_type = local.directory_services.directory_type
		domain_name = local.directory_services.domain_name
		url = local.directory_services.url  
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}

	}`, filepath)
}
