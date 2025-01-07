package foundation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFoundationDiscoverNodesDataSource(t *testing.T) {
	name := "nodes"
	resourcePath := "data.nutanix_foundation_discover_nodes." + name
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
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
