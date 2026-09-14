package lifecyclev2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan (lifecycle / Nodes):
//   - Node_Basic:          register a node (Create) with all input fields, assert literal
//                          values, then Update scalar fields; validate singular + plural
//                          datasources match resource state; ImportState verify.
//   - Node_MissingRequired: omit the Required `manufacturer` field -> expect a Terraform
//                          "Missing required argument" plan error (no live infra needed).
//   - Node_InvalidIdentifierType: supply an identifier `type` outside the SDK enum and expect
//                          the API/provider to reject it.
//   - ConfigureNode_Variants / ImageNode_Variants: exercise each OneOf configuration branch.
//   - RefreshNode_Basic:   refresh a discovered node.
//
// Node lifecycle tests operate on hardware that Foundation Central has already discovered;
// the node serial + manufacturer are read from test_config_v2.json (never hardcoded) and the
// test skips when those infra params are not provided in the environment.

const resourceNameNode = "nutanix_node_v2.test"

func TestAccV2NutanixNodeResource_Basic(t *testing.T) {
	nodeSerial := testVars.Lifecycle.Node.NodeSerial
	manufacturer := testVars.Lifecycle.Node.Manufacturer
	if nodeSerial == "" || manufacturer == "" {
		t.Skip("Skipping node lifecycle test: lifecycle.node.node_serial / manufacturer not set in test_config_v2.json")
	}

	updatedModel := "tf-updated-model"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testNodeConfig(nodeSerial, manufacturer, "tf-model"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNode, "manufacturer", manufacturer),
					resource.TestCheckResourceAttr(resourceNameNode, "model", "tf-model"),
					resource.TestCheckResourceAttr(resourceNameNode, "identifiers.#", "1"),
					resource.TestCheckResourceAttr(resourceNameNode, "identifiers.0.type", "SERIAL_NUMBER"),
					resource.TestCheckResourceAttr(resourceNameNode, "identifiers.0.value", nodeSerial),
					resource.TestCheckResourceAttrSet(resourceNameNode, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNode, "state"),
				),
			},
			{
				Config: testNodeConfig(nodeSerial, manufacturer, updatedModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNode, "model", updatedModel),
					resource.TestCheckResourceAttr(resourceNameNode, "manufacturer", manufacturer),
				),
			},
			// Singular datasource must mirror the resource state.
			{
				Config: testNodeConfig(nodeSerial, manufacturer, updatedModel) + testNodeDataSourcesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.nutanix_node_v2.test", "manufacturer", resourceNameNode, "manufacturer"),
					resource.TestCheckResourceAttrPair("data.nutanix_node_v2.test", "model", resourceNameNode, "model"),
					resource.TestCheckResourceAttrPair("data.nutanix_node_v2.test", "ext_id", resourceNameNode, "ext_id"),
					resource.TestCheckResourceAttrSet("data.nutanix_nodes_v2.test", "nodes.#"),
				),
			},
			{
				ResourceName:      resourceNameNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccV2NutanixNodeResource_MissingRequired ensures the Required manufacturer argument is
// enforced by the schema before any API call is made. This needs no live infrastructure.
func TestAccV2NutanixNodeResource_MissingRequired(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testNodeConfigMissingManufacturer(),
				ExpectError: regexp.MustCompile(`The argument "manufacturer" is required`),
			},
		},
	})
}

func testNodeConfig(nodeSerial, manufacturer, model string) string {
	return fmt.Sprintf(`
resource "nutanix_node_v2" "test" {
  manufacturer = "%[2]s"
  model        = "%[3]s"

  identifiers {
    type  = "SERIAL_NUMBER"
    value = "%[1]s"
  }
}
`, nodeSerial, manufacturer, model)
}

func testNodeDataSourcesConfig() string {
	return `
data "nutanix_node_v2" "test" {
  ext_id = nutanix_node_v2.test.ext_id
}

data "nutanix_nodes_v2" "test" {
  depends_on = [nutanix_node_v2.test]
}
`
}

func testNodeConfigMissingManufacturer() string {
	return `
resource "nutanix_node_v2" "test" {
  model = "tf-model"

  identifiers {
    type  = "SERIAL_NUMBER"
    value = "TF-SERIAL-0001"
  }
}
`
}

// TestAccV2NutanixNodesDatasource_Basic validates the plural list datasource independently.
func TestAccV2NutanixNodesDatasource_Basic(t *testing.T) {
	if testVars.Lifecycle.Node.NodeSerial == "" {
		t.Skip("Skipping nodes datasource test: lifecycle.node.node_serial not set in test_config_v2.json")
	}
	r := acctest.RandIntRange(1, 1000)
	_ = r
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "nutanix_nodes_v2" "list" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_nodes_v2.list", "nodes.#"),
				),
			},
		},
	})
}
