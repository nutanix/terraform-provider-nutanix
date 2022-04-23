package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFoundationDiscoverNodesDataSource(t *testing.T) {
	name := "nodes"
	resourcePath := "data.nutanix_foundation_discover_nodes." + name
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccFoundationPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDiscoverNodesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourcePath, "entities.0.nodes.0.ipv6_address"),
					resource.TestCheckResourceAttrSet(resourcePath, "entities.0.nodes.0.hypervisor"),
					resource.TestCheckResourceAttrSet(resourcePath, "entities.0.block_id"),
				),
			},
		},
	})
}

func testDiscoverNodesConfig(name string) string {
	return fmt.Sprintf(`data "nutanix_foundation_discover_nodes" "%s" {}`, name)
}
