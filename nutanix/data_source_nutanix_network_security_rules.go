package nutanix

import (
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixNetworkSecurityRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixNetworkSecurityRulesRead,

		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sort_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"filter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"length": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"offset": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_security_rule_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"api_version": {
							Type: schema.TypeString,

							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_update_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"categories": categoriesSchema(),
						"owner_reference": {
							Type: schema.TypeMap,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type: schema.TypeString,
									},
									"uuid": {
										Type: schema.TypeString,
									},
									"name": {
										Type: schema.TypeString,
									},
								},
							},
						},
						"project_reference": {
							Type: schema.TypeMap,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type: schema.TypeString,
									},
									"uuid": {
										Type: schema.TypeString,
									},
									"name": {
										Type: schema.TypeString,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type: schema.TypeString,

							Computed: true,
						},
						"quarantine_rule_action": {
							Type: schema.TypeString,

							Computed: true,
						},
						"quarantine_rule_outbound_allow_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type: schema.TypeString,

										Computed: true,
									},
									"ip_subnet": {
										Type: schema.TypeString,

										Computed: true,
									},
									"ip_subnet_prefix_length": {
										Type: schema.TypeString,

										Computed: true,
									},
									"tcp_port_range_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_port": {
													Type: schema.TypeString,

													Computed: true,
												},
												"start_port": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"udp_port_range_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_port": {
													Type: schema.TypeInt,

													Computed: true,
												},
												"start_port": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"filter_kind_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"filter_type": {
										Type: schema.TypeString,

										Computed: true,
									},
									"filter_params": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"values": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"peer_specification_type": {
										Type: schema.TypeString,

										Computed: true,
									},

									"expiration_time": {
										Type: schema.TypeString,

										Computed: true,
									},
									"network_function_chain_reference": {
										Type: schema.TypeMap,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": {
													Type:     schema.TypeString,
													Required: true,
												},
												"uuid": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"icmp_type_code_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"code": {
													Type: schema.TypeString,

													Computed: true,
												},
												"type": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"quarantine_rule_target_group_default_internal_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quarantine_rule_target_group_peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quarantine_rule_target_group_filter_kind_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"quarantine_rule_target_group_filter_type": {
							Type: schema.TypeString,

							Computed: true,
						},
						"quarantine_rule_target_group_filter_params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"quarantine_rule_inbound_allow_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type: schema.TypeString,

										Computed: true,
									},
									"ip_subnet": {
										Type: schema.TypeString,

										Computed: true,
									},
									"ip_subnet_prefix_length": {
										Type: schema.TypeString,

										Computed: true,
									},
									"tcp_port_range_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_port": {
													Type: schema.TypeString,

													Computed: true,
												},
												"start_port": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"udp_port_range_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_port": {
													Type: schema.TypeInt,

													Computed: true,
												},
												"start_port": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"filter_kind_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"filter_type": {
										Type: schema.TypeString,

										Computed: true,
									},
									"filter_params": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"values": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"peer_specification_type": {
										Type: schema.TypeString,

										Computed: true,
									},

									"expiration_time": {
										Type: schema.TypeString,

										Computed: true,
									},
									"network_function_chain_reference": {
										Type: schema.TypeMap,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": {
													Type:     schema.TypeString,
													Required: true,
												},
												"uuid": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"icmp_type_code_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"code": {
													Type: schema.TypeString,

													Computed: true,
												},
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"app_rule_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_rule_outbound_allow_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type: schema.TypeString,

										Computed: true,
									},
									"ip_subnet": {
										Type: schema.TypeString,

										Computed: true,
									},
									"ip_subnet_prefix_length": {
										Type: schema.TypeString,

										Computed: true,
									},
									"tcp_port_range_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_port": {
													Type: schema.TypeInt,

													Computed: true,
												},
												"start_port": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"udp_port_range_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_port": {
													Type: schema.TypeInt,

													Computed: true,
												},
												"start_port": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"filter_kind_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"filter_type": {
										Type: schema.TypeString,

										Computed: true,
									},
									"filter_params": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"values": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"peer_specification_type": {
										Type: schema.TypeString,

										Computed: true,
									},

									"expiration_time": {
										Type: schema.TypeString,

										Computed: true,
									},
									"network_function_chain_reference": {
										Type: schema.TypeMap,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": {
													Type:     schema.TypeString,
													Required: true,
												},
												"uuid": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"icmp_type_code_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"code": {
													Type: schema.TypeString,

													Computed: true,
												},
												"type": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"app_rule_target_group_default_internal_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_rule_target_group_peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_rule_target_group_filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"app_rule_target_group_filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_rule_target_group_filter_params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"app_rule_inbound_allow_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type: schema.TypeString,

										Computed: true,
									},
									"ip_subnet": {
										Type: schema.TypeString,

										Computed: true,
									},
									"ip_subnet_prefix_length": {
										Type: schema.TypeString,

										Computed: true,
									},
									"tcp_port_range_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_port": {
													Type: schema.TypeString,

													Computed: true,
												},
												"start_port": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"udp_port_range_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_port": {
													Type: schema.TypeInt,

													Computed: true,
												},
												"start_port": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"filter_kind_list": {
										Type: schema.TypeList,

										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"filter_type": {
										Type: schema.TypeString,

										Computed: true,
									},
									"filter_params": {
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"values": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"peer_specification_type": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"expiration_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_function_chain_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": {
													Type:     schema.TypeString,
													Required: true,
												},
												"uuid": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
									"icmp_type_code_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"code": {
													Type: schema.TypeString,

													Computed: true,
												},
												"type": {
													Type: schema.TypeString,

													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"isolation_rule_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isolation_rule_first_entity_filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"isolation_rule_first_entity_filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isolation_rule_first_entity_filter_params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"isolation_rule_second_entity_filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"isolation_rule_second_entity_filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isolation_rule_second_entity_filter_params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceNutanixNetworkSecurityRulesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	resp, err := conn.V3.ListAllNetworkSecurityRule()
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})

		m, c := setRSEntityMetadata(v.Metadata)

		entity["metadata"] = m
		entity["project_reference"] = flattenReferenceValues(v.Metadata.ProjectReference)
		entity["owner_reference"] = flattenReferenceValues(v.Metadata.OwnerReference)
		entity["categories"] = c
		entity["api_version"] = utils.StringValue(v.APIVersion)
		entity["name"] = utils.StringValue(v.Spec.Name)
		entity["description"] = utils.StringValue(v.Spec.Description)

		if v.Spec.Resources.QuarantineRule != nil {
			entity["quarantine_rule_action"] = utils.StringValue(v.Spec.Resources.QuarantineRule.Action)

			if v.Spec.Resources.QuarantineRule.OutboundAllowList != nil {
				oal := v.Spec.Resources.QuarantineRule.OutboundAllowList
				qroaList := make([]map[string]interface{}, len(oal))
				for k, oa := range oal {
					qroaItem := make(map[string]interface{})
					qroaItem["protocol"] = utils.StringValue(oa.Protocol)

					if oa.IPSubnet != nil {
						qroaItem["ip_subnet"] = utils.StringValue(oa.IPSubnet.IP)
						qroaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(oa.IPSubnet.PrefixLength), 10)
					}

					if oa.TCPPortRangeList != nil {
						tcpprl := oa.TCPPortRangeList
						tcpprList := make([]map[string]interface{}, len(tcpprl))
						for i, tcp := range tcpprl {
							tcpItem := make(map[string]interface{})
							tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
							tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
							tcpprList[i] = tcpItem
						}
						qroaItem["tcp_port_range_list"] = tcpprList
					}

					if oa.UDPPortRangeList != nil {
						udpprl := oa.UDPPortRangeList
						udpprList := make([]map[string]interface{}, len(udpprl))
						for i, udp := range udpprl {
							udpItem := make(map[string]interface{})
							udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
							udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
							udpprList[i] = udpItem
						}
						qroaItem["udp_port_range_list"] = udpprList
					}

					if oa.Filter != nil {
						qroaItem["filter_kind_list"] = utils.StringValueSlice(oa.Filter.KindList)
						qroaItem["filter_type"] = utils.StringValue(oa.Filter.Type)
						qroaItem["filter_params"] = expandFilterParams(oa.Filter.Params)

					}

					qroaItem["peer_specification_type"] = utils.StringValue(oa.PeerSpecificationType)
					qroaItem["expiration_time"] = utils.StringValue(oa.ExpirationTime)
					qroaItem["network_function_chain_reference"] = flattenReferenceValues(oa.NetworkFunctionChainReference)

					if oa.IcmpTypeCodeList != nil {
						icmptcl := oa.IcmpTypeCodeList
						icmptcList := make([]map[string]interface{}, len(icmptcl))
						for i, icmp := range icmptcl {
							icmpItem := make(map[string]interface{})
							icmpItem["end_port"] = strconv.FormatInt(utils.Int64Value(icmp.Code), 10)
							icmpItem["start_port"] = strconv.FormatInt(utils.Int64Value(icmp.Type), 10)
							icmptcList[i] = icmpItem
						}
						qroaItem["icmp_type_code_list"] = icmptcList
					}

					qroaList[k] = qroaItem
				}

				// Set quarantine_rule_outbound_allow_list
				entity["quarantine_rule_outbound_allow_list"] = qroaList
			}

			if v.Spec.Resources.QuarantineRule.TargetGroup != nil {
				tg := v.Spec.Resources.QuarantineRule.TargetGroup
				entity["quarantine_rule_target_group_default_internal_policy"] = utils.StringValue(tg.DefaultInternalPolicy)
				entity["quarantine_rule_target_group_peer_specification_type"] = utils.StringValue(tg.PeerSpecificationType)

				if tg.Filter != nil {
					entity["quarantine_rule_target_group_filter_kind_list"] = utils.StringValueSlice(tg.Filter.KindList)
					entity["quarantine_rule_target_group_filter_type"] = utils.StringValue(tg.Filter.Type)
					entity["quarantine_rule_target_group_filter_params"] = expandFilterParams(tg.Filter.Params)
				}

			}

			if v.Spec.Resources.QuarantineRule.InboundAllowList != nil {
				ial := v.Spec.Resources.QuarantineRule.InboundAllowList
				qriaList := make([]map[string]interface{}, len(ial))
				for k, ia := range ial {
					qriaItem := make(map[string]interface{})
					qriaItem["protocol"] = utils.StringValue(ia.Protocol)

					if ia.IPSubnet != nil {
						qriaItem["ip_subnet"] = utils.StringValue(ia.IPSubnet.IP)
						qriaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(ia.IPSubnet.PrefixLength), 10)
					}

					if ia.TCPPortRangeList != nil {
						tcpprl := ia.TCPPortRangeList
						tcpprList := make([]map[string]interface{}, len(tcpprl))
						for i, tcp := range tcpprl {
							tcpItem := make(map[string]interface{})
							tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
							tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
							tcpprList[i] = tcpItem
						}
						qriaItem["tcp_port_range_list"] = tcpprList
					}

					if ia.UDPPortRangeList != nil {
						udpprl := ia.UDPPortRangeList
						udpprList := make([]map[string]interface{}, len(udpprl))
						for i, udp := range udpprl {
							udpItem := make(map[string]interface{})
							udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
							udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
							udpprList[i] = udpItem
						}
						qriaItem["udp_port_range_list"] = udpprList
					}

					if ia.Filter != nil {
						qriaItem["filter_kind_list"] = utils.StringValueSlice(ia.Filter.KindList)
						qriaItem["filter_type"] = utils.StringValue(ia.Filter.Type)
						qriaItem["filter_params"] = expandFilterParams(ia.Filter.Params)
					}

					qriaItem["peer_specification_type"] = utils.StringValue(ia.PeerSpecificationType)
					qriaItem["expiration_time"] = utils.StringValue(ia.ExpirationTime)
					qriaItem["network_function_chain_reference"] = flattenReferenceValues(ia.NetworkFunctionChainReference)

					if ia.IcmpTypeCodeList != nil {
						icmptcl := ia.IcmpTypeCodeList
						icmptcList := make([]map[string]interface{}, len(icmptcl))
						for i, icmp := range icmptcl {
							icmpItem := make(map[string]interface{})
							icmpItem["end_port"] = strconv.FormatInt(utils.Int64Value(icmp.Code), 10)
							icmpItem["start_port"] = strconv.FormatInt(utils.Int64Value(icmp.Type), 10)
							icmptcList[i] = icmpItem
						}
						qriaItem["icmp_type_code_list"] = icmptcList
					}

					qriaList[k] = qriaItem
				}
				// Set quarantine_rule_inbound_allow_list
				entity["quarantine_rule_inbound_allow_list"] = qriaList
			}
		} else {
			entity["quarantine_rule_inbound_allow_list"] = make([]string, 0)
			entity["quarantine_rule_outbound_allow_list"] = make([]string, 0)
			entity["quarantine_rule_target_group_filter_kind_list"] = make([]string, 0)
			entity["quarantine_rule_target_group_filter_params"] = make([]string, 0)
		}

		if v.Spec.Resources.AppRule != nil {
			entity["app_rule_action"] = utils.StringValue(v.Spec.Resources.AppRule.Action)

			if oal := v.Spec.Resources.AppRule.OutboundAllowList; oal != nil {
				aroaList := make([]map[string]interface{}, len(oal))
				for k, oa := range oal {
					aroaItem := make(map[string]interface{})
					aroaItem["protocol"] = utils.StringValue(oa.Protocol)

					if oa.IPSubnet != nil {
						aroaItem["ip_subnet"] = utils.StringValue(oa.IPSubnet.IP)
						aroaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(oa.IPSubnet.PrefixLength), 10)
					}

					if oa.TCPPortRangeList != nil {
						tcpprl := oa.TCPPortRangeList
						tcpprList := make([]map[string]interface{}, len(tcpprl))
						for i, tcp := range tcpprl {
							tcpItem := make(map[string]interface{})
							tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
							tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
							tcpprList[i] = tcpItem
						}
						aroaItem["tcp_port_range_list"] = tcpprList
					}

					if oa.UDPPortRangeList != nil {
						udpprl := oa.UDPPortRangeList
						udpprList := make([]map[string]interface{}, len(udpprl))
						for i, udp := range udpprl {
							udpItem := make(map[string]interface{})
							udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
							udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
							udpprList[i] = udpItem
						}
						aroaItem["udp_port_range_list"] = udpprList
					}

					if oa.Filter != nil {
						aroaItem["filter_kind_list"] = utils.StringValueSlice(oa.Filter.KindList)
						aroaItem["filter_type"] = utils.StringValue(oa.Filter.Type)
						aroaItem["filter_params"] = expandFilterParams(oa.Filter.Params)
					}

					aroaItem["peer_specification_type"] = utils.StringValue(oa.PeerSpecificationType)
					aroaItem["expiration_time"] = utils.StringValue(oa.ExpirationTime)
					aroaItem["network_function_chain_reference"] = flattenReferenceValues(oa.NetworkFunctionChainReference)

					if oa.IcmpTypeCodeList != nil {
						icmptcl := oa.IcmpTypeCodeList
						icmptcList := make([]map[string]interface{}, len(icmptcl))
						for i, icmp := range icmptcl {
							icmpItem := make(map[string]interface{})
							icmpItem["end_port"] = strconv.FormatInt(utils.Int64Value(icmp.Code), 10)
							icmpItem["start_port"] = strconv.FormatInt(utils.Int64Value(icmp.Type), 10)
							icmptcList[i] = icmpItem
						}
						aroaItem["icmp_type_code_list"] = icmptcList
					}

					aroaList[k] = aroaItem
				}

				// Set app_rule_outbound_allow_list
				entity["app_rule_outbound_allow_list"] = aroaList
			}

			if tg := v.Spec.Resources.AppRule.TargetGroup; tg != nil {
				entity["app_rule_target_group_default_internal_policy"] = utils.StringValue(tg.DefaultInternalPolicy)
				entity["app_rule_target_group_peer_specification_type"] = utils.StringValue(tg.PeerSpecificationType)

				if tg.Filter != nil {
					entity["app_rule_target_group_filter_kind_list"] = utils.StringValueSlice(tg.Filter.KindList)
					entity["app_rule_target_group_filter_type"] = utils.StringValue(tg.Filter.Type)
					entity["app_rule_target_group_filter_params"] = expandFilterParams(tg.Filter.Params)
				}
			}

			if ial := v.Spec.Resources.AppRule.InboundAllowList; ial != nil {
				ariaList := make([]map[string]interface{}, len(ial))
				for k, ia := range ial {
					ariaItem := make(map[string]interface{})
					ariaItem["protocol"] = utils.StringValue(ia.Protocol)

					if ia.IPSubnet != nil {
						ariaItem["ip_subnet"] = utils.StringValue(ia.IPSubnet.IP)
						ariaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(ia.IPSubnet.PrefixLength), 10)
					}

					if ia.TCPPortRangeList != nil {
						tcpprl := ia.TCPPortRangeList
						tcpprList := make([]map[string]interface{}, len(tcpprl))
						for i, tcp := range tcpprl {
							tcpItem := make(map[string]interface{})
							tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
							tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
							tcpprList[i] = tcpItem
						}
						ariaItem["tcp_port_range_list"] = tcpprList
					}

					if ia.UDPPortRangeList != nil {
						udpprl := ia.UDPPortRangeList
						udpprList := make([]map[string]interface{}, len(udpprl))
						for i, udp := range udpprl {
							udpItem := make(map[string]interface{})
							udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
							udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
							udpprList[i] = udpItem
						}
						ariaItem["udp_port_range_list"] = udpprList
					}

					if ia.Filter != nil {
						ariaItem["filter_kind_list"] = utils.StringValueSlice(ia.Filter.KindList)
						ariaItem["filter_type"] = utils.StringValue(ia.Filter.Type)
						ariaItem["filter_params"] = expandFilterParams(ia.Filter.Params)
					}

					ariaItem["peer_specification_type"] = utils.StringValue(ia.PeerSpecificationType)
					ariaItem["expiration_time"] = utils.StringValue(ia.ExpirationTime)
					ariaItem["network_function_chain_reference"] = flattenReferenceValues(ia.NetworkFunctionChainReference)

					if icmptcl := ia.IcmpTypeCodeList; icmptcl != nil {
						icmptcList := make([]map[string]interface{}, len(icmptcl))
						for i, icmp := range icmptcl {
							icmpItem := make(map[string]interface{})
							icmpItem["end_port"] = strconv.FormatInt(utils.Int64Value(icmp.Code), 10)
							icmpItem["start_port"] = strconv.FormatInt(utils.Int64Value(icmp.Type), 10)
							icmptcList[i] = icmpItem
						}
						ariaItem["icmp_type_code_list"] = icmptcList
					}

					ariaList[k] = ariaItem
				}

				// Set app_rule_inbound_allow_list
				entity["app_rule_inbound_allow_list"] = ariaList
			}

		} else {
			entity["app_rule_action"] = ""
		}

		if v.Spec.Resources.IsolationRule != nil {
			entity["isolation_rule_action"] = utils.StringValue(v.Spec.Resources.IsolationRule.Action)

			if firstFilter := v.Spec.Resources.IsolationRule.FirstEntityFilter; firstFilter != nil {
				entity["isolation_rule_first_entity_filter_kind_list"] = utils.StringValueSlice(firstFilter.KindList)
				entity["isolation_rule_first_entity_filter_type"] = utils.StringValue(firstFilter.Type)
				entity["isolation_rule_first_entity_filter_params"] = expandFilterParams(firstFilter.Params)
			}

			if secondFilter := v.Spec.Resources.IsolationRule.SecondEntityFilter; secondFilter != nil {
				entity["isolation_rule_second_entity_filter_kind_list"] = utils.StringValueSlice(secondFilter.KindList)
				entity["isolation_rule_second_entity_filter_type"] = utils.StringValue(secondFilter.Type)
				entity["isolation_rule_second_entity_filter_params"] = expandFilterParams(secondFilter.Params)
			}

		} else {
			entity["isolation_rule_first_entity_filter_kind_list"] = make([]string, 0)
			entity["isolation_rule_first_entity_filter_params"] = make([]string, 0)
			entity["isolation_rule_second_entity_filter_kind_list"] = make([]string, 0)
			entity["isolation_rule_second_entity_filter_params"] = make([]string, 0)
		}

		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId(uuid.NewV4().String())

	return nil
}
