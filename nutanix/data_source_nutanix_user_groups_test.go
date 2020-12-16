package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixUserGroupsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_user_groups.test", "entities.#"),
					resource.TestCheckResourceAttr(
						"data.nutanix_user_groups.test", "entities.#", "3"),
				),
			},
		},
	})
}

func testAccUserGroupsDataSourceConfig() string {
	return `
data "nutanix_user_groups" "test" {}
`
}
