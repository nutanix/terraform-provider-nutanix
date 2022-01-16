package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixKarbonClusterKubeConfigDataSource_basic(t *testing.T) {
	t.Skip()
	r := acctest.RandInt()
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterKubeConfigDataSourceConfig(subnetName, r, defaultContainter, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster_kubeconfig.config", "id"),
					resource.TestCheckResourceAttr(
						"data.nutanix_karbon_cluster_kubeconfig.config", "karbon_cluster_name", fmt.Sprintf("test-karbon-%d", r)),
				),
			},
		},
	})
}

func TestAccNutanixKarbonClusterKubeConfigDataSource_basicByName(t *testing.T) {
	r := acctest.RandInt()
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterKubeConfigDataSourceConfigByName(subnetName, r, defaultContainter, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster_kubeconfig.config", "id"),
					resource.TestCheckResourceAttr(
						"data.nutanix_karbon_cluster_kubeconfig.config", "karbon_cluster_name", fmt.Sprintf("test-karbon-%d", r)),
				),
			},
		},
	})
}

func testAccKarbonClusterKubeConfigDataSourceConfig(subnetName string, r int, containter string, workers int) string {
	return testAccNutanixKarbonClusterConfig(subnetName, r, containter, workers, "flannel") + `
	data "nutanix_karbon_cluster_kubeconfig" "config" {
		karbon_cluster_id = nutanix_karbon_cluster.cluster.id
	}
	`
}

func testAccKarbonClusterKubeConfigDataSourceConfigByName(subnetName string, r int, containter string, workers int) string {
	return testAccNutanixKarbonClusterConfig(subnetName, r, containter, workers, "flannel") + `
	data "nutanix_karbon_cluster_kubeconfig" "config" {
		karbon_cluster_name = nutanix_karbon_cluster.cluster.name
	}
	`
}
