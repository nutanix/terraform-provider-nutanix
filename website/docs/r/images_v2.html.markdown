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

resource "nutanix_images_v2" "img-1" {
  name        = "test-image"
  description = "img desc"
  type        = "ISO_IMAGE"
  source {
    url_source {
      url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
    }
  }
}


resource "nutanix_images_v2" "img-2"{
  name = "test-image"
  description = "img desc"
  type = "DISK_IMAGE"
  source {
    url_source {
      url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
    }
  }
  cluster_location_ext_ids = [
        "ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
  ]
}

resource "nutanix_images_v2" "object-liteStore-img" {
  name        = "image-object-lite-example"
  description = "Image created from object store"
  type        = "DISK_IMAGE"
  source {
    object_lite_source {
      key = "img-lite-key-example"
    }
  }
  lifecycle {
    ignore_changes = [
      source
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

- `name`: (Required) The user defined name of an image.
- `description`: (Optional) The user defined description of an image.
- `type`: (Required) The type of an image. Valid values "DISK_IMAGE", "ISO_IMAGE"
- `checksum`: (Optional) The checksum of an image.
- `source`: (Optional) The source of an image. It can be a VM disk or a URL.
- `category_ext_ids`: (Optional) List of category external identifiers for an image.
- `cluster_location_ext_ids`: (Optional) List of cluster external identifiers where the image is located.

### checksum

- `hex_digest`: (Required) The SHA1/SHA256 digest of an image file in hexadecimal format.
- `object_type`: sha1 or sha256 type of image

### source
The `source` supports the following:
> Only one of the following sources can be specified at a time.

- `url_source`: (Optional) The URL for creating an image.
- `vm_disk_source`: (Optional) The URL for creating an image.
- `object_lite_source`: (Optional) The URL for creating an image.


#### url_source
The `url_source` supports the following:

- `url`: (Required) The URL for creating an image.
- `should_allow_insecure_url`: (Optional) Ignore the certificate errors, if the value is true. Default is false.
- `basic_auth`: (Optional) Basic authentication credentials for image source HTTP/S URL
- `basic_auth.username`: (Required) Username for basic authentication
- `basic_auth.password`: (Required) Password for basic authentication.

#### vm_disk_source
The `vm_disk_source` supports the following:

- `ext_id`: (Required) The external identifier of VM Disk.

#### object_lite_source
The `object_lite_source` supports the following:

- `key`: (Required) Key that identifies the source object in the bucket. The resource implies the bucket, 'vmm-images' for Image and 'vmm-ovas' for OVA.

## Attributes Reference

The following attributes are exported:

- `name`: The user defined name of an image.
- `description`: The user defined description of an image.
- `type`: The type of an image.
- `checksum`: The checksum of an image.
- `size_bytes`: The size in bytes of an image file.
- `source`: The source of an image. It can be a VM disk or a URL.
- `category_ext_ids`: List of category external identifiers for an image.
- `cluster_location_ext_ids`: List of cluster external identifiers where the image is located.
- `create_time`: Create time of an image.
- `last_update_time`: Last update time of an image.
- `owner_ext_id`: External identifier of the owner of the image
- `placement_policy_status`: Status of an image placement policy.

### placement_policy_status

- `placement_policy_ext_id`: Image placement policy external identifier.
- `compliance_status`: Compliance status for a placement policy.
- `enforcement_mode`: Indicates whether the placement policy enforcement is ongoing or has failed.
- `policy_cluster_ext_ids`: List of cluster external identifiers of the image location for the enforced placement policy.
- `enforced_cluster_ext_ids`: List of cluster external identifiers for the enforced placement policy.
- `conflicting_policy_ext_ids`: List of image placement policy external identifier that conflict with the current one.

## Import

This helps to manage existing entities which are not created through terraform. Images can be imported using the `UUID`(ext_id in V4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_images_v2" "import_image"{}

// execute this command in cli, UUID can be fetched using the datasource ex: data "nutanix_images" "fetch_images"{}
terraform import nutanix_images_v2.import_image <UUID>
```

See detailed information in [Nutanix Create Image V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Images/operation/createImage)
