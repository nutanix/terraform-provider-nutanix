---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_server_configs_v2"
sidebar_current: "docs-nutanix-datasource-directory-server-configs-v2"
description: |-
  Gets the list of Directory Servers.
---

# nutanix_directory_server_configs_v2

Gets the list of Directory Servers.

## Example Usage

```hcl
# List all Directory Server configurations.
data "nutanix_directory_server_configs_v2" "configs-list" {}
```

## Argument Reference

* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.

## Attributes Reference

* `directory_server_configs` - List of Directory Server configurations. Each element exposes the same attributes as the `nutanix_directory_server_config_v2` datasource:
  * `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
  * `directory_service_reference` - The ExtID of the directory service that will be used for mapping.
  * `domain_controllers` - List of domain controllers to be used for event scraping.
  * `is_default_category_enabled` - Enablement status of the default category.
  * `matching_criterias` - The matching criteria used to determine whether an entity will be categorized.
  * `should_keep_default_category_on_login` - Retain default category on user login.
  * `project_ext_id` - The external identifier of the project associated with the configuration.
  * `links` - A HATEOAS style link for the response.
  * `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Directory Server Configs V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.3).
