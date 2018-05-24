package nutanix

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixSubnets() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNutanixSubnetsRead,
		Schema: getDataSourceSubnetsSchema(),
	}
}

func dataSourceNutanixSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	metadata := &v3.SubnetListMetadata{}

	if v, ok := d.GetOk("metadata"); ok {
		m := v.(map[string]interface{})
		metadata.Kind = utils.String("subnet")
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
	resp, err := conn.V3.ListSubnet(metadata)
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

		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)

		// set availability zone reference values
		availabilityZoneReference := make(map[string]interface{})
		if v.Status.AvailabilityZoneReference != nil {
			availabilityZoneReference["kind"] = utils.StringValue(v.Status.AvailabilityZoneReference.Kind)
			availabilityZoneReference["name"] = utils.StringValue(v.Status.AvailabilityZoneReference.Name)
			availabilityZoneReference["uuid"] = utils.StringValue(v.Status.AvailabilityZoneReference.UUID)
		}
		entity["availability_zone_reference"] = availabilityZoneReference
		// set cluster reference values
		clusterReference := make(map[string]interface{})
		if v.Status.ClusterReference != nil {
			clusterReference["kind"] = utils.StringValue(v.Status.ClusterReference.Kind)
			clusterReference["name"] = utils.StringValue(v.Status.ClusterReference.Name)
			clusterReference["uuid"] = utils.StringValue(v.Status.ClusterReference.UUID)
		}
		entity["cluster_reference"] = clusterReference
		entity["state"] = utils.StringValue(v.Status.State)

		entity["vswitch_name"] = utils.StringValue(v.Status.Resources.VswitchName)

		stype := ""
		if v.Status.Resources.SubnetType != nil {
			stype = utils.StringValue(v.Status.Resources.SubnetType)
		}
		entity["subnet_type"] = stype

		dgIP := ""
		pl := int64(0)
		sIP := ""
		address := make(map[string]interface{})
		port := int64(0)
		dhcpSA := make(map[string]interface{})
		ipcpl := make([]string, 0)
		dOptions := make(map[string]interface{})
		dnsList := make([]string, 0)
		dsList := make([]string, 0)
		poolList := make([]string, 0)

		if v.Status.Resources.IPConfig != nil {
			dgIP = utils.StringValue(v.Status.Resources.IPConfig.DefaultGatewayIP)
			pl = utils.Int64Value(v.Status.Resources.IPConfig.PrefixLength)
			sIP = utils.StringValue(v.Status.Resources.IPConfig.SubnetIP)

			if v.Status.Resources.IPConfig.DHCPServerAddress != nil {
				address["ip"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPServerAddress.IP)
				address["fqdn"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPServerAddress.FQDN)
				address["ipv6"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPServerAddress.IPV6)
				port = utils.Int64Value(v.Status.Resources.IPConfig.DHCPServerAddress.Port)
			}

			dhcpSA = address

			if v.Status.Resources.IPConfig.PoolList != nil {
				pl := v.Status.Resources.IPConfig.PoolList
				poolList = make([]string, len(pl))
				for k, v := range pl {
					poolList[k] = utils.StringValue(v.Range)
				}
				ipcpl = poolList
			}
			if v.Status.Resources.IPConfig.DHCPOptions != nil {
				dOptions["boot_file_name"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPOptions.BootFileName)
				dOptions["domain_name"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPOptions.DomainName)
				dOptions["tftp_server_name"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPOptions.TFTPServerName)

				if v.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList != nil {
					dnsl := v.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList
					dnsList = make([]string, len(dnsl))
					for k, v := range dnsl {
						dnsList[k] = utils.StringValue(v)
					}
				}
				if v.Status.Resources.IPConfig.DHCPOptions.DomainSearchList != nil {
					dnsl := v.Status.Resources.IPConfig.DHCPOptions.DomainSearchList
					dsList = make([]string, len(dnsl))
					for k, v := range dnsl {
						dsList[k] = utils.StringValue(v)
					}
				}
			}
		}
		entity["default_gateway_ip"] = dgIP
		entity["prefix_length"] = pl
		entity["subnet_ip"] = sIP
		entity["dhcp_server_address"] = dhcpSA
		entity["dhcp_server_address_port"] = port
		entity["ip_config_pool_list_ranges"] = ipcpl
		entity["dhcp_options"] = dOptions
		entity["dhcp_domain_name_server_list"] = dnsList
		entity["dhcp_domain_search_list"] = dsList

		entity["vlan_id"] = utils.Int64Value(v.Status.Resources.VlanID)

		nfcr := make(map[string]interface{})
		if v.Status.Resources.NetworkFunctionChainReference != nil {
			nfcr["kind"] = utils.StringValue(v.Status.Resources.NetworkFunctionChainReference.Kind)
			nfcr["name"] = utils.StringValue(v.Status.Resources.NetworkFunctionChainReference.Name)
			nfcr["uuid"] = utils.StringValue(v.Status.Resources.NetworkFunctionChainReference.UUID)
		}
		entity["network_function_chain_reference"] = nfcr
		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}

func getDataSourceSubnetsSchema() map[string]*schema.Schema {
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
					"api_version": {
						Type:     schema.TypeString,
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
								"kind": {
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
						Type:     schema.TypeList,
						Optional: true,
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
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
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
					"project_reference": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
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
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"availability_zone_reference": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
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
					"message_list": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"message": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"reason": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"details": {
									Type:     schema.TypeMap,
									Computed: true,
								},
							},
						},
					},
					"cluster_reference": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
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
					"vswitch_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"subnet_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"default_gateway_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"prefix_length": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"subnet_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"dhcp_server_address": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"fqdn": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"ipv6": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"dhcp_server_address_port": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"ip_config_pool_list_ranges": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"dhcp_options": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"boot_file_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"domain_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"tftp_server_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"dhcp_domain_name_server_list": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"dhcp_domain_search_list": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"vlan_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
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
				},
			},
		},
	}
}
