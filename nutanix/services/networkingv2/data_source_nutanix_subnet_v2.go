package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixSubnetV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixSubnetV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Required: true,
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
	}
}

func dataSourceNutanixSubnetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	extID := d.Get("ext_id")
	resp, err := conn.SubnetAPIInstance.GetSubnetById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching subnets : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Subnet)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subnet_type", flattenSubnetType(getResp.SubnetType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_id", getResp.NetworkId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dhcp_options", flattenDhcpOptions(getResp.DhcpOptions)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_config", flattenIPConfig(getResp.IpConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_reference", getResp.ClusterReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("virtual_switch_reference", getResp.VirtualSwitchReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_nat_enabled", getResp.IsNatEnabled); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_external", getResp.IsExternal); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("reserved_ip_addresses", flattenReservedIPAddresses(getResp.ReservedIpAddresses)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dynamic_ip_addresses", flattenReservedIPAddresses(getResp.DynamicIpAddresses)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_function_chain_reference", getResp.NetworkFunctionChainReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("bridge_name", getResp.BridgeName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_advanced_networking", getResp.IsAdvancedNetworking); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_name", getResp.ClusterName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor_type", getResp.HypervisorType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("virtual_switch", flattenVirtualSwitch(getResp.VirtualSwitch)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc", flattenVPC(getResp.Vpc)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_prefix", getResp.IpPrefix); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip_usage", flattenIPUsage(getResp.IpUsage)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("migration_state", flattenMigrationState(getResp.MigrationState)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extID.(string))
	return nil
}

func DataSourceVPCSchemaV2() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vpc_type": {
			Type:     schema.TypeString,
			Optional: true,
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
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: DatasourceMetadataSchemaV2(),
			},
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"common_dhcp_options": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"domain_name_servers": {
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
		"snat_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ipv4": SchemaForValuePrefixLength(),
					"ipv6": SchemaForValuePrefixLength(),
				},
			},
		},
		"external_subnets": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"subnet_reference": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"external_ips": {
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
					"gateway_nodes": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"active_gateway_node": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"node_id": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"node_ip_address": {
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
					"active_gateway_count": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"external_routing_domain_reference": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"externally_routable_prefixes": {
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
	}
}

func DataSourceVirtualSwitchSchemaV2() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Optional: true,
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
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: DatasourceMetadataSchemaV2(),
			},
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"is_default": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"has_deployment_error": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"mtu": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"bond_mode": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"clusters": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ext_id": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"hosts": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ext_id": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"internal_bridge_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"host_nics": {
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"ip_address": {
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
								"route_table": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"gateway_ip_address": SchemaForValuePrefixLength(),
				},
			},
		},
	}
}

func SchemaForValuePrefixLength() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"prefix_length": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func DatasourceMetadataSchemaV2() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"owner_reference_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"owner_user_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"project_reference_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"project_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"category_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func flattenLinks(pr []import2.ApiLink) []map[string]interface{} {
	if len(pr) > 0 {
		linkList := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			links := map[string]interface{}{}
			if v.Href != nil {
				links["href"] = v.Href
			}
			if v.Rel != nil {
				links["rel"] = v.Rel
			}

			linkList[k] = links
		}
		return linkList
	}
	return nil
}

func flattenDhcpOptions(pr *import1.DhcpOptions) []interface{} {
	if pr != nil {
		dhcpOps := make([]interface{}, 0)

		dhcp := make(map[string]interface{})

		dhcp["domain_name_servers"] = flattenNtpServer(pr.DomainNameServers)
		dhcp["domain_name"] = pr.DomainName
		dhcp["search_domains"] = pr.SearchDomains
		dhcp["tftp_server_name"] = pr.TftpServerName
		dhcp["boot_file_name"] = pr.BootFileName
		dhcp["ntp_servers"] = flattenNtpServer(pr.NtpServers)

		dhcpOps = append(dhcpOps, dhcp)

		return dhcpOps
	}
	return nil
}

func flattenNtpServer(pr []config.IPAddress) []map[string]interface{} {
	if len(pr) > 0 {
		ips := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ip := make(map[string]interface{})

			ip["ipv4"] = flattenIPv4(v.Ipv4)
			ip["ipv6"] = flattenIPv6(v.Ipv6)

			ips[k] = ip
		}
		return ips
	}
	return nil
}

