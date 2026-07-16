---
layout: "nutanix"
page_title: "NUTANIX: nutanix_welcome_banner_v2"
sidebar_current: "docs-nutanix-datasource-welcome-banner-v2"
description: |-
  Fetches the configured welcome banner.
---

# nutanix_welcome_banner_v2

Fetches the configured welcome banner.

## Example Usage

```hcl
data "nutanix_welcome_banner_v2" "welcome-banner" {}
```

## Argument Reference

This is a singleton data source and does not take any arguments.

## Attributes Reference

The following attributes are exported:

- `content`: - Content of the welcome banner.
- `created_time`: - Creation time of the welcome banner.
- `is_enabled`: - Flag to denote whether the welcome banner is enabled or not.
- `last_updated_time`: - Last updated time of the welcome banner.

See detailed information in [Nutanix Get Welcome Banner v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.1#tag/WelcomeBanner/operation/getWelcomeBanner).
