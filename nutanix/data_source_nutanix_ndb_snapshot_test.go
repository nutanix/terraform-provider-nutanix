package nutanix

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceNDBSnapshotName = "data.nutanix_ndb_snapshot.test"

func TestAccEraSnapshotDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSnapshotDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "owner_id"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "properties.#"),
					resource.TestCheckResourceAttr(dataSourceNDBSnapshotName, "metadata.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "snapshot_uuid"),
					resource.TestCheckResourceAttr(dataSourceNDBSnapshotName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccEraSnapshotDataSource_WithFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSnapshotDataSourceConfigWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "owner_id"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "properties.#"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "nx_cluster_id"),
					resource.TestCheckResourceAttr(dataSourceNDBSnapshotName, "metadata.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceNDBSnapshotName, "snapshot_uuid"),
					resource.TestCheckResourceAttr(dataSourceNDBSnapshotName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccEraSnapshotDataSource_WithWrongFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEraSnapshotDataSourceConfigWithWrongFilters(),
				ExpectError: regexp.MustCompile("An internal error has occurred"),
			},
		},
	})
}

func testAccEraSnapshotDataSourceConfig() string {
	return `
		data "nutanix_ndb_snapshots" "test1" {}

		data "nutanix_ndb_snapshot" "test" {
			snapshot_id = data.nutanix_ndb_snapshots.test1.snapshots.0.id
		}
	`
}

func testAccEraSnapshotDataSourceConfigWithFilters() string {
	return `
		data "nutanix_ndb_snapshots" "test1" {}

		data "nutanix_ndb_snapshot" "test" {
			snapshot_id = data.nutanix_ndb_snapshots.test1.snapshots.0.id
			filters{
				timezone= "UTC"
			}
		}
	`
}

func testAccEraSnapshotDataSourceConfigWithWrongFilters() string {
	return `
		data "nutanix_ndb_snapshots" "test1" {}

		data "nutanix_ndb_snapshot" "test" {
			snapshot_id = data.nutanix_ndb_snapshots.test1.snapshots.0.id
			filters{
				timezone= "IST"
			}
		}
	`
}
