---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_config_v2"
sidebar_current: "docs-nutanix-datasource-foundation-central-config-v2"
description: |-
  Returns the current configuration settings of the Foundation Central service.
---

# nutanix_foundation_central_config_v2

Returns the current configuration settings of the Foundation Central service.

Foundation Central exposes a single, cluster-wide configuration object, so this data
source takes no id input and returns that one object.

## Example

```hcl
data "nutanix_foundation_central_config_v2" "current" {}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

The following attributes are exported:

* `ahv_installation_timeout_minutes`: Timeout in minutes for AHV installation.
* `aos_download_timeout_minutes`: Timeout in minutes for AOS download.
* `commit_id`: Foundation Central Commit ID.
* `tls_certificate_fingerprint`: SHA-256 fingerprint of the TLS certificate used by Foundation Central service.
* `version`: Foundation Central Version.

See detailed information in [Nutanix Get Foundation Central Config V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
