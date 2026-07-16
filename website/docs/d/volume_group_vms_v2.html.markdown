---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_vms_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-vms-v2"
description: |-
  Query the list of VM attachments for a Volume Group. Deprecated: This API has been deprecated.
---

# nutanix_volume_group_vms_v2

Query the list of VM attachments for a Volume Group identified by {extId}. Deprecated: This API has been deprecated.

## Example Usage

```hcl

# List all the VM attachments for a Volume Group.
data "nutanix_volume_group_vms_v2" "list-vm-attachments" {
  volume_group_ext_id = "3770be9d-06be-4e25-b85d-3457d9b0ceb1"
}

# List all the VM attachments for a Volume Group with pagination.
data "nutanix_volume_group_vms_v2" "list-vm-attachments-paginated" {
  volume_group_ext_id = "3770be9d-06be-4e25-b85d-3457d9b0ceb1"
  page                = 0
  limit               = 10
}
```

## Argument Reference

The following arguments are supported:

* `volume_group_ext_id`: -(Required) The external identifier of a Volume Group.
* `page`: - A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit` : - A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` : - A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response.
* `orderby` : - A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default.

## Attributes Reference

The following attributes are exported:

* `attachments`: - List of VM attachments for a Volume Group identified by {extId}.

### Attachments

The `attachments` contains the list of VM attachments for the Volume Group. Each attachment contains the following attributes:

* `ext_id`: - The external identifier of the VM.
* `index`: - The index on the SCSI bus to attach the VM to the Volume Group. This is an optional field.

See detailed information in [Nutanix List VM Attachments by Volume Group Id V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.2#tag/VolumeGroups/operation/listVmAttachmentsByVolumeGroupId).
