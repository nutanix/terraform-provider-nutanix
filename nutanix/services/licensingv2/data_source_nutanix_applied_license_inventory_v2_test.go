package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceLicenseInventory = "data.nutanix_applied_license_inventory_v2.get_inventory"

func TestDataSourceLicenseInventoryV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLicenseInventoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceLicenseInventory, "acceptances.0.accepted_by.company_name", "Nutanix"),
					resource.TestCheckResourceAttr(datasourceLicenseInventory, "acceptances.0.accepted_by.job_title", "MTS"),
					resource.TestCheckResourceAttr(datasourceLicenseInventory, "acceptances.0.accepted_by.login_id", "admin"),
					resource.TestCheckResourceAttr(datasourceLicenseInventory, "acceptances.0.accepted_by.user_name", "Nutanix"),
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