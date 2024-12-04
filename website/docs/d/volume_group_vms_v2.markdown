---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_vms_v2"
sidebar_current: "docs-nutanix-datasource-volume-groups-v2"
description: |-
  Describes a list of VM attachments for a Volume Group identified by {extId}.
---

# nutanix_volume_group_vms_v2

Describes a list of VM attachments for a Volume Group identified by {extId}.

## Example Usage

```hcl
resource "nutanix_volume_group_v2" "vg"{
  name                               = "test_volume_group"
  cluster_reference                  = "<Cluster uuid>"
}

resource "nutanix_volume_group_vm_v2" "vg_vm_test"{
  volume_group_ext_id = nutanix_volume_group_v2.vg.id
  vm_ext_id           = var.vg_vm_ext_id
}

# List all the VM attachments for a Volume Group.
data "nutanix_volume_group_vms_v2" "vg_vm_test"{
  ext_id = nutanix_volume_group_v2.vg.id
}
```

##  Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of the volume group.

* `page`: - A query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.

* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.

* `filter` : A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields: clusterReference, extId.

* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields: clusterReference, extId.

## Attributes Reference
The following attributes are exported:

* `vms_attachments`: - List of Volume Groups.

### vms_attachments

The vms attachments entities attribute element contains the followings attributes:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

#### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix Volumes V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0.b1).
