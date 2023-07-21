---
layout: "nutanix"
page_title: "NUTANIX: nutanix_karbon_worker_nodepool"
sidebar_current: "docs-nutanix-resource-karbon-worker-nodepool"
description: |-
  Provides a resource to add/remove worker nodepool in an existing Nutanix Kubernetes Engine (NKE).
---

# nutanix_karbon_worker_nodepool

Provides a resource to add/remove worker nodepool in an existing Nutanix Kubernetes Engine (NKE).

## Example Usage

```hcl
    resource "nutanix_karbon_worker_nodepool" "kworkerNp" {
        cluster_name = "karbon"
        name = "workerpool1"
        num_instances = 1
        ahv_config {
            cpu= 4
            disk_mib= 122880
            memory_mib=8192
            network_uuid= "61213511-6383-4a38-9ac8-4a552c0e5865"
        }
    }
```

```hcl
    resource "nutanix_karbon_worker_nodepool" "kworkerNp" {
        cluster_name = "karbon"
        name = "workerpool1"
        num_instances = 1
        ahv_config {
            cpu= 4
            disk_mib= 122880
            memory_mib=8192
            network_uuid= "61213511-6383-4a38-9ac8-4a552c0e5865"
        }
        labels={
           k1="v1"
           k2="v2"
	    }
    }
```


## Argument Reference

The following arguments are supported:

* `cluster_name`: (Required) Kubernetes cluster name
* `name`: (Required) unique worker nodepool name
* `node_os_version`: (Optional) The version of the node OS image
* `num_instances`: (Required) number of node instances
* `ahv_config`: (Optional)  VM configuration in AHV.
* `labels`: (Optional) labels of node

### ahv_config
The following arguments are supported for ahv_config:

* `cpu`: - (Required) The number of VCPUs allocated for each VM on the PE cluster.
* `disk_mib`: - (Optional) Size of local storage for each VM on the PE cluster in MiB.
* `memory_mib`: - (Optional) Memory allocated for each VM on the PE cluster in MiB.
* `network_uuid`: - (Required) The UUID of the network for the VMs deployed with this resource configuration.
* `prism_element_cluster_uuid`: - (Optional) The unique universal identifier (UUID) of the Prism Element
* `iscsi_network_uuid`: (Optional) VM network UUID for isolating iscsi data traffic.


## Attributes Reference

The following attributes are exported:

* `nodes`: List of node details of pool.
* `nodes.hostname`: hostname of node
* `nodes.ipv4_address`: ipv4 address of node

## Timeouts

* create
* update
* delete


See detailed information in [Add Node Pool in NKE](https://www.nutanix.dev/api_references/nke/#/c78e2e6b9d9a4-add-a-node-pool-to-a-kubernetes-cluster)
