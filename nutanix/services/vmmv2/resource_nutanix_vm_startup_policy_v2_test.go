// Acceptance tests for nutanix_vm_startup_policy_v2 resource.
//
// TestAccV2NutanixVmStartupPolicyResource_Basic:
//   - Step 1: Create a VM startup policy with 2 groups, and 1 start condition (power_on).
//   - Step 2: Update name, description, and delay duration.
//   - Step 3: Validate datasource reads (single policy and list) return correct data.
//
// TestAccV2NutanixVmStartupPolicyResource_GuestBootup:
//   - Step 1: Create a VM startup policy with guest_bootup start condition (timeout_duration_secs).
//   - Step 2: Switch power_state_criteria from guest_bootup → power_on.
//   - Step 3: Switch power_state_criteria from power_on → guest_bootup (timeout_duration_secs=300).
//
// TestAccV2NutanixVmStartupPolicyResource_Validation:
//   - Step 1: groups < MinItems (1 group, need 2) → rejected.
//   - Step 2: groups > MaxItems (7 groups, max 6) → rejected.
//   - Step 3: categories > MaxItems (6 categories in one group, max 5) → rejected.
//   - Step 4: start_conditions > MaxItems (6 start conditions, max 5) → rejected.
//   - Step 5: delay_duration_secs > 600 → rejected.
//   - Step 6: delay_duration_secs < 0 → rejected.
//   - Step 7: power_state_criteria > MaxItems 1 (2 blocks) → rejected.
//   - Step 8: guest_bootup > MaxItems 1 (2 blocks) → rejected.
//
// TestAccV2NutanixVmStartupPolicyResource_GroupCategoryUpdates:
//   - Step 1: Create policy with 2 groups, 1 category each — [cat1], [cat2].
//   - Step 2: Replace a category in an existing group — [cat3], [cat2].
//   - Step 3: Add a new category to an existing group — [cat3, cat4], [cat2].
//   - Step 4: Add a new group (start_conditions scale to N-1) — [cat3, cat4], [cat2], [cat5].
//   - Step 5: Remove an existing group — [cat3, cat4], [cat5].
package vmmv2_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vmstartuppolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameVmStartupPolicy = "nutanix_vm_startup_policy_v2.test"
const datasourceNameVmStartupPolicy = "data.nutanix_vm_startup_policy_v2.test_ds"
const datasourceNameVmStartupPolicies = "data.nutanix_vm_startup_policies_v2.test_list"

func TestAccV2NutanixVmStartupPolicyResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-startup-policy-%d", r)
	updatedName := fmt.Sprintf("test-vm-startup-policy-%d-updated", r)
	desc := "test vm startup policy description"
	updatedDesc := "test vm startup policy description updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmStartupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVmStartupPolicyV2Config(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "create_time"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.0.categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.1.categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.delay_duration_secs", "30"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.guest_bootup.#", "0"),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "num_compliant_vms"),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "num_non_compliant_vms"),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "num_pending_vms"),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "num_dependency_conflicts"),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "num_start_condition_conflicts"),
				),
			},
			{
				Config: testVmStartupPolicyV2ConfigUpdated(updatedName, updatedDesc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.delay_duration_secs", "60"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.guest_bootup.#", "0"),
				),
			},
			{
				Config: testVmStartupPolicyV2ConfigWithDatasources(updatedName, updatedDesc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVmStartupPolicy, "name", updatedName),
					resource.TestCheckResourceAttr(datasourceNameVmStartupPolicy, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(datasourceNameVmStartupPolicy, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVmStartupPolicy, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVmStartupPolicies, "policies.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmStartupPolicyResource_GuestBootup(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-startup-policy-gb-%d", r)
	desc := "test vm startup policy with guest bootup"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmStartupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVmStartupPolicyV2ConfigGuestBootup(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.guest_bootup.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.guest_bootup.0.timeout_duration_secs", "120"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "0"),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "ext_id"),
				),
			},
			{
				Config: testVmStartupPolicyV2ConfigSwitchToPowerOn(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.delay_duration_secs", "30"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.guest_bootup.#", "0"),
				),
			},
			{
				Config: testVmStartupPolicyV2ConfigSwitchToGuestBootup(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.delay_duration_secs", "30"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.guest_bootup.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.guest_bootup.0.timeout_duration_secs", "300"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "0"),
				),
			},
		},
	})
}

func testAccCheckNutanixVmStartupPolicyDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_vm_startup_policy_v2" {
			continue
		}
		getRequest := import2.GetVmStartupPolicyByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := conn.VmmAPI.VmStartupPoliciesAPIInstance.GetVmStartupPolicyById(ctx, &getRequest)
		if err == nil {
			return fmt.Errorf("VM Startup Policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testVspBaseConfig(r int) string {
	return fmt.Sprintf(`

		resource "nutanix_category_v2" "cat1" {
			key            = "vm-startup-cat1-%[1]d"
			value          = "vm-startup-cat1-value-%[1]d"
		}
		resource "nutanix_category_v2" "cat2" {
			key            = "vm-startup-cat2-%[1]d"
			value          = "vm-startup-cat2-value-%[1]d"
		}
	`, r)
}

func testVmStartupPolicyV2Config(name, desc string, r int) string {
	return testVspBaseConfig(r) + fmt.Sprintf(`
		resource "nutanix_vm_startup_policy_v2" "test" {
			name           = "%[1]s"
			description    = "%[2]s"
			groups {
				categories {
					ext_id = nutanix_category_v2.cat1.id
				}
			}
			groups {
				categories {
					ext_id = nutanix_category_v2.cat2.id
				}
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					power_on {}
				}
			}
		}
`, name, desc)
}

func testVmStartupPolicyV2ConfigUpdated(name, desc string, r int) string {
	return testVspBaseConfig(r) + fmt.Sprintf(`
		resource "nutanix_vm_startup_policy_v2" "test" {
			name           = "%[1]s"
			description    = "%[2]s"
			groups {
				categories {
					ext_id = nutanix_category_v2.cat1.id
				}
			}
			groups {
				categories {
					ext_id = nutanix_category_v2.cat2.id
				}
			}
			start_conditions {
				delay_duration_secs = 60
				power_state_criteria {
					power_on {}
				}
			}
		}
`, name, desc)
}

func testVmStartupPolicyV2ConfigWithDatasources(name, desc string, r int) string {
	return testVmStartupPolicyV2ConfigUpdated(name, desc, r) + `
		data "nutanix_vm_startup_policy_v2" "test_ds" {
			ext_id = nutanix_vm_startup_policy_v2.test.id
		}
		data "nutanix_vm_startup_policies_v2" "test_list" {
			depends_on = [nutanix_vm_startup_policy_v2.test]
		}
	`
}

// TestAccV2NutanixVmStartupPolicyResource_Validation tests schema-level validation:
//   - Step 1: groups < MinItems (1 group, need 2) → rejected.
//   - Step 2: groups > MaxItems (7 groups, max 6) → rejected.
//   - Step 3: categories > MaxItems (6 categories in one group, max 5) → rejected.
//   - Step 4: start_conditions > MaxItems (6 start conditions, max 5) → rejected.
//   - Step 5: delay_duration_secs > 600 → rejected.
//   - Step 6: delay_duration_secs < 0 → rejected.
//   - Step 7: power_state_criteria > MaxItems 1 (2 blocks) → rejected.
//   - Step 8: guest_bootup > MaxItems 1 (2 blocks) → rejected.
func TestAccV2NutanixVmStartupPolicyResource_Validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testVspValidationConfig_TooFewGroups(),
				ExpectError: regexp.MustCompile(`At least 2 "groups" blocks are required`),
			},
			{
				Config:      testVspValidationConfig_TooManyGroups(),
				ExpectError: regexp.MustCompile(`No more than 6 "groups" blocks are allowed`),
			},
			{
				Config:      testVspValidationConfig_TooManyCategories(),
				ExpectError: regexp.MustCompile(`No more than 5 "categories" blocks are allowed`),
			},
			{
				Config:      testVspValidationConfig_TooManyStartConditions(),
				ExpectError: regexp.MustCompile(`No more than 5 "start_conditions" blocks are allowed`),
			},
			{
				Config:      testVspValidationConfig_DelayTooHigh(),
				ExpectError: regexp.MustCompile(`expected .* to be in the range \(0 - 600\)`),
			},
			{
				Config:      testVspValidationConfig_DelayNegative(),
				ExpectError: regexp.MustCompile(`expected .* to be in the range \(0 - 600\)`),
			},
			{
				Config:      testVspValidationConfig_DuplicatePowerStateCriteria(),
				ExpectError: regexp.MustCompile(`No more than 1 "power_state_criteria" blocks are allowed`),
			},
			{
				Config:      testVspValidationConfig_DuplicateGuestBootup(),
				ExpectError: regexp.MustCompile(`No more than 1 "guest_bootup" blocks are allowed`),
			},
		},
	})
}

