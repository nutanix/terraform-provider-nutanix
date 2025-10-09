package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseConfigurations = "data.nutanix_license_configurations_v2.get_configurations"

func TestLicensingDataSourceLicenseConfigurationsV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseConfigurationsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.0.enforcement_policy"),
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.0.has_non_compliant_features"),
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.0.has_ultimate_trail_ended"),
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.0.is_license_check_disabled"),
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.0.is_multi_cluster"),
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.0.is_stand_by"),
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.0.license_class"),
					resource.TestCheckResourceAttrSet(datasourceLicenseConfigurations, "entities.0.post_paid_config.#"),
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