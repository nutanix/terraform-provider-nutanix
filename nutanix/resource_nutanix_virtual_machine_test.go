package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spf13/cast"
)

func TestAccNutanixVirtualMachine_basic(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
				),
			},
			{
				Config: testAccNutanixVMConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
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

func TestAccNutanixVirtualMachine_WithDisk(t *testing.T) {
	r := acctest.RandInt()

	resourceName := "nutanix_virtual_machine.vm-disk"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigWithDisk(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "disk_list.#"),
					resource.TestCheckResourceAttr(resourceName, "disk_list.#", "4"),
				),
			},
			{
				Config: testAccNutanixVMConfigWithDiskUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "disk_list.#"),
					resource.TestCheckResourceAttr(resourceName, "disk_list.#", "3"),
				),
			},
			{
				ResourceName:      "nutanix_virtual_machine.vm-disk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}})
}

func TestAccNutanixVirtualMachine_hotadd(t *testing.T) {
	vmName := acctest.RandomWithPrefix("test-dou-vm")
	cpus := 1
	sockets := 1
	memory := 1024
	hotAdd := true
	imageName := acctest.RandomWithPrefix("test-dou-image")

	vmNameUpdated := acctest.RandomWithPrefix("test-dou-vm")
	cpusUpdated := 2
	socketsUpdated := 2
	memoryUpdated := 2048
	hotAddUpdated := false // To force a reboot
	resourceName := "nutanix_virtual_machine.vm10"
	imageNameUpdate := acctest.RandomWithPrefix("test-dou-image")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigHotAdd(vmName, cpus, sockets, memory, hotAdd, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", cast.ToString(memory)),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", cast.ToString(sockets)),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", cast.ToString(cpus)),
					resource.TestCheckResourceAttr(resourceName, "use_hot_add", cast.ToString(hotAdd)),
				),
			},
			{
				Config: testAccNutanixVMConfigHotAdd(vmNameUpdated, cpusUpdated, socketsUpdated, memoryUpdated, hotAddUpdated, imageNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", cast.ToString(memoryUpdated)),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", cast.ToString(socketsUpdated)),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", cast.ToString(cpusUpdated)),
					resource.TestCheckResourceAttr(resourceName, "use_hot_add", cast.ToString(hotAddUpdated))),
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

func TestAccNutanixVirtualMachine_updateFields(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm2"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigUpdatedFields(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-dou-%d", r)),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
				),
			},
			{
				Config: testAccNutanixVMConfigUpdatedFieldsUpdated(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-dou-%d-updated", r)),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "256"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "2"),
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

func TestAccNutanixVirtualMachine_WithSubnet(t *testing.T) {
	r := acctest.RandIntRange(101, 110)
	resourceName := "nutanix_virtual_machine.vm3"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigWithSubnet(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "nic_list_status.0.ip_endpoint_list.0.ip"),
				),
			},
		},
	})
}

func TestAccNutanixVirtualMachine_WithSerialPortList(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm5"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigWithSerialPortList(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "serial_port_list.0.index", "1"),
					resource.TestCheckResourceAttr(resourceName, "serial_port_list.0.is_connected", "true"),
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

func TestAccNutanixVirtualMachine_PowerStateMechanism(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm6"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigPowerStateMechanism(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "power_state_mechanism", "ACPI"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"power_state_mechanism_config", "power_state_mechanism"},
			},
		},
	})
}

func TestAccNutanixVirtualMachine_CdromGuestCustomisationReboot(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm7"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigCdromGuestCustomisationReboot(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
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

func TestAccNutanixVirtualMachine_CloudInitCustomKeyValues(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm8"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigCloudInitCustomKeyValues(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
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

func TestAccNutanixVirtualMachine_DeviceProperties(t *testing.T) {
	r := acctest.RandInt()

	resourceName := "nutanix_virtual_machine.vm9"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigWithDeviceProperties(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "disk_list.#"),
					resource.TestCheckResourceAttr(resourceName, "disk_list.#", "1"),
				),
			},
			{
				ResourceName:      "nutanix_virtual_machine.vm9",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}})
}

func TestAccNutanixVirtualMachine_cloningVM(t *testing.T) {
	r := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigCloningVM(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.vm2"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm2", "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm2", "power_state", "ON"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm2", "memory_size_mib", "186"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm2", "num_sockets", "1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm2", "num_vcpus_per_socket", "1"),
				),
			},
		},
	})
}

