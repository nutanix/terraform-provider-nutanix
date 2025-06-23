---
layout: "nutanix"
page_title: "NUTANIX: nutanix_certificate_v2 "
sidebar_current: "docs-nutanix-datasource-object-store-certificate-v2"
description: |-
  Get the details of the SSL certificate which can be used to connect to an Object store.


---

# nutanix_certificate_v2

Get the details of the SSL certificate which can be used to connect to an Object store.


## Example Usage

```hcl
data "nutanix_certificate_v2" "example"{
  object_store_ext_id = "ac91151a-28b4-4ffe-b150-6bcb2ec80cd4"
  ext_id              = "ef0a9a54-e7e1-42e2-a59f-de779ec1c9ea"
}

```

## Argument Reference

The following arguments are supported:

- `object_store_ext_id`: -(Required) The UUID of the Object store.
- `ext_id`: -(Required) The UUID of the certificate of an Object store.

## Attributes Reference

The following attributes are exported:

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

See detailed information in [Nutanix Get the details of an Object store certificate V4 ](https://developers.nutanix.com/api-reference?namespace=objects&version=v4.0#tag/ObjectStores/operation/getCertificateById).
