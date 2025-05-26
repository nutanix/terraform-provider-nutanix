---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_v2"
sidebar_current: "docs-nutanix-resource-volume-group-v2"
description: |-
  This operation submits a request to Create a new Volume Group.
---

# nutanix_volume_group_v2

Provides a resource to Create a new Volume Group.

## Example Usage

```hcl

resource "nutanix_volume_group_v2" "volume_group_example"{
  name                               = "volume_group_test"
  description                        = "Test Create Volume group with spec"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  target_name                        = "volumegroup-test-001234"
  created_by                         = "example"
  cluster_reference                  = "ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
  iscsi_features {
    enabled_authentications = "CHAP"
    target_secret           = "123456789abc"
  }

  storage_features {
    flash_mode {
      is_enabled = true
    }
  }
  usage_type = "USER"
  is_hidden  = false

  # ignore changes to target_secret, target secret will not be returned in terraform plan output
  lifecycle {
    ignore_changes = [
      iscsi_features[0].target_secret
    ]
  }
}
```

## Argument Reference
The following arguments are supported:


* `ext_id`: -(Optional) A globally unique identifier of an instance that is suitable for external consumption.
* `name`: -(Required) Volume Group name. This is an optional field.
* `description`: -(Optional) Volume Group description. This is an optional field.
* `should_load_balance_vm_attachments`: -(Optional) Indicates whether to enable Volume Group load balancing for VM attachments. This cannot be enabled if there are iSCSI client attachments already associated with the Volume Group, and vice-versa. This is an optional field.
* `sharing_status`: -(Optional) Indicates whether the Volume Group can be shared across multiple iSCSI initiators. The mode cannot be changed from SHARED to NOT_SHARED on a Volume Group with multiple attachments. Similarly, a Volume Group cannot be associated with more than one attachment as long as it is in exclusive mode. This is an optional field. Valid values are SHARED, NOT_SHARED
* `target_name`: -(Optional) Name of the external client target that will be visible and accessible to the client.
* `enabled_authentications`: -(Optional) The authentication type enabled for the Volume Group. Valid values are CHAP, NONE
* `iscsi_features`: -(Optional) iSCSI specific settings for the Volume Group.
* `created_by`: -(Optional) Service/user who created this Volume Group.
* `cluster_reference`: -(Required) The UUID of the cluster that will host the Volume Group.
* `storage_features`: -(Optional) Storage optimization features which must be enabled on the Volume Group.
* `usage_type`: -(Optional) Expected usage type for the Volume Group. This is an indicative hint on how the caller will consume the Volume Group.  Valid values are BACKUP_TARGET, INTERNAL, TEMPORARY, USER
* `attachment_type`: -(Optional) The field indicates whether a VG has a VM or an external attachment associated with it. Valid values are :
  - EXTERNAL : Volume Group has an external iSCSI or NVMf attachment.
  - NONE : Volume Group has no attachment.
  - DIRECT : Volume Group has a VM attachment.
* `protocol`: -(Optional) Type of protocol to be used for Volume Group. Valid values are :
  - NOT_ASSIGNED :  Volume Group does not use any protocol.
  - ISCSI : Volume Group uses iSCSI protocol.
  - NVMF : Volume Group uses NVMf protocol.
* `is_hidden`: -(Optional) Indicates whether the Volume Group is meant to be hidden or not.
* `disks`: -(Optional) A list of Volume Disks to be attached to the Volume Group.

### Iscsi Features

The iscsi_features attribute supports the following:

* `enabled_authentications`: - The authentication type enabled for the Volume Group.

### Storage Features

The storage features attribute supports the following:

* `flash_mode`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.

#### Flash Mode

The flash mode features attribute supports the following:

* `is_enabled`: - Indicates whether the flash mode is enabled for the Volume Group.

### Disks

The disks attribute supports the following:

* `index`: - Index of the disk in a Volume Group. This field is optional and immutable.
* `disk_size_bytes`: - ize of the disk in bytes. This field is mandatory during Volume Group creation if a new disk is being created on the storage container.
* `description`: - Volume Disk description.
* `disk_data_source_reference`: -(Required) Disk Data Source Reference.
* `disk_storage_features`: - Storage optimization features which must be enabled on the Volume Disks. This is an optional field. If omitted, the disks will honor the Volume Group specific storage features setting.


#### Disk Data Source Reference

The disk_data_source_reference attribute supports the following:

* `ext_id`: - The external identifier of the Data Source Reference.
* `name`: - The name of the Data Source Reference.bled for the Volume Group.
* `uris`: - The uri list of the Data Source Reference.
* `entity_type`: - The Entity Type of the Data Source Reference.

#### Disk Storage Features

The disk_storage_features attribute supports the following:

* `flash_mode`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.

##### Flash Mode

The flash mode features attribute supports the following:

* `is_enabled`: - Indicates whether the flash mode is enabled for the Volume Group Disk.

See detailed information in [Nutanix Create Volume Group V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/VolumeGroups/operation/createVolumeGroup).
