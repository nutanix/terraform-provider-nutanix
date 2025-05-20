---
layout: "nutanix"
page_title: "NUTANIX: nutanix_service_groups_v2"
sidebar_current: "docs-nutanix-resource-service-groups-v2"
description: |-
  This operation submits a request to create a service group based on the input parameters.
---

# nutanix_service_group

Create an service Group

## Example Usage

```hcl

# Add Service  group. with TCP and UDP
resource "nutanix_service_groups_v2" "tcp-udp-service" {
  name        = "service_group_tcp_udp"
  description = "service group description"
  tcp_services {
    start_port = "232"
    end_port   = "232"
  }
  udp_services {
    start_port = "232"
    end_port   = "232"
  }
}

# service group with ICMP
resource "nutanix_service_groups_v2" "icmp-service" {
  name        = "service_group_icmp"
  description = "service group description"
  icmp_services {
    type = 8
    code = 0
  }
}

# service group with All TCP, UDP and ICMP
resource "nutanix_service_groups_v2" "all-service" {
  name        = "service_group_udp_tcp_icmp"
  description = "service group description"
  tcp_services {
    start_port = "232"
    end_port   = "232"
  }
  udp_services {
    start_port = "232"
    end_port   = "232"
  }
  icmp_services {
    type = 8
    code = 0
  }
}

```


## Argument Reference

The following arguments are supported:

* `name`: (Required) Name of the service group
* `description`: (Optional) Description of the service group
* `tcp_services`: (Optional) List of TCP ports in the service.
* `udp_services`: (Optional) List of UDP ports in the service.
* `icmp_services`: (Optional) Icmp Type Code List.


### tcp_services, udp_services
* `start_port`: (Required) start port
* `end_port`: (Required) end port

### icmp_services
* `type`: (Optional) Icmp service Type. Ignore this field if Type has to be ANY.
* `code`: (Optional) Icmp service Code. Ignore this field if Code has to be ANY
* `is_all_allowed`: (Optional) Set this field to true if both Type and Code is ANY. Default is False.


## Attributes Reference

The following attributes are exported:

* `ext_id`: address group uuid.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `policy_references`: Reference to policy associated with Address Group.
* `created_by`: created by.
* `is_system_defined`: Service Group is system defined or not.


See detailed information in [Nutanix Service Groups V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.0#tag/ServiceGroups/operation/createServiceGroup).
