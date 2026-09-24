---
layout: "nutanix"
page_title: "NUTANIX: nutanix_nic_profile_v2"
sidebar_current: "docs-nutanix-resource-nic-profile-v2"
description: |-
  Create a NIC Profile with the given name, description and capability specifications.
---

# nutanix_nic_profile_v2

Create a NIC Profile with the given name, description and capability specifications.

## Example Usage

```hcl
resource "nutanix_nic_profile_v2" "example" {
  name        = "tf-nic-profile"
  description = "Description of the NIC Profile."
  nic_family  = "15b3:101d"
  host_nic_ext_ids = [
    "11111111-1111-1111-1111-111111111111",
  ]

  capability_config {
    capability_type = "SRIOV"
  }

  operating_mode = "ETHERNET"
}
```

## Argument Reference

The following arguments are supported:

- `name`: (Required) Name of the NIC Profile.
- `description`: (Optional) Description of the NIC Profile.
- `capability_config`: (Required) Capability configuration.
- `nic_family`: (Required) Specification for a specific device family of Host NIC. The given string must be in the format vendor_id:device_id.
- `host_nic_ext_ids`: (Optional) Set of Host NIC UUIDs that should be associated with the NIC Profile. Updating this set will associate newly added Host NICs and disassociate removed Host NICs.
- `operating_mode`: (Optional) Operating mode.
- `owner_type`: (Optional) Owner type.
- `project_ext_id`: (Optional) UUID of the project that owns this entity

### capability_config

The `capability_config` object contains the following arguments:

- `capability_type`: (Required) Capability type.

## Attribute Reference

The following attributes are exported:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `host_nic_ext_ids`: Set of Host NIC UUIDs currently associated with the NIC Profile.
- `host_nic_references`: List of host NICs references associated with the NIC Profile.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `metadata`: Metadata associated with this resource.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).

### host_nic_references

The `host_nic_references` object contains the following attributes:

- `associated_vm_nic_references`: List of VM NICs references associated with the Host Nic.
- `compliance_status`: Host NIC compliance status.
- `ext_id`: UUID of the Host Nic.
- `num_vfs`: Number of VFs associated with the Host Nic.

### links

The `links` object contains the following attributes:

- `href`: The URL at which the entity described by the link can be accessed.
- `rel`: A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### metadata

The `metadata` object contains the following attributes:

- `owner_reference_id`: A globally unique identifier that represents the owner of this resource.
- `owner_user_name`: The userName of the owner of this resource.
- `project_reference_id`: A globally unique identifier that represents the project this resource belongs to.
- `project_name`: The name of the project this resource belongs to.
- `category_ids`: A list of globally unique identifiers that represent all the categories the resource is associated with.

## Import

NIC Profile can be imported using the NIC Profile UUID:

```hcl
terraform import nutanix_nic_profile_v2.example <nic-profile-uuid>
```
