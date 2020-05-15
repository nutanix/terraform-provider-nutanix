package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixRecoveryPlan_basic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-protection-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-protection-desc-dou")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRecoveryPlanConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccNutanixRecoveryPlanConfig(nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func TestAccResourceNutanixRecoveryPlan_importBasic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRecoveryPlanConfig(name, description),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckNutanixRecoveryPlanImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixRecoveryPlanImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckNutanixRecoveryPlanExists(resourceName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[*resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", *resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixRecoveryPlanDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_recovery_plan" {
			continue
		}
		for {
			_, err := conn.API.V3.GetRecoveryPlan(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}
	}
	return nil
}

func testAccNutanixRecoveryPlanConfig(name, description string) string {
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
	`, name, description)
}
