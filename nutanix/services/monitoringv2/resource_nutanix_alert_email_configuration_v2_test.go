package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameAlertEmailConfig = "nutanix_alert_email_configuration_v2.test"

func TestAccV2NutanixAlertEmailConfigurationResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAlertEmailConfigurationResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAlertEmailConfig, "is_enabled", "true"),
				),
			},
		},
	})
}

func TestAccV2NutanixAlertEmailConfigurationResource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAlertEmailConfigurationResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAlertEmailConfig, "is_enabled", "true"),
				),
			},
			{
				Config: testAlertEmailConfigurationResourceConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAlertEmailConfig, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameAlertEmailConfig, "is_email_digest_enabled", "true"),
				),
			},
		},
	})
}

func testAlertEmailConfigurationResourceConfig() string {
	return `
resource "nutanix_alert_email_configuration_v2" "test" {
  is_enabled = true
}
`
}

func testAlertEmailConfigurationResourceConfigUpdated() string {
	return `
resource "nutanix_alert_email_configuration_v2" "test" {
  is_enabled             = true
  is_email_digest_enabled = true
}
`
}
