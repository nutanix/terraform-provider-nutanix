package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameDirectoryService = "data.nutanix_directory_service_v2.test"

func TestAccV2NutanixDirectoryServiceDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServiceDatasourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryService, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameDirectoryService, "name", testVars.Iam.DirectoryServices.Name),
					resource.TestCheckResourceAttr(datasourceNameDirectoryService, "domain_name", testVars.Iam.DirectoryServices.DomainName),
					resource.TestCheckResourceAttr(datasourceNameDirectoryService, "directory_type", "ACTIVE_DIRECTORY"),
					resource.TestCheckResourceAttr(datasourceNameDirectoryService, "url", testVars.Iam.DirectoryServices.URL),
					resource.TestCheckResourceAttr(datasourceNameDirectoryService, "service_account.0.username", testVars.Iam.DirectoryServices.ServiceAccount.Username),
					resource.TestCheckResourceAttrSet(datasourceNameDirectoryService, "service_account.0.password"),
					resource.TestCheckResourceAttr(datasourceNameDirectoryService, "white_listed_groups.0", testVars.Iam.DirectoryServices.WhiteListedGroups[0]),
				),
			},
		},
	})
}

func testDirectoryServiceDatasourceConfig(filepath string) string {
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
	  
	data "nutanix_directory_service_v2" "test" {
		ext_id     = resource.nutanix_directory_services_v2.test.id
		depends_on = [resource.nutanix_directory_services_v2.test]
	}
	`, filepath)
}
