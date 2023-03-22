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

func TestAccEraTagsDataSource_WithEntityType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagsDataSourceConfigWithEntityType(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.values"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.required"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.required"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_tags.tags", "tags.0.entity_type", "DATABASE"),
				),
			},
		},
	})
}

func testAccEraTagsDataSourceConfig() string {
	return `
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "test-tag"
			description = "test tag description"
			entity_type = "DATABASE"
			required = false
		}

		data "nutanix_ndb_tags" "tags"{ }
	`
}

func testAccEraTagsDataSourceConfigWithEntityType() string {
	return `
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "test-tag"
			description = "test tag description"
			entity_type = "DATABASE"
			required = false
		}

		data "nutanix_ndb_tags" "tags"{ 
			entity_type= "DATABASE"
		}
	`
}
