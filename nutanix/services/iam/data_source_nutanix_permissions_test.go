package iam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixPermissionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
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
