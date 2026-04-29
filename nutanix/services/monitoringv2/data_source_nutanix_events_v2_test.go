package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameEvents = "data.nutanix_events_v2.test"

func TestAccV2NutanixEventsDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixEventsDataSource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsDataSourceConfigWithLimit(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.#"),
				),
			},
		},
	})
}

func testAccEventsDataSourceConfig() string {
	return `
data "nutanix_events_v2" "test" {}
`
}

func testAccEventsDataSourceConfigWithLimit() string {
	return `
data "nutanix_events_v2" "test" {
  limit = 5
}
`
}
