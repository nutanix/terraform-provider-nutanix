package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseCompliances = "data.nutanix_license_compliances_v2.get_compliances"

func TestDataSourceLicenseCompliancesV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseCompliancesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceLicenseCompliances, "acceptances.0.accepted_by.company_name", "Nutanix"),
					resource.TestCheckResourceAttr(datasourceLicenseCompliances, "acceptances.0.accepted_by.job_title", "MTS"),
					resource.TestCheckResourceAttr(datasourceLicenseCompliances, "acceptances.0.accepted_by.login_id", "admin"),
					resource.TestCheckResourceAttr(datasourceLicenseCompliances, "acceptances.0.accepted_by.user_name", "Nutanix"),
				),
			},
		},
	})
}

func testDataSourceLicenseCompliancesConfig() string {
	return `
	data "nutanix_license_compliances_v2" "get_compliances" {}
  `
}