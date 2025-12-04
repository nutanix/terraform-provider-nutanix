---
layout: "nutanix"
page_title: "NUTANIX: nutanix_categories_v2"
sidebar_current: "docs-nutanix-datasource-categories-v2"
description: |-
  Fetch a list of categories with pagination, filtering, sorting, selection and optional expansion of associated entity counts.
---

# nutanix_categories_v2
List categories


## Example

```hcl

data "nutanix_categories_v2" "categories-list"{}

data "nutanix_categories_v2" "categories-paginated"{
  page = 1
  limit = 10
}

data "nutanix_categories_v2" "categories-sorted"{
  order_by = "key desc"
}

data "nutanix_categories_v2" "categories-filtered"{
  filter = "key eq 'key_example'"
}

output "category" {
  value = data.nutanix_categories_v2.categories-list.categories.0
}
```


## Argument Reference

The following arguments are supported:

* `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: (Optional)A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - extId
  - key
  - ownerUuid
  - type
  - value
* `order_by`: (Optional)A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - key
  - value
* `expand`: (Optional)A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are \$filter, \$select and \$orderby. The following expansion keys are supported:
  - associations
  - detailedAssociations
* `select`: (Optional)A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields:
  - description
  - extId
  - key
  - ownerUuid
  - type
  - value

## Attributes Reference
The following attributes are exported:

* `categories`: List of categories

## Categories
The `categories` contains list of categories. Each category has the following attributes:

* `ext_id`: The extID for the category.
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


See detailed information in [Nutanix List Categories v4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/Categories/operation/listCategories).
