---
layout: "nutanix"
page_title: "NUTANIX: nutanix_key_management_server_v2"
sidebar_current: "docs-nutanix-resource-key-management-server-v2"
description: |-
  Provides a Nutanix Key Management Server resource to Create a Key Management Server.
---

# nutanix_key_management_server_v2

Creates a new key management server with the specified access credentials.

## Example

```hcl
resource "nutanix_key_management_server_v2" "kms"{
  name = "tf-kms-example"
  access_information {
    endpoint_url           = "https://example-url.vault.azure.net/"
    key_id                 = "example_key_id"
    tenant_id              = "d4c8a8e5-91b3-4f7e-8c2e-77d6f4a22f11"
    client_id              = "e29f9c62-3e56-41d0-b123-7f8a22c0cdef"
    client_secret          = "7Z3uQ~vO4trhXk8B5M9qjwgT1pR2uC9yD1zF0wX3"
    credential_expiry_date = "2026-09-01"
  }
  lifecycle {
    ignore_changes = [
      access_information[0].client_secret,
      access_information[0].key_id
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

- `name`: -(Required) Name of the key management server (KMS).
- `access_information`: -(Required) Access information for the Azure Key Vault.

### access_information

The `access_information` supports the following:

- `endpoint_url`: -(Required) Endpoint URL for the Azure Key Vault.
- `key_id`: -(Required) Master key identifier for the Azure Key Vault.
- `tenant_id`: -(Required) Tetant identifier for the Azure Key Vault.
- `client_id`: -(Required) Client identifier for the Azure Key Vault.
- `client_secret`: -(Required) Client secret for the Azure Key Vault.
- `credential_expiry_date`: -(Required) When the client secret is going to expire.

See detailed information in [Nutanix Create a key management server V4](https://developers.nutanix.com/api-reference?namespace=security&version=v4.0#tag/KeyManagementServers/operation/createKeyManagementServer)
