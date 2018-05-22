package nutanix

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixSubnet() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNutanixSubnetRead,
		Schema: getDataSourceSubnetSchema(),
	}
}

func dataSourceNutanixSubnetRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	subnetID, ok := d.GetOk("subnet_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute vm_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetSubnet(subnetID.(string))

	if err != nil {
		return err
	}

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = resp.Metadata.LastUpdateTime.String()
	metadata["kind"] = utils.StringValue(resp.Metadata.Kind)
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
	if err := d.Set("api_version", utils.StringValue(resp.APIVersion)); err != nil {
		return err
	}
	if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
		return err
	}
	if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
		return err
	}
	// set availability zone reference values
	availabilityZoneReference := make(map[string]interface{})
	if resp.Status.AvailabilityZoneReference != nil {
		availabilityZoneReference["kind"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Kind)
		availabilityZoneReference["name"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Name)
		availabilityZoneReference["uuid"] = utils.StringValue(resp.Status.AvailabilityZoneReference.UUID)
	}
	if err := d.Set("availability_zone_reference", availabilityZoneReference); err != nil {
		return err
	}
	// set cluster reference values
	clusterReference := make(map[string]interface{})
	if resp.Status.ClusterReference != nil {
		clusterReference["kind"] = utils.StringValue(resp.Status.ClusterReference.Kind)
		clusterReference["name"] = utils.StringValue(resp.Status.ClusterReference.Name)
		clusterReference["uuid"] = utils.StringValue(resp.Status.ClusterReference.UUID)
	}
	if err := d.Set("cluster_reference", clusterReference); err != nil {
		return err
	}
	// set state value
	if err := d.Set("state", utils.StringValue(resp.Status.State)); err != nil {
		return err
	}
	if err := d.Set("vswitch_name", utils.StringValue(resp.Status.Resources.VswitchName)); err != nil {
		return err
	}
	if resp.Status.Resources.SubnetType != nil {
		if err := d.Set("subnet_type", utils.StringValue(resp.Status.Resources.SubnetType)); err != nil {
			return err
		}
	} else {
		if err := d.Set("subnet_type", ""); err != nil {
			return err
		}
	}
	if resp.Status.Resources.IPConfig != nil {
		if err := d.Set("default_gateway_ip", utils.StringValue(resp.Status.Resources.IPConfig.DefaultGatewayIP)); err != nil {
			return err
		}
		if err := d.Set("prefix_length", utils.Int64Value(resp.Status.Resources.IPConfig.PrefixLength)); err != nil {
			return err
		}
		if err := d.Set("subnet_ip", utils.StringValue(resp.Status.Resources.IPConfig.SubnetIP)); err != nil {
			return err
		}
		address := make(map[string]interface{})
		port := int64(0)
		if resp.Status.Resources.IPConfig.DHCPServerAddress != nil {
			//set ip_config.dhcp_server_address
			address["ip"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.IP)
			address["fqdn"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.FQDN)
			address["ipv6"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.IPV6)
			port = utils.Int64Value(resp.Status.Resources.IPConfig.DHCPServerAddress.Port)
		}
		if err := d.Set("dhcp_server_address", address); err != nil {
			return err
		}
		if err := d.Set("dhcp_server_address_port", port); err != nil {
			return err
		}
		if resp.Status.Resources.IPConfig.PoolList != nil {
			pl := resp.Status.Resources.IPConfig.PoolList
			poolList := make([]string, len(pl))
			for k, v := range pl {
				poolList[k] = utils.StringValue(v.Range)
			}
			if err := d.Set("ip_config_pool_list_ranges", poolList); err != nil {
				return err
			}
		} else {
			if err := d.Set("ip_config_pool_list_ranges", make([]string, 0)); err != nil {
				return err
			}
		}
		if resp.Status.Resources.IPConfig.DHCPOptions != nil {
			//set dhcp_options
			dOptions := make(map[string]interface{})
			dOptions["boot_file_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.BootFileName)
			dOptions["domain_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.DomainName)
			dOptions["tftp_server_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.TFTPServerName)

			if err := d.Set("dhcp_options", dOptions); err != nil {
				return err
			}

			if resp.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList != nil {
				dnsl := resp.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList
				dnsList := make([]string, len(dnsl))
				for k, v := range dnsl {
					dnsList[k] = utils.StringValue(v)
				}
				if err := d.Set("dhcp_domain_name_server_list", dnsList); err != nil {
					return err
				}
			}
			if resp.Status.Resources.IPConfig.DHCPOptions.DomainSearchList != nil {
				dnsl := resp.Status.Resources.IPConfig.DHCPOptions.DomainSearchList
				dsList := make([]string, len(dnsl))
				for k, v := range dnsl {
					dsList[k] = utils.StringValue(v)
				}
				if err := d.Set("dhcp_domain_search_list", dsList); err != nil {
					return err
				}
			}
		} else {
			if err := d.Set("dhcp_options", make(map[string]interface{})); err != nil {
				return err
			}
			if err := d.Set("dhcp_domain_name_server_list", make([]map[string]interface{}, 0)); err != nil {
				return err
			}
			if err := d.Set("dhcp_domain_search_list", make([]map[string]interface{}, 0)); err != nil {
				return err
			}
		}
	} else {
		if err := d.Set("default_gateway_ip", ""); err != nil {
			return err
		}
		if err := d.Set("prefix_length", 0); err != nil {
			return err
		}
		if err := d.Set("subnet_ip", ""); err != nil {
			return err
		}
		if err := d.Set("dhcp_server_address_port", 0); err != nil {
			return err
		}
		if err := d.Set("ip_config_pool_list_ranges", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
		if err := d.Set("dhcp_options", make(map[string]interface{})); err != nil {
			return err
		}
		if err := d.Set("dhcp_domain_name_server_list", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
		if err := d.Set("dhcp_domain_search_list", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
	}

	if err := d.Set("vlan_id", resp.Status.Resources.VlanID); err != nil {
		return err
	}
	// set network_function_chain_reference
	if resp.Status.Resources.NetworkFunctionChainReference != nil {
		nfcr := make(map[string]interface{})
		nfcr["kind"] = utils.StringValue(resp.Status.Resources.NetworkFunctionChainReference.Kind)
		nfcr["name"] = utils.StringValue(resp.Status.Resources.NetworkFunctionChainReference.Name)
		nfcr["uuid"] = utils.StringValue(resp.Status.Resources.NetworkFunctionChainReference.UUID)

		if err := d.Set("network_function_chain_reference", nfcr); err != nil {
			return err
		}
	} else {
		if err := d.Set("network_function_chain_reference", make(map[string]interface{})); err != nil {
			return err
		}
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func getDataSourceSubnetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_id": {
			Type:     schema.TypeString,
			Required: true,
		},
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
	}
}
