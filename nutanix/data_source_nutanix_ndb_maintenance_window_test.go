package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraMaintenanceWindowDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraMaintenanceWindowDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_window.test", "name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_window.test", "properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_window.test", "schedule.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_window.test", "description"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_maintenance_window.test", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccEraMaintenanceWindowDataSourceConfig() string {
	return `
		data "nutanix_ndb_maintenance_windows" "window"{ }

		data "nutanix_ndb_maintenance_window" "test"{
			id  = data.nutanix_ndb_maintenance_windows.window.maintenance_windows.0.id
		}
	`
}
