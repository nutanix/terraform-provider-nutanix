package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSnapshot = "data.nutanix_self_service_app_snapshots.test"

func TestAccNutanixCalmSnapshotGetDatasource(t *testing.T) {
	appName := testVars.SelfService.AppWithSnapshotName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnapshotDataSourceConfig(appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameSnapshot, "app_name", appName),
					resource.TestCheckResourceAttr(datasourceNameSnapshot, "entities.0.action_name", "Snapshot_s1"),
					resource.TestCheckResourceAttr(datasourceNameSnapshot, "kind", "vm_recovery_group"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshot, "entities.#"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshot, "entities.0.%"),
				),
			},
		},
	})
}

func testSnapshotDataSourceConfig(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_recovery_point" "test" {
			app_name = "%[1]s"
			action_name = "Snapshot_s1"
			recovery_point_name = "snap1"
		}
			
		data "nutanix_self_service_app_snapshots" "test" {
			app_name = "%[1]s"
			length = 250
			offset = 0
			depends_on = [nutanix_self_service_app_recovery_point.test]
		}
`, name)
}
