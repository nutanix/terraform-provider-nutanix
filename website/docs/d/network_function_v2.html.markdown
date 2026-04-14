---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_function_v2"
sidebar_current: "docs-nutanix-datasource-network-function-v2"
description: |-
  Provides a datasource to get a single Network Function corresponding to the ext_id.
---

# nutanix_network_function_v2

Get a single Network Function corresponding to the ext_id.

## Example Usage

```hcl
data "nutanix_network_function_v2" "nf" {
  ext_id = "52a4db2a-78a9-4c21-8e51-6c26a6ff92a9"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: (Required) Network Function UUID

## Attribute Reference

The following attributes are exported:

- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `ext_id`:  globally unique identifier of an instance that is suitable for external consumption.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `metadata`: Metadata associated with this resource.
- `name`: Name of the network function.
- `description`: Description of the network function.
- `failure_handling`: Failure handling behavior when network function is unhealthy. Values:

   | Value | Description |
   | --- | --- |
   | `NO_ACTION` | When network function is unhealthy, no action is taken and traffic is black-holed. This value is deprecated. If it continues to be used, it will automatically be converted to FAIL_CLOSE. |
   | `FAIL_CLOSE` | When all the network function VM(s) are down, all traffic from sources is blocked to prevent it from bypassing the security. |
   | `FAIL_OPEN` | When all the network function VM(s) are down, traffic from sources can be forwarded directly to the destinations, effectively bypassing the network function VM. |

- `high_availability_mode`: High availability configuration used between virtual NIC pairs. Values:
  
   | Value | Description |
   | --- | --- |
   | `ACTIVE_PASSIVE` | NIC pair is in ACTIVE_PASSIVE mode. In ACTIVE_PASSIVE mode, one of the NIC pairs will be selected as the ACTIVE network function and all other NIC pairs will be on STANDBY |

- `traffic_forwarding_mode`: Traffic forwarding mode. Values:
  
   | Value | Description |
   | --- | --- |
   | `INLINE` | Inline traffic redirection is applied through the network function VM to enable comprehensive inspection and policy enforcement. |
   | `VTAP` | Traffic is passively mirrored to the network function VM for out-of-band monitoring, without affecting the original traffic flow. The failureHandling or dataPlaneHealthCheckConfig or egressNicReference inside any of NicPair is not supported along with this mode. API will fail as part of validation if passed with VTAP trafficForwardingMode. |

- `data_plane_health_check_config`: Data plane health check configuration.
- `nic_pairs`: List of NIC pairs part of this network function.

### Links

The `links` attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Metadata

The `metadata` object contains the following attributes:

- `owner_reference_id` : A globally unique identifier that represents the owner of this resource.
- `owner_user_name` : The userName of the owner of this resource.
- `project_reference_id` : A globally unique identifier that represents the project this resource belongs to.
- `project_name` : The name of the project this resource belongs to.
- `category_ids` : A list of globally unique identifiers that represent all the categories the resource is associated with.

### Nic Pairs

The `nic_pairs` object contains the following attributes:

- `ingress_nic_reference` : UUID of NIC which will be used as ingress NIC..
- `egress_nic_reference` : UUID of NIC which will be used as egress NIC.
- `is_enabled` :  `Default: true`. Administrative state of the NIC pair. If it's set to False, the NIC pair will not be selected as ACTIVE network function.
- `vm_reference` : VM UUID which both ingress/egress NICs are part of.
- `data_plane_health_status` : Data plane health status of the NIC pair. Values:
  
   | Value | Description |
   | --- | --- |
   | `HEALTHY` | Entity is healthy. |
   | `UNHEALTHY` | Entity is unhealthy. |

- `high_availability_state` : High availability state of the NIC pair. Values:
  
   | Value | Description |
   | --- | --- |
   | `ACTIVE` | NIC pair is in ACTIVE mode. |
   | `PASSIVE` | NIC Pair is in standby mode.|

### Data Plane Health Check Config

The `data_plane_health_check_config` object contains the following attributes:

- `failure_threshold` :  `Default: 3`. The number of failure checks after which the target is considered unhealthy.
- `interval_secs` :  `Default: 5`. Interval in seconds between health checks.
- `success_threshold` :  `Default: 3`. The number of successful checks after which the target is considered healthy.
- `timeout_secs` :  `Default: 1`. The time, in seconds, after which a health check times out.

See detailed information in [Nutanix Network Function v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.3#tag/NetworkFunctions/operation/getNetworkFunctionById).
