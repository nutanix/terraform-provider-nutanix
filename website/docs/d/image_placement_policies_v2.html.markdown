---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_placements_v2"
sidebar_current: "docs-nutanix-datasource-image-placements-v2"
description: |-
 Describes a Image placement policies
---

# nutanix_image_v2

List image placement policies details.

## Example

```hcl
# List all image placement policies
data "nutanix_image_placement_policies_v2" "list-ipp"{}

# List image placement policies with filter, page and limit
data "nutanix_image_placement_policies_v2" "filtered-ipp"{
    filter="startswith(name,'ipp_name')"
    page=0
    limit=10
}

```


## Argument Reference

The following arguments are supported:

* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`:A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions. For example, filter '$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
    - description
    - enforcementState
    - name
* `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
    - description
    - enforcementState
    - name
* `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields:
    - createTime
    - description
    - enforcementState
    - extId
    - lastUpdateTime
    - links
    - name
    - ownerExtId
    - placementType
    - tenantId


## Attribute Reference

* `placement_policies`: List of all image placement policies

### placement_policies

The `placement_policies` object is a list of image placement policies. each image placement policy object contains the following attributes:

* `ext_id`: The external identifier of an image placement policy.
* `name`: (Required) Name of the image placement policy.
* `description`: (Optional) Description of the image placement policy.
* `placement_type`: (Required) Type of the image placement policy. Valid values:
    - HARD: Hard placement policy. Images can only be placed on clusters enforced by the image placement policy.
    - SOFT: Soft placement policy. Images can be placed on clusters apart from those enforced by the image placement policy.
* `image_entity_filter`: (Required) Category-based entity filter.
* `cluster_entity_filter`: (Required) Category-based entity filter.
* `enforcement_state`: (Optional) Enforcement status of the image placement policy. Valid values:
    - ACTIVE: The image placement policy is being actively enforced.
    - SUSPENDED: The policy enforcement for image placement is suspended.

### image_entity_filter
* `type`: (Required) Filter matching type. Valid values:
    - CATEGORIES_MATCH_ALL: Image policy only applies to the entities that are matched to all the corresponding entity categories attached to the image policy.
    - CATEGORIES_MATCH_ANY: Image policy applies to the entities that match any subset of the entity categories attached to the image policy.
* `category_ext_ids`: Array of strings

### cluster_entity_filter
* `type`: (Required) Filter matching type. Valid values:
    - CATEGORIES_MATCH_ALL: Image policy only applies to the entities that are matched to all the corresponding entity categories attached to the image policy.
    - CATEGORIES_MATCH_ANY: Image policy applies to the entities that match any subset of the entity categories attached to the image policy.
* `category_ext_ids`: Array of strings

See detailed information in [Nutanix List Image placement policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/ImagePlacementPolicies/operation/listPlacementPolicies)
