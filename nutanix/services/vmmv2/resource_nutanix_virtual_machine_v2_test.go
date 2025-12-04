package vmmv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVms = "nutanix_virtual_machine_v2.test"

func TestAccV2NutanixVmsResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_BasicUpdate(t *testing.T) {
	r := acctest.RandInt()
	desc := "test vm description"
	updatedDesc := "test vm updated description"
	name := fmt.Sprintf("tf-test-vm-%d", r)
	updatedName := fmt.Sprintf("tf-test-vm-updated-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
				),
			},
			{
				Config: testVmsV4Config(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithDisk(t *testing.T) {
	r := acctest.RandInt()
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithDisk(r, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "disks.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "project.#"),
					resource.TestCheckResourceAttrPair(resourceNameVms, "project.0.ext_id", "nutanix_project.projects", "metadata.uuid"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_DiskWithDatasource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithDiskWithDatasource(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("%s-new", name)),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "disks.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.backing_info.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.backing_info.0.vm_disk.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithNic(t *testing.T) {
	r := acctest.RandInt()
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithNic(r, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "disks.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "nics.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithNicTrunk(t *testing.T) {
	r := acctest.RandInt()
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithNicWithTrunkVlan(r, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "disks.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "nics.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "TRUNK"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.trunked_vlans.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithLegacyBootOrder(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithLegacyBoot(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.#"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.0.legacy_boot.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.0", "CDROM"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.1", "DISK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.2", "NETWORK"),
				),
			},
			{
				Config: testVmsV4ConfigWithLegacyBootWithUpdateOrder(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.#"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.0.legacy_boot.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.0", "DISK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.1", "CDROM"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.2", "NETWORK"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithUEFIBootOrder(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	nameUpdated := fmt.Sprintf("tf-test-vm-updated-%d", r)
	descUpdated := "test vm updated description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithUEFIBoot(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.#"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.0.uefi_boot.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.uefi_boot.0.boot_order.0", "CDROM"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.uefi_boot.0.boot_order.1", "DISK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.uefi_boot.0.boot_order.2", "NETWORK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.uefi_boot.0.is_secure_boot_enabled", "false"),
				),
			},
			{
				Config: testVmsV4ConfigWithUEFIBootUpdate(nameUpdated, descUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", descUpdated),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.#"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.0.uefi_boot.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.uefi_boot.0.boot_order.0", "NETWORK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.uefi_boot.0.boot_order.1", "DISK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.uefi_boot.0.boot_order.2", "CDROM"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.uefi_boot.0.is_secure_boot_enabled", "false"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithCdrom(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithCdrom(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "cd_roms.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "cd_roms.0.disk_address.0.bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameVms, "cd_roms.0.disk_address.0.index", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithCdromIDE(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithCdromIde(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "cd_roms.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "cd_roms.0.disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(resourceNameVms, "cd_roms.0.disk_address.0.index", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithCdromBackingInfo(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithCdromWithBackingInfo(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "cd_roms.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "cd_roms.0.disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(resourceNameVms, "cd_roms.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "cd_roms.0.backing_info.0.data_source.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithCloudInit(t *testing.T) {
	r := acctest.RandInt()
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithCloudInit(r, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "disks.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "nics.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithSysprep(t *testing.T) {
	r := acctest.RandInt()
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithSysprep(r, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "disks.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "nics.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "cd_roms.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "cd_roms.0.iso_type", "GUEST_CUSTOMIZATION"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithCloudInitWithCustomKeys(t *testing.T) {
	t.Skip("Failing due to issue with cloud_init user_data & custom_keys")
	r := acctest.RandInt()
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithCloudInitWithCustomKeys(r, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "disks.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "nics.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_UpdateDiskNics(t *testing.T) {
	r := acctest.RandInt()
	desc := "test vm description"
	updatedDesc := "test vm updated description"
	name := fmt.Sprintf("tf-test-vm-%d", r)
	updatedName := fmt.Sprintf("tf-test-vm-updated-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithDiskNic(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			{
				Config: testVmsV4ConfigWitUpdatedDiskNic(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "2"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.1.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "disks.1.disk_address.0.index", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithLegacyBootDevice(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithLegacyBootDevice(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.#"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "boot_config.0.legacy_boot.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_device.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_device.0.boot_device_disk.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_device.0.boot_device_disk.0.disk_address.0.index", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithCategories(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCategoriesV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "2"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "categories.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "categories.#", "3"),
				),
			},
			{
				Config: testVmsCategoriesV4ConfigUpdate(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "2"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "categories.#"),
					resource.TestCheckResourceAttr(resourceNameVms, "categories.#", "2"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithSerialPorts(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsConfigWithSerialPorts(name, desc, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "2"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttr(resourceNameVms, "serial_ports.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "serial_ports.0.index", "2"),
					resource.TestCheckResourceAttr(resourceNameVms, "serial_ports.0.is_connected", "true"),
				),
			},
			{
				Config: testVmsConfigWithSerialPorts(name, desc, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "2"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttr(resourceNameVms, "serial_ports.0.index", "2"),
					resource.TestCheckResourceAttr(resourceNameVms, "serial_ports.0.is_connected", "false"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_WithGpus(t *testing.T) {
	if testVars.VMM.GPUS[0].Vendor == "" && testVars.VMM.GPUS[0].Mode == "" && testVars.VMM.GPUS[0].DeviceID == 0 {
		t.Skip("Skipping test as no GPU devices found")
	}
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	updatedDesc := "test vm updated description"
	updatedName := fmt.Sprintf("tf-test-vm-updated-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV2WithGpus(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "2"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttr(resourceNameVms, "gpus.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "gpus.0.device_id", strconv.Itoa(testVars.VMM.GPUS[0].DeviceID)),
					resource.TestCheckResourceAttr(resourceNameVms, "gpus.0.mode", testVars.VMM.GPUS[0].Mode),
					resource.TestCheckResourceAttr(resourceNameVms, "gpus.0.vendor", testVars.VMM.GPUS[0].Vendor),
					resource.TestCheckResourceAttrSet(resourceNameVms, "gpus.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "gpus.0.name"),
				),
			},
			{
				Config: testVmsV2RemoveGpus(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "2"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttr(resourceNameVms, "gpus.#", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsResource_ClusterAutomaticSelection(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsClusterAutomaticSelectionConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVms, "protection_type", "UNPROTECTED"),
					resource.TestCheckResourceAttr(resourceNameVms, "is_agent_vm", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "machine_type", "PC"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "cluster.0.ext_id"),

					// list vms checks
					resource.TestCheckResourceAttrSet(datasourceNameVM, "vms.#"),
					resource.TestCheckResourceAttr(datasourceNameVM, "vms.0.name", name),
					resource.TestCheckResourceAttrSet(datasourceNameVM, "vms.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVM, "vms.0.cluster.0.ext_id"),

					// get vm checks
					resource.TestCheckResourceAttr(datasourceNameVMs, "name", name),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVMs, "cluster.0.ext_id"),
				),
			},
		},
	})
}

func testVmsV4Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}
`, name, desc)
}

func testVmsV4ConfigWithDisk(r int, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		data "nutanix_subnets_v2" "subnets" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			subnetExtId = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		resource "nutanix_project" "projects" {
			name          = "tf-project-%[1]d"
			description   = "project twice"
			use_project_internal = true
			api_version = "3.1"
			cluster_reference_list {
				uuid = local.cluster0
			}

			subnet_reference_list {
				uuid = local.subnetExtId
			}

			default_subnet_reference {
				uuid = local.subnetExtId
			}
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "tf-test-vm-%[1]d"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			power_state = "OFF"
			cluster {
				ext_id = local.cluster0
			}
			project {
				ext_id = nutanix_project.projects.metadata.uuid
      		}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = local.subnetExtId
					}
					vlan_mode = "ACCESS"
				}
			}
		}
`, r, desc, filepath)
}

func testVmsV4ConfigWithDiskWithDatasource(name string, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "testWithDisk"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			power_state = "OFF"
		}


		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s-new"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = 1073741824
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
						data_source{
							reference{
								vm_disk_reference {
									disk_address{
										bus_type="SCSI"
										index = 0
									}
									vm_reference{
										ext_id = resource.nutanix_virtual_machine_v2.testWithDisk.id
									}
								}
							}
						}
					}
				}
			}
			power_state = "ON"
			depends_on = [ resource.nutanix_virtual_machine_v2.testWithDisk ]
			lifecycle{
				ignore_changes = [
					disks.0.backing_info.0.vm_disk.0.data_source
				]
			}
		}
`, name, desc, filepath)
}

func testVmsV4ConfigWithNic(r int, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "tf-test-vm-%[1]d"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			power_state = "ON"
		}
`, r, desc, filepath)
}

func testVmsV4ConfigWithNicWithTrunkVlan(r int, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "tf-test-vm-%[1]d"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "TRUNK"
					trunked_vlans = ["1"]
				}
			}
			power_state = "ON"
		}
`, r, desc, filepath)
}

func testVmsV4ConfigWithLegacyBoot(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			boot_config{
				legacy_boot{
				  boot_order = ["CDROM", "DISK","NETWORK"]
				}
			}
			power_state = "ON"
		}
`, name, desc)
}

func testVmsV4ConfigWithUEFIBoot(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name                 = "%[1]s"
			description          = "%[2]s"
			num_cores_per_socket = 1
			num_sockets          = 1
			cluster {
				ext_id = local.cluster0
			}


			boot_config {
				uefi_boot {
				}
			}

			power_state = "OFF"
		}
`, name, desc)
}

func testVmsV4ConfigWithUEFIBootUpdate(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name                 = "%[1]s"
			description          = "%[2]s"
			num_cores_per_socket = 1
			num_sockets          = 1
			cluster {
				ext_id = local.cluster0
			}


			boot_config {
				uefi_boot {
					boot_order = ["NETWORK","DISK", "CDROM", ]
				}
			}

			power_state = "OFF"
		}
`, name, desc)
}

func testVmsV4ConfigWithLegacyBootWithUpdateOrder(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			boot_config{
				legacy_boot{
				  boot_order = ["DISK", "CDROM", "NETWORK"]
				}
			}
			power_state = "ON"
		}
`, name, desc)
}

func testVmsV4ConfigWithLegacyBootDevice(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}



		data "nutanix_image" "ngt-image" {
		  image_name = local.vmm.image_name
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						data_source {
							reference {
								image_reference{
									image_ext_id = data.nutanix_image.ngt-image.id
								}
							}
						}
					}
				}
			}
			boot_config{
				legacy_boot{
					boot_device{
						boot_device_disk {
							disk_address {
								bus_type = "SCSI"
								index = 0
							}
						}
				  	}
				}
			}
			power_state = "ON"
			depends_on = [data.nutanix_clusters_v2.clusters, data.nutanix_image.ngt-image]
		}
`, name, desc, filepath)
}

func testVmsV4ConfigWithCdrom(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			boot_config{
				legacy_boot{
				  boot_order = ["CDROM", "DISK","NETWORK"]
				}
			}
			cd_roms{
				disk_address{
					bus_type = "SATA"
					index= 0
				}
			}
		}
`, name, desc)
}

func testVmsV4ConfigWithCdromIde(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			boot_config{
				legacy_boot{
				  boot_order = ["CDROM", "DISK","NETWORK"]
				}
			}
			cd_roms{
				disk_address{
					bus_type = "IDE"
					index= 0
				}
			}
			power_state = "ON"
		}
`, name, desc)
}

func testVmsV4ConfigWithCdromWithBackingInfo(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "testWithDisk"{
			name= "%[1]s-second"
			description =  "%[2]s-second"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			power_state = "OFF"
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			boot_config{
				legacy_boot{
				  boot_order = ["CDROM", "DISK","NETWORK"]
				}
			}
			cd_roms{
				disk_address{
					bus_type = "IDE"
					index= 0
				}
				backing_info{
					storage_container {
						ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
					}
					data_source {
						reference{
							vm_disk_reference {
								disk_address{
									bus_type="SCSI"
									index = 0
								}
								vm_reference{
									ext_id= resource.nutanix_virtual_machine_v2.testWithDisk.id
								}
							}
						}
					}
				}
			}
			power_state = "ON"
			lifecycle{
				ignore_changes = [
					cd_roms.0.backing_info.0.data_source
				]
			}
		}
`, name, desc, filepath)
}

func testVmsV4ConfigWithCloudInit(r int, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			gs = base64encode("#cloud-config\nusers:\n  - name: ubuntu\n    ssh-authorized-keys:\n      - ssh-rsa DUMMYSSH mypass\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']")
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "tf-test-vm-%[1]d"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			guest_customization{
				config{
					cloud_init{
						cloud_init_script{
							user_data{
								value="${local.gs}"
							}
						}
					}
				}
			}

			lifecycle{
				ignore_changes = [
					guest_customization, cd_roms
				]
			}
		}
`, r, desc, filepath)
}

func testVmsV4ConfigWithSysprep(r int, desc string) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "clusters" {}

locals {
	cluster0 = [
		for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	][0]
	config = jsondecode(file("%[3]s"))
	vmm = local.config.vmm
}

data "nutanix_storage_containers_v2" "ngt-sc" {
	filter = "clusterExtId eq '${local.cluster0}'"
	limit = 1
}

data "nutanix_subnets_v2" "subnets" {
	filter = "name eq '${local.vmm.subnet_name}'"
}

resource "nutanix_virtual_machine_v2" "test"{
	name= "tf-test-vm-%[1]d"
	description =  "%[2]s"
	num_cores_per_socket = 1
	num_sockets = 1
	cluster {
		ext_id = local.cluster0
	}
	disks{
		disk_address{
			bus_type = "SCSI"
			index = 0
		}
		backing_info{
			vm_disk{
				disk_size_bytes = "1073741824"
				storage_container{
					ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
				}
			}
		}
	}
	nics{
		network_info{
			nic_type = "NORMAL_NIC"
			subnet{
				ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
			}
			vlan_mode = "ACCESS"
		}
	}
	guest_customization {
		config {
		sysprep {
			install_type = "PREPARED"
				sysprep_script {
					unattend_xml {
						value = file("%[4]s") # unattend_xml file value, value is encoded in base64
					}
			}
		}
		}
	}

	lifecycle{
		ignore_changes = [
			guest_customization, cd_roms
		]
	}
}
`, r, desc, filepath, untendedXMLFilePath)
}

func testVmsV4ConfigWithCloudInitWithCustomKeys(r int, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			gs = base64encode("#cloud-config\nusers:\n  - name: ubuntu\n    ssh-authorized-keys:\n      - ssh-rsa DUMMYSSH mypass\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']")
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "tf-test-vm-%[1]d"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			guest_customization {
				config {
					cloud_init {
						cloud_init_script {
							custom_key_values {
								key_value_pairs {
									name = "locale"
									value {
										string = "en-US"
									}
								}
							}
						}
					}
				}
			}
			lifecycle{
			ignore_changes = [
				guest_customization, cd_roms
			]
			}
		}
`, r, desc, filepath)
}

func testVmsV4ConfigWithDiskNic(name string, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = 1073741824
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			power_state = "ON"
		}
	`, name, desc, filepath)
}

func testVmsV4ConfigWitUpdatedDiskNic(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 2
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 1
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			power_state = "ON"
		}
	`, name, desc, filepath)
}

func testVmsCategoriesV4Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_categories_v2" "ctg"{}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 2
			cluster {
				ext_id = local.cluster0
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			categories{
				ext_id = data.nutanix_categories_v2.ctg.categories.2.ext_id
			}
			categories{
				ext_id = data.nutanix_categories_v2.ctg.categories.5.ext_id
			}
			categories{
				ext_id = data.nutanix_categories_v2.ctg.categories.3.ext_id
			}
			power_state = "ON"
		}
	`, name, desc, filepath)
}

func testVmsCategoriesV4ConfigUpdate(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_categories_v2" "ctg"{}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 2
			cluster {
				ext_id = local.cluster0
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			categories{
				ext_id = data.nutanix_categories_v2.ctg.categories.2.ext_id
			}
			categories{
				ext_id = data.nutanix_categories_v2.ctg.categories.5.ext_id
			}
			power_state = "ON"
		}
	`, name, desc, filepath)
}

func testVmsV2WithGpus(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		  cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		  config = jsondecode(file("%[3]s"))
		  vmm    = local.config.vmm
		}


		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 2
			cluster {
				ext_id = local.cluster0
			}

			gpus {
				device_id = local.vmm.gpus[0].device_id
				mode      = local.vmm.gpus[0].mode
				vendor    = local.vmm.gpus[0].vendor
			}
			power_state = "ON"
		}
	`, name, desc, filepath)
}

func testVmsV2RemoveGpus(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		  cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		  ][0]
		  config = jsondecode(file("%[3]s"))
		  vmm    = local.config.vmm
		}


		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 2
			cluster {
				ext_id = local.cluster0
			}

			power_state = "OFF"
		}
	`, name, desc, filepath)
}

func testVmsConfigWithSerialPorts(name, desc string, isconn bool) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = jsondecode(file("%[4]s"))
			vmm = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_categories_v2" "ctg"{}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 2
			cluster {
				ext_id = local.cluster0
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				}
			}
			serial_ports{
				index = 2
				is_connected = %[3]t
			}
			power_state = "ON"
		}
	`, name, desc, isconn, filepath)
}

func testVmsClusterAutomaticSelectionConfig(name, desc string) string {
	return fmt.Sprintf(`

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
		}

		data "nutanix_virtual_machine_v2" "test" {
			ext_id = nutanix_virtual_machine_v2.test.id
		}

		data "nutanix_virtual_machines_v2" "test" {
			filter = "name eq '${nutanix_virtual_machine_v2.test.name}'"
		}
`, name, desc)
}
