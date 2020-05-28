package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixSubnetsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_subnets.test", "entities.0.cluster_reference.name"),
				),
			},
		},
	})
}
func testAccSubnetsDataSourceConfig() string {
	return `data "nutanix_subnets" "test" {}`
}
