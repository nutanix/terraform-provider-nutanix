package nke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccKarbonClusterKubeConfigDataSource_basic(t *testing.T) {
	r := acctest.RandInt()
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	KubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterKubeConfigDataSourceConfig(subnetName, r, defaultContainter, 1, KubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster_kubeconfig.config", "id"),
					resource.TestCheckResourceAttr(
						"data.nutanix_karbon_cluster_kubeconfig.config", "name", "test-karbon"),
				),
			},
		},
	})
}

func TestAccKarbonClusterKubeConfigDataSource_basicByName(t *testing.T) {
	r := acctest.RandInt()
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	KubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterKubeConfigDataSourceConfigByName(subnetName, r, defaultContainter, 1, KubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster_kubeconfig.config", "id"),
					resource.TestCheckResourceAttr(
						"data.nutanix_karbon_cluster_kubeconfig.config", "karbon_cluster_name", "test-karbon"),
				),
			},
		},
	})
}

func testAccKarbonClusterKubeConfigDataSourceConfig(subnetName string, r int, containter string, workers int, k8s string) string {
	return `
		data "nutanix_karbon_clusters" "kclusters" {}

		data "nutanix_karbon_cluster_kubeconfig" "config" {
			karbon_cluster_id = data.nutanix_karbon_clusters.kclusters.clusters.0.uuid
		}
	`
}

func testAccKarbonClusterKubeConfigDataSourceConfigByName(subnetName string, r int, containter string, workers int, k8s string) string {
	return `
		data "nutanix_karbon_clusters" "kclusters" {}

		data "nutanix_karbon_cluster_kubeconfig" "config" {
			karbon_cluster_name = data.nutanix_karbon_clusters.kclusters.clusters.0.name
		}
	`
}
