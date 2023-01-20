package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceNameLogCatchDB = "nutanix_ndb_database_log_catchup.acctest-managed"

func TestAccEra_LogCatchUpbasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseLogCatchUpConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameLogCatchDB, "log_catchup_version", ""),
					resource.TestCheckResourceAttr(resourceNameLogCatchDB, "database_id", ""),
					resource.TestCheckResourceAttrSet(resourceNameLogCatchDB, "time_machine_id"),
				),
			},
		},
	})
}

func testAccEraDatabaseLogCatchUpConfig() string {
	return (`
		data "nutanix_ndb_time_machines" "test1" {}

		resource "nutanix_ndb_log_catchups" "name" {
			time_machine_id = data.nutanix_ndb_time_machines.test1.time_machines.0.id
		}
	`)
}
