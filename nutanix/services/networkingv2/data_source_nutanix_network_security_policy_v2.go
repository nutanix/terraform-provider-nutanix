package networkingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNetworkSecurityPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixNetworkSecurityPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"two_env_isolation_rule_spec": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"first_isolation_group": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"second_isolation_group": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"application_rule_spec": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"secured_group_category_references": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"src_allow_spec": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"dest_allow_spec": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"src_category_references": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"dest_category_references": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"src_subnet": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"dest_subnet": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"src_address_group_references": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"dest_address_group_references": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"service_group_references": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"is_all_protocol_allowed": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"tcp_services": {
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
												"udp_services": {
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
												"icmp_services": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_all_allowed": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"type": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"code": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"network_function_chain_reference": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"intra_entity_group_rule_spec": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"secured_group_action": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"secured_group_category_references": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"multi_env_isolation_rule_spec": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"spec": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"all_to_all_isolation_group": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"isolation_group": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"group_category_references": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
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
												},
											},
										},
									},
								},
							},
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
					},
				},
			},
			"is_ipv6_traffic_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_hitlog_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"secured_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_by": {
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
		},
	}
}

func DataSourceNutanixNetworkSecurityPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id")

	resp, err := conn.NetworkingSecurityInstance.GetNetworkSecurityPolicyById(utils.StringPtr((extID.(string))))
	if err != nil {
		return diag.Errorf("error while fetching network security policy: %v", err)
	}
	getResp := resp.Data.GetValue().(import1.NetworkSecurityPolicy)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", flattenSecurityPolicyType(getResp.Type)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", flattenPolicyState(getResp.State)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rules", flattenNetworkSecurityPolicyRule(getResp.Rules)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_ipv6_traffic_allowed", getResp.IsIpv6TrafficAllowed); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_hitlog_enabled", getResp.IsHitlogEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scope", flattenSecurityPolicyScope(getResp.Scope)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vpc_reference", getResp.VpcReferences); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("secured_groups", getResp.SecuredGroups); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreationTime != nil {
		t := getResp.CreationTime
		if err := d.Set("creation_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdateTime != nil {
		t := getResp.LastUpdateTime
		if err := d.Set("last_update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("is_system_defined", getResp.IsSystemDefined); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksMicroSeg(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}

func flattenNetworkSecurityPolicyRule(pr []import1.NetworkSecurityPolicyRule) []interface{} {
	if len(pr) > 0 {
		nets := make([]interface{}, len(pr))

		for k, v := range pr {
			net := make(map[string]interface{})

			if v.ExtId != nil {
				net["ext_id"] = v.ExtId
			}
			if v.Description != nil {
				net["description"] = v.Description
			}
			if v.Type != nil {
				net["type"] = flattenRuleType(v.Type)
			}
			if v.Spec != nil {
				net["spec"] = flattenOneOfNetworkSecurityPolicyRuleSpec(v.Spec)
			}
			if v.Links != nil {
				net["links"] = flattenLinksMicroSeg(v.Links)
			}
			nets[k] = net
		}
		return nets
	}
	return nil
}

func flattenOneOfNetworkSecurityPolicyRuleSpec(pr *import1.OneOfNetworkSecurityPolicyRuleSpec) []map[string]interface{} {
	if pr != nil {
		isolationRuleSpec := make(map[string]interface{})
		isolationRuleSpecList := make([]map[string]interface{}, 0)
		appRuleSpec := make(map[string]interface{})
		appRuleSpecList := make([]map[string]interface{}, 0)
		intraRuleSpec := make(map[string]interface{})
		intraRuleSpecList := make([]map[string]interface{}, 0)
		multiEnvIsolationRuleSpec := make(map[string]interface{})
		multiEnvIsolationRuleSpecList := make([]map[string]interface{}, 0)

		if *pr.ObjectType_ == "microseg.v4.config.TwoEnvIsolationRuleSpec" {
			env := make(map[string]interface{})
			envList := make([]map[string]interface{}, 0)

			isolationValue := pr.GetValue().(import1.TwoEnvIsolationRuleSpec)

			env["first_isolation_group"] = isolationValue.FirstIsolationGroup
			env["second_isolation_group"] = isolationValue.SecondIsolationGroup

			envList = append(envList, env)

			isolationRuleSpec["two_env_isolation_rule_spec"] = envList

			isolationRuleSpecList = append(isolationRuleSpecList, isolationRuleSpec)

			return isolationRuleSpecList
		}
		if *pr.ObjectType_ == "microseg.v4.config.ApplicationRuleSpec" {
			app := make(map[string]interface{})
			appList := make([]map[string]interface{}, 0)

			appRuleValue := pr.GetValue().(import1.ApplicationRuleSpec)

			if appRuleValue.SecuredGroupCategoryReferences != nil {
				app["secured_group_category_references"] = appRuleValue.SecuredGroupCategoryReferences
			}
			if appRuleValue.SrcAllowSpec != nil {
				app["src_allow_spec"] = flattenAllowType(appRuleValue.SrcAllowSpec)
			}
			if appRuleValue.DestAllowSpec != nil {
				app["dest_allow_spec"] = flattenAllowType(appRuleValue.DestAllowSpec)
			}
			if appRuleValue.SrcCategoryReferences != nil {
				app["src_category_references"] = appRuleValue.SrcCategoryReferences
			}
			if appRuleValue.DestCategoryReferences != nil {
				app["dest_category_references"] = appRuleValue.DestCategoryReferences
			}
			if appRuleValue.SrcSubnet != nil {
				app["src_subnet"] = flattenIPv4AddressMicroSegList(appRuleValue.SrcSubnet)
			}
			if appRuleValue.DestSubnet != nil {
				app["dest_subnet"] = flattenIPv4AddressMicroSegList(appRuleValue.DestSubnet)
			}
			if appRuleValue.SrcAddressGroupReferences != nil {
				app["src_address_group_references"] = appRuleValue.SrcAddressGroupReferences
			}
			if appRuleValue.DestAddressGroupReferences != nil {
				app["dest_address_group_references"] = appRuleValue.DestAddressGroupReferences
			}
			if appRuleValue.ServiceGroupReferences != nil {
				app["service_group_references"] = appRuleValue.ServiceGroupReferences
			}
			if appRuleValue.IsAllProtocolAllowed != nil {
				app["is_all_protocol_allowed"] = appRuleValue.IsAllProtocolAllowed
			}
			if appRuleValue.TcpServices != nil {
				app["tcp_services"] = flattenTCPPortRangeSpec(appRuleValue.TcpServices)
			}
			if appRuleValue.UdpServices != nil {
				app["udp_services"] = flattenUDPPortRangeSpec(appRuleValue.UdpServices)
			}
			if appRuleValue.IcmpServices != nil {
				app["icmp_services"] = flattenIcmpTypeCodeSpec(appRuleValue.IcmpServices)
			}
			if appRuleValue.NetworkFunctionChainReference != nil {
				app["network_function_chain_reference"] = appRuleValue.NetworkFunctionChainReference
			}

			appList = append(appList, app)

			appRuleSpec["application_rule_spec"] = appList

			appRuleSpecList = append(appRuleSpecList, appRuleSpec)
			return appRuleSpecList
		}
		if *pr.ObjectType_ == "microseg.v4.config.IntraEntityGroupRuleSpec" {
			intra := make(map[string]interface{})
			intraList := make([]map[string]interface{}, 0)

			intraRuleValue := pr.GetValue().(import1.IntraEntityGroupRuleSpec)

			if intraRuleValue.SecuredGroupAction != nil {
				intra["secured_group_action"] = flattenIntraEntityGroupRuleAction(intraRuleValue.SecuredGroupAction)
			}
			if intraRuleValue.SecuredGroupCategoryReferences != nil {
				intra["secured_group_category_references"] = intraRuleValue.SecuredGroupCategoryReferences
			}

			intraList = append(intraList, intra)

			intraRuleSpec["intra_entity_group_rule_spec"] = intraList
			intraRuleSpecList = append(intraRuleSpecList, intraRuleSpec)

			return intraRuleSpecList
		}
		if *pr.ObjectType_ == "microseg.v4.config.MultiEnvIsolationRuleSpec" {
			multiEenv := make(map[string]interface{})
			multiEnvList := make([]map[string]interface{}, 0)

			specMap := make([]map[string]interface{}, 0)
			allToAllIsolationGroup := make([]map[string]interface{}, 0)
			isolationGroups := make([]map[string]interface{}, 0)
			groupCategoryRef := make(map[string]interface{})

			multiEnvIsolationValue := pr.GetValue().(import1.MultiEnvIsolationRuleSpec)

			allIsolationGroupValue := multiEnvIsolationValue.Spec.GetValue().(import1.AllToAllIsolationGroup)

			for _, group := range allIsolationGroupValue.IsolationGroups {
				groupCategoryRef["group_category_reference"] = group.GroupCategoryReferences
				isolationGroups = append(isolationGroups, groupCategoryRef)
			}

			allToAllIsolationGroup = append(allToAllIsolationGroup, isolationGroups...)

			specMap = append(specMap, allToAllIsolationGroup...)

			multiEenv["multi_env_isolation_rule_spec"] = specMap

			multiEnvList = append(multiEnvList, multiEenv)

			multiEnvIsolationRuleSpec["multi_env_isolation_rule_spec"] = multiEnvList

			multiEnvIsolationRuleSpecList = append(multiEnvIsolationRuleSpecList, multiEnvIsolationRuleSpec)

			aJSON, _ := json.Marshal(multiEnvIsolationRuleSpecList)
			log.Printf("[DEBUG] multiEnvIsolationRuleSpecList: %s", string(aJSON))
			return multiEnvIsolationRuleSpecList
		}
	}
	return nil
}

func flattenIPv4AddressMicroSegList(pr *config.IPv4Address) []interface{} {
	if pr != nil {
		ipv4 := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = pr.Value
		ip["prefix_length"] = pr.PrefixLength

		ipv4 = append(ipv4, ip)

		return ipv4
	}
	return nil
}

func flattenAllowType(pr *import1.AllowType) string {
	const two, three = 2, 3
	if pr != nil {
		if *pr == import1.AllowType(two) {
			return "ALL"
		}
		if *pr == import1.AllowType(three) {
			return "NONE"
		}
	}
	return "UNKNOWN"
}

func flattenPolicyState(pr *import1.SecurityPolicyState) string {
	const two, three, four = 2, 3, 4
	if pr != nil {
		if *pr == import1.SecurityPolicyState(two) {
			return "SAVE"
		}
		if *pr == import1.SecurityPolicyState(three) {
			return "MONITOR"
		}
		if *pr == import1.SecurityPolicyState(four) {
			return "ENFORCE"
		}
	}
	return "UNKNOWN"
}

func flattenRuleType(pr *import1.RuleType) string {
	const two, three, four, five, six = 2, 3, 4, 5, 6
	if pr != nil {
		if *pr == import1.RuleType(two) {
			return "QUARANTINE"
		}
		if *pr == import1.RuleType(three) {
			return "TWO_ENV_ISOLATION"
		}
		if *pr == import1.RuleType(four) {
			return "APPLICATION"
		}
		if *pr == import1.RuleType(five) {
			return "ENFORCE"
		}
		if *pr == import1.RuleType(six) {
			return "MULTI_ENV_ISOLATION"
		}
	}
	return "UNKNOWN"
}

func flattenSecurityPolicyType(pr *import1.SecurityPolicyType) string {
	const two, three, four = 2, 3, 4
	if pr != nil {
		if *pr == import1.SecurityPolicyType(two) {
			return "QUARANTINE"
		}
		if *pr == import1.SecurityPolicyType(three) {
			return "ISOLATION"
		}
		if *pr == import1.SecurityPolicyType(four) {
			return "APPLICATION"
		}
	}
	return "UNKNOWN"
}

func flattenSecurityPolicyScope(pr *import1.SecurityPolicyScope) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == import1.SecurityPolicyScope(two) {
			return "ALL_VLAN"
		}
		if *pr == import1.SecurityPolicyScope(three) {
			return "ALL_VPC"
		}
		if *pr == import1.SecurityPolicyScope(four) {
			return "VPC_LIST"
		}
	}
	return "UNKNOWN"
}

func flattenIntraEntityGroupRuleAction(pr *import1.IntraEntityGroupRuleAction) string {
	if pr != nil {
		const two, three = 2, 3

		if *pr == import1.IntraEntityGroupRuleAction(two) {
			return "ALLOW"
		}
		if *pr == import1.IntraEntityGroupRuleAction(three) {
			return "DENY"
		}
	}
	return "UNKNOWN"
}
