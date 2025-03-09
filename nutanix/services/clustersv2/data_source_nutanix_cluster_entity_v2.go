package clustersv2

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/response"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const defaultValue = 32

func DatasourceNutanixClusterEntityV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixClusterEntityV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number_of_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"controller_vm_ip": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLength(),
												"ipv6": SchemaForValuePrefixLength(),
											},
										},
									},
									"node_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_ip": {
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
					},
				},
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
						"external_data_services_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
						"external_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nfs_subnet_white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name_server_ip_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
									"fqdn": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"ntp_server_ip_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
									"fqdn": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"smtp_server": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_address": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLength(),
															"ipv6": SchemaForValuePrefixLength(),
															"fqdn": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
														},
													},
												},
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"username": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"password": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"masquerading_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
						"masquerading_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"management_server": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLength(),
												"ipv6": SchemaForValuePrefixLength(),
											},
										},
									},
									"type": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_drs_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_registered": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_in_use": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_management_server_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backplane": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_segmentation_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"vlan_tag": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"subnet": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"prefix_length": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"netmask": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"prefix_length": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"http_proxy_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": {
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
												"ipv6": {
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
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"proxy_types": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"http_proxy_white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"incarnation_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"build_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"build_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"full_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"commit_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"short_commit_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hypervisor_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cluster_function": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"authorized_public_key_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"redundancy_factor": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_software_map": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"software_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cluster_arch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fault_tolerance_state": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"current_max_fault_tolerance": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"desired_max_fault_tolerance": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"domain_awareness_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_cluster_fault_tolerance": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"desired_cluster_fault_tolerance": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"redundancy_status": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_cassandra_preparation_done": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_zookeeper_preparation_done": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"is_remote_support_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"operation_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_lts": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_password_remote_login_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"encryption_in_transit_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encryption_option": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"encryption_scope": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"pulse_status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"pii_scrubbing_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"is_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"upgrade_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inefficient_vm_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_profile_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_eligibility_score": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixClusterEntityV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	var expand *string

	extID := d.Get("ext_id")

	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}
	resp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(extID.(string)), expand)
	if err != nil {
		return diag.Errorf("error while fetching cluster entity : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Cluster)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nodes", flattenNodeReference(getResp.Nodes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network", flattenClusterNetworkReference(getResp.Network)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config", flattenClusterConfigReference(getResp.Config)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("upgrade_status", flattenUpgradeStatus(getResp.UpgradeStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_count", getResp.VmCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inefficient_vm_count", getResp.InefficientVmCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("container_name", getResp.ContainerName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("categories", getResp.Categories); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_profile_ext_id", getResp.ClusterProfileExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backup_eligibility_score", getResp.BackupEligibilityScore); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
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
					Default:  defaultValue,
				},
			},
		},
	}
}

