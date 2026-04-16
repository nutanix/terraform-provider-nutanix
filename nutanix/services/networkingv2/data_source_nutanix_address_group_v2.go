package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAddressGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixAddressGroupV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_addresses": SchemaForValuePrefixLength(),
			"ip_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"policy_references": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixAddressGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id")
	resp, err := conn.AddressGroupAPIInstance.GetAddressGroupById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching address group : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.AddressGroup)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ipv4_addresses", flattenIPv4AddressMicroSeg(getResp.Ipv4Addresses)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_ranges", flattenIPv4Range(getResp.IpRanges)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policy_references", flattenListofString(getResp.PolicyReferences)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksMicroSeg(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenLinksMicroSeg(pr []import2.ApiLink) []map[string]interface{} {
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

func flattenListofString(str []string) []string {
	if len(str) > 0 {
		strList := make([]string, len(str))

		strList = append(strList, str...)
		return strList
	}
	return nil
}

func flattenIPv4Range(pr []import1.IPv4Range) []interface{} {
	if len(pr) > 0 {
		ranges := make([]interface{}, len(pr))

		for k, v := range pr {
			rg := make(map[string]interface{})

			rg["start_ip"] = v.StartIp
			rg["end_ip"] = v.EndIp

			ranges[k] = rg
		}
		return ranges
	}
	return nil
}

func flattenIPv4AddressMicroSeg(pr []config.IPv4Address) []interface{} {
	if len(pr) > 0 {
		ipv4List := make([]interface{}, len(pr))

		for k, v := range pr {
			ipv4 := make(map[string]interface{})

			ipv4["value"] = v.Value
			ipv4["prefix_length"] = v.PrefixLength

			ipv4List[k] = ipv4
		}
		return ipv4List
	}
	return nil
}
