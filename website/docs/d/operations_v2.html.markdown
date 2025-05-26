---
layout: "nutanix"
page_title: "NUTANIX: nutanix_operations_v2"
sidebar_current: "docs-nutanix-datasource-operations-v2"
description: |-
  This operation retrieves a list of all the operations.
---

# nutanix_operations_v2
Lists the operations defined on the system. List of operations can be further filtered out using various filtering options.

## Example

```hcl
#list operations
data "nutanix_operations_v2" "operation-list" {}

# filtered list operation
data "nutanix_operations_v2" "operation-list-filtered" {
  filter = "displayName eq 'Create_Role'"
}

# list operations withe page and limit
data "nutanix_operations_v2" "operation-list-paginated" {
  page  = 1
  limit = 10
}

```

## Attribute Reference

The following attributes are exported:

* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. The filter can be applied to the following fields:
    - clientName
    - createdTime
    - displayName
    - entityType
    - extId
    - lastUpdatedTime
    - operationType

* `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
    - createdTime
    - displayName
    - entityType
    - extId
    - lastUpdatedTime
* `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. it can be applied to the following fields:
    - associatedEndpointList
    - clientName
    - createdTime
    - description
    - displayName
    - entityType
    - extId
    - lastUpdatedTime
    - links
    - operationType
    - relatedOperationList
    - tenantId


## Attributes Reference
The following attributes are exported:

* `operations`: List of all operations


### operations
The `operations` attribute supports the following:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `display_name`: Permission name.
* `description`: Permission description
* `created_time`: Permission creation time
* `last_updated_time`: Permission last updated time.
* `entity_type`: Type of entity associated with this Operation.
* `operation_type`: The Operation type. Currently we support INTERNAL, EXTERNAL and SYSTEM_DEFINED_ONLY.
* `client_name`: Client that created the entity.
* `related_operation_list`: List of related Operations. These are the Operations which might need to be given access to, along with the current Operation, for certain workflows to succeed.
* `associated_endpoint_list`: List of associated endpoint objects for the Operation.

### associated_endpoint_list
* `api_version`: Version of the API for the provided associated endpoint.
* `endpoint_url`: Endpoint URL.
* `http_method`: HTTP method for the provided associated endpoint.

See detailed information in [Nutanix List Operations](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Operations/operation/listOperations).
