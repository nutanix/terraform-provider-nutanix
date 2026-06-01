---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_server_config_v2"
sidebar_current: "docs-nutanix-resource-directory-server-config-v2"
description: |-
  Configures various aspects of identity categorization for a Directory Server.
---

# nutanix_directory_server_config_v2

Configures various aspects of identity categorization for a Directory Server.

## Example Usage

```hcl
resource "nutanix_directory_server_config_v2" "example" {
  directory_service_reference = "<directory-service-ext-id>"

  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "ALL"
  }

  is_default_category_enabled           = true
  should_keep_default_category_on_login = false
}
```

## Argument Reference

The following arguments are supported:

* `directory_service_reference` - (Required, ForceNew) The ExtID of the directory service that will be used for mapping.
* `domain_controllers` - (Optional) List of domain controllers to be used for event scraping. Each element supports the following:
  * `fqdn` - (Optional) Fully qualified domain name.
    * `value` - (Optional) The FQDN value.
  * `ipv4` - (Optional) IPv4 address.
    * `value` - (Optional) The IPv4 address.
    * `prefix_length` - (Optional) The prefix length.
  * `ipv6` - (Optional) IPv6 address.
    * `value` - (Optional) The IPv6 address.
    * `prefix_length` - (Optional) The prefix length.
* `is_default_category_enabled` - (Optional) Enablement status of the default category.
* `matching_criterias` - (Optional) The matching criteria used to determine whether an entity will be categorized. Each element supports the following:
  * `criteria` - (Optional) The criteria to use for matching entities.
  * `match_entity` - (Optional) The entity to match. Valid values: `VM`.
  * `match_field` - (Optional) The field to match on. Valid values: `NAME`.
  * `match_type` - (Optional) The type of match. Valid values: `CONTAINS`, `ALL`.
* `should_keep_default_category_on_login` - (Optional) Retain default category on user login.

## Attributes Reference

The following attributes are exported:

* `ext_id` - The external identifier of the directory server config.
* `tenant_id` - A globally unique identifier that represents the tenant.
* `links` - A HATEOAS style link for the response.
* `project_ext_id` - The external identifier of the associated project.

## Import

Directory Server Config can be imported using the `ext_id`:

```shell
terraform import nutanix_directory_server_config_v2.example <ext_id>
```
