---
layout: "nutanix"
page_title: "NUTANIX: nutanix_patched_image_v2"
sidebar_current: "docs-nutanix-resource-patched-image-v2"
description: |-
  Creates a patched hypervisor image from a base image, customized for specific node deployments.
---

# nutanix_patched_image_v2

Creates a patched hypervisor image from a base image, customized for specific node deployments.

A patched image references the claim token used during its creation via
`claim_token_ext_id`. Create and delete are asynchronous, task-based operations;
the image cannot be updated in place, so all input fields force recreation.

## Example

```hcl
resource "nutanix_claim_token_v2" "example" {
  name            = "example-patched-image-token"
  expiry_time     = "2030-01-01T00:00:00Z"
  max_usage_count = 5
}

resource "nutanix_patched_image_v2" "example" {
  claim_token_ext_id = nutanix_claim_token_v2.example.ext_id
  name               = "example-patched-image"
  host_type          = "AHV"

  image_details {
    local_host_image_ext_id = "00000000-0000-0000-0000-000000000000"
  }

  node_configurations {
    node_ext_id = "00000000-0000-0000-0000-000000000001"

    host_configuration {
      hostname = "patched-host-01"

      network_details {
        management {
          ip {
            ipv4 {
              value         = "10.0.0.10"
              prefix_length = 24
            }
          }
          gateway {
            ipv4 {
              value = "10.0.0.1"
            }
          }
          vlan_id = 0
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `claim_token_ext_id`: (Required) ExtId of the claim token used in the patched image.
* `name`: (Required) Patched image name.
* `host_type`: (Required) Type of the host installed on the node. One of `AHV`, `ESX`.
* `version`: (Optional) Version of the base hypervisor or host OS image used for patching.
* `image_details`: (Required) Details of the image used for patching.
* `node_configurations`: (Required) List of node configurations used for patching the image.

### image_details

* `local_host_image_ext_id`: (Required) External ID of the uploaded host image.

### node_configurations

* `node_ext_id`: (Required) External ID of the node.
* `host_configuration`: (Optional) Host configuration applied to the node.

### host_configuration

* `hostname`: (Optional) Hostname of the hypervisor or host OS.
* `nameservers`: (Optional) List of nameserver IP addresses to be used by the host.
* `ntpservers`: (Optional) List of NTP server IP addresses to be used by the host.
* `network_details`: (Optional) Network details for the host.

### network_details

* `management`: (Optional) Management network configuration (`ip`, `gateway`, `vlan_id`, `mtu_bytes`).

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `patched_iso_url`: URL for downloading the patched image.
* `patched_iso_sha256_checksum`: SHA-256 checksum of the patched image.
* `created_time`: Timestamp when the patched image was created.
* `owner_ext_id`: External ID of the owner.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.

## Import

The `nutanix_patched_image_v2` resource can be imported using its external ID:

```
terraform import nutanix_patched_image_v2.example <patched_image_ext_id>
```

See detailed information in [Nutanix Create Patched Image V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
