package prismv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameListRestorablePCs = "data.nutanix_pcs_v2.test"

func TestAccV2NutanixRestorablePcsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List Restorable pcs
			{
				Config: testAccListRestorablePCConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.#"),
					checkAttributeLength(datasourceNameListRestorablePCs, "restorable_pcs", 1),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.config.0.build_info.0.version"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.config.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.config.0.size"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.network.0.name_servers.0.ipv4.0.value"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.network.0.ntp_servers.#"),
				),
			},
		},
	})
}

func testAccListRestorablePCConfig() string {
	return `
data "nutanix_restorable_pcs_v2" "test" {}
`
}