func flattenIPv4(pr *config.IPv4Address) []interface{} {
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

func flattenIPv6(pr *config.IPv6Address) []interface{} {
	if pr != nil {
		ipv6 := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = pr.Value
		ip["prefix_length"] = pr.PrefixLength

		ipv6 = append(ipv6, ip)

		return ipv6
	}
	return nil
}

func flattenIPConfig(pr []import1.IPConfig) []map[string]interface{} {
	if len(pr) > 0 {
		ipCfgs := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ip := make(map[string]interface{})

			ip["ipv4"] = flattenIpv4Config(v.Ipv4)
			ip["ipv6"] = flattenIpv6Config(v.Ipv6)

			ipCfgs[k] = ip
		}
		return ipCfgs
	}
	return nil
}

func flattenIPv4Subnet(pr *import1.IPv4Subnet) []interface{} {
	if pr != nil {
		subs := make([]interface{}, 0)

		sub := make(map[string]interface{})

		sub["ip"] = flattenIPv4(pr.Ip)
		sub["prefix_length"] = pr.PrefixLength

		subs = append(subs, sub)
		return subs
	}
	return nil
}

func flattenIPv6Subnet(pr *import1.IPv6Subnet) []interface{} {
	if pr != nil {
		subs := make([]interface{}, 0)

		sub := make(map[string]interface{})

		sub["ip"] = flattenIPv6(pr.Ip)
		sub["prefix_length"] = pr.PrefixLength

		subs = append(subs, sub)
		return subs
	}
	return nil
}

func flattenPoolListIPv4(pr []import1.IPv4Pool) []map[string]interface{} {
	if len(pr) > 0 {
		poolList := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			pool := make(map[string]interface{})

			pool["start_ip"] = flattenIPv4(v.StartIp)
			pool["end_ip"] = flattenIPv4(v.EndIp)

			poolList[k] = pool
		}
		return poolList
	}
	return nil
}

func flattenPoolListIPv6(pr []import1.IPv6Pool) []map[string]interface{} {
	if len(pr) > 0 {
		poolList := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			pool := make(map[string]interface{})

			pool["start_ip"] = flattenIPv6(v.StartIp)
			pool["end_ip"] = flattenIPv6(v.EndIp)

			poolList[k] = pool
		}
		return poolList
	}
	return nil
}

func flattenIpv4Config(pr *import1.IPv4Config) []interface{} {
	if pr != nil {
		ipCfg := make([]interface{}, 0)

		cfg := make(map[string]interface{})

		cfg["ip_subnet"] = flattenIPv4Subnet(pr.IpSubnet)
		cfg["default_gateway_ip"] = flattenIPv4(pr.DefaultGatewayIp)
		cfg["dhcp_server_address"] = flattenIPv4(pr.DhcpServerAddress)
		cfg["pool_list"] = flattenPoolListIPv4(pr.PoolList)

		ipCfg = append(ipCfg, cfg)
		return ipCfg
	}
	return nil
}

func flattenIpv6Config(pr *import1.IPv6Config) []interface{} {
	if pr != nil {
		ipCfg := make([]interface{}, 0)

		cfg := make(map[string]interface{})

		cfg["ip_subnet"] = flattenIPv6Subnet(pr.IpSubnet)
		cfg["default_gateway_ip"] = flattenIPv6(pr.DefaultGatewayIp)
		cfg["dhcp_server_address"] = flattenIPv6(pr.DhcpServerAddress)
		cfg["pool_list"] = flattenPoolListIPv6(pr.PoolList)

		ipCfg = append(ipCfg, cfg)
		return ipCfg
	}
	return nil
}

func flattenSubnetType(sb *import1.SubnetType) string {
	const two, three = 2, 3
	if sb != nil {
		if *sb == import1.SubnetType(two) {
			return "OVERLAY"
		}
		if *sb == import1.SubnetType(three) {
			return "VLAN"
		}
	}
	return "UNKNOWN"
}

func flattenReservedIPAddresses(pr []config.IPAddress) []map[string]interface{} {
	if len(pr) > 0 {
		ipsList := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ip := make(map[string]interface{})

			ip["ipv4"] = flattenIPv4(v.Ipv4)
			ip["ipv6"] = flattenIPv6(v.Ipv6)

			ipsList[k] = ip
		}
		return ipsList
	}
	return nil
}

