package fc_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFCOnboardNodesResource(t *testing.T) {
	// Since the acceptance test environment most likely doesn't have any Intersight managed nodes
	// to onboard, we test that an onboarding config has a non-empty plan and correctly produces
	// a Node not found error
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testFCOnboardNodesResource(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      testFCOnboardNodesResource(),
				ExpectError: regexp.MustCompile(`Node not found with serial ABC12345D6E`),
			},
		},
	})
}

func testFCOnboardNodesResource() string {
	return `
		resource "nutanix_foundation_central_onboard_nodes" "node" {
			node_serial = "ABC12345D6E"
		}
  `
}
