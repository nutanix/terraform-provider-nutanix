---
layout: "nutanix"
page_title: "NUTANIX: nutanix_category_key"
sidebar_current: "docs-nutanix-datasource-category-key"
description: |-
  Describe a Nutanix Category Key and its values (if it has them).
---

# nutanix_category_key

Describe a Nutanix Category Key and its values (if it has them).

## Example Usage

```hcl
resource "nutanix_category_key" "test_key_value"{
    name = "data_source_category_key_test_values"
    description = "Data Source CategoryKey Test with Values"
}

resource "nutanix_category_value" "test_value"{
    name = nutanix_category_key.test_key_value.name
    value = "test_category_value_data_source"
    description = "Data Source CategoryValue Test with Values"
}


data "nutanix_category_key" "test_key_value" {
    name = nutanix_category_key.test_key_value.name
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) The name for the category key.

## Attributes Reference

The following attributes are exported:

* `system_defined`: - Specifying whether its a system defined category.
* `description`: - A description for category key.
* `api_version` - The version of the API.
* `values`: - A list of the values from this category key (if it has them).

See detailed information in [Nutanix Image](https://www.nutanix.dev/api_references/prism-central-v3/#/d9979ade0b152-get-a-category-key).
