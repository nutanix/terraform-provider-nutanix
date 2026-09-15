---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_rate_limit_policies_v2"
sidebar_current: "docs-nutanix-datasource-image-rate-limit-policies-v2"
description: |-
  Lists image rate limit policies created on Prism Central along with other details such as, name, description and so on. This API supports operations such as filtering, sorting, selection, and pagination.
---

# nutanix_image_rate_limit_policies_v2

Lists image rate limit policies created on Prism Central along with other details such as, name, description and so on. This API supports operations such as filtering, sorting, selection, and pagination.

```hcl
data "nutanix_image_rate_limit_policies_v2" "example" {}
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

* `rate_limit_policies`: List of image rate limit policies.

### rate_limit_policies

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `name`: Name of the image rate limit policy.
* `description`: Image rate limit policy specification.
* `rate_limit_kbps`: Network bandwidth in KBps that the rate limited image operation can utilize.
* `cluster_entity_filter`: Category-based entity filter for clusters.
* `matching_cluster_ext_ids`: External identifier of the Prism Elements where a rate limit is the effective rate limit policy.
* `owner_ext_id`: External identifier of the owner of the rate limit policy.
* `owner_name`: Name of the owner of the rate limit policy.
* `create_time`: Image rate limit policy creation time.
* `last_update_time`: Last updated time of an image rate limit policy.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

### cluster_entity_filter

* `type`: Filter matching type.
* `category_ext_ids`: Filter matches entities that have these categories attached.

See detailed information in [Nutanix Image Rate Limit Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/ImageRateLimitPolicies/operation/listRateLimitPolicies)
