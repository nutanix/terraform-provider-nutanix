package nutanix

import (
	"github.com/hashicorp/terraform/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixSubnetsRead,
		Schema: map[string]*schema.Schema{
			"metadata": getDSMetadataSchema(),
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
						"categories": categoriesSchema(),
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
								},
							},
						},
						"cluster_reference_name": {
							Type:     schema.TypeString,
							Computed: true,
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
		},
	}
}

func dataSourceNutanixSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.ListAllSubnet()
	if err != nil {
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
		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)
		entity["availability_zone_reference"] = flattenReferenceValues(v.Status.AvailabilityZoneReference)
		entity["cluster_reference"] = getClusterReferenceValues(v.Status.ClusterReference)
		entity["cluster_reference_name"] = utils.StringValue(v.Status.ClusterReference.Name)
		entity["state"] = utils.StringValue(v.Status.State)
		entity["vswitch_name"] = utils.StringValue(v.Status.Resources.VswitchName)

		stype := ""
		if v.Status.Resources.SubnetType != nil {
			stype = utils.StringValue(v.Status.Resources.SubnetType)
		}
		entity["subnet_type"] = stype

		dgIP := ""
		sIP := ""
		pl := int64(0)
		port := int64(0)
		address := make(map[string]interface{})
		dhcpSA := make(map[string]interface{})
		dOptions := make(map[string]interface{})
		ipcpl := make([]string, 0)
		dnsList := make([]string, 0)
		dsList := make([]string, 0)

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
				poolList := make([]string, len(pl))
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
					dnsList = utils.StringValueSlice(v.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList)
				}
				if v.Status.Resources.IPConfig.DHCPOptions.DomainSearchList != nil {
					dsList = utils.StringValueSlice(v.Status.Resources.IPConfig.DHCPOptions.DomainSearchList)
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
		entity["network_function_chain_reference"] = flattenReferenceValues(v.Status.Resources.NetworkFunctionChainReference)
		entities[k] = entity
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.SetId(uuid.NewV4().String())

	return d.Set("entities", entities)
}
