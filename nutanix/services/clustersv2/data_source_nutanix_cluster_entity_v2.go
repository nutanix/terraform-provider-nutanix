package clustersv2

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const ipv4PrefixLengthDefaultValue = 32

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
			"links": common.LinksSchema(),
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
										Elem:     common.SchemaForIPList(false),
									},
									"node_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_ip": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     common.SchemaForIPList(false),
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
							Elem:     common.SchemaForIPList(false),
						},
						"external_data_services_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
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
							Elem:     common.SchemaForIPList(true),
						},
						"ntp_server_ip_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(true),
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
													Elem:     common.SchemaForIPList(true),
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
							Elem:     common.SchemaForIPList(false),
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
										Elem:     common.SchemaForIPList(false),
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
													Type:     schema.TypeInt,
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
													Type:     schema.TypeInt,
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
										Elem:     common.SchemaForIPList(false),
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
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
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

	d.SetId(utils.StringValue(getResp.ExtId))
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

func flattenTargetType(targetType *import1.HttpProxyWhiteListTargetType) string {
	return common.FlattenPtrEnum(targetType)
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

func flattenProxyTypes(proxyTypes []import1.HttpProxyType) []string {
	return common.FlattenEnumValueList(proxyTypes)
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
	return common.FlattenPtrEnum(level)
}

func flattenIPAddress(addr *import4.IPAddress) []map[string]interface{} {
	if addr == nil {
		return nil
	}

	ipMap := map[string]interface{}{
		"ipv4": flattenIPv4Address(addr.Ipv4),
		"ipv6": flattenIPv6Address(addr.Ipv6),
	}

	return []map[string]interface{}{ipMap}
}

func flattenIPAddressOrFQDN(addrs []import4.IPAddressOrFQDN) []map[string]interface{} {
	if len(addrs) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(addrs))
	for i, addr := range addrs {
		ipMap := map[string]interface{}{
			"ipv4": flattenIPv4Address(addr.Ipv4),
			"ipv6": flattenIPv6Address(addr.Ipv6),
			"fqdn": flattenFQDN(addr.Fqdn),
		}
		result[i] = ipMap
	}
	return result
}

func flattenSMTPServerRef(smtpServerRef *import1.SmtpServerRef) []map[string]interface{} {
	if smtpServerRef == nil {
		return nil
	}

	smtpRef := map[string]interface{}{
		"email_address": utils.StringValue(smtpServerRef.EmailAddress),
		"server":        flattenSMTPNetwork(smtpServerRef.Server),
	}
	smtpRef["type"] = common.FlattenPtrEnum(smtpServerRef.Type)

	return []map[string]interface{}{smtpRef}
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

func flattenIPv4Address(ipv4Address *import4.IPv4Address) []interface{} {
	if ipv4Address == nil {
		return nil
	}

	ip := map[string]interface{}{
		"value":         ipv4Address.Value,
		"prefix_length": ipv4Address.PrefixLength,
	}

	return []interface{}{ip}
}

func flattenIPv6Address(ipv6Address *import4.IPv6Address) []interface{} {
	if ipv6Address == nil {
		return nil
	}

	ip := map[string]interface{}{
		"value":         ipv6Address.Value,
		"prefix_length": ipv6Address.PrefixLength,
	}

	return []interface{}{ip}
}

func flattenFQDN(pr *import4.FQDN) []interface{} {
	if pr == nil {
		return nil
	}

	f := map[string]interface{}{
		"value": pr.Value,
	}

	return []interface{}{f}
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

func flattenRedundancyStatus(redundancyStatus *import1.RedundancyStatusDetails) []map[string]interface{} {
	if redundancyStatus != nil {
		redStatus := make(map[string]interface{})

		redStatus["is_cassandra_preparation_done"] = redundancyStatus.IsCassandraPreparationDone
		redStatus["is_zookeeper_preparation_done"] = redundancyStatus.IsZookeeperPreparationDone

		return []map[string]interface{}{redStatus}
	}
	return nil
}

func flattenClusterFaultTolerance(faultTolerance *import1.ClusterFaultToleranceRef) string {
	return common.FlattenPtrEnum(faultTolerance)
}

func flattenHypervisorType(hypervisorTypes []import1.HypervisorType) []string {
	return common.FlattenEnumValueList(hypervisorTypes)
}

func flattenClusterFunctionRef(clusterFunctionRefs []import1.ClusterFunctionRef) []string {
	return common.FlattenEnumValueList(clusterFunctionRefs)
}

func flattenClusterArchReference(clusterArchReference *import1.ClusterArchReference) string {
	return common.FlattenPtrEnum(clusterArchReference)
}

func flattenOperationMode(operationMode *import1.OperationMode) string {
	return common.FlattenPtrEnum(operationMode)
}

func flattenEncryptionStatus(EncryptionStatus *import1.EncryptionStatus) string {
	return common.FlattenPtrEnum(EncryptionStatus)
}

func flattenEncryptionOptionInfo(encryptionOptionInfos []import1.EncryptionOptionInfo) []string {
	return common.FlattenEnumValueList(encryptionOptionInfos)
}

func flattenEncryptionScopeInfo(encryptionScopeInfos []import1.EncryptionScopeInfo) []string {
	return common.FlattenEnumValueList(encryptionScopeInfos)
}

func flattenKeyManagementServerType(keyManagementServerType *import1.KeyManagementServerType) string {
	return common.FlattenPtrEnum(keyManagementServerType)
}

func flattenDomainAwarenessLevel(domainAwarenessLevel *import1.DomainAwarenessLevel) string {
	return common.FlattenPtrEnum(domainAwarenessLevel)
}

func flattenUpgradeStatus(upgradeStatus *import1.UpgradeStatus) string {
	return common.FlattenPtrEnum(upgradeStatus)
}

func flattenSoftwareTypeRef(softwareTypeRef *import1.SoftwareTypeRef) string {
	return common.FlattenPtrEnum(softwareTypeRef)
}
