package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEradbserversVMDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEradbserversVMDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.metadata.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.owner_id"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.vm_info.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.vm_info.0.network_info.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.vm_cluster_uuid"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.windows_db_server"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.working_directory"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.mac_addresses.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.protection_domain_id"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.ip_addresses.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_dbservers.dbservers", "dbservers.0.era_version"),
				),
			},
		},
	})
}

func testAccEradbserversVMDataSourceConfig() string {
	return `
		data "nutanix_ndb_dbservers" "dbservers"{}
		`
}
