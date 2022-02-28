package nutanix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixProtectionRuleDataSource_basic(t *testing.T) {
	if os.Getenv("PROTECTION_RULES_TEST_FLAG") != "true" {
		t.Skip()
	}
	dataSourceName := "data.nutanix_protection_rule.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")
	aZUUIDSource := testVars.ProtectionPolicy.LocalAz.UUID
	clusterUUIDSource := testVars.ProtectionPolicy.LocalAz.ClusterUUID
	aZUUIDTarget := testVars.ProtectionPolicy.DestinationAz.UUID
	clusterUUIDTarget := testVars.ProtectionPolicy.DestinationAz.UUID

	nameUpdated := acctest.RandomWithPrefix("test-protection-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-protection-desc-dou")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixProtectionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProtectionRuleDataSourceConfig(name, description, aZUUIDSource, clusterUUIDSource, aZUUIDTarget, clusterUUIDTarget, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "description", description),
				),
			},
			{
				Config: testAccProtectionRuleDataSourceConfig(nameUpdated, descriptionUpdated, aZUUIDSource, clusterUUIDSource, aZUUIDTarget, clusterUUIDTarget, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(dataSourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func testAccProtectionRuleDataSourceConfig(name, description, aZUUIDSource, clusterUUIDSource, aZUUIDTarget, clusterUUIDTarget string, snapshots int64) string {
	return fmt.Sprintf(`
		resource "nutanix_protection_rule" "test" {
			name        = "%s"
			description = "%s"
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
						num_snapshots = %[7]d
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
						num_snapshots = %[7]d
					}
				}
			}
			category_filter {
				params {
					name = "Environment"
					values = ["Staging"]
				}
			}
		}
		data "nutanix_protection_rule" "test" {
			protection_rule_id = nutanix_protection_rule.test.id
		}
`, name, description, aZUUIDSource, clusterUUIDSource, aZUUIDTarget, clusterUUIDTarget, snapshots)
}
