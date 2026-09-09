---
layout: "nutanix"
page_title: "NUTANIX: nutanix_nodes_v2"
sidebar_current: "docs-nutanix-datasource-nodes-v2"
description: |-
  Returns a list of nodes registered with a specific claim token.
---

# nutanix_nodes_v2

Returns a list of nodes registered with a specific claim token.

## Example

```hcl
data "nutanix_nodes_v2" "example" {
  claim_token_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `claim_token_ext_id`: (Required) External ID of the claim token.
* `page`: (Optional) A URL query parameter that specifies the page number of the result set.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attributes Reference

* `nodes`: List of nodes registered with the claim token. Each element has the same attributes as the [nutanix_node_v2](node_v2.html) datasource.

See detailed information in [Nutanix List Nodes By Claim Token V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
