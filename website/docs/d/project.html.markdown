---
layout: "nutanix"
page_title: "NUTANIX: nutanix_project"
sidebar_current: "docs-nutanix-datasource-project"
description: |-
  Describe a Nutanix Project and its values (if it has them).
---

# nutanix_project

Describe a Nutanix Project and its values (if it has them).

## Example Usage

```hcl
resource "nutanix_subnet" "subnet" {
  cluster_uuid       = "<YOUR_CLUSTER_ID>"
  name               = "sunet_test_name"
  description        = "Description of my unit test VLAN"
  vlan_id            = 31
  subnet_type        = "VLAN"
  subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length      = 24

  dhcp_options = {
    boot_file_name   = "bootfile"
    domain_name      = "nutanix"
    tftp_server_name = "10.250.140.200"
  }

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

resource "nutanix_project" "project_test" {
  name        = "my-project"
  description = "This is my project"

  categories {
    name  = "Environment"
    value = "Staging"
  }

  default_subnet_reference {
    uuid = nutanix_subnet.subnet.metadata.uuid
  }

  api_version = "3.1"
}

data "nutanix_project" "test" {
    project_id = nutanix_project.project_test.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id`: - (Required) The `id` of the project.

## Attributes Reference

The following attributes are exported:

* `name` The name for the project.
* `description` A description for project.

### Resource Domain
* `resource_domain` (Deprecated) Not supported starting from provider version `2.4.0` and expected to be empty. Remove any usage from configuration/scripts.

### Account Reference List
* `account_reference_list` - List of accounts associated with the project.
* `account_reference_list.#.kind` - The kind name. Default value is `account`
* `account_reference_list.#.uuid` - The UUID of an account.
* `account_reference_list.#.name` - The name of an account.

### Environment Reference List
* `environment_reference_list` - List of environments associated with the project.
* `environment_reference_list.#.kind` - The kind name. Default value is `environment`
* `environment_reference_list.#.uuid` - The UUID of an environment.
* `environment_reference_list.#.name` - The name of an environment.

### Default Subnet Reference Map
* `default_subnet_reference` - Reference to a subnet.
* `default_subnet_reference.kind` - The kind name. Default value is `subnet`
* `default_subnet_reference.uuid` - The UUID of a subnet.
* `default_subnet_reference.name` - The name of a subnet.

### user_reference_list
* `user_reference_list` - List of users in the project.
* `user_reference_list.#.kind` - The kind name. Default value is `user`
* `user_reference_list.#.uuid` - The UUID of a user
* `user_reference_list.#.name` - The name of a user.

### External User Group Reference List
* `external_user_group_reference_list` - List of directory service user groups. These groups are not managed by Nutanix.
* `external_user_group_reference_list.#.kind` - The kind name. Default value is `user_group`
* `external_user_group_reference_list.#.uuid` - The UUID of a user_group
* `external_user_group_reference_list.#.name` - The name of a user_group

### Subnet Reference List
* `subnet_reference_list` - List of subnets for the project.
* `subnet_reference_list.#.kind` - The kind name. Default value is `subnet`
* `subnet_reference_list.#.uuid` - The UUID of a subnet
* `subnet_reference_list.#.name` - The name of a subnet.

### External Network List
* `external_network_list` - List of external networks associated with the project.
* `external_network_list.#.uuid` - The UUID of a network.
* `external_network_list.#.name` - The name of a network.

### Tunnel Reference List
* `tunnel_reference_list` - (Optional/Computed) List of tunnels associated with the project.
* `tunnel_reference_list.#.kind` - (Optional) The kind name. Default value is `tunnel`
* `tunnel_reference_list.#.uuid` - (Required) The UUID of a tunnel
* `tunnel_reference_list.#.name` - (Optional/Computed) The name of a tunnel.

### Cluster Reference List
* `cluster_reference_list` - (Optional/Computed) List of clusters associated with the project..
* `cluster_reference_list.#.kind` - (Optional) The kind name. Default value is `cluster`
* `cluster_reference_list.#.uuid` - (Required) The UUID of a cluster
* `cluster_reference_list.#.name` - (Optional/Computed) The name of a cluster.

### VPC Reference List
* `vpc_reference_list` - (Optional/Computed) List of VPCs associated with the project..
* `vpc_reference_list.#.kind` - (Optional) The kind name. Default value is `vpc`
* `vpc_reference_list.#.uuid` - (Required) The UUID of a vpc
* `vpc_reference_list.#.name` - (Optional/Computed) The name of a vpc.

### Default Environment Reference Map
* `default_environment_reference` - (Optional/Computed) Reference to a environment.
* `default_environment_reference.kind` - (Optional) The kind name. Default value is `environment`
* `default_environment_reference.uuid` - (Required) The UUID of a environment
* `default_environment_reference.name` - (Optional/Computed) The name of a environment.

### ACP
ACPs will be exported if use_project_internal flag is set.
* `name` - Name of ACP
* `description` - Description of ACP
* `user_reference_list` - List of Reference of users.
* `user_group_reference_list` - List of Reference of users groups.
* `role_reference` - Reference to role.
* `context_filter_list` - The list of context filters. These are OR filters. The scope-expression-list defines the context, and the filter works in conjunction with the entity-expression-list.

The context_list attribute supports the following:

* `scope_filter_expression_list`: - (Optional) Filter the scope of an Access Control Policy.
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

### Metadata
The metadata attribute exports the following:

* `last_update_time` - UTC date and time in RFC-3339 format when the project was last updated.
* `uuid` - Project UUID.
* `creation_time` - UTC date and time in RFC-3339 format when the project was created.
* `spec_version` - Version number of the latest spec.
* `spec_hash` - Hash of the spec. This will be returned from server.
* `name` - Project name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories
The categories attribute supports the following:

* `name` - the key name.
* `value` - value of the key.

### Reference
The `project_reference`, `owner_reference` attributes supports the following:

* `kind` - (Required) The kind name (Default value: `project`).
* `name` - the name.
* `uuid` - (Required) the UUID.

See detailed information in [Nutanix Project](https://www.nutanix.dev/api_references/prism-central-v3/#/81f93ec6d7685-get-a-existing-project).
