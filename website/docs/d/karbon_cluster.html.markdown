---
layout: "nutanix"
page_title: "NUTANIX: nutanix_karbon_cluster"
sidebar_current: "docs-nutanix-datasource-karbon-cluster"
description: |-
 Describes a Karbon Cluster
---

# nutanix_karbon_cluster

Describes a Karbon Cluster

## Example Usage

```hcl
data "nutanix_karbon_cluster" "cluster" {
   karbon_cluster_id = "<YOUR-CLUSTER-ID>"
}
```

## Argument Reference

The following arguments are supported:

* `karbon_cluster_id`: Represents karbon cluster uuid
* `karbon_cluster_name`: Represents the name of karbon cluster

## Attribute Reference

The following attributes are exported:

* `name`: - The name for the k8s cluster.
* `wait_timeout_minutes`: - Maximum wait time for the Karbon cluster to provision.
* `version`: - K8s version of the cluster.
* `single_master_config`: - Configuration of a single master node.
* `active_passive_config`: - The active passive mode uses the Virtual Router Redundancy Protocol (VRRP) protocol to provide high availability of the master.
* `external_lb_config`: - The external load balancer configuration in the case of a multi-master-external-load-balancer type master deployment.
* `private_registry`: - Allows the Karbon cluster to pull images of a list of private registries.
* `etcd_node_pool`: - Configuration of the node pools that the nodes in the etcd cluster belong to. The etcd nodes require a minimum of 8,192 MiB memory and 409,60 MiB disk space.
* `master_node_pool`: - Configuration of the master node pools.
* `cni_config`: - K8s cluster networking configuration. The flannel or the calico configuration needs to be provided.

### External LB Config

The external load balancer configuration in the case of a multi-master-external-load-balancer type master deployment.

* `external_lb_config.#.external_ipv4_address`: The external load balancer IPV4 address.
* `external_lb_config.#.master_nodes_config`: The configuration of the master nodes.
* `external_lb_config.#.master_nodes_config.ipv4_address`: The IPV4 address to assign to the master.
* `external_lb_config.#.master_nodes_config.node_pool_name`: The name of the node pool in which this master IPV4 address will be used.

### private_registry
User inputs of storage configuration parameters for VMs.

* `private_registry`: - List of private registries.
* `private_registry.registry_name`: - Name of the private registry to add to the Karbon cluster.

### Node Pool

The `etcd_node_pool`, `master_node_pool`, `worker_node_pool` attribute supports the following:

* `name`: - Unique name of the node pool.
* `node_os_version`: - The version of the node OS image.
* `num_instances`: - Number of nodes in the node pool.
* `ahv_config`: - VM configuration in AHV.
* `ahv_config.cpu`: - The number of VCPUs allocated for each VM on the PE cluster.
* `ahv_config.disk_mib`: - Size of local storage for each VM on the PE cluster in MiB.
* `ahv_config.memory_mib`: - Memory allocated for each VM on the PE cluster in MiB.
* `ahv_config.network_uuid`: - The UUID of the network for the VMs deployed with this resource configuration.
* `ahv_config.prism_element_cluster_uuid`: - The unique universal identifier (UUID) of the Prism Element cluster used to deploy VMs for this node pool.
* `nodes`: - List of the deployed nodes in the node pool.
* `nodes.hostname`: - Hostname of the deployed node.
* `nodes.ipv4_address`: - IP of the deployed node.

### cni_config

The cni_config attribute exports the following:

* `node_cidr_mask_size`: - The size of the subnet from the pod_ipv4_cidr assigned to each host. A value of 24 would allow up to 255 pods per node.
* `pod_ipv4_cidr`: - CIDR for pods in the cluster.
* `service_ipv4_cidr`: - Classless inter-domain routing (CIDR) for k8s services in the cluster.
* `flannel_config`: - Configuration of the flannel container network interface (CNI) provider.
* `calico_config`: - Configuration of the calico CNI provider.
* `calico_config.ip_pool_config`: - List of IP pools to be configured/managed by calico.
* `calico_config.ip_pool_config.cidr`: - IP range to use for this pool, it should fall within pod cidr.

See detailed information in [Nutanix Karbon Cluster](https://www.nutanix.dev/api_references/nke/#/31fc4aec34c3d-get-the-kubernetes-cluster-details).
