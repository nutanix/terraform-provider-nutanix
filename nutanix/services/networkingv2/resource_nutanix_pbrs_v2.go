package networkingv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	config "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixPbrsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixPbrsV2Create,
		ReadContext:   ResourceNutanixPbrsV2Read,
		UpdateContext: ResourceNutanixPbrsV2Update,
		DeleteContext: ResourceNutanixPbrsV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Optional: true,
				Type:     schema.TypeString,
				Computed: true,
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
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_match": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address_type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"SUBNET", "EXTERNAL", "ANY"}, false),
												},
												"subnet_prefix": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": SchemaForValuePrefixLength(),
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"ipv6": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": SchemaForValuePrefixLength(),
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Optional: true,
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
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address_type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"SUBNET", "EXTERNAL", "ANY"}, false),
												},
												"subnet_prefix": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": SchemaForValuePrefixLength(),
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"ipv6": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": SchemaForValuePrefixLength(),
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Optional: true,
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
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"TCP", "UDP", "PROTOCOL_NUMBER",
											"ANY", "ICMP",
										}, false),
									},
									"protocol_parameters": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"layer_four_protocol_object": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"source_port_ranges": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_port": {
																			Type:     schema.TypeInt,
																			Required: true,
																		},
																		"end_port": {
																			Type:     schema.TypeInt,
																			Required: true,
																		},
																	},
																},
															},
															"destination_port_ranges": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_port": {
																			Type:     schema.TypeInt,
																			Required: true,
																		},
																		"end_port": {
																			Type:     schema.TypeInt,
																			Required: true,
																		},
																	},
																},
															},
														},
													},
												},
												"icmp_object": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"icmp_type": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
															"icmp_code": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"protocol_number_object": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"protocol_number": {
																Type:     schema.TypeInt,
																Required: true,
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
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"REROUTE", "DENY", "FORWARD", "PERMIT"}, false),
									},
									"reroute_params": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"service_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLength(),
															"ipv6": SchemaForValuePrefixLength(),
														},
													},
												},
												"reroute_fallback_action": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													ValidateFunc: validation.StringInSlice([]string{"PASSTHROUGH", "NO_ACTION", "ALLOW", "DENY"}, false),
												},
												"ingress_service_ip": {
													Type:     schema.TypeList,
													Optional: true,
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
													Optional: true,
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
										Optional: true,
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
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"vpc_ext_id": {
				Type:     schema.TypeString,
				Required: true,
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

func ResourceNutanixPbrsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	inputSpec := import1.RoutingPolicy{}
	vpcRef := ""
	pbrPriority := 0
	pbrName := ""

	if name, ok := d.GetOk("name"); ok {
		inputSpec.Name = utils.StringPtr(name.(string))
		pbrName = name.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		inputSpec.Description = utils.StringPtr(description.(string))
	}
	if priority, ok := d.GetOk("priority"); ok {
		inputSpec.Priority = utils.IntPtr(priority.(int))
		pbrPriority = priority.(int)
	}
	if policies, ok := d.GetOk("policies"); ok {
		inputSpec.Policies = expandRoutingPolicyRule(policies.([]interface{}))
	}
	if vpcExtID, ok := d.GetOk("vpc_ext_id"); ok {
		inputSpec.VpcExtId = utils.StringPtr(vpcExtID.(string))
		vpcRef = vpcExtID.(string)
	}

	resp, err := conn.RoutingPolicy.CreateRoutingPolicy(&inputSpec)
	if err != nil {
		return diag.Errorf("error while deleting routing policy : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Routing Policy to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for routing policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from List Routing Policy API as Currently task entities does not return uuid
	filter := fmt.Sprintf("vpcExtId eq  '%s'", vpcRef)
	readResp, err := conn.RoutingPolicy.ListRoutingPolicies(nil, nil, &filter, nil, nil, nil)
	if err != nil {
		return diag.Errorf("error while fetching routing policies : %v", err)
	}

	for _, v := range readResp.Data.GetValue().([]import1.RoutingPolicy) {
		if utils.StringValue(v.Name) == pbrName && utils.IntValue(v.Priority) == pbrPriority {
			d.SetId(*v.ExtId)
			d.Set("ext_id", *v.ExtId)
			break
		}
	}
	return ResourceNutanixPbrsV2Read(ctx, d, meta)
}

func ResourceNutanixPbrsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.RoutingPolicy.GetRoutingPolicyById(utils.StringPtr(d.Id()))
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
	return nil
}

func ResourceNutanixPbrsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.RoutingPolicy.GetRoutingPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching routing policy : %v", err)
	}

	respVpc := resp.Data.GetValue().(import1.RoutingPolicy)

	updateSpec := respVpc

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("priority") {
		updateSpec.Priority = utils.IntPtr(d.Get("priority").(int))
	}
	if d.HasChange("policies") {
		updateSpec.Policies = expandRoutingPolicyRule(d.Get("policies").([]interface{}))
	}
	if d.HasChange("vpc_ext_id") {
		updateSpec.VpcExtId = utils.StringPtr(d.Get("vpc_ext_id").(string))
	}

	updateResp, err := conn.RoutingPolicy.UpdateRoutingPolicyById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating routing policy : %v", err)
	}
	TaskRef := updateResp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Routing Policy to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for routing policy (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixPbrsV2Read(ctx, d, meta)
}

func ResourceNutanixPbrsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.RoutingPolicy.DeleteRoutingPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting routing policy : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Subnet to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for routing policy (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandRoutingPolicyRule(pr []interface{}) []import1.RoutingPolicyRule {
	if len(pr) > 0 {
		rules := make([]import1.RoutingPolicyRule, len(pr))
		for k, v := range pr {
			val := v.(map[string]interface{})
			rule := import1.RoutingPolicyRule{}
			if policyMatch, ok := val["policy_match"]; ok && len(policyMatch.([]interface{})) > 0 {
				rule.PolicyMatch = expandRoutingPolicyMatchCondition(policyMatch)
			}
			if policyAction, ok := val["policy_action"]; ok && len(policyAction.([]interface{})) > 0 {
				rule.PolicyAction = expandRoutingPolicyAction(policyAction)
			}
			if isBidirectional, ok := val["is_bidirectional"]; ok {
				rule.IsBidirectional = utils.BoolPtr(isBidirectional.(bool))
			}
			rules[k] = rule
		}
		return rules
	}
	return nil
}

func expandRoutingPolicyMatchCondition(pr interface{}) *import1.RoutingPolicyMatchCondition {
	if pr != nil {
		match := import1.RoutingPolicyMatchCondition{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if source, ok := val["source"]; ok && len(source.([]interface{})) > 0 {
			match.Source = expandAddressTypeObject(source)
		}
		if destination, ok := val["destination"]; ok && len(destination.([]interface{})) > 0 {
			match.Destination = expandAddressTypeObject(destination)
		}
		if protocolType, ok := val["protocol_type"]; ok {
			const two, three, four, five, six = 2, 3, 4, 5, 6
			protoMap := map[string]interface{}{
				"ANY":             two,
				"ICMP":            three,
				"TCP":             four,
				"UDP":             five,
				"PROTOCOL_NUMBER": six,
			}

			pInt := protoMap[protocolType.(string)]
			p := import1.ProtocolType(pInt.(int))
			match.ProtocolType = &p
		}
		if protocolParameters, ok := val["protocol_parameters"]; ok && len(protocolParameters.([]interface{})) > 0 {
			match.ProtocolParameters = expandOneOfRoutingPolicyMatchConditionProtocolParameters(protocolParameters)
		}
		return &match
	}
	return nil
}

func expandAddressTypeObject(pr interface{}) *import1.AddressTypeObject {
	if pr != nil {
		address := import1.AddressTypeObject{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if addressType, ok := val["address_type"]; ok {
			const two, three, four = 2, 3, 4
			addMap := map[string]interface{}{
				"ANY":      two,
				"EXTERNAL": three,
				"SUBNET":   four,
			}
			pInt := addMap[addressType.(string)]
			p := import1.AddressType(pInt.(int))
			address.AddressType = &p
		}
		if subnetPrefix, ok := val["subnet_prefix"]; ok && len(subnetPrefix.([]interface{})) > 0 {
			address.SubnetPrefix = expandIPSubnetObject(subnetPrefix)
		}
		return &address
	}
	return nil
}

func expandIPSubnetObject(pr interface{}) *import1.IPSubnet {
	if pr != nil {
		subnet := import1.IPSubnet{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
			subnet.Ipv4 = expandIPv4Subnet(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
			subnet.Ipv6 = expandIPv6Subnet(ipv6)
		}
		return &subnet
	}
	return nil
}

func expandOneOfRoutingPolicyMatchConditionProtocolParameters(pr interface{}) *import1.OneOfRoutingPolicyMatchConditionProtocolParameters {
	if pr != nil {
		protoParams := import1.OneOfRoutingPolicyMatchConditionProtocolParameters{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if layerFourProtocolObject, ok := val["layer_four_protocol_object"]; ok && len(layerFourProtocolObject.([]interface{})) > 0 {
			layerFourObj := import1.NewLayerFourProtocolObject()
			layerI := layerFourProtocolObject.([]interface{})
			layerVal := layerI[0].(map[string]interface{})

			if sourcePortRanges, ok := layerVal["source_port_ranges"]; ok && len(sourcePortRanges.([]interface{})) > 0 {
				layerFourObj.SourcePortRanges = expandPortRange(sourcePortRanges.([]interface{}))
			}
			if destinationPortRanges, ok := layerVal["destination_port_ranges"]; ok && len(destinationPortRanges.([]interface{})) > 0 {
				layerFourObj.DestinationPortRanges = expandPortRange(destinationPortRanges.([]interface{}))
			}

			protoParams.SetValue(*layerFourObj)
		}

		if icmpObject, ok := val["icmp_object"]; ok && len(icmpObject.([]interface{})) > 0 {
			icmpObj := import1.NewICMPObject()
			icmpI := icmpObject.([]interface{})
			icmpVal := icmpI[0].(map[string]interface{})
			if icmpType, ok := icmpVal["icmp_type"]; ok {
				icmpObj.IcmpType = utils.IntPtr(icmpType.(int))
			}
			if icmpCode, ok := icmpVal["icmp_code"]; ok {
				icmpObj.IcmpCode = utils.IntPtr(icmpCode.(int))
			}
			protoParams.SetValue(*icmpObj)
		}

		if protoNum, ok := val["protocol_number_object"]; ok && len(protoNum.([]interface{})) > 0 {
			protoNumObj := import1.NewProtocolNumberObject()
			protoI := protoNum.([]interface{})
			protoVal := protoI[0].(map[string]interface{})
			if protocolNumber, ok := protoVal["protocol_number"]; ok {
				protoNumObj.ProtocolNumber = utils.IntPtr(protocolNumber.(int))
			}
			protoParams.SetValue(*protoNumObj)
		}
		return &protoParams
	}
	return nil
}

func expandPortRange(pr []interface{}) []import1.PortRange {
	if len(pr) > 0 {
		portRanges := make([]import1.PortRange, len(pr))
		for k, v := range pr {
			val := v.(map[string]interface{})
			portRange := import1.PortRange{}
			if startPort, ok := val["start_port"]; ok {
				portRange.StartPort = utils.IntPtr(startPort.(int))
			}
			if endPort, ok := val["end_port"]; ok {
				portRange.EndPort = utils.IntPtr(endPort.(int))
			}
			portRanges[k] = portRange
		}
		return portRanges
	}
	return nil
}

func expandRoutingPolicyAction(pr interface{}) *import1.RoutingPolicyAction {
	if pr != nil {
		action := import1.RoutingPolicyAction{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if actionType, ok := val["action_type"]; ok {
			const two, three, four, five = 2, 3, 4, 5
			actMap := map[string]interface{}{
				"PERMIT":  two,
				"DENY":    three,
				"REROUTE": four,
				"FORWARD": five,
			}
			pInt := actMap[actionType.(string)]
			p := import1.RoutingPolicyActionType(pInt.(int))
			action.ActionType = &p
		}
		if rerouteParams, ok := val["reroute_params"]; ok && len(rerouteParams.([]interface{})) > 0 {
			action.RerouteParams = expandRerouteParams(rerouteParams.([]interface{}))
		}
		if nexthopIPAddress, ok := val["nexthop_ip_address"]; ok && len(nexthopIPAddress.([]interface{})) > 0 {
			action.NexthopIpAddress = expandIPAddressObject(nexthopIPAddress)
		}
		return &action
	}
	return nil
}

func expandRerouteParams(pr []interface{}) []import1.RerouteParam {
	if len(pr) > 0 {
		reroutes := make([]import1.RerouteParam, len(pr))
		for k, v := range pr {
			val := v.(map[string]interface{})
			reroute := import1.RerouteParam{}
			if serviceIP, ok := val["service_ip"]; ok && len(serviceIP.([]interface{})) > 0 {
				reroute.ServiceIp = expandIPAddressObject(serviceIP)
			}
			if rerouteFallbackAction, ok := val["reroute_fallback_action"]; ok {
				const two, three, four, five = 2, 3, 4, 5
				actMap := map[string]interface{}{
					"ALLOW":       two,
					"DROP":        three,
					"PASSTHROUGH": four,
					"NO_ACTION":   five,
				}

				pInt := actMap[rerouteFallbackAction.(string)]
				p := import1.RerouteFallbackAction(pInt.(int))
				reroute.RerouteFallbackAction = &p
			}
			if ingressServiceIP, ok := val["ingress_service_ip"]; ok && len(ingressServiceIP.([]interface{})) > 0 {
				reroute.IngressServiceIp = expandIPAddressObject(ingressServiceIP)
			}
			if egressServiceIP, ok := val["egress_service_ip"]; ok && len(egressServiceIP.([]interface{})) > 0 {
				reroute.EgressServiceIp = expandIPAddressObject(egressServiceIP)
			}
			reroutes[k] = reroute
		}
		return reroutes
	}
	return nil
}

func expandIPAddressObject(pr interface{}) *config.IPAddress {
	if pr != nil {
		ip := config.IPAddress{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
			ip.Ipv4 = expandIPv4Address(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
			ip.Ipv6 = expandIPv6Address(ipv6)
		}
		return &ip
	}
	return nil
}
