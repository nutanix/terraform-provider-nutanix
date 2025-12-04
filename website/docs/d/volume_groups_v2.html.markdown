---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_groups_v2"
sidebar_current: "docs-nutanix-datasource-volume-groups-v2"
description: |-
  Describes a List all the Volume Groups.
---

# nutanix_volume_groups_v2

Describes a List all the Volume Groups.

## Example Usage

```hcl
# list all the Volume Groups.
data "nutanix_volume_groups_v2" "volume_groups"{}

# list all the Volume Groups with pagination.
data "nutanix_volume_groups_v2" "vg-pagination"{
  page = 1
  limit = 10
}

# list all the Volume Groups with filter.
data "nutanix_volume_groups_v2" "vg-filter"{
  filter = "name eq 'volume_group_test'"
}
```

##  Argument Reference

The following arguments are supported:

* `page`: - A query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` : A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - clusterReference
  - extId
  - name
* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - clusterReference
  - extId
  - name
* `expand` : A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are \$filter, \$select and \$orderby. The following expansion keys are supported. The expand can be applied to the following fields:
  - clusterReference
  - metadata
* `select` : A query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., \*), then all properties on the matching resource will be returned. The select can be applied to the following fields:
  - clusterReference
  - extId
  - name

## Attributes Reference
The following attributes are exported:

* `volume_groups`: - List of Volume Groups.

## Volume Groups
The `volume_groups` contains list of Volume Groups. Each Volume Group contains the following attributes:

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

#### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

#### Iscsi Features

The iscsi_features attribute supports the following:

* `enabled_authentications`: - The authentication type enabled for the Volume Group.

#### Storage Features

The storage features attribute supports the following:

* `flash_mode`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.

##### Flash Mode

The flash mode features attribute supports the following:

* `is_enabled`: - Indicates whether the flash mode is enabled for the Volume Group.

See detailed information in [Nutanix List Volume Groups V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/VolumeGroups/operation/listVolumeGroups).
