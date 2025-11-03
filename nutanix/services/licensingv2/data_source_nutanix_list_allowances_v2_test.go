package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceListAllowances = "data.nutanix_license_violations_v2.get_violations"

func TestLicensingDataSourceListAllowancesV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceListAllowancesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceListAllowances, "entities.#"),
					resource.TestCheckResourceAttrSet(datasourceListAllowances, "entities.0.feature_id"),
					resource.TestCheckResourceAttrSet(datasourceListAllowances, "entities.0.scope"),
					resource.TestCheckResourceAttrSet(datasourceListAllowances, "entities.0.value"),
					resource.TestCheckResourceAttrSet(datasourceListAllowances, "entities.0.value_type"),
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