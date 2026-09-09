---
layout: "nutanix"
page_title: "NUTANIX: nutanix_connection_node_v2"
sidebar_current: "docs-nutanix-datasource-connection-node-v2"
description: |-
  Returns details of a node discovered through a specific hardware provider connection.
---

# nutanix_connection_node_v2

Returns details of a node discovered through a specific hardware provider connection.

## Example

```hcl
data "nutanix_connection_node_v2" "example" {
  hardware_provider_ext_id = "00000000-0000-0000-0000-000000000000"
  connection_ext_id        = "11111111-1111-1111-1111-111111111111"
  ext_id                   = "22222222-2222-2222-2222-222222222222"
}
```

## Argument Reference

The following arguments are supported:

* `hardware_provider_ext_id`: (Required) External ID of the hardware provider.
* `connection_ext_id`: (Required) External ID of the hardware provider connection.
* `ext_id`: (Required) External ID of the discovered node.

## Attributes Reference

The following attributes are exported:

* `bmc_ip`: BMC IP address of the node.
* `cpu_info`: CPU information of the node.
* `identifiers`: Identifiers of the node (e.g. serial number).
* `managed_by`: List of systems managing the node.
* `manufacturer`: Manufacturer of the node.
* `memory_gb`: Memory of the node in GB.
* `model`: Model of the node.
* `provider_data`: Provider-specific data for the node.
* `last_refreshed_time`: Last time the node information was refreshed.
* `links`: A HATEOAS style list of links for the response.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
