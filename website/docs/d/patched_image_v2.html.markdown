---
layout: "nutanix"
page_title: "NUTANIX: nutanix_patched_image_v2"
sidebar_current: "docs-nutanix-datasource-patched-image-v2"
description: |-
  Returns details of a patched image identified by its external ID.
---

# nutanix_patched_image_v2

Returns details of a patched image identified by its external ID.

## Example

```hcl
data "nutanix_patched_image_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id`: (Required) A globally unique identifier of an instance that is suitable for external consumption.

## Attributes Reference

* `claim_token_ext_id`: ExtId of the claim token used in the patched image.
* `name`: Patched image name.
* `host_type`: Type of the host installed on the node.
* `version`: Version of the base hypervisor or host OS image used for patching.
* `image_details`: Details of the image used for patching (`local_host_image_ext_id`).
* `node_configurations`: List of node configurations used for patching the image.
* `patched_iso_url`: URL for downloading the patched image.
* `patched_iso_sha256_checksum`: SHA-256 checksum of the patched image.
* `created_time`: Timestamp when the patched image was created.
* `owner_ext_id`: External ID of the owner.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.

See detailed information in [Nutanix Get Patched Image V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
