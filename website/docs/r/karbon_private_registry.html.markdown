---
layout: "nutanix"
page_title: "NUTANIX: nutanix_karbon_private_registry"
sidebar_current: "docs-nutanix-resource-karbon-registry"
description: |-
  Provides a Nutanix Karbon Registry resource to Create a private registry entry in Karbon.
---

# nutanix_karbon_private_registry

Provides a Nutanix Karbon Registry resource to Create a private registry entry in Karbon.

## Example Usage

```hcl
data "nutanix_karbon_private_registry" "registries" {}

resource "nutanix_karbon_private_registry" "registry" {
}

```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) Name of the private registry configuration. **Note:** Updates to this attribute forces new resource creation.
* `cert`: - (Optional) Certificate of the private registry in format of base64-encoded byte array. **Note:** Updates to this attribute forces new resource creation.
* `url`: - (Optional) URL of the private registry. **Note:** Updates to this attribute forces new resource creation.
* `port`: - (Optional) Port of the private registry.
* `username`: - (Optional) Username for authentication to the private registry. **Note:** Updates to this attribute forces new resource creation.
* `password`: - (Optional) Password for authentication to the private registry. **Note:** Updates to this attribute forces new resource creation.


## Attributes Reference

The following attributes are exported:

* `endpoint`: - Endpoint of the private in format `url:port`.


See detailed information in [Nutanix Karbon Registry](https://www.nutanix.dev/api_references/nke/#/1ac891e63bef8-create-the-private-registry-entry-in-nke-with-the-provided-configuration-api-format-https-server-nutanix-com-9440-karbon-v1-alpha-1-registries).