func flattenVirtualSwitch(vs *import1.VirtualSwitch) []map[string]interface{} {
	if vs != nil {
		vSwitch := make([]map[string]interface{}, 0)

		v := make(map[string]interface{})

		if vs.TenantId != nil {
			v["tenant_id"] = vs.TenantId
		}

		v["ext_id"] = vs.ExtId
		v["name"] = vs.Name
		v["description"] = vs.Description
		v["is_default"] = vs.IsDefault
		v["has_deployment_error"] = vs.HasDeploymentError
		v["mtu"] = vs.Mtu
		v["bond_mode"] = flattenBondMode(vs.BondMode)
		v["clusters"] = flattenClusters(vs.Clusters)
		v["metadata"] = flattenMetadata(vs.Metadata)
		v["links"] = flattenLinks(vs.Links)

		vSwitch = append(vSwitch, v)
		return vSwitch
	}
	return nil
}

func flattenBondMode(pr *import1.BondModeType) string {
	const two, three, four, five = 2, 3, 4, 5
	if pr != nil {
		if *pr == import1.BondModeType(two) {
			return "ACTIVE_BACKUP"
		}
		if *pr == import1.BondModeType(three) {
			return "BALANCE_SLB"
		}
		if *pr == import1.BondModeType(four) {
			return "BALANCE_TCP"
		}
		if *pr == import1.BondModeType(five) {
			return "NONE"
		}
	}
	return "UNKNOWN"
}

func flattenMetadata(pr *config.Metadata) []map[string]interface{} {
	if pr != nil {
		meta := make([]map[string]interface{}, 0)

		m := make(map[string]interface{})

		m["owner_reference_id"] = utils.StringValue(pr.OwnerReferenceId)
		m["owner_user_name"] = utils.StringValue(pr.OwnerUserName)
		m["project_reference_id"] = utils.StringValue(pr.ProjectReferenceId)
		m["project_name"] = utils.StringValue(pr.ProjectName)
		m["category_ids"] = pr.CategoryIds

		meta = append(meta, m)
		return meta
	}
	return nil
}

func flattenClusters(pr []import1.Cluster) []map[string]interface{} {
	if len(pr) > 0 {
		clsList := make([]map[string]interface{}, 0)

		for k, v := range pr {
			cls := make(map[string]interface{})

			cls["ext_id"] = v.ExtId
			cls["hosts"] = flattenHosts(v.Hosts)
			cls["gateway_ip_address"] = flattenIPv4(v.GatewayIpAddress)

			clsList[k] = cls
			return clsList
		}
	}
	return nil
}

func flattenHosts(pr []import1.Host) []map[string]interface{} {
	if len(pr) > 0 {
		hosts := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			h := make(map[string]interface{})

			h["ext_id"] = v.ExtId
			h["internal_bridge_name"] = v.InternalBridgeName
			h["host_nics"] = v.HostNics
			h["ip_address"] = flattenIPv4Subnet(v.IpAddress)
			h["route_table"] = v.RouteTable

			hosts[k] = h
		}
		return hosts
	}
	return nil
}

func flattenVPC(pr *import1.Vpc) []map[string]interface{} {
	if pr != nil {
		vpcs := make([]map[string]interface{}, 0)

		vpc := make(map[string]interface{})

		if pr.TenantId != nil {
			vpc["tenant_id"] = pr.TenantId
		}
		vpc["ext_id"] = pr.ExtId
		vpc["links"] = flattenLinks(pr.Links)
		vpc["metadata"] = flattenMetadata(pr.Metadata)
		vpc["name"] = pr.Name
		vpc["description"] = pr.Description
		vpc["common_dhcp_options"] = flattenCommonDhcpOptions(pr.CommonDhcpOptions)
		vpc["snat_ips"] = flattenNtpServer(pr.SnatIps)
		vpc["external_subnets"] = flattenExternalSubnets(pr.ExternalSubnets)
		vpc["external_routing_domain_reference"] = pr.ExternalRoutingDomainReference
		vpc["externally_routable_prefixes"] = flattenExternallyRoutablePrefixes(pr.ExternallyRoutablePrefixes)

		vpcs = append(vpcs, vpc)
		return vpcs
	}
	return nil
}

