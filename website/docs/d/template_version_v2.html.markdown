---
layout: "nutanix"
page_title: "NUTANIX: nutanix_template_version_v4"
sidebar_current: "docs-nutanix-datasource-template-version-v4"
description: |-
 Retrieve the Template Version details for the given Template Version identifier.
---

# nutanix_template_version_v4

Retrieve the Template Version details for the given Template Version identifier.

## Example

```hcl
    data "nutanix_template_version_v4" "test" { 
        template_ext_id = {{ template uuid }}
        ext_id = {{ template version uuid }}
    }
```

## Argument Reference

The following arguments are supported:
* `template_ext_id`: (Required) The identifier of a Template.
* `ext_id`: (Required) The identifier of a Template Version.


## Attribute Reference

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