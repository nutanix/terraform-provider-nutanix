---
layout: "nutanix"
page_title: "NUTANIX: nutanix_key_management_server_v2"
sidebar_current: "docs-nutanix-datasource-key-management_server_v2"
description: |-
  Provides a Nutanix Key Management Server datasource to get details of a key management server by ext_id. A Key Management Server (KMS) is a centralized system that securely generates, stores, and manages cryptographic keys used for data encryption. When data-at-rest encryption is enabled on Nutanix clusters, KMS ensures that encryption keys are protected and accessible only to authorized entities, providing an additional layer of security for sensitive data.
---

# nutanix_key_management_server_v2

Fetches the details of a key management server identified by its unique identifier (ext_id).  


## Example

```hcl
data "nutanix_key_management_server_v2" "kms"{
  ext_id = "aa300b88-8560-4eb3-ba6c-49b0ff8c9cc1"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: (Required) Unique identifier for the key management server of type UUID.

## Attribute Reference

The Following attributes are exported:

- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `name`: - Name of the key management server (KMS).
- `access_information`: -KMS Access information, it can be Azure Key Vault access information or KMIP based External Key Manager Access Information.
- `creation_timestamp`: - The timestamp when the key management server was created.

### access_information

The `access_information` exports the following:

- `azure_key_vault`: - Access information for the Azure Key Vault.
- `kmip_key_vault`: - Access information for the KMIP Key Vault.

#### azure_key_vault

The `azure_key_vault` exports the following:

- `endpoint_url`: Endpoint URL for the Azure Key Vault.
- `key_id`: Master key identifier for the Azure Key Vault.
- `tenant_id`: Tetant identifier for the Azure Key Vault.
- `client_id`: Client identifier for the Azure Key Vault.
- `truncated_client_secret`: Truncated client secret for the Azure Key Vault.
- `credential_expiry_date`: When the client secret is going to expire.

#### kmip_key_vault

The `kmip_key_vault` exports the following:

- `ca_name`: Name of the Certificate Authority.
- `ca_pem`: Cert PEM File.
- `cert_pem`: Cert PEM.
- `endpoints`: List of endpoints of the External Key Manager server.
  - `ip_address`: IP address of the External Key Manager server.
    - `ipv4`: IPv4 address of the External Key Manager server.
      - `value`: IPv4 address of the External Key Manager server.
      - `prefix_length`: Prefix length of the IPv4 address.
    - `ipv6`: IPv6 address of the External Key Manager server.
      - `value`: IPv6 address of the External Key Manager server.
      - `prefix_length`: Prefix length of the IPv6 address.
    - `fqdn`: FQDN of the External Key Manager server.
      - `value`: FQDN of the External Key Manager server.
  - `port`: Port of the External Key Manager server.

See detailed information in [Nutanix Get details of a key management server V4](https://developers.nutanix.com/api-reference?namespace=security&version=v4.1#tag/KeyManagementServers/operation/getKeyManagementServerById)
