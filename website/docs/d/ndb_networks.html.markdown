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

### properties_map
* `vlan_subnet_mask`: subnet mask of vlan
* `vlan_primary_dns`: primary dns of vlan
* `vlan_secondary_dns`: secondary dns of vlan
* `vlan_gateway`: gateway of vlan