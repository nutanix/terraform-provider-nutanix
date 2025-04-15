package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRecoveryPoint = "nutanix_self_service_app_recovery_point.test"

func TestAccNutanixCalmAppCreateRecoveryPoint(t *testing.T) {
	name := testVars.SelfService.AppWithSnapshotName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppCreateRecoveryPoint(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameRecoveryPoint, "recovery_point_name", "snap1"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoint, "action_name", "Snapshot_s1"),
				),
			},
		},
	})
}

func testCalmAppCreateRecoveryPoint(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_recovery_point" "test" {
		app_name = "%[1]s"
		action_name = "Snapshot_s1"
		recovery_point_name = "snap1"
		}
`, name)
}
