package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraTagDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.status"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_tag.tag", "values", "0"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_tag.tag", "status", "ENABLED"),
				),
			},
		},
	})
}

func TestAccEraTagDataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_tags.tags", "tags.0.status"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_tag.tag", "values", "0"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_tag.tag", "status", "ENABLED"),
				),
			},
		},
	})
}

func testAccEraTagDataSourceConfig() string {
	return `
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "test-tag"
			description = "test tag description"
			entity_type = "DATABASE"
			required = false
		}

		data "nutanix_ndb_tags" "tags"{ }

		data "nutanix_ndb_tag" "tag"{
			id = data.nutanix_ndb_tags.tags.tags.0.id
		}
	`
}

func testAccEraTagDataSourceConfigByName() string {
	return `
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "test-tag-name"
			description = "test tag description"
			entity_type = "DATABASE"
			required = false
		}

		data "nutanix_ndb_tags" "tags"{ }

		data "nutanix_ndb_tag" "tag"{
			name = data.nutanix_ndb_tags.tags.tags.0.name
		}
	`
}
