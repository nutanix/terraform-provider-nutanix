---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_iscsi_clients_v2"
sidebar_current: "docs-nutanix-datasource-volume-iscsi-clients-v2"
description: |-
  Describes aList all the iSCSI clients.
---

# nutanix_volume_iscsi_clients_v2

Fetches the list of iSCSI clients.


## Example Usage

```hcl
# List all the iSCSI clients.
data "nutanix_volume_iscsi_clients_v2" "list-iscsi-clients" {}

# list iSCSI clients with a filter.
data "nutanix_volume_iscsi_clients_v2" "list-iscsi-clients-filter" {
  filter = "clusterReference eq '00057b8b-0b7b-4b7b-0000-000000000000'"
}

# list iSCSI clients with a limit and pagination.
data "nutanix_volume_iscsi_clients_v2" "list-iscsi-clients-limit" {
  page  = 2
  limit = 1
}
```

##  Argument Reference

The following arguments are supported:

* `page`: - A query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` : A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - clusterReference
  - extId
* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - clusterReference
  - extId
* `expand` : A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are \$filter, \$select and \$orderby. The following expansion keys are supported. The expand can be applied to the following fields:
  - cluster
* `select` : A query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., \*), then all properties on the matching resource will be returned. The select can be applied to the following fields:
  - clusterReference
  - extId

## Attributes Reference
The following attributes are exported:

* `iscsi_clients`: - List all the iSCSI clients.

## Iscsi Clients
The `iscsi_clients` contains list of iSCSI clients. Each iSCSI client contains the following attributes:

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



See detailed information in [Nutanix List iSCSI Clients V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/IscsiClients/operation/listIscsiClients).
