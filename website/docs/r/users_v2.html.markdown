---
layout: "nutanix"
page_title: "NUTANIX: nutanix_users_v2"
sidebar_current: "docs-nutanix-resource-users-v2"
description: |-
  Create a User.
---

# nutanix_users_v2

Provides Nutanix resource to Create a User.

## Example Usage

```hcl
resource "nutanix_users_v2" "user"{
  username = "<username>"
  user_type = "<user_type>"
  idp_id = "<idp_id>"
  display_name = "<display_name>"
  locale = "<locale>"
  region = "<region>"
  force_reset_password = <force_reset_password>
  status = "<status>"
}

# user_type of SERVICE_ACCOUNT
resource "nutanix_users_v2" "service_account" {
	username = "service_account_terraform_example"
	description = "service account tf"
	email_id = "example_terraform_plugin@domain.com"
	user_type = "SERVICE_ACCOUNT"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: -(Optional) External Identifier of the User.
- `username`: -(Required) Identifier for the User in the form an email address.
- `user_type`: -(Required) Enum: `$UNKNOWN` `$REDACTED` `LOCAL` `SAML` `LDAP` `EXTERNAL` `SERVICE_ACCOUNT`
  Type of the User.
- `idp_id`: -(Optional) Identifier of the IDP for the User.
- `display_name`: -(Optional) Display name for the User.
- `first_name`: -(Optional) First name for the User.
- `middle_initial`: -(Optional) Middle name for the User.
- `last_name`: -(Optional) Last name for the User.
- `email_id`: -(Optional) Email Id for the User.
- `locale`: -(Optional) Default locale for the User.
- `region`: -(Optional) Default Region for the User.
- `password`: -(Optional) Password for the User.
- `is_force_reset_password`: -(Optional) Flag to force the User to reset password.
- `additional_attributes`: -(Optional) Any additional attribute for the User.
- `status`: -(Optional) Status of the User. `ACTIVE`: Denotes that the local User is active. `INACTIVE`: Denotes that the local User is inactive and needs to be reactivated.
- `description`: -( Optional ) Description of the user.

### Additional Attributes

The additional_attributes attribute supports the following:

- `name`: -(Optional) The URL at which the entity described by the link can be accessed.
- `value`: -(Optional) A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

## Attributes Reference

The following attributes are exported:

- `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id` - The External Identifier of the User Group.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `username`: - Identifier for the User in the form an email address.
- `user_type`: - Enum: `$UNKNOWN` `$REDACTED` `LOCAL` `SAML` `LDAP` `EXTERNAL`
  Type of the User.
- `idp_id`: - Identifier of the IDP for the User.
- `display_name`: - Display name for the User.
- `first_name`: - First name for the User.
- `middle_initial`: - Middle name for the User.
- `last_name`: - Last name for the User.
- `email_id`: - Email Id for the User.
- `locale`: - Default locale for the User.
- `region`: - Default Region for the User.
- `is_force_reset_password`: - Flag to force the User to reset password.
- `additional_attributes`: - Any additional attribute for the User.
- `status`: - Status of the User. `ACTIVE`: Denotes that the local User is active. `INACTIVE`: Denotes that the local User is inactive and needs to be reactivated.
- `description`: - Description of the user.
- `buckets_access_keys`: - Bucket Access Keys for the User.
- `last_login_time`: - Last successful logged in time for the User.
- `created_time`: - Creation time of the User.
- `last_updated_time`: - Last updated time of the User.
- `created_by`: - User or Service who created the User.
- `last_updated_by`: - Last updated by this User ID.

### Links

The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Additional Attributes

The additional_attributes attribute supports the following:

- `name`: - The URL at which the entity described by the link can be accessed.
- `value`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Buckets Access Keys

The buckets_access_keys attribute supports the following:

- `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `access_key_name`: - Name of the Bucket Access Key.
- `secret_access_key`: - Secret Access Key, it will be returned only during Bucket Access Key creation.
- `user_id`: - User Identifier who owns the Bucket Access Key.
- `created_time`: - Creation time for the Bucket Access Key.

## Import

This helps to manage existing entities which are not created through terraform. Users can be imported using the `UUID`. (ext_id in v4 terms). eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_users_v2" "import_user" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_users_v2" "fetch_users"{}
terraform import nutanix_users_v2.import_user <UUID>
```

## Delete Behavior

Deleting users via Terraform is **not currently supported** in `nutanix_users_v2`.
If you attempt to delete a resource, Terraform will issue a **warning** and continue without removing the actual user in Prism Central.

To delete a user, you must remove it manually via the **Prism Central UI**, or as a workaround, you can:

1. **Import the user into the `nutanix_user` resource** (v3 API).
2. Use Terraform to delete the user, as delete operations are supported in the v3 API.

## References

See detailed information in [Nutanix Users v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Users/operation/createUser).
