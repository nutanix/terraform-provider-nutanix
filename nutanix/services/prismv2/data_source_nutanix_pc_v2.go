package prismv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixFetchPcV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixPcV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links":  schemaForLinks(),
			"config": schemaForReadPcConfig(),
			"is_registered_with_hosting_cluster": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network": schemaForReadPcNetwork(),
			"hosting_cluster_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"should_enable_high_availability": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"node_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func DatasourceNutanixPcV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	pcExtID := d.Get("ext_id").(string)
	resp, err := conn.DomainManagerAPIInstance.GetDomainManagerById(utils.StringPtr(pcExtID))

	if err != nil {
		return diag.Errorf("error while fetching Domain Manager Configuration Detail: %s", err)
	}

	deployPcBody := resp.Data.GetValue().(config.DomainManager)

	if err := d.Set("tenant_id", utils.StringValue(deployPcBody.TenantId)); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(deployPcBody.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("config", flattenPCConfig(deployPcBody.Config)); err != nil {
		return diag.Errorf("error setting config: %s", err)
	}
	if err := d.Set("is_registered_with_hosting_cluster", utils.BoolValue(deployPcBody.IsRegisteredWithHostingCluster)); err != nil {
		return diag.Errorf("error setting is_registered_with_hosting_cluster: %s", err)
	}
	if err := d.Set("network", flattenPCNetwork(deployPcBody.Network)); err != nil {
		return diag.Errorf("error setting network: %s", err)
	}
	if err := d.Set("hosting_cluster_ext_id", utils.StringValue(deployPcBody.HostingClusterExtId)); err != nil {
		return diag.Errorf("error setting hosting_cluster_ext_id: %s", err)
	}
	if err := d.Set("should_enable_high_availability", utils.BoolValue(deployPcBody.ShouldEnableHighAvailability)); err != nil {
		return diag.Errorf("error setting should_enable_high_availability: %s", err)
	}
	if err := d.Set("node_ext_ids", deployPcBody.NodeExtIds); err != nil {
		return diag.Errorf("error setting node_ext_ids: %s", err)
	}

	d.SetId(utils.StringValue(deployPcBody.ExtId))

	return nil
}

// schemas

func schemaForReadPcConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"should_enable_lockdown_mode": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"build_info": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"version": {
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
				"size": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"bootstrap_config": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"environment_info": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"type": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"provider_type": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"provisioning_type": {
											Type:     schema.TypeString,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
				"resource_config": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"container_ext_ids": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"num_vcpus": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"memory_size_bytes": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"data_disk_size_bytes": {
								Type:     schema.TypeInt,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func schemaForReadPcNetwork() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"external_address": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     schemaForIPAddress(),
				},
				"name_servers": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     schemaForIPAddressOrFqdn(),
				},
				"ntp_servers": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     schemaForIPAddressOrFqdn(),
				},
				"fqdn": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"external_networks": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"default_gateway": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     schemaForIPAddressOrFqdn(),
							},
							"subnet_mask": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     schemaForIPAddressOrFqdn(),
							},
							"ip_ranges": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"begin": {
											Type:     schema.TypeList,
											Computed: true,
											Elem:     schemaForIPAddress(),
										},
										"end": {
											Type:     schema.TypeList,
											Computed: true,
											Elem:     schemaForIPAddress(),
										},
									},
								},
							},
							"network_ext_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

// flattens

