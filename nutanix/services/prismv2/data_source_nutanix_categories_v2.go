package prismv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixCategoriesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixCategoriesV2Read,
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
			"categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_group": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"detailed_associations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_group": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rel": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"href": {
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

func DatasourceNutanixCategoriesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	// initialize query params
	var filter, orderBy, expand, selects *string
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
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}
	resp, err := conn.CategoriesAPIInstance.ListCategories(page, limit, filter, orderBy, expand, selects)
	if err != nil {
		return diag.Errorf("error while fetching categories : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("categories", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of categories.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.Category)

	if err := d.Set("categories", flattenCategoriesEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenCategoriesEntities(pr []import1.Category) []interface{} {
	if len(pr) > 0 {
		ctgList := make([]interface{}, len(pr))

		for k, v := range pr {
			ctg := make(map[string]interface{})

			ctg["ext_id"] = v.ExtId
			ctg["key"] = v.Key
			ctg["value"] = v.Value
			ctg["type"] = flattenCategoryType(v.Type)
			ctg["description"] = v.Description
			ctg["owner_uuid"] = v.OwnerUuid
			ctg["associations"] = flattenAssociationSummary(v.Associations)
			ctg["detailed_associations"] = flattenAssociationDetail(v.DetailedAssociations)
			ctg["tenant_id"] = v.TenantId
			ctg["links"] = flattenLinks(v.Links)

			ctgList[k] = ctg
		}
		return ctgList
	}
	return nil
}
