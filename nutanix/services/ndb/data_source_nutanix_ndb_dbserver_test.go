package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraDBServerVMDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDBServerVMDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "vm_info.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "vm_cluster_uuid"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "windows_db_server"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "working_directory"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "mac_addresses.#"),
				),
			},
		},
	})
}

func TestAccEraDBServerVMDataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDBServerVMDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "vm_info.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "vm_cluster_uuid"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "windows_db_server"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "working_directory"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbserver.dbserver", "mac_addresses.#"),
				),
			},
		},
	})
}

func testAccEraDBServerVMDataSourceConfig() string {
	return `
		data "nutanix_ndb_dbservers" "dbservers"{}

		data "nutanix_ndb_dbserver" "dbserver"{
			id = data.nutanix_ndb_dbservers.dbservers.dbservers.0.id
		}	
		`
}

func testAccEraDBServerVMDataSourceConfigByName() string {
	return `
		data "nutanix_ndb_dbservers" "dbservers"{}

		data "nutanix_ndb_dbserver" "dbserver"{
			name = data.nutanix_ndb_dbservers.dbservers.dbservers.0.name
		}	
		`
}
