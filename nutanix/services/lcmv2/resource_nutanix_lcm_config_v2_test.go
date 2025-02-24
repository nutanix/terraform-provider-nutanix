package lcmv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccV2NutanixLcmConfigUpdate(t *testing.T) {
	datasourceLcmConfigBeforeUpdate := "data.nutanix_lcm_config_v2.get_lcm_config_before_update"
	resourceLcmConfig := "nutanix_lcm_config_v2.update_lcm_config"
	datasourceLcmConfigAfterUpdate := "data.nutanix_lcm_config_v2.get_lcm_config_after_update"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLcmUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Check Auto Inventory is disabled Before Update
					resource.TestCheckResourceAttr(datasourceLcmConfigBeforeUpdate, "is_auto_inventory_enabled", "false"),

					// Update LCM Config Check
					resource.TestCheckResourceAttr(resourceLcmConfig, "is_auto_inventory_enabled", "true"),

					// Check Auto Inventory is enabled after Update
					resource.TestCheckResourceAttr(datasourceLcmConfigAfterUpdate, "is_auto_inventory_enabled", "true"),
				),
			},
		},
	})
}

func testLcmUpdateConfig() string {
	return `
# Get LCM Config before Update
data "nutanix_lcm_config_v2" "get_lcm_config_before_update" {
	lifecycle {
		postcondition {
			condition     = self.is_auto_inventory_enabled == false
			error_message = "Auto Inventory is already enabled, current value: ${self.is_auto_inventory_enabled}"
	   }
  	}
}

# Update LCM Config: Enable auto inventory
resource "nutanix_lcm_config_v2" "update_lcm_config" {
	is_auto_inventory_enabled = true
	depends_on = [data.nutanix_lcm_config_v2.get_lcm_config_before_update]
}

# Get LCM Config after Update
data "nutanix_lcm_config_v2" "get_lcm_config_after_update" {
   lifecycle {
		postcondition {
			condition     = self.is_auto_inventory_enabled == true
			error_message = "Auto Inventory is not enabled, current value: ${self.is_auto_inventory_enabled}"
	   }
  	}
	depends_on = [nutanix_lcm_config_v2.update_lcm_config]
}
`
}
