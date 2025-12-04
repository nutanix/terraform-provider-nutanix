---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pc_unregistration_v2 "
sidebar_current: "docs-nutanix-resource-pc-unregistration-v2"
description: |-
  Unregister a registered remote cluster from the local cluster
---

# nutanix_pc_unregistration_v2

Unregister a registered remote cluster from the local cluster. This process is asynchronous, creating an un-registration task and returning its UUID.

## Example Usage

```hcl

resource "nutanix_pc_unregistration_v2 " "unregister-pc"{
  pc_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17"
  ext_id = "869aa8a5-5aeb-423f-829d-f932d2656b6c"
}

```

## Argument Reference

The following arguments are supported:

- `pc_ext_id`: -(Required) The external identifier of the domain manager (Prism Central) resource
- `ext_id`: -(Required) Cluster UUID of a remote cluster.

See detailed information in [Nutanix PC Unregistration V4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/unregister).
