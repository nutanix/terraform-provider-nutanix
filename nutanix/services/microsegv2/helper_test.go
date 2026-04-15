package microsegv2_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameEntityGroupV2 = "nutanix_entity_group_v2.test"

func testEntityGroupV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	api := conn.MicroSegAPI.EntityGroupsAPIInstance

	for _, rs := range state.RootModule().Resources {
		if rs.Type == "nutanix_entity_group_v2" {
			_, err := api.GetEntityGroupById(utils.StringPtr(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("entity group v2 still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}
