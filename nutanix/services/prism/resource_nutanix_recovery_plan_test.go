package prism_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixRecoveryPlanWithStageList_basic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	nameUpdated := acctest.RandomWithPrefix("test-protection-name-dou")
	descriptionUpdated := acctest.RandomWithPrefix("test-protection-desc-dou")

	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
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
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
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
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
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

// TestAccNutanixRecoveryPlan_withAvailabilityZoneList creates a recovery plan
// with a fully populated parameters block (availability_zone_list,
// primary_location_index and network_mapping_list) and verifies the values are
// persisted and read back into state.
func TestAccNutanixRecoveryPlan_withAvailabilityZoneList(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	azSource := testVars.ProtectionPolicy.LocalAz.UUID
	azTarget := testVars.ProtectionPolicy.DestinationAz.UUID

	name := acctest.RandomWithPrefix("test-rp-name")
	description := acctest.RandomWithPrefix("test-rp-desc")
	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			if azSource == "" || azTarget == "" {
				t.Skip("skipping: protection_policy availability zones are not configured in test_config.json")
			}
		},
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRecoveryPlanConfigWithAZList(name, description, stageUUID, azSource, azTarget, testVars.SubnetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.primary_location_index", "0"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.availability_zone_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.availability_zone_list.0.availability_zone_url", azSource),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.availability_zone_list.1.availability_zone_url", azTarget),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.network_mapping_list.#", "1"),
				),
			},
		},
	})
}

// TestAccNutanixRecoveryPlan_multiStageNetworkMapping creates a recovery plan
// with multiple stages plus a network mapping across two availability zones, and
// asserts that both stages and the full network mapping round-trip through state.
func TestAccNutanixRecoveryPlan_multiStageNetworkMapping(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	azSource := testVars.ProtectionPolicy.LocalAz.UUID
	azTarget := testVars.ProtectionPolicy.DestinationAz.UUID

	name := acctest.RandomWithPrefix("test-rp-name")
	description := acctest.RandomWithPrefix("test-rp-desc")
	stageUUID1 := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"
	stageUUID2 := "bc899241-1931-5e18-b2c6-c1cb5e4b5365"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			if azSource == "" || azTarget == "" || testVars.SubnetName == "" {
				t.Skip("skipping: protection_policy availability zones / subnet_name are not configured in test_config.json")
			}
		},
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRecoveryPlanConfigMultiStageNetworkMapping(name, description, stageUUID1, stageUUID2, azSource, azTarget, testVars.SubnetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRecoveryPlanExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "stage_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.network_mapping_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.network_mapping_list.0.availability_zone_network_mapping_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.network_mapping_list.0.availability_zone_network_mapping_list.0.recovery_network.0.name", testVars.SubnetName),
				),
			},
		},
	})
}

// TestAccNutanixRecoveryPlan_negativeMissingGatewayIP is a negative test: a
// subnet_list without the required gateway_ip must fail validation before any
// API call is made.
func TestAccNutanixRecoveryPlan_negativeMissingGatewayIP(t *testing.T) {
	name := acctest.RandomWithPrefix("test-rp-name")
	description := acctest.RandomWithPrefix("test-rp-desc")
	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNutanixRecoveryPlanConfigMissingGatewayIP(name, description, stageUUID),
				ExpectError: regexp.MustCompile(`The argument "gateway_ip" is required`),
			},
		},
	})
}

// TestAccNutanixRecoveryPlan_negativeMissingIPAddress is a negative test: an
// ip_config_list entry without the required ip_address must fail validation.
func TestAccNutanixRecoveryPlan_negativeMissingIPAddress(t *testing.T) {
	name := acctest.RandomWithPrefix("test-rp-name")
	description := acctest.RandomWithPrefix("test-rp-desc")
	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixRecoveryPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNutanixRecoveryPlanConfigMissingIPAddress(name, description, stageUUID),
				ExpectError: regexp.MustCompile(`The argument "ip_address" is required`),
			},
		},
	})
}

