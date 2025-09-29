package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseEntitlements = "data.nutanix_license_entitlements_v2.get_entitlements"

func TestDataSourceLicenseEntitlementsV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseEntitlementsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceLicenseEntitlements, "acceptances.0.accepted_by.company_name", "Nutanix"),
					resource.TestCheckResourceAttr(datasourceLicenseEntitlements, "acceptances.0.accepted_by.job_title", "MTS"),
					resource.TestCheckResourceAttr(datasourceLicenseEntitlements, "acceptances.0.accepted_by.login_id", "admin"),
					resource.TestCheckResourceAttr(datasourceLicenseEntitlements, "acceptances.0.accepted_by.user_name", "Nutanix"),
				),
			},
		},
	})
}

func testDataSourceLicenseEntitlementsConfig() string {
	return `
	data "nutanix_license_entitlements_v2" "get_entitlements" {}
  `
}