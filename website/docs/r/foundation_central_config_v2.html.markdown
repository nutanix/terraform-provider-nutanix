---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_config_v2"
sidebar_current: "docs-nutanix-resource-foundation-central-config-v2"
description: |-
  Updates the configuration settings of the Foundation Central service.
---

# nutanix_foundation_central_config_v2

Updates the configuration settings of the Foundation Central service.

Foundation Central exposes a single, cluster-wide configuration object (there is no
entity identifier). This resource therefore behaves as a singleton: applying it
updates the current configuration, and destroying it is a no-op (the configuration
object itself cannot be removed).

## Example

```hcl
resource "nutanix_foundation_central_config_v2" "example" {
  ahv_installation_timeout_minutes = 60
  aos_download_timeout_minutes     = 60
}
```

## Argument Reference

The following arguments are supported:

* `ahv_installation_timeout_minutes`: (Optional) Timeout in minutes for AHV installation.
* `aos_download_timeout_minutes`: (Optional) Timeout in minutes for AOS download.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `commit_id`: Foundation Central Commit ID.
* `tls_certificate_fingerprint`: SHA-256 fingerprint of the TLS certificate used by Foundation Central service.
* `version`: Foundation Central Version.

See detailed information in [Nutanix Update Foundation Central Config V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
