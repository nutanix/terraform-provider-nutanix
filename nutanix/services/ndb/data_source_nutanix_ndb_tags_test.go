package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraTagsDataSource_basic(t *testing.T) {
	r := acc.RandIntBetween(11, 15)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagsDataSourceConfig(r),
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
	r := acc.RandIntBetween(21, 25)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagsDataSourceConfigWithEntityType(r),
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

func testAccEraTagsDataSourceConfig(r int) string {
	return fmt.Sprintf(`
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "test-tag-%[1]d"
			description = "test tag description"
			entity_type = "DATABASE"
			required = false
		}

		data "nutanix_ndb_tags" "tags"{ }
	`, r)
}

func testAccEraTagsDataSourceConfigWithEntityType(r int) string {
	return fmt.Sprintf(`
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "test-tag-%[1]d"
			description = "test tag description"
			entity_type = "DATABASE"
			required = false
		}

		data "nutanix_ndb_tags" "tags"{ 
			entity_type= "DATABASE"
		}
	`, r)
}