func testVspValidationConfig_TooFewGroups() string {
	return `
		resource "nutanix_vm_startup_policy_v2" "test" {
			name = "validation-test"
			groups {
				categories {
					ext_id = "00000000-0000-0000-0000-000000000001"
				}
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					power_on {}
				}
			}
		}
	`
}

func testVspValidationConfig_TooManyGroups() string {
	return `
		resource "nutanix_vm_startup_policy_v2" "test" {
			name = "validation-test"
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000001" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000002" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000003" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000004" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000005" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000006" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000007" }
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					power_on {}
				}
			}
		}
	`
}

func testVspValidationConfig_TooManyCategories() string {
	return `
		resource "nutanix_vm_startup_policy_v2" "test" {
			name = "validation-test"
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000001" }
				categories { ext_id = "00000000-0000-0000-0000-000000000002" }
				categories { ext_id = "00000000-0000-0000-0000-000000000003" }
				categories { ext_id = "00000000-0000-0000-0000-000000000004" }
				categories { ext_id = "00000000-0000-0000-0000-000000000005" }
				categories { ext_id = "00000000-0000-0000-0000-000000000006" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000007" }
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					power_on {}
				}
			}
		}
	`
}

func testVspValidationConfig_TooManyStartConditions() string {
	return `
		resource "nutanix_vm_startup_policy_v2" "test" {
			name = "validation-test"
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000001" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000002" }
			}
			start_conditions {
				delay_duration_secs = 10
				power_state_criteria {
					power_on {}
				}
			}
			start_conditions {
				delay_duration_secs = 20
				power_state_criteria {
					power_on {}
				}
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					power_on {}
				}
			}
			start_conditions {
				delay_duration_secs = 40
				power_state_criteria {
					power_on {}
				}
			}
			start_conditions {
				delay_duration_secs = 50
				power_state_criteria {
					power_on {}
				}
			}
			start_conditions {
				delay_duration_secs = 60
				power_state_criteria {
					power_on {}
				}
			}
		}
	`
}

func testVspValidationConfig_DelayTooHigh() string {
	return `
		resource "nutanix_vm_startup_policy_v2" "test" {
			name = "validation-test"
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000001" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000002" }
			}
			start_conditions {
				delay_duration_secs = 601
				power_state_criteria {
					power_on {}
				}
			}
		}
	`
}

func testVspValidationConfig_DelayNegative() string {
	return `
		resource "nutanix_vm_startup_policy_v2" "test" {
			name = "validation-test"
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000001" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000002" }
			}
			start_conditions {
				delay_duration_secs = -1
				power_state_criteria {
					power_on {}
				}
			}
		}
	`
}

func testVspValidationConfig_DuplicatePowerStateCriteria() string {
	return `
		resource "nutanix_vm_startup_policy_v2" "test" {
			name = "validation-test"
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000001" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000002" }
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					power_on {}
				}
				power_state_criteria {
					power_on {}
				}
			}
		}
	`
}

func testVspValidationConfig_DuplicateGuestBootup() string {
	return `
		resource "nutanix_vm_startup_policy_v2" "test" {
			name = "validation-test"
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000001" }
			}
			groups {
				categories { ext_id = "00000000-0000-0000-0000-000000000002" }
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					guest_bootup {
						timeout_duration_secs = 60
					}
					guest_bootup {
						timeout_duration_secs = 120
					}
				}
			}
		}
	`
}

