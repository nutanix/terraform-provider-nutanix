package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixFCNodeDetailsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCNodeDetailsDataSourceConfig(),
			},
		},
	})
}

func TestAccNutanixFCNodeDetailsDataSource_NodeUUID(t *testing.T) {
	// apiKeyName := acctest.RandomWithPrefix("test-key")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCNodeDetailsDataSourceConfigWithUUID(),
				Check: resource.ComposeTestCheckFunc(
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
