package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameImagePlacementPolicy = "data.nutanix_image_placement_policy_v2.test"

func TestAccV2NutanixImagePlacementDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-ipp-%d", r)
	desc := "test ipp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagePlacementDataSourceConfigV2(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameImagePlacementPolicy, "name", name),
					resource.TestCheckResourceAttr(datasourceNameImagePlacementPolicy, "placement_type", "SOFT"),
					resource.TestCheckResourceAttr(datasourceNameImagePlacementPolicy, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacementPolicy, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacementPolicy, "last_update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacementPolicy, "owner_ext_id"),
				),
			},
		},
	})
}

func testAccImagePlacementDataSourceConfigV2(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_categories_v2" "categories" {}

		locals {
			category0 = data.nutanix_categories_v2.categories.categories.0.ext_id
		}
		resource "nutanix_image_placement_policy_v2" "ipp" {
			name           = "%[1]s"
			description    = "%[2]s"
			placement_type = "SOFT"
			cluster_entity_filter {
				category_ext_ids = [
					local.category0,
				]
				type = "CATEGORIES_MATCH_ALL"
			}
			image_entity_filter {
				category_ext_ids = [
					local.category0,
				]
				type = "CATEGORIES_MATCH_ALL"
			}

			lifecycle{
				ignore_changes = [
					cluster_entity_filter,
					image_entity_filter,
				]
			}
		}

		data "nutanix_image_placement_policy_v2" "test"{
			ext_id = resource.nutanix_image_placement_policy_v2.ipp.id
		}
`, name, desc)
}
