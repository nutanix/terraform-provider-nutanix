---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_iscsi_client_v2"
sidebar_current: "docs-nutanix-datasource-volume-iscsi-client-v2"
description: |-
  Describes iSCSI client details identified by {extId}.
---

# nutanix_volume_iscsi_clients_v2

Fetches the iSCSI client details identified by {extId}.



## Example Usage

```hcl
data "nutanix_volume_iscsi_client_v2" "example"{
  ext_id = "be0e4630-23da-4b9c-a76b-f24fd64b46b6"
 }
```

##  Argument Reference
The following arguments are supported:


* `ext_id`: -(Required) A query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.

## Attributes Reference
The following attributes are exported:


* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `iscsi_initiator_name`: -iSCSI initiator name. During the attach operation, exactly one of iscsiInitiatorName and iscsiInitiatorNetworkId must be specified. This field is immutable.
* `iscsi_initiator_network_id`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
* `client_secret`: -(Optional) iSCSI initiator client secret in case of CHAP authentication. This field should not be provided in case the authentication type is not set to CHAP.
* `enabled_authentications`: -(Optional) (Optional) The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided. Valid values are CHAP, NONE
* `num_virtual_targets`: -(Optional) Number of virtual targets generated for the iSCSI target. This field is immutable.
* `attachment_site`: -(Optional) The site where the Volume Group attach operation should be processed. This is an optional field. This field may only be set if Metro DR has been configured for this Volume Group. Valid values are SECONDARY, PRIMARY.


#### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

#### iscsi initiator network id

The iscsi_initiator_network_id attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.
* `fqdn`: - A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

##### IPV4

The ipv4 attribute supports the following:

* `value`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `prefix_length`: - The prefix length of the network to which this host IPv4 address belongs.

##### IPV6

The ipv6 attribute supports the following:

* `value`: - An unique address that identifies a device on the internet or a local network in IPv6 format.
* `prefix_length`: - The prefix length of the network to which this host IPv6 address belongs.

##### FQDN

The fqdn attribute supports the following:

* `value`: - The fully qualified domain name.


#### Attached Targets

The attached_targets attribute supports the following:

* `num_virtual_targets`: - Number of virtual targets generated for the iSCSI target. This field is immutable.
* `iscsi_target_name`: - Name of the iSCSI target that the iSCSI client is connected to. This is a read-only field.



See detailed information in [Nutanix Get iSCSI Client V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/IscsiClients/operation/getIscsiClientById).
