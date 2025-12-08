---
layout: "nutanix"
page_title: "NUTANIX: nutanix_key_management_servers_v2"
sidebar_current: "docs-nutanix-datasource-key-management_servers_v2"
description: |-
  Provides a Nutanix Key Management Server datasource to list key management servers. A Key Management Server (KMS) is a centralized system that securely generates, stores, and manages cryptographic keys used for data encryption. When data-at-rest encryption is enabled on Nutanix clusters, KMS ensures that encryption keys are protected and accessible only to authorized entities, providing an additional layer of security for sensitive data.


---

# nutanix_key_management_servers_v2

A Key Management Server (KMS) is a critical component in enterprise data security that handles the lifecycle of cryptographic keys—including creation, storage, rotation, and deletion. In Nutanix environments, KMS integrates with data-at-rest encryption to protect sensitive information stored on clusters. By centralizing key management, organizations can enforce consistent security policies, meet compliance requirements (such as HIPAA, PCI-DSS, and GDPR), and maintain control over who can access encrypted data.

This datasource provides a comprehensive list of all key management servers configured in your environment, including their access details and relevant attributes.


## Example

```hcl
data "nutanix_key_management_servers_v2" "kms-list"{}
```


## Attribute Reference

The Following attributes are exported:

- `kms`: - List of key management servers (KMS).


### kms
The `kms` attribute export the following:

- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `name`: - Name of the key management server (KMS).
- `access_information`: - Access information for the Azure Key Vault.

#### access_information

The `access_information` exports the following:

- `endpoint_url`: (Required) Endpoint URL for the Azure Key Vault.
- `key_id`: (Required) Master key identifier for the Azure Key Vault.
- `tenant_id`: (Required) Tetant identifier for the Azure Key Vault.
- `client_id`: (Required) Client identifier for the Azure Key Vault.
- `client_secret`: (Required) Client secret for the Azure Key Vault.
- `credential_expiry_date`: (Required) When the client secret is going to expire.

See detailed information in [Nutanix List key management servers V4](https://developers.nutanix.com/api-reference?namespace=security&version=v4.0#tag/KeyManagementServers/operation/listKeyManagementServers)
