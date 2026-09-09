package lifecyclev2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_refresh_connection_v2 (action resource)
// ------------------------------------------------------------
// Refresh is a one-shot mutating action that refreshes nodes or resource pools from
// a hardware provider connection. The test creates a connection first, then triggers
// a refresh_resources_spec refresh and asserts the action completes.
//
//  1. TestAccV2NutanixRefreshConnectionResource_Resources: refresh IP/MAC/server-identity
//     pools for a freshly-created connection.

func TestAccV2NutanixRefreshConnectionResource_Resources(t *testing.T) {
	name := fmt.Sprintf("tf-test-conn-refresh-%d", acctest.RandInt())
	refreshResourceName := "nutanix_refresh_connection_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccConnectionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRefreshConnectionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(refreshResourceName, "id"),
					resource.TestCheckResourceAttr(refreshResourceName, "refresh_resources_spec.0.should_refresh_ip_pools", "true"),
				),
			},
		},
	})
}

func testRefreshConnectionConfig(name string) string {
	return testConnectionConfigBasic(name, "us-west", "refresh-user") + fmt.Sprintf(`
resource "nutanix_refresh_connection_v2" "test" {
  hardware_provider_ext_id = "%[1]s"
  ext_id                   = nutanix_connection_v2.test.ext_id
  refresh_resources_spec {
    should_refresh_ip_pools              = true
    should_refresh_mac_pools             = true
    should_refresh_server_identity_pools = true
  }
}
`, testVars.Lifecycle.HardwareProviders.HardwareProviderExtID)
}