func TestAccNutanixVirtualMachine_withPowerMechanism(t *testing.T) {
	resourceName := "nutanix_virtual_machine.myvm"

	vmName := acctest.RandomWithPrefix("vm_name_dou_")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigWithPowerMechanism(vmName, "true", "", "GUEST"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
				),
			},
			{
				Config: testAccNutanixVMConfigWithPowerMechanism(vmName, "", "false", "ACPI"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
				),
			},
			{
				Config: testAccNutanixVMConfigWithPowerMechanism(vmName, "false", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
				),
			},
			{
				Config: testAccNutanixVMConfigWithPowerMechanism(vmName, "true", "true", "HARD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineExists(n string) resource.TestCheckFunc {
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

func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine" {
			continue
		}
		for {
			_, err := conn.API.V3.GetVM(rs.Primary.ID)
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

func testAccNutanixVMConfig(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_virtual_machine" "vm1" {
			name = "test-dou-%d"
			cluster_uuid = "${local.cluster1}"

			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

					num_vcpus_per_socket = 1
					num_sockets          = 1
					memory_size_mib      = 186


					categories {
						name  = "Environment"
						value = "Staging"
					}
				}
	`, r)
}

func testAccNutanixVMConfigWithDisk(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}


		resource "nutanix_image" "cirros-034-disk" {
			name        = "test-image-dou-vm-create-%[1]d"
			source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
			description = "heres a tiny linux image, not an iso, but a real disk!"
		}

		resource "nutanix_virtual_machine" "vm-disk" {
			name                 = "test-dou-vm-%[1]d"
			cluster_uuid         = local.cluster1
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = nutanix_image.cirros-034-disk.id
				}

				device_properties {
					disk_address = {
						device_index = 0
						adapter_type = "SCSI"
					}
					device_type = "DISK"
				}
			}
			disk_list {
				disk_size_mib = 100
			}
			disk_list {
				disk_size_mib = 200
			}
			disk_list {
				disk_size_mib = 300
			}
		}
	`, r)
}

func testAccNutanixVMConfigWithDiskUpdate(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_image" "cirros-034-disk" {
			name        = "test-image-dou-%[1]d"
			source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
			description = "heres a tiny linux image, not an iso, but a real disk!"
		}

		resource "nutanix_virtual_machine" "vm-disk" {
			name                 = "test-dou-vm-%[1]d"
			cluster_uuid         = local.cluster1
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = nutanix_image.cirros-034-disk.id
				}

				device_properties {
					disk_address = {
						device_index = 0
						adapter_type = "SCSI"
					}
					device_type = "DISK"
				}
				disk_size_bytes = 68157440
				disk_size_mib   = 65
			}
			disk_list {
				disk_size_mib = 100
			}
			disk_list {
				disk_size_mib = 200
			}
		}
	`, r)
}

func testAccNutanixVMConfigUpdate(r int) string {
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

			categories {
				name  = "Environment"
				value = "Production"
			}
		}
	`, r)
}

func testAccNutanixVMConfigUpdatedFields(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_virtual_machine" "vm2" {
			name                 = "test-dou-%d"
			cluster_uuid         = local.cluster1
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186


			categories {
				name  = "Environment"
				value = "Staging"
			}
		}
	`, r)
}

func testAccNutanixVMConfigUpdatedFieldsUpdated(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_virtual_machine" "vm2" {
			name                 = "test-dou-%d-updated"
			cluster_uuid         = local.cluster1
			num_vcpus_per_socket = 2
			num_sockets          = 2
			memory_size_mib      = 256

			categories {
				name  = "Environment"
				value = "Production"
			}
		}
	`, r)
}

func testAccNutanixVMConfigWithSubnet(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_subnet" "sub" {
			cluster_uuid = "${local.cluster1}"

			# General Information for subnet
			name        = "terraform-vm-with-subnet-%[1]d"
			description = "Description of my unit test VLAN"
			vlan_id     = %[1]d
			subnet_type = "VLAN"

			# Provision a Managed L3 Network
			# This bit is only needed if you intend to turn on AHV's IPAM
			subnet_ip          = "10.250.140.0"
			default_gateway_ip = "10.250.140.1"
			prefix_length      = 24
			dhcp_options = {
				boot_file_name   = "bootfile"
				domain_name      = "nutanix"
				tftp_server_name = "10.250.140.200"
			}
			dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
			dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
			ip_config_pool_list_ranges   = ["10.250.140.20 10.250.140.100"]
		}

		resource "nutanix_image" "cirros-034-disk" {
			name        = "test-image-dou-%[1]d"
			source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
			description = "heres a tiny linux image, not an iso, but a real disk!"
		}

		resource "nutanix_virtual_machine" "vm3" {
			name = "test-dou-vm-%[1]d"

			categories {
				name  = "Environment"
				value = "Staging"
			}

			cluster_uuid         = "${local.cluster1}"
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = "${nutanix_image.cirros-034-disk.id}"
				}
			}

			nic_list {
				subnet_uuid = "${nutanix_subnet.sub.id}"
			}
		}

		output "ip_address" {
			value = "${lookup(nutanix_virtual_machine.vm3.nic_list_status.0.ip_endpoint_list[0], "ip")}"
		}
	`, r)
}

func testAccNutanixVMConfigWithSerialPortList(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_virtual_machine" "vm5" {
			name = "test-dou-%d"
			cluster_uuid = "${local.cluster1}"

			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			serial_port_list {
				index = 1
				is_connected = true
			}

			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			serial_port_list {
				index        = 1
				is_connected = true
			}

			categories {
				name  = "Environment"
				value = "Staging"
			}
		}
	`, r)
}

