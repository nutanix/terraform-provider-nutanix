---
layout: "nutanix"
page_title: "NUTANIX: nutanix_category_v2"
sidebar_current: "docs-nutanix-datasource-category-v2"
description: |-
  Fetch details of a category with the given external identifier.
---

# nutanix_category_v2
Fetch a category


## Example

```hcl

data "nutanix_category_v2" "get-category"{
  ext_id = "85e68112-5b2b-4220-bc8d-e529e4bf420e"
}

```


## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) The extID for the category.
* `expand`: (Optional)A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are \$filter, \$select and \$orderby. The following expansion keys are supported:
  - associations
  - detailedAssociations

## Attributes Reference

The following attributes are exported:

* `key`: The key of a category when it is represented in key:value format.
* `value`: The value of a category when it is represented in key:value format
* `type`: Denotes the type of a category.
There are three types of categories: SYSTEM, INTERNAL, and USER.
* `description`: A string consisting of the description of the category as defined by the user.
* `owner_uuid`: This field contains the UUID of a user who owns the category.
* `associations`: This field gives basic information about resources that are associated to the category.
* `detailed_associations`: This field gives detailed information about resources that are associated to the category.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.


### associations
* `category_id`: External identifier for the given category, used across all v4 apis/entities/resources where categories are referenced.
* `resource_type`: An enum denoting the associated resource types. Resource types are further grouped into 2 types - entity or a policy.
* `resource_group`: An enum denoting the resource group.
Resources can be organized into either an entity or a policy.
* `count`: Count of associations of a particular type of entity or policy

### detailed_associations
* `category_id`: External identifier for the given category, used across all v4 apis/entities/resources where categories are referenced.
* `resource_type`: An enum denoting the associated resource types. Resource types are further grouped into 2 types - entity or a policy.
* `resource_group`: An enum denoting the resource group.
Resources can be organized into either an entity or a policy.
* `resource_id`: The UUID of the entity or policy associated with the particular category.


See detailed information in [Nutanix Get Category v4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/Categories/operation/getCategoryById).
