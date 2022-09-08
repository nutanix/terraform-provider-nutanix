package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraSLAsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSLAsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_era_slas.test", "slas.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_era_slas.test", "slas.0.id"),
					resource.TestCheckResourceAttrSet("data.nutanix_era_slas.test", "slas.0.name"),
				),
			},
		},
	})
}

func testAccEraSLAsDataSourceConfig() string {
	return `
		data "nutanix_era_slas" "test" { }
	`
}
