package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameAlert = "data.nutanix_alert_v2.test"

func TestAccV2NutanixAlertDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAlertDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameAlert, "ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameAlert, "title"),
					resource.TestCheckResourceAttrSet(dataSourceNameAlert, "severity"),
					resource.TestCheckResourceAttrSet(dataSourceNameAlert, "creation_time"),
					resource.TestCheckResourceAttrSet(dataSourceNameAlert, "alert_type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAlert, "cluster_uuid"),
				),
			},
		},
	})
}

func testAlertDatasourceConfig() string {
	return `
data "nutanix_alerts_v2" "list" {}

data "nutanix_alert_v2" "test" {
	ext_id = data.nutanix_alerts_v2.list.alerts.0.ext_id
	depends_on = [data.nutanix_alerts_v2.list]
}
`
}
