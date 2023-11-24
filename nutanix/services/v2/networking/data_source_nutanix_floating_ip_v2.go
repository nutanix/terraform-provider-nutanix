package networking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/models/networking/v4/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixFloatingIPV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixFloatingIPV4Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"association": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_nic_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
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
			"floating_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(),
						"ipv6": SchemaForValuePrefixLength(),
					},
				},
			},
			"external_subnet_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_subnet": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DataSourceNutanixSubnetv4(),
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"association_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_nic_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DataSourceVPCSchemaV4(),
				},
			},
			"vm_nic": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV4(),
				},
			},
		},
	}
}

func DatasourceNutanixFloatingIPV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	extID := d.Get("ext_id")
	resp, err := conn.FloatingIpApiInstance.GetFloatingIp(utils.StringPtr(extID.(string)))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching subnets : %v", errorMessage["message"])
	}

	getResp := resp.Data.GetValue().(import1.FloatingIp)
	fmt.Println(getResp)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("association", flattenAssociation(getResp.Association)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("floating_ip", flattenFloatingIP(getResp.FloatingIp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_subnet_reference", getResp.ExternalSubnetReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_subnet", flattenExternalSubnet(getResp.ExternalSubnet)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("private_ip", getResp.PrivateIp); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("floating_ip_value", getResp.FloatingIpValue); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("association_status", getResp.AssociationStatus); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm_nic_reference", getResp.VmNicReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc", flattenVpc(getResp.Vpc)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_nic", flattenVmNic(getResp.VmNic)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extID.(string))
	return nil
}

func flattenFloatingIP(pr *import1.FloatingIPAddress) []map[string]interface{} {
	if pr != nil {
		fips := make([]map[string]interface{}, 0)

		fip := make(map[string]interface{})

		fip["ipv4"] = flattenFloatingIPv4Address(pr.Ipv4)
		fip["ipv6"] = flattenFloatingIPv6Address(pr.Ipv6)

		fips = append(fips, fip)
		return fips
	}
	return nil
}

func flattenFloatingIPv4Address(pr *import1.FloatingIPv4Address) []map[string]interface{} {
	if pr != nil {
		ips := make([]map[string]interface{}, 0)

		ip := make(map[string]interface{})

		ip["prefix_length"] = pr.PrefixLength
		ip["value"] = pr.Value

		ips = append(ips, ip)

		return ips
	}
	return nil
}

func flattenFloatingIPv6Address(pr *import1.FloatingIPv6Address) []map[string]interface{} {
	if pr != nil {
		ips := make([]map[string]interface{}, 0)

		ip := make(map[string]interface{})

		ip["prefix_length"] = pr.PrefixLength
		ip["value"] = pr.Value

		ips = append(ips, ip)

		return ips
	}
	return nil
}

func flattenVpc(pr *import1.Vpc) []map[string]interface{} {
	if pr != nil {
		vpcList := make([]map[string]interface{}, 0)

		vpc := make(map[string]interface{})

		vpc["tenant_id"] = pr.TenantId
		vpc["ext_id"] = pr.ExtId
		vpc["links"] = flattenLinks(pr.Links)
		vpc["metadata"] = flattenMetadata(pr.Metadata)
		vpc["name"] = pr.Name
		vpc["description"] = pr.Description
		vpc["common_dhcp_options"] = flattenCommonDhcpOptions(pr.CommonDhcpOptions)
		vpc["snat_ips"] = flattenNtpServer(pr.SnatIps)
		vpc["external_subnets"] = flattenExternalSubnets(pr.ExternalSubnets)
		vpc["external_routing_domain_reference"] = pr.ExternalRoutingDomainReference
		vpc["externally_routable_prefixes"] = flattenExternallyRoutablePrefixes(pr.ExternallyRoutablePrefixes)
		vpcList = append(vpcList, vpc)

		return vpcList
	}
	return nil
}

func flattenVmNic(pr *import1.VmNic) []map[string]interface{} {
	if pr != nil {
		nics := make([]map[string]interface{}, 0)
		nic := make(map[string]interface{})

		nic["private_ip"] = pr.PrivateIp

		nics = append(nics, nic)
		return nics
	}
	return nil
}

func flattenExternalSubnet(pr *import1.Subnet) []map[string]interface{} {
	if pr != nil {
		subs := make([]map[string]interface{}, 0)

		sub := make(map[string]interface{})

		sub["name"] = pr.Name
		sub["description"] = pr.Description
		sub["links"] = pr.Links
		sub["subnet_type"] = flattenSubnetType(pr.SubnetType)
		sub["network_id"] = pr.NetworkId
		sub["dhcp_options"] = flattenDhcpOptions(pr.DhcpOptions)
		sub["ip_config"] = flattenIPConfig(pr.IpConfig)
		sub["cluster_reference"] = pr.ClusterReference
		sub["virtual_switch_reference"] = pr.VirtualSwitchReference
		sub["vpc_reference"] = pr.VpcReference
		sub["is_nat_enabled"] = pr.IsNatEnabled
		sub["is_external"] = pr.IsExternal
		sub["reserved_ip_addresses"] = flattenReservedIPAddresses(pr.ReservedIpAddresses)
		sub["dynamic_ip_addresses"] = flattenReservedIPAddresses(pr.DynamicIpAddresses)
		sub["network_function_chain_reference"] = pr.NetworkFunctionChainReference
		sub["bridge_name"] = pr.BridgeName
		sub["is_advanced_networking"] = pr.IsAdvancedNetworking
		sub["cluster_name"] = pr.ClusterName
		sub["hypervisor_type"] = pr.HypervisorType
		sub["virtual_switch"] = flattenVirtualSwitch(pr.VirtualSwitch)
		sub["vpc"] = flattenVPC(pr.Vpc)
		sub["ip_prefix"] = pr.IpPrefix
		sub["ip_usage"] = pr.IpUsage
		sub["migration_state"] = pr.MigrationState

		subs = append(subs, sub)
		return subs
	}
	return nil
}

func flattenAssociation(pr *import1.OneOfFloatingIpAssociation) []map[string]interface{} {
	if pr != nil {
		resList := make([]map[string]interface{}, 0)

		assc := make(map[string]interface{})

		getVal := pr.ObjectType_

		if utils.StringValue(getVal) == "networking.v4.config.PrivateIpAssociation" {
			privateIPAssociation := pr.GetValue().(import1.PrivateIpAssociation)
			assc["private_ip"] = flattenNodeIPAddress(privateIPAssociation.PrivateIp)
			assc["vpc_reference"] = privateIPAssociation.VpcReference
		} else {
			vmNicAssociation := pr.GetValue().(import1.VmNicAssociation)
			assc["vm_nic_reference"] = vmNicAssociation.VmNicReference
			assc["vpc_reference"] = vmNicAssociation.VpcReference
		}

		resList = append(resList, assc)
		return resList
	}
	return nil
}
