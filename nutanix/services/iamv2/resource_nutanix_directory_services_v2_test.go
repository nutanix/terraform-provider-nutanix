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
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServicesMain.PrimaryAD.DirectoryServiceName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServicesMain.PrimaryAD.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServicesMain.PrimaryAD.URL),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServicesMain.PrimaryAD.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "white_listed_groups.0", "test"),
				),
			},
			{
				Config: testDirectoryServicesUpdateResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServicesMain.PrimaryAD.DirectoryServiceName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServicesMain.PrimaryAD.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "group_search_type", "NON_RECURSIVE"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServicesMain.PrimaryAD.URL),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServicesMain.PrimaryAD.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "white_listed_groups.0", "test_updated"),
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
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServicesMain.PrimaryAD.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServicesMain.PrimaryAD.URL),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServicesMain.PrimaryAD.Username),
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
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", testVars.Iam.DirectoryServicesMain.PrimaryAD.DirectoryServiceName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "domain_name", testVars.Iam.DirectoryServicesMain.PrimaryAD.Name),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "url", testVars.Iam.DirectoryServicesMain.PrimaryAD.URL),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "service_account.0.username", testVars.Iam.DirectoryServicesMain.PrimaryAD.Username),
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "service_account.0.password"),
				),
			},
			{
				Config:      testDirectoryServicesResourceConfig() + testDirectoryServicesDuplicatedResourceConfig(),
				ExpectError: regexp.MustCompile("Failed to create directory service as directory service with name " + testVars.Iam.DirectoryServicesMain.PrimaryAD.DirectoryServiceName + " already exists"),
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
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
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
		white_listed_groups = [ "test"]
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
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
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
		white_listed_groups = [ "test_updated"]
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
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
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
				user_search_base = "dc=iamopenldap,dc=com"
				username_attribute = "uid"
				user_object_class = "posixAccount"
			}
			user_group_configuration {
				group_object_class = "posixGroup"
				group_search_base = "ou=Group,dc=iamopenldap,dc=com"
				group_member_attribute = "memberUid"
				group_member_attribute_value = "uid"
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
		white_listed_groups = [ "test"]
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
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
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
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
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
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
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
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
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
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
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

const datasourceNameDSProject = "data.nutanix_directory_service_v2.project_test"
const datasourceNameDSProjectList = "data.nutanix_directory_services_v2.list_test"
const defaultProjectUUID = "00000000-0000-0000-0000-000000000000"

func TestAccV2NutanixDirectoryServicesResource_DefaultProjectAndSharing(t *testing.T) {
	r := acctest.RandInt()
	projectName := fmt.Sprintf("tf-ds-share-proj-%d", r)
	dsName := fmt.Sprintf("tf-ds-share-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixDirectoryServicesV2Destroy,
		Steps: []resource.TestStep{
			// Step 1: Create with shared_with_projects = [project_id]
			{
				Config: testDirectoryServicesProjectConfig(dsName, projectName, "RECURSIVE", "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", dsName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameDirectoryServices, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameDSProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProjectList, "directory_services.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameDSProjectList, "directory_services.0.ext_id", resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProjectList, "directory_services.0.shared_with_projects.#", "1"),
				),
			},
			// Step 2: Share all
			{
				Config: testDirectoryServicesProjectConfig(dsName, projectName, "NON_RECURSIVE", "share_all"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "group_search_type", "NON_RECURSIVE"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "share_with_all_projects", "true"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "share_with_all_projects", "true"),
				),
			},
			// Step 3: Unshare all + Share with specific project in one step
			{
				Config: testDirectoryServicesProjectConfig(dsName, projectName, "RECURSIVE", "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "group_search_type", "RECURSIVE"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameDirectoryServices, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameDSProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProjectList, "directory_services.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameDSProjectList, "directory_services.0.ext_id", resourceNameDirectoryServices, "ext_id"),
				),
			},
			// Step 4: Unshare specific project
			{
				Config: testDirectoryServicesProjectConfig(dsName, projectName, "RECURSIVE", "empty"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "group_search_type", "RECURSIVE"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "shared_with_projects.#", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesResource_NonDefaultProjectSharingFails(t *testing.T) {
	r := acctest.RandInt()
	project1Name := fmt.Sprintf("tf-ds-proj1-%d", r)
	project2Name := fmt.Sprintf("tf-ds-proj2-%d", r)
	dsName := fmt.Sprintf("tf-ds-nondefault-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixDirectoryServicesV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesNonDefaultProjectConfig(dsName, project1Name, project2Name, "none"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameDirectoryServices, "project_ext_id", "nutanix_project_v2.project1", "ext_id"),
				),
			},
			{
				Config:      testDirectoryServicesNonDefaultProjectConfig(dsName, project1Name, project2Name, "share"),
				ExpectError: regexp.MustCompile("directory service belonging to non-default project is not allowed to be shared/unshared"),
			},
		},
	})
}

