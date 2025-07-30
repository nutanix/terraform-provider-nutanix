package passwordmanagerv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamePasswordManager = "nutanix_password_change_request_v2.test"

func TestAccV2NutanixPasswordManagerResource_UpdatePasswordForAdminPCUser(t *testing.T) {
	nutanixPCuserFilter := "username eq 'nutanix' and systemType eq Clustermgmt.Config.SystemType'PC'"

	passwords := []string{
		"nutanix/4u", // Initial password
		"x.K.2.$.j.$.l.0",
		"sW3*Hj8%Gp2(",
		"tR6@Vz1#Hn5$",
		"B8!cK2*Wx4%",
		"j.Y.4.$.9.M.f.1",
		"nM2^vC7*Qs4(",
		"R5@hY1!dUo6%",
		"a.Z.9.@.S.t.p",
		"gF7!mK2#bW9@",
		"Pq#4Zx8&Lt3$",
	}
	// "D4&fQ9^mZ7!",

	// Build a sequence of steps: each step updates from passwords[i] to passwords[i+1]
	var steps []resource.TestStep
	for i := 0; i < len(passwords); i++ {
		current := passwords[i]
		next := passwords[(i+1)%len(passwords)]
		steps = append(steps, resource.TestStep{
			PreConfig: func() {
				fmt.Printf("Step %d : Updating password from '%s' to '%s'\n", i, current, next)
			},
			Config: testAccPasswordManagerResourceUpdatePasswordForAdminPCUserConfig(nutanixPCuserFilter, current, next),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("nutanix_password_change_request_v2.test", "ext_id"),
			),
		})
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps:     steps,
	})
}

func testAccPasswordManagerResourceUpdatePasswordForAdminPCUserConfig(filter, current_password, next_password string) string {
	return fmt.Sprintf(`

provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
}
data "nutanix_system_user_passwords_v2" "test" {
	filter = "%[1]s"
}

resource "nutanix_password_change_request_v2" "test" {
	ext_id = data.nutanix_system_user_passwords_v2.test.passwords.0.ext_id
	current_password = "%[2]s"
	new_password = "%[3]s"
}
`, filter, current_password, next_password)
}
