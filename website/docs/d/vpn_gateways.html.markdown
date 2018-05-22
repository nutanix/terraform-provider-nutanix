---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_gateways"
sidebar_current: "docs-outscale-datasource-vpn-gateways"
description: |-
  Describes one or more VPN gateways.
---

# outscale_vpn_gateways

Describes one or more VPN gateways.

## Example Usage

```hcl
resource "outscale_vpn_gateway" "unattached" {
    tag {
		Name = "terraform-testacc-vpn-gateway-data-source-unattached-%d"
      	ABC  = "testacc-%d"
		XYZ  = "testacc-%d"
    }
}

data "outscale_vpn_gateways" "test_by_id" {
	vpn_gateway_id = ["${outscale_vpn_gateway.unattached.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `vpn_gateway_id`	One or more virtual private gateways.

See detailed information in [Outscale VPN Gateways](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described VPN Gateways on the following properties:

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

* `vpn_gateway_set` -	Information about one or more virtual private gateways.	false	VpnGateway
* `request_id`-	The ID of the request	false	

See detailed information in [Describe VPN Gateways](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpnGateways_get.html#_api_fcu-action_describevpngateways_get).