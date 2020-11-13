package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixUserDataSource_basic(t *testing.T) {
	principalName := "dou-user@ntnxlab.local"
	expectedDisplayName := "dou-user"
	directoryServiceUUID := "dd19a896-8e72-4158-b716-98455ceda220"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(principalName, directoryServiceUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_user.user", "display_name", expectedDisplayName),
					resource.TestCheckResourceAttrSet("data.nutanix_user.user", "directory_service_user.#"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(pn, dsuuid string) string {
	return fmt.Sprintf(`
resource "nutanix_user" "user" {
	directory_service_user {
		user_principal_name = "%s"
		directory_service_reference {
		uuid = "%s"
		}
	}
}

data "nutanix_user" "user" {
	uuid = nutanix_user.user.id
}
`, pn, dsuuid)
}

func TestAccNutanixUserDataSource_byName(t *testing.T) {
	principalName := "dou-user@ntnxlab.local"
	expectedDisplayName := "dou-user"
	directoryServiceUUID := "dd19a896-8e72-4158-b716-98455ceda220"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfigByName(principalName, directoryServiceUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_user.user", "display_name", expectedDisplayName),
					resource.TestCheckResourceAttrSet("data.nutanix_user.user", "directory_service_user.#"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfigByName(pn, dsuuid string) string {
	return fmt.Sprintf(`
resource "nutanix_user" "user" {
	directory_service_user {
		user_principal_name = "%s"
		directory_service_reference {
		uuid = "%s"
		}
	}
}

data "nutanix_user" "user" {
	name = nutanix_user.user.name
}
`, pn, dsuuid)
}
