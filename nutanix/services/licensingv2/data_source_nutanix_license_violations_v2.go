package licensingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLicenseViolationsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLicenseViolationsReadV2,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
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
						"feature_violations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"feature_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"affected_entity": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"capacity_violations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"meter": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"shortfall": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"expired_licenses": {
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
									"license_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"meter": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"used_quantity": {
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

func DatasourceNutanixLicenseViolationsReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	// initialize query params
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	}

	resp, err := conn.LicensesAPIInstance.ListViolations(page, limit)
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
			Detail:   "The API returned an empty list of License Violations.",
		}}
	}

	getRes := resp.Data.GetValue().([]import1.Violation)

	if err := d.Set("entities", flattenLicenseViolations(getRes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenLicenseViolations(violations []import1.Violation) []map[string]interface{} {
	if violations == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(violations))
	for i, violation := range violations {
		m := make(map[string]interface{})
		m["tenant_id"] = violation.TenantId
		m["ext_id"] = violation.ExtId
		m["links"] = flattenLinks(violation.Links)
		m["is_multi_cluster"] = violation.IsMulticluster
		m["feature_violations"] = flattenFeatureViolations(violation.FeatureViolations)
		m["capacity_violations"] = flattenCapacityViolations(violation.CapacityViolations)
		m["expired_licenses"] = flattenExpiredLicenses(violation.ExpiredLicenses)
		out[i] = m
	}
	return out
}

func flattenFeatureViolations(featureViolations []import1.FeatureViolation) []map[string]interface{} {
	if featureViolations == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(featureViolations))
	for i, featureViolation := range featureViolations {
		m := make(map[string]interface{})
		m["feature_id"] = featureViolation.FeatureId
		m["name"] = featureViolation.Name
		m["description"] = featureViolation.Description
		m["affected_entity"] = featureViolation.AffectedEntity
		out[i] = m
	}
	return out
}

func flattenCapacityViolations(capacityViolations []import1.CapacityViolation) []map[string]interface{} {
	if capacityViolations == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(capacityViolations))
	for i, capacityViolation := range capacityViolations {
		m := make(map[string]interface{})
		m["type"] = capacityViolation.Type
		m["category"] = capacityViolation.Category
		m["meter"] = capacityViolation.Meter
		m["shortfall"] = capacityViolation.Shortfall
		out[i] = m
	}
	return out
}

func flattenExpiredLicenses(licenses []import1.ExpiredLicense) []map[string]interface{} {
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
		m["license_id"] = license.LicenseId
		m["meter"] = license.Meter
		m["used_quantity"] = license.UsedQuantity
		out[i] = m
	}
	return out
}