func testDirectoryServicesProjectConfig(dsName, projectName, groupSearchType, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `
		share_with_all_projects = false
		shared_with_projects    = [nutanix_project_v2.share_project.ext_id]`
	case "empty":
		shareBlock = `
		share_with_all_projects = false
		shared_with_projects    = []`
	case "share_all":
		shareBlock = `share_with_all_projects = true`
	case "unshare_all":
		shareBlock = `share_with_all_projects = false`
	}
	listDataSource := ""
	if shareState == "share" {
		listDataSource = `
	data "nutanix_directory_services_v2" "list_test" {
		filter     = "sharedWithProjects/any(p:p eq '${nutanix_project_v2.share_project.id}')"
		depends_on = [nutanix_directory_services_v2.test]
	}`
	}
	return fmt.Sprintf(`
	locals {
		config             = jsondecode(file("%[1]s"))
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
	}

	resource "nutanix_project_v2" "share_project" {
		name       = "%[2]s"
		project_id = "%[2]s"
		description = "project for directory service sharing test"
	}

	resource "nutanix_directory_services_v2" "test" {
		name           = "%[6]s"
		url            = local.directory_services.url
		directory_type = "ACTIVE_DIRECTORY"
		domain_name    = local.directory_services.domain_name
		group_search_type = "%[3]s"
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		white_listed_groups = ["test"]
		%[4]s
		lifecycle {
			ignore_changes = [service_account.0.password]
		}
		depends_on = [nutanix_project_v2.share_project]
	}

	data "nutanix_directory_service_v2" "project_test" {
		ext_id = nutanix_directory_services_v2.test.id
	}
	%[5]s
`, filepath, projectName, groupSearchType, shareBlock, listDataSource, dsName)
}

func testDirectoryServicesNonDefaultProjectConfig(dsName, proj1, proj2, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `shared_with_projects = [nutanix_project_v2.project2.ext_id]`
	case "empty":
		shareBlock = `shared_with_projects = []`
	}
	return fmt.Sprintf(`
	locals {
		config             = jsondecode(file("%[1]s"))
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
	}

	resource "nutanix_project_v2" "project1" {
		name       = "%[2]s"
		project_id = "%[2]s"
		description = "first project for directory service test"
	}

	resource "nutanix_project_v2" "project2" {
		name       = "%[3]s"
		project_id = "%[3]s"
		description = "second project for directory service test"
	}

	resource "nutanix_directory_services_v2" "test" {
		name           = "%[5]s"
		url            = local.directory_services.url
		directory_type = "ACTIVE_DIRECTORY"
		domain_name    = local.directory_services.domain_name
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		white_listed_groups = ["test"]
		project_ext_id = nutanix_project_v2.project1.ext_id
		%[4]s
		lifecycle {
			ignore_changes = [service_account.0.password]
		}
		depends_on = [nutanix_project_v2.project1, nutanix_project_v2.project2]
	}
`, filepath, proj1, proj2, shareBlock, dsName)
}

func TestAccV2NutanixDirectoryServicesResource_ShareWithAllProjects(t *testing.T) {
	r := acctest.RandInt()
	projectName := fmt.Sprintf("tf-ds-shareall-proj-%d", r)
	dsName := fmt.Sprintf("tf-ds-shareall-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixDirectoryServicesV2Destroy,
		Steps: []resource.TestStep{
			// Step 1: Create with share_with_all_projects = true
			{
				Config: testDirectoryServicesShareAllConfig(dsName, projectName, "RECURSIVE", "share_all"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "name", dsName),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "share_with_all_projects", "true"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "share_with_all_projects", "true"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameDSProjectList, "directory_services.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameDSProjectList, "directory_services.0.ext_id", resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProjectList, "directory_services.0.share_with_all_projects", "true"),
				),
			},
			// Step 2: Unshare all + Share with specific project in one step
			{
				Config: testDirectoryServicesShareAllConfig(dsName, projectName, "NON_RECURSIVE", "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "group_search_type", "NON_RECURSIVE"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServices, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameDirectoryServices, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameDSProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameDSProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProjectList, "directory_services.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameDSProjectList, "directory_services.0.ext_id", resourceNameDirectoryServices, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDSProjectList, "directory_services.0.share_with_all_projects", "false"),
				),
			},
		},
	})
}

