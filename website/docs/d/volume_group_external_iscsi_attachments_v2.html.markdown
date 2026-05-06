---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_external_iscsi_attachments_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-external-iscsi-attachments-v2"
description: |-
  Lists external iSCSI attachments for a Volume Group.
---

# nutanix_volume_group_external_iscsi_attachments_v2

Query the list of external iSCSI attachments for a Volume Group identified by {extId}. Deprecated: This API has been deprecated.

## Example Usage

```hcl
data "nutanix_volume_group_external_iscsi_attachments_v2" "example" {
  volume_group_ext_id = "d09aeec9-5bb7-4bfd-9717-a051178f6e7c"
}
```

## Argument Reference

The following arguments are supported:

* `volume_group_ext_id`: -(Required) The external identifier of a Volume Group.

## Attributes Reference

The following attributes are exported:

* `external_iscsi_attachments`: - List of external iSCSI attachments.

### External iSCSI Attachments

* `ext_id`: - The external identifier of an iSCSI client.
* `cluster_reference`: - The UUID of the cluster that will host the iSCSI client. This field is read-only.
