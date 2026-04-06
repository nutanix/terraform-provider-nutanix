package networkingv2_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameNetworkFunctionV2_1 = "nutanix_network_function_v2.ntf-1"

func TestAccV2NutanixNetworkFunctionResource_Egress_Ingress(t *testing.T) {
	r := acctest.RandInt()
	vmmResourceName_1 := "nutanix_virtual_machine_v2.vm-1"
	vmmResourceName_2 := "nutanix_virtual_machine_v2.vm-2"

	subnet_name := fmt.Sprintf("tf-test-subnet-%d", r)
	vmm_1_name := fmt.Sprintf("tf-test-vm-1-%d", r)
	vmm_2_name := fmt.Sprintf("tf-test-vm-2-%d", r)
	name := fmt.Sprintf("tf-test-network-function-%d", r)

	networkFunctionConfig := testAccNetworkFunctionV2ConfigPrerequisites(subnet_name, vmm_1_name, vmm_2_name) +
		testAccNetworkFunctionV2EgressIngressConfig(name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNetworkFunctionResourcesDestroy,
		Steps: []resource.TestStep{
			// Prerequisites
			{
				Config: testAccNetworkFunctionV2ConfigPrerequisites(subnet_name, vmm_1_name, vmm_2_name),
				Check: resource.ComposeTestCheckFunc(
					//  vmm 1 checks
					resource.TestCheckResourceAttrSet(vmmResourceName_1, "id"),
					resource.TestCheckResourceAttr(vmmResourceName_1, "name", vmm_1_name),
					resource.TestCheckResourceAttrSet(vmmResourceName_1, "nics.#"),
					resource.TestCheckResourceAttrSet(vmmResourceName_1, "nics.0.ext_id"),
					resource.TestCheckResourceAttr(vmmResourceName_1, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.network_function_nic_type", "INGRESS"),
					resource.TestCheckResourceAttrSet(vmmResourceName_1, "nics.1.ext_id"),
					resource.TestCheckResourceAttr(vmmResourceName_1, "nics.1.nic_network_info.0.virtual_ethernet_nic_network_info.0.network_function_nic_type", "EGRESS"),
					resource.TestCheckResourceAttrSet(vmmResourceName_1, "nics.2.ext_id"),
					resource.TestCheckResourceAttr(vmmResourceName_1, "nics.2.nic_network_info.0.virtual_ethernet_nic_network_info.0.network_function_nic_type", "INGRESS"),
					resource.TestCheckResourceAttrSet(vmmResourceName_1, "nics.3.ext_id"),
					resource.TestCheckResourceAttr(vmmResourceName_1, "nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttrSet(vmmResourceName_1, "nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.subnet.0.ext_id"),
					resource.TestCheckResourceAttrSet(vmmResourceName_1, "nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_info.0.learned_ip_addresses.0.value"),
					//  vmm 2 checks
					resource.TestCheckResourceAttrSet(vmmResourceName_2, "id"),
					resource.TestCheckResourceAttr(vmmResourceName_2, "name", vmm_2_name),
					resource.TestCheckResourceAttrSet(vmmResourceName_2, "nics.#"),
					resource.TestCheckResourceAttrSet(vmmResourceName_2, "nics.0.ext_id"),
					resource.TestCheckResourceAttr(vmmResourceName_2, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.network_function_nic_type", "INGRESS"),
					resource.TestCheckResourceAttrSet(vmmResourceName_2, "nics.1.ext_id"),
					resource.TestCheckResourceAttr(vmmResourceName_2, "nics.1.nic_network_info.0.virtual_ethernet_nic_network_info.0.network_function_nic_type", "EGRESS"),
					resource.TestCheckResourceAttrSet(vmmResourceName_2, "nics.2.ext_id"),
					resource.TestCheckResourceAttr(vmmResourceName_2, "nics.2.nic_network_info.0.virtual_ethernet_nic_network_info.0.network_function_nic_type", "INGRESS"),
					resource.TestCheckResourceAttrSet(vmmResourceName_2, "nics.3.ext_id"),
					resource.TestCheckResourceAttr(vmmResourceName_2, "nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttrSet(vmmResourceName_2, "nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.subnet.0.ext_id"),
					resource.TestCheckResourceAttrSet(vmmResourceName_2, "nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_info.0.learned_ip_addresses.0.value"),
				),
			},
			// Create Network Function
			{
				Config: networkFunctionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkFunctionV2_1, "id"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "name", name),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "description", "First Network function managed by Terraform"),
					resource.TestCheckResourceAttr(
						resourceNameNetworkFunctionV2_1,
						"description",
						"First Network function managed by Terraform",
					),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "high_availability_mode", "ACTIVE_PASSIVE"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "traffic_forwarding_mode", "INLINE"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "failure_handling", "FAIL_CLOSE"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.#", "2"),
					testAccCheckNetworkFunctionPair(
						resourceNameNetworkFunctionV2_1,
						"nutanix_virtual_machine_v2.vm-1",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					testAccCheckNetworkFunctionPair(
						resourceNameNetworkFunctionV2_1,
						"nutanix_virtual_machine_v2.vm-2",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "data_plane_health_check_config.0.failure_threshold", "3"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "data_plane_health_check_config.0.interval_secs", "4"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "data_plane_health_check_config.0.success_threshold", "3"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "data_plane_health_check_config.0.timeout_secs", "2"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "high_availability_mode", "ACTIVE_PASSIVE"),
					// wait until the network function data_plane_health_status is HEALTHY
					waitForNetworkFunctionHealth(resourceNameNetworkFunctionV2_1, "data_plane_health_status", "HEALTHY"),
				),
			},
			// Read Network Function
			{
				Config: networkFunctionConfig + `
				  data "nutanix_network_function_v2" "ntf" {
					ext_id = nutanix_network_function_v2.ntf-1.id
				  }
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_network_function_v2.ntf", "id"),
					resource.TestCheckResourceAttr("data.nutanix_network_function_v2.ntf", "name", name),
					resource.TestCheckResourceAttr("data.nutanix_network_function_v2.ntf", "high_availability_mode", "ACTIVE_PASSIVE"),
					resource.TestCheckResourceAttr("data.nutanix_network_function_v2.ntf", "nic_pairs.#", "2"),
					testAccCheckNetworkFunctionDataSourcePair(
						"data.nutanix_network_function_v2.ntf",
						"nutanix_virtual_machine_v2.vm-1",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					testAccCheckNetworkFunctionDataSourcePair(
						"data.nutanix_network_function_v2.ntf",
						"nutanix_virtual_machine_v2.vm-2",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					resource.TestCheckResourceAttr("data.nutanix_network_function_v2.ntf", "nic_pairs.0.data_plane_health_status", "HEALTHY"),
					resource.TestCheckResourceAttr("data.nutanix_network_function_v2.ntf", "nic_pairs.1.data_plane_health_status", "HEALTHY"),
					testAccCheckNetworkFunctionNICPairHAStateCounts(
						"data.nutanix_network_function_v2.ntf",
						map[string]int{
							"ACTIVE":  1,
							"PASSIVE": 1,
						},
					),
				),
			},
			// Update Network Function
			{
				Config: testAccNetworkFunctionV2ConfigPrerequisites(subnet_name, vmm_1_name, vmm_2_name) + testAccNetworkFunctionV2EgressIngressUpdateConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkFunctionV2_1, "id"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "name", name+"_updated"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "description", "First Network function managed by Terraform updated"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "high_availability_mode", "ACTIVE_PASSIVE"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "traffic_forwarding_mode", "INLINE"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "failure_handling", "FAIL_OPEN"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "data_plane_health_check_config.0.interval_secs", "3"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "data_plane_health_check_config.0.timeout_secs", "3"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "data_plane_health_check_config.0.success_threshold", "2"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "data_plane_health_check_config.0.failure_threshold", "2"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.#", "1"),
					testAccCheckNetworkFunctionPair(
						resourceNameNetworkFunctionV2_1,
						"nutanix_virtual_machine_v2.vm-1",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.0.is_enabled", "true"),
				),
			},
		},
	})
}

func TestAccV2NutanixNetworkFunctionResource_Ingress(t *testing.T) {
	r := acctest.RandInt()
	subnet_name := fmt.Sprintf("tf-test-subnet-%d", r)
	vmm_1_name := fmt.Sprintf("tf-test-vm-1-%d", r)
	vmm_2_name := fmt.Sprintf("tf-test-vm-2-%d", r)
	name := fmt.Sprintf("tf-test-network-function-%d", r)

	networkFunctionConfig := testAccNetworkFunctionV2ConfigPrerequisites(subnet_name, vmm_1_name, vmm_2_name) +
		testAccNetworkFunctionV2IngressConfig(name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNetworkFunctionResourcesDestroy,
		Steps: []resource.TestStep{
			// Prerequisites
			{
				Config: testAccNetworkFunctionV2ConfigPrerequisites(subnet_name, vmm_1_name, vmm_2_name),
			},
			{
				Config: networkFunctionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkFunctionV2_1, "id"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "name", name),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "high_availability_mode", "ACTIVE_PASSIVE"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "traffic_forwarding_mode", "VTAP"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.#", "2"),
					resource.TestCheckResourceAttrPair(
						resourceNameNetworkFunctionV2_1,
						"nic_pairs.0.vm_reference",
						"nutanix_virtual_machine_v2.vm-1",
						"id",
					),
					resource.TestCheckResourceAttrPair(
						resourceNameNetworkFunctionV2_1,
						"nic_pairs.0.ingress_nic_reference",
						"nutanix_virtual_machine_v2.vm-1",
						"nics.2.ext_id",
					),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.0.is_enabled", "true"),
					resource.TestCheckResourceAttrPair(
						resourceNameNetworkFunctionV2_1,
						"nic_pairs.1.vm_reference",
						"nutanix_virtual_machine_v2.vm-2",
						"id",
					),
					resource.TestCheckResourceAttrPair(
						resourceNameNetworkFunctionV2_1,
						"nic_pairs.1.ingress_nic_reference",
						"nutanix_virtual_machine_v2.vm-2",
						"nics.2.ext_id",
					),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.1.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "high_availability_mode", "ACTIVE_PASSIVE"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "traffic_forwarding_mode", "VTAP"),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.#", "2"),
					resource.TestCheckResourceAttrPair(
						resourceNameNetworkFunctionV2_1,
						"nic_pairs.0.vm_reference",
						"nutanix_virtual_machine_v2.vm-1",
						"id",
					),
					resource.TestCheckResourceAttrPair(
						resourceNameNetworkFunctionV2_1,
						"nic_pairs.0.ingress_nic_reference",
						"nutanix_virtual_machine_v2.vm-1",
						"nics.2.ext_id",
					),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.0.is_enabled", "true"),
					resource.TestCheckResourceAttrPair(
						resourceNameNetworkFunctionV2_1,
						"nic_pairs.1.vm_reference",
						"nutanix_virtual_machine_v2.vm-2",
						"id",
					),
					resource.TestCheckResourceAttrPair(
						resourceNameNetworkFunctionV2_1,
						"nic_pairs.1.ingress_nic_reference",
						"nutanix_virtual_machine_v2.vm-2",
						"nics.2.ext_id",
					),
					resource.TestCheckResourceAttr(resourceNameNetworkFunctionV2_1, "nic_pairs.1.is_enabled", "true"),
				),
			},
		},
	})
}

func TestAccV2NutanixNetworkFunctionResource_InvalidConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNetworkFunctionResourcesDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "nutanix_network_function_v2" "ntf-1" {
            name = "tf-test-network-function-invalid-config"
            high_availability_mode = "ACTIVE_PASSIVE"
            traffic_forwarding_mode = "VTAP"
            nic_pairs {
              ingress_nic_reference = "00000000-0000-0000-0000-000000000000"
              vm_reference = "a5555555-5555-5555-5555-555555555555"
              is_enabled = false
            }
          }
        `,
				ExpectError: regexp.MustCompile("Failed to Create Network Function due to an invalid input parameter - Failed to configure network function"),
			},
		},
	})
}

// testAccNetworkFunctionV2ConfigPrerequisitesNoPostcondition returns the same
// prerequisites as testAccNetworkFunctionV2ConfigPrerequisites but with VM
// postconditions removed. Use for tests that only need VMs and NF to exist
// (e.g. NSP with network_function_reference) and cannot wait for DHCP on the
// VM NORMAL_NIC.
func testAccNetworkFunctionV2ConfigPrerequisitesNoPostcondition(subnet_name, vmm_1_name, vmm_2_name string) string {
	s := testAccNetworkFunctionV2ConfigPrerequisites(subnet_name, vmm_1_name, vmm_2_name)
	// Replace full lifecycle block (including postcondition) with lifecycle that only has ignore_changes.
	lifecycleWithPostcondition1 := `  lifecycle {
    ignore_changes = [cd_roms, guest_customization, nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config]

    postcondition {
      condition = length([
        for nic in self.nics : 1
        if try(length(nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].ipv4_info[0].learned_ip_addresses[0].value), 0) > 0
      ]) > 0
      error_message = "The first VM must have at least one NIC with an assigned IPv4 address."
    }
  }
`
	lifecycleWithPostcondition2 := `  lifecycle {
    ignore_changes = [cd_roms, guest_customization, nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config]
    postcondition {
      condition = length([
        for nic in self.nics : 1
        if try(length(nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].ipv4_info[0].learned_ip_addresses[0].value), 0) > 0
      ]) > 0
      error_message = "The second VM must have at least one NIC with an assigned IPv4 address."
    }
  }
`
	lifecycleIgnoreOnly := `  lifecycle {
    ignore_changes = [cd_roms, guest_customization, nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config]
  }
`
	s = strings.Replace(s, lifecycleWithPostcondition1, lifecycleIgnoreOnly, 1)
	s = strings.Replace(s, lifecycleWithPostcondition2, lifecycleIgnoreOnly, 1)
	return s
}

func testAccNetworkFunctionV2ConfigPrerequisites(subnet_name, vmm_1_name, vmm_2_name string) string {
	return fmt.Sprintf(`
locals {
  config             = jsondecode(file("%[1]s"))
  aosFilter           = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

data "nutanix_clusters_v2" "clusters" {
  filter = local.aosFilter
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
  img_name       = local.config.ubuntu_image
  gz_ntf         = <<EOT
#cloud-config
chpasswd:
  list: |
    ubuntu:nutanix/4u
  expire: false
disable_root: false
ssh_pwauth:   true
runcmd:
  - last_iface=$(ls /sys/class/net/ | grep -E '^e' | sort | tail -1)
  - ip link set dev $last_iface up
  - dhclient -v $last_iface
  - apt-get update
  - apt-get install -y bridge-utils
  - iface1=$(ls /sys/class/net/ | grep -E '^e' | sort | head -1)
  - iface2=$(ls /sys/class/net/ | grep -E '^e' | sort | head -2 | tail -1)
  - brctl addbr br0
  - brctl addif br0 $iface1
  - brctl addif br0 $iface2
  - ip link set dev $iface1 up
  - ip link set dev $iface2 up
  - ip link set dev br0 up
  EOT
}

data "nutanix_images_v2" "vm_img" {
  filter = "name eq '${local.img_name}'"
}

# VLAN subnet with advanced networking (required for network function NICs)
resource "nutanix_subnet_v2" "subnet" {
  name                   = "%[2]s"
  description            = "Subnet managed by Terraform for Network Function testing"
  subnet_type            = "VLAN"
  network_id             = 800
  cluster_reference      = local.cluster_ext_id
  is_advanced_networking = true
}

# Create VM to test the first network function
resource "nutanix_virtual_machine_v2" "vm-1" {
  name        = "%[3]s"
  description = "First VM for testing the network function"
  cluster {
    ext_id = local.cluster_ext_id
  }
  num_cores_per_socket         = 2
  num_sockets                  = 2
  memory_size_bytes            = 4 * 1024 * 1024 * 1024 # 4GB
  is_agent_vm                  = false
  hardware_clock_timezone      = "UTC"
  is_memory_overcommit_enabled = false
  apc_config {
    is_apc_enabled = false
  }
  disks {
    backing_info {
      vm_disk {
        disk_size_bytes = 20 * 1024 * 1024 * 1024 #20GB
        data_source {
          reference {
            image_reference {
              image_ext_id = data.nutanix_images_v2.vm_img.images[0].ext_id
            }
          }
        }
      }
    }
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
  }
  guest_customization {
    config {
      cloud_init {
        cloud_init_script {
          user_data {
            value = base64encode(local.gz_ntf)
          }
        }
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type                  = "NETWORK_FUNCTION_NIC"
        network_function_nic_type = "INGRESS"
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type                  = "NETWORK_FUNCTION_NIC"
        network_function_nic_type = "EGRESS"
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type                  = "NETWORK_FUNCTION_NIC"
        network_function_nic_type = "INGRESS"
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = nutanix_subnet_v2.subnet.id
        }
        ipv4_config {
          should_assign_ip = true
        }
      }
    }
  }

  power_state = "ON"
  lifecycle {
    ignore_changes = [cd_roms, guest_customization, nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config]

    postcondition {
      condition = length([
        for nic in self.nics : 1
        if try(length(nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].ipv4_info[0].learned_ip_addresses[0].value), 0) > 0
      ]) > 0
      error_message = "The first VM must have at least one NIC with an assigned IPv4 address."
    }
  }
}


locals {
  vm_ext_id_1 = nutanix_virtual_machine_v2.vm-1.id
}

# Create second VM to test the second network function
# Create VM to test the first network function
resource "nutanix_virtual_machine_v2" "vm-2" {
  name        = "%[4]s"
  description = "Second VM for testing the network function"
  cluster {
    ext_id = local.cluster_ext_id
  }
  num_cores_per_socket         = 2
  num_sockets                  = 2
  memory_size_bytes            = 4 * 1024 * 1024 * 1024 # 4GB
  is_agent_vm                  = false
  hardware_clock_timezone      = "UTC"
  is_memory_overcommit_enabled = false
  apc_config {
    is_apc_enabled = false
  }
  disks {
    backing_info {
      vm_disk {
        disk_size_bytes = 20 * 1024 * 1024 * 1024 #20GB
        data_source {
          reference {
            image_reference {
              image_ext_id = data.nutanix_images_v2.vm_img.images[0].ext_id
            }
          }
        }
      }
    }
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
  }
  guest_customization {
    config {
      cloud_init {
        cloud_init_script {
          user_data {
            value = base64encode(local.gz_ntf)
          }
        }
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type                  = "NETWORK_FUNCTION_NIC"
        network_function_nic_type = "INGRESS"
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type                  = "NETWORK_FUNCTION_NIC"
        network_function_nic_type = "EGRESS"
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type                  = "NETWORK_FUNCTION_NIC"
        network_function_nic_type = "INGRESS"
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = nutanix_subnet_v2.subnet.id
        }
        ipv4_config {
          should_assign_ip = true
        }

      }
    }
  }

  power_state = "ON"
  lifecycle {
    ignore_changes = [cd_roms, guest_customization, nics.3.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config]
    postcondition {
      condition = length([
        for nic in self.nics : 1
        if try(length(nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].ipv4_info[0].learned_ip_addresses[0].value), 0) > 0
      ]) > 0
      error_message = "The second VM must have at least one NIC with an assigned IPv4 address."
    }
  }
}

locals {
  vm_ext_id_2 = nutanix_virtual_machine_v2.vm-2.id
}

# set the INGRESS and EGRESS NIC ext_id
locals {
  # Set INGRESS and EGRESS NIC ext_ids for both VMs
  vm_1_ingress_1_ext_id = element(
    [
      for nic in nutanix_virtual_machine_v2.vm-1.nics : nic
      if nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type == "INGRESS"
    ],
    0
  ).ext_id

  vm_1_ingress_2_ext_id = element(
    [
      for nic in nutanix_virtual_machine_v2.vm-1.nics : nic
      if nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type == "INGRESS"
    ],
    length([
      for nic in nutanix_virtual_machine_v2.vm-1.nics : nic
      if nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type == "INGRESS"
    ]) - 1
  ).ext_id

  vm_1_egress_1_ext_id = element(
    [
      for nic in nutanix_virtual_machine_v2.vm-1.nics : nic
      if nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type == "EGRESS"
    ],
    0
  ).ext_id

  vm_2_ingress_1_ext_id = element(
    [
      for nic in nutanix_virtual_machine_v2.vm-2.nics : nic
      if nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type == "INGRESS"
    ],
    0
  ).ext_id

  vm_2_ingress_2_ext_id = element(
    [
      for nic in nutanix_virtual_machine_v2.vm-2.nics : nic
      if nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type == "INGRESS"
    ],
    length([
      for nic in nutanix_virtual_machine_v2.vm-2.nics : nic
      if nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type == "INGRESS"
    ]) - 1
  ).ext_id

  vm_2_egress_1_ext_id = element(
    [
      for nic in nutanix_virtual_machine_v2.vm-2.nics : nic
      if nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type == "EGRESS"
    ],
    0
  ).ext_id
}

	`, filepath, subnet_name, vmm_1_name, vmm_2_name)
}

func testAccNetworkFunctionV2EgressIngressConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_network_function_v2" "ntf-1" {
  name                    = "%[1]s"
  description             = "First Network function managed by Terraform"
  high_availability_mode  = "ACTIVE_PASSIVE"
  failure_handling        = "FAIL_CLOSE"
  traffic_forwarding_mode = "INLINE"
  nic_pairs {
    ingress_nic_reference = local.vm_1_ingress_1_ext_id
    egress_nic_reference  = local.vm_1_egress_1_ext_id
    vm_reference          = local.vm_ext_id_1
    is_enabled            = true
  }
  nic_pairs {
    ingress_nic_reference = local.vm_2_ingress_1_ext_id
    egress_nic_reference  = local.vm_2_egress_1_ext_id
    vm_reference          = local.vm_ext_id_2
    is_enabled            = true
  }
  data_plane_health_check_config {
    interval_secs     = 4
    timeout_secs      = 2
    success_threshold = 3
    failure_threshold = 3
  }
}
`, name)
}

func testAccNetworkFunctionV2EgressIngressUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_network_function_v2" "ntf-1" {
  name                    = "%[1]s_updated"
  description             = "First Network function managed by Terraform updated"
  high_availability_mode  = "ACTIVE_PASSIVE"
  failure_handling        = "FAIL_OPEN"
  traffic_forwarding_mode = "INLINE"
  nic_pairs {
    ingress_nic_reference = local.vm_1_ingress_1_ext_id
    egress_nic_reference  = local.vm_1_egress_1_ext_id
    vm_reference          = local.vm_ext_id_1
    is_enabled            = true
  }
  data_plane_health_check_config {
    interval_secs     = 3
    timeout_secs      = 3
    success_threshold = 2
    failure_threshold = 2
  }
}
`, name)
}

func testAccNetworkFunctionV2IngressConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_network_function_v2" "ntf-1" {
  name                    = "%[1]s"
  description             = "Second Network function managed by Terraform"
  high_availability_mode  = "ACTIVE_PASSIVE"
  traffic_forwarding_mode = "VTAP"
  nic_pairs {
    ingress_nic_reference = local.vm_1_ingress_2_ext_id
    vm_reference          = local.vm_ext_id_1
    is_enabled            = true
  }
  nic_pairs {
    ingress_nic_reference = local.vm_2_ingress_2_ext_id
    vm_reference          = local.vm_ext_id_2
    is_enabled            = true
  }
}`, name)
}
