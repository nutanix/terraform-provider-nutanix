package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRestore = "nutanix_self_service_app_restore.test"

func TestAccNutanixCalmAppRestoreRecoveryPoint(t *testing.T) {
	name := testVars.SelfService.AppWithSnapshotName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppRestoreRecoveryPoint(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameRestore, "restore_action_name", "Restore_s1"),
					resource.TestCheckResourceAttr(resourceNameRestore, "state", "SUCCESS"),
				),
			},
		},
	})
}

func testCalmAppRestoreRecoveryPoint(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_recovery_point" "test" {
		app_name = "%[1]s"
		action_name = "Snapshot_s1"
		recovery_point_name = "snap1"
		}

		data "nutanix_self_service_app_snapshots" "snapshots" {
		app_name = "%[1]s"
		length = 250
		offset = 0
		depends_on = [nutanix_self_service_app_recovery_point.test]
		}

		locals {
			snapshot_uuid = [
			for snapshot in data.nutanix_self_service_app_snapshots.snapshots.entities :
			snapshot.uuid if snapshot.name == "Snapshot_Configs1"
			][0]
		}

		resource "nutanix_self_service_app_restore" "test" {
		restore_action_name = "Restore_s1"
		app_name = "%[1]s"
		snapshot_uuid = local.snapshot_uuid
		}
`, name)
}
