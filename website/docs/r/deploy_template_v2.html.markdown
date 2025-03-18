---
layout: "nutanix"
page_title: "NUTANIX: nutanix_deploy_templates_v2"
sidebar_current: "docs-nutanix-resource-deploy-templates-v2"
description: |-
  Deploy one or more VMs from a Template. Number of VMs to be deployed and their corresponding VM configuration overrides can be provided.
---

# nutanix_deploy_templates_v2

Deploy one or more VMs from a Template. Number of VMs to be deployed and their corresponding VM configuration overrides can be provided.

## Example

```hcl
resource "nutanix_deploy_templates_v2" "deploy-temp" {
    ext_id = "ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
    number_of_vms = 1
    cluster_reference = "0005b6b8-7b3b-4b0b-8b3b-7b3b4b0b8b3b"
    override_vm_config_map{
        name= "example-tf-temp"
        memory_size_bytes = 4294967296
        num_sockets=2
        num_cores_per_socket=1
        num_threads_per_core=1
    }
}
```


## Argument Reference

The following arguments are supported:
* `ext_id`: (Required) The identifier of a Template.
* `version_id`: (Optional) The identifier of a Template Version.
* `number_of_vms`: (Required) Number of VMs to be deployed.
* `override_vm_config_map`: (Optional) The map specifying the VM configuration overrides for each of the specified VM(s) to be created. The overrides can include the created VM Name, Configuration and Guest Customization.
* `cluster_reference`: (Required) The identifier of the Cluster where the VM(s) will be created using a Template.


### override_vm_config_map

* `name`: (Optional) VM name.
* `num_sockets`: (Optional) Number of vCPU sockets.
* `num_cores_per_socket`: (Optional) Number of cores per socket.
* `num_threads_per_core`: (Optional) Number of threads per core.
* `memory_size_bytes`: (Optional) Memory size in bytes.
* `nics`: (Optional) NICs attached to the VM.
* `guest_customization`: (Optional) Stage a Sysprep or cloud-init configuration file to be used by the guest for the next boot. Note that the Sysprep command must be used to generalize the Windows VMs before triggering this API call.


### nics

* `ext_id`: (Optional) A globally unique identifier of an instance that is suitable for external consumption.
* `backing_info`: (Optional) Defines a NIC emulated by the hypervisor
* `network_info`: (Optional) Network information for a NIC.

### nics.backing_info
* `model`: (Optional) Options for the NIC emulation. Valid values "VIRTIO" , "E1000".
* `mac_address`: (Optional) MAC address of the emulated NIC.
* `is_connected`: (Optional) Indicates whether the NIC is connected or not. Default is True.
* `num_queues`: (Optional) The number of Tx/Rx queue pairs for this NIC. Default is 1.

### nics.network_info
* `nic_type`: (Optional) NIC type. Valid values "SPAN_DESTINATION_NIC",  "NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC" .
* `network_function_chain`: (Optional) The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
* `network_function_nic_type`: (Optional) The type of this Network function NIC. Defaults to INGRESS.
* `subnet`: (Required) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC
* `vlan_mode`: (Required) all the virtual NICs are created in ACCESS mode, which permits only one VLAN per virtual network. TRUNKED mode allows multiple VLANs on a single VM NIC for network-aware user VMs.
* `trunked_vlans`: (Optional) List of networks to trunk if VLAN mode is marked as TRUNKED. If empty and VLAN mode is set to TRUNKED, all the VLANs are trunked.
* `should_allow_unknown_macs`: (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not. This is applicable only for the NICs on the overlay subnets.
* `ipv4_config`: (Optional) The IP address configurations.


### guest_customization

* `config`: (Required) The Nutanix Guest Tools customization settings.

* `config.sysprep`: (Optional) Sysprep config
* `config.cloud_init`: (Optional) CloudInit Config


### config.sysprep
* `install_type`: (Required) Indicates whether the guest will be freshly installed using this unattend configuration, or this unattend configuration will be applied to a pre-prepared image. Values allowed is 'PREPARED', 'FRESH'.

* `sysprep_script`: (Required) Object either UnattendXml or CustomKeyValues
* `sysprep_script.unattend_xml`: (Optional) xml object
* `sysprep_script.custom_key_values`: (Optional) The list of the individual KeyValuePair elements.


### config.cloud_init
* `datasource_type`: (Optional) Type of datasource. Default: CONFIG_DRIVE_V2
* `metadata`: The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded. Default value is 'CONFIG_DRIVE_V2'.
* `cloud_init_script`: (Optional) The script to use for cloud-init.
* `cloud_init_script.user_data`: (Optional) user data object
* `cloud_init_script.custom_keys`: (Optional) The list of the individual KeyValuePair elements.



See detailed information in [Nutanix Deploy Template V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Templates/operation/deployTemplate).
