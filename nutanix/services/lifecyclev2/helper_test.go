package lifecyclev2_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// testAccCheckNutanixNodeDestroy verifies that a node resource no longer exists after destroy.
// The GetNodeById API is expected to return an error once the node has been removed from
// Foundation Central management.
func testAccCheckNutanixNodeDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).LifecycleAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_node_v2" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		if _, err := conn.NodesAPIInstance.GetNodeById(utils.StringPtr(rs.Primary.ID)); err == nil {
			return fmt.Errorf("node %s still exists after destroy", rs.Primary.ID)
		}
	}
	return nil
}
