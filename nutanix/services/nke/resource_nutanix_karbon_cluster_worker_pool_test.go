package nke_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccKarbonClusterWorkerPool_basic(t *testing.T) {
	resourceName := "nutanix_karbon_worker_nodepool.nodepool"
	subnetName := testVars.SubnetName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixKarbonClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixKarbonClusterWorkerNodePoolConfig(subnetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "workerpool1"),
					resource.TestCheckResourceAttr(resourceName, "num_instances", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "ahv_config.#"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.cpu", "4"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.disk_mib", "122880"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.memory_mib", "8192"),
					resource.TestCheckResourceAttrSet(resourceName, "node_os_version"),
					resource.TestCheckResourceAttr(resourceName, "labels.k1", "v1"),
					resource.TestCheckResourceAttr(resourceName, "labels.k2", "v2"),
				),
			},
			{ // Test for non-empty plans. No modification.
				Config:   testAccNutanixKarbonClusterWorkerNodePoolConfig(subnetName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccKarbonClusterWorkerPool_Update(t *testing.T) {
	resourceName := "nutanix_karbon_worker_nodepool.nodepool"
	subnetName := testVars.SubnetName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixKarbonClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixKarbonClusterWorkerNodePoolConfig(subnetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "workerpool1"),
					resource.TestCheckResourceAttr(resourceName, "num_instances", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "ahv_config.#"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.cpu", "4"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.disk_mib", "122880"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.memory_mib", "8192"),
					resource.TestCheckResourceAttrSet(resourceName, "node_os_version"),
					resource.TestCheckResourceAttr(resourceName, "labels.k1", "v1"),
					resource.TestCheckResourceAttr(resourceName, "labels.k2", "v2"),
				),
			},
			{ // Test for non-empty plans. No modification.
				Config:   testAccNutanixKarbonClusterWorkerNodePoolConfig(subnetName),
				PlanOnly: true,
			},
			{ // Test to update labels and increase nodes
				Config: testAccNutanixKarbonClusterWorkerNodePoolConfigUpdate(subnetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "workerpool1"),
					resource.TestCheckResourceAttr(resourceName, "num_instances", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "ahv_config.#"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.cpu", "4"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.disk_mib", "122880"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.memory_mib", "8192"),
					resource.TestCheckResourceAttrSet(resourceName, "node_os_version"),
					resource.TestCheckResourceAttr(resourceName, "labels.k1", "v1"),
					resource.TestCheckResourceAttr(resourceName, "labels.k2", "v2"),
					resource.TestCheckResourceAttr(resourceName, "labels.k3", "v3"),
				),
			},
			{ // Test to decrease the number of nodes
				Config: testAccNutanixKarbonClusterWorkerNodePoolConfig(subnetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "workerpool1"),
					resource.TestCheckResourceAttr(resourceName, "num_instances", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "ahv_config.#"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.cpu", "4"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.disk_mib", "122880"),
					resource.TestCheckResourceAttr(resourceName, "ahv_config.0.memory_mib", "8192"),
					resource.TestCheckResourceAttrSet(resourceName, "node_os_version"),
					resource.TestCheckResourceAttr(resourceName, "labels.k1", "v1"),
					resource.TestCheckResourceAttr(resourceName, "labels.k2", "v2"),
				),
			},
		},
	})
}

func testAccNutanixKarbonClusterWorkerNodePoolConfig(subnetName string) string {
	return fmt.Sprintf(`

		data "nutanix_karbon_clusters" "kclusters" {}

		data "nutanix_subnet" "karbon_subnet" {
			subnet_name = "%s"
		}

		resource "nutanix_karbon_worker_nodepool" "nodepool" {
			cluster_name = data.nutanix_karbon_clusters.kclusters.clusters.0.name
			name = "workerpool1"
			num_instances = 1
			ahv_config {
				cpu= 4
				disk_mib= 122880
				memory_mib=8192
				network_uuid= data.nutanix_subnet.karbon_subnet.id
			}
			labels={
				k1="v1"
				k2="v2"
			}
			depends_on = [ data.nutanix_karbon_clusters.kclusters ]
		}

	`, subnetName)
}

func testAccNutanixKarbonClusterWorkerNodePoolConfigUpdate(subnetName string) string {
	return fmt.Sprintf(`

		data "nutanix_karbon_clusters" "kclusters" {}

		data "nutanix_subnet" "karbon_subnet" {
			subnet_name = "%s"
		}

		resource "nutanix_karbon_worker_nodepool" "nodepool" {
			cluster_name = data.nutanix_karbon_clusters.kclusters.clusters.0.name
			name = "workerpool1"
			num_instances = 2
			ahv_config {
				cpu= 4
				disk_mib= 122880
				memory_mib=8192
				network_uuid= data.nutanix_subnet.karbon_subnet.id
			}
			labels={
				k1="v1"
				k2="v2"
				k3="v3"
			}
			depends_on = [ data.nutanix_karbon_clusters.kclusters ]
		}

	`, subnetName)
}
