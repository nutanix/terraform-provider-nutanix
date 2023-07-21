---
layout: "nutanix"
page_title: "NUTANIX: nutanix_category_key"
sidebar_current: "docs-nutanix-resource-category-key"
description: |-
  Provides a Nutanix Category key resource to Create a category key name.
---

# nutanix_category_key

Provides a Nutanix Category key resource to Create a category key name.

## Example Usage

```hcl
resource "nutanix_category_key" "test"{
    name = "app-support-example"
    description = "App Support Category Key"
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) The name for the category key.
* `description`: - (Optional) A description for category key.

## Attributes Reference

The following attributes are exported:

* `system_defined`: - Specifying whether its a system defined category.
* `api_version` - (Optional) The version of the API.

See detailed information in [Nutanix Category Key](https://www.nutanix.dev/api_references/prism-central-v3/#/c67191febf803-create-or-update-a-category-key).
