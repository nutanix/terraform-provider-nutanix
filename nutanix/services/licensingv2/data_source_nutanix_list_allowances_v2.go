package licensingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixListAllowancesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixListAllowancesReadV2,
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
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
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
						"is_multi_cluster": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"details": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"feature_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"scope": {
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

func DatasourceNutanixListAllowancesReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	// initialize query params
	var filter, orderBy, expands, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	}
	if expand, ok := d.GetOk("expand"); ok {
		expands = utils.StringPtr(expand.(string))
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	}

	resp, err := conn.LicensesAPIInstance.ListAllowances(page, limit, filter, orderBy, expands, selects)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Data == nil {
		if err := d.Set("entities", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of Allowances.",
		}}
	}

	getRes := resp.Data.GetValue().([]import1.Allowance)

	if err := d.Set("entities", flattenAllowances(getRes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenAllowances(allowances []import1.Allowance) []map[string]interface{} {
	if allowances == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(allowances))
	for i, allowance := range allowances {
		m := make(map[string]interface{})
		m["tenant_id"] = allowance.TenantId
		m["ext_id"] = allowance.ExtId
		m["links"] = flattenLinks(allowance.Links)
		m["is_multi_cluster"] = allowance.IsMulticluster
		m["name"] = allowance.Name
		m["type"] = allowance.Type
		m["details"] = flattenAllowanceDetails(allowance.Details)
		out[i] = m
	}
	return out
}

func flattenAllowanceDetails(details []import1.AllowanceDetail) []map[string]interface{} {
	if details == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(details))
	for i, detail := range details {
		m := make(map[string]interface{})
		m["feature_id"] = detail.FeatureId
		m["value_type"] = detail.ValueType
		m["value"] = detail.Value
		m["scope"] = detail.Scope
		out[i] = m
	}
	return out
}