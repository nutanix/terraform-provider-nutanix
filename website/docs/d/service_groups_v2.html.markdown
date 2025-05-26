---
layout: "nutanix"
page_title: "NUTANIX: nutanix_service_groups_v2"
sidebar_current: "docs-nutanix-datasource-service-groups-v2"
description: |-
  This operation retrieves the list of service_groups.
---

# nutanix_service_groups_v2

List all the service Groups.

## Example Usage

```hcl

data "nutanix_service_groups_v2" "service_group"{}

data "nutanix_service_groups_v2" "service_group_filtered"{
    filter = "name eq 'service_group_name'"
}

```


## Argument Reference

The following arguments are supported:

* `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources. The filter can be applied to the following fields:
    - `createdBy`
    - `description`
    - `extId`
    - `isSystemDefined`
    - `name`
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
    - `description`
    - `extId`
    - `name`
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. The select can be applied to the following fields:
    - `createdBy`
    - `description`
    - `extId`
    - `icmpServices`
    - `isSystemDefined`
    - `links`
    - `name`
    - `policyReferences`
    - `tcpServices`
    - `udpServices`
    - `tenantId`
    - `udpServices`


## Attribute Reference

The following attributes are exported:

* `service_groups`: List of service groups

### Service Groups
The `service_groups` object contains the following attributes:

* `ext_id`: service group UUID.
* `name`: A short identifier for an service Group.
* `description`: A user defined annotation for an service Group.
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




See detailed information in [Nutanix List Service Groups v4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.0#tag/ServiceGroups/operation/listServiceGroups).
