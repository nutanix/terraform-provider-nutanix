---
layout: "nutanix"
page_title: "NUTANIX: nutanix_patched_images_v2"
sidebar_current: "docs-nutanix-datasource-patched-images-v2"
description: |-
  Returns a list of patched images.
---

# nutanix_patched_images_v2

Returns a list of patched images.

## Example

```hcl
data "nutanix_patched_images_v2" "example" {}
```

## Argument Reference

* `page`: (Optional) A URL query parameter that specifies the page number of the result set.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attributes Reference

* `patched_images`: List of patched images. Each element has the same attributes as the [nutanix_patched_image_v2](patched_image_v2.html) datasource.

See detailed information in [Nutanix List Patched Images V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
