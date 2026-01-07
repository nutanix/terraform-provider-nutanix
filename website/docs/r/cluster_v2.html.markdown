---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster_entity_v2"
sidebar_current: "docs-nutanix-resource-cluster-entity-v2"
description: |-
   Provides the basic infrastructure for compute, storage and networking.
---

# nutanix_cluster_v2

Represents the Cluster entity. Provides the basic infrastructure for compute, storage and networking. This includes the operations that can be carried out on cluster and its subresources - host (node), rsyslog servers etc and actions that can be performed on cluster - add a node, remove a node, attach categories.

-> **Recommendations:** It is recommended to create and register the cluster with Prism Central as part of the same workflow. Cluster updates, importing, and destruction through Terraform are supported only when the cluster is registered with Prism Central. To register a cluster with Prism Central use Terraform resource nutanix_pc_registration_v2.

-> **Note:**: Cluster resource supports add/remove node operations. However, these operations require cluster to be registered with Prism Central.

-> **Note:**: TThe cluster resource supports both associating and disassociating categories, allowing you to attach or detach categories on a cluster through Terraform. However, these operations require cluster to be registered with Prism Central.

**Note:**: The cluster resource supports both associating and disassociating cluster profile, allowing you to attach or detach cluster profile on a cluster through Terraform. However, these operations require cluster to be registered with Prism Central.

## Example Usage:

### Example 1: 1 Node Cluster Creation Example

```hcl
resource "nutanix_cluster_v2" "cluster"{
  name = "cluster-example"
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = "10.xx.xx.xx"
        }
      }
    }
  }
  config {
    cluster_function  = ["AOS"]
    redundancy_factor = 1
    cluster_arch      = "X86_64"
    fault_tolerance_state {
      domain_awareness_level = "DISK"
    }
  }
  network {
    external_address {
      ipv4 {
        value = "10.xx.xx.xx"
      }
    }
    external_data_services_ip {
      ipv4 {
        value = "10.xx.xx.xx"
      }
    }
    ntp_server_ip_list {
      fqdn {
        value = "ntp.server.nutanix.com"
      }
    }
    ntp_server_ip_list {
      fqdn {
        value = "ntp.server_1.nutanix.com"
      }
    }
    smtp_server {
      email_address = "example.ex@exmple.com"
      server {
        ip_address {
          ipv4 {
            value = "10.xx.xx.xx"
          }
        }
        port     = 123
        username = "example"
        password = "example!2134"
      }
      type = "PLAIN"
    }
  }
}
```

### Example 2: 3 Node Cluster Creation Example and Adding Nodes Example

```hcl
resource "nutanix_cluster_v2" "cluster-3nodes" {
  name   = "tf-cluster-3nodes"
  dryrun = false
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = "10.00.00.1"
        }
      }
    }
    node_list {
      controller_vm_ip {
        ipv4 {
          value = "10.00.00.2"
        }
      }
    }
    node_list {
      controller_vm_ip {
        ipv4 {
          value = "10.00.00.3"
        }
      }
    }

    ## after creating and registering the cluster with prism central,
    ## to add another node to the cluster, uncomment the following block
    # node_list {
    #   controller_vm_ip {
    #     ipv4 {
    #       value = "10.00.00.4"
    #     }
    #   }
    #   should_skip_host_networking   = false
    #   should_skip_pre_expand_checks = true
    # }
  }
  config {
    cluster_function = local.clusters.config.cluster_functions
    cluster_arch     = local.clusters.config.cluster_arch
    fault_tolerance_state {
      domain_awareness_level = "NODE"
    }
  }

  lifecycle {
    ignore_changes = [links, categories, config.0.cluster_function]
  }
}
```


### Example 3: Creating a cluster with categories

