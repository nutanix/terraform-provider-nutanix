---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_disk_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-disk-v2"
description: |-
  Describes the details of a Volume Disk.
---

# nutanix_volume_group_disk_v2

Describes a Query the Volume Disk identified by {extId} in the Volume Group identified by {volumeGroupExtId}.

## Example Usage

```hcl
# Get the details of a Volume Disk attached to the Volume Group.
data "nutanix_volume_group_disk_v2" "example"{
  volume_group_ext_id = "3770be9d-06be-4e25-b85d-3457d9b0ceb1"
  ext_id              = "1d92110d-26b5-46c0-8c93-20b8171373e0"
}
```

## Argument Reference

The following arguments are supported:

* `volume_group_ext_id `: -(Required) The external identifier of the Volume Group.
* `ext_id `: -(Required) The external identifier of the Volume Disk.


## Attributes Reference

The following attributes are exported:
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `index`: - Index of the disk in a Volume Group. This field is optional and immutable.
* `disk_size_bytes`: - ize of the disk in bytes. This field is mandatory during Volume Group creation if a new disk is being created on the storage container.
* `storage_container_id`: - Storage container on which the disk must be created. This is a read-only field.
* `description`: - Volume Disk description.
* `disk_data_source_reference`: - Disk Data Source Reference.
* `disk_storage_features`: - Storage optimization features which must be enabled on the Volume Disks. This is an optional field. If omitted, the disks will honor the Volume Group specific storage features setting.

#### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

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

See detailed information in [Nutanix Get Volume Disk V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/VolumeGroups/operation/getVolumeDiskById).
