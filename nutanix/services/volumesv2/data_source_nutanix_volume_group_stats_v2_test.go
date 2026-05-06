package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupStats = "data.nutanix_volume_group_stats_v2.test"

func TestAccV2NutanixVolumeGroupStatsDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-stats-%d", r)
	desc := "terraform test volume group stats description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupStatsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupStats, "ext_id"),
				),
			},
		},
	})
}

func testAccVolumeGroupStatsDataSourceConfig(name, desc string) string {
	return testAccVolumeGroupResourceConfig(name, desc) + `
		data "nutanix_volume_group_stats_v2" "test" {
			ext_id     = resource.nutanix_volume_group_v2.test.id
			start_time = "2024-01-01T00:00:00Z"
			end_time   = "2026-12-31T23:59:59Z"
			depends_on = [resource.nutanix_volume_group_v2.test]
		}
	`
}
