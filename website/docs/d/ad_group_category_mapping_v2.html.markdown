---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ad_group_category_mapping_v2"
sidebar_current: "docs-nutanix-datasource-ad-group-category-mapping-v2"
description: |-
  Gets the category to directory configuration information with the provided ExtID.
---

# nutanix_ad_group_category_mapping_v2

Gets the category to directory configuration information with the provided ExtID.

## Example Usage

```hcl
# Read a single Category Mapping by ext_id.
data "nutanix_ad_group_category_mapping_v2" "get-mapping" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id` - (Required) A globally unique identifier of the Category Mapping.

## Attributes Reference

* `name` - Name of Category Mapping.
* `category_name` - The name for the category that this mapping is for.
* `category_value` - The value for the category that this mapping is for.
* `ad_info` - The Active Directory object information for this mapping.
  * `directory_service_reference` - The ExtID of the directory service that will be used for mapping.
  * `object_identifier` - The objectGUID for the object in AD.
  * `object_path` - The path for the mapped object in AD.
  * `status` - The mapping status of AD Mapping.
* `project_ext_id` - The external identifier of the project associated with the Category Mapping.
* `links` - A HATEOAS style link for the response.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Directory Server Configs V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.3).