// TestAccV2NutanixVmStartupPolicyResource_GroupCategoryUpdates tests all group/category mutations:
//
// Step 1 (Create):  group0=[cat1], group1=[cat2]
// Step 2 (Update category in existing group): group0=[cat3], group1=[cat2]      — cat1 replaced by cat3 in group0
// Step 3 (Add new category to existing group): group0=[cat3,cat4], group1=[cat2] — cat4 added to group0
// Step 4 (Add new group): group0=[cat3,cat4], group1=[cat2], group2=[cat5]      — group2 added
// Step 5 (Remove existing group): group0=[cat3,cat4], group2=[cat5]             — group1 removed
func TestAccV2NutanixVmStartupPolicyResource_GroupCategoryUpdates(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vsp-grp-upd-%d", r)
	desc := "test group category updates"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmStartupPolicyDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with group0=[cat1], group1=[cat2]
			{
				Config: testVspGroupUpdateConfig(name, desc, r, [][]string{{"cat_gu1"}, {"cat_gu2"}}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVmStartupPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.0.categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.1.categories.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.0.categories.0.ext_id", "nutanix_category_v2.cat_gu1", "id"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.1.categories.0.ext_id", "nutanix_category_v2.cat_gu2", "id"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "1"),
				),
			},
			// Step 2: Update category in existing group — group0=[cat3], group1=[cat2]
			{
				Config: testVspGroupUpdateConfig(name, desc, r, [][]string{{"cat_gu3"}, {"cat_gu2"}}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.0.categories.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.0.categories.0.ext_id", "nutanix_category_v2.cat_gu3", "id"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.1.categories.0.ext_id", "nutanix_category_v2.cat_gu2", "id"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "1"),
				),
			},
			// Step 3: Add new category to existing group — group0=[cat3,cat4], group1=[cat2]
			{
				Config: testVspGroupUpdateConfig(name, desc, r, [][]string{{"cat_gu3", "cat_gu4"}, {"cat_gu2"}}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.0.categories.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.0.categories.0.ext_id", "nutanix_category_v2.cat_gu3", "id"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.0.categories.1.ext_id", "nutanix_category_v2.cat_gu4", "id"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.1.categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "1"),
				),
			},
			// Step 4: Add new group — group0=[cat3,cat4], group1=[cat2], group2=[cat5]
			{
				Config: testVspGroupUpdateConfig(name, desc, r, [][]string{{"cat_gu3", "cat_gu4"}, {"cat_gu2"}, {"cat_gu5"}}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "3"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.1.power_state_criteria.0.power_on.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.0.categories.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.1.categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.2.categories.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.2.categories.0.ext_id", "nutanix_category_v2.cat_gu5", "id"),
				),
			},
			// Step 5: Remove group1 — group0=[cat3,cat4], group2=[cat5]
			{
				Config: testVspGroupUpdateConfig(name, desc, r, [][]string{{"cat_gu3", "cat_gu4"}, {"cat_gu5"}}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "start_conditions.0.power_state_criteria.0.power_on.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.0.categories.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.0.categories.0.ext_id", "nutanix_category_v2.cat_gu3", "id"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.0.categories.1.ext_id", "nutanix_category_v2.cat_gu4", "id"),
					resource.TestCheckResourceAttr(resourceNameVmStartupPolicy, "groups.1.categories.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameVmStartupPolicy, "groups.1.categories.0.ext_id", "nutanix_category_v2.cat_gu5", "id"),
				),
			},
		},
	})
}

