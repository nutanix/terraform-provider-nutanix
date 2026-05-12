---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_vm_attachments_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-vm-attachments-v2"
description: |-
  Query the list of VM attachments for a Volume Group.
---

# nutanix_volume_group_vm_attachments_v2

Query the list of VM attachments for a Volume Group identified by {extId}. Deprecated: This API has been deprecated.

## Example Usage

```hcl
data "nutanix_volume_group_vm_attachments_v2" "example" {
  volume_group_ext_id = "d09aeec9-5bb7-4bfd-9717-a051178f6e7c"
}
```

## Argument Reference

The following arguments are supported:

* `volume_group_ext_id`: -(Required) The external identifier of a Volume Group.

## Attributes Reference

The following attributes are exported:

* `vm_attachments`: - List of VM attachments for the Volume Group.

### VM Attachments

Each element in `vm_attachments` has the following fields:

* `ext_id`: - The external identifier of the VM.
* `index`: - The index on the SCSI bus to attach the VM to the Volume Group. This is an optional field.
