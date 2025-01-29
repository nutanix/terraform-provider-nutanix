---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_vm_v2"
sidebar_current: "docs-nutanix-resource-volume-group-vm-attachments-v2"
description: |-
  This operation submits a request to Attaches VM to a Volume Group identified by {extId}.
---

# nutanix_volume_group_vm_v2

Provides a resource to Create a new Volume Group.

## Example Usage

```hcl

resource "nutanix_volume_group_vm_v2" "vg_vm_example"{
  volume_group_ext_id = "<vg uuid>"
  vm_ext_id           = var.vg_vm_ext_id
}
```

## Argument Reference
The following arguments are supported:


* `volume_group_ext_id`: -(Required) The external identifier of the volume group.
* `vm_ext_id`: -(Required) A globally unique identifier of an instance that is suitable for external consumption. 
* `index`: -(Optional) The index on the SCSI bus to attach the VM to the Volume Group. 


See detailed information in [Nutanix Volumes V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0).
