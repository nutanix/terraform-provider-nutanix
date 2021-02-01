package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixUserDataSource_basic(t *testing.T) {
	principalName := "dou-user-3@ntnxlab.local"
	expectedDisplayName := "dou-user-3"
	directoryServiceUUID := "542d7921-1385-4b6e-ab10-09f2ca4f054d"

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
	user_id = nutanix_user.user.id
}
`, pn, dsuuid)
}

func TestAccNutanixUserDataSource_byName(t *testing.T) {
	principalName := "dou-user@ntnxlab.local"
	expectedDisplayName := "dou-user"
	directoryServiceUUID := "542d7921-1385-4b6e-ab10-09f2ca4f054d"

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
	user_name = nutanix_user.user.name
}
`, pn, dsuuid)
}
