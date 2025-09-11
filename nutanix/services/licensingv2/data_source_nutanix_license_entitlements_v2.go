package licensingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixListLicenseEntitlementsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixListLicenseEntitlementsReadV2,
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
						"cluster_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_registered": {
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
									"meter": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"quantity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"earliest_expiry_date": {
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

func DatasourceNutanixListLicenseEntitlementsReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	resp, err := conn.LicensesAPIInstance.ListEntitlements(page, limit, filter, orderBy, expands, selects)
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
			Detail:   "The API returned an empty list of License Entitlements.",
		}}
	}

	getRes := resp.Data.GetValue().([]import1.Entitlement)

	if err := d.Set("entities", flattenEntitlements(getRes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenEntitlements(entitlements []import1.Entitlement) []map[string]interface{} {
	if entitlements == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(entitlements))
	for i, entitlement := range entitlements {
		m := make(map[string]interface{})
		m["tenant_id"] = entitlement.TenantId
		m["ext_id"] = entitlement.ExtId
		m["links"] = flattenLinks(entitlement.Links)
		m["is_multi_cluster"] = entitlement.IsMulticluster
		m["cluster_ext_id"] = entitlement.ClusterExtId
		m["is_registered"] = entitlement.IsRegistered
		m["name"] = entitlement.Name
		m["type"] = entitlement.Type
		m["details"] = flattenEntitlementDetails(entitlement.Details)
		out[i] = m
	}
	return out
}

func flattenEntitlementDetails(details []import1.EntitlementDetail) map[string]interface{} {
	if details == nil {
		return map[string]interface{}{}
	}
	m := make(map[string]interface{})
	for _, detail := range details {
		m["name"] = detail.Name
		m["type"] = detail.Type
		m["category"] = detail.Category
		m["sub_category"] = detail.SubCategory
		m["meter"] = detail.Meter
		m["quantity"] = detail.Quantity
		if detail.EarliestExpiryDate != nil {
			m["earliest_expiry_date"] = detail.EarliestExpiryDate.String()
		}
		m["scope"] = detail.Scope
	}
	return m
}