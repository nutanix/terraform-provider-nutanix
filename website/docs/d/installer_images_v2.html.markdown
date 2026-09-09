---
layout: "nutanix"
page_title: "NUTANIX: nutanix_installer_images_v2"
sidebar_current: "docs-nutanix-datasource-installer-images-v2"
description: |-
  Returns a paginated list of all installer images registered in Foundation Central.
---

# nutanix_installer_images_v2

Returns a paginated list of all installer images registered in Foundation Central.

## Example

```hcl
data "nutanix_installer_images_v2" "list" {}

data "nutanix_installer_images_v2" "filtered" {
  limit  = 10
  filter = "name eq 'example-aos-installer-image'"
}
```

## Argument Reference

The following arguments are supported:

* `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.

## Attributes Reference

The following attributes are exported:

* `images`: List of installer images registered in Foundation Central. Each element exports the same attributes as the [nutanix_installer_image_v2](installer_image_v2.html) datasource.

See detailed information in [Nutanix List Installer Images V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
