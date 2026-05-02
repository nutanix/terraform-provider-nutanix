package monitoringv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameFindConflicting = "nutanix_find_conflicting_uda_policies_v2.test"

func TestAccV2NutanixFindConflictingUdaPoliciesResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	title := fmt.Sprintf("tf-test-uda-conflict-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFindConflictingUdaPoliciesResourceConfig(title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameFindConflicting, "conflicting_policies.#"),
				),
			},
		},
	})
}

func testFindConflictingUdaPoliciesResourceConfig(title string) string {
	return fmt.Sprintf(`
resource "nutanix_find_conflicting_uda_policies_v2" "test" {
  title       = "%s"
  entity_type = "vm"

  trigger_conditions {
    condition {
      metric_name = "hypervisor_cpu_usage_ppm"
      operator    = "GREATER_THAN"
      threshold_value {
        int_value = 900000
      }
    }
    condition_type = "STATIC_THRESHOLD"
    severity_level = "CRITICAL"
  }
}
`, title)
}
