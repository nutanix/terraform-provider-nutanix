package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixProtectionRule_basic(t *testing.T) {
	resourceName := "nutanix_protection_rule.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")
	aZUrlSource := "c99ab7cd-9191-4fcb-8fc0-232eff76e595"
	uuidSource := "4db9adc1-8d13-4585-a901-a3ce1276ecb0"
	aZUrlTarget := "45a97947-4b09-4179-8e9b-0c2859020539"
	uuidTarget := "40cc9ba1-4c3c-4deb-a04e-a5e33c09d767"

	nameUpdated := acctest.RandomWithPrefix("test-protection-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-protection-desc-dou")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixProtectionRUleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProtectionRuleConfig(name, description, aZUrlSource, uuidSource, aZUrlTarget, uuidTarget, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProtectionRuleExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccNutanixProtectionRuleConfig(nameUpdated, descriptionUpdated, aZUrlSource, uuidSource, aZUrlTarget, uuidTarget, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProtectionRuleExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func TestAccResourceNutanixProtectionRule_importBasic(t *testing.T) {
	resourceName := "nutanix_protection_rule.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")
	aZUrlSource := "c99ab7cd-9191-4fcb-8fc0-232eff76e595"
	uuidSource := "4db9adc1-8d13-4585-a901-a3ce1276ecb0"
	aZUrlTarget := "45a97947-4b09-4179-8e9b-0c2859020539"
	uuidTarget := "40cc9ba1-4c3c-4deb-a04e-a5e33c09d767"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixProtectionRUleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProtectionRuleConfig(name, description, aZUrlSource, uuidSource, aZUrlTarget, uuidTarget, 1),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckNutanixProtectionRuleImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixProtectionRuleImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckNutanixProtectionRuleExists(resourceName *string) resource.TestCheckFunc {
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

func testAccCheckNutanixProtectionRUleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_protection_rule" {
			continue
		}
		for {
			_, err := conn.API.V3.GetProtectionRule(rs.Primary.ID)
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

func testAccNutanixProtectionRuleConfig(name, description, aZUrlSource, clusterUUIDSource, aZUrlTarget, clusterUUIDTarget string, snapshots int64) string {
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
					values = ["Dev"]
				}
			}
		}
	`, name, description, aZUrlSource, clusterUUIDSource, aZUrlTarget, clusterUUIDTarget, snapshots)
}
