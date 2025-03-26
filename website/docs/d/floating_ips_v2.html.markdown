---
layout: "nutanix"
page_title: "NUTANIX: nutanix_floating_ips_v2"
sidebar_current: "docs-nutanix-datasource-floating-ips-v2"
description: |-
  Provides a datasource to retrieve floating ip with floating_ip_uuid.
---

# nutanix_floating_ips_v2
Provides a datasource to retrieve floating IP with floating_ip_uuid .

## Example Usage

```hcl
data "nutanix_floating_ips_v2" "floating-ips"{}

data "nutanix_floating_ips_v2" "floating-ips-filter"{
    filter = "name eq 'floating_ip_example'"
}

data "nutanix_floating_ips_v2" "floating-ips-limit"{
    limit = 10
}

data "nutanix_floating_ips_v2" "floating-ips-filter-limit"{
    filter = "name eq 'floating_ip_example'"
    limit = 10
}
```

## Argument Reference
The following arguments are supported:

- `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources. The filter can be applied to the following fields:
  - externalSubnetReference
  - floatingIp/ipv4/value
  - floatingIp/ipv6/value
  - loadBalancerSessionReference
  - name
- `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
  - floatingIp/ipv4/value
  - floatingIp/ipv6/value
  - name
- `expand`: (Optional) A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. The expand can be applied to the following fields:
  - externalSubnet
  - vpc
  - vmNic


## Attribute Reference
The following attributes are exported:

- `floating_ips`: List of all Floating IPs.

### Floating IPs
The `floating_ips` object contains the following attributes:

- `ext_id`: Floating IP UUID
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



See detailed information in [Nutanix List Floating IPs v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/FloatingIps/operation/listFloatingIps).
