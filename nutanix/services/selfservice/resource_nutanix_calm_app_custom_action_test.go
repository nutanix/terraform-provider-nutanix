package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameAction = "nutanix_calm_app_custom_action.test"

func TestAccNutanixCalmAppResource_CustomAction(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"
	action_name := "custom1"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppRunCustomAction(name, desc, action_name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAction, "action_name", action_name),
				),
			},
		},
	})
}

func testCalmAppRunCustomAction(name, desc, action_name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_provision" "test" {
		bp_name         = "demo_bp"
		app_name        = "%[1]s"
		app_description = "%[2]s"
		}

		resource "nutanix_calm_app_custom_action" "test" {
		app_name        = nutanix_calm_app_provision.test.app_name
		action_name = "%[3]s"
		}
`, name, desc, action_name)
}
