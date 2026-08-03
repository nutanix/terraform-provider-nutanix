package vmmv2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/imageratelimitpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameImageRateLimitPolicy = "nutanix_image_rate_limit_policy_v2.test"

// TestAccV2NutanixImageRateLimitPolicyResource_Basic covers the following scenarios:
//
// Step 1 — Create:
//   - Creates 2 categories (cat_1, cat_2) via nutanix_category_v2.
//   - Creates a rate limit policy with name, description, rate_limit_kbps=1000,
//     and a cluster_entity_filter using CATEGORIES_MATCH_ALL with 1 category (cat_1).
//   - Verifies name, description, rate_limit_kbps, ext_id, create_time,
//     last_update_time, filter type="CATEGORIES_MATCH_ALL", and category count=1.
//
// Step 2 — Update (name + description):
//   - Updates name and description.
//   - Verifies updated values and that rate_limit_kbps remains stable.
//
// Step 3 — Update (rate_limit_kbps):
//   - Updates rate_limit_kbps from 1000 to 2000.
//   - Verifies the new rate limit value.
//
// Step 4 — Update (filter type CATEGORIES_MATCH_ALL -> CATEGORIES_MATCH_ANY + expand categories from 1 to 2):
//   - Changes cluster_entity_filter type from CATEGORIES_MATCH_ALL to CATEGORIES_MATCH_ANY.
//   - Expands category_ext_ids from 1 category (cat_1) to 2 categories (cat_1 + cat_2).
//   - Verifies filter type="CATEGORIES_MATCH_ANY" and category_ext_ids.#=2.
//   - Verifies name, description, rate_limit_kbps remain stable.
//
// Step 5 — Update (shrink categories from 2 back to 1):
//   - Removes cat_2, keeping only cat_1 in category_ext_ids.
//   - Verifies category_ext_ids.#=1 and filter type remains CATEGORIES_MATCH_ANY.
//
// Step 6 — Datasource (singular, Get by ID):
//   - Reads the policy back via data.nutanix_image_rate_limit_policy_v2.
//   - Verifies all attributes match the resource using TestCheckResourceAttrPair.
//
// Step 7 — Datasource (plural, List):
//   - Reads all policies via data.nutanix_image_rate_limit_policies_v2.
//   - Verifies the list is non-empty.
//
// Destroy:
//   - CheckDestroy calls GetRateLimitPolicyById and verifies it returns an
//     error (policy no longer exists).
func TestAccV2NutanixImageRateLimitPolicyResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-rate-limit-%d", r)
	updatedName := fmt.Sprintf("test-rate-limit-updated-%d", r)
	desc := "test rate limit policy description"
	updatedDesc := "updated rate limit policy description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixImageRateLimitPolicyDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with CATEGORIES_MATCH_ALL, 1 category
			{
				Config: testImageRateLimitPolicyConfig(r, name, desc, 1000, "CATEGORIES_MATCH_ALL", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameImageRateLimitPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "description", desc),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "rate_limit_kbps", "1000"),
					resource.TestCheckResourceAttrSet(resourceNameImageRateLimitPolicy, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameImageRateLimitPolicy, "last_update_time"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "cluster_entity_filter.0.type", "CATEGORIES_MATCH_ALL"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "cluster_entity_filter.0.category_ext_ids.#", "1"),
				),
			},
			// Step 2: Update name + description
			{
				Config: testImageRateLimitPolicyConfig(r, updatedName, updatedDesc, 1000, "CATEGORIES_MATCH_ALL", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "rate_limit_kbps", "1000"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "cluster_entity_filter.0.type", "CATEGORIES_MATCH_ALL"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "cluster_entity_filter.0.category_ext_ids.#", "1"),
				),
			},
			// Step 3: Update rate_limit_kbps 1000 -> 2000
			{
				Config: testImageRateLimitPolicyConfig(r, updatedName, updatedDesc, 2000, "CATEGORIES_MATCH_ALL", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "rate_limit_kbps", "2000"),
				),
			},
			// Step 4: Update filter type ALL -> ANY + expand categories from 1 to 2
			{
				Config: testImageRateLimitPolicyConfig(r, updatedName, updatedDesc, 2000, "CATEGORIES_MATCH_ANY", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "rate_limit_kbps", "2000"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "cluster_entity_filter.0.type", "CATEGORIES_MATCH_ANY"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "cluster_entity_filter.0.category_ext_ids.#", "2"),
				),
			},
			// Step 5: Shrink categories from 2 back to 1, keep ANY
			{
				Config: testImageRateLimitPolicyConfig(r, updatedName, updatedDesc, 2000, "CATEGORIES_MATCH_ANY", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "cluster_entity_filter.0.type", "CATEGORIES_MATCH_ANY"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "cluster_entity_filter.0.category_ext_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameImageRateLimitPolicy, "rate_limit_kbps", "2000"),
				),
			},
			// Step 6: Datasource singular
			{
				Config: testImageRateLimitPolicyDatasourceConfig(r, updatedName, updatedDesc, 2000, "CATEGORIES_MATCH_ANY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_image_rate_limit_policy_v2.test", "name", updatedName),
					resource.TestCheckResourceAttr("data.nutanix_image_rate_limit_policy_v2.test", "description", updatedDesc),
					resource.TestCheckResourceAttr("data.nutanix_image_rate_limit_policy_v2.test", "rate_limit_kbps", "2000"),
					resource.TestCheckResourceAttr("data.nutanix_image_rate_limit_policy_v2.test", "cluster_entity_filter.0.type", "CATEGORIES_MATCH_ANY"),
					resource.TestCheckResourceAttr("data.nutanix_image_rate_limit_policy_v2.test", "cluster_entity_filter.0.category_ext_ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.nutanix_image_rate_limit_policy_v2.test", "ext_id"),
					resource.TestCheckResourceAttrSet("data.nutanix_image_rate_limit_policy_v2.test", "owner_ext_id"),
				),
			},
			// Step 7: Datasource plural (list)
			{
				Config: testImageRateLimitPoliciesListDatasourceConfig(r, updatedName, updatedDesc, 2000, "CATEGORIES_MATCH_ANY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_image_rate_limit_policies_v2.test", "rate_limit_policies.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_image_rate_limit_policies_v2.test", "rate_limit_policies.0.name", updatedName),
				),
			},
		},
	})
}

