package prism_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceAccessPolicy = "nutanix_access_control_policy.test"

func TestAccNutanixAccessControlPolicy_basic(t *testing.T) {
	t.Skip() // https://jira.nutanix.com/browse/ENG-483192
	name := acctest.RandomWithPrefix("accest-access-policy")
	roleName := acctest.RandomWithPrefix("test-acc-role")
	description := "Description of my access control policy"
	nameUpdated := acctest.RandomWithPrefix("accest-access-policy")
	descriptionUpdated := "Description of my access control policy updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfig(name, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", description),
				),
			},
			{
				Config: testAccNutanixAccessControlPolicyConfig(nameUpdated, descriptionUpdated, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", descriptionUpdated),
				),
			},
			{
				ResourceName:      resourceAccessPolicy,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixAccessControlPolicy_WithUser(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	roleName := acctest.RandomWithPrefix("test-acc-role-")
	description := "Description of my access control policy"
	nameUpdated := acctest.RandomWithPrefix("accest-access-policy")
	descriptionUpdated := "Description of my access control policy updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfigWithUser(name, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", description),
				),
			},
			{
				Config: testAccNutanixAccessControlPolicyConfigWithUser(nameUpdated, descriptionUpdated, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", descriptionUpdated),
				),
			},
			{
				ResourceName:      resourceAccessPolicy,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixAccessControlPolicy_WithCategory(t *testing.T) {
	t.Skip() // https://jira.nutanix.com/browse/ENG-483192
	name := acctest.RandomWithPrefix("accest-access-policy")
	roleName := acctest.RandomWithPrefix("test-acc-role-")
	description := "Description of my access control policy"
	nameUpdated := acctest.RandomWithPrefix("accest-access-policy")
	descriptionUpdated := "Description of my access control policy updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfigWithCategory(name, description, "Production", roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					testAccCheckNutanixCategories(resourceAccessPolicy),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceAccessPolicy, "categories.0.name"),
					resource.TestCheckResourceAttrSet(resourceAccessPolicy, "categories.0.value"),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "categories.0.name", "Environment"),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "categories.0.value", "Production"),
				),
			},
			{
				Config: testAccNutanixAccessControlPolicyConfigWithCategory(nameUpdated, descriptionUpdated, "Staging", roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceAccessPolicy, "categories.0.name"),
					resource.TestCheckResourceAttrSet(resourceAccessPolicy, "categories.0.value"),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "categories.0.name", "Environment"),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "categories.0.value", "Staging"),
				),
			},
			{
				ResourceName:      resourceAccessPolicy,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixAccessControlPolicy_WithCluster(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	roleName := acctest.RandomWithPrefix("test-acc-role-")
	description := "Description of my access control policy"
	nameUpdated := acctest.RandomWithPrefix("accest-access-policy")
	descriptionUpdated := "Description of my access control policy updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfigWithCluster(name, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", description),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "context_filter_list.0.scope_filter_expression_list.0.left_hand_side", "CLUSTER"),
				),
			},
			{
				Config: testAccNutanixAccessControlPolicyConfigWithCluster(nameUpdated, descriptionUpdated, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "context_filter_list.0.scope_filter_expression_list.0.left_hand_side", "CLUSTER"),
				),
			},
		},
	})
}

func TestAccNutanixAccessControlPolicy_WithVPC(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	roleName := acctest.RandomWithPrefix("test-acc-role-")
	description := "Description of my access control policy"
	nameUpdated := acctest.RandomWithPrefix("accest-access-policy")
	descriptionUpdated := "Description of my access control policy updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfigWithVPC(name, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", description),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "context_filter_list.0.scope_filter_expression_list.0.left_hand_side", "VPC"),
				),
			},
			{
				Config: testAccNutanixAccessControlPolicyConfigWithVPC(nameUpdated, descriptionUpdated, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "context_filter_list.0.scope_filter_expression_list.0.left_hand_side", "VPC"),
				),
			},
		},
	})
}

func TestAccNutanixAccessControlPolicy_WithConnectivityType(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	roleName := acctest.RandomWithPrefix("test-acc-role-")
	description := "Description of my access control policy"
	nameUpdated := acctest.RandomWithPrefix("accest-access-policy")
	descriptionUpdated := "Description of my access control policy updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfigWithConnectivityType(name, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", description),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "context_filter_list.0.scope_filter_expression_list.0.left_hand_side", "CONNECTIVITY_TYPE"),
				),
			},
			{
				Config: testAccNutanixAccessControlPolicyConfigWithConnectivityType(nameUpdated, descriptionUpdated, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "context_filter_list.0.scope_filter_expression_list.0.left_hand_side", "CONNECTIVITY_TYPE"),
				),
			},
		},
	})
}

func testAccCheckNutanixAccessControlPolicyExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAccessPolicy]
		if !ok {
			return fmt.Errorf("not found: %s", resourceAccessPolicy)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixAccessControlPolicyDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_access_control_policy" {
			continue
		}
		if _, err := resourceNutanixAccessControlPolicyExists(conn.API, rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
	}

	return nil
}

func resourceNutanixAccessControlPolicyExists(conn *v3.Client, name string) (*string, error) {
	var accessUUID *string

	filter := fmt.Sprintf("name==%s", name)
	accessList, err := conn.V3.ListAllAccessControlPolicy(filter)
	if err != nil {
		return nil, err
	}

	for _, access := range accessList.Entities {
		if utils.StringValue(access.Status.Name) == name {
			accessUUID = access.Metadata.UUID
		}
	}
	return accessUUID, nil
}

