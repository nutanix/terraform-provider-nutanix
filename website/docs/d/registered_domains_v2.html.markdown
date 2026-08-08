---
layout: "nutanix"
page_title: "NUTANIX: nutanix_registered_domains_v2"
sidebar_current: "docs-nutanix-datasource-registered-domains-v2"
description: |-
  List all registered domains.
---

# nutanix_registered_domains_v2

List all registered domains.

## Example Usage

```hcl
data "nutanix_registered_domains_v2" "list" {}

data "nutanix_registered_domains_v2" "filtered" {
  limit  = 5
  filter = "name eq 'example-registered-domain'"
}
```

## Argument Reference

The following arguments are supported:

* `page`:- (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit`:- (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter`:- (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`:- (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `expand`:- (Optional) A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved.
* `select`:- (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.

## Attributes Reference

The following attributes are exported:

* `registered_domains`:- List of registered domains. Each element supports the same attributes as the [nutanix_registered_domain_v2](registered_domain_v2.html) datasource.

See detailed information in [Nutanix Registered Domains V2](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.0).
