---
layout: "nutanix"
page_title: "NUTANIX: nutanix_promote_protected_resource_v2"
sidebar_current: "docs-nutanix-resource-promote-protected-resource-v2"
description: |-
  Promotes the specified synced entity at the target site. This is only relevant if the synced entity is protected in a synchronous schedule.

---

# nutanix_promote_protected_resource_v2

Promotes the specified synced entity at the target site. This is only relevant if the synced entity is protected in a synchronous schedule.


## Example:

```hcl

# Promote a protected virtual machine on remote site
# This example promotes a protected virtual machine on a remote site.
# Steps:
# 1. Define the provider for the remote site
# 2. Create a category and a protection policy, on the local site
# 3. Create a virtual machine and associate it with the protection policy, on local site
# 4. Promote the protected virtual machine on the remote site

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

# promote the protected virtual machine on the remote site
resource "nutanix_promote_protected_resource_v2" "promote-example" {
  provider = nutanix.remote
  ext_id   = "d22529bb-f02d-4710-894b-d1de772d7832" # protected resource (VM or VG) ext_id
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.


See detailed information in [Nutanix Promote Protected Resource v4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/ProtectedResources/operation/promoteProtectedResource).