func flattenLinks(pr []import3.ApiLink) []map[string]interface{} {
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

func flattenNodeReference(pr *import1.NodeReference) []map[string]interface{} {
	if pr != nil {
		nodeRef := make([]map[string]interface{}, 0)
		node := make(map[string]interface{})

		node["number_of_nodes"] = pr.NumberOfNodes
		node["node_list"] = flattenNodeListItemReference(pr.NodeList)

		nodeRef = append(nodeRef, node)
		return nodeRef
	}
	return nil
}

func flattenNodeListItemReference(pr []import1.NodeListItemReference) []interface{} {
	if len(pr) > 0 {
		nodeList := make([]interface{}, len(pr))

		for k, v := range pr {
			node := make(map[string]interface{})

			node["controller_vm_ip"] = flattenIPAddress(v.ControllerVmIp)
			node["node_uuid"] = v.NodeUuid
			node["host_ip"] = flattenIPAddress(v.HostIp)

			nodeList[k] = node
		}
		return nodeList
	}
	return nil
}

func flattenClusterNetworkReference(pr *import1.ClusterNetworkReference) []map[string]interface{} {
	if pr != nil {
		clsNet := make([]map[string]interface{}, 0)

		cls := make(map[string]interface{})

		cls["external_address"] = flattenIPAddress(pr.ExternalAddress)
		cls["external_data_services_ip"] = flattenIPAddress(pr.ExternalDataServiceIp)
		cls["external_subnet"] = pr.ExternalSubnet
		cls["internal_subnet"] = pr.InternalSubnet
		cls["nfs_subnet_white_list"] = pr.NfsSubnetWhitelist
		cls["name_server_ip_list"] = flattenIPAddressOrFQDN(pr.NameServerIpList)
		cls["ntp_server_ip_list"] = flattenIPAddressOrFQDN(pr.NtpServerIpList)
		cls["smtp_server"] = flattenSMTPServerRef(pr.SmtpServer)
		cls["masquerading_ip"] = flattenIPAddress(pr.MasqueradingIp)
		cls["masquerading_port"] = pr.MasqueradingPort
		cls["management_server"] = flattenManagementServerRef(pr.ManagementServer)
		cls["fqdn"] = pr.Fqdn
		cls["key_management_server_type"] = flattenKeyManagementServerType(pr.KeyManagementServerType)
		cls["backplane"] = flattenBackplaneNetworkParams(pr.Backplane)
		cls["http_proxy_list"] = flattenHTTPProxyList(pr.HttpProxyList)
		cls["http_proxy_white_list"] = flattenHTTPProxyWhiteList(pr.HttpProxyWhiteList)

		clsNet = append(clsNet, cls)
		return clsNet
	}
	return nil
}

func flattenHTTPProxyWhiteList(proxyWhiteList []import1.HttpProxyWhiteListConfig) []interface{} {
	if len(proxyWhiteList) > 0 {
		proxyList := make([]interface{}, len(proxyWhiteList))

		for k, v := range proxyWhiteList {
			proxy := make(map[string]interface{})

			proxy["target"] = v.Target
			proxy["target_type"] = flattenTargetType(v.TargetType)

			proxyList[k] = proxy
		}
		return proxyList
	}
	return nil
}

func flattenTargetType(targetType *import1.HttpProxyWhiteListTargetType) interface{} {
	if targetType != nil {
		const (
			ipv4Address, ipv6Address, ipv4NetworkMask, domainNameSuffix, hostName = 2, 3, 4, 5, 6
		)

		switch *targetType {
		case ipv6Address:
			return "IPV6_ADDRESS"
		case hostName:
			return "HOST_NAME"
		case domainNameSuffix:
			return "DOMAIN_NAME_SUFFIX"
		case ipv4NetworkMask:
			return "IPV4_NETWORK_MASK"
		case ipv4Address:
			return "IPV4_ADDRESS"
		}
	}
	return "UNKNOWN"
}

func flattenHTTPProxyList(httpProxyList []import1.HttpProxyConfig) []interface{} {
	if len(httpProxyList) > 0 {
		proxyList := make([]interface{}, len(httpProxyList))

		for k, v := range httpProxyList {
			proxy := make(map[string]interface{})

			proxy["ip_address"] = flattenIPAddress(v.IpAddress)
			proxy["port"] = v.Port
			proxy["username"] = v.Username
			proxy["password"] = v.Password
			proxy["name"] = v.Name
			proxy["proxy_types"] = flattenProxyTypes(v.ProxyTypes)

			proxyList[k] = proxy
		}
		return proxyList
	}
	return nil
}

func flattenProxyTypes(proxyTypes []import1.HttpProxyType) []interface{} {
	if len(proxyTypes) > 0 {
		types := make([]interface{}, len(proxyTypes))
		const (
			HTTP, HTTPS, SOCKS = 2, 3, 4
		)

		for k, v := range proxyTypes {
			switch v {
			case HTTP:
				types[k] = "HTTP"
			case HTTPS:
				types[k] = "HTTPS"
			case SOCKS:
				types[k] = "SOCKS"
			default:
				types[k] = "UNKNOWN"
			}
		}
		return types
	}
	return nil
}

func flattenClusterConfigReference(pr *import1.ClusterConfigReference) []map[string]interface{} {
	if pr != nil {
		clsConfig := make([]map[string]interface{}, 0)

		cls := make(map[string]interface{})

		cls["incarnation_id"] = pr.IncarnationId
		cls["build_info"] = flattenBuildReference(pr.BuildInfo)
		cls["hypervisor_types"] = flattenHypervisorType(pr.HypervisorTypes)
		cls["cluster_function"] = flattenClusterFunctionRef(pr.ClusterFunction)
		cls["timezone"] = pr.Timezone
		cls["authorized_public_key_list"] = flattenPublicKey(pr.AuthorizedPublicKeyList)
		cls["redundancy_factor"] = pr.RedundancyFactor
		cls["cluster_software_map"] = flattenSoftwareMapReference(pr.ClusterSoftwareMap)
		cls["cluster_arch"] = flattenClusterArchReference(pr.ClusterArch)
		cls["fault_tolerance_state"] = flattenFaultToleranceState(pr.FaultToleranceState)
		cls["is_remote_support_enabled"] = pr.IsRemoteSupportEnabled
		cls["operation_mode"] = flattenOperationMode(pr.OperationMode)
		cls["is_lts"] = pr.IsLts
		cls["is_password_remote_login_enabled"] = pr.IsPasswordRemoteLoginEnabled
		cls["encryption_in_transit_status"] = flattenEncryptionStatus(pr.EncryptionInTransitStatus)
		cls["encryption_option"] = flattenEncryptionOptionInfo(pr.EncryptionOption)
		cls["encryption_scope"] = flattenEncryptionScopeInfo(pr.EncryptionScope)
		cls["pulse_status"] = flattenPulseStatus(pr.PulseStatus)
		cls["is_available"] = pr.IsAvailable

		clsConfig = append(clsConfig, cls)
		return clsConfig
	}
	return nil
}

func flattenPulseStatus(status *import1.PulseStatus) []map[string]interface{} {
	if status != nil {
		pulse := make(map[string]interface{})

		pulse["is_enabled"] = status.IsEnabled
		pulse["pii_scrubbing_level"] = flattenPulseScrubbingLevel(status.PiiScrubbingLevel)

		return []map[string]interface{}{pulse}
	}
	return nil
}

func flattenPulseScrubbingLevel(level *import1.PIIScrubbingLevel) string {
	if level != nil {
		const DEFAULT, ALL = 2, 3

		switch *level {
		case DEFAULT:
			return "DEFAULT"
		case ALL:
			return "ALL"
		}
	}
	return "UNKNOWN"
}

func flattenIPAddress(pr *import4.IPAddress) []map[string]interface{} {
	if pr != nil {
		ips := make([]map[string]interface{}, 0)
		ip := make(map[string]interface{})

		ip["ipv4"] = flattenIPv4Address(pr.Ipv4)
		ip["ipv6"] = flattenIPv6Address(pr.Ipv6)

		ips = append(ips, ip)
		return ips
	}
	return nil
}

func flattenIPAddressOrFQDN(pr []import4.IPAddressOrFQDN) []map[string]interface{} {
	if len(pr) > 0 {
		ips := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ip := make(map[string]interface{})

			ip["ipv4"] = flattenIPv4Address(v.Ipv4)
			ip["ipv6"] = flattenIPv6Address(v.Ipv6)
			ip["fqdn"] = flattenFQDN(v.Fqdn)

			ips[k] = ip
		}
		return ips
	}
	return nil
}

