package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixServiceGroupsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixServiceGroupsV2Read,
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
			"service_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
		},
	}
}

func DatasourceNutanixServiceGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

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

	resp, err := conn.ServiceGroupAPIInstance.ListServiceGroups(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching service groups : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("service_groups", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of service groups.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.ServiceGroup)
	if err := d.Set("service_groups", flattenServiceGroupsEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenServiceGroupsEntities(pr []import1.ServiceGroup) []interface{} {
	if len(pr) > 0 {
		serviceGroups := make([]interface{}, len(pr))

		for k, v := range pr {
			sg := make(map[string]interface{})

			sg["ext_id"] = v.ExtId
			sg["name"] = v.Name
			sg["description"] = v.Description
			sg["tcp_services"] = flattenTCPPortRangeSpec(v.TcpServices)
			sg["udp_services"] = flattenUDPPortRangeSpec(v.UdpServices)
			sg["icmp_services"] = flattenIcmpTypeCodeSpec(v.IcmpServices)
			sg["is_system_defined"] = v.IsSystemDefined
			if v.PolicyReferences != nil {
				sg["policy_references"] = flattenListofString(v.PolicyReferences)
			}
			if v.Links != nil {
				sg["links"] = flattenLinksMicroSeg(v.Links)
			}
			if v.CreatedBy != nil {
				sg["created_by"] = v.CreatedBy
			}
			if v.TenantId != nil {
				sg["tenant_id"] = v.TenantId
			}

			serviceGroups[k] = sg
		}
		return serviceGroups
	}
	return nil
}
