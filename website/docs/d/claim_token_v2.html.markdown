---
layout: "nutanix"
page_title: "NUTANIX: nutanix_claim_token_v2"
sidebar_current: "docs-nutanix-datasource-claim-token-v2"
description: |-
  Returns details of a claim token identified by its external ID.
---

# nutanix_claim_token_v2

Returns details of a claim token identified by its external ID.

## Example

```hcl
data "nutanix_claim_token_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id`: (Required) A globally unique identifier of an instance that is suitable for external consumption.

## Attributes Reference

* `name`: Name of the claim token.
* `expiry_time`: Expiry time of the claim token.
* `max_usage_count`: Maximum number of times the claim token can be used for registering nodes.
* `created_time`: Time when the claim token was created.
* `current_usage_count`: Number of times the claim token has been used for registering nodes.
* `owner_ext_id`: External ID of the owner.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.

See detailed information in [Nutanix Get Claim Token V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
