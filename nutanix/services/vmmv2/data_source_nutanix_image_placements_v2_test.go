package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameImagePlacements = "data.nutanix_image_placement_policies_v2.test"

func TestAccV2NutanixImagePlacementsDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-ipp-%d", r)
	desc := "test ipp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagePlacementsPreConfigV2(name, desc) + testAccImagePlacementsDataSourceConfigV2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacements, "placement_policies.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagePlacementsDatasource_WithFilters(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-ipp-%d", r)
	desc := "test ipp description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagePlacementsPreConfigV2(name, desc) + testAccImagePlacementsDataSourceConfigV2WithFilters(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacements, "placement_policies.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacements, "placement_policies.0.placement_type"),
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacements, "placement_policies.0.description"),
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacements, "placement_policies.0.create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacements, "placement_policies.0.last_update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImagePlacements, "placement_policies.0.owner_ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagePlacementsDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagePlacementsDataSourceConfigV2WithInvalidFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameImagePlacements, "placement_policies.#", "0"),
				),
			},
		},
	})
}

func testAccImagePlacementsPreConfigV2(name, desc string) string {
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
`, name, desc)
}

func testAccImagePlacementsDataSourceConfigV2() string {
	return `
		data "nutanix_image_placement_policies_v2" "test"{

			depends_on = [
				resource.nutanix_image_placement_policy_v2.ipp
			]
		}
`
}

func testAccImagePlacementsDataSourceConfigV2WithFilters(name string) string {
	return fmt.Sprintf(`

		data "nutanix_image_placement_policies_v2" "test"{
			page=0
			limit=10
			filter="name eq '%s'"

		    depends_on = [
				resource.nutanix_image_placement_policy_v2.ipp
			]
		}
`, name)
}

func testAccImagePlacementsDataSourceConfigV2WithInvalidFilters() string {
	return `
		data "nutanix_image_placement_policies_v2" "test"{
			filter="name eq 'invalid_filter'"
		}
`
}
