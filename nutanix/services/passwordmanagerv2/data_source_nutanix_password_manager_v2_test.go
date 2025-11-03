package passwordmanagerv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNamePasswordManager = "data.nutanix_system_user_passwords_v2.test"

func TestAccV2NutanixPasswordManagerDataSource_ListPasswordStatusOfAllSystemUsers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordManagerDataSourceListPasswordStatusOfAllSystemUsersConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkAttributeMinLength(datasourceNamePasswordManager, "passwords", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixPasswordManagerDataSource_ListPasswordStatusOfAllSystemUsersWithLimit(t *testing.T) {
	limit := 2
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordManagerDataSourceListPasswordStatusOfAllSystemUsersWithLimitConfig(limit),
				Check: resource.ComposeTestCheckFunc(
					checkAttributeLengthEqual(datasourceNamePasswordManager, "passwords", limit),
				),
			},
		},
	})
}

func TestAccV2NutanixPasswordManagerDataSource_ListPasswordStatusOfAllSystemUsersWithFilter(t *testing.T) {
	usernameFilter := "username eq 'admin'"
	systemTypePCFilter := "username eq 'admin' and systemType eq Clustermgmt.Config.SystemType'PC'"
	systemTypeAOSFilter := "username eq 'admin' and systemType eq Clustermgmt.Config.SystemType'AOS'"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordManagerDataSourceListPasswordStatusOfAllSystemUsersWithFilterConfig(usernameFilter),
				Check: resource.ComposeTestCheckFunc(
					checkAttributeMinLength(datasourceNamePasswordManager, "passwords", 1),
					resource.TestCheckResourceAttr(datasourceNamePasswordManager, "passwords.0.username", "admin"),
				),
			},
			{
				Config: testAccPasswordManagerDataSourceListPasswordStatusOfAllSystemUsersWithFilterConfig(systemTypePCFilter),
				Check: resource.ComposeTestCheckFunc(
					checkAttributeMinLength(datasourceNamePasswordManager, "passwords", 1),
					resource.TestCheckResourceAttr(datasourceNamePasswordManager, "passwords.0.username", "admin"),
					resource.TestCheckResourceAttr(datasourceNamePasswordManager, "passwords.0.system_type", "PC"),
				),
			},
			{
				Config: testAccPasswordManagerDataSourceListPasswordStatusOfAllSystemUsersWithFilterConfig(systemTypeAOSFilter),
				Check: resource.ComposeTestCheckFunc(
					checkAttributeMinLength(datasourceNamePasswordManager, "passwords", 1),
					resource.TestCheckResourceAttr(datasourceNamePasswordManager, "passwords.0.username", "admin"),
					resource.TestCheckResourceAttr(datasourceNamePasswordManager, "passwords.0.system_type", "AOS"),
				),
			},
		},
	})
}

func testAccPasswordManagerDataSourceListPasswordStatusOfAllSystemUsersConfig() string {
	return `
data "nutanix_system_user_passwords_v2" "test" {}
`
}

func testAccPasswordManagerDataSourceListPasswordStatusOfAllSystemUsersWithLimitConfig(limit int) string {
	return fmt.Sprintf(`
data "nutanix_system_user_passwords_v2" "test" {
  limit = %d
}
`, limit)
}

func testAccPasswordManagerDataSourceListPasswordStatusOfAllSystemUsersWithFilterConfig(filter string) string {
	return fmt.Sprintf(`
data "nutanix_system_user_passwords_v2" "test" {
  filter = "%s"
}
`, filter)
}
