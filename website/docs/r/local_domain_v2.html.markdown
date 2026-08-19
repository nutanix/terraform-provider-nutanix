---
layout: "nutanix"
page_title: "NUTANIX: nutanix_local_domain_v2"
sidebar_current: "docs-nutanix-resource-local-domain-v2"
description: |-
  Registers the domain with Nutanix Central, establishing a connection between Prism Central and their corresponding Nutanix Central components.
---

# nutanix_local_domain_v2

Registers the domain with Nutanix Central, establishing a connection between Prism Central and their corresponding Nutanix Central components. The local domain is a singleton entity: there is exactly one local domain per Prism Central. Creating this resource registers the domain and destroying it unregisters the domain from Nutanix Central.

## Example Usage

```hcl
resource "nutanix_local_domain_v2" "example" {
  credentials {
    api_key = "nutanix-central-api-key"
    key_id  = "nutanix-central-key-id"
  }
  domain_ext_id            = "00000000-0000-0000-0000-000000000000"
  identity_provider_ext_id = "11111111-1111-1111-1111-111111111111"
  target_url               = "https://nutanix-central.example.com"
  tenant_ext_id            = "22222222-2222-2222-2222-222222222222"
}
```

## Argument Reference

The following arguments are supported:

* `credentials`:- (Required) Credentials used by Prism Central to authenticate with Nutanix Central. See [Credentials](#credentials) below.
* `domain_ext_id`:- (Required) Unique identifier for the Prism Central.
* `identity_provider_ext_id`:- (Required) Identity provider unique identifier. Identity provider can be directory-service or SAML based identity provider.
* `target_url`:- (Required) Prism Central should use this URL to communicate with Nutanix Central.
* `tenant_ext_id`:- (Required) Unique identifier of the tenant.

### Credentials

The `credentials` attribute supports the following:

* `api_key`:- (Required, Sensitive) API key which is used by Prism Central as authentication for Nutanix Central communication.
* `key_id`:- (Required, Sensitive) This is used by Prism Central to form authentication token for Nutanix Central communication.

## Attributes Reference

The following attributes are exported:

* `name`:- Name of the domain on Nutanix Central.
* `target_url`:- Prism Central should use this URL to communicate with Nutanix Central.
* `connection_status`:- Functional connectivity between Prism Central and Nutanix Central (and its components). See [Connection Status](#connection-status) below.
* `registration_status`:- The registration status of the domain with Nutanix Central. See [Registration Status](#registration-status) below.
* `location`:- Location details associated with the domain. See [Location](#location) below.

### Connection Status

The `connection_status` attribute exports the following:

* `status`:- Connectivity status for the domain. One of `ACTIVE`, `INACTIVE`.
* `service_statuses`:- Connectivity status of each dependent service with Nutanix Central. See [Service Statuses](#service-statuses) below.

### Service Statuses

The `service_statuses` attribute exports the following:

* `service_name`:- Local domain dependent service name. One of `TUNNEL_CONNECTION`, `IDENTITY_MANAGEMENT`.
* `status`:- Connectivity status for the service. One of `ACTIVE`, `INACTIVE`.

### Registration Status

The `registration_status` attribute exports the following:

* `days_to_expiry`:- Indicates remaining days to domain credential expiry, 0 means expired.
* `error_stage`:- The stage for the local domain's register/unregister operation.
* `identity_provider_ext_id`:- Identity provider unique identifier. Identity provider can be directory-service or SAML based identity provider.
* `state`:- Registration state. One of `UNREGISTERED`, `REGISTERING`, `REGISTERED`, `UNREGISTERING`, `REGISTRATION_ERROR`, `UNREGISTRATION_ERROR`, `DRAFT`.

### Location

The `location` attribute exports the following:

* `ext_id`:- A globally unique identifier of the location.
* `name`:- Location name.
* `country`:- Location Country.
* `state`:- Location state.
* `coordinates`:- The location coordinates, that is latitude and longitude. See [Coordinates](#coordinates) below.

### Coordinates

The `coordinates` attribute exports the following:

* `latitude`:- The latitude part of the location.
* `longitude`:- The longitude part of the location.

See detailed information in [Nutanix Register Local Domain V4](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.0).