func testAccCheckNutanixCategories(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if val, ok := rs.Primary.Attributes["categories.0.name"]; !ok || val == "" {
			return fmt.Errorf("%s: manual Attribute '%s' expected to be set", n, "categories.0.name")
		}

		return nil
	}
}

func testAccNutanixAccessControlPolicyConfig(name, description, roleName string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[3]s"
	description = "description role"
	permission_reference_list {
		kind = "permission"
		uuid = "%[4]s"
	}
}
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference {
		kind = "role"
		uuid = nutanix_role.test.id
	}
}
`, name, description, roleName, testVars.Permissions[0].UUID)
}

func testAccNutanixAccessControlPolicyConfigWithCategory(name, description, categoryValue, roleName string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[4]s"
	description = "description role"
	permission_reference_list {
		kind = "permission"
		uuid = "%[5]s"
	}
}
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference {
		kind = "role"
		uuid = nutanix_role.test.id
	}
	categories {
		name = "Environment"
		value = "%[3]s"
	}
}
`, name, description, categoryValue, roleName, testVars.Permissions[0].UUID)
}

func testAccNutanixAccessControlPolicyConfigWithUser(name, description, roleName string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[3]s"
	description = "description role"
	permission_reference_list {
		kind = "permission"
		uuid = "%[4]s"
	}
}
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference {
		kind = "role"
		uuid = nutanix_role.test.id
	}
	user_reference_list{
		uuid = "00000000-0000-0000-0000-000000000000"
		name = "admin"
	}

	context_filter_list{
		scope_filter_expression_list{
			operator = "IN"
			left_hand_side = "PROJECT"
			right_hand_side {
				categories {
						name = "Environment"
						value = ["Dev", "Dev1"]
					}
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "all"
			right_hand_side{
				collection = "ALL"
			}
		}
	}

	context_filter_list{
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "cluster"
			right_hand_side{
				uuid_list = ["00058ef8-c31c-f0bc-0000-000000007b23"]
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "image"
			right_hand_side{
				collection = "ALL"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "category"
			right_hand_side{
				collection = "ALL"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "marketplace_item"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "app_task"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "app_variable"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
	}
}
`, name, description, roleName, testVars.Permissions[0].UUID)
}

func testAccNutanixAccessControlPolicyConfigWithCluster(name, description, roleName string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[3]s"
	description = "description role"
	permission_reference_list {
		kind = "permission"
		uuid = "%[4]s"
	}
}
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference {
		kind = "role"
		uuid = nutanix_role.test.id
	}
	user_reference_list{
		uuid = "00000000-0000-0000-0000-000000000000"
		name = "admin"
	}

	context_filter_list{
		scope_filter_expression_list{
			operator = "IN"
			left_hand_side = "CLUSTER"
			right_hand_side {
				categories {
						name = "Environment"
						value = ["Dev", "Dev1"]
					}
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "all"
			right_hand_side{
				collection = "ALL"
			}
		}
	}

	context_filter_list{
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "cluster"
			right_hand_side{
				uuid_list = ["00058ef8-c31c-f0bc-0000-000000007b23"]
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "image"
			right_hand_side{
				collection = "ALL"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "category"
			right_hand_side{
				collection = "ALL"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "marketplace_item"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "app_task"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "app_variable"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
	}
}
`, name, description, roleName, testVars.Permissions[0].UUID)
}

func testAccNutanixAccessControlPolicyConfigWithVPC(name, description, roleName string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[3]s"
	description = "description role"
	permission_reference_list {
		kind = "permission"
		uuid = "%[4]s"
	}
}
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference {
		kind = "role"
		uuid = nutanix_role.test.id
	}
	user_reference_list{
		uuid = "00000000-0000-0000-0000-000000000000"
		name = "admin"
	}

	context_filter_list{
		scope_filter_expression_list{
			operator = "IN"
			left_hand_side = "VPC"
			right_hand_side {
				categories {
						name = "Environment"
						value = ["Dev", "Dev1"]
					}
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "all"
			right_hand_side{
				collection = "ALL"
			}
		}
	}

	context_filter_list{
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "cluster"
			right_hand_side{
				uuid_list = ["00058ef8-c31c-f0bc-0000-000000007b23"]
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "image"
			right_hand_side{
				collection = "ALL"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "category"
			right_hand_side{
				collection = "ALL"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "marketplace_item"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "app_task"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "app_variable"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
	}
}
`, name, description, roleName, testVars.Permissions[0].UUID)
}

func testAccNutanixAccessControlPolicyConfigWithConnectivityType(name, description, roleName string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[3]s"
	description = "description role"
	permission_reference_list {
		kind = "permission"
		uuid = "%[4]s"
	}
}
resource "nutanix_access_control_policy" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference {
		kind = "role"
		uuid = nutanix_role.test.id
	}
	user_reference_list{
		uuid = "00000000-0000-0000-0000-000000000000"
		name = "admin"
	}

	context_filter_list{
		scope_filter_expression_list{
			operator = "IN"
			left_hand_side = "CONNECTIVITY_TYPE"
			right_hand_side {
				categories {
						name = "Environment"
						value = ["Dev", "Dev1"]
					}
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "all"
			right_hand_side{
				collection = "ALL"
			}
		}
	}

	context_filter_list{
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "cluster"
			right_hand_side{
				uuid_list = ["00058ef8-c31c-f0bc-0000-000000007b23"]
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "image"
			right_hand_side{
				collection = "ALL"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "category"
			right_hand_side{
				collection = "ALL"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "marketplace_item"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "app_task"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
		entity_filter_expression_list{
			operator = "IN"
			left_hand_side_entity_type = "app_variable"
			right_hand_side{
				collection = "SELF_OWNED"
			}
		}
	}
}
`, name, description, roleName, testVars.Permissions[0].UUID)
}
