package clustersv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameClusters = "data.nutanix_clusters_v2.test"

func TestAccNutanixClustersDataSourceV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClustersDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "cluster_entities.#"),
				),
			},
		},
	})
}

func TestAccNutanixClustersDataSourceV2_filter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClustersDataSourceConfigWithFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "cluster_entities.#"),
					resource.TestCheckResourceAttr(dataSourceNameClusters, "cluster_entities.0.config.0.cluster_function.0", "PRISM_CENTRAL"),
				),
			},
		},
	})
}

func TestAccNutanixClustersDataSourceV2_limit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClustersDataSourceConfigWithLimit(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "cluster_entities.#"),
					resource.TestCheckResourceAttr(dataSourceNameClusters, "cluster_entities.#", "1"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccClustersDataSourceConfig = `
data "nutanix_clusters_v2" "test" {}`

func testAccClustersDataSourceConfigWithFilter() string {
	return `
data "nutanix_clusters_v2" "test" {
	filter = "startswith(name, 'PC_')"
}`
}

func testAccClustersDataSourceConfigWithLimit() string {
	return ` 
data "nutanix_clusters_v2" "test" {
	limit = 1
}`
}
