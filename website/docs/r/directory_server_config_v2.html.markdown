---
layout: "nutanix"
page_title: "NUTANIX: nutanix_directory_server_config_v2"
sidebar_current: "docs-nutanix-resource-directory-server-config-v2"
description: |-
  Configures various aspects of identity categorization.
---

# nutanix_directory_server_config_v2

Configures various aspects of identity categorization for Flow Network Security ID-based security (ID Firewall). The Directory Server Config controls which entities are eligible for dynamic categorization, the domain controllers used to scrape Windows logon events, and the default-category behavior.

~> **Note:** Only one Directory Server Config is allowed per Prism Central instance.

## Example Usage

### CONTAINS match type

```hcl
resource "nutanix_directory_server_config_v2" "example" {
  directory_service_reference           = "<Directory Service ext_id>"
  is_default_category_enabled           = true
  # is_default_category_enabled = true - Allowed only for CONTAINS match type

  should_keep_default_category_on_login = false
  # should_keep_default_category_on_login = true
  # Allowed only for CONTAINS match type
  # Allowed only when is_default_category_enabled = true

  # Determine which entities are eligible for dynamic categorization.
  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "CONTAINS"
    criteria     = "DeveloperVM"
    # Criteria should be the name of the VMs on which the category should be applied.
    # Criteria is allowed only for CONTAINS match type
    # Criteria is not allowed for ALL match type

    # If you want to apply the category to all the VMs, match_type should be ALL.
  }

  # Domain controllers used to scrape Windows logon events.
  domain_controllers {
    fqdn {
      value = "dc01.example.com"
    }
  }
}
```

### ALL match type

```hcl
resource "nutanix_directory_server_config_v2" "all_vms" {
  directory_service_reference           = var.directory_service_reference
  is_default_category_enabled           = false
  should_keep_default_category_on_login = false

  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "ALL"
    # criteria must NOT be set when match_type is ALL
  }

  domain_controllers {
    ipv4 {
      value = "10.0.0.10"
    }
  }
}
```

## Argument Reference

* `directory_service_reference` - (Optional) The ExtID of the directory service that will be used for mapping.
* `domain_controllers` - (Optional) List of domain controllers to be used for event scraping.
* `is_default_category_enabled` - (Optional) Enablement status of the default category. Can only be set to `true` when `match_type` is `CONTAINS`. Must be `false` when `match_type` is `ALL`.
* `matching_criterias` - (Optional) The matching criteria list used to determine whether an entity will be categorized by identity categorization. If match type is `ALL`, all the entities will be categorized.
* `should_keep_default_category_on_login` - (Optional) Retain default category on user login. Can only be `true` when `is_default_category_enabled` is also `true` and `match_type` is `CONTAINS`. Must be `false` when `match_type` is `ALL`.
* `project_ext_id` - (Optional) The external identifier of the project associated with the configuration.

### domain_controllers

Each `domain_controllers` block is an IP address or FQDN and may contain exactly one of:

* `ipv4` - (Optional) An IPv4 address with `value` (Required) and `prefix_length` (Optional).
* `ipv6` - (Optional) An IPv6 address with `value` (Required) and `prefix_length` (Optional).
* `fqdn` - (Optional) A fully qualified domain name with `value` (Optional).

### matching_criterias

* `criteria` - (Optional) The criteria to use for matching entities to be categorized. Only allowed when `match_type` is `CONTAINS`. Must not be set when `match_type` is `ALL`.
* `match_entity` - (Optional) The entity type to match on. Valid values: `VM`.
* `match_field` - (Optional) The field to match on. Today only `NAME` is supported, which matches on an entity's name.
* `match_type` - (Optional) The type of match. Valid values: `CONTAINS`, `ALL`. `CONTAINS` performs a substring match on the given entity and field for the criteria value whereas `ALL` allows all strings to match on the given entity.

## Attributes Reference

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `links` - A HATEOAS style link for the response.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

## Import

Directory Server Config can be imported using its `ext_id`:

```bash
terraform import nutanix_directory_server_config_v2.example <directoryServerConfigUUID>
```

See detailed information in [Nutanix Directory Server Configs V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.3).
