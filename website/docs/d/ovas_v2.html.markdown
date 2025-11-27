---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ovas_v2 "
sidebar_current: "docs-nutanix-datasource-ovas-v2"
description: |-
  This lists all accessible OVAs using the default pagination, which can be customized.

---

# nutanix_ovas_v2

This lists all accessible OVAs using the default pagination, which can be customized.



## Example Usage

```hcl

// Fetch all OVAs
data "nutanix_ovas_v2" "example"{
}

// filtered ovas on disk format
data "nutanix_ovas_v2" "example_filtered_disk_format"{
  filter = "diskFormat eq Vmm.Content.OvaDiskFormat'QCOW2'"
}

// filtered ovas on name
data "nutanix_ovas_v2" "example_filtered_name"{
  filter = "name eq 'example-ova'"
}

// filtered ovas on parentVm
data "nutanix_ovas_v2" "example_filtered_parent_vm"{
  filter = "parentVm eq 'LinuxServer_VM'"
}

// filtered ovas on sizeBytes
data "nutanix_ovas_v2" "example_filtered_size"{
  filter = "sizeBytes eq 57"
}

// limit, select, orderby and select example
data "nutanix_ovas_v2" "example_2"{
  filter = "startswith(parentVm, 'Linux')"
  limit = 10
  select = "diskFormat,extId,name,vmConfig,checksum"
  order_by = "name desc"
}

```

## Attribute Reference

The following attributes are exported:

* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. The filter can be applied to the following fields:
    - diskFormat
    - name
    - parentVm
    - sizeBytes

* `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
    - createTime
    - lastUpdateTime
    - name
    - sizeBytes
* `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. it can be applied to the following fields:
    - clusterLocationExtIds
    - createTime
    - createdBy
    - diskFormat
    - extId
    - lastUpdateTime
    - links
    - name
    - parentVm
    - sizeBytes
    - tenantId
    - vmConfig


## Attributes Reference
The following attributes are exported:

* `ovas`: List of all OVAs

## Attributes Reference

The following attributes are exported:

- `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `name`: - Name of the OVA.
- `checksum`: - The checksum of an OVA.
- `size_bytes`: - Size of OVA in bytes.
- `created_by`: - Information of the user.
- `parent_vm`: - The parent VM used for creating the OVA.
- `disk_format`: - Disk format of an OVA.
  |ENUM |Description |
  |---|---|
  | VMDK | The VMDK disk format of an OVA. |
  | QCOW2 | The QCOW2 disk format of an OVA. |
- `create_time`: - Time when the OVA was created time.
- `last_update_time`: - Time when the OVA was last updated time.

### Links

The `links` attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### checksum

The checksum argument supports the following :

- `ova_sha1_checksum`: - The SHA1 checksum of the OVA file.
- `ova_sha256_checksum`: - The SHA256 checksum of the OVA file.

#### ova_sha1_checksum, ova_sha256_checksum

The `ova_sha1_checksum` and `ova_sha256_checksum` arguments support the following:

- `hex_digest`: - The hexadecimal representation of the checksum.


See detailed information in [Nutanix List Ovas V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.1#tag/Ovas/operation/listOvas).
