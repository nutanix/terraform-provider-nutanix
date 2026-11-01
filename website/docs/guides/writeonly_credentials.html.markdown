---
layout: "nutanix"
page_title: "WriteOnly Credentials"
sidebar_current: "docs-nutanix-guides-writeonly-credentials"
description: |-
  WriteOnly credential attributes in the Nutanix Terraform provider: behavior, upgrade caveats, and brownfield notes.
---

# WriteOnly Credentials

Several Nutanix Terraform provider resources mark passwords, API keys, private keys, and similar secrets as **WriteOnly**. WriteOnly attributes are accepted in configuration and sent to the API, but they are **never stored in Terraform state**.

## Behavior

- Keep WriteOnly values in configuration, variables, or CI secrets. They will not round-trip from `terraform.tfstate`.
- Re-supply WriteOnly attributes in configuration on every apply when they are required for create or update.
- Do **not** put WriteOnly attributes under `lifecycle.ignore_changes`. Terraform cannot reconcile those values from state.
- Datasources that return credentials may still use `Sensitive` (readable from state). Resources omit WriteOnly values from state after apply.

## Affected resources

| Resource | WriteOnly attributes |
|----------|----------------------|
| [`nutanix_ssl_certificate_v2`](/docs/providers/nutanix/r/ssl_certificate_v2.html) | `passphrase`, `private_key` |
| [`nutanix_foundation_image_nodes`](/docs/providers/nutanix/r/foundation_image_nodes.html) | `ipmi_password`, `ucsm_password`, `xs_master_password`, `hypervisor_password`, `unc_password`, `fc_settings.fc_metadata.api_key`, node-level `ipmi_password` |
| [`nutanix_foundation_ipmi_config`](/docs/providers/nutanix/r/foundation_ipmi_config.html) | `ipmi_password` |
| [`nutanix_registered_domain_v2`](/docs/providers/nutanix/r/registered_domain_v2.html) | `registration_config.credentials.api_key`, `registration_config.credentials.key_id` |
| [`nutanix_local_domain_v2`](/docs/providers/nutanix/r/local_domain_v2.html) | `credentials.api_key`, `credentials.key_id` |
| [`nutanix_password_change_request_v2`](/docs/providers/nutanix/r/password_change_request_v2.html) | `current_password`, `new_password` |
| [`nutanix_directory_services_v2`](/docs/providers/nutanix/r/directory_services_v2.html) | `service_account.password` |
| [`nutanix_pc_backup_target_v2`](/docs/providers/nutanix/r/pc_backup_target_v2.html) | object store `credentials.access_key_id`, `credentials.secret_access_key` |
| [`nutanix_pc_deploy_v2`](/docs/providers/nutanix/r/pc_deploy_v2.html) | `config.credentials.password` |
| [`nutanix_pc_registration_v2`](/docs/providers/nutanix/r/pc_registration_v2.html) | remote cluster `authentication.password` |
| [`nutanix_pc_restore_source_v2`](/docs/providers/nutanix/r/pc_restore_source_v2.html) | object store `credentials.access_key_id`, `credentials.secret_access_key` |
| [`nutanix_key_management_server_v2`](/docs/providers/nutanix/r/key_management_server_v2.html) | `client_secret`, `cert_pem`, `private_key`, `ca_pem` |
| [`nutanix_tenant_v2`](/docs/providers/nutanix/r/tenant_v2.html) | `admin_info.password` |

Each resource page also links this guide and marks the attributes as WriteOnly in Argument Reference.

## Upgrade caveats (brownfield)

When upgrading from an older provider that stored these secrets as Sensitive-only values in state:

- Secrets no longer persist in `terraform.tfstate`. Ensure they remain available in config, variables, or CI secrets before the next apply.
- Do not add `lifecycle.ignore_changes` for WriteOnly attributes after upgrade.
- For `nutanix_pc_backup_target_v2` and `nutanix_pc_restore_source_v2`, expect a one-time in-place plan to re-supply object-store credentials.
- For `nutanix_pc_restore_source_v2`, `ext_id` may change when credentials are re-applied (API recreate behavior).
