---
layout: "nutanix"
page_title: "NUTANIX: nutanix_category_mapping_v2"
sidebar_current: "docs-nutanix-datasource-category-mapping-v2"
description: |-
  Gets the category to directory configuration information by ExtID.
---

# nutanix_category_mapping_v2

Gets the category to directory configuration information by ExtID.

## Example Usage

```hcl
data "nutanix_category_mapping_v2" "example" {
  ext_id = "<ext-id>"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) The external identifier of the category mapping.

## Attributes Reference

The following attributes are exported:

* `tenant_id` - The tenant identifier.
* `links` - A HATEOAS style link for the response.
* `name` - Name of the category mapping.
* `category_name` - The name for the category that this mapping is for.
* `category_value` - The value for the category that this mapping is for.
* `ad_info` - Active Directory information for the mapping.
  * `directory_service_reference` - The ExtID of the directory service.
  * `object_identifier` - The objectGUID for the object in AD.
  * `object_path` - The path for the mapped object in AD.
  * `status` - The mapping status. Values: `USABLE`, `DELETED`, `DIRECTORY_NOT_CONFIGURED`.
* `project_ext_id` - The external identifier of the associated project.