func flattenSMTPServerRef(pr *import1.SmtpServerRef) []map[string]interface{} {
	if pr != nil {
		smtp := make([]map[string]interface{}, 0)
		s := make(map[string]interface{})

		s["email_address"] = pr.EmailAddress
		s["server"] = flattenSMTPNetwork(pr.Server)
		if pr.Type != nil {
			const PLAIN, STARTTLS, SSL = 2, 3, 4
			switch *pr.Type {
			case PLAIN:
				s["type"] = "PLAIN"
			case STARTTLS:
				s["type"] = "STARTTLS"
			case SSL:
				s["type"] = "SSL"
			default:
				s["type"] = "UNKNOWN"
			}
		}

		smtp = append(smtp, s)
		return smtp
	}
	return nil
}

func flattenBackplaneNetworkParams(pr *import1.BackplaneNetworkParams) []map[string]interface{} {
	if pr != nil {
		backplane := make([]map[string]interface{}, 0)

		back := make(map[string]interface{})

		back["is_segmentation_enabled"] = pr.IsSegmentationEnabled
		back["vlan_tag"] = pr.VlanTag
		back["subnet"] = flattenIPv4Address(pr.Subnet)
		back["netmask"] = flattenIPv4Address(pr.Netmask)

		backplane = append(backplane, back)
		return backplane
	}
	return nil
}

