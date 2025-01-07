package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameScaleDB = "nutanix_ndb_database_scale.acctest-managed"

func TestAccEra_Scalebasic(t *testing.T) {
	storageSize := "4"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseScaleConfig(storageSize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameScaleDB, "application_type", "postgres_database"),
					resource.TestCheckResourceAttr(resourceNameScaleDB, "data_storage_size", storageSize),
					resource.TestCheckResourceAttrSet(resourceNameScaleDB, "name"),
					resource.TestCheckResourceAttrSet(resourceNameScaleDB, "description"),
					resource.TestCheckResourceAttrSet(resourceNameScaleDB, "time_machine.#"),
				),
			},
		},
	})
}

func testAccEraDatabaseScaleConfig(size string) string {
	return fmt.Sprintf(`
		data "nutanix_ndb_databases" "test" {
			database_type = "postgres_database"
		}

		resource "nutanix_ndb_database_scale" "acctest-managed" {
			application_type = "postgres_database"
			database_uuid = data.nutanix_ndb_databases.test.database_instances.0.id
			data_storage_size = %[1]s
	  	}
	`, size)
}
