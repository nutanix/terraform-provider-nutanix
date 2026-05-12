package volumesv2_test

import (
	"fmt"
	"testing"
	"time"

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
	endTime := time.Now().UTC().Format(time.RFC3339)
	startTime := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	return testAccVolumeGroupResourceConfig(name, desc) + fmt.Sprintf(`
		data "nutanix_volume_group_stats_v2" "test" {
			ext_id     = resource.nutanix_volume_group_v2.test.id
			start_time = "%s"
			end_time   = "%s"
			depends_on = [resource.nutanix_volume_group_v2.test]
		}
	`, startTime, endTime)
}
