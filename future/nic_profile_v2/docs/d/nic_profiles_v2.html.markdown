---
layout: "nutanix"
page_title: "NUTANIX: nutanix_nic_profiles_v2"
sidebar_current: "docs-nutanix-datasource-nic-profiles-v2"
description: |-
  Lists all NIC Profiles with host NICs and capability.
---

# nutanix_nic_profiles_v2

Lists all NIC Profiles with host NICs and capability.

## Example Usage

```hcl
data "nutanix_nic_profiles_v2" "example" {}
```

## Attribute Reference

The following attributes are exported:

- `nic_profiles`: List of NIC Profiles with host NICs and capability.

### nic_profiles

The `nic_profiles` object contains the following attributes:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `name`: Name of the NIC Profile.
- `description`: Description of the NIC Profile.
- `capability_config`: Capability configuration.
- `host_nic_references`: List of host NICs references associated with the NIC Profile.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `metadata`: Metadata associated with this resource.
- `nic_family`: Specification for a specific device family of Host NIC. The given string must be in the format vendor_id:device_id.
- `operating_mode`: Operating mode.
- `owner_type`: Owner type.
- `project_ext_id`: UUID of the project that owns this entity
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
