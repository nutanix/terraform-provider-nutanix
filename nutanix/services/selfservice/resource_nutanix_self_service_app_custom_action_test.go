package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameAction = "nutanix_self_service_app_custom_action.test"

func TestAccNutanixCalmAppResource_CustomAction(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"
	actionName := "custom1"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppRunCustomAction(testVars.SelfService.BlueprintName, name, desc, actionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAction, "action_name", actionName),
				),
			},
		},
	})
}

func testCalmAppRunCustomAction(blueprintName, name, desc, actionName string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_provision" "test" {
		bp_name         = "%[1]s"
		app_name        = "%[2]s"
		app_description = "%[3]s"
		}

		resource "nutanix_self_service_app_custom_action" "test" {
		app_name        = nutanix_self_service_app_provision.test.app_name
		action_name = "%[4]s"
		}
`, blueprintName, name, desc, actionName)
}
