package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixSubnetsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixSubnetsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DataSourceNutanixSubnetV2(),
			},
		},
	}
}

func dataSourceNutanixSubnetsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	// initialize query params
	var filter, orderBy, expand, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	resp, err := conn.SubnetAPIInstance.ListSubnets(page, limit, filter, orderBy, expand, selects)
	if err != nil {
		return diag.Errorf("error while fetching subnets : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("subnets", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of subnets.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.Subnet)

	if err := d.Set("subnets", flattenSubnetEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenSubnetEntities(pr []import1.Subnet) []interface{} {
	if len(pr) > 0 {
		subnets := make([]interface{}, len(pr))

		for k, v := range pr {
			sub := make(map[string]interface{})

			sub["ext_id"] = v.ExtId
			sub["name"] = v.Name
			sub["description"] = v.Description
			sub["links"] = flattenLinks(v.Links)
			sub["subnet_type"] = flattenSubnetType(v.SubnetType)
			sub["network_id"] = v.NetworkId
			sub["dhcp_options"] = flattenDhcpOptions(v.DhcpOptions)
			sub["ip_config"] = flattenIPConfig(v.IpConfig)
			sub["cluster_reference"] = v.ClusterReference
			sub["virtual_switch_reference"] = v.VirtualSwitchReference
			sub["vpc_reference"] = v.VpcReference
			sub["is_nat_enabled"] = v.IsNatEnabled
			sub["is_external"] = v.IsExternal
			sub["reserved_ip_addresses"] = flattenReservedIPAddresses(v.ReservedIpAddresses)
			sub["dynamic_ip_addresses"] = flattenReservedIPAddresses(v.DynamicIpAddresses)
			sub["network_function_chain_reference"] = v.NetworkFunctionChainReference
			sub["bridge_name"] = v.BridgeName
			sub["is_advanced_networking"] = v.IsAdvancedNetworking
			sub["cluster_name"] = v.ClusterName
			sub["hypervisor_type"] = v.HypervisorType
			sub["virtual_switch"] = flattenVirtualSwitch(v.VirtualSwitch)
			sub["vpc"] = flattenVPC(v.Vpc)
			sub["ip_prefix"] = v.IpPrefix
			sub["ip_usage"] = flattenIPUsage(v.IpUsage)
			sub["migration_state"] = flattenMigrationState(v.MigrationState)
			sub["metadata"] = flattenMetadata(v.Metadata)
			subnets[k] = sub
		}
		return subnets
	}
	return nil
}
