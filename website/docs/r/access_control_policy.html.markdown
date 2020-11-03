---
layout: "nutanix"
page_title: "NUTANIX: nutanix_access_control_policy"
sidebar_current: "docs-nutanix-resource-access_control_policy"
description: |-
  This operation submits a request to create an access control policy based on the input parameters.
---

# nutanix_access_control_policy

Provides a resource to create an access control policy based on the input parameters.

## Example Usage

``` hcl
resource "nutanix_access_control_policy" "test" {
	name        = "NAME OF ACCESS CONTROL POLICY"
	description = "DESCRIPTION OF THE ACCESS CONTROL POLICY"
	role_reference {
		kind = "role"
		uuid = "UUID of role"
	}
	user_reference_list{
		uuid = "UUID of User existent"
		name = "admin"
	}

	context_filter_list{
        entity_filter_expression_list{
            operator = "IN"
            left_hand_side_entity_type = "cluster"
            right_hand_side{
                uuid_list = ["00058ef8-c31c-f0bc-0000-000000007b23"]
            }
        }
        entity_filter_expression_list{
            operator = "IN"
            left_hand_side_entity_type = "image"
            right_hand_side{
                collection = "ALL"
            }
        }
        entity_filter_expression_list{
            operator = "IN"
            left_hand_side_entity_type = "category"
            right_hand_side{
                collection = "ALL"
            }
        }
        entity_filter_expression_list{
            operator = "IN"
            left_hand_side_entity_type = "marketplace_item"
            right_hand_side{
                collection = "SELF_OWNED"
            }
        }
        entity_filter_expression_list{
            operator = "IN"
            left_hand_side_entity_type = "app_task"
            right_hand_side{
                collection = "SELF_OWNED"
            }
        }
        entity_filter_expression_list{
            operator = "IN"
            left_hand_side_entity_type = "app_variable"
            right_hand_side{
                collection = "SELF_OWNED"
            }
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Optional) Name of the Access Control Policy.
* `description`: - (Optional) The description of the association of a role to a user in a given context.
* `categories`: - (Optional) Categories for the Access Control Policy.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.
* `project_reference`: - (Optional) The reference to a project.
* `role_reference`: - (Required) The reference to a role.
* `user_reference_list`: - (Optional) The User(s) being assigned a given role.
* `user_group_reference_list`: - (Optional) The User group(s) being assigned a given role.
* `filter_list`: - (Optional) The list of filters, which define the entities.

### Filter List

The filter_list attribute supports the following:

* `context_list`: - (Optional) The list of context filters. These are OR filters. The scope-expression-list defines the context, and the filter works in conjunction with the entity-expression-list. NOTE: - the absence of a scope expression in a filter implies global context.

### Context List

The context_list attribute supports the following:

* `scope_filter_expression_list`: - (Optional) The device ID which is used to uniquely identify this particular disk.
* `entity_filter_expression_list` - (Required) A list of Entity filter expressions.

### Scope Filter Expression List

The scope_filter_expression_list attribute supports the following.

* `left_hand_side`: - (Optional)  The LHS of the filter expression - the scope type.
* `operator`: - (Required) The operator of the filter expression.
* `right_hand_side`: - (Required) The right hand side (RHS) of an scope expression.


### Entity Filter Expression List

The scope_filter_expression_list attribute supports the following.

* `left_hand_side_entity_type`: - (Optional)  The LHS of the filter expression - the entity type.
* `operator`: - (Required) The operator in the filter expression.
* `right_hand_side`: - (Required) The right hand side (RHS) of an scope expression.

### Right Hand Side

The right_hand_side attribute supports the following.

* `collection`: - (Optional)  A representative term for supported groupings of entities. ALL = All the entities of a given kind.
* `categories`: - (Optional) The category values represented as a dictionary of key -> list of values.
* `uuid_list`: - (Optional) The explicit list of UUIDs for the given kind.

## Attributes Reference

The following attributes are exported:

* `metadata`: - The vm kind metadata.
* `api_version` - The version of the API.
* `state`: - The state of the vm.
* `cluster_name`: - The name of the cluster.
* `host_reference`: - Reference to a host.
* `hypervisor_type`: - The hypervisor type for the hypervisor the VM is hosted on.
* `nic_list_status`: - Status NICs attached to the VM.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when vm was last updated.
* `uuid`: - vm UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when vm was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - vm name.
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

See detailed information in [Nutanix Access Control Policies](https://www.nutanix.dev/reference/prism_central/v3/api/access-control-policies/).
