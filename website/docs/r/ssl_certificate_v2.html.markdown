---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ssl_certificate_v2"
sidebar_current: "docs-nutanix-resource-ssl-certificate-v2"
description: |-
  Update the SSL certificate for a specific cluster
---

# nutanix_ssl_certificate_v2

Update the SSL certificate for a specific cluster. To update the SSL certificate for a specific cluster, you must provide a valid certificate payload in `Base64` format. You can either import a new SSL certificate or replace an existing one by supplying all necessary fields, including the Base64-encoded certificate and private key. Alternatively, you can regenerate a self-signed certificate by specifying the privateKeyAlgorithm, noting that only the RSA_2048 algorithm is supported for SSL certificate regeneration. This process helps maintain the security and integrity of your cluster's communications by allowing you to update or regenerate the SSL certificate as needed.

## Example Usage

```hcl
resource "nutanix_ssl_certificate_v2" "this" {
  cluster_ext_id = "6a68ecf4-8cac-42b3-9805-3186c2cecbd2"
  passphrase = "password"
  private_key = base64encode("private_key")
  public_certificate = base64encode("public_certificate")
  ca_chain = base64encode("ca_chain")
  private_key_algorithm = "RSA_2048"
}

```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id` (Required) — UUID of the cluster to which the host NIC belongs.
* `passphrase` (Optional) — Passphrase used for SSL certificate.
* `private_key` (Optional) — Private Key used for SSL certificate.
* `public_certificate` (Optional) — Public Certificate used for SSL certificate.
* `ca_chain` (Optional) — Description of the certificate authority (CA) chain.
* `private_key_algorithm` (Required) — Private Key Algorithm used for SSL certificate.

    | Enum         | Description                                                      |
    |--------------|------------------------------------------------------------------|
    | ECDSA_256    | Elliptic Curve Digital Signature Algorithm (256-bit)             |
    | ECDSA_384    | Elliptic Curve Digital Signature Algorithm (384-bit)             |
    | ECDSA_521    | Elliptic Curve Digital Signature Algorithm (521-bit)             |
    | JKS          | Java KeyStore (JKS) format                                       |
    | KRB_KEYTAB   | Kerberos Keytab format                                           |
    | PKCS12       | PKCS#12 format for certificate storage                           |
    | RSA_2048     | RSA 2048-bit encryption                                          |
    | RSA_4096     | RSA 4096-bit encryption                                          |
    | RSA_PUBLIC   | Public RSA key format                                            |

## API Reference

See detailed information in [Nutanix SSL Certificate v4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.1#tag/SSLCertificate/operation/updateSSLCertificate)
