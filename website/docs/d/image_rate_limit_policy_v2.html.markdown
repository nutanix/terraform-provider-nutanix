---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_rate_limit_policy_v2"
sidebar_current: "docs-nutanix-datasource-image-rate-limit-policy-v2"
description: |-
  Retrieves an image rate limit policy details for the provided external identifier.
---

# nutanix_image_rate_limit_policy_v2

Retrieves an image rate limit policy details for the provided external identifier.

```hcl
data "nutanix_image_rate_limit_policy_v2" "example" {
  ext_id = "<IMAGE_RATE_LIMIT_POLICY_UUID>"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) The external identifier of image rate limit policy.

## Attributes Reference

The following attributes are exported:

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

See detailed information in [Nutanix Image Rate Limit Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/ImageRateLimitPolicies/operation/getRateLimitPolicyById)
