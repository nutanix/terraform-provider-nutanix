---
layout: "nutanix"
page_title: "NUTANIX: nutanix_virtual_switch_v2"
sidebar_current: "docs-nutanix-datasource-virtual-switch-v2"
description: |-
  Retrieve detailed information about a specific Virtual Switch using the Nutanix v4 API.
---

# nutanix_virtual_switch_v2

The `nutanix_virtual_switch_v2` data source allows you to retrieve comprehensive configuration and status details of a specific Virtual Switch within your Nutanix environment, given its Universally Unique Identifier (UUID).

This data source is extremely useful when you need to reference an existing Virtual Switch's configuration—such as its bonding mode, physical host NIC mappings, or VLAN identifiers—to provision dependencies like Virtual Machines or VPC Mappings.

## Example Usage

```hcl
# Fetch the Virtual Switch details by its UUID
data "nutanix_virtual_switch_v2" "example_vs" {
  ext_id = "1ce108f8-bd48-4287-aead-ca0c3276fc05"
}

# Output the bond mode of the Virtual Switch
output "virtual_switch_bond_mode" {
  value = data.nutanix_virtual_switch_v2.example_vs.bond_mode
}

# Output the VLAN identifier for the first cluster attached to this switch
output "first_cluster_vlan_id" {
  value = data.nutanix_virtual_switch_v2.example_vs.clusters[0].vlan_identifier
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) The globally unique identifier (UUID) of the Virtual Switch you want to retrieve.

## Attribute Reference

The following attributes are exported:

* `name` - The human-readable name of the Virtual Switch.
* `description` - A descriptive text providing context or details about the Virtual Switch's purpose.
* `bond_mode` - The uplink bonding policy applied to this Virtual Switch. Common values include `ACTIVE_BACKUP`, `BALANCE_SLB` (Source Load Balancing), `BALANCE_TCP`, or `NONE`.
* `mtu` - The Maximum Transmission Unit (MTU) configured for the virtual switch, in bytes (e.g., `1500` or `9000` for Jumbo Frames).
* `is_quick_mode` - A boolean flag indicating if quick mode is enabled. When `true`, host nodes are not put into maintenance mode during Virtual Switch updates, which speeds up operations but may cause brief network interruptions.
* `is_default` - A boolean indicating whether this is a default system Virtual Switch (which typically cannot be deleted or heavily modified).
* `has_deployment_error` - A boolean flag that is `true` if the virtual switch configuration failed to deploy consistently across every node in the cluster.
* `has_update_in_progress` - A boolean indicating if the virtual switch is currently undergoing an update operation.
* `has_delete_in_progress` - A boolean indicating if the virtual switch is currently in the process of being deleted.
* `owner_type` - Identifies the type of entity that created or owns the Virtual Switch.
* `project_ext_id` - The UUID of the Nutanix Project that owns this entity, used for RBAC and multi-tenancy.
* `shared_with_projects` - A list of project UUIDs that this virtual switch has been shared with, granting them access to utilize the network.
* `tenant_id` - A globally unique identifier representing the tenant that owns this entity.
* `links` - A list of HATEOAS-style links related to the response, providing URLs for API navigation.
* `metadata` - A block containing administrative and organizational metadata for the mapping, such as assigned categories (tags).

### clusters

The `clusters` block defines the configuration of the Virtual Switch mapped to specific clusters. It exports the following:

* `ext_id` - The UUID of the cluster where this Virtual Switch configuration is applied.
* `gateway_ip_address` - A block containing the `value` and `prefix_length` of the Gateway IP address configured for this cluster's virtual switch.
* `vlan_identifier` - The specific VLAN ID assigned to this virtual switch within the cluster context.
* `hosts` - A list of host-specific configurations within the cluster. See the [`hosts`](#hosts) block below.

### hosts

The `hosts` block details how the Virtual Switch is configured on individual physical nodes. It exports the following:

* `ext_id` - The UUID of the physical host.
* `host_nics` - A list of physical Network Interface Cards (NICs) on the host that are bound to this virtual switch as uplinks.
* `internal_bridge_name` - The system-level logical bridge name created on the host (e.g., `br0`, `br1`).
* `ip_address` - A block containing the `ip` (`value` and `prefix_length`) and overarching `prefix_length` defining the IPv4 subnet address configured for the bridge interface on this host.
* `active_uplink` - The designated active uplink interface when using an `ACTIVE_BACKUP` bond mode.
* `route_table` - The internal routing table ID number utilized by this host for traffic forwarding.

### igmp_spec

The `igmp_spec` block outlines the Internet Group Management Protocol (IGMP) configurations for managing multicast traffic. It exports the following:

* `is_snooping_enabled` - A boolean indicating if IGMP snooping is enabled on the Virtual Switch to optimize multicast traffic delivery.
* `snooping_timeout` - The timeout value in seconds before an inactive multicast group membership expires.
* `querier_spec` - A block defining the IGMP querier settings. See the [`querier_spec`](#querier_spec) block below.

### querier_spec

The `querier_spec` block contains settings for the IGMP querier, which actively queries hosts for their multicast group memberships. It exports the following:

* `is_querier_enabled` - A boolean indicating if the Virtual Switch acts as an IGMP querier.
* `vlan_id_list` - A list of VLAN IDs on which the Virtual Switch actively sends IGMP queries.

See detailed information in [Nutanix Get Virtual Switch V2](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.4#tag/VirtualSwitches/operation/getVirtualSwitchById).
