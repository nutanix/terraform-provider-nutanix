package clusters_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixClusterDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_cluster.cluster", "id"),
				),
			},
		},
	})
}

func TestAccNutanixClusterDataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfigByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_cluster.cluster", "id"),
					resource.TestCheckResourceAttrSet(
						"data.nutanix_cluster.cluster", "name"),
					resource.TestCheckResourceAttrSet(
						"data.nutanix_cluster.cluster", "cluster_id"),
				),
			},
		},
	})
}

func TestAccNutanixClusterDataSource_ByNameNotExisting(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterDataSourceConfigByNameNotExisting,
				ExpectError: regexp.MustCompile("did not find cluster with name *"),
			},
		},
	})
}

func TestAccNutanixClusterDataSource_NameUuidConflict(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterDataSourceConfigNameUUIDConflict,
				ExpectError: regexp.MustCompile(" * conflicts with *"),
			},
		},
	})
}

const testAccClusterDataSourceConfig = `
data "nutanix_clusters" "clusters" {}


data "nutanix_cluster" "cluster" {
	cluster_id = data.nutanix_clusters.clusters.entities.0.metadata.uuid
}`

const testAccClusterDataSourceConfigByName = `
data "nutanix_clusters" "clusters" {}


data "nutanix_cluster" "cluster" {
	name = data.nutanix_clusters.clusters.entities.0.name
}`

const testAccClusterDataSourceConfigByNameNotExisting = `
data "nutanix_clusters" "clusters" {}


data "nutanix_cluster" "cluster" {
	name = "ThisDoesNotExist"
}`

const testAccClusterDataSourceConfigNameUUIDConflict = `
data "nutanix_clusters" "clusters" {}


data "nutanix_cluster" "cluster" {
	cluster_id = data.nutanix_clusters.clusters.entities.0.metadata.uuid
	name = data.nutanix_clusters.clusters.entities.0.name
}`
