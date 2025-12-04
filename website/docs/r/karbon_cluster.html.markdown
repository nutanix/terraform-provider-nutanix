---
layout: "nutanix"
page_title: "NUTANIX: nutanix_karbon_cluster"
sidebar_current: "docs-nutanix-resource-karbon-cluster"
description: |-
  Provides a Nutanix Karbon Cluster resource to Create a k8s cluster.
---

# nutanix_karbon_cluster

Provides a Nutanix Karbon Cluster resource to Create a k8s cluster.

**Note:** Minimum tested version is Karbon 2.2

**Note:** Kubernetes and Node OS upgrades are not supported using this provider.

## Example Usage

```hcl
resource "nutanix_karbon_cluster" "example_cluster" {
  name       = "example_cluster"
  version    = "1.18.15-1"
  storage_class_config {
    reclaim_policy = "Delete"
    volumes_config {
      file_system                = "ext4"
      flash_mode                 = false
      password                   = "my_pe_pw"
      prism_element_cluster_uuid = "my_pe_cluster_uuid"
      storage_container          = "my_storage_container_name"
      username                   = "my_pe_username"
    }
  }
  cni_config {
    node_cidr_mask_size = 24
    pod_ipv4_cidr       = "172.20.0.0/16"
    service_ipv4_cidr   = "172.19.0.0/16"
  }
  worker_node_pool {
    node_os_version = "ntnx-1.0"
    num_instances   = 1
    ahv_config {
      network_uuid               = "my_subnet_id"
      prism_element_cluster_uuid = "my_pe_cluster_uuid"
    }
  }
  etcd_node_pool {
    node_os_version = "ntnx-1.0"
    num_instances   = 1
    ahv_config {

      network_uuid               = "my_subnet_id"
      prism_element_cluster_uuid = "my_pe_cluster_uuid"
    }
  }
  master_node_pool {
    node_os_version = "ntnx-1.0"
    num_instances   = 1
    ahv_config {
      network_uuid               = "my_subnet_id"
      prism_element_cluster_uuid = "my_pe_cluster_uuid"
    }
  }
}

```


### resource to create karbon cluster with timeouts
```hcl
resource "nutanix_karbon_cluster" "example_cluster" {
  name       = "example_cluster"
  version    = "1.18.15-1"
  storage_class_config {
    reclaim_policy = "Delete"
    volumes_config {
      file_system                = "ext4"
      flash_mode                 = false
      password                   = "my_pe_pw"
      prism_element_cluster_uuid = "my_pe_cluster_uuid"
      storage_container          = "my_storage_container_name"
      username                   = "my_pe_username"
    }
  }
  cni_config {
    node_cidr_mask_size = 24
    pod_ipv4_cidr       = "172.20.0.0/16"
    service_ipv4_cidr   = "172.19.0.0/16"
  }
  worker_node_pool {
    node_os_version = "ntnx-1.0"
    num_instances   = 1
    ahv_config {
      network_uuid               = "my_subnet_id"
      prism_element_cluster_uuid = "my_pe_cluster_uuid"
    }
  }
  etcd_node_pool {
    node_os_version = "ntnx-1.0"
    num_instances   = 1
    ahv_config {
      network_uuid               = "my_subnet_id"
      prism_element_cluster_uuid = "my_pe_cluster_uuid"
    }
  }
  master_node_pool {
    node_os_version = "ntnx-1.0"
    num_instances   = 1
    ahv_config {
      network_uuid               = "my_subnet_id"
      prism_element_cluster_uuid = "my_pe_cluster_uuid"
    }
  }
  timeouts {
    create = "1h"
    update = "30m"
    delete = "10m"
	}
}
```


## Argument Reference

The following arguments are supported:

* `name`: - (Required) The name for the k8s cluster. **Note:** Updates to this attribute forces new resource creation.
* `wait_timeout_minutes`: - (Optional) Maximum wait time for the Karbon cluster to provision.
* `version`: - (Required) K8s version of the cluster. **Note:** Updates to this attribute forces new resource creation.
* `storage_class_config`: - (Required) Storage class configuration attribute for defining the persistent volume attributes. **Note:** Updates to this attribute forces new resource creation.
* `single_master_config`: - (Optional) Configuration of a single master node. **Note:** Updates to this attribute forces new resource creation.
* `active_passive_config`: - (Optional) The active passive mode uses the Virtual Router Redundancy Protocol (VRRP) protocol to provide high availability of the master. **Note:** Updates to this attribute forces new resource creation.
* `external_lb_config`: - (Optional) The external load balancer configuration in the case of a multi-master-external-load-balancer type master deployment. **Note:** Updates to this attribute forces new resource creation.
* `private_registry`: - (Optional) Allows the Karbon cluster to pull images of a list of private registries.
* `etcd_node_pool`: - (Required) Configuration of the node pools that the nodes in the etcd cluster belong to. The etcd nodes require a minimum of 8,192 MiB memory and 409,60 MiB disk space.
* `master_node_pool`: - (Required) Configuration of the master node pools.
* `cni_config`: - (Required) K8s cluster networking configuration. The flannel or the calico configuration needs to be provided. **Note:** Updates to this attribute forces new resource creation.

