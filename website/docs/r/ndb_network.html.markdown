---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_network"
sidebar_current: "docs-nutanix-resource-ndb-network"
description: |-
  This operation submits a request to create, update and delete networks in Nutanix database service (NDB).
---

# nutanix_ndb_network

Provides a resource to create VLANs and IP address pools that are managed both in NDB and outside NDB. 

## Example Usage

### resource to create network for NDB
```hcl
    resource "nutanix_ndb_network" "name" {
        name= "test-sub"
        type="Static"
        cluster_id = "{{ cluster_id }}"
        gateway= "{{ gatway for the vlan }}"
        subnet_mask = "{{ subnet mask for the vlan}}"
        primary_dns = " {{ primary dns for the vlan }}"
        secondary_dns= "{{secondary dns for the vlan }}"
        ip_pools{
            start_ip = "{{ starting address range}}"
            end_ip = "{{ ending address range }}"
        }
    }   
```

### resource to create network for NDB with dns domain
```hcl
    resource "nutanix_ndb_network" "name" {
        name= "test-sub"
        type="Static"
        cluster_id = "{{ cluster_id }}"
        gateway= "{{ gatway for the vlan }}"
        subnet_mask = "{{ subnet mask for the vlan}}"
        primary_dns = " {{ primary dns for the vlan }}"
        secondary_dns= "{{secondary dns for the vlan }}"
        ip_pools{
            start_ip = "{{ starting address range}}"
            end_ip = "{{ ending address range }}"
        }
        dns_domain = {{ dns domain }}
    }   
```

## Argument Reference
* `name`: (Required) Name of the vlan to be attached in NDB
* `type`: (Required) Vlan type. Supports [DHCP, Static]
* `cluster_id`: (Required) Select the Nutanix cluster on which you want to add the VLAN.
* `ip_pools`: (Optional) Manage IP Address Pool in NDB option if you want to assign static IP addresses to your database server VMs
* `gateway`: (Optional) Gateway for vlan. Supports in Static IP address assignment only 
* `subnet_mask`: (Optional) Subnet mask for vlan. (Static IP address assignment only)
* `primary_dns`: (Optional) primary dns for vlan. (Static IP address assignment only)
* `secondary_dns`: (Optional) secondary dns for vlan. (Static IP address assignment only)
* `dns_domain`: (Optional) dns domain for vlan. (Static IP address assignment only)

### ip_pools
* `start_ip`: (Required) starting IP address range for new database servers 
* `end_ip`: (Required) ending IP address range for new database servers 


## Attributes Reference

* `managed`: Managed by NDB or not
* `stretched_vlan_id`: stretched vlan id
* `properties`: properties of network
* `properties_map`: properties map of network


See detailed information in [NDB Network](https://www.nutanix.dev/api_references/ndb/#/4a4fc22c2843d-add-a-v-lan-to-ndb).