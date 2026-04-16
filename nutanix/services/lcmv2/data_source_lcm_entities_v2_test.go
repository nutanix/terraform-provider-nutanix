package lcmv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameLcmEntities = "data.nutanix_lcm_entities_v2.lcm-entities"

func TestAccV2NutanixLcmEntitiesDatasource_InvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLcmEntitiesDatasourceV4WithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameLcmEntities, "entities.#"),
					resource.TestCheckResourceAttr(datasourceNameLcmEntities, "entities.#", "0"),
				),
			},
		},
	})
}

func testLcmEntitiesDatasourceV4WithInvalidFilterConfig() string {
	return `
	data "nutanix_lcm_entities_v2" "lcm-entities" {
		filter = "entityModel eq 'invalid_model'"
	}
	`
}
