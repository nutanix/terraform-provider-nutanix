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
    azure_key_vault {
      endpoint_url           = "https://example-url.vault.azure.net/"
      key_id                 = "example_key_id"
      tenant_id              = "d4c8a8e5-91b3-4f7e-8c2e-77d6f4a22f11"
      client_id              = "e29f9c62-3e56-41d0-b123-7f8a22c0cdef"
      client_secret          = "7Z3uQ~vO4trhXk8B5M9qjwgT1pR2uC9yD1zF0wX3"
      credential_expiry_date = "2026-09-01"
    }
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

- `name`: - (Required) Name of the key management server (KMS).
- `access_information`: - (Required) KMS Access information, it can be Azure Key Vault access information or KMIP based External Key Manager Access Information.

### access_information

The `access_information` it must be one of `azure_key_vault` or `kmip_key_vault`.

- `azure_key_vault`: - (Optional) Azure Key Vault access information.
- `kmip_key_vault`: - (Optional) KMIP based External Key Manager Access Information.

#### azure_key_vault

The `azure_key_vault` supports the following:

- `endpoint_url`: - (Required) Endpoint URL for the Azure Key Vault.
- `key_id`: - (Required) Master key identifier for the Azure Key Vault.
- `tenant_id`: - (Required) Tetant identifier for the Azure Key Vault.
- `client_id`: - (Required) Client identifier for the Azure Key Vault.
- `client_secret`: - (Required) Client secret for the Azure Key Vault.
- `credential_expiry_date`: - (Required) When the client secret is going to expire.

#### kmip_key_vault

The `kmip_key_vault` supports the following:

- `ca_name`: - (Required) Name of the CA.
- `ca_pem`: - (Required) CA PEM.
- `cert_pem`: - (Required) Cert PEM.
- `private_key`: - (Required) Private key.
- `endpoints`: - (Required) List of endpoints of the External Key Manager server.
  - `ip_address`: - (Required) IP address of the External Key Manager server.
    - `ipv4`: - (Optional) IPv4 address of the External Key Manager server.
      - `value`: - (Required) IPv4 address of the External Key Manager server.
      - `prefix_length`: - (Optional) Prefix length of the IPv4 address. 
    - `ipv6`: - (Optional) IPv6 address of the External Key Manager server.
      - `value`: - (Required) IPv6 address of the External Key Manager server.
      - `prefix_length`: - (Optional) Prefix length of the IPv6 address.
    - `fqdn`: - (Optional) FQDN of the External Key Manager server.
      - `value`: - (Optional) FQDN of the External Key Manager server.
  - `port`: - (Required) Port of the External Key Manager server.

See detailed information in [Nutanix Create a key management server V4](https://developers.nutanix.com/api-reference?namespace=security&version=v4.1#tag/KeyManagementServers/operation/createKeyManagementServer)
