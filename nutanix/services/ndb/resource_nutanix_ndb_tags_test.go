package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameTag = "nutanix_ndb_tag.acctest-managed"

func TestAccEraTag_basic(t *testing.T) {
	name := "test-tag-tf"
	desc := "this is tag desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTag, "name", name),
					resource.TestCheckResourceAttr(resourceNameTag, "description", desc),
					resource.TestCheckResourceAttr(resourceNameTag, "entity_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceNameTag, "required", "false"),
					resource.TestCheckResourceAttr(resourceNameTag, "status", "ENABLED"),
					resource.TestCheckResourceAttrSet(resourceNameTag, "date_created"),
					resource.TestCheckResourceAttrSet(resourceNameTag, "date_modified"),
				),
			},
		},
	})
}

func TestAccEraTag_WithUpdate(t *testing.T) {
	name := "test-tag-tf"
	updateName := "test-tag-updated"
	desc := "this is tag desc"
	updatedDesc := "this is updated tag desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTagConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTag, "name", name),
					resource.TestCheckResourceAttr(resourceNameTag, "description", desc),
					resource.TestCheckResourceAttr(resourceNameTag, "entity_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceNameTag, "required", "false"),
					resource.TestCheckResourceAttr(resourceNameTag, "status", "ENABLED"),
					resource.TestCheckResourceAttrSet(resourceNameTag, "date_created"),
					resource.TestCheckResourceAttrSet(resourceNameTag, "date_modified"),
				),
			},
			{
				Config: testAccEraTagUpdatedConfig(updateName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTag, "name", updateName),
					resource.TestCheckResourceAttr(resourceNameTag, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameTag, "entity_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceNameTag, "required", "true"),
					resource.TestCheckResourceAttr(resourceNameTag, "status", "ENABLED"),
					resource.TestCheckResourceAttrSet(resourceNameTag, "date_created"),
					resource.TestCheckResourceAttrSet(resourceNameTag, "date_modified"),
				),
			},
		},
	})
}

func testAccEraTagConfig(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "%[1]s"
			description = "%[2]s"
			entity_type = "DATABASE"
			required = false
		}
	`, name, desc)
}

func testAccEraTagUpdatedConfig(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_ndb_tag" "acctest-managed" {
			name= "%[1]s"
			description = "%[2]s"
			entity_type = "DATABASE"
			required = true
		}
	`, name, desc)
}
