package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseViolations = "data.nutanix_license_violations_v2.get_violations"

func TestLicensingDataSourceLicenseViolationsV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseViolationsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceLicenseViolations, "entities.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseViolations, "entities.0.capacity_violations.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseViolations, "entities.0.feature_violations.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseViolations, "entities.0.expired_licenses.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseViolations, "entities.0.is_multi_cluster"),
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