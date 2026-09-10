package vmmv2_test

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					resource.TestCheckResourceAttrPair(resourceNameVms, "project.0.ext_id", "nutanix_project_v2.share_project", "ext_id"),
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
				Config: testVmsV4ConfigWithNic(r, desc, true),
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
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			{
				Config: testVmsV4ConfigWithNic(r, desc, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVms, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.model", "VIRTIO"),
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
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "TRUNK"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.trunked_vlans.#", "1"),
				),
			},
		},
	})
}

// TestAccV2NutanixVmsResource_NicAddRemove tests:
// 1. Create VM with one NIC
// 2. Update VM to add a second NIC
// 3. Remove the added NIC (back to one NIC)
// This test covers issue #994 - VM creation/update with multiple NICs
func TestAccV2NutanixVmsResource_NicAddRemove(t *testing.T) {
	r := acctest.RandIntRange(1, 1000)
	name := fmt.Sprintf("tf-test-vm-nic-%d", r)
	desc := "test vm for NIC add/remove"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmsResourceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create VM with one NIC
			{
				Config: testVmsV4ConfigWithSingleNic(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 2: Add a second NIC
			{
				Config: testVmsV4ConfigWithTwoNics(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.1.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.1.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 3: Remove the second NIC (back to one NIC)
			// Keep subnet resource to avoid deletion error while IP assignment is being released
			{
				Config: testVmsV4ConfigWithSingleNicKeepSubnet(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
		},
	})
}

// TestAccV2NutanixVmsResource_NicScenariosVlanModeAndIsConnected covers:
// 1. Create VM with new fields (nic_backing_info, nic_network_info) and legacy (backing_info, network_info) same values -> plan no changes
// 2. Update VM with new fields (vlan_mode, is_connected) -> plan no changes
// 3. Update VM with old fields (vlan_mode, is_connected) -> plan no changes
// 4. Update VM with new and old fields having different values -> new fields win, state updated for both blocks
// 5. Update VM with new and old fields with same values -> plan no changes
// 6. Destroy VM (CheckDestroy)
func TestAccV2NutanixVmsResource_NicScenariosVlanModeAndIsConnected(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-nic-scenarios-%d", r)
	desc := "test vm for NIC vlan_mode and is_connected scenarios"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmsResourceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create VM with both new and legacy NIC blocks, same values (is_connected=false, vlan_mode=TRUNK)
			{
				Config: testVmsV4ConfigNicScenariosStep1Create(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.model", "VIRTIO"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "TRUNK"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "TRUNK"),
				),
			},
			// Idempotent: same config, no change
			{
				Config: testVmsV4ConfigNicScenariosStep1Create(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "TRUNK"),
				),
			},
			// Step 3: Update VM with new fields only (is_connected=true, vlan_mode=ACCESS)
			{
				Config: testVmsV4ConfigNicScenariosStep2UpdateNewFields(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Idempotent after update with new fields
			{
				Config: testVmsV4ConfigNicScenariosStep2UpdateNewFields(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 5: Update VM with old fields only (same values is_connected=true, vlan_mode=ACCESS) -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosStep3UpdateOldFields(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Idempotent after update with old fields
			{
				Config: testVmsV4ConfigNicScenariosStep3UpdateOldFields(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 7: Update with new and old fields having different values -> new fields win
			{
				Config:             testVmsV4ConfigNicScenariosStep4DifferentValues(r, name, desc),
				ExpectNonEmptyPlan: true, // since the new fields are updated and the state is updated with new values for both blocks (expected behavior)
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 8: Update with new and old fields same values -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosStep5SameValues(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Idempotent: same config
			{
				Config: testVmsV4ConfigNicScenariosStep5SameValues(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
		},
	})
}

// TestAccV2NutanixVmsResource_NicScenariosCreateWithSameValuesThenUpdates covers:
// 13. Create the vm with new fields and old fields with same values -> new fields values -> terraform plan no changes
// 14. Update the vm with new fields -> terraform plan no changes
// 15. Update the vm with old fields -> terraform plan no changes
// 16. Update the vm with new fields and old fields with different values -> taking new fields values -> state updated for both blocks (expected behavior)
// 17. Update the vm with new fields and old fields with same values -> new fields values -> terraform plan no changes
// 18. Delete the vm -> terraform destroy (CheckDestroy)
func TestAccV2NutanixVmsResource_NicScenariosCreateWithSameValuesThenUpdates(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-nic-same-then-updates-%d", r)
	desc := "test vm for NIC scenarios 13-18 (create with same values then updates)"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmsResourceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create VM with new fields and old fields with same values -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosStep1Create(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "TRUNK"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "TRUNK"),
				),
			},
			{
				Config: testVmsV4ConfigNicScenariosStep1Create(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "TRUNK"),
				),
			},
			// Step 3: Update the vm with new fields -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosStep2UpdateNewFields(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			{
				Config: testVmsV4ConfigNicScenariosStep2UpdateNewFields(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 5: Update the vm with old fields -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosStep3UpdateOldFields(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			{
				Config: testVmsV4ConfigNicScenariosStep3UpdateOldFields(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// 7: Update with new and old different values -> new fields win, state updated for both blocks
			{
				Config:             testVmsV4ConfigNicScenariosStep4DifferentValues(r, name, desc),
				ExpectNonEmptyPlan: true, // since the new fields are updated and the state is updated with new values for both blocks (expected behavior)
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 8: Update with new and old same values -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosStep5SameValues(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			{
				Config: testVmsV4ConfigNicScenariosStep5SameValues(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
		},
	})
}

// TestAccV2NutanixVmsResource_NicScenariosCreateWithDifferentValuesThenUpdates covers (main.tf scenarios 19-23):
// 19. Create the vm with new fields and old fields with different values -> taking new fields values -> state updated with new values for both blocks (expected behavior)
// 20. Update the vm with new fields -> terraform plan no changes
// 21. Update the vm with old fields -> terraform plan no changes
// 22. Update the vm with new fields and old fields with same values -> new fields values -> terraform plan no changes
// 23. Update the vm with new fields and old fields with different values -> taking new fields values -> state updated (expected behavior)
func TestAccV2NutanixVmsResource_NicScenariosCreateWithDifferentValuesThenUpdates(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-nic-diff-then-updates-%d", r)
	desc := "test vm for NIC scenarios 19-23 (create with different values then updates)"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmsResourceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create VM with new and old fields having different values -> new wins, state updated for both blocks
			{
				Config:             testVmsV4ConfigNicScenariosStep4DifferentValues(r, name, desc),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 2: Update the vm with new fields -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosStep1Create(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "TRUNK"),
				),
			},
			{
				Config: testVmsV4ConfigNicScenariosStep1Create(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "TRUNK"),
				),
			},
			// Step 4: Update the vm with old fields -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosLegacyOnly(r, name, desc, false, "TRUNK"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "TRUNK"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "TRUNK"),
				),
			},
			{
				Config: testVmsV4ConfigNicScenariosLegacyOnly(r, name, desc, false, "TRUNK"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "TRUNK"),
				),
			},
			// Step 6: Update with new and old same values -> plan no changes
			{
				Config: testVmsV4ConfigNicScenariosStep5SameValues(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
				),
			},
			{
				Config: testVmsV4ConfigNicScenariosStep5SameValues(r, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
				),
			},
			// Step 8: Update with new and old different values -> new wins, state updated
			{
				Config:             testVmsV4ConfigNicScenariosStep4DifferentValues(r, name, desc),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_backing_info.0.virtual_ethernet_nic.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.backing_info.0.is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.network_info.0.vlan_mode", "ACCESS"),
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

func TestAccV2NutanixVmsResource_LegacyBootOrderChangeWithBootDevice(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	mac := fmt.Sprintf("50:6b:8d:%02x:%02x:%02x", r&0xff, (r>>8)&0xff, (r>>16)&0xff)
	bootDeviceNic := fmt.Sprintf(`
				  boot_device {
					boot_device_nic {
					  mac_address = "%s"
					}
				  }`, mac)
	bootDeviceDisk := `
				  boot_device {
					boot_device_disk {
					  disk_address {
						bus_type = "SATA"
						index    = 0
					  }
					}
				  }`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV4ConfigWithLegacyBootOrderChange(name, desc, mac, "", `["CDROM", "DISK", "NETWORK"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", name),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.0", "CDROM"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.1", "DISK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.2", "NETWORK"),
				),
			},
			{
				Config: testVmsV4ConfigWithLegacyBootOrderChange(name, desc, mac, bootDeviceNic, `["NETWORK", "DISK", "CDROM"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_device.0.boot_device_nic.0.mac_address", mac),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.0", "NETWORK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.1", "DISK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.2", "CDROM"),
				),
			},
			{
				Config: testVmsV4ConfigWithLegacyBootOrderChange(name, desc, mac, bootDeviceDisk, `["CDROM", "DISK", "NETWORK"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_device.0.boot_device_disk.0.disk_address.0.bus_type", "SATA"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_device.0.boot_device_disk.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_device.0.boot_device_nic.#", "0"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.0", "CDROM"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.1", "DISK"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.2", "NETWORK"),
				),
			},
			{
				Config: testVmsV4ConfigWithLegacyBootOrderChange(name, desc, mac, "", `["CDROM", "DISK", "NETWORK"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_device.#", "0"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.0", "CDROM"),
					resource.TestCheckResourceAttr(resourceNameVms, "boot_config.0.legacy_boot.0.boot_order.1", "DISK"),
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
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
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
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
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
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
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
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
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
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(resourceNameVms, "nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.vlan_mode", "ACCESS"),
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

// TestAccV2NutanixVmsResource_WaitForRoutableIP covers issue #871: on create the resource
// must not report the APIPA/link-local address (169.254.0.0/16) that guest tools surface
// transiently before DHCP completes.
//
// Step 1 uses the defaults (wait_for_ip_timeout = 5, wait_for_ip_routable = true) and
// asserts that any IPv4 address learned on the NIC is routable. Step 2 sets
// wait_for_ip_timeout = 0 to assert the wait can be disabled without failing the apply.
//
// Note: this is a timing-dependent guest behaviour, so the check only fails when an APIPA
// address is actually reported - it does not require the guest to obtain a lease (a VM
// with no guest tools simply reports no learned address and waits out the timeout, which
// is non-fatal by design).
func TestAccV2NutanixVmsResource_WaitForRoutableIP(t *testing.T) {
	r := acctest.RandInt()
	desc := "test vm waiting for a routable ip"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsV2ConfigWaitForIP(r, desc, 5, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "name", fmt.Sprintf("tf-test-vm-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVms, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceNameVms, "wait_for_ip_timeout", "5"),
					resource.TestCheckResourceAttr(resourceNameVms, "wait_for_ip_routable", "true"),
					resource.TestCheckResourceAttrSet(resourceNameVms, "nics.#"),
					testCheckLearnedIPsAreRoutable(resourceNameVms),
				),
			},
			{
				// wait_for_ip_timeout < 1 disables the wait entirely.
				Config: testVmsV2ConfigWaitForIP(r, desc, 0, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVms, "wait_for_ip_timeout", "0"),
					resource.TestCheckResourceAttr(resourceNameVms, "power_state", "ON"),
				),
			},
		},
	})
}

// testCheckLearnedIPsAreRoutable fails if any IPv4 address learned on the VM's NICs is an
// APIPA/link-local address, which is what issue #871 reported being handed to provisioners.
func testCheckLearnedIPsAreRoutable(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		for key, value := range rs.Primary.Attributes {
			if !strings.HasPrefix(key, "nics.") || !strings.HasSuffix(key, ".value") ||
				!strings.Contains(key, ".ipv4_info.") || !strings.Contains(key, ".learned_ip_addresses.") {
				continue
			}
			if ip := net.ParseIP(value); ip != nil && ip.To4() != nil && ip.IsLinkLocalUnicast() {
				return fmt.Errorf("%s = %q is an APIPA/link-local address; expected a routable IPv4 address (issue #871)", key, value)
			}
		}
		return nil
	}
}

func testVmsV2ConfigWaitForIP(r int, desc string, waitTimeout int, routable bool) string {
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

		data "nutanix_storage_containers_v2" "sc" {
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
							ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics{
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
				}
			}
			power_state = "ON"

			wait_for_ip_timeout  = %[4]d
			wait_for_ip_routable = %[5]t
		}
`, r, desc, filepath, waitTimeout, routable)
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
		data "nutanix_clusters_v2" "clusters" {
			filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"		
		}

		resource "nutanix_subnet_v2" "subnet" {
		  name                 = "vlan.888"
		  description          = "Default VLAN 888 Subnet"
		  subnet_type          = "VLAN"
		  cluster_reference    = local.cluster0
		  network_id           = 888
		  shared_with_projects = [nutanix_project_v2.share_project.ext_id]
		  # Destroy the subnet before the resource group: a subnet on the RG's
		  # cluster counts as an entity, so deleting the RG first fails with
		  # clustermgmt:10018 "resource group is not empty".
		  depends_on = [nutanix_resource_group_v2.test]
		}

		locals {
			cluster0 = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
			subnetExtId = nutanix_subnet_v2.subnet.ext_id
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		resource "nutanix_project_v2" "share_project" {
			name        = "tf-project-%[1]d"
			project_id = "tf-project-%[1]d"
			description = "project twice"
		}

		resource "nutanix_resource_group_v2" "test" {
			name           = "tf-rg-%[1]d"
			project_ext_id = nutanix_project_v2.share_project.ext_id
			placement_targets {
				cluster_ext_id = local.cluster0
				storage_containers {
					ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
				}
			}
			# The backend auto-associates a default storage container with the
			# placement target (cluster) and returns it on read, even though we only
			# set cluster_ext_id. That server-populated value would otherwise surface
			# as a perpetual "remove storage_containers" diff after apply and fail the
			# test, so ignore post-create changes to placement_targets here.
			lifecycle {
				ignore_changes = [placement_targets]
			}
		}
		
		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				ext_id = nutanix_project_v2.share_project.ext_id
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = local.subnetExtId
						}
						vlan_mode = "ACCESS"
					}
				}
			}
			# Destroy the VM (whose disk lives on the RG placement-target storage
			# container) before the resource group, otherwise the RG delete fails
			# with clustermgmt:10018 "resource group is not empty".
			depends_on = [nutanix_resource_group_v2.test]
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
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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

func testVmsV4ConfigWithNic(r int, desc string, isConnected bool) string {
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
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
				}
				nic_backing_info{
					virtual_ethernet_nic{
						is_connected = %[4]t
						model = "VIRTIO"
					}
				}
			}
			power_state = "ON"
		}
`, r, desc, filepath, isConnected)
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
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "TRUNK"
						trunked_vlans = ["1"]
					}
				}
			}
			power_state = "ON"
		}
`, r, desc, filepath)
}

// testVmsV4ConfigNicScenariosStep1Create: VM with both nic_* and legacy blocks, same values (is_connected=false, vlan_mode=TRUNK).
func testVmsV4ConfigNicScenariosStep1Create(r int, name, desc string) string {
	return testVmsV4ConfigNicScenariosWithParams(r, name, desc, false, "TRUNK", false, "TRUNK")
}

// testVmsV4ConfigNicScenariosStep2UpdateNewFields: update using new fields (is_connected=true, vlan_mode=ACCESS); both blocks same.
func testVmsV4ConfigNicScenariosStep2UpdateNewFields(r int, name, desc string) string {
	return testVmsV4ConfigNicScenariosWithParams(r, name, desc, true, "ACCESS", true, "ACCESS")
}

// testVmsV4ConfigNicScenariosStep3UpdateOldFields: config with only legacy blocks (is_connected=true, vlan_mode=ACCESS).
func testVmsV4ConfigNicScenariosStep3UpdateOldFields(r int, name, desc string) string {
	return testVmsV4ConfigNicScenariosLegacyOnly(r, name, desc, true, "ACCESS")
}

// testVmsV4ConfigNicScenariosStep4DifferentValues: new fields is_connected=true, vlan_mode=ACCESS; legacy is_connected=false, vlan_mode=TRUNK. New wins.
func testVmsV4ConfigNicScenariosStep4DifferentValues(r int, name, desc string) string {
	return testVmsV4ConfigNicScenariosWithParams(r, name, desc, true, "ACCESS", false, "TRUNK")
}

// testVmsV4ConfigNicScenariosStep5SameValues: both blocks same (is_connected=true, vlan_mode=ACCESS).
func testVmsV4ConfigNicScenariosStep5SameValues(r int, name, desc string) string {
	return testVmsV4ConfigNicScenariosWithParams(r, name, desc, true, "ACCESS", true, "ACCESS")
}

// testVmsV4ConfigNicScenariosBase returns shared VM config (data sources, cluster, disks) with the given nics block injected.
func testVmsV4ConfigNicScenariosBase(name, desc, nicsBlock string) string {
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
			filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
			limit  = 1
		}

		resource "nutanix_virtual_machine_v2" "test" {
			name                 = "%[1]s"
			description          = "%[2]s"
			num_cores_per_socket  = 1
			num_sockets           = 1
			cluster {
				ext_id = local.cluster0
			}
			disks {
				disk_address {
					bus_type = "SCSI"
					index    = 0
				}
				backing_info {
					vm_disk {
						disk_size_bytes = "1073741824"
						storage_container {
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			%[4]s
			power_state = "OFF"
		}
`, name, desc, filepath, nicsBlock)
}

