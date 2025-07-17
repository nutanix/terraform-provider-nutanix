---
layout: "nutanix"
page_title: "NUTANIX: nutanix_certificate_v2 "
sidebar_current: "docs-nutanix-resource-object-store-certificate-v2"
description: |-
  Create a SSL certificate for an Object store

---

# nutanix_certificate_v2

This operation creates a new default certificate and keys. It also creates the alternate FQDNs and alternate IPs for the Object store. The certificate of an Object store can be created when it is in a OBJECT_STORE_AVAILABLE or OBJECT_STORE_CERT_CREATION_FAILED state. If the publicCert, privateKey, and ca values are provided in the request body, these values are used to create the new certificate. If these values are not provided, a new certificate will be generated if 'shouldGenerate' is set to true and if it is set to false, the existing certificate will be used as the new certificate. Optionally, a list of additional alternate FQDNs and alternate IPs can be provided. These alternateFqdns and alternateIps must be included in the CA certificate if it has been provided.



## Example Usage

```hcl
data "nutanix_certificate_v2" "example"{
  object_store_ext_id = "ac91151a-28b4-4ffe-b150-6bcb2ec80cd4"
  ext_id              = "ef0a9a54-e7e1-42e2-a59f-de779ec1c9ea"
}

```

## Argument Reference

The following arguments are supported:

- `object_store_ext_id`: -(Required) The UUID of the Object store.
- `ext_id`: -(Required) The UUID of the certificate of an Object store.

## Attributes Reference

The following attributes are exported:

- `path`: -(Required) Path to a JSON file which contains the public certificates, private key, and CA certificate or chain, along with a list of alternate FQDNs and alternate IPs to create a certificate for the Object store.

The Content of the JSON file :
| Field           | Description                                                                                                                                                                                                 |
|----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `alternateFqdns` | The list of alternate FQDNs for accessing the Object store. The FQDNs must consist of at least 2 parts separated by a '.'. Each part can contain upper and lower case letters, digits, hyphens or underscores but must begin and end with a letter. Each part can be up to 63 characters long. For example: `objects-0.pc_nutanix.com`. |
| `alternateIps`   | A list of the IPs included as Subject Alternative Names (SANs) in the certificate. The IPs must be among the public IPs of the Object store (`publicNetworkIps`).                                        |
| `ca`             | The CA certificate or chain to upload.                                                                                                                                                                    |
| `publicCert`     | The public certificate to upload.                                                                                                                                                                          |
| `privateKey`     | The private key to upload.                                                                                                                                                                                 |
| `shouldGenerate` | If true, a new certificate is generated with the provided alternate FQDNs and IPs.                                                                                                                        |

## JSON Example
```json
{
  "alternateFqdns": [
    {
      "value": "fqdn1.example.com"
    },
    {
      "value": "fqdn2.example.com"
    }
  ],
  "alternateIps": [
    {
      "ipv4": {
        "value": "192.168.1.1"
      }
    },
    {
      "ipv4": {
         "value": "192.168.1.2"
      }
    }
  ],
  "shouldGenerate": true,
  "ca": "-----BEGIN CERTIFICATE-----\nMIIDzTCCArWgAwIBAgIUI...\n-----END CERTIFICATE-----",
  "publicCert": "-----BEGIN CERTIFICATE-----\nMIIDzTCCArWgAwIBAgIUI...\n-----END CERTIFICATE-----",
  "privateKey": "-----BEGIN RSA PRIVATE KEY-----\nMIIDzTCCArWgAwIBAgIUI...\n-----END RSA PRIVATE KEY-----"
}
```

See detailed information in [Nutanix Create a SSL certificate for an Object store V4 ](https://developers.nutanix.com/api-reference?namespace=objects&version=v4.0#tag/ObjectStores/operation/createCertificate).
