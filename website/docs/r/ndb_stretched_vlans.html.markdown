---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_stretched_vlan"
sidebar_current: "docs-nutanix-resource-ndb-stretched-vlan"
description: |-
   This operation submits a request to create, update and delete stretched vlans in Nutanix database service (NDB). We can add a stretched VLAN to NDB by selecting the existing VLANs from each Nutanix cluster.
---

# nutanix_ndb_stretched_vlan

Provides a resource to create stretched vlans based on the input parameters. 

## Example Usage

### resource to add stretched vlan in NDB
```hcl
    resource "nutanix_ndb_stretched_vlan" "name" {
        name = "test-stretcName"
        description = "vlan desc updated"
        type = "Static"
        vlan_ids = [
            "{{ vlan_id_1 }}",
            "{{ vlan_id_2 }}"
        ]
    }
```

### resource to update the strteched vlan with new gateway and subnet mask
```hcl
    resource "nutanix_ndb_stretched_vlan" "name" {
        name = "test-stretcName"
        description = "vlan desc updated"
        type = "Static"
        vlan_ids = [
            "{{ vlan_id_1 }}",
            "{{ vlan_id_2 }}"
        ]
        metadata{
            gateway = "{{ gateway of vlans }}"
            subnet_mask = "{{ subnet mask of vlans }}"
        }
    }
```


## Argument Reference

* `name`: (Required) name for the stretched VLAN
* `description`: (Optional) Description of stretched vlan
* `type`: (Required) type of vlan. static VLANs that are managed in NDB can be added to a stretched VLAN. 
* `vlan_ids`: (Required) list of vlan ids to be added in NDB

* `metadata`: (Optional) Update the stretched VLAN Gateway and Subnet Mask IP address
* `metadata.gateway`: Update the gateway of stretched vlan
* `metadata.subnet_mask`: Update the subnet_mask of stretched vlan

## Attributes Reference
The following attributes are exported:

* `vlans_list`: properties of vlans

### vlans_list
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