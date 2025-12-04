package nke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccKarbonClusterDataSource_basic(t *testing.T) {
	r := acctest.RandInt()
	dataSourceName := "data.nutanix_karbon_cluster.kcluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	KubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterDataSourceConfig(subnetName, r, defaultContainter, 1, KubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster.kcluster", "id"),
				),
			},
		},
	})
}

func TestAccKarbonClusterDataSource_basicByName(t *testing.T) {
	r := acctest.RandInt()
	//resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	KubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterDataSourceConfigByName(subnetName, r, defaultContainter, 1, KubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster.kcluster", "id"),
				),
			},
		},
	})
}

func testAccKarbonClusterDataSourceConfig(subnetName string, r int, containter string, workers int, k8s string) string {
	return `
		data "nutanix_karbon_clusters" "kclusters" {}

		data "nutanix_karbon_cluster" "kcluster" {
			karbon_cluster_id = data.nutanix_karbon_clusters.kclusters.clusters.0.uuid
		}
	`
}

func testAccKarbonClusterDataSourceConfigByName(subnetName string, r int, containter string, workers int, k8s string) string {
	return `
		data "nutanix_karbon_clusters" "kclusters" {}

		data "nutanix_karbon_cluster" "kcluster" {
			karbon_cluster_name = data.nutanix_karbon_clusters.kclusters.clusters.0.name
		}

	`
}
