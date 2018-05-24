---
layout: "nutanix"
page_title: "NUTANIX: nutanix_subnet"
sidebar_current: "docs-nutanix-datasource-subnet"
description: |-
  This operation retrieves a subnet based on the input parameters. A subnet is a block of IP addresses.
---

# nutanix_subnet

Provides a resource to create a subnet based on the input parameters. A subnet is a block of IP addresses.

## Example Usage

```hcl
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"

	cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
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

data "nutanix_subnet" "test" {
	subnet_id = "${nutanix_subnet.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `subnet_id`: - (Required) The name for the subnet.

## Attributes Reference

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
* `prefix_length`: - (Optional).
* `subnet_ip`: - (Optional) Subnet IP address.
* `dhcp_server_address`: - (Optional) Host address.
* `dhcp_server_address_port`: - (Optional) Port Number.
* `dhcp_options`: - (Optional) Spec for defining DHCP options.
* `dhcp_domain_search_list`: - (Optional).
* `vlan_id`: - (Optional).
* `network_function_chain_reference`: - (Optional) The reference to a network_function_chain.

## Attributes Reference

The following attributes are exported:

* `metadata`: - The subnet kind metadata.
* `state`: -

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `uuid`: - subnet uuid.
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
* `uuid`: - the uuid(Required).

Note: `cluster_reference`, `subnet_reference` does not support the attribute `name`

See detailed information in [Nutanix Subnet](http://developer.nutanix.com/reference/prism_central/v3/#definitions-subnet_resources).
