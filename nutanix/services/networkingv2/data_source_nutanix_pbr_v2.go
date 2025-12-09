package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixPbrV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixPbrV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV2(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_match": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"subnet_prefix": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": SchemaForValuePrefixLength(),
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
															"ipv6": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": SchemaForValuePrefixLength(),
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"destination": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"subnet_prefix": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": SchemaForValuePrefixLength(),
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
															"ipv6": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": SchemaForValuePrefixLength(),
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"protocol_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol_parameters": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"layer_four_protocol_object": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"source_port_ranges": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_port": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																		"end_port": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
															"destination_port_ranges": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_port": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																		"end_port": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
														},
													},
												},
												"icmp_object": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"icmp_type": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"icmp_code": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"protocol_number_object": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"protocol_number": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"policy_action": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reroute_params": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"service_ip": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLength(),
															"ipv6": SchemaForValuePrefixLength(),
														},
													},
												},
												"reroute_fallback_action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ingress_service_ip": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLength(),
															"ipv6": SchemaForValuePrefixLength(),
														},
													},
												},
												"egress_service_ip": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLength(),
															"ipv6": SchemaForValuePrefixLength(),
														},
													},
												},
											},
										},
									},
									"nexthop_ip_address": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLength(),
												"ipv6": SchemaForValuePrefixLength(),
											},
										},
									},
								},
							},
						},
						"is_bidirectional": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"vpc_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixPbrV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	extID := d.Get("ext_id")
	resp, err := conn.RoutingPolicy.GetRoutingPolicyById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching routing policy : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.RoutingPolicy)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("priority", getResp.Priority); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_ext_id", getResp.VpcExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policies", flattenPolicies(getResp.Policies)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc", flattenVpcName(getResp.Vpc)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenPolicies(pr []import1.RoutingPolicyRule) []interface{} {
	if len(pr) > 0 {
		policies := make([]interface{}, len(pr))

		for k, v := range pr {
			policy := make(map[string]interface{})

			policy["policy_match"] = flattenPolicyMatch(v.PolicyMatch)
			policy["policy_action"] = flattenRoutingPolicyAction(v.PolicyAction)
			policy["is_bidirectional"] = v.IsBidirectional

			policies[k] = policy
		}
		return policies
	}
	return nil
}

func flattenPolicyMatch(pr *import1.RoutingPolicyMatchCondition) []map[string]interface{} {
	if pr != nil {
		policyMatches := make([]map[string]interface{}, 0)
		match := make(map[string]interface{})

		match["source"] = flattenAddressTypeObject(pr.Source)
		match["destination"] = flattenAddressTypeObject(pr.Destination)
		match["protocol_type"] = flattenProtocolType(pr.ProtocolType)
		match["protocol_parameters"] = flattenOneOfRoutingPolicyMatchConditionProtocolParameters(pr.ProtocolParameters)

		policyMatches = append(policyMatches, match)
		return policyMatches
	}
	return nil
}

func flattenAddressTypeObject(pr *import1.AddressTypeObject) []map[string]interface{} {
	if pr != nil {
		addressType := make([]map[string]interface{}, 0)

		address := make(map[string]interface{})

		address["address_type"] = flattenAddressType(pr.AddressType)
		address["subnet_prefix"] = flattenIPSubnet(pr.SubnetPrefix)

		addressType = append(addressType, address)
		return addressType
	}
	return nil
}

func flattenIPSubnet(pr *import1.IPSubnet) []map[string]interface{} {
	if pr != nil {
		ipsubnets := make([]map[string]interface{}, 0)
		subnet := make(map[string]interface{})

		subnet["ipv4"] = flattenIPv4Subnet(pr.Ipv4)
		subnet["ipv6"] = flattenIPv6Subnet(pr.Ipv6)

		ipsubnets = append(ipsubnets, subnet)
		return ipsubnets
	}
	return nil
}

