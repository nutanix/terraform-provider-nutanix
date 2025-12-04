---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_v2"
sidebar_current: "docs-nutanix-datasource-image-v2"
description: |-
 Describes a Image
---

# nutanix_image_v2

Retrieve the image details for the provided external identifier.

## Example

```hcl
data "nutanix_image_v2" "get-image"{
    ext_id = "0005a7b1-0b3b-4b3b-8b3b-0b3b4b3b4b3b"
}

```


## Argument Reference

The following arguments are supported:

* `ext_id`: The external identifier of an image.

## Attribute Reference

The following attributes are exported:

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

See detailed information in [Nutanix Get Image](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Images/operation/getImageById)
