---
layout: "nutanix"
page_title: "NUTANIX: nutanix_nodes_v2"
sidebar_current: "docs-nutanix-datasource-nodes-v2"
description: |-
  Returns a paginated list of all nodes managed by Foundation Central.
---

# nutanix_nodes_v2

Returns a paginated list of all nodes managed by Foundation Central.

## Example

```hcl
data "nutanix_nodes_v2" "example" {}

data "nutanix_nodes_v2" "filtered" {
  limit  = 10
  filter = "manufacturer eq 'Nutanix'"
}
```

## Argument Reference

* `page` - (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `expand` - (Optional) A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved.
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.

## Attributes Reference

* `nodes` - List of nodes managed by Foundation Central. Each element has the same attributes as the [nutanix_node_v2](node_v2.html) datasource.

See detailed information in [Nutanix List Nodes V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