func flattenPCConfig(pcConfig *config.DomainManagerClusterConfig) []map[string]interface{} {
	if pcConfig != nil {
		pcConfigMap := make(map[string]interface{})
		if pcConfig.ShouldEnableLockdownMode != nil {
			pcConfigMap["should_enable_lockdown_mode"] = *pcConfig.ShouldEnableLockdownMode
		}
		if pcConfig.BuildInfo != nil {
			buildInfo := make(map[string]interface{})
			if pcConfig.BuildInfo.Version != nil {
				buildInfo["version"] = utils.StringValue(pcConfig.BuildInfo.Version)
			}
			pcConfigMap["build_info"] = []map[string]interface{}{buildInfo}
		}
		if pcConfig.Name != nil {
			pcConfigMap["name"] = *pcConfig.Name
		}
		if pcConfig.Size != nil {
			pcConfigMap["size"] = flattenClusterSize(*pcConfig.Size)
		}
		if pcConfig.BootstrapConfig != nil {
			pcConfigMap["bootstrap_config"] = flattenBootstrapConfig(pcConfig.BootstrapConfig)
		}
		if pcConfig.ResourceConfig != nil {
			pcConfigMap["resource_config"] = flattenResourceConfig(pcConfig.ResourceConfig)
		}
		return []map[string]interface{}{pcConfigMap}
	}
	return nil
}

func flattenClusterSize(size config.Size) string {
	const STARTER, SMALL, LARGE, EXTRALARGE = 2, 3, 4, 5

	switch size {
	case STARTER:
		return "STARTER"
	case SMALL:
		return "SMALL"
	case LARGE:
		return "LARGE"
	case EXTRALARGE:
		return "EXTRALARGE"
	default:
		return "UNKNOWN"
	}
}

func flattenBootstrapConfig(bootstrapConfig *config.BootstrapConfig) []map[string]interface{} {
	if bootstrapConfig != nil {
		bootstrapConfigMap := make(map[string]interface{})
		if bootstrapConfig.EnvironmentInfo != nil {
			bootstrapConfigMap["environment_info"] = flattenEnvironmentInfo(bootstrapConfig.EnvironmentInfo)
		}
		return []map[string]interface{}{bootstrapConfigMap}
	}
	return nil
}

func flattenResourceConfig(resourceConfig *config.DomainManagerResourceConfig) []map[string]interface{} {
	if resourceConfig != nil {
		resourceConfigMap := make(map[string]interface{})
		if resourceConfig.NumVcpus != nil {
			resourceConfigMap["num_vcpus"] = utils.IntValue(resourceConfig.NumVcpus)
		}
		if resourceConfig.MemorySizeBytes != nil {
			resourceConfigMap["memory_size_bytes"] = utils.Int64Value(resourceConfig.MemorySizeBytes)
		}
		if resourceConfig.DataDiskSizeBytes != nil {
			resourceConfigMap["data_disk_size_bytes"] = utils.Int64Value(resourceConfig.DataDiskSizeBytes)
		}
		if resourceConfig.ContainerExtIds != nil {
			resourceConfigMap["container_ext_ids"] = resourceConfig.ContainerExtIds
		}
		return []map[string]interface{}{resourceConfigMap}
	}
	return nil
}

func flattenEnvironmentInfo(environmentInfo *config.EnvironmentInfo) []map[string]interface{} {
	if environmentInfo != nil {
		info := make(map[string]interface{})
		if environmentInfo.Type != nil {
			info["type"] = flattenEnvironmentType(*environmentInfo.Type)
		}
		if environmentInfo.ProviderType != nil {
			info["provider_type"] = flattenEnvironmentProviderType(*environmentInfo.ProviderType)
		}
		if environmentInfo.ProvisioningType != nil {
			info["provisioning_type"] = flattenEnvironmentProvisioningType(*environmentInfo.ProvisioningType)
		}
		return []map[string]interface{}{info}
	}
	return nil
}

func flattenEnvironmentType(environmentType config.EnvironmentType) string {
	const ONPREM, NtnxCloud = 2, 3

	switch environmentType {
	case ONPREM:
		return "ONPREM"
	case NtnxCloud:
		return "NTNX_CLOUD"
	default:
		return "UNKNOWN"
	}
}

func flattenEnvironmentProvisioningType(provisioningType config.ProvisioningType) string {
	const NTNX, NATIVE = 2, 3
	switch provisioningType {
	case NTNX:
		return "NTNX"
	case NATIVE:
		return "NATIVE"
	default:
		return "UNKNOWN"
	}
}

