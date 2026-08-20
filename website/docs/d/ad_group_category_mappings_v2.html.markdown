---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ad_group_category_mappings_v2"
sidebar_current: "docs-nutanix-datasource-ad-group-category-mappings-v2"
description: |-
  Gets the list of Directory Server Category Mappings.
---

# nutanix_ad_group_category_mappings_v2

Gets the list of Directory Server Category Mappings.

## Example Usage

```hcl
# List all Category Mappings, optionally filtered.
data "nutanix_ad_group_category_mappings_v2" "mappings-list" {
  limit  = 10
  filter = "name eq '<catgeory mapping name>'"
}
```

## Argument Reference

* `page` - (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.

## Attributes Reference

* `category_mappings` - List of Category Mappings. Each element exposes the same attributes as the `nutanix_ad_group_category_mapping_v2` datasource:
  * `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
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

See detailed information in [Nutanix Directory Server Configs V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.3#tag/DirectoryServerConfigs/operation/listDirectoryServerConfigs).