func flattenManagementServerRef(pr *import1.ManagementServerRef) []map[string]interface{} {
	if pr != nil {
		mgmServer := make([]map[string]interface{}, 0)

		mgm := make(map[string]interface{})

		mgm["ip"] = flattenIPAddress(pr.Ip)
		mgm["type"] = pr.Type
		mgm["is_drs_enabled"] = pr.IsDrsEnabled
		mgm["is_registered"] = pr.IsRegistered
		mgm["is_in_use"] = pr.IsInUse

		mgmServer = append(mgmServer, mgm)
		return mgmServer
	}
	return nil
}

func flattenSMTPNetwork(pr *import1.SmtpNetwork) []map[string]interface{} {
	if pr != nil {
		smtp := make([]map[string]interface{}, 0)

		s := make(map[string]interface{})
		ipAddressOrFQDN := make([]import4.IPAddressOrFQDN, 1)
		ipAddressOrFQDN[0] = *pr.IpAddress
		s["ip_address"] = flattenIPAddressOrFQDN(ipAddressOrFQDN)
		s["port"] = pr.Port
		s["username"] = pr.Username
		s["password"] = pr.Password

		smtp = append(smtp, s)
		return smtp
	}
	return nil
}

func flattenIPv4Address(pr *import4.IPv4Address) []interface{} {
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

func flattenIPv6Address(pr *import4.IPv6Address) []interface{} {
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

func flattenFQDN(pr *import4.FQDN) []interface{} {
	if pr != nil {
		fqdn := make([]interface{}, 0)

		f := make(map[string]interface{})

		f["value"] = pr.Value

		fqdn = append(fqdn, f)

		return fqdn
	}
	return nil
}

func flattenBuildReference(pr *import1.BuildReference) []map[string]interface{} {
	if pr != nil {
		buildRef := make([]map[string]interface{}, 0)
		build := make(map[string]interface{})

		build["build_type"] = pr.BuildType
		build["version"] = pr.Version
		build["full_version"] = pr.FullVersion
		build["commit_id"] = pr.CommitId
		build["short_commit_id"] = pr.ShortCommitId

		buildRef = append(buildRef, build)
		return buildRef
	}
	return nil
}

func flattenPublicKey(pr []import1.PublicKey) []map[string]interface{} {
	if len(pr) > 0 {
		pubKey := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			pub := make(map[string]interface{})

			pub["key"] = utils.StringValue(v.Key)
			pub["name"] = utils.StringValue(v.Name)

			pubKey[k] = pub
		}

		sort.SliceStable(pubKey, func(i, j int) bool {
			return pubKey[i]["key"].(string) < pubKey[j]["key"].(string)
		})

		return pubKey
	}
	return nil
}

func flattenSoftwareMapReference(pr []import1.SoftwareMapReference) []interface{} {
	if len(pr) > 0 {
		mapRef := make([]interface{}, len(pr))

		for k, v := range pr {
			m := make(map[string]interface{})

			m["software_type"] = flattenSoftwareTypeRef(v.SoftwareType)
			m["version"] = v.Version

			mapRef[k] = m
		}
		return mapRef
	}
	return nil
}

func flattenFaultToleranceState(pr *import1.FaultToleranceState) []map[string]interface{} {
	if pr != nil {
		faultTol := make([]map[string]interface{}, 0)

		fault := make(map[string]interface{})

		fault["current_max_fault_tolerance"] = pr.CurrentMaxFaultTolerance
		fault["desired_max_fault_tolerance"] = pr.DesiredMaxFaultTolerance
		fault["domain_awareness_level"] = flattenDomainAwarenessLevel(pr.DomainAwarenessLevel)
		fault["current_cluster_fault_tolerance"] = flattenClusterFaultTolerance(pr.CurrentClusterFaultTolerance)
		fault["desired_cluster_fault_tolerance"] = flattenClusterFaultTolerance(pr.DesiredClusterFaultTolerance)
		fault["redundancy_status"] = flattenRedundancyStatus(pr.RedundancyStatus)

		faultTol = append(faultTol, fault)
		return faultTol
	}
	return nil
}

func flattenRedundancyStatus(redundancyStatus *import1.RedundancyStatusDetails) []interface{} {
	if redundancyStatus != nil {
		redStatus := make(map[string]interface{})

		redStatus["is_cassandra_preparation_done"] = redundancyStatus.IsCassandraPreparationDone
		redStatus["is_zookeeper_preparation_done"] = redundancyStatus.IsZookeeperPreparationDone

		return []interface{}{redStatus}
	}
	return nil
}

func flattenClusterFaultTolerance(faultTolerance *import1.ClusterFaultToleranceRef) string {
	if faultTolerance != nil {
		const two, three, four, five = 2, 3, 4, 5

		switch *faultTolerance {
		case two:
			return "CFT_0N_AND_0D"
		case three:
			return "CFT_1N_OR_1D"
		case four:
			return "CFT_2N_OR_2D"
		case five:
			return "CFT_1N_AND_1D"
		default:
			return "UNKNOWN"
		}
	}
	return "UNKNOWN"
}

func flattenHypervisorType(pr []import1.HypervisorType) []string {
	if len(pr) > 0 {
		hyperTypes := make([]string, len(pr))
		const AHV, ESX, HYPERV, XEN, NATIVEHOST = 2, 3, 4, 5, 6

		for i, v := range pr {
			if v == import1.HypervisorType(AHV) {
				hyperTypes[i] = "AHV"
			}
			if v == import1.HypervisorType(ESX) {
				hyperTypes[i] = "ESX"
			}
			if v == import1.HypervisorType(HYPERV) {
				hyperTypes[i] = "HYPERV"
			}
			if v == import1.HypervisorType(XEN) {
				hyperTypes[i] = "XEN"
			}
			if v == import1.HypervisorType(NATIVEHOST) {
				hyperTypes[i] = "NATIVEHOST"
			}
		}
		return hyperTypes
	}
	return []string{"UNKNOWN"}
}

func flattenClusterFunctionRef(pr []import1.ClusterFunctionRef) []string {
	if len(pr) > 0 {
		clsFuncs := make([]string, len(pr))

		const two, three, four, five, six, seven, eight = 2, 3, 4, 5, 6, 7, 8
		for i, v := range pr {
			if v == import1.ClusterFunctionRef(two) {
				clsFuncs[i] = "AOS"
			}
			if v == import1.ClusterFunctionRef(three) {
				clsFuncs[i] = "PRISM_CENTRAL"
			}
			if v == import1.ClusterFunctionRef(four) {
				clsFuncs[i] = "CLOUD_DATA_GATEWAY"
			}
			if v == import1.ClusterFunctionRef(five) {
				clsFuncs[i] = "AFS"
			}
			if v == import1.ClusterFunctionRef(six) {
				clsFuncs[i] = "ONE_NODE"
			}
			if v == import1.ClusterFunctionRef(seven) {
				clsFuncs[i] = "TWO_NODE"
			}
			if v == import1.ClusterFunctionRef(eight) {
				clsFuncs[i] = "ANALYTICS_PLATFORM"
			}
		}
		return clsFuncs
	}
	return nil
}

func flattenClusterArchReference(pr *import1.ClusterArchReference) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == import1.ClusterArchReference(two) {
			return "X86_64"
		}
		if *pr == import1.ClusterArchReference(three) {
			return "PPC64LE"
		}
	}
	return "UNKNOWN"
}

