package vmmv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	resourceNameNGTInsertISO = "nutanix_ngt_insert_iso_v2.test"
	resourceVm               = "nutanix_virtual_machine_v2.ngt-vm"
	datasourceVmNGT          = "data.nutanix_virtual_machine_v2.ngt-vm-refresh"
	timeSleep                = 2 * time.Minute
)

func TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmHaveNGTTest(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Log("Creating and Powering on the VM")
				},
				Config: testPreEnvConfig(vmName, r),
			},
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to power on")
					time.Sleep(timeSleep)
					t.Log("Installing NGT")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.1", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "IMMEDIATE"),
				),
			},
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to reboot")
					time.Sleep(timeSleep)
					t.Log("Inserting NGT Iso")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTInsertIsoConfig("true", "insert"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "guest_os_version"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmHaveNGTIsConfigFalse(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Log("Creating and Powering on the VM")
				},
				Config: testPreEnvConfig(vmName, r),
			},
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to power on")
					time.Sleep(timeSleep)
					t.Log("Installing NGT")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.1", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "IMMEDIATE"),
				),
			},
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to reboot")
					time.Sleep(timeSleep)
					t.Log("Inserting NGT Iso")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTInsertIsoConfig("false", "insert"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "guest_os_version"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmDoseNotHaveNGT(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Log("Creating and Powering on the VM")
				},
				Config: testPreEnvConfig(vmName, r),
			},
			{
				Config: testPreEnvConfig(vmName, r) + testNGTInsertIsoConfig("true", "insert"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "available_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "guest_os_version", ""),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_installed", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_reachable", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmDoseNotHaveNGTIsConfigFalse(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineV2Destroy,
		Steps: []resource.TestStep{
			// Step 1: create the VM
			{
				PreConfig: func() {
					fmt.Println("Step 1: Creating and Powering on the VM")
				},
				Config: testVmConfig(vmName),
			},
			// Step 2: insert the NGT ISO on vm
			{
				PreConfig: func() {
					fmt.Println("Step 2: Inserting the NGT ISO on vm resource")
				},
				Config: testVmConfig(vmName) + testNGTInsertIsoConfig("false", "insert"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "guest_tools.0.is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "guest_tools.0.is_installed", "false"),
					resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_enabled", datasourceVmNGT, "guest_tools.0.is_enabled"),
					resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_installed", datasourceVmNGT, "guest_tools.0.is_installed"),
					resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_reachable", datasourceVmNGT, "guest_tools.0.is_reachable"),
					resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "version", datasourceVmNGT, "guest_tools.0.version"),
					resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "guest_os_version", datasourceVmNGT, "guest_tools.0.guest_os_version"),
					resource.TestCheckResourceAttrPair(datasourceVmNGT, "guest_tools.0.capabilities.#", resourceNameNGTInsertISO, "capablities.#"),
					resource.TestCheckResourceAttrPair(datasourceVmNGT, "guest_tools.0.capabilities.0", resourceNameNGTInsertISO, "capablities.0"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "cd_roms.0.iso_type", "GUEST_TOOLS"),
				),
			},
			// Step 3: NGT Installation
			{
				PreConfig: func() {
					fmt.Println("Step 3: NGT Installation")
				},
				Config: testVmConfig(vmName) +  testNGTInsertIsoConfig("false", "insert") + testNGTInstallationConfig(),
			},
			// Step 4: check the NGT is installed
			{
				PreConfig: func() {
					fmt.Println("Step 4: Checking the NGT is installed")
				},
				Config: testVmConfig(vmName) +  testNGTInsertIsoConfig("false", "insert") + testNGTInstallationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_installed", "true"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceVm, "guest_tools.0.version"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_iso_inserted", "false"),
					resource.TestCheckResourceAttrSet(resourceVm, "guest_tools.0.guest_os_version"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.capabilities.#", "1"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.capabilities.1", "VSS_SNAPSHOT"),
				),
			},
			// Step 5: eject the NGT ISO from the vm
			{
				PreConfig: func() {
					fmt.Println("Step 5: Ejecting the NGT ISO from the vm resource")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInsertIsoConfig("false", "eject"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceVmNGT, "guest_tools.0.is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "guest_tools.0.is_installed", "true"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "guest_tools.0.is_reachable", "true"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "guest_tools.0.is_enabled", "true"),
					resource.TestCheckResourceAttrSet(datasourceVmNGT, "guest_tools.0.version"),
					resource.TestCheckResourceAttrSet(datasourceVmNGT, "guest_tools.0.guest_os_version"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "guest_tools.0.capabilities.#", "1"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "guest_tools.0.capabilities.0", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(datasourceVmNGT, "cd_roms.0.iso_type", "OTHER"),
				),
			},
			// Step 6: check the NGT ISO is ejected from the vm resource
			{
				PreConfig: func() {
					fmt.Println("Step 6: Checking the NGT ISO is ejected from the vm resource")
				},
				Config: testVmConfig(vmName) + testNGTInsertIsoConfig("false", "eject"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_installed", "true"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.is_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceVm, "guest_tools.0.version"),
					resource.TestCheckResourceAttrSet(resourceVm, "guest_tools.0.guest_os_version"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.capabilities.#", "2"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.capabilities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceVm, "guest_tools.0.capabilities.1", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(resourceVm, "cd_roms.0.iso_type", "OTHER"),
				),
			},
		},
	})
}

