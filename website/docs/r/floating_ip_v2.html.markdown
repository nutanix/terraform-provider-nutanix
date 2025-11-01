---
layout: "nutanix"
page_title: "NUTANIX: nutanix_floating_ip_v2"
sidebar_current: "docs-nutanix-resource-floating-ip-v2"
description: |-
  Create Floating IPs .
---

# nutanix_floating_ip_v2

Provides Nutanix resource to create Floating IPs.

##  Example1 :  create Floating IP with External Subnet

```hcl

# create Floating IP with External Subnet UUID
resource "nutanix_floating_ip_v2" "fip-ext-subnet"{
  name                      = "example-fip"
  description               = "example fip  description"
  external_subnet_reference = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
}

```

## Example2 :  create Floating IP with External Subnet with vm association

```hcl
resource "nutanix_floating_ip_v2" "fip-ext-subnet-vm"{
  name                      = "example-fip"
  description               = "example fip  description"
  external_subnet_reference = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
  association {
    vm_nic_association {
      vm_nic_reference = "31e4b3b1-4b3b-4b3b-4b3b-4b3b4b3b4b3b"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `name`: (Required) Name of the floating IP.
- `description`: (Optional) Description for the Floating IP.
- `association`: (Optional) Association of the Floating IP with either NIC or Private IP
- `floating_ip`: (Optional) Floating IP address.
- `external_subnet_reference`: (Optional) External subnet reference for the Floating IP to be allocated in on-prem only.
- `vpc_reference`: (Optional) VPC reference UUID
- `vm_nic_reference`: (Optional) VM NIC reference.

### association

- `vm_nic_association`: (Optional) Association of Floating IP with nic
- `vm_nic_association.vm_nic_reference`: (Required) VM NIC reference.
- `vm_nic_association.vpc_reference`: (Optional) VPC reference to which the VM NIC subnet belongs.

- `private_ip_association`: (Optional) Association of Floating IP with private IP
- `private_ip_association.vpc_reference`: (Required) VPC in which the private IP exists.
- `private_ip_association.private_ip`: (Required) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

### floating_ip

- `ipv4`: Reference to IP Configuration
- `ipv6`: Reference to IP Configuration

### ipv4, ipv6 (Reference to IP Configuration)

- `value`: value of address
- `prefix_length`: Prefix length of the network to which this host IPv4 address belongs. Default value is 32.

## Attributes Reference

The following attributes are exported:

- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `metadata`: Metadata associated with this resource.
- `association_status`: Association status of floating IP.
- `external_subnet`: Networking common base object
- `vpc`: Networking common base object
- `vm_nic`: Virtual NIC for projections

## Import

This helps to manage existing entities which are not created through terraform. Floating IPs can be imported using the `UUID`. (ext_id in v4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_floating_ip_v2" "floating_ip"{}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_floating_ips_v2" "fetch_fips"{}
terraform import nutanix_floating_ip_v2.floating_ip <UUID>
```

See detailed information in [Nutanix Floating IP v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/FloatingIps/operation/createFloatingIp).
