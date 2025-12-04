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

func DataSourceNutanixSubnetsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixSubnetsV2Read,
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
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dhcp_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_name_servers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLength(),
												"ipv6": SchemaForValuePrefixLength(),
											},
										},
									},
									"domain_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"search_domains": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"tftp_server_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"boot_file_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ntp_servers": {
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
						"ip_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_subnet": {
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
												"default_gateway_ip":  SchemaForValuePrefixLength(),
												"dhcp_server_address": SchemaForValuePrefixLength(),
												"pool_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_ip": {
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
															"end_ip": {
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
														},
													},
												},
											},
										},
									},
									"ipv6": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_subnet": {
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
												"default_gateway_ip":  SchemaForValuePrefixLength(),
												"dhcp_server_address": SchemaForValuePrefixLength(),
												"pool_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_ip": {
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
															"end_ip": {
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
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"cluster_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_switch_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_nat_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_external": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"reserved_ip_addresses": SchemaForValuePrefixLength(),
						"dynamic_ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
						"network_function_chain_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bridge_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_advanced_networking": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_switch": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DataSourceVirtualSwitchSchemaV2(),
							},
						},
						"vpc": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DataSourceVPCSchemaV2(),
							},
						},
						"ip_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_usage": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_macs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"num_free_ips": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"num_assigned_ips": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ip_pool_usages": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_free_ips": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"num_total_ips": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"range": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_ip": {
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
															"end_ip": {
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
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"migration_state": {
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

func dataSourceNutanixSubnetsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	// initialize query params
	var filter, orderBy, expand, selects *string
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
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	resp, err := conn.SubnetAPIInstance.ListSubnets(page, limit, filter, orderBy, expand, selects)
	if err != nil {
		return diag.Errorf("error while fetching subnets : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("subnets", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of subnets.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.Subnet)

	if err := d.Set("subnets", flattenSubnetEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenSubnetEntities(pr []import1.Subnet) []interface{} {
	if len(pr) > 0 {
		subnets := make([]interface{}, len(pr))

		for k, v := range pr {
			sub := make(map[string]interface{})

			sub["ext_id"] = v.ExtId
			sub["name"] = v.Name
			sub["description"] = v.Description
			sub["links"] = flattenLinks(v.Links)
			sub["subnet_type"] = flattenSubnetType(v.SubnetType)
			sub["network_id"] = v.NetworkId
			sub["dhcp_options"] = flattenDhcpOptions(v.DhcpOptions)
			sub["ip_config"] = flattenIPConfig(v.IpConfig)
			sub["cluster_reference"] = v.ClusterReference
			sub["virtual_switch_reference"] = v.VirtualSwitchReference
			sub["vpc_reference"] = v.VpcReference
			sub["is_nat_enabled"] = v.IsNatEnabled
			sub["is_external"] = v.IsExternal
			sub["reserved_ip_addresses"] = flattenReservedIPAddresses(v.ReservedIpAddresses)
			sub["dynamic_ip_addresses"] = flattenReservedIPAddresses(v.DynamicIpAddresses)
			sub["network_function_chain_reference"] = v.NetworkFunctionChainReference
			sub["bridge_name"] = v.BridgeName
			sub["is_advanced_networking"] = v.IsAdvancedNetworking
			sub["cluster_name"] = v.ClusterName
			sub["hypervisor_type"] = v.HypervisorType
			sub["virtual_switch"] = flattenVirtualSwitch(v.VirtualSwitch)
			sub["vpc"] = flattenVPC(v.Vpc)
			sub["ip_prefix"] = v.IpPrefix
			sub["ip_usage"] = flattenIPUsage(v.IpUsage)
			sub["migration_state"] = flattenMigrationState(v.MigrationState)

			subnets[k] = sub
		}
		return subnets
	}
	return nil
}
