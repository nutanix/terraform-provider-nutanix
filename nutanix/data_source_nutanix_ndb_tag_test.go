package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraTagDataSource_basic(t *testing.T) {
	r := randIntBetween(10, 20)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagDataSourceConfig(r),
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
	r := randIntBetween(21, 30)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagDataSourceConfigByName(r),
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

func testAccEraTagDataSourceConfig(r int) string {
	return fmt.Sprintf(`
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "test-tag-%[1]d"
			description = "test tag description"
			entity_type = "DATABASE"
			required = false
		}

		data "nutanix_ndb_tags" "tags"{ }

		data "nutanix_ndb_tag" "tag"{
			id = data.nutanix_ndb_tags.tags.tags.0.id
		}
	`, r)
}

func testAccEraTagDataSourceConfigByName(r int) string {
	return fmt.Sprintf(`
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "test-tag-name-%[1]d"
			description = "test tag description"
			entity_type = "DATABASE"
			required = false
		}

		data "nutanix_ndb_tags" "tags"{ }

		data "nutanix_ndb_tag" "tag"{
			name = data.nutanix_ndb_tags.tags.tags.0.name
		}
	`, r)
}
