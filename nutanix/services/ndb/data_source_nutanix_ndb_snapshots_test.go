package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNDBSnapshotsName = "data.nutanix_ndb_snapshots.test"

func TestAccEraSnapshotsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSnapshotsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotsName, "snapshots.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotsName, "snapshots.0.properties.#"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotsName, "snapshots.0.snapshot_uuid"),
					resource.TestCheckResourceAttr(dataSourceNDBSnapshotsName, "snapshots.0.status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccEraSnapshotsDataSource_WithFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSnapshotsDataSourceConfigWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotsName, "snapshots.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotsName, "snapshots.0.properties.#"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotsName, "snapshots.0.snapshot_uuid"),
					resource.TestCheckResourceAttr(dataSourceNDBSnapshotsName, "snapshots.0.status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccEraSnapshotsDataSourceConfig() string {
	return `
		data "nutanix_ndb_snapshots" "test" {}
	`
}

func testAccEraSnapshotsDataSourceConfigWithFilters() string {
	return `
		data "nutanix_ndb_time_machines" "test1" {}

		data "nutanix_ndb_snapshots" "test" {
			filters{
				time_machine_id = data.nutanix_ndb_time_machines.test1.time_machines.0.id
			}
		}
	`
}
