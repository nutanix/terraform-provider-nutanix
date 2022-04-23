package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFoundationNodeNetworkDetailsDataSource(t *testing.T) {
	name := "nodes"
	resourcePath := "data.nutanix_foundation_node_network_details." + name
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccFoundationPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNodeNetworkDetailsConfig(name, foundationVars.IPv6Addresses),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "nodes.#", "2"),
					resource.TestCheckResourceAttrSet(resourcePath, "nodes.0.ipmi_ip"),
					resource.TestCheckResourceAttrSet(resourcePath, "nodes.1.ipmi_ip"),
				),
			},
		},
	})
}

func testNodeNetworkDetailsConfig(name string, ipv6Addr []string) string {
	return fmt.Sprintf(`
	data "nutanix_foundation_node_network_details" "%s" {
		ipv6_addresses = ["%s", "%s"]
	}`, name, ipv6Addr[0], ipv6Addr[1])
}
