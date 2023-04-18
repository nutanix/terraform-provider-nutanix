---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_networks"
sidebar_current: "docs-nutanix-datasource-ndb-networks"
description: |-
 List of networks in Nutanix Database Service
---

# nutanix_ndb_networks

 List of networks in Nutanix Database Service

## Example Usage

```hcl
    data "nutanix_ndb_networks" "nw" { }
```

## Attribute Reference
The following attributes are exported:

* `networks`: List of networks in NDB

### networks

* `id`: network id
* `name`: network name
* `managed`: network managed by NDB or not
* `type`: type of network
* `cluster_id`: cluster id where network is present
* `stretched_vlan_id`: stretched vlan id
* `properties`: properties of network
* `properties_map`: properties map of network
* `ip_addresses`: IP addresses of network
* `ip_pools`: IP Pools of network

### ip_addresses
* `ip`: ip
* `status`: status of ip
* `dbserver_id`: dbserver id
* `dbserver_name`: dbserver name

### ip_pools
* `start_ip`: start ip 
* `end_ip`: end ip
* `addresses`: address of ips ranges

### addresses
* `ip`: ip of pool
* `status`: ip status

### properties_map
* `vlan_subnet_mask`: subnet mask of vlan
* `vlan_primary_dns`: primary dns of vlan
* `vlan_secondary_dns`: secondary dns of vlan
* `vlan_gateway`: gateway of vlan


See detailed information in [List of NDB Networks](https://www.nutanix.dev/api_references/ndb/#/283556b78730b-get-vlans).