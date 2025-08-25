---
layout: "nutanix"
page_title: "NUTANIX: nutanix_system_user_passwords_v2 "
sidebar_current: "docs-nutanix-datasource-system-user-passwords-v2"
description: |-
  Lists password status of system user accounts on supported products.


---

# nutanix_system_user_passwords_v2

Lists password status of system user accounts on supported products.



## Example Usage

```hcl

# List Password Status Of All System Users
data "nutanix_system_user_passwords_v2" "passwords" {
}

# List Password Status Of All System Users With Limit
data "nutanix_system_user_passwords_v2" "limited_passwords" {
  limit  = 10
}

# List Password Status Of All System Users With Filter
data "nutanix_system_user_passwords_v2" "filtered_passwords" {
  filter = "systemType eq Clustermgmt.Config.SystemType'PC'"
}

# List Password Status Of Admin PC User
data "nutanix_system_user_passwords_v2" "admin_pc_passwords" {
  filter = "username eq 'admin' and systemType eq Clustermgmt.Config.SystemType'PC'"
}


```

## Argument Reference

The following arguments are supported:

e following attributes are exported:

- `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. The filter can be applied to the following fields:

  - clusterExtId : `filter="clusterExtId eq '8a72db6b-83f3-47b2-a65c-f5de2e50efb9'"`
  - hasHspInUse : `filter="hasHspInUse eq false"`
  - hostIp/value: `filter="hostIp/value eq '240.29.254.180'"`
  - status : `filter="status eq Clustermgmt.Config.PasswordStatus'DEFAULT'"`
  - systemType : `filter="systemType eq Clustermgmt.Config.SystemType'PC'"`
  - username : `filter="username eq 'admin'"`


- `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
  - username : `orderby="username desc"`

- `expand`:A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are \$filter, \$select and \$orderby. The following expansion keys are supported:

  - cluster : `expand="cluster"`

- `select`: URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned.
The following selection keys are supported:
  - hasHspInUse : `select="hasHspInUse"`
  - hostIp : `select="hostIp"`
  - status : `select="status"`
  - systemType : `select="systemType"`
  - username : `select="username"`



## Attributes Reference

The following attributes are exported:

- `passwords`: - List password status of system users

### Passwords
The `passwords` attribute supports the following:


- `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `username`: - Username.
- `host_ip`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
- `cluster_ext_id`: - UUID of the cluster to which the host NIC belongs.
- `last_update_time`: - Timestamp of last password change.
- `expiry_time`: - Expiry of a new password.
- `status`: - Contains possible values of password status.
  - `MULTIPLE_ISSUES`: - Some user accounts have default password or no password set.
  - `SECURE`: - Secure password is set.
  - `NOPASSWD`: - No password is set.
  - `DEFAULT`: - Default password is set.
- `system_type`: - Contains supported variants of the system products.
  - `IPMI`: - The product is of IPMI type.
  - `PC`: - The product is of Prism Central type.
  - `AOS`: - The product is of AOS type.
  - `AHV`: - The product is of AHV type.
- `has_hsp_in_use`: - Indicates whether the high-strength password is in use or not.

### Links

The `links` attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.


### Host IP
The `host_ip` attribute supports the following:

- `value`: - The IPv4 address of the host.
- `prefix_length`: - The prefix length of the network to which this host IPv4/IPv6 address belongs.

See detailed information in [Nutanix List password status of system users V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.1#tag/PasswordManager/operation/listSystemUserPasswords).
