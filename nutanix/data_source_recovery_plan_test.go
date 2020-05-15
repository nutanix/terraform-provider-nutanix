package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixRecoveryPlanDataSource_basic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.recovery_plan_test"

	name := acctest.RandomWithPrefix("test-recovery-name-dou")
	description := acctest.RandomWithPrefix("test-recovery-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-recovery-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-recovery-desc-dou")
	zone := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecoveryPlanDataSourceConfig(name, description, zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccRecoveryPlanDataSourceConfig(nameUpdated, descriptionUpdated, zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func testAccRecoveryPlanDataSourceConfig(name, description, zoneUUID string) string {
	return fmt.Sprintf(`
        resource "nutanix_recovery_plan" "recovery_plan_test" {
			name        = "%s"
			description = "%s"
			stage_list {
				stage_work{
					recover_entities{
						entity_info_list{
							categories {
								name = "Environment"
								value = "Dev"
							}
						}
					}
				}
				stage_uuid = "%[3]s"
				delay_time_secs = 0
			}
		}
		data "nutanix_recovery_plan" "test" {
			recovery_plan_id = nutanix_recovery_plan.recovery_plan_test.id
		}
`, name, description, zoneUUID)
}
