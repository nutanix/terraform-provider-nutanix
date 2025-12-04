package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNetworkSecurityPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixNetworkSecurityPoliciesV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
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
																Optional: true,
																Computed: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"dest_category_references": {
																Type:     schema.TypeList,
																Optional: true,
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
				},
			},
		},
	}
}

func DataSourceNutanixNetworkSecurityPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	// initialize query params
	var filter, orderBy, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	resp, err := conn.NetworkingSecurityInstance.ListNetworkSecurityPolicies(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching network security policy: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("network_policies", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of network security policies.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.NetworkSecurityPolicy)
	if err := d.Set("network_policies", flattenNetworkSecurityPolicy(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenNetworkSecurityPolicy(pr []import1.NetworkSecurityPolicy) []interface{} {
	if len(pr) > 0 {
		nets := make([]interface{}, len(pr))

		for k, v := range pr {
			net := make(map[string]interface{})

			net["ext_id"] = v.ExtId
			net["name"] = v.Name
			net["type"] = flattenSecurityPolicyType(v.Type)
			net["description"] = v.Description
			net["state"] = flattenPolicyState(v.State)
			net["rules"] = flattenNetworkSecurityPolicyRule(v.Rules)
			net["is_ipv6_traffic_allowed"] = v.IsIpv6TrafficAllowed
			net["is_hitlog_enabled"] = v.IsHitlogEnabled
			if v.Scope != nil {
				net["scope"] = flattenSecurityPolicyScope(v.Scope)
			}
			if v.VpcReferences != nil {
				net["vpc_reference"] = v.VpcReferences
			}
			if v.SecuredGroups != nil {
				net["secured_groups"] = v.SecuredGroups
			}
			if v.LastUpdateTime != nil {
				t := v.LastUpdateTime
				net["last_update_time"] = t.String()
			}
			if v.CreationTime != nil {
				t := v.CreationTime
				net["creation_time"] = t.String()
			}
			net["is_system_defined"] = v.IsSystemDefined
			net["created_by"] = v.CreatedBy

			if v.TenantId != nil {
				net["tenant_id"] = v.TenantId
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
