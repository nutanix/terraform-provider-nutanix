---
layout: "nutanix"
page_title: "NUTANIX: nutanix_virtual_switches_v2"
sidebar_current: "docs-nutanix-datasource-virtual-switches-v2"
description: |-
  Retrieve the list of Virtual Switches using the Nutanix v4 API.
---

# nutanix_virtual_switches_v2

The `nutanix_virtual_switches_v2` data source allows you to retrieve the list of Virtual Switches configured within your Nutanix environment, optionally scoped to a specific Prism Element cluster.

This data source is useful when you need to discover available Virtual Switches and reference their configuration—such as bonding mode, physical host NIC mappings, or VLAN identifiers—to provision dependencies like Virtual Machines or VPC Mappings.

## Example Usage

```hcl
# Fetch all Virtual Switches
data "nutanix_virtual_switches_v2" "example" {}

# Fetch the Virtual Switches scoped to a specific Prism Element cluster
data "nutanix_virtual_switches_v2" "by_cluster" {
  cluster_id = "00065264-61ec-9c95-185b-ac1f6b6f97e2"
}

# Output the name of the first Virtual Switch in the list
output "first_virtual_switch_name" {
  value = data.nutanix_virtual_switches_v2.example.virtual_switches[0].name
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) The UUID of the Prism Element cluster used to filter the returned Virtual Switches.
* `page` - (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between `0` and the maximum number of pages that are available for that resource.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set. The default value is `50` records.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources, following the [OData V4.01 URL conventions](https://docs.oasis-open.org/odata/odata/v4.01/odata-v4.01-part2-url-conventions.html).
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.

## Attribute Reference

The following attributes are exported:

* `virtual_switches` - The list of Virtual Switch entities. See the [`virtual_switches`](#virtual_switches) block below for details.

### virtual_switches

Each entry in the `virtual_switches` list exports the following attributes:

* `ext_id` - The globally unique identifier (UUID) of the Virtual Switch.
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
* `metadata` - A block containing administrative and organizational metadata for the entity, such as assigned categories (tags).

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

See detailed information in [Nutanix List Virtual Switches V2](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.4#tag/VirtualSwitches/operation/listVirtualSwitches).