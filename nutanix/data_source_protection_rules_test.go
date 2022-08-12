package nutanix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixProtectionRulesDataSource_basic(t *testing.T) {
	if os.Getenv("PROTECTION_RULES_TEST_FLAG") != "true" {
		t.Skip()
	}
	dataSourceName := "data.nutanix_protection_rules.test"
	aZUUIDSource := testVars.ProtectionPolicy.LocalAz.UUID
	clusterUUIDSource := testVars.ProtectionPolicy.LocalAz.ClusterUUID
	aZUUIDTarget := testVars.ProtectionPolicy.DestinationAz.UUID
	clusterUUIDTarget := testVars.ProtectionPolicy.DestinationAz.UUID

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProtectionRulesDataSourceConfig(aZUUIDSource, clusterUUIDSource, aZUUIDTarget, clusterUUIDTarget, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "entities.1.metadata.uuid"),
				),
			},
		},
	})
}

func testAccProtectionRulesDataSourceConfig(aZUUIDSource, clusterUUIDSource, aZUUIDTarget, clusterUUIDTarget string, snapshots int64) string {
	return fmt.Sprintf(`
		locals{
			category = "AnalyticsExclusions"
			keys = ["EfficiencyMeasurement", "AnomalyDetection"]
		}
		resource "nutanix_protection_rule" "test" {
			count = 2
			name = "test-rule-${count.index}"
			description = "test-rule-desc-${count.index}"
			ordered_availability_zone_list{
				availability_zone_url = "%s"
				cluster_uuid = "%s"
			}
			ordered_availability_zone_list{
				availability_zone_url = "%s"
				cluster_uuid = "%s"
			}

			availability_zone_connectivity_list{
				source_availability_zone_index = 0
				destination_availability_zone_index = 1
				snapshot_schedule_list{
					recovery_point_objective_secs = 3600
					snapshot_type= "CRASH_CONSISTENT"
					local_snapshot_retention_policy {
						num_snapshots = %[5]d
					}
				}
			}
			availability_zone_connectivity_list{
				source_availability_zone_index = 1
				destination_availability_zone_index = 0
				snapshot_schedule_list{
					recovery_point_objective_secs = 3600
					snapshot_type= "CRASH_CONSISTENT"
					local_snapshot_retention_policy {
						num_snapshots = %[5]d
					}
				}
			}
			category_filter {
				params {
					name = local.category
					values = [local.keys[(count.index)]]
				}
			}
		}
		data "nutanix_protection_rules" "test" {
			depends_on = [nutanix_protection_rule.test]
		}
	`, aZUUIDSource, clusterUUIDSource, aZUUIDTarget, clusterUUIDTarget, snapshots)
}
