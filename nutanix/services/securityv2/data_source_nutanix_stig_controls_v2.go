// Package securityv2 provides resources for managing security-related configurations in Nutanix.
package securityv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/security/v4/report"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixStigsControlsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixStigsControlsV2Read,
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
			"stigs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": common.LinksSchema(),
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stig_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identifiers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"affected_clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fix_text": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"benchmark_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixStigsControlsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	// initialize query params
	var filter, orderBy, selectQ *string
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
	if selectVal, ok := d.GetOk("select"); ok {
		selectQ = utils.StringPtr(selectVal.(string))
	} else {
		selectQ = nil
	}

	resp, err := conn.STIGsAPI.ListStigs(page, limit, filter, orderBy, selectQ)
	if err != nil {
		return diag.Errorf("error while fetching STIGs : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("stigs", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resource.UniqueId())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of STIGs.",
		}}
	}

	value := resp.Data.GetValue()
	stigsData, ok := value.([]report.Stig)
	if !ok {
		return diag.Errorf("unexpected response type: expected []report.Stig, got %T", value)
	}

	aJSON, _ := json.MarshalIndent(flattenSTIGs(stigsData), "", "  ")
	log.Printf("[DEBUG] STIGs fetched: %s", string(aJSON))

	// set the fetched STIGs to the resource data
	if err := d.Set("stigs", flattenSTIGs(stigsData)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenSTIGs(stigsData []report.Stig) []map[string]interface{} {
	if len(stigsData) == 0 {
		return make([]map[string]interface{}, 0)
	}

	stigsList := make([]map[string]interface{}, len(stigsData))
	for i, stig := range stigsData {
		severity := ""
		status := ""

		if stig.Severity != nil {
			severity = stig.Severity.GetName()
		}
		if stig.Status != nil {
			status = stig.Status.GetName()
		}
		stigsList[i] = map[string]interface{}{
			"ext_id":            stig.ExtId,
			"tenant_id":         stig.TenantId,
			"links":             flattenLinks(stig.Links),
			"title":             utils.StringValue(stig.Title),
			"rule_id":           utils.StringValue(stig.RuleId),
			"stig_version":      utils.StringValue(stig.StigVersion),
			"identifiers":       stig.Identifiers,
			"affected_clusters": stig.AffectedClusters,
			"severity":          severity,
			"status":            status,
			"comments":          utils.StringValue(stig.Comments),
			"fix_text":          utils.StringValue(stig.FixText),
			"benchmark_id":      utils.StringValue(stig.BenchmarkId),
		}
	}
	return stigsList
}

func flattenLinks(links []response.ApiLink) []interface{} {
	if len(links) > 0 {
		flattenedLinks := make([]interface{}, len(links))

		for k, v := range links {
			link := make(map[string]interface{})

			if v.Href != nil {
				link["href"] = v.Href
			}
			if v.Rel != nil {
				link["rel"] = v.Rel
			}
			flattenedLinks[k] = link
		}
		return flattenedLinks
	}
	return make([]interface{}, 0)
}
