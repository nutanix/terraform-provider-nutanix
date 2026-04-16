package iam_test

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

const resourceRole = "nutanix_role.test"

func TestAccNutanixRole_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-role")
	description := "Description of my role"
	nameUpdated := acctest.RandomWithPrefix("accest-role")
	descriptionUpdated := "Description of my role updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRoleExists(),
					resource.TestCheckResourceAttr(resourceRole, "name", name),
					resource.TestCheckResourceAttr(resourceRole, "description", description),
				),
			},
			{
				Config: testAccNutanixRoleConfig(nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRoleExists(),
					resource.TestCheckResourceAttr(resourceRole, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceRole, "description", descriptionUpdated),
				),
			},
			{
				ResourceName:      resourceRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixRole_WithCategory(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-role")
	description := "Description of my role"
	nameUpdated := acctest.RandomWithPrefix("accest-role")
	descriptionUpdated := "Description of my role updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleConfigWithCategory(name, description, "Production"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRoleExists(),
					testAccCheckNutanixCategories(resourceRole),
					resource.TestCheckResourceAttr(resourceRole, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceRole, "categories.0.name"),
					resource.TestCheckResourceAttrSet(resourceRole, "categories.0.value"),
					resource.TestCheckResourceAttr(resourceRole, "categories.0.name", "Environment"),
					resource.TestCheckResourceAttr(resourceRole, "categories.0.value", "Production"),
				),
			},
			{
				Config: testAccNutanixRoleConfigWithCategory(nameUpdated, descriptionUpdated, "Staging"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRoleExists(),
					resource.TestCheckResourceAttr(resourceRole, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceRole, "categories.0.name"),
					resource.TestCheckResourceAttrSet(resourceRole, "categories.0.value"),
					resource.TestCheckResourceAttr(resourceRole, "categories.0.name", "Environment"),
					resource.TestCheckResourceAttr(resourceRole, "categories.0.value", "Staging"),
				),
			},
			{
				ResourceName:      resourceRole,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixRoleExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceRole]
		if !ok {
			return fmt.Errorf("not found: %s", resourceRole)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixRoleDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_role" {
			continue
		}
		if _, err := resourceNutanixRoleExists(conn.API, rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
	}

	return nil
}

func resourceNutanixRoleExists(conn *v3.Client, name string) (*string, error) {
	var accessUUID *string

	filter := fmt.Sprintf("name==%s", name)
	accessList, err := conn.V3.ListAllRole(filter)
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

func testAccNutanixRoleConfig(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	permission_reference_list {
		kind = "permission"
		uuid = "%[3]s"
	}
}
`, name, description, testVars.Permissions[0].UUID)
}

func testAccNutanixRoleConfigWithCategory(name, description, categoryValue string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	permission_reference_list {
		kind = "permission"
		uuid = "%[4]s"
	}
	categories {
		name = "Environment"
		value = "%[3]s"
	}
}
`, name, description, categoryValue, testVars.Permissions[0].UUID)
}
