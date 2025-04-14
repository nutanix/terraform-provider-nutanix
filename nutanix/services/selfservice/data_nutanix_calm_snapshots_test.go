package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSnapshot = "data.nutanix_calm_app_snapshots.test"

func TestCalmSnapshotGetDatasource(t *testing.T) {
	app_name := "test_terraform_snapshot_restore_app"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnapshotDataSourceConfig(app_name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSnapshot, "entities.#"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshot, "entities.0.%"),
				),
			},
		},
	})
}

func testSnapshotDataSourceConfig(name string) string {
	return fmt.Sprintf(`
		data "nutanix_calm_app_snapshots" "test" {
			app_name = "%[1]s"
			length = 250
			offset = 0
		}
`, name)
}
