package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraTagsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.values"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.required"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.required"),
				),
			},
		},
	})
}

func testAccEraTagsDataSourceConfig() string {
	return `
		data "nutanix_ndb_tags" "tags"{ }
	`
}
