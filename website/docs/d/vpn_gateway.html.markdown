---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_gateway"
sidebar_current: "docs-outscale-datasource-vpn-gateway"
description: |-
  Describes one or more VPN gateway.
---

# outscale_volume

  Describes one or more VPN gateway.

## Example Usage

```hcl
resource "outscale_vpn_gateway" "unattached" {
    tag {
		Name = "terraform-testacc-vpn-gateway-data-source-unattached-%d"
      	ABC  = "testacc-%d"
		XYZ  = "testacc-%d"
    }
}

data "outscale_vpn_gateway" "test_by_id" {
	vpn_gateway_id = "${outscale_vpn_gateway.unattached.id}"
}
```

## Argument Reference

The following arguments are supported:

* `VpnGatewayId`	One or more virtual private gateways.	false	string

See detailed information in [Outscale VPN Gateway](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described VPN Gateway on the following properties:

* `attachment.state`: -	The current attachment state between the gateway and the VPC (attaching | attached | detaching | detached).	ReadVpnGateways	link.state
* `attachment.vpc-id`: -	The ID of the VPC the virtual private gateway is attached to.	ReadVpnGateways	link.lin-id
* `state`: -	The state of the virtual private gateway (pending | available | deleting | deleted).	ReadVpnGateways	state
* `tag`: -	The key/value combination of a tag associated with the resource.	ReadVpnGateways	tag
* `tag-key`: -	The key of a tag associated with the resource.	ReadVpnGateways	tag-key
* `tag-value`: -	The value of a tag associated with the resource.	ReadVpnGateways	tag-value
* `type`: -	The type of virtual private gateway (only ipsec.1 is supported).	ReadVpnGateways	type
* `vpn-gateway-id`: -	The ID of the virtual private gateway.	ReadVpnGateways	vpn-gateway-id


## Attributes Reference

The following attributes are exported:

* `attachments` -	The VPC to which the virtual private gateway is attached.	false	VpcAttachment
* `state`	- The state of the virtual private gateway (pending | available | deleting | deleted).	false	string
* `tag_set`	- One or more tags associated with the virtual private gateway.	false	Tag
* `type`	- The type of VPN connection supported by the virtual private gateway (only ipsec.1 is supported) .	false	string
* `vpn_gateway_id` -	The ID of the virtual private gateway.	false	string
* `request_id` -	The ID of the resquest	false	string

See detailed information in [Describe VPN Gateway](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpnGateways_get.html#_api_fcu-action_describevpngateways_get).