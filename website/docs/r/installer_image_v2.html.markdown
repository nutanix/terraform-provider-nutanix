---
layout: "nutanix"
page_title: "NUTANIX: nutanix_installer_image_v2"
sidebar_current: "docs-nutanix-resource-installer-image-v2"
description: |-
  Registers a new installer image from an external URL in Foundation Central.
---

# nutanix_installer_image_v2

Registers a new installer image from an external URL in Foundation Central. An
installer image bundles AOS and hypervisor software used for node provisioning.

## Example

```hcl
resource "nutanix_installer_image_v2" "example" {
  name    = "example-aos-installer-image"
  type    = "AOS"
  source  = "REMOTE_URL"
  url     = "https://example.com/aos/nos.tar.gz"
  version = "6.5.1"
}
```

## Argument Reference

The following arguments are supported:

* `name`: (Required) Name of the image. This field is required when creating an image.
* `type`: (Required) Type of the installer image. One of `AOS`, `AHV` or `ESX`.
* `source`: (Required) Source of the image. One of `LOCAL` or `REMOTE_URL`.
* `url`: (Optional) URL from where the image can be downloaded. This can be provided only through the update operation and only for images with `REMOTE_URL` as the source type.
* `version`: (Optional) Version of the image.
* `certificate_chain`: (Optional) Certificate chain for the image URL.
* `metadata_download_url`: (Optional) URL from where the image metadata can be downloaded for an AOS image.
* `checksum`: (Optional) Checksum value of the image. This is required and applicable only for AHV images. See [Checksum](#checksum) below.

### Checksum

The `checksum` block supports exactly one of the following nested blocks:

* `sha256`: (Optional) SHA-256 checksum of the image.
  * `hex_digest`: (Required) SHA-256 checksum value in hexadecimal format (64 characters).
* `md5`: (Optional) MD5 checksum of the image.
  * `hex_digest`: (Required) MD5 checksum value in hexadecimal format (32 characters).

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `file_status`: Status of the image file on the cluster.
* `metadata_status`: Status of the image metadata file on the cluster.
* `created_time`: Time when the image entity was created.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
  * `href`: The URL at which the entity described by the link can be accessed.
  * `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.

## Import

Installer images can be imported using their external ID, e.g.

```
terraform import nutanix_installer_image_v2.example <ext_id>
```

See detailed information in [Nutanix Create Installer Image V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
