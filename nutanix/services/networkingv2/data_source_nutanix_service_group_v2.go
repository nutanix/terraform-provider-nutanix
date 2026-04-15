package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixServiceGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixServiceGroupV2Read,
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
			"is_system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tcp_services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"udp_services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"icmp_services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_all_allowed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeInt,
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

func DatasourceNutanixServiceGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id")

	resp, err := conn.ServiceGroupAPIInstance.GetServiceGroupById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching service group : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.ServiceGroup)

	if err := d.Set("name", utils.StringValue(getResp.Name)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", utils.StringValue(getResp.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_system_defined", utils.BoolValue(getResp.IsSystemDefined)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tcp_services", flattenTCPPortRangeSpec(getResp.TcpServices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("udp_services", flattenUDPPortRangeSpec(getResp.UdpServices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("icmp_services", flattenIcmpTypeCodeSpec(getResp.IcmpServices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policy_references", flattenListofString(getResp.PolicyReferences)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_by", utils.StringValue(getResp.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksMicroSeg(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(getResp.TenantId)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenUDPPortRangeSpec(pr []import1.UdpPortRangeSpec) []interface{} {
	if len(pr) > 0 {
		ranges := make([]interface{}, len(pr))

		for k, v := range pr {
			rg := make(map[string]interface{})

			rg["start_port"] = utils.IntValue(v.StartPort)
			rg["end_port"] = utils.IntValue(v.EndPort)

			ranges[k] = rg
		}
		return ranges
	}
	return nil
}

func flattenTCPPortRangeSpec(pr []import1.TcpPortRangeSpec) []interface{} {
	if len(pr) > 0 {
		ranges := make([]interface{}, len(pr))

		for k, v := range pr {
			rg := make(map[string]interface{})

			rg["start_port"] = utils.IntValue(v.StartPort)
			rg["end_port"] = utils.IntValue(v.EndPort)

			ranges[k] = rg
		}
		return ranges
	}
	return nil
}

func flattenIcmpTypeCodeSpec(pr []import1.IcmpTypeCodeSpec) []interface{} {
	if len(pr) > 0 {
		ranges := make([]interface{}, len(pr))

		for k, v := range pr {
			rg := make(map[string]interface{})

			rg["is_all_allowed"] = utils.BoolValue(v.IsAllAllowed)
			rg["type"] = utils.IntValue(v.Type)
			rg["code"] = utils.IntValue(v.Code)

			ranges[k] = rg
		}
		return ranges
	}
	return nil
}