func flattenOperationMode(pr *import1.OperationMode) string {
	if pr != nil {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		if *pr == import1.OperationMode(two) {
			return "NORMAL"
		}
		if *pr == import1.OperationMode(three) {
			return "READ_ONLY"
		}
		if *pr == import1.OperationMode(four) {
			return "STAND_ALONE"
		}
		if *pr == import1.OperationMode(five) {
			return "SWITCH_TO_TWO_NODE"
		}
		if *pr == import1.OperationMode(six) {
			return "OVERRIDE"
		}
	}
	return "UNKNOWN"
}

func flattenEncryptionStatus(pr *import1.EncryptionStatus) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == import1.EncryptionStatus(two) {
			return "ENABLED"
		}
		if *pr == import1.EncryptionStatus(three) {
			return "DISABLED"
		}
	}
	return "UNKNOWN"
}

func flattenEncryptionOptionInfo(pr []import1.EncryptionOptionInfo) []string {
	if len(pr) > 0 {
		enInfo := make([]string, len(pr))
		const two, three, four = 2, 3, 4

		for i, v := range pr {
			if v == import1.EncryptionOptionInfo(two) {
				enInfo[i] = "SOFTWARE"
			}
			if v == import1.EncryptionOptionInfo(three) {
				enInfo[i] = "HARDWARE"
			}
			if v == import1.EncryptionOptionInfo(four) {
				enInfo[i] = "SOFTWARE_AND_HARDWARE"
			}
		}
		return enInfo
	}
	return nil
}

