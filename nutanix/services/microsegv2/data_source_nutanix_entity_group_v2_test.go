package microsegv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameEntityGroupV2 = "data.nutanix_entity_group_v2.test"

func TestAccNutanixEntityGroupV2Datasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-entity-group-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testEntityGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityGroupV2DatasourceConfig(r, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameEntityGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameEntityGroupV2, "name", name),
				),
			},
		},
	})
}

// TestAccNutanixEntityGroupV2Datasource_WrongExtID tests that the data source fails
// as expected when given a non-existent ext_id.
func TestAccNutanixEntityGroupV2Datasource_WrongExtID(t *testing.T) {
	wrongExtID := "83cbf00d-782f-4efc-87c8-4129f5942aaa"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEntityGroupV2DatasourceConfigWrongExtID(wrongExtID),
				ExpectError: regexp.MustCompile(`error while fetching Entity Group|entity not found|Failed to get Entity Group`),
			},
		},
	})
}

func testAccEntityGroupV2DatasourceConfig(r int, name string) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_ds_%[1]d_${count.index}_key"
  value       = "tf_entity_group_ds_%[1]d_${count.index}_value"
  description = "tf_entity_group_ds_%[1]d_${count.index}_description"
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%s"
  description = "terraform test entity group for datasource"

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

data "nutanix_entity_group_v2" "test" {
  ext_id = nutanix_entity_group_v2.test.id
}
`, r, name)
}

func testAccEntityGroupV2DatasourceConfigWrongExtID(extID string) string {
	return fmt.Sprintf(`
data "nutanix_entity_group_v2" "test" {
  ext_id = "%s"
}
`, extID)
}
