---
layout: "nutanix"
page_title: "NUTANIX: nutanix_effective_image_rate_limit_policies_v2"
sidebar_current: "docs-nutanix-datasource-effective-image-rate-limit-policies-v2"
description: |-
  The effective rate limit for the Prism Elements. If no rate limit applies to a particular cluster, no entry is returned for that cluster. The API supports operations such as filtering, sorting, selection, and pagination.
---

# nutanix_effective_image_rate_limit_policies_v2

The effective rate limit for the Prism Elements. If no rate limit applies to a particular cluster, no entry is returned for that cluster. The API supports operations such as filtering, sorting, selection, and pagination.

```hcl
data "nutanix_effective_image_rate_limit_policies_v2" "example" {}
```

## Argument Reference

The following arguments are supported:

* `page`: (Optional) A URL query parameter that specifies the page number of the result set.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attributes Reference

The following attributes are exported:

* `effective_rate_limit_policies`: List of effective image rate limit policies.

### effective_rate_limit_policies

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `cluster_ext_id`: Cluster external identifier.
* `rate_limit_ext_id`: The external identifier of image rate limit policy.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Image Rate Limit Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/ImageRateLimitPolicies/operation/listEffectiveRateLimitPolicies)
