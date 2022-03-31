---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_image"
sidebar_current: "docs-nutanix-resource-foundation-image"
description: |-
  Uploads hypervisor or AOS image to foundation.
---

# nutanix_foundation_image

Uploads hypervisor or AOS image to foundation.

## Example Usage

```hcl
resource "nutanix_foundation_image" "nos-image" {
  source = "../../../files/nutanix_installer_x86_64.tar"
  filename = "nos_image.tar"
  installer_type = "nos"
}
resource "nutanix_foundation_image" "hypervisor-image" {
  source = "../../../files/VMware-Installer.x86_64.iso"
  filename = "esx_image.iso"
  installer_type = "esx"
}
```

## Argument Reference

The following arguments are supported:

* `installer_type`: - (Required) One of "kvm", "esx", "hyperv", "xen", or "nos".
* `filename`: - (Required) Name of installer file to be kept in foundation vm.
* `source`: - (Required) Complete path to the file in machine where the .tf  files runs.

## Attributes Reference

The following attributes are exported:

* `md5sum` : - md5sum of the ISO.
* `name` :- file location in foundation vm
* `in_whitelist` :- If hypervisor ISO is in whitelist.

## lifecycle

* `Update` : - Resource will trigger new resource create call for any kind of update in resource config and delete existing image from foundation vm.

See detailed information in [Nutanix Foundation Image](https://www.nutanix.dev/api_references/foundation/#/b3A6MjIyMjM0MDQ-upload-hypervisor-or-aos-image-to-foundation).
