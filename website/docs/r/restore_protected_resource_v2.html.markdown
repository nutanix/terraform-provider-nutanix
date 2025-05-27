---
layout: "nutanix"
page_title: "NUTANIX: nutanix_restore_protected_resource_v2"
sidebar_current: "docs-nutanix-resource-restore-protected-resource-v2"
description: |-
  Restore the specified protected resource from its state at the given timestamp on the given cluster. This is only relevant if the entity is protected in a minutely schedule at the given timestamp.



---

# nutanix_restore_protected_resource_v2

Restore the specified protected resource from its state at the given timestamp on the given cluster. This is only relevant if the entity is protected in a minutely schedule at the given timestamp.


## Example 1: Restore Virtual Machine

```hcl
# Restore a protected virtual machine on remote site
# This example demonstrates how to restore a protected virtual machine on remote site.
# steps:
# 1. Define the provider for the remote site
# 2. Create a category and a protection policy, on the local site
# 3. Create a virtual machine and associate it with the protection policy, on local site
# 4. Restore the virtual machine on the remote site


# define another alias for the provider, this time for the remote PC
provider "nutanix" {
  alias    = "remote"
  username = var.nutanix_remote_username
  password = var.nutanix_remote_password
  endpoint = var.nutanix_remote_endpoint
  insecure = true
  port     = 9440
}

# create a category and a protection policy on the local site

# restore the protected virtual machine on the remote site
resource "nutanix_restore_protected_resource_v2" "rp-vm" {
  provider       = nutanix.remote
  ext_id         = "d22529bb-f02d-4710-894b-d1de772d7832" # protected vm ext_id
  cluster_ext_id = "0005b6b1-1b16-4983-b5ff-204840f85e07" # remote cluster ext_id
}

```

## Example 2: Restore Volume Group

```hcl
# Restore a protected volume group on remote site
# This example demonstrates how to restore a protected volume group on remote site.
# steps:
# 1. Define the provider for the remote site
# 2. Create a category and a protection policy, on the local site
# 3. Create a volume group and associate it with the category on the local site
# 4. Restore the volume group on the remote site


# define another alias for the provider, this time for the remote PC
provider "nutanix" {
  alias    = "remote"
  username = var.nutanix_remote_username
  password = var.nutanix_remote_password
  endpoint = var.nutanix_remote_endpoint
  insecure = true
  port     = 9440
}

# create a category , a protection policy and VG on the local site

# restore the protected volume group on the remote site
resource "nutanix_restore_protected_resource_v2" "rp-vg" {
  provider       = nutanix.remote
  ext_id         = "246c651a-1b16-4983-b5ff-204840f85e07" # protected volume group ext_id
  cluster_ext_id = "0005b6b1-1b16-4983-b5ff-204840f85e07" # remote cluster ext_id
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.
* `cluster_ext_id`: -(Required) The external identifier of the cluster on which the entity has valid restorable time ranges. The restored entity will be created on the same cluster.
* `restore_time`: -(Optional) UTC date and time in ISO 8601 format representing the time from when the state of the entity should be restored. This needs to be a valid time within the restorable time range(s) for the protected resource.


See detailed information in [Nutanix Restore Protected Resource v4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/ProtectedResources/operation/restoreProtectedResourcen).

