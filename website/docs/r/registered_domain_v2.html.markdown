---
layout: "nutanix"
page_title: "NUTANIX: nutanix_registered_domain_v2"
sidebar_current: "docs-nutanix-resource-registered-domain-v2"
description: |-
  Creates a new registered domain and saves it to database.
---

# nutanix_registered_domain_v2

Creates a new registered domain and saves it to database. A registered domain represents a Prism Central (PC) instance registered with a Nutanix Central (NC) instance.

> **WriteOnly credentials:** `registration_config.credentials.api_key`, `registration_config.credentials.key_id` are write-only. They are never stored in Terraform state. Re-supply them in configuration on every apply. Do not use `lifecycle.ignore_changes` on these attributes. See the [WriteOnly Credentials](/docs/providers/nutanix/guides/writeonly_credentials.html) guide for upgrade caveats.

## Example Usage

```hcl
resource "nutanix_registered_domain_v2" "example" {
  name               = "example-registered-domain"
  fqdn               = "pc.example.com"
  dns_name           = "pc.example.com"
  registration_state = "REGISTERED"

  platform_data {
    pc_version          = "pc.7.6"
    connectivity_status = "REACHABLE"
    cluster_count       = 1
  }

  registration_config {
    target_url = "https://nc.example.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name`:- (Optional) Name of the registered domain.
* `dns_name`:- (Optional) DNS name of the registered domain.
* `fqdn`:- (Optional) Fully qualified domain name of the registered domain.
* `domain_ext_id`:- (Optional) The external ID of the domain.
* `location_ext_id`:- (Optional) Location UUID of the domain.
* `credential_issuer`:- (Optional) The entity that issued the credentials for the domain to use as authentication for Nutanix central communication.
* `created_by`:- (Optional) The user that created the registered domain on Nutanix Central.
* `api_credential_status`:- (Optional) Status of the API credentials. One of `VALID`, `NEAR_EXPIRY`, `EXPIRED`.
* `registration_state`:- (Optional) Registration state of the domain. One of `UNREGISTERED`, `REGISTERING`, `REGISTERED`, `UNREGISTERING`, `REGISTRATION_ERROR`, `UNREGISTRATION_ERROR`, `DRAFT`.
* `registered_time`:- (Optional) Date and time of domain registration (RFC3339).
* `platform_data`:- (Optional) Platform data of the domain. See [Platform Data](#platform-data).
* `registration_config`:- (Optional) Registration configuration of the domain. See [Registration Config](#registration-config).
* `tunnel`:- (Optional) Tunnel configuration of the domain. See [Tunnel](#tunnel).

### Platform Data

The `platform_data` attribute supports the following:

* `cluster_count`:- (Optional) Number of clusters in the domain.
* `connectivity_status`:- (Optional) Domain connectivity status.
* `pc_env_type`:- (Optional) PC Environment Type of the domain.
* `pc_provider_type`:- (Optional) PC Provider Type of the domain.
* `pc_version`:- (Optional) PC Version of the domain.

### Registration Config

The `registration_config` attribute supports the following:

* `domain_ext_id`:- (Optional) Unique identifier for the Prism Central.
* `target_url`:- (Optional) Prism Central should use this URL to communicate with Nutanix Central.
* `tenant_ext_id`:- (Optional) Unique identifier of the tenant.
* `credentials`:- (Optional) Registration credentials. See [Credentials](#credentials).

### Credentials

The `credentials` attribute supports the following:

* `api_key`:- (Optional, WriteOnly) API key which is used by Prism Central as authentication for Nutanix Central communication.
* `key_id`:- (Optional, WriteOnly) This is used by Prism Central to form authentication token for Nutanix Central communication.

### Tunnel

The `tunnel` attribute supports the following:

* `endpoint`:- (Optional) Tunnel endpoint for connectivity. Prism Central make use of this endpoint to reach tunnel.
* `provider`:- (Optional) Tunnel provider. One of `IKAT`.

## Attributes Reference

The following attributes are exported:

* `ext_id`:- A globally unique identifier of the registered domain.
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity.
* `links`:- A HATEOAS style link for the response.

### Links

The `links` attribute supports the following:

* `href`:- The URL at which the entity described by the link can be accessed.
* `rel`:- A name that identifies the relationship of the link to the object.

## Import

Registered domains can be imported using their `ext_id`, e.g.

```
terraform import nutanix_registered_domain_v2.example <ext_id>
```

See detailed information in [Nutanix Registered Domain V2](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.4.b1#tag/RegisteredDomains/operation/createRegisteredDomain).
