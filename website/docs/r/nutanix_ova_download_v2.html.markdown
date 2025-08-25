---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ova_download_v2 "
sidebar_current: "docs-nutanix-resource-ova-download-v2"
description: |-
  Downloads an OVA based on the given external identifier. This is a stream download of the OVA file.
---

# nutanix_ova_download_v2

Downloads an OVA based on the given external identifier. This is a stream download of the OVA file..

## Example Usage

```hcl
// download ova file
data "nutanix_ova_download_v2" "example"{
  ova_ext_id = "8cf09a55-6ee3-45dc-bd67-239244dbecf7"
}
```

## Argument Reference

The following arguments are supported:

- `ova_ext_id`: -(Required) The external identifier for an OVA.

## Attributes Reference
The following attributes are exported:

- `ova-file-path`: The file path where the OVA is downloaded.


See detailed information in [Nutanix Download an Ova V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.1#tag/Ovas/operation/getFileByOvaId).