const datasourceNameDSSearch = "data.nutanix_directory_service_users_search_v2.search_test"

// ---------------------------------------------------------------------------
// User search — wildcard search for "administrator" (known to exist in AD)
// Validates: directory_service_ext_id, query, is_wildcard_search,
//            domain_name, search_results.#, search_results.0.*
// ---------------------------------------------------------------------------

func TestAccV2NutanixDirectoryServiceSearch_UserSearch(t *testing.T) {
	dsExtID := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSSearchConfig_Basic(dsExtID, "ssptest1@qa.nucalm.io", true),
				Check: resource.ComposeTestCheckFunc(
					// --- input attributes stored correctly ---
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "directory_service_ext_id", dsExtID),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "query", "ssptest1@qa.nucalm.io"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "is_wildcard_search", "true"),

					// --- computed attributes ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "domain_name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.#"),

					// --- first result entity ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.entity_type"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.#"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.values.#"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// User group search — wildcard search for "dnd_approval_group_3"
// Validates: all basic output fields for a group entity
// ---------------------------------------------------------------------------

func TestAccV2NutanixDirectoryServiceSearch_UserGroupSearch(t *testing.T) {
	dsExtID := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSSearchConfig_Basic(dsExtID, "dnd_approval_group_3", true),
				Check: resource.ComposeTestCheckFunc(
					// --- input attributes ---
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "directory_service_ext_id", dsExtID),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "query", "dnd_approval_group_3"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "is_wildcard_search", "true"),

					// --- computed attributes ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "domain_name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.#"),

					// --- first result entity ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.entity_type"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.#"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.values.#"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// returned_attributes — search for "dnd_approval_group_4" and request
// specific AD attributes back (cn, memberOf, objectGUID)
// ---------------------------------------------------------------------------

func TestAccV2NutanixDirectoryServiceSearch_ReturnedAttributes(t *testing.T) {
	dsExtID := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSSearchConfig_ReturnedAttrs(dsExtID, "dnd_approval_group_4"),
				Check: resource.ComposeTestCheckFunc(
					// --- input attributes ---
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "directory_service_ext_id", dsExtID),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "query", "dnd_approval_group_4"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "is_wildcard_search", "true"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "returned_attributes.#", "3"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "returned_attributes.0", "cn"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "returned_attributes.1", "memberOf"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "returned_attributes.2", "objectGUID"),

					// --- computed attributes ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "domain_name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.#"),

					// --- first result entity ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.entity_type"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.#"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.values.#"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// searched_attributes — search "administrator" against the "cn" attribute
// ---------------------------------------------------------------------------

func TestAccV2NutanixDirectoryServiceSearch_SearchedAttributes(t *testing.T) {
	dsExtID := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSSearchConfig_SearchedAttrs(dsExtID, "ssptest1"),
				Check: resource.ComposeTestCheckFunc(
					// --- input attributes ---
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "directory_service_ext_id", dsExtID),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "query", "ssptest1"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "is_wildcard_search", "true"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "searched_attributes.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "searched_attributes.0", "cn"),

					// --- computed attributes ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "domain_name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.#"),

					// --- first result entity ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.entity_type"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.#"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.values.#"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Combined: searched_attributes + returned_attributes together
// Search "dnd_approval_group_3" in cn/sAMAccountName, return cn + objectGUID
// ---------------------------------------------------------------------------

func TestAccV2NutanixDirectoryServiceSearch_SearchedAndReturnedAttributes(t *testing.T) {
	dsExtID := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSSearchConfig_SearchedAndReturnedAttrs(dsExtID, "dnd_approval_group_3"),
				Check: resource.ComposeTestCheckFunc(
					// --- input attributes ---
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "directory_service_ext_id", dsExtID),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "query", "dnd_approval_group_3"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "is_wildcard_search", "true"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "searched_attributes.#", "2"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "searched_attributes.0", "cn"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "searched_attributes.1", "sAMAccountName"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "returned_attributes.#", "2"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "returned_attributes.0", "cn"),
					resource.TestCheckResourceAttr(datasourceNameDSSearch, "returned_attributes.1", "objectGUID"),

					// --- computed attributes ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "domain_name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.#"),

					// --- first result entity ---
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.entity_type"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.#"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDSSearch, "search_results.0.attributes.0.values.#"),
				),
			},
		},
	})
}

