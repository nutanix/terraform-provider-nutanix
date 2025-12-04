---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_entities_v2"
sidebar_current: "docs-nutanix-datasource-lcm-entities-v2"
description: |-
  Get details about all LCM entities.

---

# nutanix_lcm_entities_v2
Get details about all LCM entities.


## Example

```hcl
data "nutanix_lcm_entity_v2" "entities" {}

data "nutanix_lcm_entities_v2" "lcm-entities-filtered" {
  filter = "entityModel eq 'Calm Policy Engine'"
}

data "nutanix_lcm_entities_v2" "lcm-entities-limit" {
  limit = 5
}

```

## Argument Reference

The following arguments are supported:

* `page`: - A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` :A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields: clientName, createdBy, extId, createdTime, displayName, extId, isSystemDefined, lastUpdatedTime.
    * The filter can be applied to the following fields:
        * `clusterExtId`
        * `entityClass`
        * `entityModel`
        * `entityType`
        * `entityVersion`
        * `hardwareVendor`
* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields: createdTime, distinguishedName, displayName, extId, lastUpdatedTime.
    * The orderby can be applied to the following fields:
        * `entityClass`
        * `entityModel`
        * `entityType`
        * `entityVersion`
        * `hardwareVendor`
* `select` : A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. following fields: accessibleClients, accessibleEntityTypes, assignedUserGroupsCount, assignedUsersCount, clientName, createdBy, createdTime, description, displayName, extId, isSystemDefined, lastUpdatedTime, links, operations, tenantId.
    * The select can be applied to the following fields:
        * `entityClass`
        * `entityModel`
        * `entityType`
        * `entityVersion`
        * `hardwareVendor`

## Attributes Reference
The following attributes are exported:

* `entities`: List of LCM entities.

### Entities
The `entities` attribute supports the following:

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

See detailed information in [Nutanix LCM Entities V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Entities/operation/listEntities).
