package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSnapshotPolicy = "data.nutanix_self_service_snapshot_policy_list.test"

func TestAccNutanixCalmSnapshotPolicyGetDatasource(t *testing.T) {
	blueprintName := testVars.SelfService.BlueprintWithSnapshotName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnapshotPolicyDataSourceConfig(blueprintName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameSnapshotPolicy, "bp_name", blueprintName),
					resource.TestCheckResourceAttr(datasourceNameSnapshotPolicy, "policy_list.0.snapshot_config_name", "Snapshot_Configs1"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshotPolicy, "policy_list.0.snapshot_config_uuid"),
					resource.TestCheckResourceAttr(datasourceNameSnapshotPolicy, "policy_list.0.policy_name", "test_local_snapshot_policy_local_account"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshotPolicy, "policy_list.0.policy_uuid"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshotPolicy, "policy_list.0.snapshot_config_uuid"),
				),
			},
		},
	})
}

func testSnapshotPolicyDataSourceConfig(name string) string {
	return fmt.Sprintf(`
		data "nutanix_self_service_snapshot_policy_list" "test"{
			bp_name = "%[1]s"
			length = 250
			offset = 0
		}
`, name)
}
