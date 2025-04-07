---
layout: "nutanix"
page_title: "NUTANIX: nutanix_floating_ip_v2"
sidebar_current: "docs-nutanix-datasource-floating-ip-v2"
description: |-
  Provides a datasource to retrieve the floating IP for this extId.
---

# nutanix_floating_ip_v2

Provides a datasource to retrieve the floating IP for this extId.

## Example Usage

```hcl
data "nutanix_floating_ip_v2" "get-floating-ip"{
  ext_id ="cf1a7643-605e-4622-9f2a-b3fdccdbd072"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id` - (Required) Floating IP UUID

## Attribute Reference

The following attributes are exported:

- `name`: Name of the floating IP.
- `description`: Description for the Floating IP.
- `association`: Association of the Floating IP with either NIC or Private IP
- `floating_ip`: Floating IP address.
- `external_subnet_reference`: External subnet reference for the Floating IP to be allocated in on-prem only.
- `external_subnet`: Networking common base object
- `private_ip`: Private IP value in string
- `floating_ip_value`: Floating IP value in string
- `association_status`: Association status of floating IP.
- `vpc_reference`: VPC reference UUID
- `vm_nic_reference`: VM NIC reference.
- `vpc`: Networking common base object
- `vm_nic`: Virtual NIC for projections
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `metadata`: Metadata associated with this resource.

### association

- `vm_nic_association`: Association of Floating IP with nic
- `vm_nic_association.vm_nic_reference`: VM NIC reference.
- `vm_nic_association.vpc_reference`: VPC reference to which the VM NIC subnet belongs.

- `private_ip_association`: Association of Floating IP with private IP
- `private_ip_association.vpc_reference`: VPC in which the private IP exists.
- `private_ip_association.private_ip`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

### floating_ip

- `ipv4`: Reference to IP Configuration
- `ipv6`: Reference to IP Configuration

### ipv4, ipv6 (Reference to IP Configuration)

- `value`: value of address
- `prefix_length`: Prefix length of the network to which this host IPv4 address belongs. Default value is 32.

### metadata

The `metadata` object contains the following attributes:

- `owner_reference_id` : A globally unique identifier that represents the owner of this resource.
- `owner_user_name` : The userName of the owner of this resource.
- `project_reference_id` : A globally unique identifier that represents the project this resource belongs to.
- `project_name` : The name of the project this resource belongs to.
- `category_ids` : A list of globally unique identifiers that represent all the categories the resource is associated with.

### Links

The `links` attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.


See detailed information in [Nutanix Get Floating IP v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/FloatingIps/operation/getFloatingIpById).
