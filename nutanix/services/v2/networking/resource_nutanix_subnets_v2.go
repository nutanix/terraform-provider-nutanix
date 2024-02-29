package networking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/models/common/v1/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/models/networking/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/models/prism/v4/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixSubnetv4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixSubnetv4Create,
		ReadContext:   ResourceNutanixSubnetv4Read,
		UpdateContext: ResourceNutanixSubnetv4Update,
		DeleteContext: ResourceNutanixSubnetv4Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OVERLAY", "VLAN"}, false),
			},
			"network_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dhcp_options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
						"domain_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"search_domains": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tftp_server_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"boot_file_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ntp_servers": {
							Type:     schema.TypeList,
							Optional: true,
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_subnet": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": SchemaForValuePrefixLength(),
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"default_gateway_ip":  SchemaForValuePrefixLength(),
									"dhcp_server_address": SchemaForValuePrefixLength(),
									"pool_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"end_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
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
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_subnet": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": SchemaForValuePrefixLength(),
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"default_gateway_ip":  SchemaForValuePrefixLength(),
									"dhcp_server_address": SchemaForValuePrefixLength(),
									"pool_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"end_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
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
				Optional: true,
			},
			"virtual_switch_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_nat_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_external": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reserved_ip_addresses": SchemaForValuePrefixLength(),
			"dynamic_ip_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(),
						"ipv6": SchemaForValuePrefixLength(),
					},
				},
			},
			"network_function_chain_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bridge_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_advanced_networking": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hypervisor_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual_switch": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: DataSourceVirtualSwitchSchemaV4(),
				},
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: DataSourceVPCSchemaV4(),
				},
			},
			"ip_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_usage": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num_macs": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"num_free_ips": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"num_assigned_ips": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ip_pool_usages": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_free_ips": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"num_total_ips": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"range": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"end_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixSubnetv4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	inputSpec := import1.Subnet{}

	if name, nok := d.GetOk("name"); nok {
		inputSpec.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		inputSpec.Description = utils.StringPtr(desc.(string))
	}
	if subType, ok := d.GetOk("subnet_type"); ok {
		subMap := map[string]interface{}{
			"OVERLAY": 2,
			"VLAN":    3,
		}
		pInt := subMap[subType.(string)]
		p := import1.SubnetType(pInt.(int))
		inputSpec.SubnetType = &p
	}

	if dhcp, ok := d.GetOk("dhcp_options"); ok {
		inputSpec.DhcpOptions = expandDhcpOptions(dhcp.([]interface{}))
	}
	if clsRef, ok := d.GetOk("cluster_reference"); ok {
		inputSpec.ClusterReference = utils.StringPtr(clsRef.(string))
	}
	if vsRef, ok := d.GetOk("virtual_switch_reference"); ok {
		inputSpec.VirtualSwitchReference = utils.StringPtr(vsRef.(string))
	}
	if vpcRef, ok := d.GetOk("vpc_reference"); ok {
		inputSpec.VirtualSwitchReference = utils.StringPtr(vpcRef.(string))
	}
	if isNat, ok := d.GetOk("is_nat_enabled"); ok {
		inputSpec.IsNatEnabled = utils.BoolPtr(isNat.(bool))
	}
	if isExt, ok := d.GetOk("is_external"); ok {
		inputSpec.IsExternal = utils.BoolPtr(isExt.(bool))
	}
	if reservedIPAdd, ok := d.GetOk("reserved_ip_addresses"); ok {
		inputSpec.ReservedIpAddresses = expandIPAddress(reservedIPAdd.([]interface{}))
	}
	if dynamicIPAdd, ok := d.GetOk("dynamic_ip_addresses"); ok {
		inputSpec.DynamicIpAddresses = expandIPAddress(dynamicIPAdd.([]interface{}))
	}
	if ntwfuncRef, ok := d.GetOk("network_function_chain_reference"); ok {
		inputSpec.NetworkFunctionChainReference = utils.StringPtr(ntwfuncRef.(string))
	}
	if bridgeName, ok := d.GetOk("bridge_name"); ok {
		inputSpec.BridgeName = utils.StringPtr(bridgeName.(string))
	}
	if isAdvNet, ok := d.GetOk("is_advanced_networking"); ok {
		inputSpec.IsAdvancedNetworking = utils.BoolPtr(isAdvNet.(bool))
	}
	if clsName, ok := d.GetOk("cluster_name"); ok {
		inputSpec.ClusterName = utils.StringPtr(clsName.(string))
	}
	if hypervisorType, ok := d.GetOk("hypervisor_type"); ok {
		inputSpec.HypervisorType = utils.StringPtr(hypervisorType.(string))
	}
	if vswitch, ok := d.GetOk("virtual_switch"); ok {
		inputSpec.VirtualSwitch = expandVirtualSwitch(vswitch)
	}
	if vpc, ok := d.GetOk("vpc"); ok {
		inputSpec.Vpc = expandVpc(vpc)
	}
	if ipPrefix, ok := d.GetOk("ip_prefix"); ok {
		inputSpec.IpPrefix = utils.StringPtr(ipPrefix.(string))
	}
	if ipUsage, ok := d.GetOk("ip_usage"); ok {
		inputSpec.IpUsage = exapndIPUsage(ipUsage)
	}

	resp, err := conn.SubnetApiInstance.CreateSubnet(&inputSpec)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while creating subnets : %v", errorMessage["message"])
	}

	getResp := resp.Data.GetValue().(import4.TaskReference)

	fmt.Println(getResp)
	return nil
}

