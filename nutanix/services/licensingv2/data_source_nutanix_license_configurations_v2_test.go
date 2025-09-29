package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseConfigurations = "data.nutanix_license_configurations_v2.get_configurations"

func TestDataSourceLicenseConfigurationsV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseConfigurationsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceLicenseConfigurations, "acceptances.0.accepted_by.company_name", "Nutanix"),
					resource.TestCheckResourceAttr(datasourceLicenseConfigurations, "acceptances.0.accepted_by.job_title", "MTS"),
					resource.TestCheckResourceAttr(datasourceLicenseConfigurations, "acceptances.0.accepted_by.login_id", "admin"),
					resource.TestCheckResourceAttr(datasourceLicenseConfigurations, "acceptances.0.accepted_by.user_name", "Nutanix"),
				),
			},
		},
	})
}

func testDataSourceLicenseConfigurationsConfig() string {
	return `
	data "nutanix_license_configurations_v2" "get_configurations" {}
  `
}