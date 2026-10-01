package lifecyclev2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan (ConfigureNode action resource):
//   - ConfigureServer variant: configure a discovered node in the hardware provider.
//   - UnConfigureServer variant: remove the configuration from a node.
//   - InvalidFecMode: reject an invalid network_adaptor_fec_mode enum (no live infra needed).
//
// Configure tests require a real discovered node ext_id from test_config_v2.json.

func TestAccV2NutanixConfigureNodeResource_ConfigureServer(t *testing.T) {
	nodeExtID := testVars.Lifecycle.Node.ExtID
	if nodeExtID == "" {
		t.Skip("Skipping configure node test: lifecycle.node.ext_id not set in test_config_v2.json")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigureNodeConfigureServerConfig(nodeExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_configure_node_v2.test", "configure_server.0.node_ext_id", nodeExtID),
					resource.TestCheckResourceAttrSet("nutanix_configure_node_v2.test", "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixConfigureNodeResource_UnConfigureServer(t *testing.T) {
	nodeExtID := testVars.Lifecycle.Node.ExtID
	if nodeExtID == "" {
		t.Skip("Skipping unconfigure node test: lifecycle.node.ext_id not set in test_config_v2.json")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigureNodeUnConfigureServerConfig(nodeExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_configure_node_v2.test", "un_configure_server.0.node_ext_id", nodeExtID),
				),
			},
		},
	})
}

// TestAccV2NutanixConfigureNodeResource_InvalidFecMode validates schema-level enum rejection.
func TestAccV2NutanixConfigureNodeResource_InvalidFecMode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testConfigureNodeInvalidFecModeConfig(),
				ExpectError: regexp.MustCompile(`expected .* to be one of`),
			},
		},
	})
}

func testConfigureNodeConfigureServerConfig(nodeExtID string) string {
	return fmt.Sprintf(`
resource "nutanix_configure_node_v2" "test" {
  configure_server {
    node_ext_id = "%[1]s"

    server_settings_config {
      network_adaptor_fec_mode = "OFF"
    }
  }
}
`, nodeExtID)
}

func testConfigureNodeUnConfigureServerConfig(nodeExtID string) string {
	return fmt.Sprintf(`
resource "nutanix_configure_node_v2" "test" {
  un_configure_server {
    node_ext_id = "%[1]s"
  }
}
`, nodeExtID)
}

func testConfigureNodeInvalidFecModeConfig() string {
	return `
resource "nutanix_configure_node_v2" "test" {
  configure_server {
    node_ext_id = "00000000-0000-0000-0000-000000000000"

    server_settings_config {
      network_adaptor_fec_mode = "INVALID"
    }
  }
}
`
}