// ===========================================================================
// Config generators — Directory Service Users Search
// ===========================================================================

func testDSSearchConfig_Basic(dsExtID, query string, wildcard bool) string {
	return fmt.Sprintf(`
data "nutanix_directory_service_users_search_v2" "search_test" {
  directory_service_ext_id = "%[1]s"
  query                    = "%[2]s"
  is_wildcard_search       = %[3]t
}
`, dsExtID, query, wildcard)
}

func testDSSearchConfig_ReturnedAttrs(dsExtID, query string) string {
	return fmt.Sprintf(`
data "nutanix_directory_service_users_search_v2" "search_test" {
  directory_service_ext_id = "%[1]s"
  query                    = "%[2]s"
  is_wildcard_search       = true
  returned_attributes      = ["cn", "memberOf", "objectGUID"]
}
`, dsExtID, query)
}

func testDSSearchConfig_SearchedAttrs(dsExtID, query string) string {
	return fmt.Sprintf(`
data "nutanix_directory_service_users_search_v2" "search_test" {
  directory_service_ext_id = "%[1]s"
  query                    = "%[2]s"
  is_wildcard_search       = true
  searched_attributes      = ["cn"]
}
`, dsExtID, query)
}

func testDSSearchConfig_SearchedAndReturnedAttrs(dsExtID, query string) string {
	return fmt.Sprintf(`
data "nutanix_directory_service_users_search_v2" "search_test" {
  directory_service_ext_id = "%[1]s"
  query                    = "%[2]s"
  is_wildcard_search       = true
  searched_attributes      = ["cn", "sAMAccountName"]
  returned_attributes      = ["cn", "objectGUID"]
}
`, dsExtID, query)
}

func testDirectoryServicesShareAllConfig(dsName, projectName, groupSearchType, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `
		share_with_all_projects = false
		shared_with_projects    = [nutanix_project_v2.share_project.ext_id]`
	case "share_all":
		shareBlock = `share_with_all_projects = true`
	case "unshare_all":
		shareBlock = `share_with_all_projects = false`
	}
	return fmt.Sprintf(`
	locals {
		config             = jsondecode(file("%[1]s"))
		directory_services = {
			name            = local.config.iam.directory_services_main.primary_ad.directory_service_name
			domain_name     = local.config.iam.directory_services_main.primary_ad.name
			url             = local.config.iam.directory_services_main.primary_ad.url
			directory_type  = "ACTIVE_DIRECTORY"
			service_account = {
				username = local.config.iam.directory_services_main.primary_ad.username
				password = local.config.iam.directory_services_main.primary_ad.password
			}
		}
	}

	resource "nutanix_project_v2" "share_project" {
		name       = "%[2]s"
		project_id = "%[2]s"
		description = "project for directory service share-all test"
	}

	resource "nutanix_directory_services_v2" "test" {
		name           = "%[5]s"
		url            = local.directory_services.url
		directory_type = "ACTIVE_DIRECTORY"
		domain_name    = local.directory_services.domain_name
		group_search_type = "%[3]s"
		service_account {
			username = local.directory_services.service_account.username
			password = local.directory_services.service_account.password
		}
		white_listed_groups = ["test"]
		%[4]s
		lifecycle {
			ignore_changes = [service_account.0.password]
		}
		depends_on = [nutanix_project_v2.share_project]
	}

	data "nutanix_directory_service_v2" "project_test" {
		ext_id = nutanix_directory_services_v2.test.id
	}

	data "nutanix_directory_services_v2" "list_test" {
		filter     = "name eq '${nutanix_directory_services_v2.test.name}'"
		depends_on = [nutanix_directory_services_v2.test]
	}
`, filepath, projectName, groupSearchType, shareBlock, dsName)
}
