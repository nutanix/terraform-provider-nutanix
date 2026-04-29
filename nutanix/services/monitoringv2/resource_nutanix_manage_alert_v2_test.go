package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameManageAlert = "nutanix_manage_alert_v2.test"

func TestAccV2NutanixManageAlertResource_Acknowledge(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testManageAlertResourceConfig_Acknowledge(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameManageAlert, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameManageAlert, "action_type", "ACKNOWLEDGE"),
				),
			},
		},
	})
}

func TestAccV2NutanixManageAlertResource_Resolve(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testManageAlertResourceConfig_Resolve(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameManageAlert, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameManageAlert, "action_type", "RESOLVE"),
				),
			},
		},
	})
}

func testManageAlertResourceConfig_Acknowledge() string {
	return `
data "nutanix_alerts_v2" "list" {}

resource "nutanix_manage_alert_v2" "test" {
  ext_id      = data.nutanix_alerts_v2.list.alerts.0.ext_id
  action_type = "ACKNOWLEDGE"
  depends_on  = [data.nutanix_alerts_v2.list]
}
`
}

func testManageAlertResourceConfig_Resolve() string {
	return `
data "nutanix_alerts_v2" "list" {}

resource "nutanix_manage_alert_v2" "test" {
  ext_id      = data.nutanix_alerts_v2.list.alerts.0.ext_id
  action_type = "RESOLVE"
  depends_on  = [data.nutanix_alerts_v2.list]
}
`
}