func testAccCheckNutanixImageRateLimitPolicyDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_image_rate_limit_policy_v2" {
			continue
		}
		req := import3.GetRateLimitPolicyByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := conn.VmmAPI.ImageRateLimitPoliciesAPIInstance.GetRateLimitPolicyById(ctx, &req)
		if err == nil {
			return fmt.Errorf("image rate limit policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccIRLPCategoryConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "cat_1" {
  key   = "tf-irlp-cat1-%[1]d"
  value = "cat-val-1"
}

resource "nutanix_category_v2" "cat_2" {
  key   = "tf-irlp-cat2-%[1]d"
  value = "cat-val-2"
}
`, r)
}

func testImageRateLimitPolicyConfig(r int, name, desc string, rateLimitKbps int, filterType string, twoCategories bool) string {
	categoryExpr := `[nutanix_category_v2.cat_1.ext_id]`
	if twoCategories {
		categoryExpr = `[nutanix_category_v2.cat_1.ext_id, nutanix_category_v2.cat_2.ext_id]`
	}
	return testAccIRLPCategoryConfig(r) + fmt.Sprintf(`
resource "nutanix_image_rate_limit_policy_v2" "test" {
  name            = "%[1]s"
  description     = "%[2]s"
  rate_limit_kbps = %[3]d
  cluster_entity_filter {
    category_ext_ids = %[4]s
    type             = "%[5]s"
  }
}
`, name, desc, rateLimitKbps, categoryExpr, filterType)
}

func testImageRateLimitPolicyDatasourceConfig(r int, name, desc string, rateLimitKbps int, filterType string) string {
	return testImageRateLimitPolicyConfig(r, name, desc, rateLimitKbps, filterType, false) + `
data "nutanix_image_rate_limit_policy_v2" "test" {
  ext_id = nutanix_image_rate_limit_policy_v2.test.ext_id
}
`
}

func testImageRateLimitPoliciesListDatasourceConfig(r int, name, desc string, rateLimitKbps int, filterType string) string {
	return testImageRateLimitPolicyConfig(r, name, desc, rateLimitKbps, filterType, false) + `
data "nutanix_image_rate_limit_policies_v2" "test" {
  filter     = "name eq '${nutanix_image_rate_limit_policy_v2.test.name}'"
  depends_on = [nutanix_image_rate_limit_policy_v2.test]
}
`
}
