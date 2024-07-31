package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameSnapshotDB = "nutanix_ndb_database_snapshot.acctest-managed"

func TestAccEra_Snapshotbasic(t *testing.T) {
	name := "test-acc-snapshot"
	removalIndays := "2"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseSnapshotConfig(name, removalIndays),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSnapshotDB, "name", name),
					resource.TestCheckResourceAttr(resourceNameSnapshotDB, "remove_schedule_in_days", removalIndays),
					resource.TestCheckResourceAttr(resourceNameSnapshotDB, "database_snapshot", "false"),
				),
			},
		},
	})
}

func TestAccEra_Snapshot_ReplicateToClusters(t *testing.T) {
	name := "test-acc-snapshot"
	removalIndays := "2"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseSnapshotConfigReplicateToClusters(name, removalIndays),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSnapshotDB, "name", name),
					resource.TestCheckResourceAttr(resourceNameSnapshotDB, "remove_schedule_in_days", removalIndays),
					resource.TestCheckResourceAttr(resourceNameSnapshotDB, "database_snapshot", "false"),
					resource.TestCheckResourceAttr(resourceNameSnapshotDB, "replicate_to_clusters.#", "2"),
				),
			},
		},
	})
}

func testAccEraDatabaseSnapshotConfig(name, removalIndays string) string {
	return fmt.Sprintf(`
		data "nutanix_ndb_time_machines" "test1" {}

		data "nutanix_ndb_time_machine" "test"{
			time_machine_name = data.nutanix_ndb_time_machines.test1.time_machines.0.name
		}

		resource "nutanix_ndb_database_snapshot" "acctest-managed" {
			time_machine_id = data.nutanix_ndb_time_machine.test.id
			name = "%[1]s"
			remove_schedule_in_days = "%[2]s"
		}
	`, name, removalIndays)
}

func testAccEraDatabaseSnapshotConfigReplicateToClusters(name, removalIndays string) string {
	return fmt.Sprintf(`
		data "nutanix_ndb_time_machines" "test1" {}

		data "nutanix_ndb_time_machine" "test"{
			time_machine_name = data.nutanix_ndb_time_machines.test1.time_machines.0.name
		}

		data "nutanix_ndb_clusters" "test" { }

		resource "nutanix_ndb_database_snapshot" "acctest-managed" {
			time_machine_id = data.nutanix_ndb_time_machine.test.id
			name = "%[1]s"
			remove_schedule_in_days = "%[2]s"
			replicate_to_clusters = [
				data.nutanix_ndb_clusters.test.clusters.0.id, data.nutanix_ndb_clusters.test.clusters.1.id
			]
		}
	`, name, removalIndays)
}
