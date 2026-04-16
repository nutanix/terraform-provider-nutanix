package microsegv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
)

const dataSourceNameEntityGroupsV2 = "data.nutanix_entity_groups_v2.test"

func TestAccNutanixEntityGroupsV2Datasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-entity-groups-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testEntityGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityGroupsV2DatasourceConfig(r, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameEntityGroupsV2, "id"),
					common.CheckAttributeLength(dataSourceNameEntityGroupsV2, "entity_groups", 1),
					resource.TestCheckResourceAttrSet(dataSourceNameEntityGroupsV2, "entity_groups.0.ext_id"),
				),
			},
		},
	})
}

func TestAccNutanixEntityGroupsV2Datasource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-entity-groups-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testEntityGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityGroupsV2DatasourceConfigWithFilter(r, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameEntityGroupsV2, "id"),
					common.CheckAttributeLengthEqual(dataSourceNameEntityGroupsV2, "entity_groups", 1),
					resource.TestCheckResourceAttrSet(dataSourceNameEntityGroupsV2, "entity_groups.0.ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameEntityGroupsV2, "entity_groups.0.name", name),
				),
			},
		},
	})
}

func TestAccNutanixEntityGroupsV2Datasource_WithLimit(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-entity-groups-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testEntityGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityGroupsV2DatasourceConfigWithLimit(r, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameEntityGroupsV2, "id"),
					resource.TestCheckResourceAttr(dataSourceNameEntityGroupsV2, "entity_groups.#", "1"),
				),
			},
		},
	})
}

func testAccEntityGroupsV2DatasourceConfig(r int, name string) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_ds_%[1]d_${count.index}_key"
  value       = "tf_entity_group_ds_%[1]d_${count.index}_value"
  description = "tf_entity_group_ds_%[1]d_${count.index}_description"
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%s"
  description = "terraform test entity group for list datasource"

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

data "nutanix_entity_groups_v2" "test" {
  depends_on = [nutanix_entity_group_v2.test]
}
`, r, name)
}

func testAccEntityGroupsV2DatasourceConfigWithFilter(r int, name string) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_ds_%[1]d_${count.index}_key"
  value       = "tf_entity_group_ds_%[1]d_${count.index}_value"
  description = "tf_entity_group_ds_%[1]d_${count.index}_description"
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%s"
  description = "terraform test entity group for list datasource"

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

data "nutanix_entity_groups_v2" "test" {
  filter = "name eq '${nutanix_entity_group_v2.test.name}'"
  depends_on = [nutanix_entity_group_v2.test]
}
`, r, name)
}

func testAccEntityGroupsV2DatasourceConfigWithLimit(r int, name string) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_ds_%[1]d_${count.index}_key"
  value       = "tf_entity_group_ds_%[1]d_${count.index}_value"
  description = "tf_entity_group_ds_%[1]d_${count.index}_description"
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%s"
  description = "terraform test entity group for list datasource with limit"

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

data "nutanix_entity_groups_v2" "test" {
  limit      = 1
  depends_on = [nutanix_entity_group_v2.test]
}
`, r, name)
}