func ResourceNutanixSubnetv4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixSubnetv4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixSubnetv4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandDhcpOptions(pr []interface{}) *import1.DhcpOptions {
	if len(pr) > 0 {
		dhcpOps := import1.DhcpOptions{}

		val := pr[0].(map[string]interface{})

		if bootfn, ok := val["boot_file_name"]; ok {
			dhcpOps.BootFileName = utils.StringPtr(bootfn.(string))
		}
		if dns, ok := val["domain_name_servers"]; ok {
			dhcpOps.DomainNameServers = expandIPAddress(dns.([]interface{}))
		}
		if dn, ok := val["domain_name"]; ok {
			dhcpOps.DomainName = utils.StringPtr(dn.(string))
		}
		if searchDomain, ok := val["search_domains"]; ok {
			dhcpOps.SearchDomains = utils.StringValueSlice(searchDomain.([]*string))
		}
		if tftp, ok := val["tftp_server_name"]; ok {
			dhcpOps.TftpServerName = utils.StringPtr(tftp.(string))
		}
		if ntp, ok := val["ntp_servers"]; ok {
			dhcpOps.NtpServers = expandIPAddress(ntp.([]interface{}))
		}
		return &dhcpOps
	}
	return nil
}

func expandIPAddress(pr []interface{}) []config.IPAddress {
	if len(pr) > 0 {
		configList := make([]config.IPAddress, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			config := config.IPAddress{}

			if ipv4, ok := val["ipv4"]; ok {
				config.Ipv4 = expandIPv4Address(ipv4)
			}
			if ipv6, ok := val["ipv4"]; ok {
				config.Ipv6 = expandIPv6Address(ipv6)
			}

			configList[k] = config
		}
		return configList
	}
	return nil
}

