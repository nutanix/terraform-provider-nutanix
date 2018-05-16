package nutanix

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func resourceNutanixNetworkSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixNetworkSecurityRuleCreate,
		Read:   resourceNutanixNetworkSecurityRuleRead,
		Update: resourceNutanixNetworkSecurityRuleUpdate,
		Delete: resourceNutanixNetworkSecurityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getNetworkSecurityRuleSchema(),
	}
}

func resourceNutanixNetworkSecurityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection
	conn := meta.(*NutanixClient).API

	// Prepare request
	request := &v3.NetworkSecurityRuleIntentInput{}
	spec := &v3.NetworkSecurityRule{}
	metadata := &v3.Metadata{}
	networkSecurityRule := &v3.NetworkSecurityRuleResources{}

	//Read arguments and set request values
	name, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")

	if !nok {
		return fmt.Errorf("Please provide the required attribute name")
	}

	// Read arguments and set request values
	if v, ok := d.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}

	//only set kind
	if errMetad := getMetadataAttributes(d, metadata, "network_security_rule"); errMetad != nil {
		return errMetad
	}

	if descok {
		spec.Description = utils.String(desc.(string))
	}

	//get resources
	if err := getNetworkSecurityRuleResources(d, networkSecurityRule); err != nil {
		return err
	}

	if descok {
		spec.Description = utils.String(desc.(string))
	}

	networkSecurityRueUUID, err := resourceNutanixNetworkSecurityRuleExists(conn, d.Get("name").(string))

	if err != nil {
		return err
	}

	if networkSecurityRueUUID != nil {
		return fmt.Errorf("Network Security Rule already with name %s exists in the given cluster, UUID %s", d.Get("name").(string), *networkSecurityRueUUID)
	}

	//set request

	spec.Resources = networkSecurityRule

	spec.Name = utils.String(name.(string))

	// set request attrs
	request.Metadata = metadata
	request.Spec = spec

	// Make request to API
	resp, err := conn.V3.CreateNetworkSecurityRule(request)

	if err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    networkSecurityRuleStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for network_security_rule (%s) to create: %s", d.Id(), err)
	}

	return resourceNutanixNetworkSecurityRuleRead(d, meta)
}

func resourceNutanixNetworkSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Network Security Rule: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*NutanixClient).API

	// Make request to the API
	resp, err := conn.V3.GetNetworkSecurityRule(d.Id())
	if err != nil {
		return err
	}

	// set metadata values
	if resp.Metadata != nil {
		metadata := make(map[string]interface{})

		metadata["last_update_time"] = resp.Metadata.LastUpdateTime.String()
		//metadata["kind"] = utils.StringValue(resp.Metadata.Kind)
		metadata["uuid"] = utils.StringValue(resp.Metadata.UUID)
		metadata["creation_time"] = resp.Metadata.CreationTime.String()
		metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(resp.Metadata.SpecVersion)))
		metadata["spec_hash"] = utils.StringValue(resp.Metadata.SpecHash)
		metadata["name"] = utils.StringValue(resp.Metadata.Name)

		if err := d.Set("metadata", metadata); err != nil {
			return err
		}

		if err := d.Set("categories", resp.Metadata.Categories); err != nil {
			return err
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

	if resp.Status.QuarantineRule != nil {
		if err := d.Set("quarantine_rule_action", utils.StringValue(resp.Status.QuarantineRule.Action)); err != nil {
			return err
		}

		if resp.Status.QuarantineRule.OutboundAllowList != nil {
			oal := resp.Status.QuarantineRule.OutboundAllowList
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

		if resp.Status.QuarantineRule.TargetGroup != nil {
			if err := d.Set("quarantine_rule_target_group_default_internal_policy",
				utils.StringValue(resp.Status.QuarantineRule.TargetGroup.DefaultInternalPolicy)); err != nil {
				return err
			}
			if err := d.Set("quarantine_rule_target_group_peer_specification_type",
				utils.StringValue(resp.Status.QuarantineRule.TargetGroup.PeerSpecificationType)); err != nil {
				return err
			}

			if resp.Status.QuarantineRule.TargetGroup.Filter != nil {
				v := resp.Status.QuarantineRule.TargetGroup
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

		if resp.Status.QuarantineRule.InboundAllowList != nil {
			ial := resp.Status.QuarantineRule.InboundAllowList
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

	if resp.Status.AppRule != nil {
		if err := d.Set("app_rule_action", utils.StringValue(resp.Status.AppRule.Action)); err != nil {
			return err
		}

		if resp.Status.AppRule.OutboundAllowList != nil {
			oal := resp.Status.AppRule.OutboundAllowList
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

		if resp.Status.AppRule.TargetGroup != nil {
			if err := d.Set("app_rule_target_group_default_internal_policy",
				utils.StringValue(resp.Status.AppRule.TargetGroup.DefaultInternalPolicy)); err != nil {
				return err
			}
			if err := d.Set("app_rule_target_group_peer_specification_type",
				utils.StringValue(resp.Status.AppRule.TargetGroup.PeerSpecificationType)); err != nil {
				return err
			}

			if resp.Status.AppRule.TargetGroup.Filter != nil {
				v := resp.Status.AppRule.TargetGroup
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

		if resp.Status.AppRule.InboundAllowList != nil {
			ial := resp.Status.AppRule.InboundAllowList
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

	}

	if resp.Status.IsolationRule != nil {
		if err := d.Set("isolation_rule_action", utils.StringValue(resp.Status.IsolationRule.Action)); err != nil {
			return err
		}

		if resp.Status.IsolationRule.FirstEntityFilter != nil {
			firstFilter := resp.Status.IsolationRule.FirstEntityFilter
			if firstFilter.KindList != nil {
				log.Println("KIND LIST##")
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

		if resp.Status.IsolationRule.SecondEntityFilter != nil {
			secondFilter := resp.Status.IsolationRule.SecondEntityFilter
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

	return nil
}

func resourceNutanixNetworkSecurityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection
	conn := meta.(*NutanixClient).API

	// Prepare request
	request := &v3.NetworkSecurityRuleIntentInput{}
	spec := &v3.NetworkSecurityRule{}
	metadata := &v3.Metadata{}
	networkSecurityRule := &v3.NetworkSecurityRuleResources{}

	response, err := conn.V3.GetNetworkSecurityRule(d.Id())

	if err != nil {
		return err
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			networkSecurityRule = response.Spec.Resources
		}
	}

	if d.HasChange("categories") {
		p := d.Get("categories").(map[string]interface{})
		labels := map[string]string{}
		for k, v := range p {
			labels[k] = v.(string)
		}
		metadata.Categories = labels
	}

	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(or["kind"].(string)),
			UUID: utils.String(or["uuid"].(string)),
			Name: utils.String(or["name"].(string)),
		}
		metadata.OwnerReference = r
	}

	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pr["kind"].(string)),
			UUID: utils.String(pr["uuid"].(string)),
			Name: utils.String(pr["name"].(string)),
		}
		metadata.ProjectReference = r
	}

	if d.HasChange("name") {
		spec.Name = utils.String(d.Get("name").(string))
	}
	if d.HasChange("description") {
		spec.Description = utils.String(d.Get("description").(string))
	}

	//TODO: Change
	if d.HasChange("quarantine_rule_action") ||
		d.HasChange("quarantine_rule_outbound_allow_list") ||
		d.HasChange("quarantine_rule_target_group_default_internal_policy") ||
		d.HasChange("quarantine_rule_target_group_peer_specification_type") ||
		d.HasChange("quarantine_rule_target_group_filter_kind_list") ||
		d.HasChange("quarantine_rule_target_group_filter_type") ||
		d.HasChange("quarantine_rule_target_group_filter_params") ||
		d.HasChange("quarantine_rule_inbound_allow_list") ||
		d.HasChange("app_rule_action") ||
		d.HasChange("app_rule_outbound_allow_list") ||
		d.HasChange("app_rule_target_group_default_internal_policy") ||
		d.HasChange("app_rule_target_group_peer_specification_type") ||
		d.HasChange("app_rule_target_group_filter_kind_list") ||
		d.HasChange("app_rule_target_group_filter_type") ||
		d.HasChange("app_rule_target_group_filter_params") ||
		d.HasChange("app_rule_inbound_allow_list") ||
		d.HasChange("isolation_rule_action") ||
		d.HasChange("isolation_rule_first_entity_filter_kind_list") ||
		d.HasChange("isolation_rule_first_entity_filter_type") ||
		d.HasChange("isolation_rule_first_entity_filter_params") ||
		d.HasChange("isolation_rule_second_entity_filter_kind_list") ||
		d.HasChange("isolation_rule_second_entity_filter_type") ||
		d.HasChange("isolation_rule_second_entity_filter_params") {

		if err := getNetworkSecurityRuleResources(d, networkSecurityRule); err != nil {
			return err
		}
		spec.Resources = networkSecurityRule
	}

	request.Spec = spec
	request.Metadata = metadata

	utils.PrintToJSON(request, "UPDATE REQUEST ###")

	_, errUpdate := conn.V3.UpdateNetworkSecurityRule(d.Id(), request)

	if errUpdate != nil {
		return errUpdate
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    networkSecurityRuleStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for network_security_rule (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixNetworkSecurityRuleRead(d, meta)

}

func resourceNutanixNetworkSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting Network Security Rule: %s", d.Get("name").(string))

	conn := meta.(*NutanixClient).API
	UUID := d.Id()

	if err := conn.V3.DeleteNetworkSecurityRule(UUID); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "DELETE_IN_PROGRESS", "COMPLETE"},
		Target:     []string{"DELETED"},
		Refresh:    networkSecurityRuleStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for network_security_rule (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixNetworkSecurityRuleExists(conn *v3.Client, name string) (*string, error) {
	log.Printf("[DEBUG] Get Network Security Rule Existance : %s", name)

	subnetEntities := &v3.ListMetadata{}
	var nsrUUID *string

	networkSecurityRuleList, err := conn.V3.ListNetworkSecurityRule(subnetEntities)

	if err != nil {
		return nil, err
	}

	for _, nsr := range networkSecurityRuleList.Entities {
		if nsr.Metadata.Name == utils.String(name) {
			nsrUUID = nsr.Metadata.UUID
		}
	}
	return nsrUUID, nil
}

func getNetworkSecurityRuleResources(d *schema.ResourceData, networkSecurityRule *v3.NetworkSecurityRuleResources) error {
	isolationRule := &v3.NetworkSecurityRuleIsolationRule{}
	quarantineRule := &v3.NetworkSecurityRuleResourcesRule{}
	qRuleTargetGroup := &v3.TargetGroup{}
	qRuleTargetGroupFilter := &v3.CategoryFilter{}
	appRule := &v3.NetworkSecurityRuleResourcesRule{}
	aRuleTargetGroup := &v3.TargetGroup{}
	aRuleTargetGroupFilter := &v3.CategoryFilter{}
	iRuleFirstEntityFilter := &v3.CategoryFilter{}
	iRuleSecondEntityFilter := &v3.CategoryFilter{}

	if qra, ok := d.GetOk("quarantine_rule_action"); ok && qra.(string) != "" {
		quarantineRule.Action = utils.String(qra.(string))
	}

	if qroal, ok := d.GetOk("quarantine_rule_outbound_allow_list"); ok && qroal != nil {
		oal := qroal.([]interface{})
		outbound := make([]*v3.NetworkRule, len(oal))

		for k, v := range oal {
			nr := v.(map[string]interface{})
			nrItem := &v3.NetworkRule{}
			iPSubnet := &v3.IPSubnet{}
			filter := &v3.CategoryFilter{}

			if proto, pok := nr["protocol"]; pok && proto.(string) != "" {
				nrItem.Protocol = utils.String(proto.(string))
			}

			if ip, ipok := nr["ip_subnet"]; ipok && ip.(string) != "" {
				iPSubnet.IP = utils.String(ip.(string))
			}

			if ippl, ipok := nr["ip_subnet_prefix_length"]; ipok && ippl.(string) != "" {
				if i, err := strconv.Atoi(ippl.(string)); err != nil {
					iPSubnet.PrefixLength = utils.Int64(int64(i))
				}
			}

			if t, tcpok := nr["tcp_port_range_list"]; tcpok {
				tcplist := t.([]interface{})
				tcpPorts := make([]*v3.PortRange, len(tcplist))

				for i, p := range tcplist {
					tcpp := p.(map[string]interface{})
					portRange := &v3.PortRange{}

					if endp, epok := tcpp["end_port"]; epok {
						if i, err := strconv.Atoi(endp.(string)); err != nil {
							portRange.EndPort = utils.Int64(int64(i))
						}
					}

					if stp, stpok := tcpp["start_port"]; stpok {
						if i, err := strconv.Atoi(stp.(string)); err != nil {
							portRange.StartPort = utils.Int64(int64(i))
						}
					}
					tcpPorts[i] = portRange
				}
				nrItem.TCPPortRangeList = tcpPorts
			}

			if u, udpok := nr["udp_port_range_list"]; udpok {
				udplist := u.([]interface{})
				udpPorts := make([]*v3.PortRange, len(udplist))

				for i, p := range udplist {
					udpp := p.(map[string]interface{})
					portRange := &v3.PortRange{}

					if endp, epok := udpp["end_port"]; epok {
						if i, err := strconv.Atoi(endp.(string)); err != nil {
							portRange.EndPort = utils.Int64(int64(i))
						}
					}

					if stp, stpok := udpp["start_port"]; stpok {
						if i, err := strconv.Atoi(stp.(string)); err != nil {
							portRange.StartPort = utils.Int64(int64(i))
						}
					}
					udpPorts[i] = portRange
				}
				nrItem.UDPPortRangeList = udpPorts
			}

			if f, fok := nr["filter_kind_list"]; fok {
				fkl := f.([]interface{})
				fkList := make([]*string, len(fkl))
				for k, v := range fkl {
					fkList[k] = utils.String(v.(string))
				}
				filter.KindList = fkList
			}

			if ft, ftok := nr["filter_type"]; ftok {
				filter.Type = utils.String(ft.(string))
			}

			if fp, fpok := nr["filter_params"]; fpok {
				fpl := fp.([]interface{})

				if len(fpl) > 0 {
					fl := make(map[string][]string)
					for _, v := range fpl {
						item := v.(map[string]interface{})

						if i, ok := item["name"]; ok && i.(string) != "" {
							if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
								var values []string
								for _, item := range k.([]interface{}) {
									values = append(values, item.(string))
								}
								fl[i.(string)] = values
							}

						}
					}
					filter.Params = fl
				} else {
					filter.Params = nil
				}

			}

			if pet, petok := nr["peer_specification_type"]; petok && pet.(string) != "" {
				nrItem.PeerSpecificationType = utils.String(pet.(string))
			}

			if et, etok := nr["expiration_time"]; etok && et.(string) != "" {
				nrItem.ExpirationTime = utils.String(et.(string))
			}

			if nfcr, nfcrok := nr["network_function_chain_reference"]; nfcrok && len(nfcr.(map[string]interface{})) > 0 {
				a := nfcr.(map[string]interface{})
				r := &v3.Reference{
					Kind: utils.String(a["kind"].(string)),
					UUID: utils.String(a["uuid"].(string)),
				}
				if v, ok := a["name"]; ok {
					r.Name = utils.String(v.(string))
				}
				nrItem.NetworkFunctionChainReference = r
			}

			if icmp, icmpok := nr["icmp_type_code_list"]; icmpok {
				ic := icmp.([]interface{})

				if len(ic) > 0 {
					icmpList := make([]*v3.NetworkRuleIcmpTypeCodeList, len(ic))

					for k, v := range ic {
						icmpm := v.(map[string]interface{})
						icmpItem := &v3.NetworkRuleIcmpTypeCodeList{}

						if c, cok := icmpm["code"]; cok && c.(string) != "" {

							if i, err := strconv.Atoi(c.(string)); err != nil {
								icmpItem.Code = utils.Int64(int64(i))
							}
						}

						if t, tok := icmpm["type"]; tok && t.(string) != "" {
							if i, err := strconv.Atoi(t.(string)); err != nil {
								icmpItem.Type = utils.Int64(int64(i))
							}
						}
						icmpList[k] = icmpItem
					}
					nrItem.IcmpTypeCodeList = icmpList
				} else {
					nrItem.IcmpTypeCodeList = nil
				}

			}

			nrItem.IPSubnet = iPSubnet
			nrItem.Filter = filter
			outbound[k] = nrItem
		}
		quarantineRule.OutboundAllowList = outbound
	}

	if qroal, ok := d.GetOk("quarantine_rule_target_group_default_internal_policy"); ok && qroal.(string) != "" {
		qRuleTargetGroup.DefaultInternalPolicy = utils.String(qroal.(string))
	}

	if qroal, ok := d.GetOk("quarantine_rule_target_group_peer_specification_type"); ok && qroal.(string) != "" {
		qRuleTargetGroup.PeerSpecificationType = utils.String(qroal.(string))
	}

	if f, fok := d.GetOk("quarantine_rule_target_group_filter_kind_list"); fok && f != nil {
		fkl := f.([]interface{})
		fkList := make([]*string, len(fkl))
		for k, v := range fkl {
			fkList[k] = utils.String(v.(string))
		}
		qRuleTargetGroupFilter.KindList = fkList
	}

	if ft, ftok := d.GetOk("quarantine_rule_target_group_filter_type"); ftok && ft.(string) != "" {
		qRuleTargetGroupFilter.Type = utils.String(ft.(string))
	}

	if fp, fpok := d.GetOk("quarantine_rule_target_group_filter_params"); fpok {
		fpl := fp.([]interface{})

		if len(fpl) > 0 {
			fl := make(map[string][]string)
			for _, v := range fpl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
						var values []string
						for _, item := range k.([]interface{}) {
							values = append(values, item.(string))
						}
						fl[i.(string)] = values
					}

				}
			}
			qRuleTargetGroupFilter.Params = fl
		} else {
			qRuleTargetGroupFilter.Params = nil
		}

	}

	if qrial, ok := d.GetOk("quarantine_rule_inbound_allow_list"); ok {
		oal := qrial.([]interface{})
		inbound := make([]*v3.NetworkRule, len(oal))

		for k, v := range oal {
			nr := v.(map[string]interface{})
			nrItem := &v3.NetworkRule{}
			iPSubnet := &v3.IPSubnet{}
			filter := &v3.CategoryFilter{}

			if proto, pok := nr["protocol"]; pok && proto.(string) != "" {
				nrItem.Protocol = utils.String(proto.(string))
			}

			if ip, ipok := nr["ip_subnet"]; ipok && ip.(string) != "" {
				iPSubnet.IP = utils.String(ip.(string))
			}

			if ippl, ipok := nr["ip_subnet_prefix_length"]; ipok && ippl.(string) != "" {
				if i, err := strconv.Atoi(ippl.(string)); err != nil {
					iPSubnet.PrefixLength = utils.Int64(int64(i))
				}
			}

			if t, tcpok := nr["tcp_port_range_list"]; tcpok {
				tcplist := t.([]interface{})
				tcpPorts := make([]*v3.PortRange, len(tcplist))

				for i, p := range tcplist {
					tcpp := p.(map[string]interface{})
					portRange := &v3.PortRange{}

					if endp, epok := tcpp["end_port"]; epok {
						if i, err := strconv.Atoi(endp.(string)); err != nil {
							portRange.EndPort = utils.Int64(int64(i))
						}
					}

					if stp, stpok := tcpp["start_port"]; stpok {
						if i, err := strconv.Atoi(stp.(string)); err != nil {
							portRange.StartPort = utils.Int64(int64(i))
						}
					}
					tcpPorts[i] = portRange
				}
				nrItem.TCPPortRangeList = tcpPorts
			}

			if u, udpok := nr["udp_port_range_list"]; udpok {
				udplist := u.([]interface{})
				udpPorts := make([]*v3.PortRange, len(udplist))

				for i, p := range udplist {
					udpp := p.(map[string]interface{})
					portRange := &v3.PortRange{}

					if endp, epok := udpp["end_port"]; epok {
						if i, err := strconv.Atoi(endp.(string)); err != nil {
							portRange.EndPort = utils.Int64(int64(i))
						}
					}

					if stp, stpok := udpp["start_port"]; stpok {
						if i, err := strconv.Atoi(stp.(string)); err != nil {
							portRange.StartPort = utils.Int64(int64(i))
						}
					}
					udpPorts[i] = portRange
				}
				nrItem.UDPPortRangeList = udpPorts
			}

			if f, fok := nr["filter_kind_list"]; fok {
				fkl := f.([]interface{})
				fkList := make([]*string, len(fkl))
				for k, v := range fkl {
					fkList[k] = utils.String(v.(string))
				}
				filter.KindList = fkList
			}

			if ft, ftok := nr["filter_type"]; ftok {
				filter.Type = utils.String(ft.(string))
			}

			if fp, fpok := nr["filter_params"]; fpok {
				fpl := fp.([]interface{})

				if len(fpl) > 0 {
					fl := make(map[string][]string)
					for _, v := range fpl {
						item := v.(map[string]interface{})

						if i, ok := item["name"]; ok && i.(string) != "" {
							if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
								var values []string
								for _, item := range k.([]interface{}) {
									values = append(values, item.(string))
								}
								fl[i.(string)] = values
							}

						}
					}
					filter.Params = fl
				} else {
					filter.Params = nil
				}

			}

			if pet, petok := nr["peer_specification_type"]; petok && pet.(string) != "" {
				nrItem.PeerSpecificationType = utils.String(pet.(string))
			}

			if et, etok := nr["expiration_time"]; etok && et.(string) != "" {
				nrItem.ExpirationTime = utils.String(et.(string))
			}

			if nfcr, nfcrok := nr["network_function_chain_reference"]; nfcrok && nfcr.(string) != "" {
				a := nfcr.(map[string]interface{})
				r := &v3.Reference{
					Kind: utils.String(a["kind"].(string)),
					UUID: utils.String(a["uuid"].(string)),
				}
				if v, ok := a["name"]; ok {
					r.Name = utils.String(v.(string))
				}
				nrItem.NetworkFunctionChainReference = r
			}

			if icmp, icmpok := nr["icmp_type_code_list"]; icmpok {
				ic := icmp.([]interface{})

				if len(ic) > 0 {
					icmpList := make([]*v3.NetworkRuleIcmpTypeCodeList, len(ic))

					for k, v := range ic {
						icmpm := v.(map[string]interface{})
						icmpItem := &v3.NetworkRuleIcmpTypeCodeList{}

						if c, cok := icmpm["code"]; cok && c.(string) != "" {

							if i, err := strconv.Atoi(c.(string)); err != nil {
								icmpItem.Code = utils.Int64(int64(i))
							}
						}

						if t, tok := icmpm["type"]; tok && t.(string) != "" {
							if i, err := strconv.Atoi(t.(string)); err != nil {
								icmpItem.Type = utils.Int64(int64(i))
							}
						}
						icmpList[k] = icmpItem
					}
					nrItem.IcmpTypeCodeList = icmpList
				} else {
					nrItem.IcmpTypeCodeList = nil
				}

			}

			nrItem.IPSubnet = iPSubnet
			nrItem.Filter = filter
			inbound[k] = nrItem
		}
		quarantineRule.InboundAllowList = inbound
	}

	if ara, ok := d.GetOk("app_rule_action"); ok && ara.(string) != "" {
		appRule.Action = utils.String(ara.(string))
	}

	if qroal, ok := d.GetOk("app_rule_outbound_allow_list"); ok {
		oal := qroal.([]interface{})
		outbound := make([]*v3.NetworkRule, len(oal))

		for k, v := range oal {
			nr := v.(map[string]interface{})
			nrItem := &v3.NetworkRule{}
			iPSubnet := &v3.IPSubnet{}
			filter := &v3.CategoryFilter{}

			if proto, pok := nr["protocol"]; pok && proto.(string) != "" {
				nrItem.Protocol = utils.String(proto.(string))
			}

			if ip, ipok := nr["ip_subnet"]; ipok && ip.(string) != "" {
				iPSubnet.IP = utils.String(ip.(string))
			}

			if ippl, ipok := nr["ip_subnet_prefix_length"]; ipok && ippl.(string) != "" {
				if i, err := strconv.Atoi(ippl.(string)); err != nil {
					iPSubnet.PrefixLength = utils.Int64(int64(i))
				}
			}

			if t, tcpok := nr["tcp_port_range_list"]; tcpok {
				tcplist := t.([]interface{})
				tcpPorts := make([]*v3.PortRange, len(tcplist))

				for i, p := range tcplist {
					tcpp := p.(map[string]interface{})
					portRange := &v3.PortRange{}

					if endp, epok := tcpp["end_port"]; epok {
						if i, err := strconv.Atoi(endp.(string)); err != nil {
							portRange.EndPort = utils.Int64(int64(i))
						}
					}

					if stp, stpok := tcpp["start_port"]; stpok {
						if i, err := strconv.Atoi(stp.(string)); err != nil {
							portRange.StartPort = utils.Int64(int64(i))
						}
					}
					tcpPorts[i] = portRange
				}
				nrItem.TCPPortRangeList = tcpPorts
			}

			if u, udpok := nr["udp_port_range_list"]; udpok {
				udplist := u.([]interface{})
				udpPorts := make([]*v3.PortRange, len(udplist))

				for i, p := range udplist {
					udpp := p.(map[string]interface{})
					portRange := &v3.PortRange{}

					if endp, epok := udpp["end_port"]; epok {
						if i, err := strconv.Atoi(endp.(string)); err != nil {
							portRange.EndPort = utils.Int64(int64(i))
						}
					}

					if stp, stpok := udpp["start_port"]; stpok {
						if i, err := strconv.Atoi(stp.(string)); err != nil {
							portRange.StartPort = utils.Int64(int64(i))
						}
					}
					udpPorts[i] = portRange
				}
				nrItem.UDPPortRangeList = udpPorts
			}

			if f, fok := nr["filter_kind_list"]; fok {
				fkl := f.([]interface{})
				fkList := make([]*string, len(fkl))
				for k, v := range fkl {
					fkList[k] = utils.String(v.(string))
				}
				filter.KindList = fkList
			}

			if ft, ftok := nr["filter_type"]; ftok {
				filter.Type = utils.String(ft.(string))
			}

			if fp, fpok := nr["filter_params"]; fpok {
				fpl := fp.([]interface{})

				if len(fpl) > 0 {
					fl := make(map[string][]string)
					for _, v := range fpl {
						item := v.(map[string]interface{})

						if i, ok := item["name"]; ok && i.(string) != "" {
							if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
								var values []string
								for _, item := range k.([]interface{}) {
									values = append(values, item.(string))
								}
								fl[i.(string)] = values
							}

						}
					}
					filter.Params = fl
				} else {
					filter.Params = nil
				}

			}

			if pet, petok := nr["peer_specification_type"]; petok && pet.(string) != "" {
				nrItem.PeerSpecificationType = utils.String(pet.(string))
			}

			if et, etok := nr["expiration_time"]; etok && et.(string) != "" {
				nrItem.ExpirationTime = utils.String(et.(string))
			}

			if nfcr, nfcrok := nr["network_function_chain_reference"]; nfcrok && len(nfcr.(map[string]interface{})) > 0 {
				a := nfcr.(map[string]interface{})
				r := &v3.Reference{
					Kind: utils.String(a["kind"].(string)),
					UUID: utils.String(a["uuid"].(string)),
				}
				if v, ok := a["name"]; ok {
					r.Name = utils.String(v.(string))
				}
				nrItem.NetworkFunctionChainReference = r
			}

			if icmp, icmpok := nr["icmp_type_code_list"]; icmpok {
				ic := icmp.([]interface{})

				if len(ic) > 0 {
					icmpList := make([]*v3.NetworkRuleIcmpTypeCodeList, len(ic))

					for k, v := range ic {
						icmpm := v.(map[string]interface{})
						icmpItem := &v3.NetworkRuleIcmpTypeCodeList{}

						if c, cok := icmpm["code"]; cok && c.(string) != "" {

							if i, err := strconv.Atoi(c.(string)); err != nil {
								icmpItem.Code = utils.Int64(int64(i))
							}
						}

						if t, tok := icmpm["type"]; tok && t.(string) != "" {
							if i, err := strconv.Atoi(t.(string)); err != nil {
								icmpItem.Type = utils.Int64(int64(i))
							}
						}
						icmpList[k] = icmpItem
					}
					nrItem.IcmpTypeCodeList = icmpList
				} else {
					nrItem.IcmpTypeCodeList = nil
				}

			}

			nrItem.IPSubnet = iPSubnet
			nrItem.Filter = filter
			outbound[k] = nrItem
		}
		appRule.OutboundAllowList = outbound
	}

	if qroal, ok := d.GetOk("app_rule_target_group_default_internal_policy"); ok && qroal != nil {
		aRuleTargetGroup.DefaultInternalPolicy = utils.String(qroal.(string))
	}

	if qroal, ok := d.GetOk("app_rule_target_group_peer_specification_type"); ok && qroal != nil {
		aRuleTargetGroup.PeerSpecificationType = utils.String(qroal.(string))
	}

	if f, fok := d.GetOk("app_rule_target_group_filter_kind_list"); fok && f != nil {
		fkl := f.([]interface{})
		fkList := make([]*string, len(fkl))
		for k, v := range fkl {
			fkList[k] = utils.String(v.(string))
		}
		aRuleTargetGroupFilter.KindList = fkList
	}

	if ft, ftok := d.GetOk("app_rule_target_group_filter_type"); ftok && ft.(string) != "" {
		aRuleTargetGroupFilter.Type = utils.String(ft.(string))
	}

	if fp, fpok := d.GetOk("app_rule_target_group_filter_params"); fpok {
		fpl := fp.([]interface{})

		if len(fpl) > 0 {
			fl := make(map[string][]string)
			for _, v := range fpl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
						var values []string
						for _, item := range k.([]interface{}) {
							values = append(values, item.(string))
						}
						fl[i.(string)] = values
					}

				}
			}
			aRuleTargetGroupFilter.Params = fl
		} else {
			aRuleTargetGroupFilter.Params = nil
		}

	}

	if qrial, ok := d.GetOk("app_rule_inbound_allow_list"); ok {
		oal := qrial.([]interface{})
		inbound := make([]*v3.NetworkRule, len(oal))

		for k, v := range oal {
			nr := v.(map[string]interface{})
			nrItem := &v3.NetworkRule{}
			iPSubnet := &v3.IPSubnet{}
			filter := &v3.CategoryFilter{}

			if proto, pok := nr["protocol"]; pok && proto.(string) != "" {
				nrItem.Protocol = utils.String(proto.(string))
			}

			if ip, ipok := nr["ip_subnet"]; ipok && ip.(string) != "" {
				iPSubnet.IP = utils.String(ip.(string))
			}

			if ippl, ipok := nr["ip_subnet_prefix_length"]; ipok && ippl.(string) != "" {
				if i, err := strconv.Atoi(ippl.(string)); err != nil {
					iPSubnet.PrefixLength = utils.Int64(int64(i))
				}
			}

			if t, tcpok := nr["tcp_port_range_list"]; tcpok {
				tcplist := t.([]interface{})
				tcpPorts := make([]*v3.PortRange, len(tcplist))

				for i, p := range tcplist {
					tcpp := p.(map[string]interface{})
					portRange := &v3.PortRange{}

					if endp, epok := tcpp["end_port"]; epok {
						if i, err := strconv.Atoi(endp.(string)); err != nil {
							portRange.EndPort = utils.Int64(int64(i))
						}
					}

					if stp, stpok := tcpp["start_port"]; stpok {
						if i, err := strconv.Atoi(stp.(string)); err != nil {
							portRange.StartPort = utils.Int64(int64(i))
						}
					}
					tcpPorts[i] = portRange
				}
				nrItem.TCPPortRangeList = tcpPorts
			}

			if u, udpok := nr["udp_port_range_list"]; udpok {
				udplist := u.([]interface{})
				udpPorts := make([]*v3.PortRange, len(udplist))

				for i, p := range udplist {
					udpp := p.(map[string]interface{})
					portRange := &v3.PortRange{}

					if endp, epok := udpp["end_port"]; epok {
						if i, err := strconv.Atoi(endp.(string)); err != nil {
							portRange.EndPort = utils.Int64(int64(i))
						}
					}

					if stp, stpok := udpp["start_port"]; stpok {
						if i, err := strconv.Atoi(stp.(string)); err != nil {
							portRange.StartPort = utils.Int64(int64(i))
						}
					}
					udpPorts[i] = portRange
				}
				nrItem.UDPPortRangeList = udpPorts
			}

			if f, fok := nr["filter_kind_list"]; fok {
				fkl := f.([]interface{})
				fkList := make([]*string, len(fkl))
				for k, v := range fkl {
					fkList[k] = utils.String(v.(string))
				}
				filter.KindList = fkList
			}

			if ft, ftok := nr["filter_type"]; ftok {
				filter.Type = utils.String(ft.(string))
			}

			if fp, fpok := nr["filter_params"]; fpok {
				fpl := fp.([]interface{})

				if len(fpl) > 0 {
					fl := make(map[string][]string)
					for _, v := range fpl {
						item := v.(map[string]interface{})

						if i, ok := item["name"]; ok && i.(string) != "" {
							if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
								var values []string
								for _, item := range k.([]interface{}) {
									values = append(values, item.(string))
								}
								fl[i.(string)] = values
							}

						}
					}
					filter.Params = fl
				} else {
					filter.Params = nil
				}

			}

			if pet, petok := nr["peer_specification_type"]; petok && pet.(string) != "" {
				nrItem.PeerSpecificationType = utils.String(pet.(string))
			}

			if et, etok := nr["expiration_time"]; etok && et.(string) != "" {
				nrItem.ExpirationTime = utils.String(et.(string))
			}

			if nfcr, nfcrok := nr["network_function_chain_reference"]; nfcrok && len(nfcr.(map[string]interface{})) > 0 {
				a := nfcr.(map[string]interface{})
				r := &v3.Reference{
					Kind: utils.String(a["kind"].(string)),
					UUID: utils.String(a["uuid"].(string)),
				}
				if v, ok := a["name"]; ok {
					r.Name = utils.String(v.(string))
				}
				nrItem.NetworkFunctionChainReference = r
			}

			if icmp, icmpok := nr["icmp_type_code_list"]; icmpok {
				ic := icmp.([]interface{})

				if len(ic) > 0 {
					icmpList := make([]*v3.NetworkRuleIcmpTypeCodeList, len(ic))

					for k, v := range ic {
						icmpm := v.(map[string]interface{})
						icmpItem := &v3.NetworkRuleIcmpTypeCodeList{}

						if c, cok := icmpm["code"]; cok && c.(string) != "" {

							if i, err := strconv.Atoi(c.(string)); err != nil {
								icmpItem.Code = utils.Int64(int64(i))
							}
						}

						if t, tok := icmpm["type"]; tok && t.(string) != "" {
							if i, err := strconv.Atoi(t.(string)); err != nil {
								icmpItem.Type = utils.Int64(int64(i))
							}
						}
						icmpList[k] = icmpItem
					}
					nrItem.IcmpTypeCodeList = icmpList
				} else {
					nrItem.IcmpTypeCodeList = nil
				}

			}

			nrItem.IPSubnet = iPSubnet
			nrItem.Filter = filter
			inbound[k] = nrItem
		}
		appRule.InboundAllowList = inbound
	}

	if ira, ok := d.GetOk("isolation_rule_action"); ok && ira.(string) != "" {
		isolationRule.Action = utils.String(ira.(string))
	}

	if f, fok := d.GetOk("isolation_rule_first_entity_filter_kind_list"); fok && f != nil {
		fkl := f.([]interface{})
		fkList := make([]*string, len(fkl))
		for k, v := range fkl {
			fkList[k] = utils.String(v.(string))
		}
		iRuleFirstEntityFilter.KindList = fkList
	}

	if ft, ftok := d.GetOk("isolation_rule_first_entity_filter_type"); ftok && ft.(string) != "" {
		iRuleFirstEntityFilter.Type = utils.String(ft.(string))
	}

	if fp, fpok := d.GetOk("isolation_rule_first_entity_filter_params"); fpok {
		fpl := fp.([]interface{})

		if len(fpl) > 0 {
			fl := make(map[string][]string)
			for _, v := range fpl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
						var values []string
						for _, item := range k.([]interface{}) {
							values = append(values, item.(string))
						}
						fl[i.(string)] = values
					}

				}
			}
			iRuleFirstEntityFilter.Params = fl
		} else {
			iRuleFirstEntityFilter.Params = nil
		}

	}

	if f, fok := d.GetOk("isolation_rule_second_entity_filter_kind_list"); fok && f != nil {
		fkl := f.([]interface{})
		fkList := make([]*string, len(fkl))
		for k, v := range fkl {
			fkList[k] = utils.String(v.(string))
		}
		iRuleSecondEntityFilter.KindList = fkList
	}

	if ft, ftok := d.GetOk("isolation_rule_second_entity_filter_type"); ftok && ft.(string) != "" {
		iRuleSecondEntityFilter.Type = utils.String(ft.(string))
	}

	if fp, fpok := d.GetOk("isolation_rule_second_entity_filter_params"); fpok {
		fpl := fp.([]interface{})

		if len(fpl) > 0 {
			fl := make(map[string][]string)
			for _, v := range fpl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
						var values []string
						for _, item := range k.([]interface{}) {
							values = append(values, item.(string))
						}
						fl[i.(string)] = values
					}

				}
			}
			iRuleSecondEntityFilter.Params = fl
		} else {
			iRuleSecondEntityFilter.Params = nil
		}

	}

	if !reflect.DeepEqual(*qRuleTargetGroupFilter, (v3.CategoryFilter{})) {
		qRuleTargetGroup.Filter = qRuleTargetGroupFilter
	}

	if !reflect.DeepEqual(*qRuleTargetGroup, (v3.TargetGroup{})) {
		quarantineRule.TargetGroup = qRuleTargetGroup
	}

	if !reflect.DeepEqual(*aRuleTargetGroupFilter, (v3.CategoryFilter{})) {
		aRuleTargetGroup.Filter = aRuleTargetGroupFilter
	}

	if !reflect.DeepEqual(*aRuleTargetGroup, (v3.TargetGroup{})) {
		appRule.TargetGroup = aRuleTargetGroup
	}

	if !reflect.DeepEqual(*quarantineRule, (v3.NetworkSecurityRuleResourcesRule{})) {
		networkSecurityRule.QuarantineRule = quarantineRule
	}

	if !reflect.DeepEqual(*appRule, (v3.NetworkSecurityRuleResourcesRule{})) {
		networkSecurityRule.AppRule = appRule
	}

	if !reflect.DeepEqual(*iRuleFirstEntityFilter, (v3.CategoryFilter{})) {
		isolationRule.FirstEntityFilter = iRuleFirstEntityFilter
	}

	if !reflect.DeepEqual(*iRuleSecondEntityFilter, (v3.CategoryFilter{})) {
		isolationRule.SecondEntityFilter = iRuleSecondEntityFilter
	}

	if !reflect.DeepEqual(*isolationRule, (v3.NetworkSecurityRuleIsolationRule{})) {
		networkSecurityRule.IsolationRule = isolationRule
	}
	return nil
}

