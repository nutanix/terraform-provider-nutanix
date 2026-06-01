---
layout: "nutanix"
page_title: "NUTANIX: nutanix_category_mappings_v2"
sidebar_current: "docs-nutanix-datasource-category-mappings-v2"
description: |-
  Gets the list of Directory Server Category Mappings.
---

# nutanix_category_mappings_v2

Gets the list of Directory Server Category Mappings.

## Example Usage

```hcl
data "nutanix_category_mappings_v2" "all" {}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) Page number for pagination.
* `limit` - (Optional) Number of items per page.
* `filter` - (Optional) OData filter expression.
* `order_by` - (Optional) OData order by expression.
* `select` - (Optional) OData select expression.

## Attributes Reference

The following attributes are exported:

* `category_mappings` - List of category mappings. Each element contains:
  * `ext_id` - The external identifier.
  * `tenant_id` - The tenant identifier.
  * `links` - A HATEOAS style link for the response.
  * `name` - Name of the category mapping.
  * `category_name` - The name for the category.
  * `category_value` - The value for the category.
  * `ad_info` - Active Directory information for the mapping.
  * `project_ext_id` - The external identifier of the associated project.
