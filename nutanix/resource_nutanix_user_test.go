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

// func resourceNutanixUserExists(conn *v3.Client, name string) (*string, error) {
// 	var userUUID *string

// 	filter := fmt.Sprintf("name==%s", name)
// 	userList, err := conn.V3.ListAllUser(filter)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, user := range userList.Entities {
// 		if utils.StringValue(user.Status.Name) == name {
// 			userUUID = user.Metadata.UUID
// 		}
// 	}
// 	return userUUID, nil
// }

func testAccCheckNutanixUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		//pretty, _ := json.MarshalIndent(rs, "", "  ")
		//fmt.Print("\n\n[DEBUG] State of User", string(pretty))

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