func testVmConfig(vmName string) string {
	return fmt.Sprintf(`
	locals {
		config = (jsondecode(file("%[1]s")))
		vmm    = local.config.vmm
		aosFilter           = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
	}
	data "nutanix_clusters_v2" "clusters" {
	  filter = local.aosFilter
	}
	
	locals {
	  clusterUUID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
	}
	
	data "nutanix_images_v2" "ngt-image" {
	  filter = "name eq '${local.vmm.image_name}'"
	}
	
	data "nutanix_storage_containers_v2" "ngt-sc" {
	  filter = "clusterExtId eq '${local.clusterUUID}'"
	  limit  = 1
	}
	
	data "nutanix_subnets_v2" "subnet" {
	  filter = "name eq '${local.vmm.subnet_name}'"
	}
	
	resource "nutanix_virtual_machine_v2" "ngt-vm" {
	  name                 = "%[2]s"
	  description          = "vm to test ngt installation and insert iso"
	  num_cores_per_socket = 1
	  num_sockets          = 1
	  memory_size_bytes    = 4 * 1024 * 1024 * 1024
	  cluster {
		ext_id = local.clusterUUID
	  }
	
	  disks {
		disk_address {
		  bus_type = "SCSI"
		  index    = 0
		}
		backing_info {
		  vm_disk {
			data_source {
			  reference {
				image_reference {
				  image_ext_id = data.nutanix_images_v2.ngt-image.images[0].ext_id
				}
			  }
			}
			disk_size_bytes = 20 * 1024 * 1024 * 1024
		  }
		}
	  }
	
	  cd_roms {
		disk_address {
		  bus_type = "IDE"
		  index    = 0
		}
	  }

	  cd_roms {
		disk_address {
		  bus_type = "IDE"
		  index    = 1
		}
	  }

	  cd_roms {
		disk_address {
		  bus_type = "SATA"
		  index    = 2
		}
	  }
	
	  nics {
		network_info {
		  nic_type = "NORMAL_NIC"
		  subnet {
			ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
		  }
		  vlan_mode = "ACCESS"
		}
	  }

	    # ---- Guest Customization Section ----
		guest_customization {
			config {
			sysprep {
				install_type = "PREPARED"
				sysprep_script {
				unattend_xml {
					# Base64-encoded unattend.xml file content
					value = base64encode(file("%[3]s"))
				}
				}
			}
			}
		}
	
	  boot_config {
		legacy_boot {
		  boot_order = ["CDROM", "DISK", "NETWORK"]
		}
	  }
	  power_state = "ON"
	
	  lifecycle {
		ignore_changes = [
			cd_roms,
			disks,
			guest_customization,
			guest_tools
		]
	  }
	}
			`, filepath, vmName, untendedWindowsXMLFilePath)
}

func testNGTInstallationConfig() string {
	return `
	resource "nutanix_ngt_installation_v2" "test" {
	  ext_id = nutanix_virtual_machine_v2.ngt-vm.id
	  credential {
		username = local.vmm.ngt.credential.username
		password = local.vmm.ngt.credential.password
	  }
	  reboot_preference {
		schedule_type = "SKIP"
	  }
	  capablities = ["VSS_SNAPSHOT"]
	}
	`
}

func testNGTInsertIsoConfig(configMode, action string) string {
	return fmt.Sprintf(`
	resource "nutanix_ngt_insert_iso_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
		capablities = ["VSS_SNAPSHOT"]
		is_config_only = %s
		action = "%s"
	}

	data "nutanix_virtual_machine_v2" "ngt-vm-refresh" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
		depends_on = [nutanix_ngt_insert_iso_v2.test]
	}	

		`, configMode, action)
}
