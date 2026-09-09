---
layout: "nutanix"
page_title: "NUTANIX: nutanix_installer_image_v2"
sidebar_current: "docs-nutanix-datasource-installer-image-v2"
description: |-
  Returns details of an installer image identified by its external ID.
---

# nutanix_installer_image_v2

Returns details of an installer image identified by its external ID.

## Example

```hcl
data "nutanix_installer_image_v2" "get" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) A globally unique identifier of an instance that is suitable for external consumption.

## Attributes Reference

The following attributes are exported:

* `name`: Name of the image.
* `type`: Type of the installer image. One of `AOS`, `AHV` or `ESX`.
* `source`: Source of the image. One of `LOCAL` or `REMOTE_URL`.
* `url`: URL from where the image can be downloaded.
* `version`: Version of the image.
* `certificate_chain`: Certificate chain for the image URL.
* `metadata_download_url`: URL from where the image metadata can be downloaded for an AOS image.
* `checksum`: Checksum value of the image. This is applicable only for AHV images.
  * `sha256`: SHA-256 checksum of the image.
    * `hex_digest`: SHA-256 checksum value in hexadecimal format (64 characters).
  * `md5`: MD5 checksum of the image.
    * `hex_digest`: MD5 checksum value in hexadecimal format (32 characters).
* `file_status`: Status of the image file on the cluster.
* `metadata_status`: Status of the image metadata file on the cluster.
* `created_time`: Time when the image entity was created.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.
  * `href`: The URL at which the entity described by the link can be accessed.
  * `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.

See detailed information in [Nutanix Get Installer Image V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
