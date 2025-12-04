---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_disk_v2"
sidebar_current: "docs-nutanix-resource-volume-group-disk-v2"
description: |-
  This operation submits a request to Creates a new Volume Disk.
---

# nutanix_volume_group_v2

Provides a resource to Creates a new Volume Disk.

## Example Usage

```hcl

# create new volume group disk  and attached it to the previous volume group
resource "nutanix_volume_group_disk_v2" "example"{
  volume_group_ext_id = "cf7de8b9-88ed-477d-a602-c34ab7174c01"
  index               = 1
  description         = "create volume disk example"
  disk_size_bytes     = 5368709120

  disk_data_source_reference {
    name        = "disk1"
    ext_id      = "1d92110d-26b5-46c0-8c93-20b8171373e0"
    entity_type = "STORAGE_CONTAINER"
    uris        = ["uri1", "uri2"]
  }

  disk_storage_features {
    flash_mode {
      is_enabled = false
    }
  }

  lifecycle {
    ignore_changes = [
      disk_data_source_reference, links
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

- `volume_group_ext_id `: -(Required) The external identifier of the Volume Group.

- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.

- `index`: - Index of the disk in a Volume Group. This field is optional and immutable.

- `disk_size_bytes`: - Size of the disk in bytes. This field is mandatory during Volume Group creation if a new disk is being created on the storage container.

- `description`: - Volume Disk description.

- `disk_data_source_reference`: -(Required) Disk Data Source Reference.
- `disk_storage_features`: - Storage optimization features which must be enabled on the Volume Disks. This is an optional field. If omitted, the disks will honor the Volume Group specific storage features setting.

#### Disk Data Source Reference

The disk_data_source_reference attribute supports the following:

- `ext_id`: - The external identifier of the Data Source Reference.
- `name`: - The name of the Data Source Reference for the Volume Group.
- `uris`: - The uri list of the Data Source Reference.
- `entity_type`: - The Entity Type of the Data Source Reference. valid values are:
  - STORAGE_CONTAINER
  - VM_DISK
  - VOLUME_DISK
  - DISK_RECOVERY_POINT

#### Disk Storage Features

The disk_storage_features attribute supports the following:

- `flash_mode`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.

##### Flash Mode

The flash mode features attribute supports the following:

- `is_enabled`: - Indicates whether the flash mode is enabled for the Volume Group Disk.

## Import

This helps to manage existing entities which are not created through terraform. Volume group disk can be imported using the `volume_group_ext_id/disk_ext_id`. (ext_id in v4 API context). eg,

**Note**:To import Volume Group Disk, you need to have the Volume Group Disk UUID, and provide it in the format mentioned above while importing.

```hcl
// create its configuration in the root module. For example:
resource "nutanix_volume_group_disk_v2" "import_vg_disk"{}

// execute the below command. UUID can be fetched using datasource. Example:

// list volume groups
data "nutanix_volume_groups_v2" "fetch_vgs"{}

// list disks for a volume group
data "nutanix_volume_group_disks_v2" "fetch_vg_disks"{
    volume_group_ext_id = data.nutanix_volume_groups_v2.fetch_vgs.volume_groups[0].ext_id
}

terraform import nutanix_volume_group_disk_v2.import_vg_disk <volume_group_ext_id>/<disk_ext_id>
```

See detailed information in [Nutanix Create Volume Disk V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/VolumeGroups/operation/createVolumeDisk).
