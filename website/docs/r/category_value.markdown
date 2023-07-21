---
layout: "nutanix"
page_title: "NUTANIX: nutanix_category_value"
sidebar_current: "docs-nutanix-resource-category-value"
description: |-
  Provides a Nutanix Category value resource to Create a category value.
---

# nutanix_category_value

Provides a Nutanix Category value resource to Create a category value.

## Example Usage

```hcl
resource "nutanix_category_key" "test-category-key"{
    name = "app-support-1"
    description = "App Support Category Key"
}


resource "nutanix_category_value" "test"{
    name = nutanix_category_key.test-category-key.id
    description = "Test Category Value"
    value = "test-value"
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) The category_key name for the category value.
* `value` - (Required) The value for the category value.
* `description`: - (Optional) A description for category value.

## Attributes Reference

The following attributes are exported:

* `system_defined`: - Specifying whether its a system defined category.
* `api_version` - (Optional) The version of the API.

See detailed information in [Nutanix Category Value](https://www.nutanix.dev/api_references/prism-central-v3/#/7032d646db45b-create-or-update-a-category-value).
