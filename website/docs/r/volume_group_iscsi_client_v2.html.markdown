---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_iscsi_client_v2"
sidebar_current: "docs-nutanix-resource-volume-group-iscsi-client-v2"
description: |-
  This operation submits a request to Attaches iSCSI initiator to a Volume Group identified by {extId}.
---

# nutanix_volume_group_iscsi_client_v2
Attaches iSCSI initiator to a Volume Group identified by {extId}.

## Example Usage

```hcl

#list iscsi clients
data "nutanix_volume_iscsi_clients_v2" "list-iscsi-clients"{}

# attach iscsi client to the volume group
resource "nutanix_volume_group_iscsi_clients_v2" "vg_iscsi_example"{
  vg_ext_id            = "1cdb5b48-fb2c-41b6-b751-b504117ee3e2"
  ext_id               = data.nutanix_volume_iscsi_clients_v2.list-iscsi-clients.iscsi_clients.0.ext_id
  iscsi_initiator_name = data.nutanix_volume_iscsi_clients_v2.list-iscsi-clients.iscsi_clients.0.iscsi_initiator_name
}
```

## Argument Reference
The following arguments are supported:


* `vg_ext_id`: -(Required) The external identifier of the volume group.
* `ext_id`: -(Required) A globally unique identifier of an instance that is suitable for external consumption.
* `iscsi_initiator_name`: -iSCSI initiator name. During the attach operation, exactly one of iscsiInitiatorName and iscsiInitiatorNetworkId must be specified. This field is immutable.
* `iscsi_initiator_network_id`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
* `client_secret`: -(Optional) iSCSI initiator client secret in case of CHAP authentication. This field should not be provided in case the authentication type is not set to CHAP.
* `enabled_authentications`: -(Optional) (Optional) The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided. Valid values are CHAP, NONE
* `num_virtual_targets`: -(Optional) Number of virtual targets generated for the iSCSI target. This field is immutable.
* `attachment_site`: -(Optional) The site where the Volume Group attach operation should be processed. This is an optional field. This field may only be set if Metro DR has been configured for this Volume Group. Valid values are SECONDARY, PRIMARY.

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


See detailed information in [Nutanix Attach an iSCSI Client to Volume Group V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/VolumeGroups/operation/attachIscsiClient).
