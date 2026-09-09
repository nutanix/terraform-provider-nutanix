---
layout: "nutanix"
page_title: "NUTANIX: nutanix_claim_tokens_v2"
sidebar_current: "docs-nutanix-datasource-claim-tokens-v2"
description: |-
  Returns a paginated list of all claim tokens.
---

# nutanix_claim_tokens_v2

Returns a paginated list of all claim tokens.

## Example

```hcl
data "nutanix_claim_tokens_v2" "example" {
  filter = "name eq 'example-claim-token'"
}
```

## Argument Reference

* `page`: (Optional) A URL query parameter that specifies the page number of the result set.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attributes Reference

* `claim_tokens`: List of claim tokens. Each element has the same attributes as the [nutanix_claim_token_v2](claim_token_v2.html) datasource.

See detailed information in [Nutanix List Claim Tokens V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
