---
layout: "nutanix"
page_title: "NUTANIX: nutanix_node_v2"
sidebar_current: "docs-nutanix-datasource-node-v2"
description: |-
  Returns the hardware, network, and state details of a node by its external ID.
---

# nutanix_node_v2

Returns the hardware, network, and state details of a node by its external ID.

## Example

```hcl
data "nutanix_node_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id` - (Required) External ID of the node.

## Attributes Reference

The following attributes are exported:

* `identifiers` - Unique identifier array for the node. Each item has `type` and `value`.
* `manufacturer` - Manufacturer of the node.
* `model` - Model of the node.
* `hostname` - Hostname of the node.
* `host_type` - Type of the host installed on the node.
* `host_version` - Version of the host installed on the node.
* `aos_version` - Version of the AOS installed on the node.
* `block_serial_number` - Block serial number of the node.
* `memory_gb` - Total memory available in the node (in GB).
* `custom_attributes` - User tags of the node.
* `owner_ext_id` - External ID of the owner.
* `provider_ext_id` - External ID of the provider of the node.
* `provider_connection_ext_id` - External ID of the provider connection of the node.
* `cpu_info` - CPU details of the node.
* `network_details` - Network details of the node (`bmc`, `cvm`, `host`).
* `cvm_connectivity_status` - Connectivity status of the CVM.
* `host_connectivity_status` - Connectivity status of the host.
* `state` - State of the node.
* `created_time` - Time when the node was discovered in Foundation Central.
* `provider_data` - Extended provider data discovered for the node.
* `links` - A HATEOAS style link for the response.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Get Node V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
