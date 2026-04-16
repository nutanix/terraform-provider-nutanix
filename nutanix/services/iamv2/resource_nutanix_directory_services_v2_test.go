package iamv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameDirectoryServices = "nutanix_directory_services_v2.test"

func TestAccV2NutanixDirectoryServicesResource_CreateACTIVE_DIRECTORYService(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixDirectoryServicesV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServices.URL),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "white_listed_groups.0", testVars.Iam.DirectoryServices.WhiteListedGroups[0]),
				),
			},
			{
				Config: testDirectoryServicesUpdateResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "group_search_type", "NON_RECURSIVE"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServices.URL),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "white_listed_groups.0", testVars.Iam.DirectoryServices.WhiteListedGroups[1]),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesResource_CreateOpenLDAPService(t *testing.T) {
	t.Skip("Skipping test as OpenLDAP waiting for LDAP configuration")

	name := fmt.Sprintf("tf-test-openldap-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixDirectoryServicesV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryOpenLDAPServicesResourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServices.URL),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesResource_CreateACTIVE_DIRECTORYAlreadyExists(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixDirectoryServicesV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServices.URL),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
				),
			},
			{
				Config:      testDirectoryServicesResourceConfig() + testDirectoryServicesDuplicatedResourceConfig(),
				ExpectError: regexp.MustCompile("Failed to create directory service as directory service with name " + testVars.Iam.DirectoryServices.Name + " already exists"),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesResource_WithNoName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutNameConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesResource_WithNoUrl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutURLConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesResource_WithNoDomainName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutDomainNameConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesResource_WithNoDirectoryType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutDirectoryTypeConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesResource_WithNoServiceAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDirectoryServicesResourceWithoutServiceAccountConfig(),
				ExpectError: regexp.MustCompile("Insufficient service_account blocks"),
			},
		},
	})
}

func testDirectoryServicesResourceConfig() string {
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

func testDirectoryServicesUpdateResourceConfig() string {
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
		group_search_type = "NON_RECURSIVE"
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

func testDirectoryOpenLDAPServicesResourceConfig(name string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%[2]s")))
		directory_services = local.config.iam.directory_services
	}

	resource "nutanix_directory_services_v2" "test" {
		name = "%[1]s"
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
	}`, name, filepath)
}

func testDirectoryServicesDuplicatedResourceConfig() string {
	return `
	resource "nutanix_directory_services_v2" "test_1" {
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
	}`
}

func testDirectoryServicesResourceWithoutNameConfig() string {
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
		directory_type = "ACTIVE_DIRECTORY"
		domain_name = local.directory_services.domain_name
		url = local.directory_services.url
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesResourceWithoutURLConfig() string {
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
		directory_type = "ACTIVE_DIRECTORY"
		domain_name = local.directory_services.domain_name
		lifecycle {
			ignore_changes = [
			  service_account.0.password,
			]
	  	}
	}`, filepath)
}

func testDirectoryServicesResourceWithoutDomainNameConfig() string {
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

func testDirectoryServicesResourceWithoutDirectoryTypeConfig() string {
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

func testDirectoryServicesResourceWithoutServiceAccountConfig() string {
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
