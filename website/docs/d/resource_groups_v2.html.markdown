---
layout: "nutanix"
page_title: "NUTANIX: nutanix_resource_groups_v2"
sidebar_current: "docs-nutanix-datasource-resource-groups-v2"
description: |-
  List the resource groups defined on the system.
---

# nutanix_resource_groups_v2

List the resource groups defined on the system.

## Example Usage

```hcl
data "nutanix_resource_groups_v2" "example" {}
```

## Argument Reference
The following arguments are supported:

- `page`:- A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`:- A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`:- A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - createTime
  - createdBy
  - lastUpdateTime
  - lastUpdatedBy
  - name
  - placementTargets/clusterExtId
  - placementTargets/storageContainers/extId
  - projectExtId
- `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - createTime
  - createdBy
  - extId
  - lastUpdateTime
  - lastUpdatedBy
  - name
- `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields:
  - name
  - projectExtId

## Attributes Reference

The following attributes are exported:

* `resource_groups`:- List of resource groups.

## Resource Groups

The `resource_groups` attribute is a list of resource group objects. Each resource group supports the following attributes:

* `ext_id`:- A globally unique identifier of the resource group.
* `name`:- Name of the resource group.
* `project_ext_id`:- External identifier of the project.
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity.
* `created_by`:- User who created the resource group.
* `last_updated_by`:- User who last updated the resource group.
* `create_time`:- Creation time (RFC3339).
* `last_update_time`:- Last update time (RFC3339).
* `placement_targets`:- List of placement targets.
* `links`:- A HATEOAS style link for the response.

## Placement Targets

The `placement_targets` attribute supports the following:

* `cluster_ext_id`:- UUID of the AOS cluster.
* `storage_containers`:- List of storage containers available for this cluster target.

## Storage Containers

The `storage_containers` attribute supports the following:
* `ext_id`:- UUID of the storage container.

## Links

The `links` attribute supports the following:

* `href`:- The URL at which the entity described by the link can be accessed.
* `rel`:- A name that identifies the relationship of the link to the object.