```hcl
resource "nutanix_cluster_v2" "cluster-with-categories" {
  name = "cluster-example"
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = "10.xx.xx.xx"
        }
      }
    }
  }
  config {
    cluster_function  = ["AOS"]
    redundancy_factor = 1
    cluster_arch      = "X86_64"
    fault_tolerance_state {
      domain_awareness_level = "DISK"
    }
  }
  network {
    external_address {
      ipv4 {
        value = "10.xx.xx.xx"
      }
    }
    external_data_services_ip {
      ipv4 {
        value = "10.xx.xx.xx"
      }
    }
    ntp_server_ip_list {
      fqdn {
        value = "ntp.server.nutanix.com"
      }
    }
    ntp_server_ip_list {
      fqdn {
        value = "ntp.server_1.nutanix.com"
      }
    }
    smtp_server {
      email_address = "example.ex@exmple.com"
      server {
        ip_address {
          ipv4 {
            value = "10.xx.xx.xx"
          }
        }
        port     = 123
        username = "example"
        password = "example!2134"
      }
      type = "PLAIN"
    }
  }
  ## after creating and registering the cluster with prism central,
  ## to associate categories to the cluster, uncomment the following block
  ## and to disassociate categories from the cluster, remove it from the categories list.
  # categories = ["72a80be7-a6d7-4ad4-8071-3d92f40f1e7d", "45588de3-7c18-4230-a147-7e26ad92d8a6"]
}
```

## Argument Reference

The following arguments are supported:
> after creating the cluster, you need to register the cluster with prism central to be able to use it.

* `dryrun`: - (Optional) parameter that allows long-running operations to execute in a dry-run mode providing ability to identify trouble spots and system failures without performing the actual operation. Additionally this mode also offers a summary snapshot of the resultant system in order to better understand how things fit together. The operation runs in dry-run mode only if the provided value is true.
* `name`: - (Required) The name for the vm.
* `nodes`: - (Optional) The reference to a node and remove node parameters.
* `config`: - (Optional) Cluster configuration details.
* `network`: - (Optional) Network details of a cluster.
* `upgrade_status`: - (Optional) Upgrade status of a cluster.
    Valid values are:
     - "CANCELLED"	The cluster upgrade is cancelled.
     - "FAILED"	The cluster upgrade failed.
     - "QUEUED"	The cluster upgrade is in the queue.
     - "SUCCEEDED"	The cluster was upgraded successfully.
     - "DOWNLOADING"	The cluster upgrade is downloading.
     - "PENDING"	The cluster upgrade is in pending state.
     - "UPGRADING"	The cluster is in upgrade state.
     - "PREUPGRADE"	The cluster is in pre-upgrade state.
     - "SCHEDULED"	The cluster upgrade is in scheduled state.
* `container_name`: - (Optional) The name of the default container created as part of cluster creation. This is part of payload for cluster create operation only.
* `categories`: - (Optional) The reference to a project.
* `cluster_profile_ext_id`: - (Optional) The reference to a cluster profile.

### Nodes

The nodes attribute supports the following:

* `node_list`: - (Optional) List of nodes in a cluster.
* `remove_node_params`: - (Optional) Parameters for removing nodes. Supports:

#### remove_node_params Attributes

* `should_skip_remove`: - (Optional, default false) Skip remove operation for the node.
* `should_skip_prechecks`: - (Optional, default false) Skip remove prechecks.
* `extra_params`: - (Optional) Extra parameters for removing nodes. Supports:
  * `should_skip_upgrade_check`: - (Optional, default false) Skip upgrade check during node removal.
  * `skip_space_check`: - (Optional, default false) Skip space check during node removal.
  * `should_skip_add_check`:- (Optional, default false) Skip add check during node removal.

### Node List

The nodes attribute supports the following:

