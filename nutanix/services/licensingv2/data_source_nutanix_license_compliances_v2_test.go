package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseCompliances = "data.nutanix_license_compliances_v2.get_compliances"

func TestLicensingDataSourceLicenseCompliancesV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseCompliancesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceLicenseCompliances, "entities.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseCompliances, "entities.0.services.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseCompliances, "entities.0.type"),
					resource.TestCheckResourceAttrSet(datasourceLicenseCompliances, "entities.0.is_multi_cluster"),
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