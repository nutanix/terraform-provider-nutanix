---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_network_available_ips"
sidebar_current: "docs-nutanix-datasource-ndb-network-available-ips"
description: |-
 List of available IPs in Network
---

# nutanix_ndb_network_available_ips

 List of available IPs in Network

## Example Usage

```hcl
    data "nutanix_ndb_network_available_ips" "network"{ 
        profile_id = "{{ network_profile_id }}"
    } 
```


## Attribute Reference

The following attributes are exported:

* `profile_id`: (Required) Network Profile id.
* `available_ips`: List of network available ips

### available_ips

* `id`: network profile id
* `name`: Network Name
* `property_name`: property name of vlan
* `type`: type of network 
* `managed`: managed by ndb or not
* `ip_addresses`: list of available ips in network
* `cluster_id`: cluster id
* `cluster_name`: cluster name