* `controller_vm_ip`: - (Required) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `host_ip`: - (Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `should_skip_host_networking`: - (Optional, default false) Flag to indicate if host networking needs to be skipped during node addition.
* `should_skip_pre_expand_checks`: - (Optional, default false) Flag to indicate if pre expand checks needs to be skipped during node addition.
* `should_skip_add_node`: - (Optional, default false) Flag to indicate if add node operation needs to be skipped during node addition.

### Controller VM IP

The controller_vm_ip attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.

### Host IP

The host_ip attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.

### Config

The config attribute supports the following:

* `build_info`: - (Optional) Build information details.
* `cluster_function`: - (Optional) Cluster function. This is part of payload for cluster
  create operation only (allowed enum values for creation are AOS, ONE_NODE & TWO_NODE only).
    Valid values are:
     - "AOS"
     - "PRISM_CENTRAL"
     - "CLOUD_DATA_GATEWAY"
     - "AFS"
     - "ONE_NODE"
     - "TWO_NODE"
     - "ANALYTICS_PLATFORM"
* `authorized_public_key_list`: - (Optional) Public ssh key details. This is part of payload for cluster update operation only.
* `redundancy_factor`: - (Optional) Redundancy factor of a cluster. This is part of payload for both cluster create & update operations.
* `cluster_arch`: - (Optional) Cluster arch.
    Valid values are:
     - "PPC64LE" PPC64LE cluster architecture type.
     - "X86_64" X86_64 cluster architecture type.
* `fault_tolerance_state`: - (Optional) Fault tolerant state of cluster.
* `operation_mode`: - (Optional) Cluster operation mode. This is part of payload for cluster
  update operation only.
    Valid values are:
     - "OVERRIDE"	Override operation mode.
     - "STAND_ALONE"	Stand-alone operation mode.
     - "SWITCH_TO_TWO_NODE"	Switch to two-node operation mode.
     - "NORMAL"	Normal operation mode.
     - "READ_ONLY"	Read-only operation mode.
* `encryption_in_transit_status`: - (Optional) Encryption in transit Status.
    Valid values are:
     - "DISABLED"	Disabled encryption status.
     - "ENABLED"	Enabled encryption status.
* `pulse_status`: - (Optional) Pulse status for a cluster. `supported only for update operations and not available during creation.`

### Build info

The build_info attribute supports the following:

* `build_type` Software build type.
* `version` Software version.
* `full_version` Full name of software version.
* `commit_id` Commit Id used for version.
* `short_commit_id` Short commit Id used for version.

### Authorized Public Key List

The authorized_public_key_list attribute supports the following:

* `name` (required) Ssh key name.
* `key` (required) Ssh key value.

### Fault Tolerance State
The fault_tolerance_state attribute supports the following:

* `domain_awareness_level` Domain awareness level corresponds to unit of cluster group. This is part of payload for both cluster create & update operations.
    Valid values are:
    - "RACK"	Rack level awareness.
    - "NODE"	Node level awareness.
    - "BLOCK"	Block level awareness.
    - "DISK"	Disk level awareness.

### Pulse Status
The pulse_status attribute supports the following:
* `is_enabled`: - (Optional) Flag to indicate if pulse is enabled or not.
* `pii_scrubbing_level`: - (Optional) PII scrubbing level.
    Valid values are:
    - "ALL" :	Scrub All PII Information from Pulse including data like entity names and IP addresses.
    - "DEFAULT":	Default PII Scrubbing level. Data like entity names and IP addresses will not be scrubbed from Pulse.

### Network

* `external_address` An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `external_data_services_ip` An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `nfs_subnet_white_list` NFS subnet whitelist addresses. This is part of payload for cluster update operation only.
* `name_server_ip_list` List of name servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
* `ntp_server_ip_list` List of NTP servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
* `smtp_server` SMTP servers on a cluster. This is part of payload for cluster update operation only.
* `masquerading_ip` An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `management_server` Management server information.
* `fqdn` Cluster fully qualified domain name. This is part of payload for cluster update operation only.
* `key_management_server_type` Management server type.
    Valid values are:
    - "PRISM_CENTRAL"	Prism Central management server.
    - "EXTERNAL"	External management server.
    - "LOCAL"	Local management server.
* `backplane` Params associated to the backplane network segmentation. This is part of payload for cluster(create operation only.)
* `http_proxy_list` List of HTTP Proxy server configuration needed to access a cluster which is hosted behind a HTTP Proxy to not reveal its identity.
* `https_proxy_white_list` Targets HTTP traffic to which is exempted from going through the configured HTTP Proxy.

### External Address

The external_address attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.

### External Data Services IP

The external_data_services_ip attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.

### Name Server IP List

The name_server_ip_list attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.
* `fqdn`: - (Optional) A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

### Ntp Server IP List

The ntp_server_ip_list attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.
* `fqdn`: - (Optional) A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

### SMTP server

The smtp_server attribute supports the following:

* `email_address` (required) SMTP email address.

* `server` (required) SMTP network details.

* `type` Type of SMTP server.
    Valid values are:
    - "PLAIN"	Plain type SMTP server.
    - "STARTTLS"	Start TLS type SMTP server.
    - "SSL"	SSL type SMTP server.

### Server

The server attribute supports the following:

* `ip_address` (required) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `port` SMTP port.
* `username` SMTP server user name.
* `password` SMTP server password.

### IP Address

The ip_address attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.
* `fqdn`: - (Optional) A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

### Masquerading IP

The masquerading_ip attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.

### Management Server

The management_server attribute supports the following:

* `ip` An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `type` Management server type.
    Valid values are:
    - "VCENTER"	Vcenter management server.
* `drs_enabled` Indicates whether it is DRS enabled or not.
* `is_registered` Indicates whether it is registered or not.
* `in_use` Indicates whether the host is managed by an entity or not.

### IP

The ip attribute supports the following:

* `ipv4`: - (Optional) ip address params.
* `ipv6`: - (Optional) Ip address params.

### Backplane

The backplane attribute supports the following:

* `is_segmentation_enabled` Flag to indicate if the backplane segmentation needs to be enabled or not.
* `vlan_tag` VLAN Id tagged to the backplane network on the cluster. This is part of cluster create payload.
* `subnet` Subnet configs.
* `netmask` Netmask configs.

### Subnet

The subnet attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - (Required) Ip address.

### Netmask

The netmask attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - (Required) Ip address.

### HTTP Proxy List
The http_proxy_list attribute supports the following:

* `ip_address`: - (Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `port`: - (Optional) HTTP Proxy server port configuration needed to access a cluster which is hosted behind a HTTP Proxy to not reveal its identity.
* `username`: - (Optional) HTTP Proxy server username needed to access a cluster which is hosted behind a HTTP Proxy to not reveal its identity.
* `password`: - (Optional) HTTP Proxy server password needed to access a cluster which is hosted behind a HTTP Proxy to not reveal its identity.
* `name`: - (Required) HTTP Proxy server name configuration needed to access a cluster which is hosted behind a HTTP Proxy to not reveal its identity.
* `proxy_type`: - (Optional) HTTP Proxy server type.
    Valid values are:
    - "HTTP"	HTTP Proxy server type.
    - "HTTPS"	HTTPS Proxy server type.
    - "SOCKS"	SOCKS Proxy server type.

### HTTPS Proxy White List
The https_proxy_white_list attribute supports the following:

* `target_type`: - (Optional) Target type.
    Valid values are:
    - "IPV6_ADDRESS"	IPV6 address.
    - "HOST_NAME"	Name of the host.
    - "IPV4_ADDRESS"	IPV4 address.
    - "DOMAIN_NAME_SUFFIX"	Domain Name Suffix required for http proxy whitelist.
    - "IPV4_NETWORK_MASK" Network Mask of the IpV4 family.

* `target`: - (Required) Target's identifier which is exempted from going through the configured HTTP Proxy.
### IPV4

The ipv4 attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - (Required) Ip address.

### IPV6

The ipv6 attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - (Required) Ip address.

## Import

This helps to manage existing entities which are not created through terraform. Clusters can be imported using the `UUID`.  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_cluster_v2" "import_cluster" {}

// execute this command in cli
terraform import nutanix_cluster_v2.import_cluster <UUID>
```

## What happens when you do terraform destroy for nutanix_clusters_v2?  First thing, inorder to destroy the cluster from Terraform it need to be registered.
```
1. Internally, It deregisters the cluster and its nodes from Prism Central

* The cluster and all its member nodes are gracefully deregistered from Prism Central.

* This ensures the cluster is no longer visible in the PC inventory or associated with any management plane workflows.

2. Remove the clustering configuration on each node

* Once deregistration completes, teardown of the cluster configuration on each node begins.

* This effectively disbands the cluster, converting all participating nodes back into standalone (unconfigured) nodes.

* At this stage, the nodes no longer share storage, network, or management metadata — they’re ready to be re-clustered or repurposed.

```

See detailed information in [Nutanix Create Cluster V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0#tag/Clusters/operation/createCluster).
