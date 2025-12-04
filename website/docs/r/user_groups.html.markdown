---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_groups"
sidebar_current: "docs-nutanix-resource-user-groups"
description: |-
  This operation add a User group to the system.
---

# nutanix_user_groups

Provides a resource to add a User group to the system..

## Example Usage

```hcl
resource "nutanix_user_groups" "user_grp" {
	directory_service_user_group{
		distinguished_name = "<distinguished name for the user group>"
	}
}
```


```hcl
resource "nutanix_user_groups" "user_grp" {
	saml_user_group{
    name = "<name of saml group>"
    idp_uuid = "<idp uuid of the group>"
  }
}
```

## Argument Reference

The following arguments are supported:

* `directory_service_user_group`: - (Optional) A Directory Service user group.
* `directory_service_ou`: - (Optional) A Directory Service organizational unit.
* `saml_user_group`: - (Optional) A SAML Service user group.

### directory_service_user_group , directory_service_ou

A Directory Service user group supports the following.

* `distinguished_name`: - (Required) The Distinguished name for the user group. 

### saml_user_group

A SAML Service user group supports the following.

* `name` :- (Required) The name of the SAML group which the IDP provides. 
* `idp_uuid` :- (Required) The UUID of the Identity Provider that the group belongs to. 


## Attributes Reference

The following attributes are exported:

* `metadata` - The user_group kind metadata.
* `api_version` - The version of the API.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: - subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

See detailed information in [Nutanix User Groups](https://www.nutanix.dev/api_references/prism-central-v3/#/2fb233cea33f8-add-a-user-group).
