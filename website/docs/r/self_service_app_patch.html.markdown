---
layout: "nutanix"
page_title: "NUTANIX: nutanix_self_service_app_patch"
sidebar_current: "docs-nutanix_self_service_app"
description: |-
  Run the specified patch on the application.
---

# nutanix_self_service_app_patch

Run the specified patch on the application by running patch action to update vm configuration, add nics, add disks, add/delete categories.

## Example 1: Update VM Configuration

This will run set patch config action in application.

```hcl
# Provision Application
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION"
    app_description = "DESCRIPTION OF APPLICATION"
}

# Run patch config (update config)
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = "NAME OF PATCH ACTION"
    config_name = "NAME OF PATCH CONFIG"
}
```

## Example 2: Update VM Configuration with runtime editable

```hcl
# Provision Application
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION"
    app_description = "DESCRIPTION OF APPLICATION"
}

# Run patch config (update config)
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = "NAME OF PATCH ACTION"
    config_name = "NAME OF PATCH CONFIG"
    vm_config {
        memory_size_mib = "SIZE IN MiB"
        num_sockets = "vCPU count"
        num_vcpus_per_socket = "NUMBER OF CORES VCPU"
    }
}
```

## Example 3: Add Category

```hcl
# Provision Application
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION"
    app_description = "DESCRIPTION OF APPLICATION"
}

# Run patch config (update config)
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = "NAME OF PATCH ACTION"
    config_name = "NAME OF PATCH CONFIG"
    categories {
        value = "CATEGORY TO BE ADDED (KEY:VALUE PAIR)"
        operation = "add"
    }
}
```

## Example 4: Delete Category

```hcl
# Provision Application
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION"
    app_description = "DESCRIPTION OF APPLICATION"
}

# Run patch config (update config)
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = "NAME OF PATCH ACTION"
    config_name = "NAME OF PATCH CONFIG"
    categories {
        value = "CATEGORY TO BE ADDED (KEY:VALUE PAIR)"
        operation = "delete"
    }
}
```

## Example 5: Add Disk

```hcl
# Provision Application
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION"
    app_description = "DESCRIPTION OF APPLICATION"
}

# Run patch config (update config)
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = "NAME OF PATCH ACTION"
    config_name = "NAME OF PATCH CONFIG"
    disks {
        disk_size_mib = "SIZE OF DISK IN MiB"
        operation = "add"
    }
}
```

## Example 6: Add Nic

```hcl
# Provision Application
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION"
    app_description = "DESCRIPTION OF APPLICATION"
}

# Run patch config (update config)
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = "NAME OF PATCH ACTION"
    config_name = "NAME OF PATCH CONFIG"
    nics {
        index = "DUMMY INDEX VALUE"
        operation = "add"
        subnet_uuid = "VALID SUBNET UUID IN PROJECT ATTACHED TO APP"
    }
}
```



## Argument Reference

The following arguments are supported:

* `app_uuid`: - (Required) The UUID of the application.
* `patch_name`: - (Required) The name of the patch to be applied. This is used to identify the action name which needs to be executed to update an application.
* `config_name`: - (Required) The name of the patch configuration. (<b>Same as patch_name for SINGLE VM)</b>


## Attribute Reference

The following attributes are exported:

* `runlog_uuid`: - (Computed) The UUID of the runlog that records the patch operation's execution details.

### vm_config

A list of virtual machine configuration changes. You can modify the VM's memory, number of sockets, and vCPUs per socket.

* `memory_size_mib`: - (Optional) The amount of memory (in MiB) to allocate for the VM.
* `num_sockets`: - (Optional) The number of vCPUs to assign.
* `num_vcpus_per_socket`: - (Optional) The number of cores per vCPU to assign to the VM.

### vm_config

A list of virtual machine configuration changes. You can modify the VM's memory, number of sockets, and vCPUs per socket.

* `memory_size_mib`: - (Optional, integer) The amount of memory (in MiB) to allocate for the VM.
* `num_sockets`: - (Optional, integer) The number of vCPUs to assign.
* `num_vcpus_per_socket`: - (Optional, integer) The number of cores per vCPU to assign to the VM.

### nics

A list of network interface changes.

* `index`: - (Optional, string) The index of the NIC. A dummy string for now.
* `operation`: - (Optional, string) The operation to perform on the NIC.
* `subnet_uuid`: - (Optional, string) The UUID of the subnet to which the NIC should be attached. 

### disks

A list of disk changes.

* `operation`: - (Optional, string) The operation to perform on the disk.
* `disk_size_mib`: - (Optional, integer) The size of the disk to allocate (in MiB).

### categories

A list of category modifications.

* `operation`: - (Optional) The operation to perform on the category. (e.g. "add", "delete")
* `value`: - (Optional, string) The value of the category. A Key:Value pair (e.g. "AppType:Oracle_DB"). There should not be any space in value.


See detailed information in [Run patch in app](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/Apps/paths/~1apps~1%7Buuid%7D~1patch~1%7Bpatch_uuid%7D~1run/post).