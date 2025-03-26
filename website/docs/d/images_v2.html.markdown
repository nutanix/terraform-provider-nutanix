---
layout: "nutanix"
page_title: "NUTANIX: nutanix_images_v2"
sidebar_current: "docs-nutanix-datasource-images-v2"
description: |-
 List of all images
---

# nutanix_images_v2

List images owned by Prism Central along with the image details like name, description, type, etc. This operation supports filtering, sorting, selection & pagination.

## Example

```hcl
data "nutanix_images_v2" "images"{}
```

## Attribute Reference
The following attributes are exported:

* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions
* `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default
* `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. 
* `images`: List of all images


## images
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `name`: The user defined name of an image.
* `description`: The user defined description of an image.
* `type`: The type of an image.
* `checksum`: The checksum of an image.
* `size_bytes`: The size in bytes of an image file.
* `source`: The source of an image. It can be a VM disk or a URL.
* `category_ext_ids`: List of category external identifiers for an image.
* `cluster_location_ext_ids`: List of cluster external identifiers where the image is located.
* `create_time`: Create time of an image.
* `last_update_time`: Last update time of an image.
* `owner_ext_id`: External identifier of the owner of the image
* `placement_policy_status`: Status of an image placement policy.


### source
* `ext_id`: The external identifier of VM Disk.
* `url`: The URL for creating an image.
* `basic_auth`: Basic authentication credentials for image source HTTP/S URL.
* `basic_auth.username`: Username for basic authentication.
* `basic_auth.password`: Password for basic authentication.


### placement_policy_status
* `placement_policy_ext_id`: Image placement policy external identifier.
* `compliance_status`: Compliance status for a placement policy.
* `enforcement_mode`: Indicates whether the placement policy enforcement is ongoing or has failed.
* `policy_cluster_ext_ids`: List of cluster external identifiers of the image location for the enforced placement policy.
* `enforced_cluster_ext_ids`: List of cluster external identifiers for the enforced placement policy.
* `conflicting_policy_ext_ids`: List of image placement policy external identifier that conflict with the current one.

See detailed information in [Nutanix Image V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0)