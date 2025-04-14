package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameApp = "data.nutanix_calm_app.test"

func TestCalmAppGetDatasource(t *testing.T) {
	r := acctest.RandInt()
	bp_name := "demo_bp"
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Test App created using Nutanix Terraform Plugin"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionConfig(bp_name, name, desc) + testAppReadDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameApp, "app_name", name),
					resource.TestCheckResourceAttr(datasourceNameApp, "state", "running"),
				),
			},
		},
	})
}

func testAppReadDataSourceConfig() string {
	return (`
		data "nutanix_calm_app" "test"{
			app_uuid = nutanix_calm_app_provision.test.id
		}
`)
}