func TestAccNutanixRecoveryPlanWithStageList_importBasic(t *testing.T) {
	resourceName := "nutanix_recovery_plan.test"

	name := acctest.RandomWithPrefix("test-protection-name-dou")
	description := acctest.RandomWithPrefix("test-protection-desc-dou")

	stageUUID := "ab788130-0820-4d07-a1b5-b0ba4d3a4254"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
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
	conn := acc.TestAccProvider.Meta().(*conns.Client)

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

// testAccNutanixRecoveryPlanConfigWithAZList builds a cross-AZ recovery plan that
// populates the parameters block (availability_zone_list, primary_location_index
// and network_mapping_list).
//
// Note: cluster_reference_list is intentionally NOT set here. For a cross-AZ
// recovery plan the API rejects cluster information with
// "Specifying cluster information for across Availability Zones is not supported".
// cluster_reference_list is only valid for same-AZ (self-AZ) recovery plans.
func testAccNutanixRecoveryPlanConfigWithAZList(name, description, stageUUID, azSource, azTarget, subnetName string) string {
	return fmt.Sprintf(`
		resource "nutanix_recovery_plan" "test" {
			name        = "%[1]s"
			description = "%[2]s"
			stage_list {
				stage_work {
					recover_entities {
						entity_info_list {
							categories {
								name  = "Environment"
								value = "Dev"
							}
						}
					}
				}
				stage_uuid      = "%[3]s"
				delay_time_secs = 0
			}
			parameters {
				primary_location_index = 0
				availability_zone_list {
					availability_zone_url = "%[4]s"
				}
				availability_zone_list {
					availability_zone_url = "%[5]s"
				}
				network_mapping_list {
					availability_zone_network_mapping_list {
						availability_zone_url = "%[4]s"
						recovery_network {
							name = "%[6]s"
							subnet_list {
								gateway_ip    = "10.38.2.129"
								prefix_length = 26
							}
						}
						test_network {
							name = "%[6]s"
							subnet_list {
								gateway_ip    = "10.38.2.129"
								prefix_length = 26
							}
						}
					}
					availability_zone_network_mapping_list {
						availability_zone_url = "%[5]s"
						recovery_network {
							name = "%[6]s"
							subnet_list {
								gateway_ip    = "10.38.4.65"
								prefix_length = 26
							}
						}
						test_network {
							name = "%[6]s"
							subnet_list {
								gateway_ip    = "10.38.4.65"
								prefix_length = 26
							}
						}
					}
				}
			}
		}
	`, name, description, stageUUID, azSource, azTarget, subnetName)
}

// testAccNutanixRecoveryPlanConfigMultiStageNetworkMapping builds a recovery
// plan with two stages and a network mapping across two availability zones.
func testAccNutanixRecoveryPlanConfigMultiStageNetworkMapping(name, description, stageUUID1, stageUUID2, azSource, azTarget, subnetName string) string {
	return fmt.Sprintf(`
		resource "nutanix_recovery_plan" "test" {
			name        = "%[1]s"
			description = "%[2]s"
			stage_list {
				stage_work {
					recover_entities {
						entity_info_list {
							categories {
								name  = "Environment"
								value = "Dev"
							}
						}
					}
				}
				stage_uuid      = "%[3]s"
				delay_time_secs = 0
			}
			stage_list {
				stage_work {
					recover_entities {
						entity_info_list {
							categories {
								name  = "Environment"
								value = "Production"
							}
						}
					}
				}
				stage_uuid      = "%[4]s"
				delay_time_secs = 0
			}
			parameters {
				primary_location_index = 0
				availability_zone_list {
					availability_zone_url = "%[5]s"
				}
				availability_zone_list {
					availability_zone_url = "%[6]s"
				}
				network_mapping_list {
					availability_zone_network_mapping_list {
						availability_zone_url = "%[5]s"
						recovery_network {
							name = "%[7]s"
							subnet_list {
								gateway_ip    = "10.38.2.129"
								prefix_length = 26
							}
						}
						test_network {
							name = "%[7]s"
							subnet_list {
								gateway_ip    = "10.38.2.129"
								prefix_length = 26
							}
						}
					}
					availability_zone_network_mapping_list {
						availability_zone_url = "%[6]s"
						recovery_network {
							name = "%[7]s"
							subnet_list {
								gateway_ip    = "10.38.4.65"
								prefix_length = 26
							}
						}
						test_network {
							name = "%[7]s"
							subnet_list {
								gateway_ip    = "10.38.4.65"
								prefix_length = 26
							}
						}
					}
				}
			}
		}
	`, name, description, stageUUID1, stageUUID2, azSource, azTarget, subnetName)
}

// testAccNutanixRecoveryPlanConfigMissingGatewayIP omits the required gateway_ip
// inside a subnet_list; used by the negative validation test.
func testAccNutanixRecoveryPlanConfigMissingGatewayIP(name, description, stageUUID string) string {
	return fmt.Sprintf(`
		resource "nutanix_recovery_plan" "test" {
			name        = "%[1]s"
			description = "%[2]s"
			stage_list {
				stage_work {
					recover_entities {
						entity_info_list {
							categories {
								name  = "Environment"
								value = "Dev"
							}
						}
					}
				}
				stage_uuid      = "%[3]s"
				delay_time_secs = 0
			}
			parameters {
				network_mapping_list {
					availability_zone_network_mapping_list {
						availability_zone_url = "c99ab7cd-9191-4fcb-8fc0-232eff76e595"
						recovery_network {
							name = "vlan0"
							subnet_list {
								prefix_length = 24
							}
						}
					}
				}
			}
		}
	`, name, description, stageUUID)
}

// testAccNutanixRecoveryPlanConfigMissingIPAddress omits the required ip_address
// inside an ip_config_list; used by the negative validation test.
func testAccNutanixRecoveryPlanConfigMissingIPAddress(name, description, stageUUID string) string {
	return fmt.Sprintf(`
		resource "nutanix_recovery_plan" "test" {
			name        = "%[1]s"
			description = "%[2]s"
			stage_list {
				stage_work {
					recover_entities {
						entity_info_list {
							categories {
								name  = "Environment"
								value = "Dev"
							}
						}
					}
				}
				stage_uuid      = "%[3]s"
				delay_time_secs = 0
			}
			parameters {
				network_mapping_list {
					availability_zone_network_mapping_list {
						availability_zone_url = "c99ab7cd-9191-4fcb-8fc0-232eff76e595"
						recovery_ip_assignment_list {
							vm_reference {
								kind = "vm"
								uuid = "00000000-0000-0000-0000-000000000000"
							}
							ip_config_list {
							}
						}
					}
				}
			}
		}
	`, name, description, stageUUID)
}
