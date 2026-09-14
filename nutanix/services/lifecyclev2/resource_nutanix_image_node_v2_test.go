package lifecyclev2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan (ImageNode action resource):
//   - HostInstallation variant: install a patched host/hypervisor image on a node.
//   - AosInstallation variant: install AOS on a node.
//   - MissingConfiguration: no configuration block supplied -> expect an error from Create.
//
// Imaging tests require a real discovered node ext_id + image ext_ids from test_config_v2.json.

func TestAccV2NutanixImageNodeResource_HostInstallation(t *testing.T) {
	nodeExtID := testVars.Lifecycle.Node.ExtID
	patchedImageExtID := testVars.Lifecycle.Node.PatchedImageExtID
	if nodeExtID == "" || patchedImageExtID == "" {
		t.Skip("Skipping image node host installation test: lifecycle.node.ext_id / patched_image_ext_id not set in test_config_v2.json")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImageNodeHostInstallationConfig(nodeExtID, patchedImageExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_image_node_v2.test", "ext_id", nodeExtID),
					resource.TestCheckResourceAttr("nutanix_image_node_v2.test", "host_installation.0.patched_image_ext_id", patchedImageExtID),
				),
			},
		},
	})
}

func TestAccV2NutanixImageNodeResource_AosInstallation(t *testing.T) {
	nodeExtID := testVars.Lifecycle.Node.ExtID
	aosImageExtID := testVars.Lifecycle.Node.AosImageExtID
	if nodeExtID == "" || aosImageExtID == "" {
		t.Skip("Skipping image node AOS installation test: lifecycle.node.ext_id / aos_image_ext_id not set in test_config_v2.json")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImageNodeAosInstallationConfig(nodeExtID, aosImageExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_image_node_v2.test", "ext_id", nodeExtID),
					resource.TestCheckResourceAttr("nutanix_image_node_v2.test", "aos_installation.0.aos_image_ext_id", aosImageExtID),
				),
			},
		},
	})
}

// TestAccV2NutanixImageNodeResource_MissingConfiguration ensures Create errors when no
// configuration branch is supplied.
func TestAccV2NutanixImageNodeResource_MissingConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testImageNodeMissingConfigurationConfig(),
				ExpectError: regexp.MustCompile(`exactly one configuration block`),
			},
		},
	})
}

func testImageNodeHostInstallationConfig(nodeExtID, patchedImageExtID string) string {
	return fmt.Sprintf(`
resource "nutanix_image_node_v2" "test" {
  ext_id = "%[1]s"

  host_installation {
    patched_image_ext_id = "%[2]s"
  }
}
`, nodeExtID, patchedImageExtID)
}

func testImageNodeAosInstallationConfig(nodeExtID, aosImageExtID string) string {
	return fmt.Sprintf(`
resource "nutanix_image_node_v2" "test" {
  ext_id = "%[1]s"

  aos_installation {
    aos_image_ext_id = "%[2]s"
    cvm_memory_gb    = 32

    management_network {
      vlan_id = 0

      ip {
        ipv4 {
          value         = "192.0.2.10"
          prefix_length = 24
        }
      }

      gateway {
        ipv4 {
          value         = "192.0.2.1"
          prefix_length = 24
        }
      }
    }
  }
}
`, nodeExtID, aosImageExtID)
}

func testImageNodeMissingConfigurationConfig() string {
	return `
resource "nutanix_image_node_v2" "test" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
`
}
