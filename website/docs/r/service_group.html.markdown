---
layout: "nutanix"
page_title: "NUTANIX: nutanix_service_group"
sidebar_current: "docs-nutanix-resource-service-group"
description: |-
  This operation submits a request to create a service group based on the input parameters.
---

# nutanix_service_group

Provides a resource to create a service group based on the input parameters.

## Example Usage

```hcl
resource "nutanix_service_group" "test" {
		name = "test_service_gp"
		description = "this is service group"

		service_list {
			protocol = "TCP"
			tcp_port_range_list {
				start_port = 22
				end_port = 22
			}
			tcp_port_range_list {
				start_port = 2222
				end_port = 2222
			}
		}
	}
```


## Argument Reference

The following arguments are supported:

* `name`: - (Required) Name of the service group
* `description`: - (Optional) Description of the service group
* `service_list`: - (Required) list of services which have protocol (TCP / UDP / ICMP) along with port details
* `system_defined`: - (ReadOnly) boolean value to denote if the service group is system defined

### Service List

The service_list argument supports the following:

* `protocol`: - (Optional) The UserPrincipalName of the user from the directory service.
* `icmp_type_code_list`: - (Optional) ICMP type code list
* `tcp_port_range_list`: - (Optional) TCP Port range list
* `udp_port_range_list`: - (Optional) UDP port range list 

#### ICMP Port range list

The icmp_type_code_list argument supports the following:

* `code`: - (Optional) Code as text
* `type`: - (Optional) Type as text

#### TCP Port Range

The tcp_port_range_list attribute supports the following:

* `start_port`: - (Optional) Start Port (Int)
* `end_port` - (Optional) End Port (Int)

#### UDP Port Range

The udp_port_range_list attribute supports the following:

* `start_port`: - (Optional) Start Port (Int)
* `end_port` - (Optional) End Port (Int)

See detailed information in [Nutanix Service Groups](https://www.nutanix.dev/api_references/prism-central-v3/#/38492a5cb53e2-create-a-new-service-group).
