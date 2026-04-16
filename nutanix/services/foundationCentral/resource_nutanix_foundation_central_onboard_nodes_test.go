package fc_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFCOnboardNodesResource(t *testing.T) {
	// Since the acceptance test environment most likely doesn't have any Intersight managed nodes
	// to onboard, we test that an onboarding config has a non-empty plan and correctly produces
	// a Node not found error

	nodeSerial := foundationVars.OnboardNodes[0].NodeSerial

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testFCOnboardNodesResource(nodeSerial),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      testFCOnboardNodesResource(nodeSerial),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Node not found with serial %[1]s`, nodeSerial)),
			},
		},
	})
}

func testFCOnboardNodesResource(nodeSerial string) string {
	return fmt.Sprintf(`
		resource "nutanix_foundation_central_onboard_nodes" "node" {
			node_serial = "%[1]s"
		}
  `, nodeSerial)
}
