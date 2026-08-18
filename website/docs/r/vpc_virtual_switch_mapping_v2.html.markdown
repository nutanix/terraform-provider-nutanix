---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vpc_virtual_switch_mapping_v2"
sidebar_current: "docs-nutanix-resource-vpc-virtual-switch-mapping-v2"
description: |-
  Create and manage VPC to Virtual Switch mappings and traffic configurations in the Nutanix v4 API.
---

# nutanix_vpc_virtual_switch_mapping_v2

The `nutanix_vpc_virtual_switch_mapping_v2` resource allows you to configure and manage the mapping between a Virtual Private Cloud (VPC) and a Virtual Switch within the Nutanix infrastructure. 

This resource is primarily used to control East-West traffic configurations at the Virtual Switch level. It allows administrators to dictate whether all network traffic is permitted through the switch or if the flow should be restricted, alongside assigning appropriate cluster boundaries and metadata (such as project ownership and categories).

## Example Usage

```hcl
resource "nutanix_vpc_virtual_switch_mapping_v2" "example" {
  mappings {
    virtual_switch_uuid      = "1ce108f8-bd48-4287-aead-ca0c3276fc05"
    cluster_uuids            = ["00065264-61ec-9c95-185b-ac1f6b6f97e2"]
    is_all_traffic_permitted = true

    metadata {
      owner_reference_id   = "44444444-4444-4444-4444-444444444444"
      owner_user_name      = "admin_user"
      project_reference_id = "55555555-5555-5555-5555-555555555555"
      project_name         = "production-project"
      category_ids = [
        "9ab5eda1-da67-4058-6c9e-66bd466a585d",
        "77d123d8-c5f5-4eda-596d-ea1440ba77bb",
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `mappings` - (Required, ForceNew) A block that defines the specific configuration for mapping a VPC to a virtual switch. This block encapsulates the target switch, traffic policies, and organizational metadata.

### mappings

The `mappings` block supports the following configuration arguments:

* `virtual_switch_uuid` - (Required, ForceNew) The Universally Unique Identifier (UUID) of the target Virtual Switch. This is the switch where the VPC traffic policies will be applied.
* `cluster_uuids` - (Optional, ForceNew) A list of cluster UUIDs where this Virtual Switch mapping is effective. This defines the physical cluster boundaries for the networking configuration.
* `is_all_traffic_permitted` - (Optional, ForceNew) A boolean flag indicating the traffic permission policy. When set to `true`, all network traffic is permitted to flow through the virtual switch. When set to `false`, the switch drops standard traffic and only permits basic ICMP and essential statistics collection requests.
* `metadata` - (Optional, ForceNew) A block containing administrative and organizational metadata for the mapping. See the [`metadata`](#metadata) block below for details.

### metadata

The `metadata` block allows you to attach logical grouping and RBAC (Role-Based Access Control) properties to the mapping. It supports the following arguments:

* `owner_reference_id` - (Optional, ForceNew) A globally unique identifier (UUID) that represents the user, group, or service account that owns this mapping resource.
* `owner_user_name` - (Optional, ForceNew) The human-readable username of the resource owner.
* `project_reference_id` - (Optional, ForceNew) A globally unique identifier (UUID) of the Nutanix Project this resource belongs to. This is heavily used for multi-tenancy and resource isolation.
* `project_name` - (Optional, ForceNew) The human-readable name of the Nutanix Project this resource is assigned to.
* `category_ids` - (Optional, ForceNew) A list of globally unique identifiers (UUIDs) representing the categories (tags) associated with the mapping. Categories are useful for policy enforcement, granular filtering, or billing attribution.

See detailed information in [Nutanix VPC Virtual Switch Mapping V4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.3#tag/VpcVirtualSwitchMappings/operation/createVpcVirtualSwitchMapping).