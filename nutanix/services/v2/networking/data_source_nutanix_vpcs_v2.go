package networking

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixVPCsv4() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixVPCsv4Read,
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
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DataSourceVPCSchemaV4(),
				},
			},
		},
	}
}

func dataSourceNutanixVPCsv4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	// initialize query params
	var filter, orderBy, selects *string
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
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}
	resp, err := conn.VpcAPIInstance.ListVpcs(page, limit, filter, orderBy, selects)
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
	getResp := resp.Data.GetValue().([]import1.Vpc)

	if err := d.Set("vpcs", flattenVPCsEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVPCsEntities(pr []import1.Vpc) []map[string]interface{} {
	if len(pr) > 0 {
		vpcs := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			vpc := make(map[string]interface{})

			if v.TenantId != nil {
				vpc["tenant_id"] = v.TenantId
			}
			vpc["ext_id"] = v.ExtId
			vpc["links"] = flattenLinks(v.Links)
			vpc["metadata"] = flattenMetadata(v.Metadata)
			vpc["name"] = v.Name
			vpc["description"] = v.Description
			vpc["common_dhcp_options"] = flattenCommonDhcpOptions(v.CommonDhcpOptions)
			vpc["snat_ips"] = flattenNtpServer(v.SnatIps)
			vpc["external_subnets"] = flattenExternalSubnets(v.ExternalSubnets)
			vpc["external_routing_domain_reference"] = v.ExternalRoutingDomainReference
			vpc["externally_routable_prefixes"] = flattenExternallyRoutablePrefixes(v.ExternallyRoutablePrefixes)

			vpcs[k] = vpc
		}
		return vpcs
	}
	return nil
}
