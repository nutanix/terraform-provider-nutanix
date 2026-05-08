package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameEvents = "data.nutanix_events_v2.test"

func TestAccV2NutanixEventsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.#"),
					checkAttributeLength(datasourceNameEvents, "events", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixEventsDatasource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsDatasourceWithLimitConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.event_type"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.creation_time"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.cluster_uuid"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.affected_entities.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.classifications.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.parameters.#"),
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.0.source_entity.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixEventsDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsDatasourceWithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEvents, "events.#"),
					resource.TestCheckResourceAttr(datasourceNameEvents, "events.#", "0"),
				),
			},
		},
	})
}

func testAccEventsDatasourceConfig() string {
	return `
data "nutanix_events_v2" "test" {}
`
}

func testAccEventsDatasourceWithLimitConfig() string {
	return `
data "nutanix_events_v2" "test" {
  limit = 1
}
`
}

func testAccEventsDatasourceWithInvalidFilterConfig() string {
	return `
data "nutanix_events_v2" "test" {
  filter = "eventType eq 'NONEXISTENT_EVENT_TYPE_ZZZZZ'"
}
`
}
