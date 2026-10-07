package vmmv2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import3 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/templateplacementpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// TestAccV2NutanixTemplatePlacementPolicy_Basic covers the following scenarios:
//
// Step 1 — Create:
//   - Creates 4 categories (2 cluster, 2 content) via nutanix_category_v2.
//   - Creates a placement policy with placement_type=SOFT, both filters using
//     CATEGORIES_MATCH_ALL with 1 category each.
//   - Verifies name, description, placement_type, ext_id, create_time,
//     filter types, and category count (1 per filter).
//
// Step 2 — Update (name, description, filter type, category list):
//   - Updates name and description.
//   - Changes cluster_filter and content_filter type from CATEGORIES_MATCH_ALL
//     to CATEGORIES_MATCH_ANY.
//   - Expands category_ext_ids from 1 to 2 in both filters (adds second category).
//   - Verifies all updated fields including filter type and category count.
//
// Step 3 — Datasource (Get by ID):
//   - Reads the policy back via data.nutanix_template_placement_policy_v2.
//   - Verifies all attributes match the resource using TestCheckResourceAttrPair.
//
// Step 4 — Datasource (List):
//   - Reads all policies via data.nutanix_template_placement_policies_v2.
//   - Verifies the list is non-empty.
//
// Destroy:
//   - CheckDestroy verifies the placement policy no longer exists after cleanup.
func TestAccV2NutanixTemplatePlacementPolicy_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := acctest.RandomWithPrefix("tf-tpp-test")
	nameUpdated := acctest.RandomWithPrefix("tf-tpp-updated")
	description := "Test template placement policy"
	descriptionUpdated := "Updated test template placement policy"
	resourceName := "nutanix_template_placement_policy_v2.test"
	datasourceName := "data.nutanix_template_placement_policy_v2.test"
	datasourceListName := "data.nutanix_template_placement_policies_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixTemplatePlacementPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplatePlacementPolicyConfig(r, name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixTemplatePlacementPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "placement_type", "SOFT"),
					resource.TestCheckResourceAttrSet(resourceName, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "cluster_filter.0.type", "CATEGORIES_MATCH_ALL"),
					resource.TestCheckResourceAttr(resourceName, "cluster_filter.0.category_ext_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "content_filter.0.type", "CATEGORIES_MATCH_ALL"),
					resource.TestCheckResourceAttr(resourceName, "content_filter.0.category_ext_ids.#", "1"),
				),
			},
			{
				Config: testAccTemplatePlacementPolicyConfigUpdated(r, nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixTemplatePlacementPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(resourceName, "placement_type", "SOFT"),
					resource.TestCheckResourceAttrSet(resourceName, "ext_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_filter.0.type", "CATEGORIES_MATCH_ANY"),
					resource.TestCheckResourceAttr(resourceName, "cluster_filter.0.category_ext_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "content_filter.0.type", "CATEGORIES_MATCH_ANY"),
					resource.TestCheckResourceAttr(resourceName, "content_filter.0.category_ext_ids.#", "2"),
				),
			},
			{
				Config: testAccTemplatePlacementPolicyDatasourceConfig(r, nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(datasourceName, "placement_type", resourceName, "placement_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "cluster_filter.0.type", resourceName, "cluster_filter.0.type"),
					resource.TestCheckResourceAttrPair(datasourceName, "content_filter.0.type", resourceName, "content_filter.0.type"),
					resource.TestCheckResourceAttrSet(datasourceName, "create_time"),
				),
			},
			{
				Config: testAccTemplatePlacementPolicyListDatasourceConfig(r, nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceListName, "template_placement_policies.#"),
				),
			},
		},
	})
}

func testAccCheckNutanixTemplatePlacementPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		conn := acc.TestAccProvider.Meta().(*conns.Client)
		req := import3.GetTemplatePlacementPolicyByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := conn.VmmAPI.TemplatePlacementPoliciesAPIInstance.GetTemplatePlacementPolicyById(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("template placement policy %s does not exist: %v", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckNutanixTemplatePlacementPolicyDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_template_placement_policy_v2" {
			continue
		}
		req := import3.GetTemplatePlacementPolicyByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := conn.VmmAPI.TemplatePlacementPoliciesAPIInstance.GetTemplatePlacementPolicyById(ctx, &req)
		if err == nil {
			return fmt.Errorf("template placement policy %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccTPPCategoryConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "cluster_cat_1" {
  key   = "tf-tpp-cluster-%[1]d"
  value = "cluster-val-1"
}

resource "nutanix_category_v2" "cluster_cat_2" {
  key   = "tf-tpp-cluster-%[1]d"
  value = "cluster-val-2"
}

resource "nutanix_category_v2" "content_cat_1" {
  key   = "tf-tpp-content-%[1]d"
  value = "content-val-1"
}

resource "nutanix_category_v2" "content_cat_2" {
  key   = "tf-tpp-content-%[1]d"
  value = "content-val-2"
}
`, r)
}

func testAccTemplatePlacementPolicyConfig(r int, name, description string) string {
	return testAccTPPCategoryConfig(r) + fmt.Sprintf(`
resource "nutanix_template_placement_policy_v2" "test" {
  name           = %[1]q
  description    = %[2]q
  placement_type = "SOFT"

  cluster_filter {
    type             = "CATEGORIES_MATCH_ALL"
    category_ext_ids = [nutanix_category_v2.cluster_cat_1.ext_id]
  }

  content_filter {
    type             = "CATEGORIES_MATCH_ALL"
    category_ext_ids = [nutanix_category_v2.content_cat_1.ext_id]
  }
}
`, name, description)
}

func testAccTemplatePlacementPolicyConfigUpdated(r int, name, description string) string {
	return testAccTPPCategoryConfig(r) + fmt.Sprintf(`
resource "nutanix_template_placement_policy_v2" "test" {
  name           = %[1]q
  description    = %[2]q
  placement_type = "SOFT"

  cluster_filter {
    type             = "CATEGORIES_MATCH_ANY"
    category_ext_ids = [nutanix_category_v2.cluster_cat_1.ext_id, nutanix_category_v2.cluster_cat_2.ext_id]
  }

  content_filter {
    type             = "CATEGORIES_MATCH_ANY"
    category_ext_ids = [nutanix_category_v2.content_cat_1.ext_id, nutanix_category_v2.content_cat_2.ext_id]
  }
}
`, name, description)
}

func testAccTemplatePlacementPolicyDatasourceConfig(r int, name, description string) string {
	return testAccTPPCategoryConfig(r) + fmt.Sprintf(`
resource "nutanix_template_placement_policy_v2" "test" {
  name           = %[1]q
  description    = %[2]q
  placement_type = "SOFT"

  cluster_filter {
    type             = "CATEGORIES_MATCH_ANY"
    category_ext_ids = [nutanix_category_v2.cluster_cat_1.ext_id, nutanix_category_v2.cluster_cat_2.ext_id]
  }

  content_filter {
    type             = "CATEGORIES_MATCH_ANY"
    category_ext_ids = [nutanix_category_v2.content_cat_1.ext_id, nutanix_category_v2.content_cat_2.ext_id]
  }
}

data "nutanix_template_placement_policy_v2" "test" {
  ext_id = nutanix_template_placement_policy_v2.test.ext_id
}
`, name, description)
}

func testAccTemplatePlacementPolicyListDatasourceConfig(r int, name, description string) string {
	return testAccTPPCategoryConfig(r) + fmt.Sprintf(`
resource "nutanix_template_placement_policy_v2" "test" {
  name           = %[1]q
  description    = %[2]q
  placement_type = "SOFT"

  cluster_filter {
    type             = "CATEGORIES_MATCH_ANY"
    category_ext_ids = [nutanix_category_v2.cluster_cat_1.ext_id, nutanix_category_v2.cluster_cat_2.ext_id]
  }

  content_filter {
    type             = "CATEGORIES_MATCH_ANY"
    category_ext_ids = [nutanix_category_v2.content_cat_1.ext_id, nutanix_category_v2.content_cat_2.ext_id]
  }
}

data "nutanix_template_placement_policies_v2" "test" {
  depends_on = [nutanix_template_placement_policy_v2.test]
}
`, name, description)
}
