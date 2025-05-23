package iam_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameUser = "nutanix_user.user"

func TestAccNutanixUser_basic(t *testing.T) {
	principalName := testVars.Users[2].PrincipalName
	directoryServiceUUID := testVars.Users[2].DirectoryServiceUUID
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
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
	t.Skip()
	username := "dou-user-2@ntnxlab.local"
	identityProviderUUID := "02316a2c-cc8c-41de-9abb-f07c4da58fda"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixUserConfigIdentityProvider(username, identityProviderUUID),
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
	conn := acc.TestAccProvider.Meta().(*conns.Client)

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

func testAccNutanixUserConfigIdentityProvider(username, ipuuid string) string {
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
