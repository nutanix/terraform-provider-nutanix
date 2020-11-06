package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixAccessControlPoliciesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessControlPoliciesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_access_control_policies.test", "entities.0.name"),
				),
			},
		},
	})
}
func testAccAccessControlPoliciesDataSourceConfig() string {
	return `data "nutanix_access_control_policies" "test" {}`
}
