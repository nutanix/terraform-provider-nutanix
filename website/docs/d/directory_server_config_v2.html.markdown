---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_server_config_v2"
sidebar_current: "docs-nutanix-datasource-directory-server-config-v2"
description: |-
  Gets the Directory Server configuration with the provided ExtID.
---

# nutanix_directory_server_config_v2

Gets the Directory Server configuration with the provided ExtID.

## Example Usage

```hcl
# Read a single Directory Server configuration by ext_id.
data "nutanix_directory_server_config_v2" "get-config" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id` - (Required) A globally unique identifier of the Directory Server Config.

## Attributes Reference

* `directory_service_reference` - The ExtID of the directory service that will be used for mapping.
* `domain_controllers` - List of domain controllers to be used for event scraping. Each entry contains one of `ipv4`, `ipv6` or `fqdn`.
* `is_default_category_enabled` - Enablement status of the default category.
* `matching_criterias` - The matching criteria used to determine whether an entity will be categorized by identity categorization.
  * `criteria` - The criteria to use for matching entities to be categorized.
  * `match_entity` - The entity type to match on.
  * `match_field` - The field to match on.
  * `match_type` - The type of match.
* `should_keep_default_category_on_login` - Retain default category on user login.
* `project_ext_id` - The external identifier of the project associated with the configuration.
* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `links` - A HATEOAS style link for the response.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Directory Server Configs V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.3).