// testVmsV4ConfigNicScenariosNicsBlock returns the nics block for both nic_* and legacy (same structure as main.tf).
func testVmsV4ConfigNicScenariosNicsBlock(newIsConnected bool, newVlanMode string, legacyIsConnected bool, legacyVlanMode string) string {
	return fmt.Sprintf(`nics {
				nic_backing_info {
					virtual_ethernet_nic {
						model        = "VIRTIO"
						is_connected = %[1]t
					}
				}
				nic_network_info {
					virtual_ethernet_nic_network_info {
						nic_type = "NORMAL_NIC"
						subnet {
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "%[2]s"
					}
				}
				backing_info {
					model        = "VIRTIO"
					is_connected = %[3]t
				}
				network_info {
					nic_type = "NORMAL_NIC"
					subnet {
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "%[4]s"
				}
			}`, newIsConnected, newVlanMode, legacyIsConnected, legacyVlanMode)
}

// testVmsV4ConfigNicScenariosNicsBlockLegacyOnly returns the nics block with only legacy backing_info and network_info.
func testVmsV4ConfigNicScenariosNicsBlockLegacyOnly(isConnected bool, vlanMode string) string {
	return fmt.Sprintf(`nics {
				backing_info {
					model        = "VIRTIO"
					is_connected = %[1]t
				}
				network_info {
					nic_type = "NORMAL_NIC"
					subnet {
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}
					vlan_mode = "%[2]s"
				}
			}`, isConnected, vlanMode)
}

