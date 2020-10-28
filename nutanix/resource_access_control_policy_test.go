package nutanix

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const resourceAccessPolicy = "nutanix_access_control_policy.accest-access-policy"

func TestAccNutanixAccessControlPolicy_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	description := "Description of my access control policy"
	uuidRole := ""

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfig(uuidRole, name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(resourceAccessPolicy),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", "Description of my access control policy"),
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

func TestAccNutanixAccessControlPolicy_Update(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-access-policy")
	description := "Description of my access control policy"
	uuidRole := ""
	nameUpdated := acctest.RandomWithPrefix("accest-access-policy updated")
	descriptionUpdated := "Description of my access control policy updated"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfig(uuidRole, name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(resourceAccessPolicy),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", "Description of my access control policy"),
				),
			},
			{
				Config: testAccNutanixAccessControlPolicyConfig(uuidRole, nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(resourceAccessPolicy),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "name", fmt.Sprintf("%s updated", nameUpdated)),
					resource.TestCheckResourceAttr(resourceAccessPolicy, "description", fmt.Sprintf("%s updated", descriptionUpdated)),
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
	resourceName := "nutanix_access_control_policy.accest-access-policy-categories"

	name := acctest.RandomWithPrefix("accest-access-policy")
	description := "Description of my access control policy"
	uuidRole := ""
	nameUpdated := acctest.RandomWithPrefix("accest-access-policy updated")
	descriptionUpdated := "Description of my access control policy updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixAccessControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAccessControlPolicyConfigWithCategory(uuidRole, name, description, "Production"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(resourceName),
					testAccCheckNutanixCategories(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "categories.2228745532.name"),
					resource.TestCheckResourceAttrSet(resourceName, "categories.2228745532.value"),
					resource.TestCheckResourceAttr(resourceName, "categories.2228745532.name", "Environment"),
					resource.TestCheckResourceAttr(resourceName, "categories.2228745532.value", "Production"),
				),
			},
			{
				Config: testAccNutanixAccessControlPolicyConfigWithCategory(uuidRole, nameUpdated, descriptionUpdated, "Staging"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixAccessControlPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "categories.2940305446.name"),
					resource.TestCheckResourceAttrSet(resourceName, "categories.2940305446.value"),
					resource.TestCheckResourceAttr(resourceName, "categories.2940305446.name", "Environment"),
					resource.TestCheckResourceAttr(resourceName, "categories.2940305446.value", "Staging"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
		},
	})
}

func testAccCheckNutanixAccessControlPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		pretty, _ := json.MarshalIndent(rs, "", "  ")
		fmt.Print("\n\n[DEBUG] State of AccessControlPolicy", string(pretty))

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixAccessControlPolicyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

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

func testAccNutanixAccessControlPolicyConfig(uuidRole, name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_access_control_policy" "accest-access-policy" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference{
		kind = "role"
		uuid = "%[3]s"
	}
}
`, name, description, uuidRole)
}

func testAccNutanixAccessControlPolicyConfigWithCategory(uuidRole, name, description, categoryValue string) string {
	return fmt.Sprintf(`
resource "nutanix_access_control_policy" "accest-access-policy-categories" {
	name        = "%[1]s"
	description = "%[2]s"
	role_reference{
		kind = "role"
		uuid = "%[3]s"
	}
	categories {
		name = "Environment"
		value = "%[4]s"
	}
}
`, name, description, uuidRole, categoryValue)
}
