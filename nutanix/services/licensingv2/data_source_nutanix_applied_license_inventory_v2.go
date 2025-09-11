package licensingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAppliedLicenseInventoryV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixAppliedLicenseInventoryReadV2,
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
			"expand": {
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"salesforce_license_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"meter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"consumption_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"quantity_used": {
										Type:     schema.TypeInt,
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

func DatasourceNutanixAppliedLicenseInventoryReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	// initialize query params
	var filter, orderBy, expandf, selects *string
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
		expandf = utils.StringPtr(expand.(string))
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	}

	resp, err := conn.LicensesAPIInstance.ListLicenses(page, limit, filter, orderBy, expandf, selects)
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
			Detail:   "The API returned an empty list of applied license inventory.",
		}}
	}

	getRes := resp.Data.GetValue().([]import1.License)

	if err := d.Set("entities", flattenAppliedLicenseInventory(getRes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenAppliedLicenseInventory(licenses []import1.License) []map[string]interface{} {
	if licenses == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(licenses))
	for i, license := range licenses {
		m := make(map[string]interface{})
		m["tenant_id"] = license.TenantId
		m["ext_id"] = license.ExtId
		m["links"] = flattenLinks(license.Links)
		m["name"] = license.Name
		m["type"] = license.Type
		m["category"] = license.Category
		m["sub_category"] = license.SubCategory
		m["expiry_date"] = license.ExpiryDate
		m["salesforce_license_id"] = license.SalesforceLicenseId
		m["meter"] = license.Meter
		m["scope"] = license.Scope
		m["quantity"] = license.Quantity
		m["consumption_details"] = flattenLicenseConsumptionDetails(license.ConsumptionDetails)
		out[i] = m
	}
	return out
}

func flattenLicenseConsumptionDetails(details []import1.Consumption) []map[string]interface{} {
	if details == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(details))
	for i, detail := range details {
		m := make(map[string]interface{})
		m["cluster_ext_id"] = detail.ClusterExtId
		m["quantity_used"] = detail.QuantityUsed
		out[i] = m
	}
	return out
}
