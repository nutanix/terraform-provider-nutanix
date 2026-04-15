package microsegv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixEntityGroupV2Resource_Basic(t *testing.T) {
	r := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-entity-group-%d", r)
	description := fmt.Sprintf("tf-entity-group-%d_desc", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testEntityGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityGroupV2ResourceConfig(r, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameEntityGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "description", description),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.selected_by", "CATEGORY_EXT_ID"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.type", "VM"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.reference_ext_ids.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNameEntityGroupV2, "allowed_config.0.entities.0.reference_ext_ids.0", "nutanix_category_v2.categories.0", "id"),
					resource.TestCheckResourceAttrPair(resourceNameEntityGroupV2, "allowed_config.0.entities.0.reference_ext_ids.1", "nutanix_category_v2.categories.1", "id"),
				),
			},
			{
				Config: testAccEntityGroupV2ResourceUpdateConfig(r, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameEntityGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "description", description+" updated"),
				),
			},
		},
	})
}

func TestAccNutanixEntityGroupV2Resource_WithoutName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEntityGroupV2ResourceConfigWithoutName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixEntityGroupV2Resource_WithWrongReferenceExtIds(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-entity-group-wrong-ref-%d", r)
	description := "entity_group_wrong_ref_ext_ids_desc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEntityGroupV2ResourceConfigWithWrongReferenceExtIds(name, description),
				ExpectError: regexp.MustCompile("categories were not found"),
			},
		},
	})
}

func testAccEntityGroupV2ResourceConfig(r int, name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_%[1]d_${count.index}_key"
  value       = "tf_entity_group_%[1]d_${count.index}_value"
  description = "tf_entity_group_%[1]d_${count.index}_description"
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%[2]s"
  description = "%[3]s"

  allowed_config {
    entities {
      type             = "VM"
      selected_by      = "CATEGORY_EXT_ID"
      reference_ext_ids = [
        nutanix_category_v2.categories[0].id,
        nutanix_category_v2.categories[1].id
      ]
    }
  }
}
`, r, name, description)
}

func testAccEntityGroupV2ResourceUpdateConfig(r int, name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_%[1]d_${count.index}_key"
  value       = "tf_entity_group_%[1]d_${count.index}_value"
  description = "tf_entity_group_%[1]d_${count.index}_description"
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%s-updated"
  description = "%s updated"

  allowed_config {
    entities {
      type             = "VM"
      selected_by      = "CATEGORY_EXT_ID"
      reference_ext_ids = [
        nutanix_category_v2.categories[0].id,
        nutanix_category_v2.categories[1].id
      ]
    }
  }
}
`, r, name, description)
}

func testAccEntityGroupV2ResourceConfigWithoutName() string {
	return `
resource "nutanix_entity_group_v2" "test" {
  description = "entity_group_without_name_desc"
}
`
}

func testAccEntityGroupV2ResourceConfigWithWrongReferenceExtIds(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_entity_group_v2" "test" {
  name        = "%s"
  description = "%s"

  allowed_config {
    entities {
      type            = "VM"
      selected_by     = "CATEGORY_EXT_ID"
      reference_ext_ids = [
        "83cbf00d-782f-4efc-87c8-4129f5942aaa",
        "83cbf00d-782f-4efc-87c8-4129f5942bbb"
      ]
    }
  }
}
`, name, description)
}
