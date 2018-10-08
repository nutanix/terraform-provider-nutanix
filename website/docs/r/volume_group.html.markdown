---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group"
sidebar_current: "docs-nutanix-resource-volume-group"
description: |-
  Provides a Nutanix Virtual Machine resource to Create a volume group.
---

# nutanix_volume_group

Provides a Nutanix Volume Group resource to Create a volume_group.

## Example Usage

```hcl
resource "nutanix_volume_group" "test_volume" {
  name        = "Test Volume Group"
  description = "Test Volume Group Description"
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) The name for the volume_group.
* `categories`: - (Optional) Categories for the volume_group.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.
* `api_version` - (Optional) Version of the API.
* `description`: - (Optional) A description for volume_group.
* `flash_mode`: - (Optional) Flash Mode, if enabled all volume disks of the VG will be pinned to SSD tier.
* `file_system_type`: - (Optional) File system to be used for volume.
* `sharing_status`: - (Optiional) Whether the volume group can be shared across multiple iSCSI initiators.
* `attachment_list`: - (Optional) VMs attached to volume group.
* `disk_list`: - (Optional) Volume group disk specification.
* `iscsi_target_prefix`: - (Optional) iSCSI target prefix-name.

### Disk List

The disk_list attribute supports the following:

* `vmdisk_UUID`: - (Optional) The UUID of this volume disk.
* `index`: - (Optional) Index of the volume disk in the group.
* `data_source_reference`: - (Optional) Reference to a kind
* `disk_size_mib`: - (Optional) Size of the disk in MiB.
* `storage_container_UUID`: - (Optional) Container UUID on which to create the disk.

### Attachment List

The attachment_list attribute supports the following:

* `vm_reference`: - (Optional) Reference to a kind.
* `iscsi_initiator_name`: - (Optional) Name of the iSCSI initiator of the workload outside Nutanix cluster.

## Attributes Reference

The following attributes are exported:

* `metadata`: - The volume_group kind metadata.
* `state`: - The state of the volume group.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when volume_group was last updated.
* `UUID`: - volume_group UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when volume_group was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - volume_group name.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, attributes supports the following:

* `kind`: - The kind name (Default value: project)(Required).
* `name`: - the name(Optional).
* `UUID`: - the UUID(Required).

Note: `vm_reference` and `data_source_reference` don't support `name` argument.

See detailed information in [Nutanix Volume Group](https://nutanix.github.io/Automation/experimental/swagger-redoc-sandbox/#tag/volume_group).