package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameEvent = "data.nutanix_event_v2.test"

func TestAccV2NutanixEventDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "event_type"),
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "creation_time"),
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "cluster_uuid"),
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "affected_entities.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "classifications.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "parameters.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvent, "source_entity.#"),
				),
			},
		},
	})
}

func testAccEventDatasourceConfig() string {
	return `
data "nutanix_events_v2" "events" {
  limit = 1
}

data "nutanix_event_v2" "test" {
  ext_id = data.nutanix_events_v2.events.events[0].ext_id
  depends_on = [data.nutanix_events_v2.events]
}
`
}
