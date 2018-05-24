package nutanix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixNetworkSecurityRule() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNutanixNetworkSecurityRuleRead,
		Schema: getDataSourceNetworkSecurityRuleSchema(),
	}
}

func dataSourceNutanixNetworkSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Network Security Rule: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	networkSecurityRuleID, ok := d.GetOk("network_security_rule_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute network_security_rule_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetNetworkSecurityRule(networkSecurityRuleID.(string))

	if err != nil {
		return err
	}

	// set metadata values
	if resp.Metadata != nil {
		metadata := make(map[string]interface{})

		metadata["last_update_time"] = resp.Metadata.LastUpdateTime.String()
		metadata["uuid"] = utils.StringValue(resp.Metadata.UUID)
		metadata["creation_time"] = resp.Metadata.CreationTime.String()
		metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(resp.Metadata.SpecVersion)))
		metadata["spec_hash"] = utils.StringValue(resp.Metadata.SpecHash)
		metadata["name"] = utils.StringValue(resp.Metadata.Name)

		if err := d.Set("metadata", metadata); err != nil {
			return err
		}

		if resp.Metadata.Categories != nil {
			categories := resp.Metadata.Categories
			var catList []map[string]interface{}

			for name, values := range categories {
				catItem := make(map[string]interface{})
				catItem["name"] = name
				catItem["value"] = values
				catList = append(catList, catItem)
			}
			if err := d.Set("categories", catList); err != nil {
				return err
			}
		}

		or := make(map[string]interface{})
		if resp.Metadata.OwnerReference != nil {
			or["kind"] = utils.StringValue(resp.Metadata.OwnerReference.Kind)
			or["name"] = utils.StringValue(resp.Metadata.OwnerReference.Name)
			or["uuid"] = utils.StringValue(resp.Metadata.OwnerReference.UUID)
		}

		if err := d.Set("owner_reference", or); err != nil {
			return err
		}

		pr := make(map[string]interface{})
		if resp.Metadata.ProjectReference != nil {
			or["kind"] = utils.StringValue(resp.Metadata.ProjectReference.Kind)
			or["name"] = utils.StringValue(resp.Metadata.ProjectReference.Name)
			or["uuid"] = utils.StringValue(resp.Metadata.ProjectReference.UUID)
		}

		if err := d.Set("project_reference", pr); err != nil {
			return err
		}
	}

	if err := d.Set("api_version", utils.StringValue(resp.APIVersion)); err != nil {
		return err
	}

	if err := d.Set("name", utils.StringValue(resp.Spec.Name)); err != nil {
		return err
	}

	if err := d.Set("description", utils.StringValue(resp.Spec.Description)); err != nil {
		return err
	}

	if resp.Spec.Resources.QuarantineRule != nil {
		if err := d.Set("quarantine_rule_action", utils.StringValue(resp.Spec.Resources.QuarantineRule.Action)); err != nil {
			return err
		}

		if resp.Spec.Resources.QuarantineRule.OutboundAllowList != nil {
			oal := resp.Spec.Resources.QuarantineRule.OutboundAllowList
			qroaList := make([]map[string]interface{}, len(oal))
			for k, v := range oal {
				qroaItem := make(map[string]interface{})
				qroaItem["protocol"] = utils.StringValue(v.Protocol)

				if v.IPSubnet != nil {
					qroaItem["ip_subnet"] = utils.StringValue(v.IPSubnet.IP)
					qroaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(v.IPSubnet.PrefixLength), 10)
				}

				if v.TCPPortRangeList != nil {
					tcpprl := v.TCPPortRangeList
					tcpprList := make([]map[string]interface{}, len(tcpprl))
					for i, tcp := range tcpprl {
						tcpItem := make(map[string]interface{})
						tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
						tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
						tcpprList[i] = tcpItem
					}
					qroaItem["tcp_port_range_list"] = tcpprList
				}

				if v.UDPPortRangeList != nil {
					udpprl := v.UDPPortRangeList
					udpprList := make([]map[string]interface{}, len(udpprl))
					for i, udp := range udpprl {
						udpItem := make(map[string]interface{})
						udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
						udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
						udpprList[i] = udpItem
					}
					qroaItem["udp_port_range_list"] = udpprList
				}

				if v.Filter != nil {
					if v.Filter.KindList != nil {
						fkl := v.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						qroaItem["filter_kind_list"] = fkList
					}

					qroaItem["filter_type"] = utils.StringValue(v.Filter.Type)

					if v.Filter.Params != nil {
						fp := v.Filter.Params
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

				qroaItem["peer_specification_type"] = utils.StringValue(v.PeerSpecificationType)
				qroaItem["expiration_time"] = utils.StringValue(v.ExpirationTime)

				// set network_function_chain_reference
				if v.NetworkFunctionChainReference != nil {
					nfcr := make(map[string]interface{})
					nfcr["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
					nfcr["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
					nfcr["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
					qroaItem["network_function_chain_reference"] = nfcr
				}

				if v.IcmpTypeCodeList != nil {
					icmptcl := v.IcmpTypeCodeList
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
			if err := d.Set("quarantine_rule_outbound_allow_list", qroaList); err != nil {
				return err
			}
		}

		if resp.Spec.Resources.QuarantineRule.TargetGroup != nil {
			if err := d.Set("quarantine_rule_target_group_default_internal_policy",
				utils.StringValue(resp.Spec.Resources.QuarantineRule.TargetGroup.DefaultInternalPolicy)); err != nil {
				return err
			}
			if err := d.Set("quarantine_rule_target_group_peer_specification_type",
				utils.StringValue(resp.Spec.Resources.QuarantineRule.TargetGroup.PeerSpecificationType)); err != nil {
				return err
			}

			if resp.Spec.Resources.QuarantineRule.TargetGroup.Filter != nil {
				v := resp.Spec.Resources.QuarantineRule.TargetGroup
				if v.Filter != nil {
					if v.Filter.KindList != nil {
						fkl := v.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						if err := d.Set("quarantine_rule_target_group_filter_kind_list", fkList); err != nil {
							return err
						}
					}

					if err := d.Set("quarantine_rule_target_group_filter_type", utils.StringValue(v.Filter.Type)); err != nil {
						return err
					}

					if v.Filter.Params != nil {
						fp := v.Filter.Params
						var fpList []map[string]interface{}

						for name, values := range fp {
							fpItem := make(map[string]interface{})
							fpItem["name"] = name
							fpItem["values"] = values
							fpList = append(fpList, fpItem)
						}

						if err := d.Set("quarantine_rule_target_group_filter_params", fpList); err != nil {
							return err
						}
					}

				}
			}

		}

		if resp.Spec.Resources.QuarantineRule.InboundAllowList != nil {
			ial := resp.Spec.Resources.QuarantineRule.InboundAllowList
			qriaList := make([]map[string]interface{}, len(ial))
			for k, v := range ial {
				qriaItem := make(map[string]interface{})
				qriaItem["protocol"] = utils.StringValue(v.Protocol)

				if v.IPSubnet != nil {
					qriaItem["ip_subnet"] = utils.StringValue(v.IPSubnet.IP)
					qriaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(v.IPSubnet.PrefixLength), 10)
				}

				if v.TCPPortRangeList != nil {
					tcpprl := v.TCPPortRangeList
					tcpprList := make([]map[string]interface{}, len(tcpprl))
					for i, tcp := range tcpprl {
						tcpItem := make(map[string]interface{})
						tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
						tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
						tcpprList[i] = tcpItem
					}
					qriaItem["tcp_port_range_list"] = tcpprList
				}

				if v.UDPPortRangeList != nil {
					udpprl := v.UDPPortRangeList
					udpprList := make([]map[string]interface{}, len(udpprl))
					for i, udp := range udpprl {
						udpItem := make(map[string]interface{})
						udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
						udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
						udpprList[i] = udpItem
					}
					qriaItem["udp_port_range_list"] = udpprList
				}

				if v.Filter != nil {
					if v.Filter.KindList != nil {
						fkl := v.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						qriaItem["filter_kind_list"] = fkList
					}

					qriaItem["filter_type"] = utils.StringValue(v.Filter.Type)

					if v.Filter.Params != nil {
						fp := v.Filter.Params
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

				qriaItem["peer_specification_type"] = utils.StringValue(v.PeerSpecificationType)
				qriaItem["expiration_time"] = utils.StringValue(v.ExpirationTime)

				// set network_function_chain_reference
				if v.NetworkFunctionChainReference != nil {
					nfcr := make(map[string]interface{})
					nfcr["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
					nfcr["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
					nfcr["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
					qriaItem["network_function_chain_reference"] = nfcr
				}

				if v.IcmpTypeCodeList != nil {
					icmptcl := v.IcmpTypeCodeList
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
			if err := d.Set("quarantine_rule_inbound_allow_list", qriaList); err != nil {
				return err
			}
		}

	} else {
		if err := d.Set("quarantine_rule_inbound_allow_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("quarantine_rule_outbound_allow_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("quarantine_rule_target_group_filter_kind_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("quarantine_rule_target_group_filter_params", make([]string, 0)); err != nil {
			return err
		}
	}

	if resp.Spec.Resources.AppRule != nil {
		if err := d.Set("app_rule_action", utils.StringValue(resp.Spec.Resources.AppRule.Action)); err != nil {
			return err
		}

		if resp.Spec.Resources.AppRule.OutboundAllowList != nil {
			oal := resp.Spec.Resources.AppRule.OutboundAllowList
			aroaList := make([]map[string]interface{}, len(oal))
			for k, v := range oal {
				aroaItem := make(map[string]interface{})
				aroaItem["protocol"] = utils.StringValue(v.Protocol)

				if v.IPSubnet != nil {
					aroaItem["ip_subnet"] = utils.StringValue(v.IPSubnet.IP)
					aroaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(v.IPSubnet.PrefixLength), 10)
				}

				if v.TCPPortRangeList != nil {
					tcpprl := v.TCPPortRangeList
					tcpprList := make([]map[string]interface{}, len(tcpprl))
					for i, tcp := range tcpprl {
						tcpItem := make(map[string]interface{})
						tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
						tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
						tcpprList[i] = tcpItem
					}
					aroaItem["tcp_port_range_list"] = tcpprList
				}

				if v.UDPPortRangeList != nil {
					udpprl := v.UDPPortRangeList
					udpprList := make([]map[string]interface{}, len(udpprl))
					for i, udp := range udpprl {
						udpItem := make(map[string]interface{})
						udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
						udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
						udpprList[i] = udpItem
					}
					aroaItem["udp_port_range_list"] = udpprList
				}

				if v.Filter != nil {
					if v.Filter.KindList != nil {
						fkl := v.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						aroaItem["filter_kind_list"] = fkList
					}

					aroaItem["filter_type"] = utils.StringValue(v.Filter.Type)

					if v.Filter.Params != nil {
						fp := v.Filter.Params
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

				aroaItem["peer_specification_type"] = utils.StringValue(v.PeerSpecificationType)
				aroaItem["expiration_time"] = utils.StringValue(v.ExpirationTime)

				// set network_function_chain_reference
				if v.NetworkFunctionChainReference != nil {
					nfcr := make(map[string]interface{})
					nfcr["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
					nfcr["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
					nfcr["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
					aroaItem["network_function_chain_reference"] = nfcr
				}

				if v.IcmpTypeCodeList != nil {
					icmptcl := v.IcmpTypeCodeList
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
			if err := d.Set("app_rule_outbound_allow_list", aroaList); err != nil {
				return err
			}
		}

		if resp.Spec.Resources.AppRule.TargetGroup != nil {
			if err := d.Set("app_rule_target_group_default_internal_policy",
				utils.StringValue(resp.Spec.Resources.AppRule.TargetGroup.DefaultInternalPolicy)); err != nil {
				return err
			}
			if err := d.Set("app_rule_target_group_peer_specification_type",
				utils.StringValue(resp.Spec.Resources.AppRule.TargetGroup.PeerSpecificationType)); err != nil {
				return err
			}

			if resp.Spec.Resources.AppRule.TargetGroup.Filter != nil {
				v := resp.Spec.Resources.AppRule.TargetGroup
				if v.Filter != nil {
					if v.Filter.KindList != nil {
						fkl := v.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						if err := d.Set("app_rule_target_group_filter_kind_list", fkList); err != nil {
							return err
						}
					}

					if err := d.Set("app_rule_target_group_filter_type", utils.StringValue(v.Filter.Type)); err != nil {
						return err
					}

					if v.Filter.Params != nil {
						fp := v.Filter.Params
						var fpList []map[string]interface{}

						for name, values := range fp {
							fpItem := make(map[string]interface{})
							fpItem["name"] = name
							fpItem["values"] = values
							fpList = append(fpList, fpItem)
						}

						if err := d.Set("app_rule_target_group_filter_params", fpList); err != nil {
							return err
						}
					}

				}
			}

		}

		if resp.Spec.Resources.AppRule.InboundAllowList != nil {
			ial := resp.Spec.Resources.AppRule.InboundAllowList
			ariaList := make([]map[string]interface{}, len(ial))
			for k, v := range ial {
				ariaItem := make(map[string]interface{})
				ariaItem["protocol"] = utils.StringValue(v.Protocol)

				if v.IPSubnet != nil {
					ariaItem["ip_subnet"] = utils.StringValue(v.IPSubnet.IP)
					ariaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(v.IPSubnet.PrefixLength), 10)
				}

				if v.TCPPortRangeList != nil {
					tcpprl := v.TCPPortRangeList
					tcpprList := make([]map[string]interface{}, len(tcpprl))
					for i, tcp := range tcpprl {
						tcpItem := make(map[string]interface{})
						tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
						tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
						tcpprList[i] = tcpItem
					}
					ariaItem["tcp_port_range_list"] = tcpprList
				}

				if v.UDPPortRangeList != nil {
					udpprl := v.UDPPortRangeList
					udpprList := make([]map[string]interface{}, len(udpprl))
					for i, udp := range udpprl {
						udpItem := make(map[string]interface{})
						udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
						udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
						udpprList[i] = udpItem
					}
					ariaItem["udp_port_range_list"] = udpprList
				}

				if v.Filter != nil {
					if v.Filter.KindList != nil {
						fkl := v.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						ariaItem["filter_kind_list"] = fkList
					}

					ariaItem["filter_type"] = utils.StringValue(v.Filter.Type)

					if v.Filter.Params != nil {
						fp := v.Filter.Params
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

				ariaItem["peer_specification_type"] = utils.StringValue(v.PeerSpecificationType)
				ariaItem["expiration_time"] = utils.StringValue(v.ExpirationTime)

				// set network_function_chain_reference
				if v.NetworkFunctionChainReference != nil {
					nfcr := make(map[string]interface{})
					nfcr["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
					nfcr["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
					nfcr["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
					ariaItem["network_function_chain_reference"] = nfcr
				}

				if v.IcmpTypeCodeList != nil {
					icmptcl := v.IcmpTypeCodeList
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
			if err := d.Set("app_rule_inbound_allow_list", ariaList); err != nil {
				return err
			}
		}

	} else {
		if err := d.Set("app_rule_action", ""); err != nil {
			return err
		}
	}

	if resp.Spec.Resources.IsolationRule != nil {
		if err := d.Set("isolation_rule_action", utils.StringValue(resp.Spec.Resources.IsolationRule.Action)); err != nil {
			return err
		}

		if resp.Spec.Resources.IsolationRule.FirstEntityFilter != nil {
			firstFilter := resp.Spec.Resources.IsolationRule.FirstEntityFilter
			if firstFilter.KindList != nil {
				fkl := firstFilter.KindList
				fkList := make([]string, len(fkl))
				for i, f := range fkl {
					fkList[i] = utils.StringValue(f)
				}
				if err := d.Set("isolation_rule_first_entity_filter_kind_list", fkList); err != nil {
					return err
				}
			} else {
				if err := d.Set("isolation_rule_first_entity_filter_kind_list", make([]string, 0)); err != nil {
					return err
				}
			}

			if err := d.Set("isolation_rule_first_entity_filter_type", utils.StringValue(firstFilter.Type)); err != nil {
				return err
			}

			if firstFilter.Params != nil {
				fp := firstFilter.Params
				var fpList []map[string]interface{}

				for name, values := range fp {
					fpItem := make(map[string]interface{})
					fpItem["name"] = name
					fpItem["values"] = values
					fpList = append(fpList, fpItem)
				}

				if err := d.Set("isolation_rule_first_entity_filter_params", fpList); err != nil {
					return err
				}
			}

		}

		if resp.Spec.Resources.IsolationRule.SecondEntityFilter != nil {
			secondFilter := resp.Spec.Resources.IsolationRule.SecondEntityFilter
			if secondFilter.KindList != nil {
				fkl := secondFilter.KindList
				fkList := make([]string, len(fkl))
				for i, f := range fkl {
					fkList[i] = utils.StringValue(f)
				}
				if err := d.Set("isolation_rule_second_entity_filter_kind_list", fkList); err != nil {
					return err
				}
			}

			if err := d.Set("isolation_rule_second_entity_filter_type", utils.StringValue(secondFilter.Type)); err != nil {
				return err
			}

			if secondFilter.Params != nil {
				fp := secondFilter.Params
				var fpList []map[string]interface{}

				for name, values := range fp {
					fpItem := make(map[string]interface{})
					fpItem["name"] = name
					fpItem["values"] = values
					fpList = append(fpList, fpItem)
				}

				if err := d.Set("isolation_rule_second_entity_filter_params", fpList); err != nil {
					return err
				}
			}

		}

	} else {
		if err := d.Set("isolation_rule_first_entity_filter_kind_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("isolation_rule_first_entity_filter_params", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("isolation_rule_second_entity_filter_kind_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("isolation_rule_second_entity_filter_params", make([]string, 0)); err != nil {
			return err
		}
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func getDataSourceNetworkSecurityRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}
