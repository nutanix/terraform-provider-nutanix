package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixPbrsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixPbrsV2Read,
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
			"routing_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
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
				},
			},
		},
	}
}

func DatasourceNutanixPbrsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	// initialize query params
	var filter, orderBy *string
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

	resp, err := conn.RoutingPolicy.ListRoutingPolicies(page, limit, filter, orderBy, nil, nil)
	if err != nil {
		return diag.Errorf("error while fetching routing policies : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("routing_policies", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of routing policies.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.RoutingPolicy)

	if err := d.Set("routing_policies", flattenRoutingEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRoutingEntities(pr []import1.RoutingPolicy) []interface{} {
	if len(pr) > 0 {
		routingEnts := make([]interface{}, len(pr))

		for k, v := range pr {
			entity := make(map[string]interface{})

			entity["ext_id"] = v.ExtId
			entity["tenant_id"] = v.TenantId
			entity["links"] = flattenLinks(v.Links)
			entity["metadata"] = flattenMetadata(v.Metadata)
			entity["name"] = v.Name
			entity["description"] = v.Description
			entity["priority"] = v.Priority
			entity["policies"] = flattenPolicies(v.Policies)
			entity["vpc_ext_id"] = v.VpcExtId
			entity["vpc"] = flattenVpcName(v.Vpc)

			routingEnts[k] = entity
		}
		return routingEnts
	}
	return nil
}
