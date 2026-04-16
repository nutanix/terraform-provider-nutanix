package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameDirectoryServices = "data.nutanix_directory_services_v2.test"

func TestAccV2NutanixDirectoryServicesDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesDatasourceV2Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryServices, "directory_services.#"),
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryServices, "directory_services.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryServices, "directory_services.0.url"),
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryServices, "directory_services.0.domain_name"),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesDatasource_WithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesDatasourceV2WithFilterConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryServices, "directory_services.#"),
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryServices, "directory_services.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDirectoryServices, "directory_services.0.name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(datasourceNameDirectoryServices, "directory_services.0.domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(datasourceNameDirectoryServices, "directory_services.0.directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(datasourceNameDirectoryServices, "directory_services.0.url", testVars.Iam.DirectoryServices.URL),
					resource.TestCheckResourceAttr(datasourceNameDirectoryServices, "directory_services.0.service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryServices, "directory_services.0.service_account.0.password"),
					resource.TestCheckResourceAttr(datasourceNameDirectoryServices, "directory_services.0.white_listed_groups.0", testVars.Iam.DirectoryServices.WhiteListedGroups[0]),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesDatasource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesDatasourceV2WithLimitConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryServices, "directory_services.#"),
					resource.TestCheckResourceAttr(datasourceNameDirectoryServices, "directory_services.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServicesDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServicesDatasourceV2WithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameDirectoryServices, "directory_services.#", "0"),
				),
			},
		},
	})
}

func testDirectoryServicesDatasourceV2Config(filepath string) string {
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
		}
		data "nutanix_directory_services_v2" "test"{
			depends_on = [resource.nutanix_directory_services_v2.test]
		}
	`, filepath)
}

func testDirectoryServicesDatasourceV2WithFilterConfig(filepath string) string {
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
	}

	data "nutanix_directory_services_v2" "test" {
		filter     = "name eq '${resource.nutanix_directory_services_v2.test.name}'"
		depends_on = [resource.nutanix_directory_services_v2.test]
	}
	`, filepath)
}

func testDirectoryServicesDatasourceV2WithLimitConfig(filepath string) string {
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
		}

		data "nutanix_directory_services_v2" "test" {
			limit     = 1
			depends_on = [resource.nutanix_directory_services_v2.test]
		}
	`, filepath)
}

func testDirectoryServicesDatasourceV2WithInvalidFilterConfig() string {
	return `
	data "nutanix_directory_services_v2" "test" {
		filter = "name eq 'invalid_filter'"
	}
	`
}
