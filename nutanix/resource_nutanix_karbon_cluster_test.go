package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNutanixKarbonCluster_basic(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixKarbonClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixKarbonClusterConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
				),
			},
			{
				Config: testAccNutanixKarbonClusterConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixKarbonClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_list"},
			},
		},
	})
}

func testAccCheckNutanixKarbonClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_karbon_cluster" {
			continue
		}
		for {
			_, err := conn.KarbonAPI.Cluster.GetKarbonCluster(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
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

func testAccNutanixKarbonClusterConfig(subnetName string, r int) string {
	return fmt.Sprintf(`
	locals {
		cluster_id = data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL" ? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid
	  }
	  
	data "nutanix_clusters" "clusters" {}
	  
	data "nutanix_subnet" "karbon_subnet" {
		subnet_name = %s
	}

	resource "nutanix_karbon_cluster" "cluster" {
		# depends_on = [nutanix_karbon_private_registry.registry]
		name    = var.karbon_cluster_name
		version = var.k8sversion
	  
		# private_registry {
		#   registry_name = nutanix_karbon_private_registry.registry.name
		# }
	  
		dynamic "active_passive_config" {
		  for_each = var.deployment_type == "active-passive" ? [1] : []
		  content {
			external_ipv4_address = var.master_vip
		  }
		}
		dynamic "external_lb_config" {
		  for_each = var.deployment_type == "active-active" ? [1] : []
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
			storage_container          = var.storage_container
		  }
		}
		cni_config {
		  dynamic "calico_config" {
			for_each = var.cni == "calico" ? [1] : []
			content {
			  ip_pool_config {
				cidr = "172.20.0.0/16"
			  }
			}
		  }
		}
		worker_node_pool {
		  node_os_version = var.node_os_version
		  num_instances   = var.amount_of_workers
		  ahv_config {
			cpu                        = 8
			disk_mib                   = 122880
			memory_mib                 = 8192
			network_uuid               = data.nutanix_subnet.karbon_subnet.id
			prism_element_cluster_uuid = local.cluster_id
		  }
		}
		etcd_node_pool {
		  node_os_version = var.node_os_version
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
		  node_os_version = var.node_os_version
		  num_instances   = var.amount_of_masters
		  ahv_config {
			cpu                        = 2
			disk_mib                   = 122880
			memory_mib                 = 4096
			network_uuid               = data.nutanix_subnet.karbon_subnet.id
			prism_element_cluster_uuid = local.cluster_id
		  }
		}
	  }


	`, subnetName)
}

func testAccNutanixKarbonClusterConfigUpdate(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_virtual_machine" "vm1" {
			name                 = "test-dou-%d"
			cluster_uuid         = "${local.cluster1}"
			num_vcpus_per_socket = 1
			num_sockets          = 2
			memory_size_mib      = 186

			boot_device_order_list = ["DISK", "CDROM"]

			categories {
				name  = "Environment"
				value = "Production"
			}
		}
	`, r)
}
