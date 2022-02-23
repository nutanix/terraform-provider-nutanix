package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRecoveryPlanDataSourceConfig_WithID(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-recovery-name-dou")
	description := acctest.RandomWithPrefix("test-recovery-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-recovery-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-recovery-desc-dou")

	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecoveryPlanDataSourceConfig_WithID(name, description, stageUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccRecoveryPlanDataSourceConfig_WithID(nameUpdated, descriptionUpdated, stageUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func TestAccRecoveryPlanDataSourceConfig_WithName(t *testing.T) {

	name := acctest.RandomWithPrefix("test-recovery-name")
	nameUpdated := acctest.RandomWithPrefix("test-recovery-name")
	description := acctest.RandomWithPrefix("test-recovery-desc")
	descriptionUpdated := acctest.RandomWithPrefix("test-recovery-desc")
	stageUUID := "bb261302-94db-11ec-b909-0242ac120002"

	resourceName := "nutanix_recovery_plan.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecoveryPlanDataSourceConfig_WithName(name, description, stageUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccRecoveryPlanDataSourceConfig_WithName(nameUpdated, descriptionUpdated, stageUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func testAccRecoveryPlanDataSourceConfig_WithID(name, description, stageUUID string) string {
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
				stage_uuid = "%s"
				delay_time_secs = 0
			}
			parameters{}
		}
		data "nutanix_recovery_plan" "test" {
			recovery_plan_id = nutanix_recovery_plan.test.id
		}

`, name, description, stageUUID)
}

func testAccRecoveryPlanDataSourceConfig_WithName(name, description, stageUUID string) string {
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
				stage_uuid = "%s"
				delay_time_secs = 0
			}
			parameters{}
		}
		data "nutanix_recovery_plan" "test" {
			recovery_plan_name = nutanix_recovery_plan.test.name
		}

`, name, description, stageUUID)
}
