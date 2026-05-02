package monitoringv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameUdaPolicies = "data.nutanix_uda_policies_v2.test"

func TestAccV2NutanixUdaPoliciesDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	title := fmt.Sprintf("tf-test-uda-policy-list-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUdaPoliciesDatasourceConfig(title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUdaPolicies, "uda_policies.#"),
				),
			},
		},
	})
}

func testUdaPoliciesDatasourceConfig(title string) string {
	return fmt.Sprintf(`
resource "nutanix_uda_policy_v2" "list_test" {
  title       = "%s"
  entity_type = "vm"
  is_enabled  = true
  trigger_wait_period = 600

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

data "nutanix_uda_policies_v2" "test" {
  depends_on = [nutanix_uda_policy_v2.list_test]
}
`, title)
}
