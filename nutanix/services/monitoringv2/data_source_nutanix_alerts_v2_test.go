package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameAlerts = "data.nutanix_alerts_v2.test"

func TestAccV2NutanixAlertsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAlertsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameAlerts, "alerts.#"),
				),
			},
		},
	})
}

func testAlertsDatasourceConfig() string {
	return `
data "nutanix_alerts_v2" "test" {}
`
}
