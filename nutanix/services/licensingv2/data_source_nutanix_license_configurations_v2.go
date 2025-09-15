package licensingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLicenseConfigurationV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLicenseConfigurationReadV2,
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
						"logical_version": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"license_verison": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"is_stand_by": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"has_non_compliant_features": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_license_check_disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"license_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enforcement_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"license_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_ultimate_trail_ended": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"post_paid_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_pulse_required": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"billing_plan": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"consumption_type": {
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

func DatasourceNutanixLicenseConfigurationReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	// initialize query params
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	}

	resp, err := conn.LicensesAPIInstance.ListSettings(page, limit)
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
			Detail:   "The API returned an empty list of License Configuration.",
		}}
	}

	getRes := resp.Data.GetValue().([]import1.Setting)

	if err := d.Set("entities", flattenLicenseConfiguration(getRes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenLicenseConfiguration(settings []import1.Setting) []map[string]interface{} {
	if settings == nil {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, len(settings))
	for i, setting := range settings {
		m := make(map[string]interface{})
		m["tenant_id"] = setting.TenantId
		m["ext_id"] = setting.ExtId
		m["links"] = flattenLinks(setting.Links)
		m["is_multi_cluster"] = setting.IsMulticluster
		m["logical_version"] = flattenLogicalVersion(setting.LogicalVersion)
		m["is_stand_by"] = setting.IsStandby
		m["has_non_compliant_features"] = setting.HasNonCompliantFeatures
		m["is_license_check_disabled"] = setting.IsLicenseCheckDisabled
		m["license_class"] = setting.LicenseClass
		m["enforcement_policy"] = setting.EnforcementPolicy
		m["license_key"] = setting.LicenseKey
		m["has_ultimate_trail_ended"] = setting.HasUltimateTrialEnded
		m["post_paid_config"] = flattenPostPaidConfig(setting.PostPaidConfig)
		out[i] = m
	}
	return out
}

func flattenPostPaidConfig(postPaidConfig *import1.PostPaidConfig) map[string]interface{} {
	if postPaidConfig == nil {
		return map[string]interface{}{}
	}
	m := make(map[string]interface{})
	m["id"] = postPaidConfig.Id
	m["category"] = postPaidConfig.Category
	m["is_pulse_required"] = postPaidConfig.IsPulseRequired
	m["billing_plan"] = postPaidConfig.BillingPlan
	m["consumption_type"] = postPaidConfig.ConsumptionType
	return m
}

func flattenLogicalVersion(logicalVersion *import1.LogicalVersion) map[string]interface{} {
	if logicalVersion == nil {
		return map[string]interface{}{}
	}
	m := make(map[string]interface{})
	m["cluster_version"] = logicalVersion.ClusterVersion
	m["license_version"] = logicalVersion.LicenseVersion
	return m
}
