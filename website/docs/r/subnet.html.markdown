---
layout: "nutanix"
page_title: "NUTANIX: nutanix_subnet"
sidebar_current: "docs-nutanix-resource-subnet"
description: |-
  This operation submits a request to create a subnet based on the input parameters. A subnet is a block of IP addresses.
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
  value = data.nutanix_clusters.clusters.entities.0.metadata.uuid
}

resource "nutanix_subnet" "next-iac-managed" {
  # What cluster will this VLAN live on?
  cluster_uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"

  # General Information
  name        = "next-iac-managed-example"
  vlan_id     = 101
  subnet_type = "VLAN"

  # Managed L3 Networks
  # This bit is only needed if you intend to turn on IPAM
  prefix_length = 20

  default_gateway_ip = "10.5.80.1"
  subnet_ip          = "10.5.80.0"

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "eng.nutanix.com"]
}
```

## Argument Reference

* `metadata`: - (Required) The subnet kind metadata.
* `availability_zone_reference`: - (Optional) The reference to a availability_zone.
* `cluster_uuid`: - (Required) The UUID of the cluster.
* `description`: - (Optional) A description for subnet.
* `name`: - (Optional) Subnet name (Readonly).
* `categories`: - (Optional) The categories of the resource.
* `owner_reference`: - (Optional) The reference to a user.
* `project_reference`: - (Optional) The reference to a project.
* `vswitch_name`: - (Optional).
* `subnet_type`: - (Optional). Valid Types are ["VLAN", "OVERLAY"]
* `default_gateway_ip`: - (Optional) Default gateway IP address.
* `prefix_length`: - (Optional).
* `subnet_ip`: - (Optional) Subnet IP address.
* `dhcp_server_address`: - (Optional) Host address.
* `dhcp_server_address_port`: - (Optional) Port Number.
* `ip_config_pool_list_ranges`: -(Optional) Range of IPs.
* `dhcp_options`: - (Optional) Spec for defining DHCP options.
* `dhcp_domain_search_list`: - (Optional).The DNS domain search list .
* `dhcp_domain_name_server_list`: - (Optional). List of Domain Name Server addresses .
* `vlan_id`: - (Optional). For VLAN subnet.
* `network_function_chain_reference`: - (Optional) The reference to a network_function_chain.
* `vpc_reference_uuid`: (Optional) VPC reference uuid
* `is_external`: - (Optional) Whether the subnet is external subnet or not.
* `enable_nat`: - (Optional) Whether NAT should be performed for VPCs attaching to the subnet. This field is supported only for external subnets. NAT is enabled by default on external subnets.

## Attributes Reference

The following attributes are exported:

* `metadata`: - The vm kind metadata.
* `state`: - The state of the subnet.
* `api_version` - The version of the API.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `uuid`: - The subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `network_function_chain_reference`, `subnet_reference`.

attributes supports the following:

* `kind`: - The kind name (Default value: project)(Required).
* `name`: - the name(Optional).
* `uuid`: - the UUID(Required).

Note: `subnet_reference` does not support the attribute `name`

See detailed information in [Nutanix Subnet](https://www.nutanix.dev/api_references/prism-central-v3/#/0cc5a30420b29-create-a-new-subnet).
