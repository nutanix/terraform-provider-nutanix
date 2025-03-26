---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_entity_v2"
sidebar_current: "docs-nutanix-datasource-lcm-entity-v2"
description: |-
  Get details about an LCM entity.
---

# nutanix_lcm_entity_v2
Get details about an LCM entity.

## Example

```hcl
data "nutanix_lcm_entity_v2" "entity-before-upgrade" {
  ext_id = "613no9d0-7caf-49y7-k582-1db5a5df580c"
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`:(Required) ExtId of the LCM entity.

## Attributes Reference
The fooling attributes are exported:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `entity_class`: LCM entity class.
* `entity_model`: LCM entity model.
* `entity_type`: Type of an LCM entity.  Enum Values:
    * `FIRMWARE`: LCM entity type firmware.
    * `SOFTWARE`: LCM entity type software.
* `entity_version`: Current version of an LCM entity.
* `hardware_family`: A hardware family for a LCM entity.
* `entity_description`: Description of an LCM entity.
* `location_info`: Location info corresponds to a tuple of location type (either node/cluster) and ExtID
* `target_version`: The requested update version of an LCM entity.
* `last_updated_time`: UTC date and time in RFC-3339 format when the task was last updated.
* `device_id`: Unique identifier of an LCM entity e.g. "HDD serial number".
* `group_uuid`: UUID of the group that this LCM entity is part of.
* `entity_details`: Detailed information for the LCM entity. For example, firmware entities contain additional information about NIC and so on.
* `child_entities`: Component information for the payload based entity.
* `available_versions`: List of available versions for an LCM entity to update.
* `sub_entities`: A list of sub-entities applicable to the entity.
* `cluster_ext_id`: Cluster uuid on which the resource is present or operation is being performed.
* `hardware_vendor`: Hardware vendor information.

### Location Info
The `location` attribute supports the following

* `uuid`: Location UUID of the resource.
* `location_type`: Scope of entity represented in LCM. This could be either Node or cluster type. Enum Values:
    * `PC`: Entity for which the scope is Prism Central wide.
    * `NODE`: Entity that belongs to a node in the cluster.
    * `CLUSTER`: Entity for which the scope is cluster wide.

### Entity Details
The `entity_details` attribute supports the following

* `name`: The key of the key-value pair.
* `value`: The value associated with the key for this key-value pair, string or integer or boolean or Array of strings or object or Array of MapOfString(objects) or Array of integers

### Available Versions
The `available_versions` attribute supports the following

* `version`: Version of the LCM entity.
* `status`: Status of the LCM entity. Enum Values:
    * `AVAILABLE`: Available version.
    * `EMERGENCY`: Emergency version.
    * `RECOMMENDED`: Deprecated version.
    * `STS`: Short-term supported version.
    * `LTS`: Long-term supported version.
    * `LATEST`: Latest version.
    * `DEPRECATED`: Deprecated version.
    * `ESTS`: Extended short-term supported version.
    * `CRITICAL`: Critical version.
* `is_enabled`: Indicates if the available update is enabled.
* `available_version_uuid`: Available version UUID.
* `order`: Order of the available version.
* `disablement_reason`: Reason for disablement of the available version.
* `release_notes`: Release notes for the available version.
* `release_date`: Release date of the available version.
* `custom_message`: Custom message associated with the available version.
* `dependencies`: List of dependencies for the available version.

#### Dependencies
The `dependencies` attribute supports the following

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `entity_class`: LCM entity class.
* `entity_model`: LCM entity model.
* `entity_type`: Type of an LCM entity.  Enum Values:
    * `FIRMWARE`: LCM entity type firmware.
    * `SOFTWARE`: LCM entity type software.
* `entity_version`: Current version of an LCM entity.
* `hardware_family`: A hardware family for a LCM entity.
* `dependent_versions`: Information of the dependent entity versions for this available entity.

#### Dependent Versions
The `dependent_versions` attribute supports the following

* `name`: The key of the key-value pair.
* `value`: The value associated with the key for this key-value pair, string or integer or boolean or Array of strings or object or Array of MapOfString(objects) or Array of integers

### Sub Entities
The `sub_entities` attribute supports the following

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `entity_class`: LCM entity class.
* `entity_model`: LCM entity model.
* `entity_type`: Type of an LCM entity.  Enum Values:
    * `FIRMWARE`: LCM entity type firmware.
    * `SOFTWARE`: LCM entity type software.
* `entity_version`: Current version of an LCM entity.
* `hardware_family`: A hardware family for a LCM entity.

See detailed information in [Nutanix LCM Entity V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Entities/operation/getEntityById).
