package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraDatabaseDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_ndb_database.test", "metadata.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_database.test", "time_zone", "UTC"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_database.test", "placeholder"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_database.test", "name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_database.test", "linked_databases.#"),
				),
			},
		},
	})
}

func testAccEraDatabaseDataSourceConfig() string {
	return `
	data "nutanix_ndb_databases" "test1" {}

	data "nutanix_ndb_database" "test" {
		database_id = data.nutanix_ndb_databases.test1.database_instances.0.id
	}
`
}
