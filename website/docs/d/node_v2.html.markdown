---
layout: "nutanix"
page_title: "NUTANIX: nutanix_node_v2"
sidebar_current: "docs-nutanix-datasource-node-v2"
description: |-
  Returns details of a node registered with a specific claim token.
---

# nutanix_node_v2

Returns details of a node registered with a specific claim token.

## Example

```hcl
data "nutanix_node_v2" "example" {
  claim_token_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id             = "00000000-0000-0000-0000-000000000001"
}
```

## Argument Reference

* `claim_token_ext_id`: (Required) External ID of the claim token.
* `ext_id`: (Required) External ID of the node.

## Attributes Reference

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `created_time`: Time when the node was discovered in Foundation Central.
* `host_type`: Type of host installed on the node.
* `host_version`: Version of the host installed on the node.
* `manufacturer`: Manufacturer of the node hardware.
* `model`: Model of the node.
* `memory_gb`: Total memory available in the node (in GB).
* `managed_by`: Systems that manage the node.
* `cpu_info`: CPU information of the node (`capacity_ghz`, `logical_core_count`, `manufacturer`, `model`, `socket_count`).
* `identifiers`: Unique identifier array for the node (`type`, `value`).
* `network_details`: Network details of the node (`bmc`, `cvm`, `host`).
* `links`: A HATEOAS style link for the response.

See detailed information in [Nutanix Get Node By Claim Token V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
