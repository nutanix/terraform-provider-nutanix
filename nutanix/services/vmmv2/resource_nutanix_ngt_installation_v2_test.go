package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	resourceNameNGTInstallation = "nutanix_ngt_installation_v2.test"
)

func TestAccV2NutanixNGTInstallationResource_InstallNGTWithRebootPreferenceSetToIMMEDIATE(t *testing.T) {
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
			// get NGT Configuration for the VM after installing NGT
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to reboot")
					time.Sleep(timeSleep)
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTConfiguration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNGTConfiguration, "guest_os_version"),
					resource.TestCheckResourceAttrSet(datasourceNameNGTConfiguration, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_installed", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_reachable", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_enabled", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.#", "2"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccV2NutanixNGTInstallationResource_InstallNGTWithRebootPreferenceSetToLATER(t *testing.T) {
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
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigLATERReboot(),
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
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "LATER"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule.0.start_time", "2026-08-01T00:00:00Z"),
				),
			},
		},
	})
}

func TestAccV2NutanixNGTInstallationResource_InstallNGTWithRebootPreferenceSetToSKIP(t *testing.T) {
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
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigSKIPReboot(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.1", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "SKIP"),
				),
			},
		},
	})
}

func TestAccV2NutanixNGTInstallationResource_WithNoVmExtId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testNGTInstallationResourceWithoutVMExtIDConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixNGTInstallationResource_UpdateNGT(t *testing.T) {
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
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceUpdateConfig(`["SELF_SERVICE_RESTORE","VSS_SNAPSHOT"]`, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.1", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "IMMEDIATE"),
				),
			},
			// test update, change capablities, remove SELF_SERVICE_RESTORE
			{
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceUpdateConfig(`["VSS_SNAPSHOT"]`, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "1"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.0", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "IMMEDIATE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_enabled", "true"),
				),
			},
			// test update, change capabilities, remove VSS_SNAPSHOT
			{
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceUpdateConfig(`[]`, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "0"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "IMMEDIATE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_enabled", "true"),
				),
			},
			// test update, change is_enabled set to false
			{
				SkipFunc: func() (bool, error) {
					t.Skip("Skipping test as it is failing due to issue in the API")
					return true, nil
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceUpdateConfig(`[]`, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "0"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "IMMEDIATE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_enabled", "false"),
				),
			},
		},
	})
}

func testNGTInstallationResourceConfigIMMEDIATEReboot() string {
	return `
	resource "nutanix_ngt_installation_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
		credential {
			username = local.vmm.ngt.credential.username
			password = local.vmm.ngt.credential.password
		}
		reboot_preference {
			schedule_type = "IMMEDIATE"
		}
		capablities = ["SELF_SERVICE_RESTORE","VSS_SNAPSHOT"]
	}`
}

func testNGTInstallationResourceUpdateConfig(capabilities, isEnabled string) string {
	return fmt.Sprintf(`

	resource "nutanix_ngt_installation_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
		credential {
			username = local.vmm.ngt.credential.username
			password = local.vmm.ngt.credential.password
		}
		reboot_preference {
			schedule_type = "IMMEDIATE"
		}
		capablities = %s
		is_enabled = %s
	}
	`, capabilities, isEnabled)
}

func testNGTInstallationResourceConfigSKIPReboot() string {
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
		capablities = ["SELF_SERVICE_RESTORE","VSS_SNAPSHOT"]
	}`
}

func testNGTInstallationResourceConfigLATERReboot() string {
	return `

	resource "nutanix_ngt_installation_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
		credential {
			username = local.vmm.ngt.credential.username
			password = local.vmm.ngt.credential.password
		}
		reboot_preference {
			schedule_type = "LATER"
			schedule {
				start_time = "2026-08-01T00:00:00Z"
			}
		}
		capablities = ["SELF_SERVICE_RESTORE","VSS_SNAPSHOT"]
	}`
}

func testNGTInstallationResourceWithoutVMExtIDConfig() string {
	return `
		resource "nutanix_ngt_installation_v2" "test" {
			credential {
				username = "username"
				password = "password"
			}
			reboot_preference {
				schedule_type = "IMMEDIATE"
			}
			capablities = ["VSS_SNAPSHOT"]
		}`
}

// this config import image, create subnet, create vm
func testPreEnvConfig(vmName string, r int) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		  clusterUUID = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		  config = (jsondecode(file("%[1]s")))
		  vmm    = local.config.vmm
		}

		data "nutanix_images_v2" "ngt-image" {
		  filter = "name eq '${local.vmm.image_name}'"
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.clusterUUID}'"
		  limit = 1
		}

		data "nutanix_subnets_v2" "subnet" {
		 filter = "name eq '${local.vmm.subnet_name}'"
		}

		resource "nutanix_virtual_machine_v2" "ngt-vm" {
		  name                 = "%[3]s"
		  description          = "vm to test ngt installation"
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

		  nics {
			network_info {
			  nic_type = "NORMAL_NIC"
			  subnet {
				ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
			  }
			  vlan_mode = "ACCESS"
			}
		  }

		  boot_config {
			legacy_boot {
			  boot_order = ["CDROM", "DISK", "NETWORK"]
			}
		  }
		  power_state = "ON"

		  lifecycle {
			ignore_changes = [guest_tools]
		  }

		  depends_on = [data.nutanix_clusters_v2.clusters, data.nutanix_images_v2.ngt-image, data.nutanix_storage_containers_v2.ngt-sc]
		}

`, filepath, r, vmName)
}

var testNGTConfiguration = `
	data "nutanix_ngt_configuration_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
	}`
