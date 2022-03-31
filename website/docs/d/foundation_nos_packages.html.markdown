---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_nos_packages"
sidebar_current: "docs-nutanix-datasource-foundation-nos-packages"
description: |-
 Describes a list of nos packages present in foundation vm
---

# nutanix_clusters

Describes a list of nos (aos) packages present in foundation vm

## Example Usage

```hcl
data "nutanix_foundation_nos_packages" "nos_packages" {}
```

## Argument Reference

No arguments are supported

## Attribute Reference

The following attributes are exported:

* `entities`: List of nos packages file names present in foundation vm

## Note
* This data source only lists .tar file names.

See detailed information in [Nutanix Foundation Nos Packages](https://www.nutanix.dev/api_references/foundation/#/b3A6MjIyMjMzODg-get-list-of-aos-packages-available-in-foundation).