// testVspGroupUpdateConfig generates a config that always declares 5 categories
// and builds the groups block dynamically from the given layout.
// Each element of `groups` is a slice of category resource names (e.g. ["cat_gu1", "cat_gu3"]).
func testVspGroupUpdateConfig(name, desc string, r int, groups [][]string) string {
	catBlock := fmt.Sprintf(`
		resource "nutanix_category_v2" "cat_gu1" {
			key            = "vsp-gu-cat1-%[1]d"
			value          = "vsp-gu-cat1-val-%[1]d"
		}
		resource "nutanix_category_v2" "cat_gu2" {
			key            = "vsp-gu-cat2-%[1]d"
			value          = "vsp-gu-cat2-val-%[1]d"
		}
		resource "nutanix_category_v2" "cat_gu3" {
			key            = "vsp-gu-cat3-%[1]d"
			value          = "vsp-gu-cat3-val-%[1]d"
		}
		resource "nutanix_category_v2" "cat_gu4" {
			key            = "vsp-gu-cat4-%[1]d"
			value          = "vsp-gu-cat4-val-%[1]d"
		}
		resource "nutanix_category_v2" "cat_gu5" {
			key            = "vsp-gu-cat5-%[1]d"
			value          = "vsp-gu-cat5-val-%[1]d"
		}
	`, r)

	groupsBlock := ""
	for _, g := range groups {
		catsBlock := ""
		for _, catName := range g {
			catsBlock += fmt.Sprintf(`
				categories {
					ext_id = nutanix_category_v2.%s.id
				}`, catName)
		}
		groupsBlock += fmt.Sprintf(`
			groups {%s
			}`, catsBlock)
	}

	startConditionsBlock := ""
	for i := 0; i < len(groups)-1; i++ {
		startConditionsBlock += fmt.Sprintf(`
			start_conditions {
				delay_duration_secs = %d
				power_state_criteria {
					power_on {}
				}
			}`, 30+i*10)
	}

	return fmt.Sprintf(`
		%s
		resource "nutanix_vm_startup_policy_v2" "test" {
			name           = "%s"
			description    = "%s"
			%s
			%s
		}
	`, catBlock, name, desc, groupsBlock, startConditionsBlock)
}

func testVspGuestBootupBaseConfig(r int) string {
	return fmt.Sprintf(`
		resource "nutanix_category_v2" "cat_gb1" {
			key            = "vm-startup-gb-cat1-%[1]d"
			value          = "vm-startup-gb-cat1-value-%[1]d"
		}
		resource "nutanix_category_v2" "cat_gb2" {
			key            = "vm-startup-gb-cat2-%[1]d"
			value          = "vm-startup-gb-cat2-value-%[1]d"
		}
	`, r)
}

func testVmStartupPolicyV2ConfigGuestBootup(name, desc string, r int) string {
	return testVspGuestBootupBaseConfig(r) + fmt.Sprintf(`
		resource "nutanix_vm_startup_policy_v2" "test" {
			name           = "%[1]s"
			description    = "%[2]s"
			groups {
				categories {
					ext_id = nutanix_category_v2.cat_gb1.id
				}
			}
			groups {
				categories {
					ext_id = nutanix_category_v2.cat_gb2.id
				}
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					guest_bootup {
						timeout_duration_secs = 120
					}
				}
			}
		}
`, name, desc)
}

func testVmStartupPolicyV2ConfigSwitchToPowerOn(name, desc string, r int) string {
	return testVspGuestBootupBaseConfig(r) + fmt.Sprintf(`
		resource "nutanix_vm_startup_policy_v2" "test" {
			name           = "%[1]s"
			description    = "%[2]s"
			groups {
				categories {
					ext_id = nutanix_category_v2.cat_gb1.id
				}
			}
			groups {
				categories {
					ext_id = nutanix_category_v2.cat_gb2.id
				}
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					power_on {}
				}
			}
		}
`, name, desc)
}

func testVmStartupPolicyV2ConfigSwitchToGuestBootup(name, desc string, r int) string {
	return testVspGuestBootupBaseConfig(r) + fmt.Sprintf(`
		resource "nutanix_vm_startup_policy_v2" "test" {
			name           = "%[1]s"
			description    = "%[2]s"
			groups {
				categories {
					ext_id = nutanix_category_v2.cat_gb1.id
				}
			}
			groups {
				categories {
					ext_id = nutanix_category_v2.cat_gb2.id
				}
			}
			start_conditions {
				delay_duration_secs = 30
				power_state_criteria {
					guest_bootup {
						timeout_duration_secs = 300
					}
				}
			}
		}
`, name, desc)
}