func networkSecurityRuleStateRefreshFunc(client *v3.Client, uuid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.V3.GetNetworkSecurityRule(uuid)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return v, "DELETED", nil
			}
			log.Printf("ERROR %s", err)
			return nil, "", err
		}

		return v, *v.Status.State, nil
	}
}

func getNetworkSecurityRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_version": {
			Type:     schema.TypeString,
			Optional: true,
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
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"owner_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"project_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
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
		"quarantine_rule_action": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"quarantine_rule_outbound_allow_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"protocol": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_subnet": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_subnet_prefix_length": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"tcp_port_range_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"end_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"start_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"udp_port_range_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"end_port": {
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"start_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"filter_kind_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"filter_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"filter_params": {
						Type:     schema.TypeList,
						Optional: true,
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
						Optional: true,
						Computed: true,
					},

					"expiration_time": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Optional: true,
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
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"icmp_type_code_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"code": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Optional: true,
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
			Optional: true,
		},
		"quarantine_rule_target_group_peer_specification_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"quarantine_rule_target_group_filter_kind_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"quarantine_rule_target_group_filter_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"quarantine_rule_target_group_filter_params": {
			Type:     schema.TypeList,
			Optional: true,
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
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"protocol": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_subnet": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_subnet_prefix_length": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"tcp_port_range_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"end_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"start_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"udp_port_range_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"end_port": {
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"start_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"filter_kind_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"filter_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"filter_params": {
						Type:     schema.TypeList,
						Optional: true,
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
						Optional: true,
						Computed: true,
					},

					"expiration_time": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Optional: true,
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
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"icmp_type_code_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"code": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Optional: true,
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
			Optional: true,
			Computed: true,
		},
		"app_rule_outbound_allow_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"protocol": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_subnet": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_subnet_prefix_length": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"tcp_port_range_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"end_port": {
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"start_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"udp_port_range_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"end_port": {
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"start_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"filter_kind_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"filter_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"filter_params": {
						Type:     schema.TypeList,
						Optional: true,
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
						Optional: true,
						Computed: true,
					},

					"expiration_time": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Optional: true,
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
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"icmp_type_code_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"code": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Optional: true,
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
			Optional: true,
		},
		"app_rule_target_group_peer_specification_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"app_rule_target_group_filter_kind_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"app_rule_target_group_filter_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"app_rule_target_group_filter_params": {
			Type:     schema.TypeList,
			Optional: true,
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
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"protocol": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_subnet": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_subnet_prefix_length": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"tcp_port_range_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"end_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"start_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"udp_port_range_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"end_port": {
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"start_port": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"filter_kind_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"filter_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"filter_params": {
						Type:     schema.TypeList,
						Optional: true,
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
						Optional: true,
						Computed: true,
					},

					"expiration_time": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Optional: true,
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
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"icmp_type_code_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"code": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Optional: true,
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
			Optional: true,
			Computed: true,
		},
		"isolation_rule_first_entity_filter_kind_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"isolation_rule_first_entity_filter_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"isolation_rule_first_entity_filter_params": {
			Type:     schema.TypeList,
			Optional: true,
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
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"isolation_rule_second_entity_filter_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"isolation_rule_second_entity_filter_params": {
			Type:     schema.TypeList,
			Optional: true,
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
	}
}
