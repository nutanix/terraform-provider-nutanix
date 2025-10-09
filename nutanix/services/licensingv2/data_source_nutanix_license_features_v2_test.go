package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseFeatures = "data.nutanix_license_features_v2.get_features"

func TestLicensingDataSourceLicenseFeaturesV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseFeaturesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceLicenseFeatures, "entities.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseFeatures, "entities.0.license_category"),
					resource.TestCheckResourceAttrSet(datasourceLicenseFeatures, "entities.0.license_type"),
					resource.TestCheckResourceAttrSet(datasourceLicenseFeatures, "entities.0.name"),
					resource.TestCheckResourceAttrSet(datasourceLicenseFeatures, "entities.0.scope"),
					resource.TestCheckResourceAttrSet(datasourceLicenseFeatures, "entities.0.value"),
					resource.TestCheckResourceAttrSet(datasourceLicenseFeatures, "entities.0.value_type"),
				),
			},
		},
	})
}

func testDataSourceLicenseFeaturesConfig() string {
	return `
	data "nutanix_license_features_v2" "get_features" {}
  `
}