// testVmsV4ConfigNicScenariosWithParams: VM with both nic_* and legacy blocks; newIsConnected/newVlanMode for nic_*, legacyIsConnected/legacyVlanMode for legacy.
func testVmsV4ConfigNicScenariosWithParams(r int, name, desc string, newIsConnected bool, newVlanMode string, legacyIsConnected bool, legacyVlanMode string) string {
	return testVmsV4ConfigNicScenariosBase(name, desc, testVmsV4ConfigNicScenariosNicsBlock(newIsConnected, newVlanMode, legacyIsConnected, legacyVlanMode))
}

// testVmsV4ConfigNicScenariosLegacyOnly: VM with only legacy backing_info and network_info (for "update with old fields" scenario).
func testVmsV4ConfigNicScenariosLegacyOnly(r int, name, desc string, isConnected bool, vlanMode string) string {
	return testVmsV4ConfigNicScenariosBase(name, desc, testVmsV4ConfigNicScenariosNicsBlockLegacyOnly(isConnected, vlanMode))
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

func testVmsV4ConfigWithLegacyBootOrderChange(name, desc, macAddress, bootDeviceBlock, bootOrder string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = jsondecode(file("%[6]s"))
			vmm = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name                 = "%[1]s"
			description          = "%[2]s"
			num_cores_per_socket = 1
			num_sockets          = 1
			power_state          = "OFF"
			cluster {
				ext_id = local.cluster0
			}
			nics {
				nic_backing_info {
					virtual_ethernet_nic {
						mac_address = "%[3]s"
					}
				}
				nic_network_info {
					virtual_ethernet_nic_network_info {
						nic_type = "NORMAL_NIC"
						subnet {
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
				}
			}
			cd_roms {
				disk_address {
					bus_type = "SATA"
					index    = 0
				}
			}
			boot_config {
				legacy_boot {
%[4]s
					boot_order = %[5]s
				}
			}
		}
`, name, desc, macAddress, bootDeviceBlock, bootOrder, filepath)
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



		data "nutanix_images_v2" "ngt-image" {
		  filter = "name eq '${local.config.images.ngt_image}'"
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
									image_ext_id = data.nutanix_images_v2.ngt-image.images[0].ext_id
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
			depends_on = [data.nutanix_clusters_v2.clusters, data.nutanix_images_v2.ngt-image]
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
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
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
	filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
		nic_network_info{
			virtual_ethernet_nic_network_info{
				nic_type = "NORMAL_NIC"
				subnet{
					ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
				}
				vlan_mode = "ACCESS"
			}
		}
	}
	guest_customization {
		config {
		sysprep {
			install_type = "PREPARED"
				sysprep_script {
					unattend_xml {
						value = base64encode(file("%[4]s")) # unattend_xml file value, value is encoded in base64
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
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
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
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
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
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
				}
			}
			nics{
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
						}
						vlan_mode = "ACCESS"
					}
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

func testVmsV4ConfigWithSingleNic(name, desc string, r int) string {
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

		data "nutanix_subnets_v2" "subnet1" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "sc" {
			filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
			limit = 1
		}

		resource "nutanix_virtual_machine_v2" "test" {
			name                 = "%[1]s"
			description          = "%[2]s"
			num_cores_per_socket = 1
			num_sockets          = 1
			cluster {
				ext_id = local.cluster0
			}
			disks {
				disk_address {
					bus_type = "SCSI"
					index    = 0
				}
				backing_info {
					vm_disk {
						disk_size_bytes = "1073741824"
						storage_container {
							ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
						}
					}
				}
			}
			nics {
				nic_network_info {
					virtual_ethernet_nic_network_info {
						nic_type  = "NORMAL_NIC"
						vlan_mode = "ACCESS"
						subnet {
							ext_id = data.nutanix_subnets_v2.subnet1.subnets[0].ext_id
						}
					}
				}
			}
			power_state = "OFF"
		}
`, name, desc, r, filepath)
}

// testVmsV4ConfigWithSingleNicKeepSubnet is used for Step 3 to remove the second NIC
// but keep the test subnet resource to avoid deletion errors due to IP assignments
func testVmsV4ConfigWithSingleNicKeepSubnet(name, desc string, r int) string {
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

		data "nutanix_subnets_v2" "subnet1" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "sc" {
			filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
			limit = 1
		}

		# Keep the subnet resource to avoid deletion while IP might still be assigned.
		# Must match testVmsV4ConfigWithTwoNics exactly (Basic VLAN mirroring the
		# first NIC's subnet) so it is not recreated between steps.
		resource "nutanix_subnet_v2" "test_subnet" {
			name                   = "tf-test-subnet-%[3]d"
			description            = "Subnet for multi-NIC test"
			cluster_reference      = local.cluster0
			subnet_type            = "VLAN"
			network_id             = local.vmm.subnet.network_id + %[3]d %% 100
			is_advanced_networking = data.nutanix_subnets_v2.subnet1.subnets[0].is_advanced_networking
			bridge_name            = data.nutanix_subnets_v2.subnet1.subnets[0].bridge_name
		}

		resource "nutanix_virtual_machine_v2" "test" {
			name                 = "%[1]s"
			description          = "%[2]s"
			num_cores_per_socket = 1
			num_sockets          = 1
			cluster {
				ext_id = local.cluster0
			}
			disks {
				disk_address {
					bus_type = "SCSI"
					index    = 0
				}
				backing_info {
					vm_disk {
						disk_size_bytes = "1073741824"
						storage_container {
							ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
						}
					}
				}
			}
			# Only one NIC - removed the second NIC
			nics {
				nic_network_info {
					virtual_ethernet_nic_network_info {
						nic_type  = "NORMAL_NIC"
						vlan_mode = "ACCESS"
						subnet {
							ext_id = data.nutanix_subnets_v2.subnet1.subnets[0].ext_id
						}
					}
				}
			}
			power_state = "OFF"
		}
`, name, desc, r, filepath)
}

func testVmsV4ConfigWithTwoNics(name, desc string, r int) string {
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

		data "nutanix_subnets_v2" "subnet1" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "sc" {
			filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
			limit = 1
		}

		# Create a second subnet for multi-NIC testing. It must be the same network
		# type as the VM's first NIC subnet (local.vmm.subnet_name): AHV rejects a
		# VM whose NICs mix "Basic" and "Advanced" networks. That first subnet is a
		# Basic VLAN (is_advanced_networking=false, on bridge br0), so mirror those
		# values here rather than letting the subnet default to Advanced networking.
		resource "nutanix_subnet_v2" "test_subnet" {
			name                   = "tf-test-subnet-%[3]d"
			description            = "Subnet for multi-NIC test"
			cluster_reference      = local.cluster0
			subnet_type            = "VLAN"
			network_id             = local.vmm.subnet.network_id + %[3]d %% 100
			is_advanced_networking = data.nutanix_subnets_v2.subnet1.subnets[0].is_advanced_networking
			bridge_name            = data.nutanix_subnets_v2.subnet1.subnets[0].bridge_name
		}

		resource "nutanix_virtual_machine_v2" "test" {
			name                 = "%[1]s"
			description          = "%[2]s"
			num_cores_per_socket = 1
			num_sockets          = 1
			cluster {
				ext_id = local.cluster0
			}
			disks {
				disk_address {
					bus_type = "SCSI"
					index    = 0
				}
				backing_info {
					vm_disk {
						disk_size_bytes = "1073741824"
						storage_container {
							ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
						}
					}
				}
			}
			# First NIC - existing subnet
			nics {
				nic_network_info {
					virtual_ethernet_nic_network_info {
						nic_type  = "NORMAL_NIC"
						vlan_mode = "ACCESS"
						subnet {
							ext_id = data.nutanix_subnets_v2.subnet1.subnets[0].ext_id
						}
					}
				}
			}
			# Second NIC - new test subnet
			nics {
				nic_network_info {
					virtual_ethernet_nic_network_info {
						nic_type  = "NORMAL_NIC"
						vlan_mode = "ACCESS"
						subnet {
							ext_id = nutanix_subnet_v2.test_subnet.ext_id
						}
					}
				}
			}
			power_state = "OFF"
		}
`, name, desc, r, filepath)
}
