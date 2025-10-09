package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceGetEula = "data.nutanix_eula_v2.get_eula"

func TestAcceptEndUserLicenseAgreementV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLicenseAcceptEndUserLicenseAgreementConfig(),
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

func testLicenseAcceptEndUserLicenseAgreementConfig() string {
	return `
	resource "nutanix_accept_eula_v2" "accept_eula" {
		user_name   = "Nutanix"
		job_title   = "MTS"
		login_id    = "admin"
		company_name = "Nutanix"
	}
	data "nutanix_eula_v2" "get_eula" {}
  `
}