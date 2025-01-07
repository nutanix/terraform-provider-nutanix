package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraDatabaseDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_ndb_database.test", "time_zone", "UTC"),
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
