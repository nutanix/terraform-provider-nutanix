---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_rate_limit_policy_v2"
sidebar_current: "docs-nutanix-resource-image-rate-limit-policy-v2"
description: |-
  Creates an image rate limit policy using the provided request body. The name, rate limit Kbps and cluster entity filter are mandatory fields for creating an image rate limit.
---

# nutanix_image_rate_limit_policy_v2

Creates an image rate limit policy using the provided request body. The name, rate limit Kbps and cluster entity filter are mandatory fields for creating an image rate limit.

```hcl
resource "nutanix_image_rate_limit_policy_v2" "example" {
  name            = "example-rate-limit-policy"
  description     = "Example image rate limit policy"
  rate_limit_kbps = 1000

  cluster_entity_filter {
    category_ext_ids = ["ab520e1d-4950-1db1-917f-a9e2ea35b8e3"]
    type             = "CATEGORIES_MATCH_ALL"
  }
  # type             = "CATEGORIES_MATCH_ALL"
  # Cluster matching with all the categories specified in the cluster_entity_filter --> category_ext_ids, will be enrolled for the rate limit policy.

  # type             = "CATEGORIES_MATCH_ANY"
  # Cluster matching with any of the categories specified in the cluster_entity_filter --> category_ext_ids, will be enrolled for the rate limit policy.
}
```

## Argument Reference

The following arguments are supported:

* `name`: (Required) Name of the image rate limit policy.
* `description`: (Optional) Image rate limit policy specification.
* `rate_limit_kbps`: (Required) Network bandwidth in KBps that the rate limited image operation can utilize.
* `cluster_entity_filter`: (Required) Category-based entity filter for clusters.

### cluster_entity_filter

* `type`: (Required) Filter matching type. Valid values: "CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY".
* `category_ext_ids`: (Required) Filter matches entities that have these categories attached.

## Attributes Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of the image rate limit policy.
* `matching_cluster_ext_ids`: External identifier of the Prism Elements where a rate limit is the effective rate limit policy.
* `owner_ext_id`: External identifier of the owner of the rate limit policy.
* `owner_name`: Name of the owner of the rate limit policy.
* `create_time`: Image rate limit policy creation time.
* `last_update_time`: Last updated time of an image rate limit policy.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

## Import

Image Rate Limit Policies can be imported using the `UUID` (ext_id in v4 API context). For example:

```hcl
resource "nutanix_image_rate_limit_policy_v2" "import_policy" {}
```

```shell
terraform import nutanix_image_rate_limit_policy_v2.import_policy <UUID>
```

See detailed information in [Nutanix Image Rate Limit Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/ImageRateLimitPolicies/operation/createRateLimitPolicy)
