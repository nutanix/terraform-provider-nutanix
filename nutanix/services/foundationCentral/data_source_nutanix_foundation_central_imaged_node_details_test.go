package fc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFCNodeDetailsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCNodeDetailsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_imaged_nodes_list.cls", "imaged_nodes.#"),
				),
			},
		},
	})
}

func TestAccFCNodeDetailsDataSource_NodeUUID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCNodeDetailsDataSourceConfigWithUUID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_imaged_nodes_list.cls", "imaged_nodes.#"),
					resource.TestCheckResourceAttr("data.nutanix_foundation_central_imaged_node_details.k1", "cvm_vlan_id", "0"),
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_imaged_node_details.k1", "imaged_node_uuid"),
				),
			},
		},
	})
}

func testAccFCNodeDetailsDataSourceConfig() string {
	return `
	data "nutanix_foundation_central_imaged_nodes_list" "cls" {}
	`
}

func testAccFCNodeDetailsDataSourceConfigWithUUID() string {
	return `
	data "nutanix_foundation_central_imaged_nodes_list" "cls" {}
	
	data "nutanix_foundation_central_imaged_node_details" "k1"{
		imaged_node_uuid = "${data.nutanix_foundation_central_imaged_nodes_list.cls.imaged_nodes[0].imaged_node_uuid}"
	}
	`
}
