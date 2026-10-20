package clustersv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestResourceNutanixClusterCategoriesV2(t *testing.T) {
	const resourceNameClusterCategories = "nutanix_cluster_category_associations_v2.cluster_categories"

	r1, r2 := acc.RandIntBetween(1, 500), acc.RandIntBetween(501, 1000)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixClusterCategoriesV2Config("nutanix_category_v2.category-1.id", r1, r2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_cluster_category_associations_v2.cluster_categories", "categories.#", "1"),
					resource.TestCheckResourceAttrPair("nutanix_cluster_category_associations_v2.cluster_categories", "categories.0", "nutanix_category_v2.category-1", "id"),
					resource.TestCheckResourceAttrPair("nutanix_cluster_category_associations_v2.cluster_categories", "cluster_ext_id", "data.nutanix_clusters_v2.clusters", "cluster_entities.0.ext_id"),
				),
			},
			{
				ResourceName:      resourceNameClusterCategories,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNutanixClusterCategoriesV2Config("nutanix_category_v2.category-1.id, nutanix_category_v2.category-2.id", r1, r2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_cluster_category_associations_v2.cluster_categories", "categories.#", "2"),
					resource.TestCheckResourceAttrPair("nutanix_cluster_category_associations_v2.cluster_categories", "cluster_ext_id", "data.nutanix_clusters_v2.clusters", "cluster_entities.0.ext_id"),
				),
			},
			{
				Config: testAccNutanixClusterCategoriesV2Config("nutanix_category_v2.category-2.id", r1, r2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_cluster_category_associations_v2.cluster_categories", "categories.#", "1"),
					resource.TestCheckResourceAttrPair("nutanix_cluster_category_associations_v2.cluster_categories", "categories.0", "nutanix_category_v2.category-2", "id"),
					resource.TestCheckResourceAttrPair("nutanix_cluster_category_associations_v2.cluster_categories", "cluster_ext_id", "data.nutanix_clusters_v2.clusters", "cluster_entities.0.ext_id"),
				),
			},
		},
	})
}

func testAccNutanixClusterCategoriesV2Config(categories string, r1 int, r2 int) string {
	return fmt.Sprintf(`
	resource "nutanix_category_v2" "category-1" {
		description = "Test category"
		key = "test-category-%[2]d"
		value = "test-category-%[2]d"
	}

	resource "nutanix_category_v2" "category-2" {
		description = "Test category"
		key = "test-category-%[3]d"
		value = "test-category-%[3]d"
	}

	data "nutanix_clusters_v2" "clusters" {
		filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
	}

	resource "nutanix_cluster_category_associations_v2" "cluster_categories" {
		cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
		categories = [%[1]s]
	}
	`, categories, r1, r2)
}
