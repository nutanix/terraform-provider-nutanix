---
layout: "nutanix"
page_title: "NUTANIX: nutanix_claim_token_v2"
sidebar_current: "docs-nutanix-resource-claim-token-v2"
description: |-
  Creates a claim token with a specified name, expiry time, and maximum usage count.
---

# nutanix_claim_token_v2

Creates a claim token with a specified name, expiry time, and maximum usage count.

A claim token is a Foundation Central object used to register nodes. Create, update
and delete are asynchronous, task-based operations.

## Example

```hcl
resource "nutanix_claim_token_v2" "example" {
  name            = "example-claim-token"
  expiry_time     = "2030-01-01T00:00:00Z"
  max_usage_count = 5
}
```

## Argument Reference

The following arguments are supported:

* `name`: (Required) Name of the claim token.
* `expiry_time`: (Required) Expiry time of the claim token (RFC3339 format).
* `max_usage_count`: (Required) Maximum number of times the claim token can be used for registering nodes.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `created_time`: Time when the claim token was created.
* `current_usage_count`: Number of times the claim token has been used for registering nodes.
* `owner_ext_id`: External ID of the owner.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.

### links

* `href`: The URL at which the entity described by the link can be accessed.
* `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.

## Import

The `nutanix_claim_token_v2` resource can be imported using its external ID:

```
terraform import nutanix_claim_token_v2.example <claim_token_ext_id>
```

See detailed information in [Nutanix Create Claim Token V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
