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

func TestAccNutanixRecoveryPlanWithStageList_basic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-protection-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-protection-desc-dou")

	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRecoveryPlanConfigWithStageList(name, description, stageUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccNutanixRecoveryPlanConfigWithStageList(nameUpdated, descriptionUpdated, stageUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPlanWithStageListDynamic_basic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-protection-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-protection-desc-dou")

	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"
	entity := `
 entity_info_list {
	categories {
		name = "Environment"
		value = "Dev"
	}
}
`
	entityUpdated := `
 entity_info_list {
	any_entity_reference_kind = "vm"
	any_entity_reference_uuid = "2457b73a-9ace-4c92-959d-dc24e09e0846"
	any_entity_reference_name = "terratest-drrunbook-1337"
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRecoveryPlanConfigWithStageListDynamic(name, description, stageUUID, entity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccNutanixRecoveryPlanConfigWithStageListDynamic(nameUpdated, descriptionUpdated, stageUUID, entityUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
			{
				Config: testAccNutanixRecoveryPlanConfigWithStageListDynamic(name, description, stageUUID, entity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPlanWithNetwork_basic(t *testing.T) {
	t.Skip()
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-protection-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-protection-desc-dou")

	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"
	azURLSource := "c99ab7cd-9191-4fcb-8fc0-232eff76e595"
	azURLTarget := "c7926832-4976-4fe4-bead-7e508e03e3ec"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRecoveryPlanConfigWithNetwork(name, description, stageUUID, azURLSource, azURLTarget),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccNutanixRecoveryPlanConfigWithNetwork(nameUpdated, descriptionUpdated, stageUUID, azURLSource, azURLTarget),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func TestAccResourceNutanixRecoveryPlanWithStageList_importBasic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRecoveryPlanConfigWithStageList(name, description, stageUUID),
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

func testAccNutanixRecoveryPlanConfigWithStageList(name, description, stageUUID string) string {
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
	`, name, description, stageUUID)
}

func testAccNutanixRecoveryPlanConfigWithNetwork(name, description, stageUUID, aZUrlSource, aZUrlTarget string) string {
	return fmt.Sprintf(`
		resource "nutanix_recovery_plan" "test" {
			name        = "%s"
			description = "%s"
			stage_list {
				stage_work{
					recover_entities{
						entity_info_list{
							any_entity_reference_name = "yst-leap-test-vm"
							any_entity_reference_kind = "vm"
							any_entity_reference_uuid = "d0e42d78-8b0f-4a6e-9eb4-93609de2403c"
						}
					}
				}
				stage_uuid = "%s"
				delay_time_secs = 0
			}
			parameters{
				network_mapping_list{
					availability_zone_network_mapping_list{
						availability_zone_url = "%s"
						recovery_network{
							name = "%[5]s"
							subnet_list {
								gateway_ip = "10.38.2.129"
								prefix_length = 26
								external_connectivity_state = "DISABLED"
							}
						}
						test_network{
							name = "%[5]s"
							subnet_list {
								gateway_ip = "192.168.0.1"
								prefix_length = 24
								external_connectivity_state = "DISABLED"
							}
						}
					}
					availability_zone_network_mapping_list{
						availability_zone_url = "%s"
						recovery_network{
							name = "%[5]s"
							subnet_list {
								gateway_ip = "10.38.4.65"
								prefix_length = 26
								external_connectivity_state = "DISABLED"
							}
						}
						test_network{
							name = "%[5]s"
							subnet_list {
								gateway_ip = "192.168.0.1"
								prefix_length = 24
								external_connectivity_state = "DISABLED"
							}
						}
					}
				}
			}
		}
	`, name, description, stageUUID, aZUrlSource, aZUrlTarget, testVars.SubnetName)
}

func testAccNutanixRecoveryPlanConfigWithStageListDynamic(name, description, stageUUID, categories string) string {
	return fmt.Sprintf(`
		resource "nutanix_recovery_plan" "test" {
			name        = "%s"
			description = "%s"
			stage_list {
				stage_work{
					recover_entities{
						   %s
					}
				}
				stage_uuid = "%s"
				delay_time_secs = 0
			}
			parameters{}
		}
	`, name, description, categories, stageUUID)
}
