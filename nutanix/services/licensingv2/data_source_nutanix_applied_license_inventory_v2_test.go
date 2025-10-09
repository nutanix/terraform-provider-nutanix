package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseInventory = "data.nutanix_applied_license_inventory_v2.get_inventory"

func TestLicensingDataSourceLicenseInventoryV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseInventoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceLicenseInventory, "entities.#"),
				),
			},
		},
	})
}

func  testDataSourceLicenseInventoryConfig() string {
	return `
	data "nutanix_applied_license_inventory_v2" "get_inventory" {}
  `
}