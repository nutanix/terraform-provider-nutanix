---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ova_vm_deploy_v2 "
sidebar_current: "docs-nutanix-resource-ova-download-v2"
description: |-
  Deploys a VM from an OVA, allowing you to override the VM configuration if needed.


---

# nutanix_ova_vm_deploy_v2
Deploys a VM from an OVA, allowing you to override the VM configuration if needed.



## Example Usage

```hcl
// Deploy vm from ova
resource "nutanix_ova_vm_deploy_v2" "test" {
  ext_id = "42e1fc04-8aa5-4572-8fa4-416f23767adb"
  override_vm_config {
    name              = "vm-from-ova"
    memory_size_bytes = 8 * 1024 * 1024 * 1024 # 8 GiB
    nics {
      backing_info {
        is_connected = true
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = "9bd6cbc2-a592-4728-ab89-473612f46b99"
        }
        vlan_mode     = "TRUNK"
        trunked_vlans = ["1"]
      }
    }
  }
  cluster_location_ext_id = "639ed004-31a6-4ff6-b06e-825617292811"
}
```

## Argument Reference

The following arguments are supported:

- `extId`: -(Required) The external identifier for an OVA.
- `override_vm_config`: -(Required) VM config override spec for OVA VM deploy endpoint
- `cluster_location_ext_id`: -(Optional) Cluster identifier to deploy VM from OVA. This field is required when deploying an OVA and must be a part of the OVA location list.

### Override VM Config
The `override_vm_config` arguments are support the following :

* `name`: (Optional) VM name.
* `num_sockets`: (Required) Number of vCPU sockets. Value should be at least 1.
* `num_cores_per_socket`: (Optional) Number of cores per socket. Value should be at least 1.
* `num_threads_per_core`: (Optional) Number of threads per core. Value should be at least 1.
* `memory_size_bytes`: (Required) Memory size in bytes.
* `nics`: (Optional) NICs attached to the VM.
* `cd_roms`: (Optional) CD-ROMs attached to the VM.
* `categories`: (Optional) Categories for the VM.


#### NICs
The `nics` attribute supports the following:

* `nic_backing_info`: (Optional) New NIC backing info (v2.4.1+). One of `virtual_ethernet_nic`, `sriov_nic`, `dp_offload_nic`.
* `nic_network_info`: (Optional) New NIC network info (v2.4.1+). One of `virtual_ethernet_nic_network_info`, `sriov_nic_network_info`, `dp_offload_nic_network_info`.
* `backing_info`: (Optional, Deprecated) Use `nic_backing_info.virtual_ethernet_nic` instead.
* `network_info`: (Optional, Deprecated) Use `nic_network_info.virtual_ethernet_nic_network_info` instead.

##### nics.backing_info
* `model`: (Optional) Options for the NIC emulation. Valid values "VIRTIO" , "E1000".
* `mac_address`: (Optional) MAC address of the emulated NIC.
* `is_connected`: (Optional) Indicates whether the NIC is connected or not. Default is True.
* `num_queues`: (Optional) The number of Tx/Rx queue pairs for this NIC. Default is 1.

##### nics.network_info
* `nic_type`: (Optional) NIC type. Valid values "SPAN_DESTINATION_NIC",  "NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC" .
* `network_function_chain`: (Optional) The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
* `network_function_nic_type`: (Optional) The type of this Network function NIC. Defaults to INGRESS.
* `subnet`: (Optional) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC
* `subnet.ext_id`: (Optional) The globally unique identifier of a subnet of type UUID.
* `vlan_mode`: (Optional) all the virtual NICs are created in ACCESS mode, which permits only one VLAN per virtual network. TRUNKED mode allows multiple VLANs on a single VM NIC for network-aware user VMs.
* `trunked_vlans`: (Optional) List of networks to trunk if VLAN mode is marked as TRUNKED. If empty and VLAN mode is set to TRUNKED, all the VLANs are trunked.
* `should_allow_unknown_macs`: (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not. This is applicable only for the NICs on the overlay subnets.
* `ipv4_config`: (Optional) The IP address configurations.

###### nics.ipv4_config
* `should_assign_ip`: If set to true (default value), an IP address must be assigned to the VM NIC - either the one explicitly specified by the user or allocated automatically by the IPAM service by not specifying the IP address. If false, then no IP assignment is required for this VM NIC.
* `ip_address`: The IP address of the NIC.
* `secondary_ip_address_list`: Secondary IP addresses for the NIC.

###### ip_address, secondary_ip_address_list
* `value`: The IPv4 address of the host.
* `prefix_length`: The prefix length of the IP address.

#### CD-ROMs
The `cd_roms` attribute supports the following:

* `disk_address`: (Optional) Virtual Machine disk (VM disk).
* `backing_info`: (Optional) Storage provided by Nutanix ADSF
* `iso_type`: Type of ISO image inserted in CD-ROM. Valid values "OTHER", "GUEST_TOOLS", "GUEST_CUSTOMIZATION" .


#### Categories
The `categories` attribute supports the following:

* `ext_id`: A globally unique identifier of a VM category of type UUID.


See detailed information in [Nutanix Deploy VMs from an OVA V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.1#tag/Ovas/operation/deployOva).
