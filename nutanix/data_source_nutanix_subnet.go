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
	conn := meta.(*NutanixClient).API

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
	if err := d.Set("categories", resp.Metadata.Categories); err != nil {
		return err
	}

	or := make(map[string]interface{})
	or["kind"] = utils.StringValue(resp.Metadata.OwnerReference.Kind)
	or["name"] = utils.StringValue(resp.Metadata.OwnerReference.Name)
	or["uuid"] = utils.StringValue(resp.Metadata.OwnerReference.UUID)
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
	clusterReference["kind"] = utils.StringValue(resp.Status.ClusterReference.Kind)
	clusterReference["name"] = utils.StringValue(resp.Status.ClusterReference.Name)
	clusterReference["uuid"] = utils.StringValue(resp.Status.ClusterReference.UUID)
	if err := d.Set("cluster_reference", clusterReference); err != nil {
		return err
	}
	// set message list values
	if resp.Status.MessageList != nil {
		messages := make([]map[string]interface{}, len(resp.Status.MessageList))
		for k, v := range resp.Status.MessageList {
			message := make(map[string]interface{})
			message["message"] = utils.StringValue(v.Message)
			message["reason"] = utils.StringValue(v.Reason)
			message["details"] = v.Details
			messages[k] = message
		}
		if err := d.Set("message_list", messages); err != nil {
			return err
		}
	}
	// set state value
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}
	if err := d.Set("vswitch_name", resp.Status.Resources.VswitchName); err != nil {
		return err
	}
	if err := d.Set("subnet_type", resp.Status.Resources.SubnetType); err != nil {
		return err
	}
	if err := d.Set("default_gateway_ip", resp.Status.Resources.IPConfig.DefaultGatewayIP); err != nil {
		return err
	}
	if err := d.Set("prefix_length", resp.Status.Resources.IPConfig.PrefixLength); err != nil {
		return err
	}
	if err := d.Set("subnet_ip", resp.Status.Resources.IPConfig.SubnetIP); err != nil {
		return err
	}
	if resp.Status.Resources.IPConfig.DHCPServerAddress != nil {
		//set ip_config.dhcp_server_address
		address := make(map[string]interface{})
		address["ip"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.IP)
		address["fqdn"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.FQDN)
		address["ipv6"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.IPV6)
		if err := d.Set("dhcp_server_address", address); err != nil {
			return err
		}
		if err := d.Set("dhcp_server_address_port", utils.Int64Value(resp.Status.Resources.IPConfig.DHCPServerAddress.Port)); err != nil {
			return err
		}
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
	}

	return nil
}

func getDataSourceSubnetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"metadata": &schema.Schema{
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"creation_time": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"spec_version": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"spec_hash": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"categories": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"owner_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"project_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"availability_zone_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"message_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"reason": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"details": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"vswitch_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet_type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"default_gateway_ip": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"prefix_length": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"subnet_ip": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"dhcp_server_address": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ip": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"fqdn": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ipv6": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"dhcp_server_address_port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"ip_config_pool_list_ranges": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dhcp_options": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"boot_file_name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"domain_name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"tftp_server_name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"dhcp_domain_name_server_list": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dhcp_domain_search_list": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"vlan_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"network_function_chain_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}
}
