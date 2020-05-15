package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixProtectionRuleDataSource_basic(t *testing.T) {
	resourceName := "nutanix_protection_rule.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-protection-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-protection-desc-dou")

	zone := "ab788130-0820-4d07-a1b5-b0ba4d3ard54"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixProtectionRUleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProtectionRuleDataSourceConfig(name, description, zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccProtectionRuleDataSourceConfig(nameUpdated, descriptionUpdated, zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func testAccProtectionRuleDataSourceConfig(name, description, zone string) string {
	return fmt.Sprintf(`
		resource "nutanix_protection_rule" "test" {
			name        = "%s"
			description = "%s"
			ordered_availability_zone_list{
				availability_zone_url = "%s"
			}

			availability_zone_connectivity_list{
				snapshot_schedule_list{
					recovery_point_objective_secs = 3600
					snapshot_type= "CRASH_CONSISTENT"
					local_snapshot_retention_policy = {
						num_snapshots = 1
					}
				}
			}
			category_filter {
				params {
					name = "Environment"
					values = ["Dev"]
				}
			}
		}
		data "nutanix_protection_rule" "test" {
			protection_rule_id = nutanix_protection_rule.test.id
		}
`, name, description, zone)
}
