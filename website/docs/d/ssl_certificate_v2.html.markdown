---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ssl_certificate_v2"
sidebar_current: "docs-nutanix-datasource-ssl-certificate-v2"
description: |-
  Get SSL certificate details for a specified cluster.
---

# nutanix_ssl_certificate_v2

Provides detailed information about the SSL certificate in privacy-enhanced mail (.pem) format for the specified cluster.

## Example Usage

```hcl
data "nutanix_ssl_certificate_v2" "this" {
  cluster_ext_id = "6a68ecf4-8cac-42b3-9805-3186c2cecbd2"
}

```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id` (Required) — UUID of the cluster to which the host NIC belongs.

## Attribute Reference

The following attributes are exported:

* `public_certificate` — Details about the public SSL certificate.
* `private_key_algorithm` — Private Key Algorithm used for SSL certificate.

    | Enum                      | Description                                |
    |---------------------------|--------------------------------------------|
    | ECDSA_256 | Elliptic Curve Digital Signature Algorithm (256-bit) description.         |
    | JKS           | Java KeyStore (JKS) format description.                   |
    | RSA_2048          | Description for RSA 2048-bit encryption.                  |
    | ECDSA_521          | ECDSA (521-bit) description.                  |
    | KRB_KEYTAB                | Kerberos Keytab format description.                 |
    | PKCS12          | PKCS#12 format description for certificate storage.                  |
    | RSA_4096       | Description for RSA 4096-bit encryption.               |
    | RSA_PUBLIC | Public RSA key format description.               |
    | ECDSA_384 | ECDSA (384-bit) description.            |

## API Reference

See detailed information in [Nutanix SSL Certificates v4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.1#tag/SSLCertificate/operation/getSSLCertificate)
