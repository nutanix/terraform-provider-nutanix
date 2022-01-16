package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixKarbonClusterDataSource_basic(t *testing.T) {
	t.Skip()
	r := acctest.RandInt()
	dataSourceName := "data.nutanix_karbon_cluster.kcluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterDataSourceConfig(subnetName, r, defaultContainter, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster.kcluster", "id"),
				),
			},
		},
	})
}

func TestAccNutanixKarbonClusterDataSource_basicByName(t *testing.T) {
	r := acctest.RandInt()
	//resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterDataSourceConfigByName(subnetName, r, defaultContainter, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_karbon_cluster.kcluster", "id"),
				),
			},
		},
	})
}

func testAccKarbonClusterDataSourceConfig(subnetName string, r int, containter string, workers int) string {
	return testAccNutanixKarbonClusterConfig(subnetName, r, containter, workers, "flannel") + `
	data "nutanix_karbon_cluster" "kcluster" {
		karbon_cluster_id = nutanix_karbon_cluster.cluster.id
	}
	`
}

func testAccKarbonClusterDataSourceConfigByName(subnetName string, r int, containter string, workers int) string {
	return testAccNutanixKarbonClusterConfig(subnetName, r, containter, workers, "flannel") + `
	data "nutanix_karbon_cluster" "kcluster" {
		karbon_cluster_name = nutanix_karbon_cluster.cluster.name
	}
	`
}
