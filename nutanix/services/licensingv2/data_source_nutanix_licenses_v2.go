// Package licensingv2 provides a client for interacting with the Nutanix licensing API.
package licensingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixListLicensesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixListLicensesV2Read,
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
			"licenses": {
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
							Type:     schema.TypeFloat,
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
										Type:     schema.TypeFloat,
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

func DatasourceNutanixListLicensesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	// initialize query params
	var filter, orderBy, expand, selectQ *string
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
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}
	if selectQy, ok := d.GetOk("select"); ok {
		selectQ = utils.StringPtr(selectQy.(string))
	} else {
		selectQ = nil
	}

	resp, err := conn.LicensesAPIInstance.ListLicenses(page, limit, filter, orderBy, expand, selectQ)
	if err != nil {
		return diag.Errorf("error while fetching licenses : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("licenses", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of licenses.",
		}}
	}

	licensesData := resp.Data.GetValue().([]config.LicenseProjection)

	aJSON, _ := json.MarshalIndent(flattenLicenses(licensesData), "", "  ")
	log.Printf("[DEBUG] Licenses fetched: %s", string(aJSON))

	// set the fetched Licenses to the resource data
	if err := d.Set("licenses", flattenLicenses(licensesData)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenLicenses(licensesData []config.LicenseProjection) []map[string]interface{} {
	if len(licensesData) == 0 {
		return make([]map[string]interface{}, 0)
	}

	licensesList := make([]map[string]interface{}, len(licensesData))
	for i, license := range licensesData {
		licenseMap := map[string]interface{}{
			"ext_id":                utils.StringValue(license.ExtId),
			"tenant_id":             utils.StringValue(license.TenantId),
			"links":                 flattenLinks(license.Links),
			"name":                  utils.StringValue(license.Name),
			"type":                  common.SafeGetName(license.Type),
			"category":              common.SafeGetName(license.Category),
			"sub_category":          common.SafeGetName(license.SubCategory),
			"expiry_date":           utils.TimeValue(license.ExpiryDate),
			"salesforce_license_id": utils.StringValue(license.SalesforceLicenseId),
			"meter":                 common.SafeGetName(license.Meter),
			"scope":                 common.SafeGetName(license.Scope),
			"quantity":              utils.Float64Value(license.Quantity),
			"consumption_details":   flattenConsumptionDetails(license.ConsumptionDetails),
		}
		licensesList[i] = licenseMap
	}
	return licensesList
}

func flattenConsumptionDetails(consumption []config.Consumption) []map[string]interface{} {
	if len(consumption) == 0 {
		return make([]map[string]interface{}, 0)
	}

	consumptionList := make([]map[string]interface{}, len(consumption))
	for i, c := range consumption {
		consumptionMap := map[string]interface{}{
			"cluster_ext_id": utils.StringValue(c.ClusterExtId),
			"quantity_used":  utils.Float64Value(c.QuantityUsed),
		}
		consumptionList[i] = consumptionMap
	}
	return consumptionList
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
	return nil
}
