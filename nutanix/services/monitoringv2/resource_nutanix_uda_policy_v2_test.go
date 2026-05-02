package monitoringv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameUdaPolicy = "nutanix_uda_policy_v2.test"

func TestAccV2NutanixUdaPolicyResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	title := fmt.Sprintf("tf-test-uda-policy-%d", r)
	updatedTitle := fmt.Sprintf("tf-test-uda-policy-%d-updated", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testUdaPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUdaPolicyResourceConfig(title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUdaPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "title", title),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "entity_type", "vm"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "trigger_conditions.0.condition.0.metric_name", "hypervisor_cpu_usage_ppm"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "trigger_conditions.0.condition.0.operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "trigger_conditions.0.condition_type", "STATIC_THRESHOLD"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "trigger_conditions.0.severity_level", "CRITICAL"),
					resource.TestCheckResourceAttrSet(resourceNameUdaPolicy, "is_enabled"),
					resource.TestCheckResourceAttrSet(resourceNameUdaPolicy, "trigger_wait_period"),
				),
			},
			{
				Config: testUdaPolicyResourceConfigUpdated(updatedTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUdaPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "title", updatedTitle),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "entity_type", "vm"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "trigger_conditions.0.condition.0.metric_name", "hypervisor_cpu_usage_ppm"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "trigger_conditions.0.condition.0.operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "trigger_conditions.0.condition_type", "STATIC_THRESHOLD"),
					resource.TestCheckResourceAttr(resourceNameUdaPolicy, "trigger_conditions.0.severity_level", "WARNING"),
				),
			},
		},
	})
}

func testUdaPolicyV2CheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_uda_policy_v2" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}
	}
	return nil
}

func testUdaPolicyResourceConfig(title string) string {
	return fmt.Sprintf(`
resource "nutanix_uda_policy_v2" "test" {
  title       = "%s"
  entity_type = "vm"
  description = "test uda policy"
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
`, title)
}

func testUdaPolicyResourceConfigUpdated(title string) string {
	return fmt.Sprintf(`
resource "nutanix_uda_policy_v2" "test" {
  title       = "%s"
  entity_type = "vm"
  description = "updated description"
  is_enabled  = true
  trigger_wait_period = 600

  trigger_conditions {
    condition {
      metric_name = "hypervisor_cpu_usage_ppm"
      operator    = "GREATER_THAN"
      threshold_value {
        int_value = 800000
      }
    }
    condition_type = "STATIC_THRESHOLD"
    severity_level = "WARNING"
  }
}
`, title)
}
