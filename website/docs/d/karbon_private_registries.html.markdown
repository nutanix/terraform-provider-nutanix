---
layout: "nutanix"
page_title: "NUTANIX: nutanix_karbon_private_registries"
sidebar_current: "docs-nutanix-datasource-karbon-private-registry"
description: |-
 Describes a List of Karbon private registry entry
---

# nutanix_karbon_private_registries

Describes a List of Karbon private registry entry

## Example Usage

```hcl
data "nutanix_karbon_private_registries" "registry" {
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


See detailed information in [Nutanix Karbon Registries](https://www.nutanix.dev/api_references/nke/#/6542bb676c318-list-the-private-registry-configurations-api-format-https-server-nutanix-com-9440-karbon-v1-alpha-1-registries). 