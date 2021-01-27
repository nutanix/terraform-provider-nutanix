package nutanix

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const resourceRole = "nutanix_role.test"

func TestAccNutanixRole_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("accest-role")
	description := "Description of my role"
	nameUpdated := acctest.RandomWithPrefix("accest-role")
	descriptionUpdated := "Description of my role updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleConfigWithCategory(name, description, "Production"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRoleExists(),
					testAccCheckNutanixCategories(resourceRole),
					resource.TestCheckResourceAttr(resourceRole, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceRole, "categories.2228745532.name"),
					resource.TestCheckResourceAttrSet(resourceRole, "categories.2228745532.value"),
					resource.TestCheckResourceAttr(resourceRole, "categories.2228745532.name", "Environment"),
					resource.TestCheckResourceAttr(resourceRole, "categories.2228745532.value", "Production"),
				),
			},
			{
				Config: testAccNutanixRoleConfigWithCategory(nameUpdated, descriptionUpdated, "Staging"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixRoleExists(),
					resource.TestCheckResourceAttr(resourceRole, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceRole, "categories.2940305446.name"),
					resource.TestCheckResourceAttrSet(resourceRole, "categories.2940305446.value"),
					resource.TestCheckResourceAttr(resourceRole, "categories.2940305446.name", "Environment"),
					resource.TestCheckResourceAttr(resourceRole, "categories.2940305446.value", "Staging"),
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
	conn := testAccProvider.Meta().(*Client)

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

func testAccNutanixRoleConfig(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	permission_reference_list {
		kind = "permission"
		uuid = "d08ea95c-8221-4590-a77a-52d69639959a"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "1a8a65c0-4333-42c6-9039-fd2585ceead7"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "bea75573-e8fe-42a3-817a-bd1bd98ab110"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "93e1cc93-d799-4f44-84ad-534814f6db0d"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "62f53a1a-324c-4da6-bcb8-2cecc07b2cb7"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "62f53a1a-324c-4da6-bcb8-2cecc07b2cb7"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "ef38a553-a20f-4a2b-b12d-bb9cca03cbdd"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "6e768a07-21ef-4615-84d0-7ec442ec942f"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "91b77724-b163-473f-94a8-d016e75c18bd"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "491ae1d0-5a8f-4bcc-9cee-068cd01c9274"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "1dbfb7b4-9896-4c2a-b6fe-fbf113bae306"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "740d29f7-18ae-4d07-aeef-3fc901c1887a"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "b28f35be-6561-4a4a-9d90-a298d2de33d7"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "428fad6c-8735-4a7d-bad3-8497bef051c8"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "ea445ec5-f9bb-4af6-92e8-0d72d11ada85"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "85a24ad8-67b6-4b63-b30f-96da1baca161"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "d370823b-82d8-4518-a486-b75ba8e130d6"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "2e9988df-47ae-44ae-9114-ada346657b90"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "4e8e9007-8fbe-4709-a069-278259238e55"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "2e9988df-47ae-44ae-9114-ada346657b90"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "2e9988df-47ae-44ae-9114-ada346657b90"
	}
}
`, name, description)
}

func testAccNutanixRoleConfigWithCategory(name, description, categoryValue string) string {
	return fmt.Sprintf(`
resource "nutanix_role" "test" {
	name        = "%[1]s"
	description = "%[2]s"
	permission_reference_list {
		kind = "permission"
		uuid = "2e9988df-47ae-44ae-9114-ada346657b90"
	}
	categories {
		name = "Environment"
		value = "%[3]s"
	}
}
`, name, description, categoryValue)
}
