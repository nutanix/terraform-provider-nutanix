package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraSLAsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSLAsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test", "slas.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test", "slas.0.id"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test", "slas.0.name"),
				),
			},
		},
	})
}

func testAccEraSLAsDataSourceConfig() string {
	return `
		data "nutanix_ndb_slas" "test" { }
	`
}
