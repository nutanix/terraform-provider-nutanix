---
layout: "nutanix"
page_title: "NUTANIX: nutanix_unregister_cluster_v2 "
sidebar_current: "docs-nutanix-unregister-cluster-v2"
description: |-
  Unregister a registered remote cluster from the local cluster


---

# nutanix_unregister_cluster_v2 

Unregister a registered remote cluster from the local cluster. This process is asynchronous, creating an unregisteration task and returning its UUID.


## Example Usage

```hcl

resource "nutanix_unregister_cluster_v2 " "pc"{
  pc_ext_id = "<PC_UUID>"
  ext_id = "cluster uuid"
}

```

## Argument Reference
The following arguments are supported:


* `pc_ext_id`: -(Required) The external identifier of the domain manager (Prism Central) resource
* `ext_id`: -(Required) Cluster UUID of a remote cluster.

See detailed information in [Nutanix Unregister a Remote Cluster from Local Cluster Docs](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/unregister).
