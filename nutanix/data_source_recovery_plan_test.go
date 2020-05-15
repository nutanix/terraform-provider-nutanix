package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixRecoveryPlanDataSource_basic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-recovery-name-dou")
	description := acctest.RandomWithPrefix("test-recovery-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-recovery-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-recovery-desc-dou")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecoveryPlanDataSourceConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccRecoveryPlanDataSourceConfig(nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func testAccRecoveryPlanDataSourceConfig(name, description string) string {
	return fmt.Sprintf(`
        resource "nutanix_recovery_plan" "test" {
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
				stage_uuid = "ab788130-0820-4d07-a1b5-b0ba4d3a4254"
				delay_time_secs = 0
			}
			parameters{}
		}
		data "nutanix_recovery_plan" "test" {
			recovery_plan_id = nutanix_recovery_plan.test.id
		}

`, name, description)
}
