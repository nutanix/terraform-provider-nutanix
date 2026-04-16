package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraDatabasesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabasesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_databases.test", "database_instances.0.time_zone"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_databases.test", "database_instances.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_databases.test", "database_instances.0.id"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_databases.test", "database_instances.0.linked_databases.#"),
				),
			},
		},
	})
}

func TestAccEraDatabasesDataSource_ByFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabasesDataSourceConfigByFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_databases.test", "database_instances.0.time_zone"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_databases.test", "database_instances.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_databases.test", "database_instances.0.id"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_databases.test", "database_instances.0.linked_databases.#"),
				),
			},
		},
	})
}

func testAccEraDatabasesDataSourceConfig() string {
	return `
	data "nutanix_ndb_databases" "test" {}
`
}

func testAccEraDatabasesDataSourceConfigByFilters() string {
	return `
	data "nutanix_ndb_databases" "test" {
		database_type = "postgres_database"
	}
`
}
