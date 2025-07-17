---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_keys_v2"
sidebar_current: "docs-nutanix-datasource-user-keys-v2"
description: |-
  List all keys identified by the external identifier of a user.
---

# nutanix_user_keys_v2

List all keys identified by the external identifier of a user.

## Example Usage

```hcl
# Data source to fetch the list of keys
data "nutanix_user_keys_v2" "get_keys" {
  user_ext_id = "<SERVICE_ACCOUNT_UUID>"
}

# Data source to fetch the key by name
data "nutanix_user_keys_v2" "get_keys_filter" {
  user_ext_id = "<SERVICE_ACCOUNT_UUID>"
  filter = "name eq '<NAME_OF_API_KEY>'"
}
```

##  Argument Reference

The following arguments are supported:

* `user_ext_id`: - ( Required ) External Identifier of the User.
* `page`:- (Optional)A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`:- (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` :- (Optional) A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
    * assignedTo
    * creationType
    * extId
    * keyType
    * lastUpdatedBy
    * name
    * status
    * tenantId  
* `orderby` :- (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:     * createdBy
    * assignedTo
    * createdBy
    * createdTime
    * expiryTime
    * keyType
    * lastUpdatedTime
    * lastUsedTime
    * name
    * status
    * tenantId
* `select` :- (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. following fields:
    * assignedTo
    * createdBy
    * createdTime
    * creationType
    * description
    * expiryTime
    * extId
    * keyType
    * lastUpdatedBy
    * lastUpdatedTime
    * lastUsedTime
    * name
    * status
    * tenantId

## Attributes Reference

The following attributes are exported:


* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id` - The External Identifier of the User Group.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`: - Identifier for the key in the form of a name.
* `description`: - Brief description of the key.
* `key_type`: - The type of key.
* `created_time`: - The creation time of the key.
* `last_updated_by`: - User who updated the key.
* `creation_type`: - The creation mechanism of this entity.
* `expiry_time`: - The time when the key will expire.
* `status`: - The status of the key.
* `created_by`: - User or service who created the key.
* `last_updated_time`: - The time when the key was updated.
* `assigned_to`: - External client to whom the given key is allocated.
* `last_used_time`: - The time when the key was last used.
* `key_details`: - Details specific to type of the key.



See detailed information in [Nutanix List keys for the user V4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Users/operation/listUserKeys)