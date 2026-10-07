---
layout: "nutanix"
page_title: "NUTANIX: nutanix_snmp_user_v2"
sidebar_current: "docs-nutanix-resource-snmp-user-v2"
description: |-
  Adds SNMP user configuration to the cluster identified by {clusterExtId}.
---

# nutanix_snmp_user_v2

Adds SNMP user configuration to the cluster identified by {clusterExtId}.
The resource manages the full lifecycle of an SNMP user: create, read, update,
and delete.

## Example Usage

```hcl
resource "nutanix_snmp_user_v2" "example" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  username       = "tf-snmp-user"
  auth_type      = "MD5"
  auth_key       = "auth-key-12345678"
  priv_type      = "DES"
  priv_key       = "priv-key-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: - (Required, ForceNew) Indicates the UUID of a cluster.
* `username`: - (Required, ForceNew) SNMP username. For SNMP trap v3 version,
  SNMP username is a required parameter. The upstream API does not support
  renaming an existing SNMP user — changing this attribute forces resource
  replacement.
* `auth_type`: - (Required) SNMP user authentication type. One of `MD5`, `SHA`, `SHA224`, `SHA256`, `SHA384`, `SHA512`.
* `auth_key`: - (Required, Sensitive) SNMP user authentication key.
* `priv_type`: - (Optional, Computed) SNMP user encryption type. One of `DES`, `AES`, `AES192`, `AES256`.
* `priv_key`: - (Optional, Computed, Sensitive) SNMP user encryption key.

## Attribute Reference

The following attributes are exported in addition to the arguments above:

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

### Links

The `links` block exports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix Create SNMP User](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.3#tag/Clusters/operation/createSnmpUser).
