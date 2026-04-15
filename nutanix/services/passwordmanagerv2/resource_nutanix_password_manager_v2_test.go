package passwordmanagerv2_test

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamePasswordManager = "nutanix_password_change_request_v2.test"

func TestAccV2NutanixPasswordManagerResource_UpdatePasswordForAdminAOSUser(t *testing.T) {
	systemTypeAOSFilter := "username eq 'admin' and systemType eq Clustermgmt.Config.SystemType'AOS'"

	passwords := []string{
		os.Getenv("NUTANIX_PASSWORD"), // Initial password
	}

	// Generate 5 new passwords:
	for i := 0; i < 5; i++ {
		pwd, err := GeneratePassword(passwords)
		if err != nil {
			log.Fatalf("could not generate password #%d: %v", i+1, err)
		}
		fmt.Printf("New password #%d: %s\n", i+1, pwd)
		passwords = append(passwords, pwd)
	}

	steps := []resource.TestStep{}

	// Build a sequence of steps: each step updates from passwords[i] to passwords[i+1]
	for i := 0; i < len(passwords); i++ {
		current := passwords[i]
		next := passwords[(i+1)%len(passwords)]

		steps = append(steps, resource.TestStep{
			PreConfig: func() {
				log.Printf("[DEBUG] Updating password from '%s' to '%s'\n", current, next)
			},
			Config: testAccPasswordManagerResourceUpdatePasswordForAdminAOSUserConfig(systemTypeAOSFilter, current, next),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceNamePasswordManager, "ext_id"),
				resource.TestCheckResourceAttr(resourceNamePasswordManager, "current_password", current),
				resource.TestCheckResourceAttr(resourceNamePasswordManager, "new_password", next),
			),
		})
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps:     steps,
	})
}

func TestAccV2NutanixPasswordManagerResource_UpdatePasswordForAdminPCUserWrongCurrentPass(t *testing.T) {
	systemTypePCFilter := "username eq 'admin' and systemType eq Clustermgmt.Config.SystemType'PC'"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordManagerResourceUpdatePasswordForAdminPCUserConfig(systemTypePCFilter, "wrong_current_password", "new_password"),
				// Error details vary across PC versions/setups (e.g. "Unauthorised (PAM authentication failed)" / "Account locked"),
				// but the task failure prefix is stable.
				ExpectError: regexp.MustCompile("Failed to change system user password due to"),
			},
		},
	},
	)
}

func TestAccV2NutanixPasswordManagerResource_UpdatePasswordForAdminAOSUserWrongCurrentPass(t *testing.T) {
	systemTypePCFilter := "username eq 'admin' and systemType eq Clustermgmt.Config.SystemType'AOS'"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordManagerResourceUpdatePasswordForAdminPCUserConfig(systemTypePCFilter, "wrong_current_password", "new_password"),
				// Error details vary across AOS versions/setups (e.g. "Password change operation failed due to some internal issue." / "RPC call for password change failed"),
				// but the task failure prefix is stable.
				ExpectError: regexp.MustCompile("Failed to change system user password due to"),
			},
		},
	},
	)
}

func testAccPasswordManagerResourceUpdatePasswordForAdminPCUserConfig(filter, currentPassword, nextPassword string) string {
	return fmt.Sprintf(`


data "nutanix_system_user_passwords_v2" "test" {
	filter = "%[1]s"
}

resource "nutanix_password_change_request_v2" "test" {
	ext_id = data.nutanix_system_user_passwords_v2.test.passwords.0.ext_id
	current_password = "%[2]s"
	new_password = "%[3]s"
}
`, filter, currentPassword, nextPassword)
}

func testAccPasswordManagerResourceUpdatePasswordForAdminAOSUserConfig(filter, currentPassword, nextPassword string) string {
	return fmt.Sprintf(`


data "nutanix_system_user_passwords_v2" "test" {
	filter = "%[1]s"
}

resource "nutanix_password_change_request_v2" "test" {
	ext_id = data.nutanix_system_user_passwords_v2.test.passwords.0.ext_id
	current_password = "%[2]s"
	new_password = "%[3]s"
}
`, filter, currentPassword, nextPassword)
}
