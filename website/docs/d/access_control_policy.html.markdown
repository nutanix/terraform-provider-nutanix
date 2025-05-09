---
layout: "nutanix"
page_title: "NUTANIX: nutanix_access_control_policy"
sidebar_current: "docs-nutanix-resource-access_control_policy"
description: |-
  Describes an Access Control Policy
---

# nutanix_access_control_policy

Describes an Access Control Policy.

## Example Usage

```hcl
resource "nutanix_access_control_policy" "test" {
	name        = "NAME OF ACCESS CONTROL POLICY"
	description = "DESCRIPTION OF THE ACCESS CONTROL POLICY"
	role_reference {
		kind = "role"
		uuid = "UUID of role"
	}
}
data "nutanix_access_control_policy" "test" {
    access_control_policy_id = nutanix_access_control_policy.test.id
}

```

## Argument Reference

The following arguments are supported:

* `access_control_policy_id`: - (Required) The UUID of an access control policy.

## Attribute Reference

The following attributes are exported:

* `name`: - Name of the Access Control Policy.
* `description`: - The description of the Access Control Policy.
* `categories`: - Categories for the Access Control Policy.
* `project_reference`: - The reference to a project.
* `owner_reference`: - The reference to a user.
* `role_reference`: - The reference to a role.
* `user_reference_list`: - The User(s) being assigned a given role.
* `user_group_reference_list`: - The User group(s) being assigned a given role.
* `filter_list`: - The list of filters, which define the entities.

### Filter List

The filter_list attribute supports the following:

* `context_list`: - The list of context filters. These are OR filters. The scope-expression-list defines the context, and the filter works in conjunction with the entity-expression-list. NOTE: - the absence of a scope expression in a filter implies global context.

### Context List

The context_list attribute supports the following:

* `scope_filter_expression_list`: - The device ID which is used to uniquely identify this particular disk.
* `entity_filter_expression_list` - A list of Entity filter expressions.

### Scope Filter Expression List

The scope_filter_expression_list attribute supports the following.

* `left_hand_side`: -  The LHS of the filter expression - the scope type.
* `operator`: - The operator of the filter expression.
* `right_hand_side`: - The right hand side (RHS) of an scope expression.


### Entity Filter Expression List

The scope_filter_expression_list attribute supports the following.

* `left_hand_side_entity_type`: -  The LHS of the filter expression - the entity type.
* `operator`: - The operator in the filter expression.
* `right_hand_side`: - The right hand side (RHS) of an scope expression.

### Right Hand Side

The right_hand_side attribute supports the following.

* `collection`: -  A representative term for supported groupings of entities. ALL = All the entities of a given kind.
* `categories`: - The category values represented as a dictionary of key -> list of values.
* `uuid_list`: - The explicit list of UUIDs for the given kind.

## Attributes Reference

The following attributes are exported:

* `metadata`: - The Access Control Policy kind metadata.
* `api_version` - The version of the API.
* `state`: - The state of the Access Control Policy.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when Access Control Policy was last updated.
* `uuid`: - Access Control Policy UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when Access Control Policy was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - Access Control Policy name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Reference

The `project_reference`, `owner_reference`, `role_reference` attributes supports the following:

* `kind`: - The kind name (Default value: project)(Required).
* `name`: - the name(Optional).
* `uuid`: - the UUID(Required).

For `user_reference_list` and `user_group_reference_list` are the same as reference but used as array.

See detailed information in [Nutanix Access Control Policies](https://www.nutanix.dev/api_references/prism-central-v3/#/15eefbae9e620-get-a-existing-access-control-policy).
