package lifecyclev2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan (RefreshNode action resource):
//   - RefreshNode_Basic: refresh a discovered node's information from the hardware provider.
//
// Refresh requires a real discovered node ext_id from test_config_v2.json.

func TestAccV2NutanixRefreshNodeResource_Basic(t *testing.T) {
	nodeExtID := testVars.Lifecycle.Node.ExtID
	if nodeExtID == "" {
		t.Skip("Skipping refresh node test: lifecycle.node.ext_id not set in test_config_v2.json")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRefreshNodeConfig(nodeExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_refresh_node_v2.test", "ext_id", nodeExtID),
				),
			},
		},
	})
}

func testRefreshNodeConfig(nodeExtID string) string {
	return fmt.Sprintf(`
resource "nutanix_refresh_node_v2" "test" {
  ext_id = "%[1]s"
}
`, nodeExtID)
}
