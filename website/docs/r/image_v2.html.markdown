---
layout: "nutanix"
page_title: "NUTANIX: nutanix_images_v2"
sidebar_current: "docs-nutanix-resource-images-v2"
description: |-
  Provides a Nutanix Image resource to Create a Image.
---

# nutanix_images_v2

Create an image using the provided request body. Name, type and source are mandatory fields to create an image.


```hcl

    resource "nutanix_images_v4" "test" {
        name = "test-image"
        description = "img desc"
        type = "ISO_IMAGE"
        source{
            url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
        }
    }

    data "nutanix_clusters" "clusters" {}

    locals {
    cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
    }

    resource "nutanix_images_v4" "test" {
        name = "test-image"
        description = "img desc"
        type = "DISK_IMAGE"
        source{
            url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
        }
        cluster_location_ext_ids = [
            local.cluster0
        ]
    }
```

## Argument Reference

The following arguments are supported:
* `name`: (Required) The user defined name of an image.
* `description`: (Optional) The user defined description of an image.
* `type`: (Required) The type of an image. Valid values "DISK_IMAGE", "ISO_IMAGE"
* `checksum`: (Optional) The checksum of an image.
* `source`: (Optional) The source of an image. It can be a VM disk or a URL.
* `category_ext_ids`: (Optional) List of category external identifiers for an image.
* `cluster_location_ext_ids`: (Optional) List of cluster external identifiers where the image is located.

### checksum
* `hex_digest`: (Required) The SHA1/SHA256 digest of an image file in hexadecimal format.
* `object_type`: sha1 or sha256 type of image


### source
* `url`: (Optional) The URL for creating an image.
* `should_allow_insecure_url`: (Optional) Ignore the certificate errors, if the value is true. Default is false.
* `basic_auth`: (Optional) Basic authentication credentials for image source HTTP/S URL
*  `basic_auth.username`: (Required) Username for basic authentication
* `basic_auth.password`: (Required) Password for basic authentication.
* `ext_id`: (Optional) The external identifier of VM Disk.



## Attributes Reference

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

### placement_policy_status
* `placement_policy_ext_id`: Image placement policy external identifier.
* `compliance_status`: Compliance status for a placement policy.
* `enforcement_mode`: Indicates whether the placement policy enforcement is ongoing or has failed.
* `policy_cluster_ext_ids`: List of cluster external identifiers of the image location for the enforced placement policy.
* `enforced_cluster_ext_ids`: List of cluster external identifiers for the enforced placement policy.
* `conflicting_policy_ext_ids`: List of image placement policy external identifier that conflict with the current one.

See detailed information in [Nutanix Image](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1)