package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameImagePlacementPolicy = "nutanix_image_placement_policy_v2.test"

func TestAccV2NutanixImagesPlacementPolicyResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-ipp-%d", r)
	desc := "test ipp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagesPlacementPolicyV4Config(name, desc, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameImagePlacementPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameImagePlacementPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameImagePlacementPolicy, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameImagePlacementPolicy, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameImagePlacementPolicy, "create_time"),
					resource.TestCheckResourceAttr(resourceNameImagePlacementPolicy, "enforcement_state", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceNameImagePlacementPolicy, "placement_type", "SOFT"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagesPlacementPolicyResource_SuspendAndResume(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-ipp-%d", r)
	desc := "test ipp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// create image placement policy
			{
				Config: testImagesPlacementPolicyV4Config(name, desc, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameImagePlacementPolicy, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameImagePlacementPolicy, "create_time"),
					resource.TestCheckResourceAttr(resourceNameImagePlacementPolicy, "enforcement_state", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceNameImagePlacementPolicy, "placement_type", "SOFT"),
				),
			},
			// suspend image placement policy
			{
				Config: testImagesPlacementPolicyV4Config(name, desc, `action = "SUSPEND"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameImagePlacementPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameImagePlacementPolicy, "action", "SUSPEND"),
				),
			},
			// resume image placement policy
			{
				Config: testImagesPlacementPolicyV4Config(name, desc, `action = "RESUME"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameImagePlacementPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameImagePlacementPolicy, "action", "RESUME"),
				),
			},
		},
	})
}

func testImagesPlacementPolicyV4Config(name, desc, action string) string {
	return fmt.Sprintf(`

		data "nutanix_categories_v2" "categories"{}

		locals {
			category0 = data.nutanix_categories_v2.categories.categories.0.ext_id
		}

		resource "nutanix_image_placement_policy_v2" "test" {
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

			%[3]s

			lifecycle{
				ignore_changes = [
					cluster_entity_filter,
					image_entity_filter,
				]
			}
			depends_on = [data.nutanix_categories_v2.categories]
		}


`, name, desc, action)
}
