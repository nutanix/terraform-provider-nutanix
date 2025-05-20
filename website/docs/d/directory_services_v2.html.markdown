---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_services_v2"
sidebar_current: "docs-nutanix-datasource-nutanix-directory-services-v2"
description: |-
    This operation retrieves a list of all Directory Service(s).
---

# nutanix_pbr

Provides a datasource to retrieve all Directory Service(s).

## Example Usage

```hcl
data "nutanix_directory_services_v2" "example"{}
```

## Argument Reference
The following arguments are supported:


* `page`: -(Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: -(Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: -(Optional) A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions. For example, filter '$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
    - createdBy
    - domainName
    - extId
    - name
* `order_by`: -(Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
    - createdTime
    - domainName
    - lastUpdatedTime
    - name
* `select`: -(Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. it can be applied to the following fields:
    - createdBy
    - createdTime
    - directoryType
    - domainName
    - extId
    - groupSearchType
    - lastUpdatedTime
    - links
    - name
    - openLdapConfiguration/userConfiguration
    - openLdapConfiguration/userGroupConfiguration
    - secondaryUrls
    - serviceAccount/password
    - serviceAccount/username
    - tenantId
    - url
    - whiteListedGroups




## Attributes Reference
The following attributes are exported:

* `directory_services`: - list of all Directory Service(s).


### Directory Services
The directory_services attribute supports the following:


* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`: - Name for the Directory Service.
* `url`: - URL for the Directory Service.
* `secondary_urls`: - Secondary URL for the Directory Service.
* `domain_name`: - Domain name for the Directory Service.
* `directory_type`: - Type of Directory Service, Supported values are: "ACTIVE_DIRECTORY" (Directory Service type is Active Directory.) and "OPEN_LDAP" (Directory Service type is Open LDAP.)
* `service_account`: - Information of Service account to connect to the Directory Service.
* `open_ldap_configuration`: - Configuration for OpenLDAP Directory Service.
* `group_search_type`: - Group membership search type for the Directory Service. Supported values are: "NON_RECURSIVE" (Doesn't search recursively within groups.) and "RECURSIVE" (Searches recursively within groups.)
* `white_listed_groups`: - List of allowed User Groups for the Directory Service.
* `created_time`: - Creation time of the Directory Service.
* `last_updated_time`: - Last updated time of the Directory Service.
* `created_by`: - User or Service who created the Directory Service.

#### Service Account

The service_account attribute supports the following:

* `username`: - Username to connect to the Directory Service.
* `password`: - Password to connect to the Directory Service.


#### Open Ldap Configuration

The open_ldap_configuration attribute supports the following:

* `user_configuration`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.
* `user_group_configuration`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.

##### User Configuration

The user_configuration attribute supports the following:

* `user_object_class`: - Object class in the OpenLDAP system that corresponds to Users.
* `user_search_base`: - Base DN for User search.
* `username_attribute`: - Unique Identifier for each User which can be used in Authentication.

##### User Group Configuration

The user_group_configuration attribute supports the following:

* `group_object_class`: - Object class in the OpenLDAP system that corresponds to groups.
* `group_search_base`: - Base DN for group search.
* `group_member_attribute`: - Attribute in a group that associates Users to the group.
* `group_member_attribute_value`: - User attribute value that will be used in group entity to associate User to the group.


See detailed information in [Nutanix List Directory Services v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/DirectoryServices/operation/listDirectoryServices).