func expandIPv4Address(pr interface{}) *config.IPv4Address {
	if pr != nil {
		ipv4 := &config.IPv4Address{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv4.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv4
	}
	return nil
}

func expandIPv6Address(pr interface{}) *config.IPv6Address {
	if pr != nil {
		ipv6 := &config.IPv6Address{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv6.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv6.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv6
	}
	return nil
}

func expandVirtualSwitch(pr interface{}) *import1.VirtualSwitch {
	if pr != nil {
		vSwitch := &import1.VirtualSwitch{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if extID, ok := val["ext_id"]; ok {
			vSwitch.ExtId = utils.StringPtr(extID.(string))
		}
		if name, ok := val["name"]; ok {
			vSwitch.Name = utils.StringPtr(name.(string))
		}
		if desc, ok := val["description"]; ok {
			vSwitch.Description = utils.StringPtr(desc.(string))
		}
		if isDefault, ok := val["is_default"]; ok {
			vSwitch.IsDefault = utils.BoolPtr(isDefault.(bool))
		}
		if hasDepErr, ok := val["has_deployment_error"]; ok {
			vSwitch.HasDeploymentError = utils.BoolPtr(hasDepErr.(bool))
		}
		if mtu, ok := val["mtu"]; ok {
			vSwitch.Mtu = utils.Int64Ptr(mtu.(int64))
		}
		if bondMode, ok := val["bond_mode"]; ok {
			bondMap := map[string]interface{}{
				"ACTIVE_BACKUP": 2,
				"BALANCE_SLB":   3,
				"BALANCE_TCP":   4,
				"NONE":          5,
			}
			pInt := bondMap[bondMode.(string)]
			p := import1.BondModeType(pInt.(int))
			vSwitch.BondMode = &p
		}
		if cls, ok := val["clusters"]; ok {
			vSwitch.Clusters = expandCluster(cls.([]interface{}))
		}
		if name, ok := val["name"]; ok {
			vSwitch.Name = utils.StringPtr(name.(string))
		}
	}
	return nil
}

func expandCluster(pr []interface{}) []import1.Cluster {
	if len(pr) > 0 {
		clsList := make([]import1.Cluster, len(pr))

		for k, v := range pr {
			cls := import1.Cluster{}
			val := v.(map[string]interface{})

			if extID, ok := val["ext_id"]; ok {
				cls.ExtId = utils.StringPtr(extID.(string))
			}
			if hosts, ok := val["hosts"]; ok {
				cls.Hosts = expandHost(hosts.([]interface{}))
			}
			if gateway, ok := val["gateway_ip_address"]; ok {
				cls.GatewayIpAddress = expandIPv4Address(gateway)
			}
			clsList[k] = cls
		}
		return clsList
	}
	return nil
}

func expandHost(pr []interface{}) []import1.Host {
	if len(pr) > 0 {
		hosts := make([]import1.Host, len(pr))

		for k, v := range pr {
			host := import1.Host{}
			val := v.(map[string]interface{})

			if extID, ok := val["ext_id"]; ok {
				host.ExtId = utils.StringPtr(extID.(string))
			}
			if hostNics, ok := val["host_nics"]; ok {
				host.HostNics = utils.StringValueSlice(hostNics.([]*string))
			}
			if ipAdd, ok := val["ip_address"]; ok {
				host.IpAddress = expandIPv4Subnet(ipAdd)
			}

			hosts[k] = host
		}
		return hosts
	}
	return nil
}

func expandIPv4Subnet(pr interface{}) *import1.IPv4Subnet {
	if pr != nil {
		ipv4Subs := &import1.IPv4Subnet{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ip, ok := val["ip"]; ok {
			ipv4Subs.Ip = expandIPv4Address(ip)
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4Subs.PrefixLength = utils.IntPtr(prefix.(int))
		}

		return ipv4Subs
	}
	return nil
}

func expandIPv6Subnet(pr interface{}) *import1.IPv6Subnet {
	if pr != nil {
		ipv6Subs := &import1.IPv6Subnet{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ip, ok := val["ip"]; ok {
			ipv6Subs.Ip = expandIPv6Address(ip)
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv6Subs.PrefixLength = utils.IntPtr(prefix.(int))
		}

		return ipv6Subs
	}
	return nil
}

func expandVpc(pr interface{}) *import1.Vpc {
	if pr != nil {
		vpc := &import1.Vpc{}
		prI := pr.([]interface{})

		val := prI[0].(map[string]interface{})

		if ext, ok := val["ext_id"]; ok {
			vpc.ExtId = utils.StringPtr(ext.(string))
		}
		if vpcType, ok := val["vpc_type"]; ok {
			vpcMap := map[string]interface{}{
				"REGULAR": 2,
				"TRANSIT": 3,
			}
			pInt := vpcMap[vpcType.(string)]
			p := import1.VpcType(pInt.(int))
			vpc.VpcType = &p
		}
		if desc, ok := val["description"]; ok {
			vpc.Description = utils.StringPtr(desc.(string))
		}
		if dhcpOps, ok := val["common_dhcp_options"]; ok {
			vpc.CommonDhcpOptions = expandVpcDhcpOptions(dhcpOps)
		}
		if extSubs, ok := val["external_subnets"]; ok {
			vpc.ExternalSubnets = expandExternalSubnet(extSubs.([]interface{}))
		}
		if extRoutingDomainRef, ok := val["external_routing_domain_reference"]; ok {
			vpc.ExternalRoutingDomainReference = utils.StringPtr(extRoutingDomainRef.(string))
		}
		if extRoutablePrefix, ok := val["externally_routable_prefixes"]; ok {
			vpc.ExternallyRoutablePrefixes = expandIPSubnet(extRoutablePrefix.([]interface{}))
		}

	}
	return nil
}

func expandVpcDhcpOptions(pr interface{}) *import1.VpcDhcpOptions {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		vpc := &import1.VpcDhcpOptions{}

		if dns, ok := val["domain_name_servers"]; ok {
			vpc.DomainNameServers = expandIPAddress(dns.([]interface{}))
		}
		return vpc
	}
	return nil
}

func expandExternalSubnet(pr []interface{}) []import1.ExternalSubnet {
	if len(pr) > 0 {
		extSubs := make([]import1.ExternalSubnet, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			sub := import1.ExternalSubnet{}

			if subRef, ok := val["subnet_reference"]; ok {
				sub.SubnetReference = utils.StringPtr(subRef.(string))
			}
			if extips, ok := val["external_ips"]; ok {
				sub.ExternalIps = expandIPAddress(extips.([]interface{}))
			}
			if gatewayNodes, ok := val["gateway_nodes"]; ok {
				sub.GatewayNodes = utils.StringValueSlice(gatewayNodes.([]*string))
			}
			if activeGatewayNode, ok := val["active_gateway_node"]; ok {
				sub.ActiveGatewayNode = expandGatewayNodeReference(activeGatewayNode)
			}
			extSubs[k] = sub
		}
		return extSubs
	}
	return nil
}

func expandGatewayNodeReference(pr interface{}) *import1.GatewayNodeReference {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		gatewayNodesRef := &import1.GatewayNodeReference{}

		if nodeID, ok := val["node_id"]; ok {
			gatewayNodesRef.NodeId = utils.StringPtr(nodeID.(string))
		}
		if nodeipAdd, ok := val["node_ip_address"]; ok {
			gatewayNodesRef.NodeIpAddress = expandIPAddressMap(nodeipAdd)
		}

		return gatewayNodesRef
	}
	return nil
}

func expandIPAddressMap(pr interface{}) *config.IPAddress {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		ipAdd := &config.IPAddress{}

		if ipv4, ok := val["ipv4"]; ok {
			ipAdd.Ipv4 = expandIPv4AddressMap(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok {
			ipAdd.Ipv6 = expandIPv6AddressMap(ipv6)
		}

		return ipAdd
	}
	return nil
}

func expandIPv4AddressMap(pr interface{}) *config.IPv4Address {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		ipv4Add := &config.IPv4Address{}

		if value, ok := val["value"]; ok {
			ipv4Add.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4Add.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv4Add
	}
	return nil
}

func expandIPv6AddressMap(pr interface{}) *config.IPv6Address {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		ipv6Add := &config.IPv6Address{}

		if value, ok := val["value"]; ok {
			ipv6Add.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv6Add.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv6Add
	}
	return nil
}

func expandIPSubnet(pr []interface{}) []import1.IPSubnet {
	if len(pr) > 0 {
		ips := make([]import1.IPSubnet, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			ip := import1.IPSubnet{}

			if ipv4, ok := val["ipv4"]; ok {
				ip.Ipv4 = expandIPv4Subnet(ipv4)
			}
			if ipv6, ok := val["ipv6"]; ok {
				ip.Ipv6 = expandIPv6Subnet(ipv6)
			}
			ips[k] = ip
		}

		return ips
	}
	return nil
}

func exapndIPUsage(pr interface{}) *import1.IPUsage {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		ipUsage := &import1.IPUsage{}

		if numMacs, ok := val["num_macs"]; ok {
			ipUsage.NumMacs = utils.Int64Ptr(numMacs.(int64))
		}
		if numFreeIPS, ok := val["num_free_ips"]; ok {
			ipUsage.NumFreeIPs = utils.Int64Ptr(numFreeIPS.(int64))
		}
		if numAssgIPs, ok := val["num_assigned_ips"]; ok {
			ipUsage.NumAssignedIPs = utils.Int64Ptr(numAssgIPs.(int64))
		}
		return ipUsage

	}
	return nil
}