### Storage Class Config

The storage_class_config attribute supports the following: 


* `name`: - (Required) The name of the storage class.
* `reclaim_policy` - (Optional) Reclaim policy for persistent volumes provisioned using the specified storage class.
* `volumes_config.#.file_system` - (Optional) Karbon uses either the ext4 or xfs file-system on the volume disk.
* `volumes_config.#.flash_mode` - (Optional) Pins the persistent volumes to the flash tier in case of a `true` value.
* `volumes_config.#.password` - (Required) The password of the Prism Element user that the API calls use to provision volumes.
* `volumes_config.#.prism_element_cluster_uuid` - (Required) The universally unique identifier (UUID) of the Prism Element cluster.
* `volumes_config.#.storage_container` - (Required) Name of the storage container the storage container uses to provision volumes.
* `volumes_config.#.username` - (Required) Username of the Prism Element user that the API calls use to provision volumes.

**Note:** Updates to this attribute forces new resource creation.
### Single Master Config

The `single_master_config` defines the deployment of a Karbon cluster with a single master setup. This is the default behavior unless the `active_passive_config` or the `external_lb_config` attributes are passed.

**Note:** Updates to this attribute forces new resource creation.
### Active-Passive Config

The `active_passive_config` attribute can be used in case a multi-master active-passive deployment is required. The external_ipv4_address should be an IP address in the Karbon cluster subnet range. The Virtual Router Redundancy Protocol (VRRP) protocol is used to provide high availability of the master.

* `active_passive_config.#.external_ipv4_address`: (Required) The VRRP IPV4 address to be used by the masters.

**Note:** Updates to this attribute forces new resource creation.
### External LB Config

The external load balancer configuration in the case of a multi-master-external-load-balancer type master deployment.

* `external_lb_config.#.external_ipv4_address`: (Required) The external load balancer IPV4 address.
* `external_lb_config.#.master_nodes_config`: (Required) The configuration of the master nodes.
* `external_lb_config.#.master_nodes_config.ipv4_address`: (Required) The IPV4 address to assign to the master.
* `external_lb_config.#.master_nodes_config.node_pool_name`: (Optional) The name of the node pool in which this master IPV4 address will be used.

**Note:** Updates to this attribute forces new resource creation.
### Private Registry
User inputs of storage configuration parameters for VMs.

* `private_registry`: - (Optional) List of private registries.
* `private_registry.registry_name`: - (Required) Name of the private registry to add to the Karbon cluster.

### Node Pool

The `etcd_node_pool`, `master_node_pool`, `worker_node_pool` attribute supports the following:

* `name`: - (Optional) Unique name of the node pool. **Note:** Updates to this attribute forces new resource creation.
* `node_os_version`: - (Required) The version of the node OS image. **Note:** Updates to this attribute forces new resource creation.
* `num_instances`: - (Required) Number of nodes in the node pool. **Note:** Updates to etcd or master node pool forces new resource creation.
* `ahv_config`: - (Optional) VM configuration in AHV. **Note:** Updates to this attribute forces new resource creation.
* `ahv_config.cpu`: - (Required) The number of VCPUs allocated for each VM on the PE cluster.
* `ahv_config.disk_mib`: - (Optional) Size of local storage for each VM on the PE cluster in MiB.
* `ahv_config.memory_mib`: - (Optional) Memory allocated for each VM on the PE cluster in MiB.
* `ahv_config.network_uuid`: - (Required) The UUID of the network for the VMs deployed with this resource configuration.
* `ahv_config.prism_element_cluster_uuid`: - (Required) The unique universal identifier (UUID) of the Prism Element cluster used to deploy VMs for this node pool.
* `nodes`: - List of the deployed nodes in the node pool.
* `nodes.hostname`: - Hostname of the deployed node.
* `nodes.ipv4_address`: - IP of the deployed node.

### CNI Config

 The boot_device_disk_address attribute supports the following:

* `node_cidr_mask_size`: - (Optional) The size of the subnet from the pod_ipv4_cidr assigned to each host. A value of 24 would allow up to 255 pods per node.
* `pod_ipv4_cidr`: - (Optional) CIDR for pods in the cluster.
* `service_ipv4_cidr`: - (Optional) Classless inter-domain routing (CIDR) for k8s services in the cluster.
* `flannel_config`: - (Optional) Configuration of the flannel container network interface (CNI) provider.
* `calico_config`: - (Optional) Configuration of the calico CNI provider.
* `calico_config.ip_pool_config`: - (Optional) List of IP pools to be configured/managed by calico.
* `calico_config.ip_pool_config.cidr`: - (Optional) IP range to use for this pool, it should fall within pod cidr.

* `timeouts`: timeouts can customize the default timeout on CRUD functions with default timeouts. Supports "h", "m" or "s" . 

**Note:** Updates to this attribute forces new resource creation.

See detailed information in [Nutanix Karbon Cluster](https://www.nutanix.dev/api_references/nke/#/895c7a174c68b-create-a-new-kubernetes-cluster).
