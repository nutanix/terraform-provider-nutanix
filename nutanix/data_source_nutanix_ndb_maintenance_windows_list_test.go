package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraMaintenanceWindowsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraMaintenanceWindowsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_windows.test", "maintenance_windows.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_windows.test", "maintenance_windows.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_windows.test", "maintenance_windows.0.properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_windows.test", "maintenance_windows.0.schedule.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_maintenance_windows.test", "maintenance_windows.0.description"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_maintenance_windows.test", "maintenance_windows.0.status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccEraMaintenanceWindowsDataSourceConfig() string {
	return `
		data "nutanix_ndb_maintenance_windows" "test"{ }
	`
}
