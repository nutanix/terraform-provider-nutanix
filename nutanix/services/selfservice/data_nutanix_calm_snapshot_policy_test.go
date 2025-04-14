package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSnapshotPolicy = "data.nutanix_calm_snapshot_policy_list.test"

func TestCalmSnapshotPolicyGetDatasource(t *testing.T) {
	bpName := "demo_bp2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnapshotPolicyDataSourceConfig(bpName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSnapshotPolicy, "policy_list.0.snapshot_config_name"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshotPolicy, "policy_list.0.snapshot_config_uuid"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshotPolicy, "policy_list.0.policy_name"),
					resource.TestCheckResourceAttrSet(datasourceNameSnapshotPolicy, "policy_list.0.policy_uuid"),
				),
			},
		},
	})
}

func testSnapshotPolicyDataSourceConfig(name string) string {
	return fmt.Sprintf(`
		data "nutanix_calm_snapshot_policy_list" "test"{
			bp_name = "%[1]s"
			length = 250
			offset = 0
		}
`, name)
}
