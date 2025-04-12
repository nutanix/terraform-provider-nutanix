package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameApp = "nutanix_calm_app_provision.test"

func TestAccNutanixCalmAppProvisionResource_Launch(t *testing.T) {
	r := acctest.RandInt()
	bp_name := "demo_bp"
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionConfig(bp_name, name, desc),
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
	bp_name := "demo_bp"
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"
	system_action1 := "stop"
	system_action2 := "start"
	system_action3 := "restart"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionConfig(bp_name, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "app_description", desc),
				),
			},
			{
				Config: testCalmAppExecuteAction(name, system_action1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "action", system_action1),
					resource.TestCheckResourceAttr(resourceNameApp, "state", "stopped"),
				),
			},
			{
				Config: testCalmAppExecuteAction(name, system_action2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "action", system_action2),
					resource.TestCheckResourceAttr(resourceNameApp, "state", "running"),
				),
			},
			{
				Config: testCalmAppExecuteAction(name, system_action3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(resourceNameApp, "action", system_action3),
					resource.TestCheckResourceAttr(resourceNameApp, "state", "running"),
				),
			},
		},
	})
}

func TestAccNutanixCalmAppProvisionResource_SoftDelete(t *testing.T) {
	r := acctest.RandInt()
	bp_name := "demo_bp"
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionConfig(bp_name, name, desc),
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

func testCalmAppProvisionConfig(bp_name, name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_provision" "test" {
		bp_name         = "%[1]s"
		app_name        = "%[2]s"
		app_description = "%[3]s"
		}
`, bp_name, name, desc)
}

func testCalmAppExecuteAction(name, system_action string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_provision" "test" {
		app_name        = "%[1]s"
		action = "%[2]s"
		}
`, name, system_action)
}

func testCalmAppExecuteSoftDelete(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_provision" "test" {
		app_name        = "%[1]s"
		soft_delete     = true
		}
`, name)
}
