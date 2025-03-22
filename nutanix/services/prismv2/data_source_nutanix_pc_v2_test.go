package prismv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameListPC = "data.nutanix_pc_v2.test"

func TestAccV2NutanixPcDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List pcs
			{
				Config: testAccFetchPCConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameListPC, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameListPC, "config.0.build_info.0.version"),
					resource.TestCheckResourceAttrSet(datasourceNameListPC, "config.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameListPC, "config.0.size"),
					resource.TestCheckResourceAttrSet(datasourceNameListPC, "network.0.name_servers.0.ipv4.0.value"),
					resource.TestCheckResourceAttrSet(datasourceNameListPC, "network.0.ntp_servers.#"),
				),
			},
		},
	})
}

func testAccFetchPCConfig() string {
	return `
data "nutanix_pcs_v2" "test" {}

data "nutanix_pc_v2" "test" {
  ext_id = data.nutanix_pcs_v2.test.pcs.0.ext_id
}


`
}
