---
layout: "nutanix"
page_title: "NUTANIX: nutanix_local_domain_v2"
sidebar_current: "docs-nutanix-datasource-local-domain-v2"
description: |-
  Returns the domain entity.
---

# nutanix_local_domain_v2

Returns the domain entity. The local domain is a singleton entity per Prism Central, so this data source takes no input arguments.

## Example Usage

```hcl
data "nutanix_local_domain_v2" "get-local-domain" {}
```

## Argument Reference

This data source takes no arguments — the local domain is a singleton.

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

See detailed information in [Nutanix Get Local Domain V4](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.0).
