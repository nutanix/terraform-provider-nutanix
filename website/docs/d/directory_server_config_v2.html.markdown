---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_server_config_v2"
sidebar_current: "docs-nutanix-datasource-directory-server-config-v2"
description: |-
  Gets the Directory Server configuration by ExtID.
---

# nutanix_directory_server_config_v2

Gets the Directory Server configuration by ExtID.

## Example Usage

```hcl
data "nutanix_directory_server_config_v2" "example" {
  ext_id = "<ext-id>"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) The external identifier of the directory server config.

## Attributes Reference

The following attributes are exported:

* `tenant_id` - A globally unique identifier that represents the tenant.
* `links` - A HATEOAS style link for the response.
* `directory_service_reference` - The ExtID of the directory service used for mapping.
* `domain_controllers` - List of domain controllers used for event scraping.
  * `fqdn` - Fully qualified domain name.
    * `value` - The FQDN value.
  * `ipv4` - IPv4 address.
    * `value` - The IPv4 address.
    * `prefix_length` - The prefix length.
  * `ipv6` - IPv6 address.
    * `value` - The IPv6 address.
    * `prefix_length` - The prefix length.
* `is_default_category_enabled` - Enablement status of the default category.
* `matching_criterias` - The matching criteria for identity categorization.
  * `criteria` - The criteria string for matching.
  * `match_entity` - The entity to match.
  * `match_field` - The field to match on.
  * `match_type` - The type of match.
* `should_keep_default_category_on_login` - Retain default category on user login.
* `project_ext_id` - The external identifier of the associated project.
