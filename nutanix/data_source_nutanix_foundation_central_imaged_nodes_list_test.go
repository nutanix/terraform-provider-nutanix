package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFCNodesListDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCNodeListDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_imaged_nodes_list.cls", "imaged_nodes.#"),
				),
			},
		},
	})
}

func testAccFCNodeListDataSourceConfig() string {
	return `
	data "nutanix_foundation_central_imaged_nodes_list" "cls" {}
	`
}
