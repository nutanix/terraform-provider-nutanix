package nutanix

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const resourceNameUser = "nutanix_user.user"

func TestAccNutanixUser_basic(t *testing.T) {
	principalName := "dou-user@ntnxlab.local"
	directoryServiceUUID := "dd19a896-8e72-4158-b716-98455ceda220"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixUserConfig(principalName, directoryServiceUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixUserExists(resourceNameUser),
					resource.TestCheckResourceAttr(resourceNameUser, "name", principalName),
					resource.TestCheckResourceAttr(resourceNameUser, "directory_service_user.#", "1"),
				),
			},
			{
				Config: testAccNutanixUserConfig(principalName, directoryServiceUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixUserExists(resourceNameUser),
					resource.TestCheckResourceAttr(resourceNameUser, "name", principalName),
					resource.TestCheckResourceAttr(resourceNameUser, "directory_service_user.#", "1"),
				),
			},
			{
				ResourceName:      resourceNameUser,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixUser_IdentityProvider(t *testing.T) {
	username := "dou-user-2@ntnxlab.local"
	identityProviderUUID := "02316a2c-cc8c-41de-9abb-f07c4da58fda"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixUserConfig_IdentityProvider(username, identityProviderUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixUserExists(resourceNameUser),
					resource.TestCheckResourceAttr(resourceNameUser, "name", username),
					resource.TestCheckResourceAttr(resourceNameUser, "identity_provider_user.#", "1"),
				),
			},
			{
				ResourceName:      resourceNameUser,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixUserDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_user" {
			continue
		}
		if _, err := conn.API.V3.GetUser(rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccCheckNutanixUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccNutanixUserConfig(pn, dsuuid string) string {
	return fmt.Sprintf(`
resource "nutanix_user" "user" {
	directory_service_user {
		user_principal_name = "%s"
		directory_service_reference {
		  uuid = "%s"
		}
	}
}
`, pn, dsuuid)
}

func testAccNutanixUserConfig_IdentityProvider(username, ipuuid string) string {
	return fmt.Sprintf(`
resource "nutanix_user" "user" {
	identity_provider_user {
		username = "%s"
		identity_provider_reference {
		  uuid = "%s"
		}
	}
}
`, username, ipuuid)
}
