---
layout: "nutanix"
page_title: "NUTANIX: nutanix_subnets"
sidebar_current: "docs-nutanix-datasource-subnets"
description: |-
 Describes a list of subnets
---

# nutanix_subnets

Describes a list of subnets

## Example Usage

```hcl
data "nutanix_subnets" "subnets" {}

data "nutanix_subnets" "test" {
    metadata {
        filter = "name==vlan0_test_2"
    }
}
```

## Attribute Reference

The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of Subnets

# Entities

The entities attribute element contains the followings attributes:

The following attributes are exported:

* `metadata`: The subnet kind metadata.
* `availability_zone_reference`: The reference to a availability_zone.
* `cluster_reference`: The reference to a cluster.
* `cluster_name`: The name of a cluster.
* `description`: A description for subnet.
* `name`: Subnet name (Readonly).
* `categories`: The API Version.
* `owner_reference`: The reference to a user.
* `project_reference`: The reference to a project.
* `vswitch_name`: The name of the vswitch.
* `subnet_type`: The type of the subnet.
* `default_gateway_ip`: Default gateway IP address.
* `prefix_length`: -. IP prefix length of the Subnet.
* `subnet_ip`: Subnet IP address.
* `dhcp_server_address`: Host address.
* `dhcp_server_address_port`: Port Number.
* `dhcp_options`: Spec for defining DHCP options.
* `dhcp_domain_search_list`: DHCP domain search list for a subnet.
* `vlan_id`: VLAN assigned to the subnet.
* `network_function_chain_reference`: The reference to a network_function_chain.
* `state`: The state of the subnet.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: subnet UUID.
* `creation_time`: UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: Version number of the latest spec.
* `spec_hash`: Hash of the spec. This will be returned from server.
* `name`: subnet name.
* `should_force_translate`: Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: the key name.
* `value`: value of the key.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, `network_function_chain_reference`, `subnet_reference`.

attributes supports the following:

* `kind`: The kind name (Default value: project.
* `name`: the name.
* `uuid`: the UUID.

Note: `cluster_reference`, `subnet_reference` does not support the attribute `name`

See detailed information in [Nutanix Subnets](https://www.nutanix.dev/api_references/prism-central-v3/#/30ce5964c8d60-get-a-list-of-existing-subnets).
