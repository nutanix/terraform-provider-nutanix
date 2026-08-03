---
layout: "nutanix"
page_title: "NUTANIX: nutanix_node_schedulable_statuses_v2"
sidebar_current: "docs-nutanix-datasource-node-schedulable-statuses-v2"
description: |-
  Check to see whether a node in a cluster is a storage-only node or not.
---

# nutanix_node_schedulable_statuses_v2

Check to see whether a node in a cluster is a storage-only node or not.

## Example Usage

```hcl
data "nutanix_node_schedulable_statuses_v2" "example" {}
```

## Argument Reference

* `cluster_id` - (Optional) Prism Element cluster reference.

## Attribute Reference

* `node_schedulable_statuses` - List of node schedulable status entries.

### node_schedulable_statuses

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `is_never_schedulable` - The boolean value to indicate whether or not node is a storage only node.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.
