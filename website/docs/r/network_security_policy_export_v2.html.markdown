---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_policy_export_v2"
sidebar_current: "docs-nutanix-resource-network-security-policy-export-v2"
description: |-
  Exports a point-in-time snapshot of the specified Network Security Policies.
---

# nutanix_network_security_policy_export_v2

Exports a point-in-time snapshot of the Network Security Policies identified by `policy_ext_ids`.

The serialized payload returned by the cluster is a **binary octet-stream**. It is exposed in two ways:

* `exported_payload` — the payload **base64-encoded** (so it is a valid, round-trippable Terraform string).
* `file_path` — the raw payload is written to this local path. The resulting file can be passed directly to the `path` argument of `nutanix_network_security_policy_import_v2` (for example, to replicate policies onto another cluster).

Because all inputs are `ForceNew`, changing the set of policies destroys and recreates the snapshot. Destroying this resource is a **no-op on the cluster** — it only removes the snapshot from Terraform state and never deletes the source policies.

## Example Usage

```hcl
# Export a snapshot of two network security policies and write the raw payload
# to a file that can be re-imported on another cluster.
resource "nutanix_network_security_policy_export_v2" "snapshot" {
  policy_ext_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002",
  ]
  file_path = "/tmp/policy_export.bin"
}

# Replicate the exported policies onto the target cluster.
# Terraform automatically infers the dependency via the file_path reference.
resource "nutanix_network_security_policy_import_v2" "replicated" {
  path = nutanix_network_security_policy_export_v2.snapshot.file_path
}

# Export scoped to a specific project.
resource "nutanix_network_security_policy_export_v2" "project_snapshot" {
  policy_ext_ids = ["00000000-0000-0000-0000-000000000001"]
  project_ext_id = "00000000-0000-0000-0000-000000000000"
  file_path      = "/tmp/project_export.bin"
}

# Export every network security policy on the cluster by omitting policy_ext_ids.
resource "nutanix_network_security_policy_export_v2" "all" {
  file_path = "/tmp/policy_export_all.bin"
}

output "exported_payload" {
  value = nutanix_network_security_policy_export_v2.snapshot.exported_payload
}
```

A complete, runnable example is available in the provider repository under `examples/network_security_policy_export_v2`.

## Argument Reference

The following arguments are supported:

* `policy_ext_ids` - (Optional, ForceNew) The list of network security policy external identifiers (UUIDs) to export. If omitted, all network security policies are exported. *Note: When omitted, it captures a snapshot of all policies at the time of creation. It will not automatically re-export if new policies are created out-of-band.*
* `project_ext_id` - (Optional, ForceNew) The external identifier of the project associated with the policies being exported.
* `file_path` - (Required, ForceNew) Local file path where the raw exported payload is written. The resulting file can be fed directly to the `path` argument of `nutanix_network_security_policy_import_v2`. *Note: Do not delete this file manually while the resource is managed by Terraform, as it may cause issues during the plan phase of dependent import resources.*

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `exported_payload` - The exported network security policy payload returned by the cluster, base64-encoded (the raw payload is a binary octet-stream).
* `task_ext_id` - A globally unique identifier for the task created by the export operation.

## Lifecycle

* **Create** — triggers the asynchronous export task, waits for it to succeed, then downloads and stores the serialized payload in `exported_payload`.
* **Read** — verifies that the source policies still exist on the cluster. If all of them have been deleted out-of-band, the resource is removed from state so it can be re-exported.
* **Update** — not applicable; all inputs are `ForceNew`.
* **Delete** — no-op on the cluster; only removes the snapshot from Terraform state.

See detailed information in [Nutanix Network Security Policy Export V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.4#tag/NetworkSecurityPolicies/operation/exportNetworkSecurityPolicy).