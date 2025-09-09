package licensingv2

import (
    "context"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/terraform-providers/terraform-provider-nutanix/utils"
    conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
    import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/config"
)

func DatasourceNutanixLicenseFeaturesV2() *schema.Resource {
    return &schema.Resource{
        ReadContext: DatasourceNutanixLicenseFeaturesReadV2,
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
                        "value_type": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "value": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "license_type": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "license_category": {
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
    }
}

func DatasourceNutanixLicenseFeaturesReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    conn := meta.(*conns.Client).LicensingAPI

    // initialize query params
    var filter, orderBy, selects *string
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
    if selectf, ok := d.GetOk("select"); ok {
        selects = utils.StringPtr(selectf.(string))
    }

    resp, err := conn.LicensesAPIInstance.ListFeatures(page, limit, filter, orderBy, selects)
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
				Detail:   "The API returned an empty list of Licensing Features.",
			}}
		}

    getRes := resp.Data.GetValue().([]import1.Feature)

    if err := d.Set("entities", flattenLicenseFeatures(getRes)); err != nil {
        return diag.FromErr(err)
    }

    d.SetId(utils.GenUUID())
    return nil
}


func flattenLicenseFeatures(features []import1.Feature) []map[string]interface{} {
		if features == nil {
				return []map[string]interface{}{}
		}
		out := make([]map[string]interface{}, len(features))
		for i, feature := range features {
				m := make(map[string]interface{})
				m["tenant_id"] = feature.TenantId
				m["ext_id"] = feature.ExtId
				m["links"] = flattenLinks(feature.Links)
				m["name"] = feature.Name
				m["value_type"] = feature.ValueType
				m["value"] = feature.Value
				m["license_type"] = feature.LicenseType
				m["license_category"] = feature.LicenseCategory
				m["scope"] = feature.Scope
				out[i] = m
		}
		return out
}