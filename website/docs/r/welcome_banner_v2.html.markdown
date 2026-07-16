---
layout: "nutanix"
page_title: "NUTANIX: nutanix_welcome_banner_v2"
sidebar_current: "docs-nutanix-resource-welcome-banner-v2"
description: |-
  Updates the welcome banner.
---

# nutanix_welcome_banner_v2

Provides Nutanix resource to update the welcome banner.

The welcome banner is a cluster-wide singleton configuration. There is exactly one welcome banner per Prism Central, so this resource does not use an external identifier. Creating or destroying the resource updates the same underlying banner configuration.

## Example Usage

```hcl
resource "nutanix_welcome_banner_v2" "example" {
  content    = "Welcome to the Nutanix cluster. Authorized access only."
  is_enabled = true
}
```

## Argument Reference

The following arguments are supported:

- `content`: - (Optional) Content of the welcome banner.
- `is_enabled`: - (Optional) Flag to denote whether the welcome banner is enabled or not.

## Attributes Reference

The following attributes are exported:

- `content`: - Content of the welcome banner.
- `is_enabled`: - Flag to denote whether the welcome banner is enabled or not.
- `created_time`: - Creation time of the welcome banner.
- `last_updated_time`: - Last updated time of the welcome banner.

## Import

This helps to manage the existing welcome banner configuration which is not created through terraform. The welcome banner is a singleton, so it can be imported using a fixed identifier `welcome_banner`. eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_welcome_banner_v2" "import_banner" {}

// execute the below command:
terraform import nutanix_welcome_banner_v2.import_banner welcome_banner
```

See detailed information in [Nutanix Update Welcome Banner V4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.1).
