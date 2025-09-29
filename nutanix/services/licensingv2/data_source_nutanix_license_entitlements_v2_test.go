package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseEntitlements = "data.nutanix_license_entitlements_v2.get_entitlements"

func TestLicensingDataSourceLicenseEntitlementsV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseEntitlementsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.0.details.#"),
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.0.details.0.category"),
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.0.details.0.meter"),
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.0.details.0.name"),
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.0.details.0.quantity"),
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.0.details.0.scope"),
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.0.details.0.sub_category"),
					resource.TestCheckResourceAttrSet(datasourceLicenseEntitlements, "entities.0.details.0.type"),
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