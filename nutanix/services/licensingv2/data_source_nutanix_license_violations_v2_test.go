package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseViolations = "data.nutanix_license_violations_v2.get_violations"

func TestDataSourceLicenseViolationsV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseViolationsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceGetEula, "acceptances.0.accepted_by.company_name", "Nutanix"),
					resource.TestCheckResourceAttr(datasourceGetEula, "acceptances.0.accepted_by.job_title", "MTS"),
					resource.TestCheckResourceAttr(datasourceGetEula, "acceptances.0.accepted_by.login_id", "admin"),
					resource.TestCheckResourceAttr(datasourceGetEula, "acceptances.0.accepted_by.user_name", "Nutanix"),
				),
			},
		},
	})
}

func testDataSourceLicenseViolationsConfig() string {
	return `
	data "nutanix_license_violations_v2" "get_violations" {}
  `
}