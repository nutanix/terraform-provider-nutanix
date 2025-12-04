---
layout: "nutanix"
page_title: "NUTANIX: nutanix_service_group_v2"
sidebar_current: "docs-nutanix-datasource-service-group-v2"
description: |-
  This operation retrieves an service_group.
---

# nutanix_service_group_v2

Get an service Group by ExtID

## Example Usage

```hcl
data "nutanix_service_group_v2" "service_group" {
  ext_id = "07167778-266d-4052-9992-f30cbfd52e83"
}
```


## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) service group UUID.

## Attribute Reference

The following attributes are exported:

* `name`: A short identifier for a Service Group.
* `description`: A user defined annotation for a Service Group.
* `is_system_defined`: Service Group is system defined or not.
* `tcp_services`: List of TCP ports in the service.
* `udp_services`: List of UDP ports in the service.
* `icmp_services`: Icmp Type Code List.
* `policy_references`: Reference to policy associated with Service Group.
* `created_by`: created by.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.


### tcp_services, udp_services
* `start_port`: start port
* `end_port`: end port

### icmp_services
* `type`: Icmp service Type. Ignore this field if Type has to be ANY.
* `code`: Icmp service Code. Ignore this field if Code has to be ANY
* `is_all_allowed`: Set this field to true if both Type and Code is ANY.

See detailed information in [Nutanix Get Service Group v4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.0#tag/ServiceGroups/operation/getServiceGroupById).
