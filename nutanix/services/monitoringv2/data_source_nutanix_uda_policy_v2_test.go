package monitoringv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameUdaPolicy = "data.nutanix_uda_policy_v2.test"

func TestAccV2NutanixUdaPolicyDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	title := fmt.Sprintf("tf-test-uda-policy-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUdaPolicyDatasourceConfig(title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUdaPolicy, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameUdaPolicy, "title", title),
					resource.TestCheckResourceAttr(datasourceNameUdaPolicy, "entity_type", "vm"),
					resource.TestCheckResourceAttr(datasourceNameUdaPolicy, "trigger_conditions.0.condition.0.metric_name", "hypervisor_cpu_usage_ppm"),
					resource.TestCheckResourceAttr(datasourceNameUdaPolicy, "trigger_conditions.0.condition.0.operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr(datasourceNameUdaPolicy, "trigger_conditions.0.condition_type", "STATIC_THRESHOLD"),
					resource.TestCheckResourceAttr(datasourceNameUdaPolicy, "trigger_conditions.0.severity_level", "CRITICAL"),
				),
			},
		},
	})
}

func testUdaPolicyDatasourceConfig(title string) string {
	return fmt.Sprintf(`
resource "nutanix_uda_policy_v2" "ds_test" {
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

data "nutanix_uda_policy_v2" "test" {
  ext_id = nutanix_uda_policy_v2.ds_test.id
}
`, title)
}
