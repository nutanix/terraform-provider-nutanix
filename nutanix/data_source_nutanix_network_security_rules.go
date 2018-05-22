package nutanix

import (
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixNetworkSecurityRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixNetworkSecurityRulesRead,

		Schema: getDataSourceNetworkSecurityRulesSchema(),
	}
}

func dataSourceNutanixNetworkSecurityRulesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*NutanixClient).API

	metadata := &v3.ListMetadata{}

	if v, ok := d.GetOk("metadata"); ok {
		m := v.(map[string]interface{})
		metadata.Kind = utils.String("network_security_rule")
		if mv, mok := m["sort_attribute"]; mok {
			metadata.SortAttribute = utils.String(mv.(string))
		}
		if mv, mok := m["filter"]; mok {
			metadata.Filter = utils.String(mv.(string))
		}
		if mv, mok := m["length"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Length = utils.Int64(int64(i))
		}
		if mv, mok := m["sort_order"]; mok {
			metadata.SortOrder = utils.String(mv.(string))
		}
		if mv, mok := m["offset"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Offset = utils.Int64(int64(i))
		}
	}

	// Make request to the API
	resp, err := conn.V3.ListNetworkSecurityRule(metadata)
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})
		// set metadata values
		metadata := make(map[string]interface{})
		metadata["last_update_time"] = utils.TimeValue(v.Metadata.LastUpdateTime).String()
		metadata["kind"] = utils.StringValue(v.Metadata.Kind)
		metadata["uuid"] = utils.StringValue(v.Metadata.UUID)
		metadata["creation_time"] = utils.TimeValue(v.Metadata.CreationTime).String()
		metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.Metadata.SpecVersion)))
		metadata["spec_hash"] = utils.StringValue(v.Metadata.SpecHash)
		metadata["name"] = utils.StringValue(v.Metadata.Name)
		entity["metadata"] = metadata

		if v.Metadata.Categories != nil {
			categories := v.Metadata.Categories
			var catList []map[string]interface{}

			for name, values := range categories {
				catItem := make(map[string]interface{})
				catItem["name"] = name
				catItem["value"] = values
				catList = append(catList, catItem)
			}
			entity["categories"] = catList
		}

		entity["api_version"] = utils.StringValue(v.APIVersion)

		pr := make(map[string]interface{})
		if v.Metadata.ProjectReference != nil {
			pr["kind"] = utils.StringValue(v.Metadata.ProjectReference.Kind)
			pr["name"] = utils.StringValue(v.Metadata.ProjectReference.Name)
			pr["uuid"] = utils.StringValue(v.Metadata.ProjectReference.UUID)
		}
		entity["project_reference"] = pr

		or := make(map[string]interface{})
		if v.Metadata.OwnerReference != nil {
			or["kind"] = utils.StringValue(v.Metadata.OwnerReference.Kind)
			or["name"] = utils.StringValue(v.Metadata.OwnerReference.Name)
			or["uuid"] = utils.StringValue(v.Metadata.OwnerReference.UUID)
		}
		entity["owner_reference"] = or

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
						if oa.Filter.KindList != nil {
							fkl := oa.Filter.KindList
							fkList := make([]string, len(fkl))
							for i, f := range fkl {
								fkList[i] = utils.StringValue(f)
							}
							qroaItem["filter_kind_list"] = fkList
						}

						qroaItem["filter_type"] = utils.StringValue(oa.Filter.Type)

						if oa.Filter.Params != nil {
							fp := oa.Filter.Params
							var fpList []map[string]interface{}

							for name, values := range fp {
								fpItem := make(map[string]interface{})
								fpItem["name"] = name
								fpItem["values"] = values
								fpList = append(fpList, fpItem)
							}
							qroaItem["filter_params"] = fpList
						}

					}

					qroaItem["peer_specification_type"] = utils.StringValue(oa.PeerSpecificationType)
					qroaItem["expiration_time"] = utils.StringValue(oa.ExpirationTime)

					// set network_function_chain_reference
					if oa.NetworkFunctionChainReference != nil {
						nfcr := make(map[string]interface{})
						nfcr["kind"] = utils.StringValue(oa.NetworkFunctionChainReference.Kind)
						nfcr["name"] = utils.StringValue(oa.NetworkFunctionChainReference.Name)
						nfcr["uuid"] = utils.StringValue(oa.NetworkFunctionChainReference.UUID)
						qroaItem["network_function_chain_reference"] = nfcr
					}

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

				//Set quarantine_rule_outbound_allow_list
				entity["quarantine_rule_outbound_allow_list"] = qroaList
			}

			if v.Spec.Resources.QuarantineRule.TargetGroup != nil {
				tg := v.Spec.Resources.QuarantineRule.TargetGroup
				entity["quarantine_rule_target_group_default_internal_policy"] = utils.StringValue(tg.DefaultInternalPolicy)
				entity["quarantine_rule_target_group_peer_specification_type"] = utils.StringValue(tg.PeerSpecificationType)

				if tg.Filter != nil {
					if tg.Filter.KindList != nil {
						fkl := tg.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						entity["quarantine_rule_target_group_filter_kind_list"] = fkList
					}

					entity["quarantine_rule_target_group_filter_type"] = utils.StringValue(tg.Filter.Type)

					if tg.Filter.Params != nil {
						fp := tg.Filter.Params
						var fpList []map[string]interface{}

						for name, values := range fp {
							fpItem := make(map[string]interface{})
							fpItem["name"] = name
							fpItem["values"] = values
							fpList = append(fpList, fpItem)
						}

						entity["quarantine_rule_target_group_filter_params"] = fpList
					}

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
						if ia.Filter.KindList != nil {
							fkl := ia.Filter.KindList
							fkList := make([]string, len(fkl))
							for i, f := range fkl {
								fkList[i] = utils.StringValue(f)
							}
							qriaItem["filter_kind_list"] = fkList
						}

						qriaItem["filter_type"] = utils.StringValue(ia.Filter.Type)

						if ia.Filter.Params != nil {
							fp := ia.Filter.Params
							var fpList []map[string]interface{}

							for name, values := range fp {
								fpItem := make(map[string]interface{})
								fpItem["name"] = name
								fpItem["values"] = values
								fpList = append(fpList, fpItem)
							}
							qriaItem["filter_params"] = fpList
						}

					}

					qriaItem["peer_specification_type"] = utils.StringValue(ia.PeerSpecificationType)
					qriaItem["expiration_time"] = utils.StringValue(ia.ExpirationTime)

					// set network_function_chain_reference
					if ia.NetworkFunctionChainReference != nil {
						nfcr := make(map[string]interface{})
						nfcr["kind"] = utils.StringValue(ia.NetworkFunctionChainReference.Kind)
						nfcr["name"] = utils.StringValue(ia.NetworkFunctionChainReference.Name)
						nfcr["uuid"] = utils.StringValue(ia.NetworkFunctionChainReference.UUID)
						qriaItem["network_function_chain_reference"] = nfcr
					}

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
				//Set quarantine_rule_inbound_allow_list
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
						if oa.Filter.KindList != nil {
							fkl := oa.Filter.KindList
							fkList := make([]string, len(fkl))
							for i, f := range fkl {
								fkList[i] = utils.StringValue(f)
							}
							aroaItem["filter_kind_list"] = fkList
						}

						aroaItem["filter_type"] = utils.StringValue(oa.Filter.Type)

						if oa.Filter.Params != nil {
							fp := oa.Filter.Params
							var fpList []map[string]interface{}

							for name, values := range fp {
								fpItem := make(map[string]interface{})
								fpItem["name"] = name
								fpItem["values"] = values
								fpList = append(fpList, fpItem)
							}
							aroaItem["filter_params"] = fpList
						}

					}

					aroaItem["peer_specification_type"] = utils.StringValue(oa.PeerSpecificationType)
					aroaItem["expiration_time"] = utils.StringValue(oa.ExpirationTime)

					// set network_function_chain_reference
					if oa.NetworkFunctionChainReference != nil {
						nfcr := make(map[string]interface{})
						nfcr["kind"] = utils.StringValue(oa.NetworkFunctionChainReference.Kind)
						nfcr["name"] = utils.StringValue(oa.NetworkFunctionChainReference.Name)
						nfcr["uuid"] = utils.StringValue(oa.NetworkFunctionChainReference.UUID)
						aroaItem["network_function_chain_reference"] = nfcr
					}

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

				//Set app_rule_outbound_allow_list
				entity["app_rule_outbound_allow_list"] = aroaList
			}

			if tg := v.Spec.Resources.AppRule.TargetGroup; tg != nil {
				entity["app_rule_target_group_default_internal_policy"] = utils.StringValue(tg.DefaultInternalPolicy)
				entity["app_rule_target_group_peer_specification_type"] = utils.StringValue(tg.PeerSpecificationType)

				if tg.Filter != nil {
					if tg.Filter.KindList != nil {
						fkl := tg.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						entity["app_rule_target_group_filter_kind_list"] = fkList
					}

					entity["app_rule_target_group_filter_type"] = utils.StringValue(tg.Filter.Type)

					if tg.Filter.Params != nil {
						fp := tg.Filter.Params
						var fpList []map[string]interface{}

						for name, values := range fp {
							fpItem := make(map[string]interface{})
							fpItem["name"] = name
							fpItem["values"] = values
							fpList = append(fpList, fpItem)
						}

						entity["app_rule_target_group_filter_params"] = fpList
					}

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
						if ia.Filter.KindList != nil {
							fkl := ia.Filter.KindList
							fkList := make([]string, len(fkl))
							for i, f := range fkl {
								fkList[i] = utils.StringValue(f)
							}
							ariaItem["filter_kind_list"] = fkList
						}

						ariaItem["filter_type"] = utils.StringValue(ia.Filter.Type)

						if ia.Filter.Params != nil {
							fp := ia.Filter.Params
							var fpList []map[string]interface{}

							for name, values := range fp {
								fpItem := make(map[string]interface{})
								fpItem["name"] = name
								fpItem["values"] = values
								fpList = append(fpList, fpItem)
							}
							ariaItem["filter_params"] = fpList
						}

					}

					ariaItem["peer_specification_type"] = utils.StringValue(ia.PeerSpecificationType)
					ariaItem["expiration_time"] = utils.StringValue(ia.ExpirationTime)

					// set network_function_chain_reference
					if ia.NetworkFunctionChainReference != nil {
						nfcr := make(map[string]interface{})
						nfcr["kind"] = utils.StringValue(ia.NetworkFunctionChainReference.Kind)
						nfcr["name"] = utils.StringValue(ia.NetworkFunctionChainReference.Name)
						nfcr["uuid"] = utils.StringValue(ia.NetworkFunctionChainReference.UUID)
						ariaItem["network_function_chain_reference"] = nfcr
					}

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

				//Set app_rule_inbound_allow_list
				entity["app_rule_inbound_allow_list"] = ariaList
			}

		} else {
			entity["app_rule_action"] = ""
		}

		if v.Spec.Resources.IsolationRule != nil {
			entity["isolation_rule_action"] = utils.StringValue(v.Spec.Resources.IsolationRule.Action)

			if firstFilter := v.Spec.Resources.IsolationRule.FirstEntityFilter; firstFilter != nil {
				if firstFilter.KindList != nil {
					fkl := firstFilter.KindList
					fkList := make([]string, len(fkl))
					for i, f := range fkl {
						fkList[i] = utils.StringValue(f)
					}
					entity["isolation_rule_first_entity_filter_kind_list"] = fkList
				} else {
					entity["isolation_rule_first_entity_filter_kind_list"] = make([]string, 0)
				}

				entity["isolation_rule_first_entity_filter_type"] = utils.StringValue(firstFilter.Type)

				if fp := firstFilter.Params; fp != nil {
					var fpList []map[string]interface{}

					for name, values := range fp {
						fpItem := make(map[string]interface{})
						fpItem["name"] = name
						fpItem["values"] = values
						fpList = append(fpList, fpItem)
					}

					entity["isolation_rule_first_entity_filter_params"] = fpList
				}

			}

			if secondFilter := v.Spec.Resources.IsolationRule.SecondEntityFilter; secondFilter != nil {
				if secondFilter.KindList != nil {
					fkl := secondFilter.KindList
					fkList := make([]string, len(fkl))
					for i, f := range fkl {
						fkList[i] = utils.StringValue(f)
					}
					entity["isolation_rule_second_entity_filter_kind_list"] = fkList
				}

				entity["isolation_rule_second_entity_filter_type"] = utils.StringValue(secondFilter.Type)

				if secondFilter.Params != nil {
					fp := secondFilter.Params
					var fpList []map[string]interface{}

					for name, values := range fp {
						fpItem := make(map[string]interface{})
						fpItem["name"] = name
						fpItem["values"] = values
						fpList = append(fpList, fpItem)
					}

					entity["isolation_rule_second_entity_filter_params"] = fpList
				}

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
	d.SetId(resource.UniqueId())

	return nil
}

func getDataSourceNetworkSecurityRulesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
					"categories": {
						Type: schema.TypeList,

						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
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
	}
}
