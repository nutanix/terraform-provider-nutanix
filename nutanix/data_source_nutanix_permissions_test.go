package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixPermissionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_permissions.test", "entities.#"),
					resource.TestCheckResourceAttrSet(
						"data.nutanix_permissions.test", "entities.#"),
				),
			},
		},
	})
}

func testAccPermissionsDataSourceConfig() string {
	return `
data "nutanix_permissions" "test" {}
`
}