func flattenEnvironmentProviderType(providerType config.ProviderType) string {
	const NTNX, AZURE, AWS, GCP, VSPHERE = 2, 3, 4, 5, 6
	switch providerType {
	case NTNX:
		return "NTNX"
	case AZURE:
		return "AZURE"
	case AWS:
		return "AWS"
	case GCP:
		return "GCP"
	case VSPHERE:
		return "VSPHERE"
	default:
		return "UNKNOWN"
	}
}

// flatten PC network

// network flattens

func flattenPCNetwork(network *config.DomainManagerNetwork) []map[string]interface{} {
	if network == nil {
		return nil
	}
	networkMap := make(map[string]interface{})
	if network.ExternalAddress != nil {
		networkMap["external_address"] = flattenIPAddress(network.ExternalAddress)
	}
	if network.NameServers != nil {
		networkMap["name_servers"] = flattenIPAddressOrFQDN(network.NameServers)
	}
	if network.NtpServers != nil {
		networkMap["ntp_servers"] = flattenIPAddressOrFQDN(network.NtpServers)
	}
	if network.Fqdn != nil {
		networkMap["fqdn"] = utils.StringValue(network.Fqdn)
	}
	//if network.InternalNetworks != nil {
	//	networkMap["internal_networks"] = flattenInternalNetworks(network.InternalNetworks)
	//}
	if network.ExternalNetworks != nil {
		networkMap["external_networks"] = flattenExternalNetworks(network.ExternalNetworks)
	}
	return []map[string]interface{}{networkMap}
}

func flattenIPAddress(ipAddress *commonConfig.IPAddress) []map[string]interface{} {
	if ipAddress == nil {
		return nil
	}
	ipAddressMap := make(map[string]interface{})
	if ipAddress.Ipv4 != nil {
		ipAddressMap["ipv4"] = flattenIPv4Address(ipAddress.Ipv4)
	}
	if ipAddress.Ipv6 != nil {
		ipAddressMap["ipv6"] = flattenIPv6Address(ipAddress.Ipv6)
	}
	return []map[string]interface{}{ipAddressMap}
}

func flattenIPv4Address(pr *commonConfig.IPv4Address) []interface{} {
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

func flattenIPv6Address(pr *commonConfig.IPv6Address) []interface{} {
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

func flattenIPAddressOrFQDN(pr []commonConfig.IPAddressOrFQDN) interface{} {
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

func flattenFQDN(pr *commonConfig.FQDN) []interface{} {
	if pr != nil {
		fqdn := make([]interface{}, 0)

		f := make(map[string]interface{})

		f["value"] = pr.Value

		fqdn = append(fqdn, f)

		return fqdn
	}
	return nil
}

func flattenIPRanges(ipRanges []commonConfig.IpRange) []map[string]interface{} {
	if len(ipRanges) > 0 {
		ipRangesMap := make([]map[string]interface{}, len(ipRanges))

		for k, v := range ipRanges {
			ipRangeMap := make(map[string]interface{})

			ipRangeMap["begin"] = flattenIPAddress(v.Begin)
			ipRangeMap["end"] = flattenIPAddress(v.End)

			ipRangesMap[k] = ipRangeMap
		}
		return ipRangesMap
	}
	return nil
}

func flattenExternalNetworks(externalNetworks []config.ExternalNetwork) []map[string]interface{} {
	if len(externalNetworks) > 0 {
		externalNetworksMap := make([]map[string]interface{}, len(externalNetworks))

		for k, v := range externalNetworks {
			externalNetworkMap := make(map[string]interface{})

			externalNetworkMap["default_gateway"] = flattenIPAddressOrFQDN([]commonConfig.IPAddressOrFQDN{*v.DefaultGateway})
			externalNetworkMap["subnet_mask"] = flattenIPAddressOrFQDN([]commonConfig.IPAddressOrFQDN{*v.SubnetMask})
			externalNetworkMap["ip_ranges"] = flattenIPRanges(v.IpRanges)
			externalNetworkMap["network_ext_id"] = v.NetworkExtId

			externalNetworksMap[k] = externalNetworkMap
		}
		return externalNetworksMap
	}
	return nil
}
