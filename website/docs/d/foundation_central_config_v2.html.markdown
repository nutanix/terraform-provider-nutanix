---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_config_v2"
sidebar_current: "docs-nutanix-datasource-foundation-central-config-v2"
description: |-
  Returns the current configuration settings of the Foundation Central service.
---

# nutanix_foundation_central_config_v2

Returns the current configuration settings of the Foundation Central service.

~> **NOTE:** The Foundation Central configuration APIs are served by the Foundation Central (FCVM) service, not Prism Central. Configure the provider with `foundation_endpoint` and `foundation_port` so the data source talks to the correct host.

## Example

```hcl
data "nutanix_foundation_central_config_v2" "current" {}
```

## Argument Reference

This data source is a singleton and takes no arguments.

## Attribute Reference

The following attributes are exported:

* `ahv_installation_timeout_minutes`: Timeout in minutes for AHV installation.
* `aos_download_timeout_minutes`: Timeout in minutes for AOS download.
* `commit_id`: Foundation Central Commit ID.
* `tls_certificate_fingerprint`: SHA-256 fingerprint of the TLS certificate used by Foundation Central service.
* `version`: Foundation Central Version.

See detailed information in [Nutanix Get Foundation Central Config V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
