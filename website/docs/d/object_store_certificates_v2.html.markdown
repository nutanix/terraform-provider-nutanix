---
layout: "nutanix"
page_title: "NUTANIX: nutanix_certificates_v2 "
sidebar_current: "docs-nutanix-datasource-object-store-certificates-v2"
description: |-
  Get a list of the SSL certificates which can be used to access an Object store.


---

# nutanix_certificates_v2

Get a list of the SSL certificates which can be used to access an Object store.




## Example Usage

```hcl
data "nutanix_certificates_v2" "example"{
  object_store_ext_id = "ac91151a-28b4-4ffe-b150-6bcb2ec80cd4"
}

```

## Argument Reference

The following arguments are supported:

- `object_store_ext_id`: -(Required) The UUID of the Object store.
* `page`: -(Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results. Default value is 0.
* `limit`: -(Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set. Default value is 50.
* `filter`: -(Optional) A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. The filter can be applied to the following fields:
    - alternateFqdns/value
    - alternateIps/ipv4/value
* `select`: -(Optional)  URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned.
    - alternateFqdns
    - alternateFqdns/value
    - alternateIps
    - alternateIps/ipv4/value

## Attributes Reference

The following attributes are exported:

- `certificates`: - list of the SSL certificates which can be used to access an Object store.

### Certificates
The `certificates` argument exports the following:

- `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `metadata`: - Metadata associated with this resource.
- `alternate_fqdns`: - The list of alternate FQDNs for accessing the Object store. The FQDNs must consist of at least 2 parts separated by a '.'. Each part can contain upper and lower case letters, digits, hyphens or underscores but must begin and end with a letter. Each part can be up to 63 characters long. For e.g 'objects-0.pc_nutanix.com'.
- `alternate_ips`: - A list of the IPs included as Subject Alternative Names (SANs) in the certificate. The IPs must be among the public IPs of the Object store (publicNetworkIps).

### Links
The `links` argument exports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Metadata
The `metadata` argument exports the following:

- `owner_reference_id`: - A globally unique identifier that represents the owner of this resource.
- `owner_user_name`: - The userName of the owner of this resource.
- `project_reference_id`: - A globally unique identifier that represents the project this resource belongs to.
- `project_name`: - The name of the project this resource belongs to.
- `category_ids`: - A list of globally unique identifiers that represent all the categories the resource is associated with.


### Alternate FQDNs
The `alternate_fqdns` argument exports the following:

- `value`: - The fully qualified domain name of the host.

### Alternate IPs
The `alternate_ips` argument exports the following:
- `ipv4`: An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv6`: An unique address that identifies a device on the internet or a local network in IPv6 format.


### IPv4, IPv6
The `ipv4` and `ipv6` argument exports the following:
- `value`: - The IPv4/IPv6 address of the host.
- `prefix_length`: - The prefix length of the network to which this host IPv4 address belongs. Default for IPv4 is 32 and for IPv6 is 128.

See detailed information in [Nutanix Get a list of the SSL certificates of an Object store V4 ](https://developers.nutanix.com/api-reference?namespace=objects&version=v4.0#tag/ObjectStores/operation/listCertificatesByObjectstoreId).
