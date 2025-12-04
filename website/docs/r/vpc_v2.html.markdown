---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vpc_v2"
sidebar_current: "docs-nutanix-resource-vpc-v2"
description: |-
  Create Virtual Private Cloud .
---

# nutanix_vpc_v2

Provides Nutanix resource to create VPC.

## Example

```hcl
resource "nutanix_vpc_v2" "vpc" {
  name        = "vpc-example"
  description = "VPC for example"
  external_subnets {
    subnet_reference = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
  }
}

# creating VPC with external routable prefixes
resource "nutanix_vpc_v2" "external-vpc-routable-vpc" {
  name        = "tf-vpc-example"
  description = "VPC "
  external_subnets {
    subnet_reference = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
    external_ips {
      ipv4 {
        value         = "192.168.0.24"
        prefix_length = 32
      }
    }
    external_ips {
      ipv4 {
        value         = "192.168.0.25"
        prefix_length = 32
      }
    }
  }
  externally_routable_prefixes {
    ipv4 {
      ip {
        value         = "172.30.0.0"
        prefix_length = 32
      }
      prefix_length = 16
    }
  }
}

// creating VPC with transit type
resource "nutanix_vpc_v2" "transit-vpc" {
  name        = "vpc-transit"
  description = "VPC for transit type"
  external_subnets {
    subnet_reference = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
  }
  vpc_type = "TRANSIT"
}

```

## Argument Reference

The following arguments are supported:

- `name`: (Required) Name of the VPC.
- `description`: (Optional) Description of the VPC.
- `vpc_type`: (Optional) Type of VPC. Acceptable values are "REGULAR" , "TRANSIT".
- `common_dhcp_options`: (Optional) List of DHCP options to be configured.
- `external_subnets`: (Optional) List of external subnets that the VPC is attached to.
- `external_routing_domain_reference`: (Optional) External routing domain associated with this route table
- `externally_routable_prefixes`: (Optional) CIDR blocks from the VPC which can talk externally without performing NAT. This is applicable when connecting to external subnets which have disabled NAT.

### common_dhcp_options

- `domain_name_servers`: (Optional) List of Domain Name Server addresses
- `domain_name_servers.ipv4`:(Optional) Reference to address configuration
- `domain_name_servers.ipv6`: (Optional) Reference to address configuration

### external_subnets

- `subnet_reference`: (Required) External subnet reference.
- `external_ips`: (Optional) List of IP Addresses used for SNAT, if NAT is enabled on the external subnet. If NAT is not enabled, this specifies the IP address of the VPC port connected to the external gateway.
- `gateway_nodes`: (Optional) List of gateway nodes that can be used for external connectivity.
- `active_gateway_node`: (Optional) Maximum number of active gateway nodes for the VPC external subnet association.


### external_ips

- `ipv4`: (Optional) Reference to address configuration
- `ipv6`: (Optional) Reference to address configuration

### externally_routable_prefixes

- `ipv4`: (Optional) IP V4 Configuration
- `ipv4.ip`: (Required) Reference to address configuration
- `ipv4.prefix_length`: (Required) The prefix length of the network

- `ipv6`: (Optional) IP V6 Configuration
- `ipv6.ip`: (Required) Reference to address configuration
- `ipv6.prefix_length`: (Required) The prefix length of the network

## Attributes Reference

The following attributes are exported:

- `ext_id`: the vpc uuid.
- `metadata`: The vpc kind metadata.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

## Import

This helps to manage existing entities which are not created through terraform. VPC can be imported using the `UUID`. (ext_id in v4 terms). eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_vpc_v2" "import_vpc" {}

// execute this command in cli
terraform import nutanix_vpc_v2.import_vpc <UUID>
```

See detailed information in [Nutanix VPC v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0).
