package clustersv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameHostEntities = "data.nutanix_hosts_v2.test"

func TestAccNutanixHostEntitiesV2Datasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testHostEntitiesDatasourceV4Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameHostEntities, "host_entities.#"),
					resource.TestCheckResourceAttrSet(datasourceNameHostEntities, "host_entities.0.ext_id"),
				),
			},
		},
	})
}

func TestAccNutanixHostEntitiesV2Datasource_WithLimit(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testHostEntitiesDatasourceV4WithLimitConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameHostEntities, "host_entities.#"),
					resource.TestCheckResourceAttrSet(datasourceNameHostEntities, "host_entities.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameHostEntities, "host_entities.#", "1"),
				),
			},
		},
	})
}

func testHostEntitiesDatasourceV4Config() string {
	return `
	data "nutanix_hosts_v2" "test"{}
	`
}

func testHostEntitiesDatasourceV4WithLimitConfig() string {
	return `
		data "nutanix_hosts_v2" "test" {
			limit     = 1
		}
	`
}
