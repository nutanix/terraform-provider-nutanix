package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceNameTag = "nutanix_ndb_tags.acctest-managed"

func TestAccEra_Tagbasic(t *testing.T) {
	name := "test-tag-tf"
	desc := "this is tag desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
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

func TestAccEra_TagWithUpdate(t *testing.T) {
	name := "test-tag-tf"
	updateName := "test-tag-updated"
	desc := "this is tag desc"
	updatedDesc := "this is updated tag desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
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
		resource "nutanix_ndb_tags" "acctest-managed" {
			name= "%[1]s"
			description = "%[2]s"
			entity_type = "DATABASE"
			required = false
		}
	`, name, desc)
}

func testAccEraTagUpdatedConfig(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_ndb_tags" "acctest-managed" {
			name= "%[1]s"
			description = "%[2]s"
			entity_type = "DATABASE"
			required = true
		}
	`, name, desc)
}
