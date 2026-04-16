---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_groups_v2"
sidebar_current: "docs-nutanix-datasource-user-groups-v2"
description: |-
  Provides a datasource to retrieve all the user groups.
---

# nutanix_user_groups_v2

Provides a datasource to retrieve all the user groups.

## Example Usage

```hcl
data "nutanix_user_groups_v2" "user-groups"{}
```

##  Argument Reference

The following arguments are supported:

* `page`: - A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` :A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - createdBy
  - distinguishedName
  - extId
  - groupType
  - idpId
  - name
* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - createdTime
  - distinguishedName
  - groupType
  - lastUpdatedTime
  - name
* `select` : A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned.The select can be applied to the following fields:
  - createdBy
  - createdTime
  - distinguishedName
  - extId
  - groupType
  - idpId
  - lastUpdatedTime
  - links
  - name
  - tenantId


## Attributes Reference
The following attributes are exported:

* `user_groups` : List all User Group(s).

### User Groups

The `user_groups`  attribute element contains the following attributes:

The following attributes are exported:

* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id` - The External Identifier of the User Group.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `group_type`: - Type of the User Group. LDAP (User Group belonging to a Directory Service (Open LDAP/AD)),  SAML (User Group belonging to a SAML IDP.)
* `idp_id`: - Identifier of the IDP for the User Group.
* `name`: - Common Name of the User Group.
* `distinguished_name`: - Identifier for the User Group in the form of a distinguished name.
* `created_time`: - Creation time of the User Group.
* `last_updated_time`: - Last updated time of the User Group.
* `created_by`: - User or Service who created the User Group.


#### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.


See detailed information in [Nutanix List User Groups v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/UserGroups/operation/listUserGroups).
