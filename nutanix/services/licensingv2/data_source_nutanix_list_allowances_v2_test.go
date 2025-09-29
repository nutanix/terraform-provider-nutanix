package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceListAllowances = "data.nutanix_license_violations_v2.get_violations"

func TestDataSourceListAllowancesV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceListAllowancesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceListAllowances, "acceptances.0.accepted_by.company_name", "Nutanix"),
					resource.TestCheckResourceAttr(datasourceListAllowances, "acceptances.0.accepted_by.job_title", "MTS"),
					resource.TestCheckResourceAttr(datasourceListAllowances, "acceptances.0.accepted_by.login_id", "admin"),
					resource.TestCheckResourceAttr(datasourceListAllowances, "acceptances.0.accepted_by.user_name", "Nutanix"),
				),
			},
		},
	})
}

func testDataSourceListAllowancesConfig() string {
	return `
	data "nutanix_list_allowances_v2" "get_allowances" {}
  `
}