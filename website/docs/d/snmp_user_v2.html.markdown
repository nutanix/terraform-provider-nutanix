---
layout: "nutanix"
page_title: "NUTANIX: nutanix_snmp_user_v2"
sidebar_current: "docs-nutanix-datasource-snmp-user-v2"
description: |-
  Fetches SNMP user configuration details identified by {extId} associated with the cluster identified by {clusterExtId}.
---

# nutanix_snmp_user_v2

Fetches SNMP user configuration details identified by {extId} associated with the cluster identified by {clusterExtId}.

## Example Usage

```hcl
data "nutanix_snmp_user_v2" "example" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id         = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: - (Required) Indicates the UUID of a cluster.
* `ext_id`: - (Required) SNMP user UUID.

## Attribute Reference

The following attributes are exported:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `username`: - SNMP username. For SNMP trap v3 version, SNMP username is required parameter.
* `auth_type`: - SNMP user authentication type. One of `MD5`, `SHA`, `SHA224`, `SHA256`, `SHA384`, `SHA512`.
* `auth_key`: - SNMP user authentication key.
* `priv_type`: - SNMP user encryption type. One of `DES`, `AES`, `AES192`, `AES256`.
* `priv_key`: - SNMP user encryption key.

### Links

The `links` block exports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix Get SNMP User](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.3#tag/Clusters/operation/getSnmpUserById).
