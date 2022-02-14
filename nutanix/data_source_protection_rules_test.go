package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixProtectionRulesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProtectionRulesDataSourceConfig(),
			},
		},
	})
}

func testAccProtectionRulesDataSourceConfig() string {
	return `
		data "nutanix_protection_rules" "test" {}
	`
}
