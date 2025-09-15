package licensingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixListLicenseCompliancesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixListLicenseCompliancesReadV2,
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"services": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"license_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_compliant": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enforcement_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enforcement_actions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"violations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"start_date": {
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
				},
			},
		},
	}
}

func DatasourceNutanixListLicenseCompliancesReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	resp, err := conn.LicensesAPIInstance.ListCompliances(page, limit, filter, orderBy, expands, selects)
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
			Detail:   "The API returned an empty list of License Compliances.",
		}}
	}

	getRes := resp.Data.GetValue().([]import1.Compliance)

	if err := d.Set("entities", flattenCompliances(getRes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenCompliances(compliances []import1.Compliance) []map[string]interface{} {
	if compliances == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(compliances))
	for i, compliance := range compliances {
		m := make(map[string]interface{})
		m["tenant_id"] = compliance.TenantId
		m["ext_id"] = compliance.ExtId
		m["links"] = flattenLinks(compliance.Links)
		m["is_multi_cluster"] = compliance.IsMulticluster
		m["cluster_ext_id"] = compliance.ClusterExtId
		m["type"] = compliance.Type
		m["services"] = flattenComplianceServices(compliance.Services)
		out[i] = m
	}
	return out
}

func flattenComplianceServices(services []import1.Service) []map[string]interface{} {
	if services == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(services))
	for i, service := range services {
		m := make(map[string]interface{})
		m["name"] = service.Name
		m["license_type"] = service.LicenseType
		m["is_compliant"] = service.IsCompliant
		m["enforcement_level"] = service.EnforcementLevel
		m["enforcement_actions"] = service.EnforcementActions
		m["violations"] = flattenServiceViolations(service.Violations)
		out[i] = m
	}
	return out
}

func flattenServiceViolations(violations []import1.ServiceViolation) []map[string]interface{} {
	if violations == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(violations))
	for i, violation := range violations {
		m := make(map[string]interface{})
		m["type"] = violation.Type
		m["start_date"] = violation.StartDate
		out[i] = m
	}
	return out
}