package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameApp = "nutanix_self_service_app_provision.test"

func TestAccNutanixCalmAppProvisionResource_Launch(t *testing.T) {
	r := acctest.RandInt()
	blueprintName := testVars.SelfService.BlueprintName
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionConfig(blueprintName, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "app_description", desc),
				),
			},
		},
	})
}

func TestAccNutanixCalmAppProvisionResource_SystemAction(t *testing.T) {
	r := acctest.RandInt()
	blueprintName := testVars.SelfService.BlueprintName
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"
	systemActionStop := "stop"
	systemActionStart := "start"
	systemActionRestart := "restart"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionConfig(blueprintName, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "app_description", desc),
				),
			},
			{
				Config: testCalmAppExecuteAction(name, systemActionStop),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "action", systemActionStop),
					resource.TestCheckResourceAttr(resourceNameApp, "state", "stopped"),
				),
			},
			{
				Config: testCalmAppExecuteAction(name, systemActionStart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "action", systemActionStart),
					resource.TestCheckResourceAttr(resourceNameApp, "state", "running"),
				),
			},
			{
				Config: testCalmAppExecuteAction(name, systemActionRestart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "action", systemActionRestart),
					resource.TestCheckResourceAttr(resourceNameApp, "state", "running"),
				),
			},
		},
	})
}

func TestAccNutanixCalmAppProvisionResource_SoftDelete(t *testing.T) {
	r := acctest.RandInt()
	blueprintName := testVars.SelfService.BlueprintName
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionConfig(blueprintName, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "app_description", desc),
				),
			},
			{
				Config: testCalmAppExecuteSoftDelete(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "soft_delete", "true"),
				),
			},
		},
	})
}

func testCalmAppProvisionConfig(blueprintName, name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_provision" "test" {
		bp_name         = "%[1]s"
		app_name        = "%[2]s"
		app_description = "%[3]s"
		}
`, blueprintName, name, desc)
}

func testCalmAppExecuteAction(name, systemAction string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_provision" "test" {
		app_name        = "%[1]s"
		action = "%[2]s"
		}
`, name, systemAction)
}

func testCalmAppExecuteSoftDelete(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_provision" "test" {
		app_name        = "%[1]s"
		soft_delete     = true
		}
`, name)
}
