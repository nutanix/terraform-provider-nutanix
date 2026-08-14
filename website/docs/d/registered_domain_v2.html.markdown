---
layout: "nutanix"
page_title: "NUTANIX: nutanix_registered_domain_v2"
sidebar_current: "docs-nutanix-datasource-registered-domain-v2"
description: |-
  Retrieves the registered domain by its external identifier.
---

# nutanix_registered_domain_v2

Retrieves the registered domain by its external identifier.

## Example Usage

```hcl
data "nutanix_registered_domain_v2" "get" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`:- (Required) External identifier of the registered domain.

## Attributes Reference

The following attributes are exported:

* `name`:- Name of the registered domain.
* `dns_name`:- DNS name of the registered domain.
* `fqdn`:- Fully qualified domain name of the registered domain.
* `domain_ext_id`:- The external ID of the domain.
* `location_ext_id`:- Location UUID of the domain.
* `credential_issuer`:- The entity that issued the credentials for the domain.
* `created_by`:- The user that created the registered domain on Nutanix Central.
* `api_credential_status`:- Status of the API credentials.
* `registration_state`:- Registration state of the domain.
* `registered_time`:- Date and time of domain registration.
* `platform_data`:- Platform data of the domain. See [Platform Data](#platform-data).
* `registration_config`:- Registration configuration of the domain. See [Registration Config](#registration-config).
* `tunnel`:- Tunnel configuration of the domain. See [Tunnel](#tunnel).
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity.
* `links`:- A HATEOAS style link for the response.

### Platform Data

The `platform_data` attribute exports the following:

* `cluster_count`:- Number of clusters in the domain.
* `connectivity_status`:- Domain connectivity status.
* `pc_env_type`:- PC Environment Type of the domain.
* `pc_provider_type`:- PC Provider Type of the domain.
* `pc_version`:- PC Version of the domain.

### Registration Config

The `registration_config` attribute exports the following:

* `domain_ext_id`:- Unique identifier for the Prism Central.
* `target_url`:- Prism Central should use this URL to communicate with Nutanix Central.
* `tenant_ext_id`:- Unique identifier of the tenant.
* `credentials`:- Registration credentials. See [Credentials](#credentials).

### Credentials

The `credentials` attribute exports the following:

* `api_key`:- API key which is used by Prism Central as authentication for Nutanix Central communication.
* `key_id`:- This is used by Prism Central to form authentication token for Nutanix Central communication.

### Tunnel

The `tunnel` attribute exports the following:

* `endpoint`:- Tunnel endpoint for connectivity.
* `provider`:- Tunnel provider.

### Links

The `links` attribute exports the following:

* `href`:- The URL at which the entity described by the link can be accessed.
* `rel`:- A name that identifies the relationship of the link to the object.

See detailed information in [Nutanix Registered Domain V2](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.0).
