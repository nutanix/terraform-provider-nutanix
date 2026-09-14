---
layout: "nutanix"
page_title: "NUTANIX: nutanix_node_v2"
sidebar_current: "docs-nutanix-resource-node-v2"
description: |-
  Registers and onboards a new node for management by Foundation Central.
---

# nutanix_node_v2

Registers and onboards a new node for management by Foundation Central.

## Example

```hcl
resource "nutanix_node_v2" "example" {
  manufacturer = "Nutanix"
  model        = "NX-3060-G8"

  identifiers {
    type  = "SERIAL_NUMBER"
    value = "ZM00000000000000"
  }
}
```

## Argument Reference

The following arguments are supported:

* `identifiers` - (Required) Unique identifier array for the node (currently limited to 1 item with NODE_SERIAL type). See [Identifiers](#identifiers) below.
* `manufacturer` - (Required) Manufacturer of the node.
* `model` - (Optional) Model of the node.
* `hostname` - (Optional) Hostname of the node.
* `host_type` - (Optional) Type of the host installed on the node. One of `AHV`, `ESX`.
* `host_version` - (Optional) Version of the host installed on the node.
* `aos_version` - (Optional) Version of the AOS installed on the node.
* `block_serial_number` - (Optional) Block serial number of the node.
* `memory_gb` - (Optional) Total memory available in the node (in GB).
* `custom_attributes` - (Optional) User tags of the node.
* `owner_ext_id` - (Optional) External ID of the owner.
* `provider_ext_id` - (Optional) External ID of the provider of the node.
* `provider_connection_ext_id` - (Optional) External ID of the provider connection of the node.
* `cpu_info` - (Optional) CPU details of the node. See [CPU Info](#cpu-info) below.
* `network_details` - (Optional) Network details of the node. See [Network Details](#network-details) below.

### Identifiers

* `type` - (Required) Type of the identifier. One of `SERIAL_NUMBER`.
* `value` - (Required) Unique identifier value of the node.

### CPU Info

* `capacity_g_hz` - (Optional) CPU clock speed in GHz.
* `logical_core_count` - (Optional) Number of logical CPU cores.
* `manufacturer` - (Optional) Manufacturer of the CPU (e.g., Intel, AMD).
* `model` - (Optional) Model name of the CPU.
* `socket_count` - (Optional) Number of CPU sockets in the node.

### Network Details

* `bmc` - (Optional) BMC network details. See [Network Component](#network-component) below.
* `cvm` - (Optional) CVM network details. Contains `management_network`. See [Network Component](#network-component) below.
* `host` - (Optional) Host network details. Contains `management_network`. See [Network Component](#network-component) below.

### Network Component

* `gateway` - (Optional) Gateway IP address (`ipv4` / `ipv6` with `value` and `prefix_length`).
* `ip` - (Optional) IP address (`ipv4` / `ipv6` with `value` and `prefix_length`).
* `vlan_id` - (Optional) VLAN ID of the network.

## Attributes Reference

In addition to the arguments above, the following computed attributes are exported:

* `ext_id` - A globally unique identifier of the node.
* `cvm_connectivity_status` - Connectivity status of the CVM.
* `host_connectivity_status` - Connectivity status of the host.
* `state` - State of the node (e.g., `DISCOVERED`, `ONBOARDED`, `CONFIGURED_HARDWARE_PROVIDER`).
* `created_time` - Time when the node was discovered in Foundation Central.
* `provider_data` - Extended provider data discovered for the node.
* `links` - A HATEOAS style link for the response.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Register Node V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
