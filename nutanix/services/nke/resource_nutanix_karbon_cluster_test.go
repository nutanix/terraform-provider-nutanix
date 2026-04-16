package nke_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccKarbonCluster_basic(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	kubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixKarbonClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixKarbonClusterConfig(subnetName, r, defaultContainter, 1, "flannel", kubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-karbon-%d", r)),
					resource.TestCheckResourceAttr(resourceName, "etcd_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "master_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_class_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.0.num_instances", "1"),
				),
			},
			{
				Config: testAccNutanixKarbonClusterConfig(subnetName, r, defaultContainter, 2, "flannel", kubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-karbon-%d", r)),
					resource.TestCheckResourceAttr(resourceName, "etcd_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "master_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_class_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.0.num_instances", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version", "master_node_pool", "worker_node_pool", "storage_class_config", "wait_timeout_minutes"}, //Wil be fixed on future API versions
			},
		},
	})
}

func TestAccKarbonCluster_scaleDown(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	kubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixKarbonClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixKarbonClusterConfig(subnetName, r, defaultContainter, 3, "flannel", kubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-karbon-%d", r)),
					resource.TestCheckResourceAttr(resourceName, "etcd_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "master_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_class_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.0.num_instances", "3"),
				),
			},
			{
				Config: testAccNutanixKarbonClusterConfig(subnetName, r, defaultContainter, 1, "flannel", kubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-karbon-%d", r)),
					resource.TestCheckResourceAttr(resourceName, "etcd_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "master_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_class_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.0.num_instances", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version", "master_node_pool", "worker_node_pool", "storage_class_config", "wait_timeout_minutes"}, //Wil be fixed on future API versions
			},
		},
	})
}

func TestAccKarbonCluster_updateCNI(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_karbon_cluster.cluster"
	subnetName := testVars.SubnetName
	defaultContainter := testVars.DefaultContainerName
	kubernetesVersion := testVars.KubernetesVersion
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixKarbonClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixKarbonClusterConfig(subnetName, r, defaultContainter, 1, "flannel", kubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-karbon-%d", r)),
					resource.TestCheckResourceAttr(resourceName, "etcd_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "master_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_class_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.0.num_instances", "1"),
				),
			},
			{
				Config: testAccNutanixKarbonClusterConfig(subnetName, r, defaultContainter, 2, "calico", kubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-karbon-%d", r)),
					resource.TestCheckResourceAttr(resourceName, "etcd_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "master_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_class_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "worker_node_pool.0.num_instances", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version", "master_node_pool", "worker_node_pool", "storage_class_config", "wait_timeout_minutes"}, //Wil be fixed on future API versions
			},
		},
	})
}

func testAccCheckNutanixKarbonClusterDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_karbon_cluster" {
			continue
		}
		for {
			_, err := conn.KarbonAPI.Cluster.GetKarbonCluster(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "Not Found:K8s cluster not found.") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}
	}

	return nil
}

func testAccCheckNutanixKarbonClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccNutanixKarbonClusterConfig(subnetName string, r int, containter string, workers int, cni, k8sVersion string) string {
	return fmt.Sprintf(`
	locals {
		cluster_id = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		node_os_version   = "%[1]s"
		deployment_type   = ""
		amount_of_workers = %[2]d
		amount_of_masters = 1
		cni               = "%[3]s"
		master_vip        = ""
	}

	data "nutanix_clusters" "clusters" {}

	data "nutanix_subnet" "karbon_subnet" {
		subnet_name = "%[4]s"
	}

	resource "nutanix_karbon_cluster" "cluster" {
		name    = "test-karbon-%[5]d"
		version = "%[7]s"

		dynamic "active_passive_config" {
		  for_each = local.deployment_type == "active-passive" ? [1] : []
		  content {
			external_ipv4_address = local.master_vip
		  }
		}
		dynamic "external_lb_config" {
		  for_each = local.deployment_type == "active-active" ? [1] : []
		  content {
			external_ipv4_address = "10.10.30.228"
			master_nodes_config {
			  ipv4_address   = "10.10.100.171"
			  node_pool_name = "master_node_pool"
			}
			master_nodes_config {
			  ipv4_address   = "10.10.100.172"
			  node_pool_name = "master_node_pool"
			}
		  }
		}

		storage_class_config {
		  reclaim_policy = "Delete"
		  volumes_config {
			flash_mode                 = false
			prism_element_cluster_uuid = local.cluster_id
			storage_container          = "%[6]s"
		  }
		}
		cni_config {
		  dynamic "calico_config" {
			for_each = local.cni == "calico" ? [1] : []
			content {
			  ip_pool_config {
				cidr = "172.20.0.0/16"
			  }
			}
		  }
		}
		worker_node_pool {
		  node_os_version = local.node_os_version
		  num_instances   = local.amount_of_workers
		  ahv_config {
			cpu                        = 8
			disk_mib                   = 122880
			memory_mib                 = 8192
			network_uuid               = data.nutanix_subnet.karbon_subnet.id
			prism_element_cluster_uuid = local.cluster_id
		  }
		}
		etcd_node_pool {
		  node_os_version = local.node_os_version
		  num_instances   = 1
		  ahv_config {
			cpu                        = 4
			disk_mib                   = 40960
			memory_mib                 = 8192
			network_uuid               = data.nutanix_subnet.karbon_subnet.id
			prism_element_cluster_uuid = local.cluster_id
		  }
		}
		master_node_pool {
		  node_os_version = local.node_os_version
		  num_instances   = local.amount_of_masters
		  ahv_config {
			cpu                        = 2
			disk_mib                   = 122880
			memory_mib                 = 4096
			network_uuid               = data.nutanix_subnet.karbon_subnet.id
			prism_element_cluster_uuid = local.cluster_id
		  }
		}
	  }

	`, testVars.NodeOsVersion, workers, cni, subnetName, r, containter, k8sVersion)
}
