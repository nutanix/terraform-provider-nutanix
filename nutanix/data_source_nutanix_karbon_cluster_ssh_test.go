package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKarbonClusterSSHDataSource_basicx(t *testing.T) {
	r := acctest.RandInt()
	//resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	KubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterSSHDataSourceConfig(subnetName, r, defaultContainter, 1, KubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster_ssh.ssh", "id"),
					resource.TestCheckResourceAttr(
						"data.nutanix_karbon_cluster_ssh.ssh", "username", "admin"),
				),
			},
		},
	})
}

func TestAccKarbonClusterSSHDataSource_basicByName(t *testing.T) {
	r := acctest.RandInt()
	//resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	KubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterSSHDataSourceConfigByName(subnetName, r, defaultContainter, 1, KubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster_ssh.ssh", "id"),
					resource.TestCheckResourceAttr(
						"data.nutanix_karbon_cluster_ssh.ssh", "username", "admin"),
				),
			},
		},
	})
}

func testAccKarbonClusterSSHDataSourceConfig(subnetName string, r int, containter string, workers int, k8s string) string {
	return `
		data "nutanix_karbon_clusters" "kclusters" {}

		data "nutanix_karbon_cluster_ssh" "ssh" {
			karbon_cluster_id = data.nutanix_karbon_clusters.kclusters.clusters.0.uuid
		}
	`
}

func testAccKarbonClusterSSHDataSourceConfigByName(subnetName string, r int, containter string, workers int, k8s string) string {
	return `
		data "nutanix_karbon_clusters" "kclusters" {}

		data "nutanix_karbon_cluster_ssh" "ssh" {
			karbon_cluster_name = data.nutanix_karbon_clusters.kclusters.clusters.0.name
		}
	`
}
