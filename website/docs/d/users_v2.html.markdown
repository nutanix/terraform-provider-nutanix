---
layout: "nutanix"
page_title: "NUTANIX: nutanix_users_v2"
sidebar_current: "docs-nutanix-datasource-users-v2"
description: |-
  Provides a datasource to retrieve all User(s).
---

# nutanix_users_v2

Provides a datasource to retrieve all User(s).

## Example Usage

```hcl
# list all users
data "nutanix_users_v2" "list-users"{}


data "nutanix_users_v2" "filtered-users" {
  filter = "username eq 'username-example'"
}

```

##  Argument Reference

The following arguments are supported:

* `page`: - A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` :A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
    * createdBy
    * displayName
    * emailId
    * extId
    * firstName
    * idpId
    * lastName
    * lastUpdatedBy
    * status
    * userType
    * username
* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:     * createdBy
    * createdTime
    * displayName
    * emailId
    * extId
    * firstName
    * lastLoginTime
    * lastName
    * lastUpdatedTime
    * userType
    * username
* `select` : A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. following fields:
    * additionalAttributes
    * bucketsAccessKeys
    * createdBy
    * createdTime
    * displayName
    * emailId
    * extId
    * firstName
    * idpId
    * isForceResetPasswordEnabled
    * lastLoginTime
    * lastName
    * lastUpdatedBy
    * lastUpdatedTime
    * links
    * locale
    * middleInitial
    * region
    * status
    * tenantId
    * userType
    * username

## Attributes Reference
The following attributes are exported:

* `users` : List all User(s).

### User Groups

The `users`  attribute element contains the following attributes:

* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id` - The External Identifier of the User Group.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `username`: - Identifier for the User in the form an email address.
* `user_type`: - Enum: `$UNKNOWN` `$REDACTED` `LOCAL` `SAML` `LDAP` `EXTERNAL`
Type of the User.
* `idp_id`: - Identifier of the IDP for the User.
* `display_name`: - Display name for the User.
* `first_name`: - First name for the User.
* `middle_initial`: - Middle name for the User.
* `last_name`: - Last name for the User.
* `email_id`: - Email Id for the User.
* `locale`: - Default locale for the User.
* `region`: - Default Region for the User.
* `is_force_reset_password`: - Flag to force the User to reset password.
* `additional_attributes`: -  Any additional attribute for the User.
* `status`: - Status of the User. `ACTIVE`: Denotes that the local User is active. `INACTIVE`: Denotes that the local User is inactive and needs to be reactivated.
* `buckets_access_keys`: - Bucket Access Keys for the User.
* `last_login_time`: - Last successful logged in time for the User.
* `created_time`: - Creation time of the User.
* `last_updated_time`: - Last updated time of the User.
* `created_by`: - User or Service who created the User.
* `last_updated_by`: - Last updated by this User ID.


### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.


### Additional Attributes

The additional_attributes attribute supports the following:

* `name`: - The URL at which the entity described by the link can be accessed.
* `value`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Buckets Access Keys

The buckets_access_keys attribute supports the following:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `access_key_name`: - Name of the Bucket Access Key.
* `secret_access_key`: - Secret Access Key, it will be returned only during Bucket Access Key creation.
* `user_id`: - User Identifier who owns the Bucket Access Key.
* `created_time`: - Creation time for the Bucket Access Key.


See detailed information in [Nutanix Users v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Users/operation/listUsers).
