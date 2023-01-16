package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceNameScaleDB = "nutanix_ndb_database_scale.acctest-managed"

func TestAccEra_Scalebasic(t *testing.T) {
	storageSize := "4"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseScaleConfig(storageSize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameScaleDB, "application_type", "postgres_database"),
					resource.TestCheckResourceAttr(resourceNameScaleDB, "data_storage_size", storageSize),
					resource.TestCheckResourceAttr(resourceNameScaleDB, "metadata.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameScaleDB, "name"),
					resource.TestCheckResourceAttrSet(resourceNameScaleDB, "description"),
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
			database_uuid = data.nutanix_ndb_databases.test.database_instances.1.id
			data_storage_size = %[1]s
	  	}
	`, size)
}
