---
layout: "nutanix"
page_title: "NUTANIX: nutanix_clusters_v2"
sidebar_current: "docs-nutanix-datasource-clusters-v2"
description: |-
  Lists all cluster entities registered to Prism Central.
---

# nutanix_clusters_v2

Lists all cluster entities registered to Prism Central.

## Example Usage

```hcl
data "nutanix_clusters_v2" "cls" {
}

data "nutanix_clusters_v2" "filtered-cls"{
  filter = "name eq 'cluster-1'"
}

data "nutanix_clusters_v2" "paged-cls" {
  page = 1
  limit = 10
}

```

## Argument Reference

The following arguments are supported:

* `page`: -(Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: -(Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: -(Optional) A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'.
   The filter can be applied to the following fields:
    - backupEligibilityScore
    - clusterProfileExtId
    - config/buildInfo/version
    - config/clusterFunction
    - config/encryptionInTransitStatus
    - config/encryptionOption
    - config/encryptionScope
    - config/hypervisorTypes
    - config/isAvailable
    - extId
    - name
    - network/keyManagementServerType
    - upgradeStatus
* `order_by`: -(Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order.
   The orderby can be applied to the following fields:
    - backupEligibilityScore
    - config/buildInfo/version
    - config/isAvailable
    - extId
    - inefficientVmCount
    - name
    - network/keyManagementServerType
    - nodes/numberOfNodes
    - upgradeStatus
    - vmCount
* `apply`: -(Optional) A URL query parameter that allows clients to specify a sequence of transformations to the entity set, such as groupby, filter, aggregate etc. As of now only support for groupby exists.For example '\$apply=groupby((templateName))' would get all templates grouped by templateName.
   The apply can be applied on the following fields:
   - config/buildInfo/version
   - nodes/numberOfNodes
* `expand`: -(Optional) A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are \$filter, \$select and \$orderby.
   The `expand` can be applied on the following fields:
   - clusterProfile
   - storageSummary
* `select`: -(Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned.
   The select  can be applied to the following fields:
    - backupEligibilityScore
    - inefficientVmCount
    - name
    - upgradeStatus
    - vmCount



## Attribute Reference

The following attributes are exported:

* `cluster_entities`: - List of cluster entities.

### Cluster Entities
The `cluster_entities` contains list of cluster entities. Each cluster entity supports the following:

* `tenant_id`: -  globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: -  A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`: -  Cluster name. This is part of payload for both cluster create & update operations.
* `nodes`: -  Node reference for a cluster.
* `config`: -  Cluster configuration details.
* `network`: -  Network details of a cluster.
* `upgrade_status`: -  Upgrade status of a cluster.
  Valid values are:
    - "CANCELLED"	The cluster upgrade is cancelled.
    - "FAILED"	The cluster upgrade failed.
    - "QUEUED"	The cluster upgrade is in the queue.
    - "SUCCEEDED"	The cluster was upgraded successfully.
    - "DOWNLOADING" The luster upgrade is downloading.
    - "PENDING"The cluster upgrade is in pending state.
    - "UPGRADING" The cluster is in upgrade state.
    - "PREUPGRADE" The cluster is in pre-upgrade state.
    - "SCHEDULED" The cluster upgrade is in scheduled state.
* `vm_count`: -  Number of VMs in the cluster.
* `inefficient_vm_count`: -  Number of inefficient VMs in the cluster.
* `container_name`: -  The name of the default container created as part of cluster creation. This is part of payload for cluster create operation only.
* `categories`: -  List of categories associated to the PE cluster.
* `cluster_profile_ext_id`: -  Cluster profile UUID.
* `backup_eligibility_score`: -  Score to indicate how much cluster is eligible for storing domain manager backup.

### Nodes

The `nodes` attributes supports the following:

* `number_of_nodes`: - Number of nodes in a cluster.
* `node_list`: - List of nodes in a cluster.

#### Node List

The `node_list` attribute supports the following:

* `controller_vm_ip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `node_uuid`: - UUID of the host.
* `host_ip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

### Config
The `config` attributes supports the following:

* `incarnation_id`: - Cluster incarnation Id. This is part of payload for cluster update operation only.
* `build_info`: - Build information details.
* `hypervisor_types`: - Hypervisor types.
    Valid values are:
    - "AHV".
    - "ESX".
    - "HYPERV".
    - "XEN".
    - "NATIVEHOST".
* `cluster_function`: - Cluster function. This is part of payload for cluster
  create operation only (allowed enum values for creation are AOS, ONE_NODE & TWO_NODE only).
  Valid values are:
    - "AOS"
    - "PRISM_CENTRAL"
    - "CLOUD_DATA_GATEWAY"
    - "AFS"
    - "ONE_NODE"
    - "TWO_NODE"
    - "ANALYTICS_PLATFORM"
* `timezone`: - Time zone on a cluster.
* `authorized_public_key_list`: - Public ssh key details. This is part of payload for cluster update operation only.
* `redundancy_factor`: - Redundancy factor of a cluster. This is part of payload for both cluster create & update operations.
* `cluster_software_map`: - Cluster software version details.
* `cluster_arch`: - Cluster arch.
  Valid values are:
    - "PPC64LE" PPC64LE cluster architecture type.
    - "X86_64" X86_64 cluster architecture type.
* `fault_tolerance_state`: - Fault tolerance state of a cluster.
* `is_remote_support_enabled`: - Remote support status.
* `operation_mode`: - Cluster operation mode. This is part of payload for cluster
  update operation only.
  Valid values are:
    - "OVERRIDE"	Override operation mode.
    - "STAND_ALONE"	Stand-alone operation mode.
    - "SWITCH_TO_TWO_NODE"	Switch to two-node operation mode.
    - "NORMAL"	Normal operation mode.
    - "READ_ONLY"	Read-only operation mode.
* `is_lts`: - Indicates whether the release is categorized as Long-term or not.
* `is_password_remote_login_enabled`: - Indicates whether the password ssh into the cluster is enabled or not.
* `encryption_in_transit_status`: - Encryption in transit Status.
  Valid values are:
    - "DISABLED"	Disabled encryption status.
    - "ENABLED" 	Enabled encryption status.
* `encryption_option`: - Encryption option.
  Valid values are:
    - "SOFTWARE".
    - "HARDWARE".
    - "SOFTWARE_AND_HARDWARE"
* `encryption_scope`: - Encryption scope.
  Valid values are:
    - "CLUSTER".
    - "CONTAINER".
* `pulse_status`: - Pulse status for a cluster.
* `is_available`: - Indicates if cluster is available to contact or not.


#### Build info

The build_info attribute supports the following:

* `build_type`: - Software build type.
* `version`: - Software version.
* `full_version`: - Full name of software version.
* `commit_id`: - Commit ID used for version.
* `short_commit_id`: - Short commit Id used for version.

#### Authorized Public Key List

The authorized_public_key_list attribute supports the following:

* `name`: - SSH key name.
* `key`: - SSH key value.

#### Cluster Software Map

The cluster_software_map attribute supports the following:

* `software_type`: - Software type. This is part of payload for cluster create operation only.
  Valid values are:
    - "PRISM_CENTRAL": Prism Central software type.
    - "NOS": NOS software.
    - "NCC": NCC software.
* `version`: - Software version.

#### Fault Tolerance State
The fault_tolerance_state attribute supports the following:

* `current_max_fault_tolerance`: - Maximum fault tolerance that is supported currently.
* `desired_max_fault_tolerance`: - Maximum fault tolerance desired.
* `domain_awareness_level`: - Domain awareness level corresponds to unit of cluster group. This is part of payload for both cluster create & update operations.
  Valid values are:
    - "RACK"	Rack level awareness.
    - "NODE"	Node level awareness.
    - "BLOCK"	Block level awareness.
    - "DISK"	Disk level awareness.

* `current_cluster_fault_tolerance`: - Cluster Fault tolerance. Set desiredClusterFaultTolerance for cluster create and update.
   Valid values are:
    - "CFT_1N_OR_1D":     - System can handle fault of one node or one disk.
    - "CFT_2N_OR_2D":     - System can handle fault of two nodes or two disks.
    - "CFT_1N_AND_1D":    - System can handle fault of one node and one disk on the other node simultaneously.
    - "CFT_0N_AND_0D":    - System can not handle any fault with a node or a disk.

* `desired_cluster_fault_tolerance`: - Cluster Fault tolerance. Set desiredClusterFaultTolerance for cluster create and update.
  Valid values are:
    - "CFT_1N_OR_1D":     - System can handle fault of one node or one disk.
    - "CFT_2N_OR_2D":     - System can handle fault of two nodes or two disks.
    - "CFT_1N_AND_1D":    - System can handle fault of one node and one disk on the other node simultaneously.
    - "CFT_0N_AND_0D":    - System can not handle any fault with a node or a disk.

* `redundancy_status`: - Redundancy Status of the cluster

##### Redundancy Status
The redundancy_status attribute supports the following:

* `is_cassandra_preparation_done`: - Boolean flag to indicate if Cassandra ensemble can meet the desired FT.
* `is_zookeeper_preparation_done`: - Boolean flag to indicate if Zookeeper ensemble can meet the desired FT.


#### Pulse Status
The pulse_status attribute supports the following:
* `is_enabled`: - (Optional) Flag to indicate if pulse is enabled or not.
* `pii_scrubbing_level`: - (Optional) PII scrubbing level.
  Valid values are:
    - "ALL" :	Scrub All PII Information from Pulse including data like entity names and IP addresses.
    - "DEFAULT":	Default PII Scrubbing level. Data like entity names and IP addresses will not be scrubbed from Pulse.

### Network
The `network` attributes supports the following:

* `external_address`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `external_data_services_ip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `external_subnet`: - Cluster external subnet address.
* `internal_subnet`: - Cluster internal subnet address.
* `nfs_subnet_whitelist`: - NFS subnet whitelist addresses. This is part of payload for cluster update operation only.
* `name_server_ip_list`: - List of name servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently
* `ntp_server_ip_list`: - List of NTP servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
* `smtp_server`: - SMTP servers on a cluster. This is part of payload for cluster update operation only.
* `masquerading_ip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `masquerading_port`: - The port to connect to the cluster when using masquerading IP.
* `management_server`: - Management server information.
* `fqdn`: - Cluster fully qualified domain name. This is part of payload for cluster update operation only.
* `key_management_server_type`: - Key management server type.
  Valid values are:
    - "PRISM_CENTRAL"	Prism Central management server.
    - "EXTERNAL"	External management server.
    - "LOCAL"	Local management server.
* `backplane`: - Params associated to the backplane network segmentation. This is part of payload for cluster create operation only.
* `http_proxy_list`: - List of HTTP Proxy server configuration needed to access a cluster which is hosted behind a HTTP Proxy to not reveal its identity.
* `https_proxy_white_list`: - Targets HTTP traffic to which is exempted from going through the configured HTTP Proxy.

#### SMTP Server
The `smtp_server` attribute supports the following:

* `email_address`: - SMTP email address.
* `server`: - SMTP network details.
* `type`: - Type of SMTP server.
  Valid values are:
    - "PLAIN"   	Plain type SMTP server.
    - "STARTTLS"	Start TLS type SMTP server.
    - "SSL" 	    SSL type SMTP server.

##### Server
The `server` attribute supports the following:

* `ip_address`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
* `port`: - SMTP port.
* `username`: - SMTP server user name.
* `password`: - SMTP server password.

#### Management Server
The `management_server` attribute supports the following:

* `ip`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
* `type`: - Type of management server.
  Valid values are:
    - "VCENTER"   	Vcenter management server.
* `is_drs_enabled`: - Indicates whether it is DRS enabled or not.
* `is_registered`: - Indicates whether it is registered or not.
* `is_in_use`: - Indicates whether the host is managed by an entity or not.

#### Backplane
The `backplane` attribute supports the following:

* `is_segmentation_enabled`: - Flag to indicate if the backplane segmentation needs to be enabled or not.
* `vlan_tag`: - VLAN Id tagged to the backplane network on the cluster. This is part of cluster create payload.
* `subnet`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `netmask`: - An unique address that identifies a device on the internet or a local network in IPv4 format.

##### Subnet, Netmask
The `subnet`, `netmask` attributes supports the following:

* `value`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format.
* `prefix_length`: - The prefix length of the network to which this host IPv4/IPv6 address belongs.

#### Http Proxy List
The `http_proxy_list` attribute supports the following:

* `ip_address`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
* `port`: - SMTP port.
* `username`: - SMTP server user name.
* `password`: - SMTP server password.
* `name`: - HTTP Proxy server name configuration needed to access a cluster which is hosted behind a HTTP Proxy to not reveal its identity.
* `proxy_type`: - Type of HTTP Proxy server.
  Valid values are:
    - "HTTP".
    - "HTTPS".
    - "SOCKS".

#### Https Proxy White List
The `https_proxy_white_list` attribute supports the following:

* `target_type`: - Type of the target which is exempted from going through the configured HTTP Proxy.
  Valid values are:
    - "IPV6_ADDRESS"	IPV6 address.
    - "HOST_NAME"	Name of the host.
    - "DOMAIN_NAME_SUFFIX" Domain Name Suffix required for http proxy whitelist.
    - "IPV4_NETWORK_MASK" Network Mask of the IpV4 family.
    - "IPV4_ADDRESS" IPV4 address.
* `target`: - Target's identifier which is exempted from going through the configured HTTP Proxy.

### Ip Address Attributes
The `nodes.host_ip`, `nodes.controller_vm_ip`, `network.external_address`,
`network.external_data_services_ip`, `network.smtp_server.server.ip_address`,
`network.management_server.ip`

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.

### Ip Address and FQDN Attributes
The `network.name_server_ip_list`, `network.ntp_server_ip_list` attributes supports the following:
#### IPV4, IPV6
The `ipv4`, `ipv6` attributes supports the following:

* `value`: - The IPv4/IPv6 address of the host.
* `prefix_length`: - The prefix length of the network to which this host IPv4/IPv6 address belongs.

#### FQDN
The `fqdn` attribute supports the following:

* `value`: - The fully qualified domain name of the host.

See detailed information in [Nutanix List Clusters V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0#tag/Clusters/operation/listClusters).
