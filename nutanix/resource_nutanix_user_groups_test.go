package nutanix

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceNameUserGroups = "nutanix_user_groups.acctest-managed"

func TestAccNutanixUserGroups_basic(t *testing.T) {
	directoryServiceDistName := testVars.UserGroupWithDistinguishedName[1].DistinguishedName
	updatedDirectoryServiceDistName := testVars.UserGroupWithDistinguishedName[1].DistinguishedName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixUserGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixUserGroupsConfig(directoryServiceDistName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_user_group.#", "1"),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_user_group.0.distinguished_name", directoryServiceDistName),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_ou.#", "0"),
				),
			},
			{
				Config: testAccNutanixUserGroupsConfigUpdated(updatedDirectoryServiceDistName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_user_group.#", "1"),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_user_group.0.distinguished_name", updatedDirectoryServiceDistName),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_ou.#", "0"),
				),
			},
		},
	})
}

func TestAccNutanixUserGroups_WithOrgUnit(t *testing.T) {
	directoryServiceOUDistName := testVars.UserGroupWithDistinguishedName[3].DistinguishedName
	updatedDirectoryServiceUserGroupName := testVars.UserGroupWithDistinguishedName[2].DistinguishedName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixUserGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixUserGroupsConfigWithOrg(directoryServiceOUDistName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_user_group.#", "0"),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_ou.#", "1"),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_ou.0.distinguished_name", directoryServiceOUDistName),
				),
			},
			{
				Config: testAccNutanixUserGroupsConfigUpdated(updatedDirectoryServiceUserGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_user_group.#", "1"),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_ou.#", "0"),
					resource.TestCheckResourceAttr(resourceNameUserGroups, "directory_service_user_group.0.distinguished_name", updatedDirectoryServiceUserGroupName),
				),
			},
		},
	})
}

func TestAccNutanixUserGroups_DuplicateEntity(t *testing.T) {
	directoryServiceDistName := testVars.UserGroupWithDistinguishedName[0].DistinguishedName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixUserGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNutanixUserGroupsConfig(directoryServiceDistName),
				ExpectError: regexp.MustCompile("DUPLICATE_ENTITY"),
			},
		},
	})
}

func testAccCheckNutanixUserGroupsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_user_groups" {
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

func testAccNutanixUserGroupsConfig(dsuuid string) string {
	return fmt.Sprintf(`
	resource "nutanix_user_groups" "acctest-managed" {
		directory_service_user_group {
			distinguished_name = "%s"
		}
	}
`, dsuuid)
}

func testAccNutanixUserGroupsConfigUpdated(dsuuid string) string {
	return fmt.Sprintf(`
	resource "nutanix_user_groups" "acctest-managed" {
		directory_service_user_group {
			distinguished_name = "%s"
		}
	}
`, dsuuid)
}

func testAccNutanixUserGroupsConfigWithOrg(dsuuid string) string {
	return fmt.Sprintf(`
	resource "nutanix_user_groups" "acctest-managed" {
		directory_service_ou {
			distinguished_name = "%s"
		}
	}
`, dsuuid)
}
