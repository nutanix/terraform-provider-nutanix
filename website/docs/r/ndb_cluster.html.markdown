---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_cluster"
sidebar_current: "docs-nutanix-resource-ndb-cluster"
description: |-
  This operation submits a request to add a Nutanix cluster to Nutanix database service (NDB).
---

# nutanix_ndb_cluster

Provides a resource to add a Nutanix cluster based on the input parameters. 

## Example Usage

```hcl
    resource "nutanix_ndb_clusters" "clsname" {
        name= "{{ test-cluster }}"
        description = "test description"
        cluster_ip = "{{ cluster_ip }}"
        username= "{{ username of cluster }}"
        password = "{{ password of cluster }}"
        storage_container = "{{ storage_container }}"
        agent_network_info{
            dns = "{{ DNS servers available in the }}"
            ntp = "{{ NTP servers available }}"
        }
        networks_info{
            type = "DHCP"
            network_info{
                vlan_name = "vlan_static"
                static_ip = "{{ static_ip }}"
                gateway = "{{ gateway }}"
                subnet_mask="{{ subnet_mask }}"
            }
            access_type = [
                "PRISM",
                "DSIP",
                "DBSERVER"
            ]
        }
    }
```

## Argument Reference

* `name`: (Required) name of the cluster to be registered
* `description`: (Optional) description of cluster
* `cluster_ip`: (Required) Prism Element IP address
* `username`: (Required) username of the Prism Element administrator
* `password`: (Required) Prism Element password
* `storage_container`: (Required) select a storage container which is used for performing database operations in the cluster
* `agent_network_info`: (Required) agent network info to register cluster 
* `networks_info`: (Required) network segmentation to segment the network traffic of the agent VM.


### agent_network_info
* `dns` : string of DNS servers(comma separted).
* `ntp`: string of NTP servers(comma separted).

### networks_info
* `type`: type of vlan. Supported [DHCP, Static, IPAM]
* `network_info`: network segmentation to segment the network traffic
* `access_type`: VLAN access types for which you want to configure network segmentation. Supports [PRISM, DSIP, DBSERVER ]. 
Prism Element: Select this VLAN access type to configure a VLAN that the NDB agent VM can use to communicate with Prism.
Prism iSCSI Data Service. Select this VLAN access type to configure a VLAN that the agent VM can use to make connection requests to the iSCSI data services IP.
DBServer Access from NDB server. Select this VLAN access type to configure a VLAN that is used for communications between the NDB agent VM and the database server VM on the newly registered NDB server cluster.

### network_info
* `vlan_name`: vlan name
* `static_ip`: static ip of agent network
* `gateway`: gateway of agent network
* `subnet_mask`: subnet mask of agent network



## Attributes Reference
The following attributes are exported:

* `id`: - id of cluster
* `name`: - name of cluster
* `unique_name`: - unique name of cluster
* `ip_addresses`: - IP address
* `fqdns`: - fqdn
* `nx_cluster_uuid`: - nutanix cluster uuid
* `description`: - description
* `cloud_type`: - cloud type
* `date_created`: - creation date
* `date_modified`: - date modified
* `version`: - version
* `owner_id`: - owner UUID
* `status`: - current status
* `hypervisor_type`: - hypervisor type
* `hypervisor_version`: - hypervisor version
* `properties`: - list of properties
* `reference_count`: - NA
* `username`: - username 
* `password`: - password
* `cloud_info`: - cloud info
* `resource_config`: - resource related consumption info
* `management_server_info`: - NA
* `entity_counts`: - no. of entities related 
* `healthy`: - if healthy status

See detailed information in [NDB Cluster](https://www.nutanix.dev/api_references/ndb/#/1f392bec2e58b-update-the-given-cluster).
