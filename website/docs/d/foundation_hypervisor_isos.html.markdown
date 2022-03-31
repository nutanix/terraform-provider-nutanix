---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_hypervisor_isos"
sidebar_current: "docs-nutanix-datasource-foundation-hypervisor-isos"
description: |-
 Describes a list of hypervisor isos packages present in foundation vm
---

# nutanix_foundation_hypervisor_isos

Describes a list of hypervisor isos image file details present in foundation vm

## Example Usage

```hcl
data "nutanix_foundation_hypervisor_isos" "hypervisor_isos" {}
```

## Argument Reference

No arguments are supported

## Attribute Reference

The following attributes are exported:

* `hyperv`: List of hyperv isos and their details present in foundation vm
* `kvm`: List of kvm isos and their details present in foundation vm
* `linux`: List of linux isos and their details present in foundation vm
* `esx`: List of esx isos and theirdetails present in foundation vm
* `xen`: List of esx isos and theirdetails present in foundation vm

### hyperv
* `filename`: Name of installer.
* `supported`: Whether front-end should treat hyp as supported.

### kvm
* `filename`: Name of installer.
* `supported`: Whether front-end should treat hyp as supported.

### linux
* `filename`: Name of installer.
* `supported`: Whether front-end should treat hyp as supported.

### esx
* `filename`: Name of installer.
* `supported`: Whether front-end should treat hyp as supported.

### xen
* `filename`: Name of installer.
* `supported`: Whether front-end should treat hyp as supported.

## Note
* This data source only lists .iso files details.

See detailed information in [Nutanix Foundation Hypervisor Isos](https://www.nutanix.dev/api_references/foundation/#/b3A6MjIyMjM0MDE-list-hypervisor-images-available-in-foundation).
