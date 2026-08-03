---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_policy_import_v2"
sidebar_current: "docs-nutanix-resource-network-security-policy-import-v2"
description: |-
  Imports all the Network Security Policies specified by the data file.
---

# nutanix_network_security_policy_import_v2

Imports all the Network Security Policies specified by the data file. This is a resource backed by an asynchronous task.

The data file is produced by [`nutanix_network_security_policy_export_v2`](network_security_policy_export_v2.html). `path` must point to an existing, readable file on the machine running Terraform; if it does not, the apply fails with `The provided path '<path>' is not a valid file` before any task is started.

The policies created by the import are tracked in the computed `imported_policy_ext_ids` attribute. Because all inputs are `ForceNew`, changing the data file destroys the previously imported policies and runs a fresh import. Destroying this resource deletes the imported policies from the cluster.

## Example Usage

```hcl
# Bulk import of network security policies from a data file.
resource "nutanix_network_security_policy_import_v2" "example" {
  path                = "/path/to/policy_export.bin"
  ntnx_purge_policies = false
}

# Import with project scoping and purge of existing policies.
resource "nutanix_network_security_policy_import_v2" "with_project" {
  path                = "/path/to/policy_export.bin"
  ntnx_purge_policies = true
  ntnx_project_ext_id = "a2b3c4d5-e6f7-8901-2345-678901234567"
}
```

### Export then import (round trip)

```hcl
# Export a policy to a file ...
resource "nutanix_network_security_policy_export_v2" "export" {
  policy_ext_ids = ["00000000-0000-0000-0000-000000000001"]
  file_path      = "${path.module}/nsp_export.bin"
}

# ... then import that file back (for example, onto another cluster).
# Terraform automatically infers the dependency via the file_path reference.
resource "nutanix_network_security_policy_import_v2" "import" {
  path                = nutanix_network_security_policy_export_v2.export.file_path
  ntnx_purge_policies = false
}
```

A complete, runnable example is available in the provider repository under `examples/network_security_policy_import_v2`.

## Argument Reference

The following arguments are supported:

* `path` - (Required, ForceNew) The local file path of the data file to import network security policies. It must be an existing, readable file.
* `ntnx_purge_policies` - (Optional, ForceNew) Specifies whether the existing policies are deleted (`true`) or retained (`false`) upon network security policy import. Defaults to `false`.
* `ntnx_project_ext_id` - (Optional, ForceNew) The project external identifier to associate with the policies being imported.
* `dryrun` - (Optional, ForceNew) When `true`, the import is validated without persisting any changes. *Note: When enabled, `imported_policy_ext_ids` will be empty, and destroying the resource will be a no-op on the cluster.*

## Attributes Reference

* `task_ext_id` - A globally unique identifier for the task created by the import operation.
* `imported_policy_ext_ids` - The list of network security policy external identifiers (UUIDs) that were created by this import.

## Lifecycle

* **Create** — validates that `path` is a readable file, triggers the asynchronous import task, waits for it to succeed, and records the created policies in `imported_policy_ext_ids`.
* **Read** — verifies the imported policies still exist; prunes any deleted out-of-band, and removes the resource from state if all are gone.
* **Update** — not applicable; all inputs are `ForceNew`.
* **Delete** — deletes each imported network security policy from the cluster and waits for the deletion tasks to complete.

See detailed information in [Nutanix Network Security Policy Import V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.4#tag/NetworkSecurityPolicies/operation/applyNetworkSecurityPolicyImport).