func flattenEncryptionScopeInfo(pr []import1.EncryptionScopeInfo) []string {
	if len(pr) > 0 {
		enScope := make([]string, len(pr))
		const two, three = 2, 3

		for i, v := range pr {
			if v == import1.EncryptionScopeInfo(two) {
				enScope[i] = "CLUSTER"
			}
			if v == import1.EncryptionScopeInfo(three) {
				enScope[i] = "CONTAINER"
			}
		}
		return enScope
	}
	return nil
}

func flattenKeyManagementServerType(pr *import1.KeyManagementServerType) string {
	if pr != nil {
		const two, three, four = 2, 3, 4

		if *pr == import1.KeyManagementServerType(two) {
			return "LOCAL"
		}
		if *pr == import1.KeyManagementServerType(three) {
			return "PRISM_CENTRAL"
		}
		if *pr == import1.KeyManagementServerType(four) {
			return "EXTERNAL"
		}
	}
	return "UNKNOWN"
}

func flattenDomainAwarenessLevel(pr *import1.DomainAwarenessLevel) string {
	if pr != nil {
		const two, three, four, five = 2, 3, 4, 5
		if *pr == import1.DomainAwarenessLevel(two) {
			return "NODE"
		}
		if *pr == import1.DomainAwarenessLevel(three) {
			return "BLOCK"
		}
		if *pr == import1.DomainAwarenessLevel(four) {
			return "RACK"
		}
		if *pr == import1.DomainAwarenessLevel(five) {
			return "DISK"
		}
	}
	return "UNKNOWN"
}

func flattenUpgradeStatus(pr *import1.UpgradeStatus) string {
	if pr != nil {
		const two, three, four, five, six, seven, eight, nine, ten = 2, 3, 4, 5, 6, 7, 8, 9, 10

		if *pr == import1.UpgradeStatus(two) {
			return "PENDING"
		}
		if *pr == import1.UpgradeStatus(three) {
			return "DOWNLOADING"
		}
		if *pr == import1.UpgradeStatus(four) {
			return "QUEUED"
		}
		if *pr == import1.UpgradeStatus(five) {
			return "PREUPGRADE"
		}
		if *pr == import1.UpgradeStatus(six) {
			return "UPGRADING"
		}
		if *pr == import1.UpgradeStatus(seven) {
			return "SUCCEEDED"
		}
		if *pr == import1.UpgradeStatus(eight) {
			return "FAILED"
		}
		if *pr == import1.UpgradeStatus(nine) {
			return "CANCELED"
		}
		if *pr == import1.UpgradeStatus(ten) {
			return "SCHEDULED"
		}
	}
	return "UNKNOWN"
}

func flattenSoftwareTypeRef(pr *import1.SoftwareTypeRef) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == import1.SoftwareTypeRef(two) {
			return "NOS"
		}
		if *pr == import1.SoftwareTypeRef(three) {
			return "NCC"
		}
		if *pr == import1.SoftwareTypeRef(four) {
			return "PRISM_CENTRAL"
		}
	}
	return "UNKNOWN"
}