func testAccNutanixVMConfigPowerStateMechanism(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_virtual_machine" "vm6" {
			name         = "test-dou-%d"
			cluster_uuid = "${local.cluster1}"

			num_vcpus_per_socket  = 1
			num_sockets           = 1
			memory_size_mib       = 186
			power_state_mechanism = "ACPI"
		}
	`, r)
}

func testAccNutanixVMConfigCdromGuestCustomisationReboot(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_virtual_machine" "vm7" {
			name         = "test-dou-%d"
			cluster_uuid = "${local.cluster1}"

			num_vcpus_per_socket                     = 1
			num_sockets                              = 1
			memory_size_mib                          = 186
			guest_customization_cloud_init_user_data = base64encode("#cloud-config\nfqdn: test.domain.local")
		}
	`, r)
}

func testAccNutanixVMConfigCloudInitCustomKeyValues(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_virtual_machine" "vm8" {
			name = "test-dou-%d"
			cluster_uuid = "${local.cluster1}"

			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			guest_customization_cloud_init_custom_key_values = {
				"username" = "myuser"
				"password" = "mypassword"
			}
		}
	`, r)
}

func testAccNutanixVMConfigWithDeviceProperties(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_image" "cirros-034-disk" {
			name        = "test-image-dou-vm-create-%[1]d"
			source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
			description = "heres a tiny linux image, not an iso, but a real disk!"
		}

		resource "nutanix_virtual_machine" "vm9" {
			name                 = "test-dou-vm-%[1]d"
			cluster_uuid         = local.cluster1
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = nutanix_image.cirros-034-disk.id
				}

				device_properties {
					device_type = "DISK"
				}
			}
		}
	`, r)
}

func testAccNutanixVMConfigCloningVM(r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_image" "cirros-034-disk" {
			name        = "test-image-dou-vm-create-%[1]d"
			source_uri  = "http://packages.oit.ncsu.edu/centos/7.7.1908/isos/x86_64/CentOS-7-x86_64-NetInstall-1908.iso"
			description = "heres a tiny linux image, not an iso, but a real disk!"
		}

		resource "nutanix_virtual_machine" "vm1" {
			name         = "test-dou-%[1]d-vm1"
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186
			cluster_uuid = "${local.cluster1}"


			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = nutanix_image.cirros-034-disk.id
				}

				device_properties {
					disk_address = {
						device_index = 0
						adapter_type = "IDE"
					}
					device_type = "CDROM"
				}
			}
			disk_list {
				disk_size_mib = 100
			}
			disk_list {
				disk_size_mib = 200
			}
			disk_list {
				disk_size_mib = 300
			}
		}

		data "nutanix_virtual_machine" "vmds" {
			vm_id = "${nutanix_virtual_machine.vm1.id}"
		}

		resource "nutanix_virtual_machine" "vm2" {
			name         = "test-dou-%[1]d-vm2"
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186
			cluster_uuid = "${local.cluster1}"


			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = "${data.nutanix_virtual_machine.vmds.disk_list.0.data_source_reference.uuid}"
				}
			}
		}
	`, r)
}

func testAccNutanixVMConfigHotAdd(vmName string, cpus, sockets, memory int, hotAdd bool, imageName string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_image" "cirros-034-disk" {
			name        = "%[6]s"
			source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
			description = "heres a tiny linux image, not an iso, but a real disk!"
		}

		resource "nutanix_virtual_machine" "vm10" {
			name         = "%[1]s"
			cluster_uuid = "${local.cluster1}"
			num_vcpus_per_socket  = %[2]d
			num_sockets           = %[3]d
			memory_size_mib       = %[4]d
			power_state_mechanism = "ACPI"
			use_hot_add           = %[5]v

			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = nutanix_image.cirros-034-disk.id
				}

				device_properties {
					disk_address = {
						device_index = 0
						adapter_type = "SCSI"
					}
					device_type = "DISK"
				}
			}
		}
	`, vmName, cpus, sockets, memory, hotAdd, imageName)
}

func testAccNutanixVMConfigWithPowerMechanism(vmName, scriptFail, scriptExec, mechanism string) string {
	var comFail, comExec, comMec, transitionConfig string
	if scriptFail != "" {
		comFail = fmt.Sprintf("should_fail_on_script_failure = %s", scriptFail)
	}
	if scriptExec != "" {
		comExec = fmt.Sprintf("enable_script_exec            = %s", scriptExec)
	}
	if scriptFail != "" || scriptExec != "" {
		transitionConfig = fmt.Sprintf(`
		guest_transition_config {
			%s
			%s
		}
			`, comFail, comExec)
	}
	if mechanism != "" {
		comMec = fmt.Sprintf(`mechanism = "%s"`, mechanism)
	}

	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
				cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
				? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_virtual_machine" "myvm" {
			name         = "%[1]s"
			cluster_uuid = "${local.cluster1}"

			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			power_state_mechanism_config {
				%s

				%s
			}

		}
	`, vmName, transitionConfig, comMec)
}
