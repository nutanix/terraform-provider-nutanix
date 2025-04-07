---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_disks_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-disks-v2"
description: |-
  Describes a List all the Volume Disks attached to the Volume Group.
---

# nutanix_volume_group_disks_v2

Query the list of disks corresponding to a Volume Group identified by {volumeGroupExtId}.
## Example Usage

```hcl

# List all the Volume Disks attached to the Volume Group.
data "nutanix_volume_group_disks_v2" "list-volume-disks"{
  volume_group_ext_id = "3770be9d-06be-4e25-b85d-3457d9b0ceb1"
}

# list all the Volume Disks attached to the Volume Group with pagination.
data "nutanix_volume_group_disks_v2" "list-volume-disks"{
  volume_group_ext_id = "3770be9d-06be-4e25-b85d-3457d9b0ceb1"
  page = 1
  limit = 10
}

# list all the Volume Disks attached to the Volume Group with filter.
data "nutanix_volume_group_disks_v2" "list-volume-disks"{
  volume_group_ext_id = "3770be9d-06be-4e25-b85d-3457d9b0ceb1"
  filter = "storageContainerId eq '07c2da68-bb67-4535-9b2a-81504f6bb2e3'"
}
```

##  Argument Reference

The following arguments are supported:

* `volume_group_ext_id`: -(Required) The external identifier of the Volume Group.
* `page`: - A query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` : A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields: storageContainerId.
* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields: diskSizeBytes.
* `expand` : A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are \$filter, \$select and \$orderby. The following expansion keys are supported. The expand can be applied to the following fields: clusterReference, metadata.
* `select` : A query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., \*), then all properties on the matching resource will be returned. The select can be applied to the following fields: extId, storageContainerId.

## Attributes Reference
The following attributes are exported:

* `disks`: - List of disks corresponding to a Volume Group identified by {volumeGroupExtId}.

### Disks
The `disks` contains list of Volume Disks attached to the Volume Group. Each Volume Disk contains the following attributes:

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

See detailed information in [Nutanix List all the Volume Disks attached to the Volume Group V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/VolumeGroups/operation/listVolumeDisksByVolumeGroupId).
