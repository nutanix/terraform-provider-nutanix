---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_services_v2"
sidebar_current: "docs-nutanix-resource-directory-services-v2"
description: |-
  This operation submits a request to Create a Directory Service.
---

# nutanix_directory_services_v2

Provides a resource to Create a Directory Service.

## Example Usage

```hcl
# Add Directory Service .
resource "nutanix_directory_services_v2" "active-directory" {
  name           = "example_active_directory"
  url            = "ldap://10.xx.xx.xx:xxxx"
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = "nutanix.com"
  service_account {
    username = "username"
    password = "password"
  }
  white_listed_groups = ["example"]
  lifecycle {
    ignore_changes = [
      service_account.0.password,
    ]
  }
}
```

## Argument Reference
The following arguments are supported:


* `ext_id`: -(Optional) A globally unique identifier of an instance that is suitable for external consumption.
* `name`: -(Required) Name for the Directory Service.
* `url`: -(Required) URL for the Directory Service.
* `secondary_urls`: -(Optional) Secondary URL for the Directory Service.
* `domain_name`: -(Required) Domain name for the Directory Service.
* `directory_type`: -(Required) Type of Directory Service, Supported values are: "ACTIVE_DIRECTORY" (Directory Service type is Active Directory.) and "OPEN_LDAP" (Directory Service type is Open LDAP.)
* `service_account`: -(Required) Information of Service account to connect to the Directory Service.
* `open_ldap_configuration`: -(Optional) Configuration for OpenLDAP Directory Service.
* `group_search_type`: -(Optional) Group membership search type for the Directory Service. Supported values are: "NON_RECURSIVE" (Doesn't search recursively within groups.) and "RECURSIVE" (Searches recursively within groups.)
* `white_listed_groups`: -(Optional) List of allowed User Groups for the Directory Service.


### Service Account

The service_account attribute supports the following:

* `username`: -(Required) Username to connect to the Directory Service.
* `password`: -(Required) Password to connect to the Directory Service.


### Open Ldap Configuration

The open_ldap_configuration attribute supports the following:

* `user_configuration`: -(Required) this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.
* `user_group_configuration`: -(Required) this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.

#### User Configuration

The user_configuration attribute supports the following:

* `user_object_class`: -(Required) Object class in the OpenLDAP system that corresponds to Users.
* `user_search_base`: -(Required) Base DN for User search.
* `username_attribute`: -(Required) Unique Identifier for each User which can be used in Authentication.

#### User Group Configuration

The user_group_configuration attribute supports the following:

* `group_object_class`: -(Required) Object class in the OpenLDAP system that corresponds to groups.
* `group_search_base`: -(Required) Base DN for group search.
* `group_member_attribute`: -(Required) Attribute in a group that associates Users to the group.
* `group_member_attribute_value`: -(Required) User attribute value that will be used in group entity to associate User to the group.


## Attributes Reference
The following attributes are exported:


* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`: - Name for the Directory Service.
* `url`: - URL for the Directory Service.
* `secondary_urls`: - Secondary URL for the Directory Service.
* `domain_name`: - Domain name for the Directory Service.
* `directory_type`: - Type of Directory Service, Supported values are: "ACTIVE_DIRECTORY" (Directory Service type is Active Directory.) and "OPEN_LDAP" (Directory Service type is Open LDAP.)
* `service_account`: - Information of Service account to connect to the Directory Service.
* `open_ldap_configuration`: - Configuration for OpenLDAP Directory Service.
* `group_search_type`: - Group membership search type for the Directory Service. Supported values are: "NON_RECURSIVE" (Doesn't search recursively within groups.) and "RECURSIVE" (Searches recursively within groups.)
* `white_listed_groups`: - List of allowed User Groups for the Directory Service.
* `created_time`: - Creation time of the Directory Service.
* `last_updated_time`: - Last updated time of the Directory Service.
* `created_by`: - User or Service who created the Directory Service.

### Service Account

The service_account attribute supports the following:

* `username`: - Username to connect to the Directory Service.
* `password`: - Password to connect to the Directory Service.


### Open Ldap Configuration

The open_ldap_configuration attribute supports the following:

* `user_configuration`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.
* `user_group_configuration`: - this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.

#### User Configuration

The user_configuration attribute supports the following:

* `user_object_class`: - Object class in the OpenLDAP system that corresponds to Users.
* `user_search_base`: - Base DN for User search.
* `username_attribute`: - Unique Identifier for each User which can be used in Authentication.

#### User Group Configuration

The user_group_configuration attribute supports the following:

* `group_object_class`: - Object class in the OpenLDAP system that corresponds to groups.
* `group_search_base`: - Base DN for group search.
* `group_member_attribute`: - Attribute in a group that associates Users to the group.
* `group_member_attribute_value`: - User attribute value that will be used in group entity to associate User to the group.


See detailed information in [Nutanix Directory Services v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/DirectoryServices/operation/createDirectoryService).
