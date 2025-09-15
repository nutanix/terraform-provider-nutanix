package licensingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixClusterLicenseRecommendationsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixClusterLicenseRecommendationsReadV2,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
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
						"details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scope": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operation": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"license_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"license_expiry_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"comment": {
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

func DatasourceNutanixClusterLicenseRecommendationsReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	// initialize query params
	var selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	}

	resp, err := conn.LicensesAPIInstance.ListRecommendations(page, limit, selects)
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
			Detail:   "The API returned an empty list of License Recommendations.",
		}}
	}

	getRes := resp.Data.GetValue().([]import1.Recommendation)

	if err := d.Set("entities", flattenRecommendations(getRes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenRecommendations(recommendations []import1.Recommendation) []map[string]interface{} {
	if recommendations == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(recommendations))
	for i, recommendation := range recommendations {
		m := make(map[string]interface{})
		m["tenant_id"] = recommendation.TenantId
		m["ext_id"] = recommendation.ExtId
		m["links"] = flattenLinks(recommendation.Links)
		m["details"] = flattenRecommendationDetails(recommendation.Details)
		out[i] = m
	}
	return out
}

func flattenRecommendationDetails(details []import1.RecommendationDetail) []map[string]interface{} {
	if details == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(details))
	for i, detail := range details {
		m := make(map[string]interface{})
		m["scope"] = detail.Scope
		m["operation"] = detail.Operation
		m["license_id"] = detail.LicenseId
		m["license_expiry_date"] = detail.LicenseExpiryDate
		m["comment"] = detail.Comment
		out[i] = m
	}
	return out
}