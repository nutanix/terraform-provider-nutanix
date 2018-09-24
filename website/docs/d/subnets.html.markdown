---
layout: "nutanix"
page_title: "NUTANIX: nutanix_subnets"
sidebar_current: "docs-nutanix-datasource-images"
description: |-
 Describes a List of Subnets
---

# nutanix_subnets

Describes a List of Subnets

## Example Usage

```hcl
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.UUID}"
}

resource "nutanix_subnet" "test" {
    name = "dou_vlan0_test_%d"

    cluster_reference = {
      kind = "cluster"
      UUID = "${data.nutanix_clusters.clusters.entities.0.metadata.UUID}"
    }

    vlan_id = 201
    subnet_type = "VLAN"

    prefix_length = 24
    default_gateway_ip = "192.168.0.1"
    subnet_ip = "192.168.0.0"
    #ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

    dhcp_options {
        boot_file_name = "bootfile"
        tftp_server_name = "192.168.0.252"
        domain_name = "nutanix"
    }

    dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
    dhcp_domain_search_list = ["nutanix.com", "calm.io"]

}

data "nutanix_subnets" "test" {}
```

## Argument Reference

The following arguments are supported:

* `metadata`: Represents virtual machine UUID

### Metadata Argument

The metadata attribute supports the following:

* `kind`: - The kind name.
* `sort_attribute`: The attribute to perform sort on.
* `filter`: - The filter in FIQL syntax used for the results.
* `length`: - The number of records to retrieve relative to the offset.
* `sort_order`: - The sort order in which results are returned
* `offset`: - Offset from the start of the entity list

## Attribute Reference

The following attributes are exported:

* `entities`: - A list of virtual machines.

### Entities Attribute

The entities attribute supports the following:

* `metadata`: - (Required) The subnet kind metadata.
* `availability_zone_reference`: - (Optional) The reference to a availability_zone.
* `cluster_reference`: - (Optional) The reference to a cluster.
* `cluster_name`: - (Optional) The name of a cluster.
* `description`: - (Optional) A description for subnet.
* `name`: - (Optional) Subnet name (Readonly).
* `categories`: - (Optional) The API Version.
* `owner_reference`: - (Optional) The reference to a user.
* `project_reference`: - (Optional) The reference to a project.
* `vswitch_name`: - (Optional).
* `subnet_type`: - (Optional).
* `default_gateway_ip`: - (Optional) Default gateway IP address.
* `prefix_length`: - (Optional). IP prefix length of the subnet.
* `subnet_ip`: - (Optional) Subnet IP address.
* `dhcp_server_address`: - (Optional) Host address.
* `dhcp_server_address_port`: - (Optional) Port Number.
* `dhcp_options`: - (Optional) Spec for defining DHCP options.
* `dhcp_domain_search_list`: - (Optional).
* `vlan_id`: - (Optional). The VLAN ID of the subnet.
* `network_function_chain_reference`: - (Optional) The reference to a network_function_chain.
* `state`: - The state of the subnet.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: - subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, `network_function_chain_reference`, `subnet_reference`.

attributes supports the following:

* `kind`: - The kind name (Default value: project)(Required).
* `name`: - the name(Optional).
* `UUID`: - the UUID(Required).

Note: `cluster_reference`, `subnet_reference` does not support the attribute `name`

See detailed information in [Nutanix Subnet](http://developer.nutanix.com/reference/prism_central/v3/#definitions-subnet_resources).