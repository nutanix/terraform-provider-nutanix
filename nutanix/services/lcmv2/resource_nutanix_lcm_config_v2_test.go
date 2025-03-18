package lcmv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccV2NutanixLcmConfigUpdate(t *testing.T) {
	resourceLcmConfig := "nutanix_lcm_config_v2.update_lcm_config"
	datasourceLcmConfigAfterUpdate := "data.nutanix_lcm_config_v2.get_lcm_config_after_update"
	datasourceLcmConfigBeforeUpdate := "data.nutanix_lcm_config_v2.get_lcm_config_before_update"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLcmUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Check LCM Config before Update
					resource.TestCheckResourceAttr(datasourceLcmConfigBeforeUpdate, "is_auto_inventory_enabled", "false"),

					// Update LCM Config Check
					resource.TestCheckResourceAttr(resourceLcmConfig, "is_auto_inventory_enabled", "true"),

					// Check Auto Inventory is enabled after Update
					resource.TestCheckResourceAttr(datasourceLcmConfigAfterUpdate, "is_auto_inventory_enabled", "true"),

					resource.TestCheckResourceAttr(datasourceLcmConfigAfterUpdate, "auto_inventory_schedule", "16:30"),
				),
			},
		},
	})
}

func testLcmUpdateConfig() string {
	return `
# Get LCM Config before Update
data "nutanix_lcm_config_v2" "get_lcm_config_before_update" {

}

# Update LCM Config: Enable auto inventory
resource "nutanix_lcm_config_v2" "update_lcm_config" {
	is_auto_inventory_enabled = true
	auto_inventory_schedule = "16:30"
	depends_on = [data.nutanix_lcm_config_v2.get_lcm_config_before_update]
}

# Get LCM Config after Update
data "nutanix_lcm_config_v2" "get_lcm_config_after_update" {
	depends_on = [nutanix_lcm_config_v2.update_lcm_config]
}

# Update LCM Config: Enable auto inventory
resource "nutanix_lcm_config_v2" "update_lcm_config_reset_auto_inventory" {
	is_auto_inventory_enabled = false
	depends_on = [data.nutanix_lcm_config_v2.get_lcm_config_after_update]
}
`
}
