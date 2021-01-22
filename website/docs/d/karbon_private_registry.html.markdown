---
layout: "nutanix"
page_title: "NUTANIX: nutanix_karbon_private_registry"
sidebar_current: "docs-nutanix-datasource-karbon-private-registry"
description: |-
 Describes a Karbon private registry entry
---

# nutanix_karbon_private_registry

Describes Karbon private registry entry

## Example Usage

```hcl
data "nutanix_karbon_private_registry" "registry" {
   cluster_id = "<YOUR-CLUSTER-ID>"
}
```

## Argument Reference

The following arguments are supported:

* `private_registry_id`: Represents karbon private registry uuid
* `private_registry_name`: Represents the name of karbon private registry

## Attributes Reference

The following attributes are supported:

* `name`: - Name of the private registry.
* `uuid`: - UUID of the private registry.
* `endpoint`: - Endpoint of the private in format `url:port`.


See detailed information in [Nutanix Karbon Registry](https://www.nutanix.dev/reference/karbon/api-reference/registry/getprivateregistry/).