func flattenOneOfRoutingPolicyMatchConditionProtocolParameters(pr *import1.OneOfRoutingPolicyMatchConditionProtocolParameters) []map[string]interface{} {
	if pr != nil {
		layerFour := make(map[string]interface{})
		layerFourList := make([]map[string]interface{}, 0)
		icmp := make(map[string]interface{})
		icmpList := make([]map[string]interface{}, 0)
		protoNum := make(map[string]interface{})
		protoNumList := make([]map[string]interface{}, 0)

		if *pr.ObjectType_ == "networking.v4.config.LayerFourProtocolObject" {
			layer := make(map[string]interface{})
			layerList := make([]map[string]interface{}, 0)

			ip := pr.GetValue().(import1.LayerFourProtocolObject)

			layer["source_port_ranges"] = flattenPortRange(ip.SourcePortRanges)
			layer["destination_port_ranges"] = flattenPortRange(ip.DestinationPortRanges)

			layerList = append(layerList, layer)

			layerFour["layer_four_protocol_object"] = layerList

			layerFourList = append(layerFourList, layerFour)

			return layerFourList
		}

		if *pr.ObjectType_ == "networking.v4.config.ICMPObject" {
			obj := make(map[string]interface{})
			objList := make([]map[string]interface{}, 0)

			ip := pr.GetValue().(import1.ICMPObject)

			obj["icmp_code"] = ip.IcmpCode
			obj["icmp_type"] = ip.IcmpType

			objList = append(objList, obj)

			icmp["icmp_object"] = objList

			icmpList = append(icmpList, icmp)

			return icmpList
		}
		proto := make(map[string]interface{})
		protoList := make([]map[string]interface{}, 0)

		vm := pr.GetValue().(import1.ProtocolNumberObject)

		proto["protocol_number"] = vm.ProtocolNumber

		protoList = append(protoList, proto)

		protoNum["protocol_number_object"] = protoList

		protoNumList = append(protoNumList, protoNum)

		return protoNumList
	}
	return nil
}

func flattenRoutingPolicyAction(pr *import1.RoutingPolicyAction) []map[string]interface{} {
	if pr != nil {
		policyAction := make([]map[string]interface{}, 0)
		policy := make(map[string]interface{})

		policy["action_type"] = flattenRoutingPolicyActionType(pr.ActionType)
		policy["reroute_params"] = flattenRerouteParam(pr.RerouteParams)
		policy["nexthop_ip_address"] = flattenIPAddress(pr.NexthopIpAddress)

		policyAction = append(policyAction, policy)
		return policyAction
	}
	return nil
}

func flattenRerouteParam(pr []import1.RerouteParam) []interface{} {
	if len(pr) > 0 {
		routeParams := make([]interface{}, len(pr))

		for k, v := range pr {
			param := make(map[string]interface{})

			param["service_ip"] = flattenNodeIPAddress(v.ServiceIp)
			param["reroute_fallback_action"] = flattenRerouteFallbackAction(v.RerouteFallbackAction)

			routeParams[k] = param
		}
		return routeParams
	}
	return nil
}

func flattenPortRange(pr []import1.PortRange) []interface{} {
	if len(pr) > 0 {
		ports := make([]interface{}, len(pr))

		for k, v := range pr {
			port := make(map[string]interface{})

			port["start_port"] = v.StartPort
			port["end_port"] = v.EndPort

			ports[k] = port
		}
		return ports
	}
	return nil
}

func flattenAddressType(pr *import1.AddressType) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == import1.AddressType(two) {
			return "ANY"
		}
		if *pr == import1.AddressType(three) {
			return "EXTERNAL"
		}
		if *pr == import1.AddressType(four) {
			return "SUBNET"
		}
	}
	return "UNKNOWN"
}

func flattenRoutingPolicyActionType(pr *import1.RoutingPolicyActionType) string {
	if pr != nil {
		const two, three, four, five = 2, 3, 4, 5
		if *pr == import1.RoutingPolicyActionType(two) {
			return "PERMIT"
		}
		if *pr == import1.RoutingPolicyActionType(three) {
			return "DENY"
		}
		if *pr == import1.RoutingPolicyActionType(four) {
			return "REROUTE"
		}
		if *pr == import1.RoutingPolicyActionType(five) {
			return "FORWARD"
		}
	}
	return "UNKNOWN"
}

func flattenProtocolType(pr *import1.ProtocolType) string {
	if pr != nil {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		if *pr == import1.ProtocolType(two) {
			return "ANY"
		}
		if *pr == import1.ProtocolType(three) {
			return "ICMP"
		}
		if *pr == import1.ProtocolType(four) {
			return "TCP"
		}
		if *pr == import1.ProtocolType(five) {
			return "UDP"
		}
		if *pr == import1.ProtocolType(six) {
			return "PROTOCOL_NUMBER"
		}
	}
	return "UNKNOWN"
}

func flattenVpcName(pr *import1.VpcName) []map[string]interface{} {
	if pr != nil {
		vpcs := make([]map[string]interface{}, 0)
		vpc := make(map[string]interface{})

		vpc["name"] = pr.Name

		vpcs = append(vpcs, vpc)
		return vpcs
	}
	return nil
}

func flattenRerouteFallbackAction(pr *import1.RerouteFallbackAction) string {
	if pr != nil {
		const two, three, four, five = 2, 3, 4, 5
		if *pr == import1.RerouteFallbackAction(two) {
			return "ALLOW"
		}
		if *pr == import1.RerouteFallbackAction(three) {
			return "DROP"
		}
		if *pr == import1.RerouteFallbackAction(four) {
			return "PASSTHROUGH"
		}
		if *pr == import1.RerouteFallbackAction(five) {
			return "NO_ACTION"
		}
	}
	return "UNKNOWN"
}
