---
layout: "nutanix"
page_title: "NUTANIX: nutanix_template_versions_v4"
sidebar_current: "docs-nutanix-datasource-template-versions-v4"
description: |-
 List Versions with details like name, description, VM configuration, etc.
---

# nutanix_template_versions_v4

List Versions with details like name, description, VM configuration, etc. This operation supports filtering, sorting & pagination.

## Example

```hcl
    data "nutanix_template_version_v4" "test" { 
        ext_id = {{ template uuid }}
    }
```

## Argument Reference

The following arguments are supported:
* `ext_id`: (Required) The identifier of a Template.

* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results. Default is 0.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions.
* `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default.

* `template_versions`: List of template versions

## template_versions Reference

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `version_name`: The user defined name of a Template Version.
* `version_description`: The user defined description of a Template Version.
* `vm_spec`: VM configuration details.
* `create_time`: Time when the Template was created.
* `created_by`: Information of the User.
* `is_active_version`: Specify whether to mark the Template Version as active or not. The newly created Version during Template Creation, Updation or Guest OS Updation is set to Active by default unless specified otherwise.
* `is_gc_override_enabled`: Allow or disallow override of the Guest Customization during Template deployment.


### created_by

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `username`: Identifier for the User in the form an email address.
* `user_type`: Type of the User.
* `idp_id`: Identifier of the IDP for the User.
* `display_name`: Display name for the User.
* `first_name`: First name for the User.
* `middle_initial`: Middle name for the User.
* `last_name`: Last name for the User.
* `email_id`: Email Id for the User.
* `locale`: Default locale for the User.
* `region`: Default Region for the User.
* `is_force_reset_password_enabled`: Flag to force the User to reset password.
* `additional_attributes`: Any additional attribute for the User.
* `status`: Status of the User.
* `buckets_access_keys`: Bucket Access Keys for the User.
* `last_login_time`: Last successful logged in time for the User.
* `created_time`: Creation time of the User.
* `last_updated_time`: Last updated time of the User.
* `created_by`: User or Service who created the User.
* `last_updated_by`: Last updated by this User ID.


See detailed information in [Nutanix Template](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).