package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixAccessControlPolicysDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessControlPolicysDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_access_control_policies.test", "entities.0.name"),
				),
			},
		},
	})
}
func testAccAccessControlPolicysDataSourceConfig() string {
	return `data "nutanix_access_control_policies" "test" {}`
}
