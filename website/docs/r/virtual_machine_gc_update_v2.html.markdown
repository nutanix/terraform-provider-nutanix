---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_gc_update_v2"
sidebar_current: "docs-nutanix-resource-vm-gc-update-v2"
description: |-
  Provides a Nutanix Virtual Machine resource to Create a virtual machine guest customization update.
---

# nutanix_vm_gc_update_v2

Provides a Nutanix Virtual Machine resource to Create a virtual machine guest customization update.

## Example Usage

```hcl
data "nutanix_virtual_machines_v2" "vm-list"{}

resource "nutanix_vm_gc_update_v2" "vm-gc-update"{
  ext_id = data.nutanix_virtual_machines_v2.vm-list.vms.0.data.ext_id
  config{
    cloud_init{
      cloud_init_script{
        user_data{
          value="IyEvYmluL2Jhc2gKZWNobyAiSGVsbG8gV29ybGQiCg=="
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: - (Required) The globally unique identifier of a VM. It should be of type UUID.
* `config`: - (Optional) The Nutanix Guest Tools customization settings.

### Config

The config attribute supports the following:

* `sysprep`: - (Optional) VM guests may be customized at boot time using one of several different methods. Currently, cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting to change it after creation will result in an error. Additional properties can be specified. For example - in the context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own custom script.

* `cloud_init`: - (Optional) VM guests may be customized at boot time using one of several different methods. Currently, cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting to change it after creation will result in an error. Additional properties can be specified. For example - in the context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own custom script.

### Sysprep

The sysprep attribute supports the following:

* `install_type`: - (Optional) Whether the guest will be freshly installed using this unattend configuration, or whether this unattend configuration will be applied to a pre-prepared image. Default is `PREPARED`.
    Valid values are:
    - `PREPARED` is done when sysprep is used to finalize Windows installation from an installed Windows and file name it is searching `unattend.xml` for `unattend_xml` parameter
    - `FRESH` is done when sysprep is used to install Windows from ISO and file name it is searching `autounattend.xml` for `unattend_xml` parameter
* `unattend_xml`: - (Optional) Generic key value pair used for custom attributes.

### Cloud Init

The cloud_init attribute supports the following:

* `datasource_type`: - (Optional) Type of datasource.
Default: CONFIG_DRIVE_V2Default is `CONFIG_DRIVE_V2`.
    Valid values are:
    - `CONFIG_DRIVE_V2` The type of datasource for cloud-init is Config Drive V2.
* `metadata` - (Optional) The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded.
* `cloud_init_script`: - (Optional) The script to use for cloud-init.

### Cloud Init Script

The cloud_init_script attribute supports the following:

* `user_data`: - (Optional) The contents of the user_data configuration for cloud-init. This can be formatted as YAML, JSON, or could be a shell script. The value must be base64 encoded.
* `custom_key_values`: - (Optional) Generic key value pair used for custom attributes in cloud init.

### User Data

The user_data attribute supports the following:

* `value`: - (Optional) The value for the cloud-init user_data.

### Custom Key Values

The custom_key_values attribute supports the following:

* `key_value_pairs`: - (Optional) The list of the individual KeyValuePair elements.

### Key Value Pairs

The key_value_pairs attribute supports the following:

* `name`: - (Optional) The key of this key-value pair
* `value`: - (Optional) The value associated with the key for this key-value pair.

See detailed information in [Nutanix Customize Gest VM V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/customizeGuestVm).

