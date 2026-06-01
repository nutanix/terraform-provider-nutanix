---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_server_configs_v2"
sidebar_current: "docs-nutanix-datasource-directory-server-configs-v2"
description: |-
  Gets the list of Directory Server configurations.
---

# nutanix_directory_server_configs_v2

Gets the list of Directory Server configurations.

## Example Usage

```hcl
data "nutanix_directory_server_configs_v2" "all" {}
```

## Argument Reference

The following arguments are supported:

* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attributes Reference

The following attributes are exported:

* `directory_server_configs` - List of directory server configurations. Each element contains:
  * `ext_id` - The external identifier.
  * `tenant_id` - The tenant identifier.
  * `links` - A HATEOAS style link for the response.
  * `directory_service_reference` - The ExtID of the directory service.
  * `domain_controllers` - List of domain controllers.
  * `is_default_category_enabled` - Enablement status of the default category.
  * `matching_criterias` - The matching criteria for identity categorization.
  * `should_keep_default_category_on_login` - Retain default category on login.
  * `project_ext_id` - The external identifier of the associated project.
