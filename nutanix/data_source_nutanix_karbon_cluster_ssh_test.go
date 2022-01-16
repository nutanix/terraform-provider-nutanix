package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNutanixKarbonClusterSSHDataSource_basicx(t *testing.T) {
	t.Skip()
	r := acctest.RandInt()
	//resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterSSHDataSourceConfig(subnetName, r, defaultContainter, 1),
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

func TestAccNutanixKarbonClusterSSHDataSource_basicByName(t *testing.T) {
	r := acctest.RandInt()
	//resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKarbonClusterSSHDataSourceConfigByName(subnetName, r, defaultContainter, 1),
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

func testAccKarbonClusterSSHDataSourceConfig(subnetName string, r int, containter string, workers int) string {
	return testAccNutanixKarbonClusterConfig(subnetName, r, containter, workers, "flannel") + `
	data "nutanix_karbon_cluster_ssh" "ssh" {
		karbon_cluster_id = nutanix_karbon_cluster.cluster.id
	}
	`
}

func testAccKarbonClusterSSHDataSourceConfigByName(subnetName string, r int, containter string, workers int) string {
	return testAccNutanixKarbonClusterConfig(subnetName, r, containter, workers, "flannel") + `
	data "nutanix_karbon_cluster_ssh" "ssh" {
		karbon_cluster_name = nutanix_karbon_cluster.cluster.name
	}
	`
}