func flattenCommonDhcpOptions(pr *import1.VpcDhcpOptions) []map[string]interface{} {
	if pr != nil {
		dhcp := make([]map[string]interface{}, 0)
		d := make(map[string]interface{})

		d["domain_name_servers"] = flattenNtpServer(pr.DomainNameServers)

		dhcp = append(dhcp, d)
		return dhcp
	}
	return nil
}

func flattenExternalSubnets(pr []import1.ExternalSubnet) []map[string]interface{} {
	if len(pr) > 0 {
		extSubs := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			sub := make(map[string]interface{})
			sub["subnet_reference"] = v.SubnetReference
			sub["external_ips"] = flattenNtpServer(v.ExternalIps)
			sub["gateway_nodes"] = v.GatewayNodes
			sub["active_gateway_node"] = flattenActiveGatewayNode(v.ActiveGatewayNodes)
			sub["active_gateway_count"] = v.ActiveGatewayCount

			extSubs[k] = sub
		}
		return extSubs
	}
	return nil
}

func flattenActiveGatewayNode(pr []import1.GatewayNodeReference) []map[string]interface{} {
	if len(pr) > 0 {
		nodes := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			node := make(map[string]interface{})

			node["node_id"] = v.NodeId
			node["node_ip_address"] = flattenNodeIPAddress(v.NodeIpAddress)

			nodes[k] = node
		}
		// n := make(map[string]interface{})
		//
		//n["node_id"] =
		//n["node_ip_address"] = flattenNodeIPAddress(pr.NodeIpAddress)
		//
		//nodes = append(nodes, n)
		return nodes
	}
	return nil
}

func flattenNodeIPAddress(pr *config.IPAddress) []map[string]interface{} {
	if pr != nil {
		ips := make([]map[string]interface{}, 0)
		ip := make(map[string]interface{})

		ip["ipv4"] = flattenIPv4(pr.Ipv4)
		ip["ipv6"] = flattenIPv6(pr.Ipv6)

		ips = append(ips, ip)
		return ips
	}
	return nil
}

func flattenExternallyRoutablePrefixes(pr []import1.IPSubnet) []map[string]interface{} {
	if len(pr) > 0 {
		exts := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ext := make(map[string]interface{})

			ext["ipv4"] = flattenIPv4Subnet(v.Ipv4)
			ext["ipv6"] = flattenIPv6Subnet(v.Ipv6)

			exts[k] = ext
		}
		return exts
	}
	return nil
}

func flattenIPUsage(pr *import1.IPUsage) []map[string]interface{} {
	if pr != nil {
		usage := make([]map[string]interface{}, 0)

		ip := make(map[string]interface{})

		ip["num_macs"] = pr.NumMacs
		ip["num_free_ips"] = pr.NumFreeIPs
		ip["num_assigned_ips"] = pr.NumAssignedIPs
		ip["ip_pool_usages"] = flattenIPPoolUsages(pr.IpPoolUsages)

		usage = append(usage, ip)
		return usage
	}
	return nil
}

func flattenIPPoolUsages(pr []import1.IPPoolUsage) []map[string]interface{} {
	if len(pr) > 0 {
		ips := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ip := make(map[string]interface{})

			ip["num_free_ips"] = v.NumFreeIPs
			ip["num_total_ips"] = v.NumTotalIPs
			ip["range"] = flattenIPv4Pool(v.Range)

			ips[k] = ip
		}
		return ips
	}
	return nil
}

func flattenIPv4Pool(pr *import1.IPv4Pool) []map[string]interface{} {
	if pr != nil {
		pool := make([]map[string]interface{}, 0)

		ip := make(map[string]interface{})

		ip["start_ip"] = flattenIPv4(pr.StartIp)
		ip["end_ip"] = flattenIPv4(pr.EndIp)

		pool = append(pool, ip)
		return pool
	}
	return nil
}

func flattenMigrationState(pr *import1.MigrationState) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == import1.MigrationState(two) {
			return "IN_PROGRESS"
		}
		if *pr == import1.MigrationState(three) {
			return "FAILED"
		}
	}
	return "UNKNOWN"
}
