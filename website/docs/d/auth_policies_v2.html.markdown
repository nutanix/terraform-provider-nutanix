---
layout: "nutanix"
page_title: "NUTANIX: nutanix_authorization_policies_v2"
sidebar_current: "docs-nutanix-datasource-authorization-policies-v2"
description: |-
  This operation retrieves the list of existing Authorization Policies.
---

# nutanix_authorization_policies_v2

Provides a datasource to retrieve the list of existing Authorization Policies.

## Example Usage

```hcl
#list of authorization policies, with limit and filter
data "nutanix_authorization_policies_v2" "filtered-ap"{
  filter = "displayName eq 'auth_policy_example'"
  limit  = 2
}

# list of authorization policies, with select
data "nutanix_authorization_policies_v2" "select-ap"{
  select     = "extId,displayName,description,authorizationPolicyType"
}
```

## Argument Reference

The following arguments are supported:

- `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources. The filter can be applied to the following fields:
  - authorizationPolicyType
  - clientName
  - createdBy
  - createdTime
  - displayName
  - extId
  - isSystemDefined
  - lastUpdatedTime
  - role
- `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
  - createdTime
  - displayName
  - extId
  - lastUpdatedTime
  - role
- `expand`: (Optional) A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are \$filter, \$select and \$orderby. The following expansion keys are supported:
  - role
- `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. The select can be applied to the following fields:
  - authorizationPolicyType
  - authorizationPolicyType
  - clientName
  - createdBy
  - createdTime
  - description
  - displayName
  - entities
  - extId
  - identities
  - isSystemDefined
  - lastUpdatedTime
  - links
  - role
  - tenantId

## Attribute Reference

The following attributes are exported:

- `auth_policies`: List of all existing Authorization Policies.

## Authorization Policies

The following attributes are exported for each Authorization Policy:

- `ext_id`: ext_id of Authorization policy.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `display_name`: Name of the Authorization Policy.
- `description`: Description of the Authorization Policy.
- `client_name`: Client that created the entity.
- `identities`: The identities for which the Authorization Policy is created.
- `entities`: The entities being qualified by the Authorization Policy.
- `role`: The Role associated with the Authorization Policy.
- `created_time`: The creation time of the Authorization Policy.
- `last_updated_time`: The time when the Authorization Policy was last updated.
- `created_by`: User or Service Name that created the Authorization Policy.
- `is_system_defined`: Flag identifying if the Authorization Policy is system defined or not.
- `authorization_policy_type`: Type of Authorization Policy.
  - `PREDEFINED_READ_ONLY` : System-defined read-only ACP, i.e. no modifications allowed.
  - `SERVICE_DEFINED_READ_ONLY` : Read-only ACP defined by a service.
  - `PREDEFINED_UPDATE_IDENTITY_ONLY` : System-defined ACP prohibiting any modifications from customer.
  - `SERVICE_DEFINED` : ACP defined by a service.
  - `USER_DEFINED` : ACP defined by an User.

### Links

The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object

See detailed information in [Nutanix List Authorization Policies v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/AuthorizationPolicies/operation/listAuthorizationPolicies).
