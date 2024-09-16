---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_vm_v2"
sidebar_current: "docs-nutanix-resource-volume-group-vm-attachments-v4"
description: |-
  This operation submits a request to Attaches VM to a Volume Group identified by {extId}.
---

# nutanix_volume_group_vm_v2

Provides a resource to Create a new Volume Group.

## Example Usage

``` hcl
data "nutanix_clusters" "clusters"{
}

#get desired cluster data from setup
locals {
  cluster1 = [
    for cluster in data.nutanix_clusters.clusters.entities :
    cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
  ][0]
}

resource "nutanix_volume_group_v2" "test"{
  name                               = "test_volume_group"
  cluster_reference                  = local.cluster1
}

resource "nutanix_volume_group_vm_v2" "vg_vm_example"{
  volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
  vm_ext_id           = var.vg_vm_ext_id
}

```

## Argument Reference
The following arguments are supported:


* `volume_group_ext_id`: -(Required) The external identifier of the volume group.
* `vm_ext_id`: -(Required) A globally unique identifier of an instance that is suitable for external consumption. 
* `index`: -(Optional) The index on the SCSI bus to attach the VM to the Volume Group. 


See detailed information in [Nutanix Volumes](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0.b1).
