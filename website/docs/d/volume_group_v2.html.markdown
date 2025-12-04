---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-v2"
description: |-
  Describes a Volume Group.
---

# nutanix_volume_group_v2

Query the Volume Group identified by {extId}.


## Example Usage

```hcl
data "nutanix_volume_group_v2" "volume_group"{
  ext_id = "d09aeec9-5bb7-4bfd-9717-a051178f6e7c"
}

```

## Argument Reference

The following arguments are supported:

* `ext_id `: -(Required) The external identifier of the Volume Group.


## Attributes Reference

The following attributes are exported:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`: -(Required) Volume Group name. This is an optional field.
* `description`: - Volume Group description. This is an optional field.
* `should_load_balance_vm_attachments`: - Indicates whether to enable Volume Group load balancing for VM attachments. This cannot be enabled if there are iSCSI client attachments already associated with the Volume Group, and vice-versa. This is an optional field.
* `sharing_status`: - Indicates whether the Volume Group can be shared across multiple iSCSI initiators. The mode cannot be changed from SHARED to NOT_SHARED on a Volume Group with multiple attachments. Similarly, a Volume Group cannot be associated with more than one attachment as long as it is in exclusive mode. This is an optional field. Valid values are SHARED, NOT_SHARED
* `target_name`: - Name of the external client target that will be visible and accessible to the client.
* `enabled_authentications`: - The authentication type enabled for the Volume Group. Valid values are CHAP, NONE
* `iscsi_features`: - iSCSI specific settings for the Volume Group.
* `created_by`: - Service/user who created this Volume Group.
* `cluster_reference`: - The UUID of the cluster that will host the Volume Group.
* `storage_features`: - Storage optimization features which must be enabled on the Volume Group.
* `usage_type`: - Expected usage type for the Volume Group. This is an indicative hint on how the caller will consume the Volume Group.  Valid values are BACKUP_TARGET, INTERNAL, TEMPORARY, USER
* `is_hidden`: - Indicates whether the Volume Group is meant to be hidden or not.

### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Iscsi Features

The iscsi_features attribute supports the following:

* `enabled_authentications`: - The authentication type enabled for the Volume Group.

### Storage Features

The storage features attribute supports the following:

* `flash_mode`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.

#### Flash Mode

The flash mode features attribute supports the following:

* `is_enabled`: - Indicates whether the flash mode is enabled for the Volume Group.

See detailed information in [Nutanix Get Volume Group v4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/VolumeGroups/operation/getVolumeGroupById).
