---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_groups_v2"
sidebar_current: "docs-nutanix-resource-user-groups-v2"
description: |-
  This operation add a User group to the system.
---

# nutanix_user_groups_v2

Provides a resource to add a User group to the system..

## Example Usage

```hcl
resource "nutanix_user_groups_v2" "usr_group"{
  group_type         = "LDAP"
  idp_id             = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
  name               = "group_0664229e"
  # IAM stores distinguished name in lowerCase, even user send them as upperCase. Suggest user to use lowerCase letters.
  distinguished_name = "cn=group_0664229e,ou=group,dc=devtest,dc=local"
}

# Saml User group
resource "nutanix_user_groups_v2" "saml-ug" {
  group_type = "SAML"
  idp_id     = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
  name       = "adfs19admingroup"
}

```

## Argument Reference

The following arguments are supported:

- `ext_id` -(Optional) The External Identifier of the User Group.
- `group_type`: -(Required) Type of the User Group. LDAP (User Group belonging to a Directory Service (Open LDAP/AD)), SAML (User Group belonging to a SAML IDP.)
- `idp_id`: -(Required) Identifier of the IDP for the User Group.
- `name`: -(Optional) Common Name of the User Group.
- `distinguished_name`: -(Optional) Identifier for the User Group in the form of a distinguished name.

## Attributes Reference

The following attributes are exported:

- `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id` - The External Identifier of the User Group.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `group_type`: - Type of the User Group. LDAP (User Group belonging to a Directory Service (Open LDAP/AD)), SAML (User Group belonging to a SAML IDP.)
- `idp_id`: - Identifier of the IDP for the User Group.
- `name`: - Common Name of the User Group.
- `distinguished_name`: - Identifier for the User Group in the form of a distinguished name.
- `created_time`: - Creation time of the User Group.
- `last_updated_time`: - Last updated time of the User Group.
- `created_by`: - User or Service who created the User Group.

### Links

The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

## Import

This helps to manage existing entities which are not created through terraform. User Group can be imported using the `UUID`. (ext_id in v4 API context). eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_user_groups_v2" "import_ug" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_user_groups_v2" "fetch_ugs"{}
terraform import nutanix_user_groups_v2.import_ug <UUID>
```

See detailed information in [Nutanix Create User Group v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/UserGroups/operation/createUserGroup).
