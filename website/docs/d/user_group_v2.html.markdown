---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_group_v2"
sidebar_current: "docs-nutanix-datasource-user-group-v2"
description: |-
  This operation retrieves a User Group based on the External Identifier of the User Group.
---

# nutanix_user_group_v2

Provides a datasource to retrieve a user group based on the External Identifier of the User Group.

## Example Usage

```hcl

data "nutanix_user_group_v2" "get-ug"{
	ext_id = "a2a8650a-358a-4791-90c9-7a8b6e2989d6"
}

```




## Argument Reference

The following arguments are supported:

* `ext_id`: - (Required) The External Identifier of the User Group.

## Attributes Reference

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

### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.


See detailed information in [Nutanix Get User Group v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/UserGroups/operation/getUserGroupById).
