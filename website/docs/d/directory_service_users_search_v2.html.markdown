---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_service_users_search_v2"
sidebar_current: "docs-nutanix-datasource-directory-service-users-search-v2"
description: |-
  Searches for users or user groups in a directory service.
---

# nutanix_directory_service_users_search_v2

Searches for users or user groups in a directory service through its external identifier. This datasource calls the `POST /api/iam/v4.1/authn/directory-services/{extId}/$actions/search` API.

## Example Usage

### Search for a user by exact name

```hcl
# First, create or reference a directory service
resource "nutanix_directory_services_v2" "ad" {
  name           = "example_active_directory"
  url            = "ldap://10.xx.xx.xx:xxxx"
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = "example.com"
  service_account {
    username = "admin"
    password = "password"
  }
  white_listed_groups = ["Admins"]
  lifecycle {
    ignore_changes = [service_account.0.password]
  }
}

# Search for a specific user
data "nutanix_directory_service_users_search_v2" "find_user" {
  directory_service_ext_id = nutanix_directory_services_v2.ad.id
  query                    = "john.doe"
  is_wildcard_search       = false
}

output "search_results" {
  value = data.nutanix_directory_service_users_search_v2.find_user.search_results
}
```

### Wildcard search with returned attributes

```hcl
# Search with wildcard and request specific attributes
data "nutanix_directory_service_users_search_v2" "find_admins" {
  directory_service_ext_id = nutanix_directory_services_v2.ad.id
  query                    = "admin"
  is_wildcard_search       = true
  returned_attributes      = ["cn", "memberOf", "mail"]
}
```

### Search for a user group

```hcl
# Search for a group using specific search attributes
data "nutanix_directory_service_users_search_v2" "find_group" {
  directory_service_ext_id = nutanix_directory_services_v2.ad.id
  query                    = "DeveloperGroup"
  is_wildcard_search       = false
  searched_attributes      = ["cn"]
  returned_attributes      = ["cn", "member", "description"]
}
```

## Argument Reference

The following arguments are supported:

* `directory_service_ext_id`: -(Required) External identifier of the directory service to search.
* `query`: -(Required) Query string for directory service search.
* `is_wildcard_search`: -(Optional) Flag indicating whether the search should be a wildcard search or not. Defaults to `true`.
* `searched_attributes`: -(Optional) Attributes for search operation. By default, the search will be performed with a common name.
* `returned_attributes`: -(Optional) Attributes returned by the search operation.

## Attributes Reference

The following attributes are exported:

* `domain_name`: - Domain name for the directory service.
* `search_results`: - List of search result entities. Each entry contains the following:

### Search Results

* `entity_type`: - Type of entity, either user or group.
* `name`: - Name of the entity in canonical format.
* `attributes`: - List of attributes for the search entity.

#### Attributes

* `name`: - Name of the attribute.
* `values`: - List of values for the attribute.

See detailed information in [Nutanix Directory Service Search v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/DirectoryServices/operation/searchDirectoryService).
