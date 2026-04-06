package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamConfig "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEntitiesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixEntitiesV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Description: "A URL query parameter that specifies the page number of the result set.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"limit": {
				Description: "A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"filter": {
				Description: "OData filter expression for filtering entities.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"order_by": {
				Description: "OData orderby expression for sorting entities.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"select": {
				Description: "OData select expression to specify which fields to return.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"entities": {
				Description: "List of IAM entities.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        DatasourceNutanixEntityV2(),
			},
		},
	}
}

func DatasourceNutanixEntitiesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	var page, limit *int
	var filter, orderBy, selectParam *string

	if v, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		selectParam = utils.StringPtr(v.(string))
	}

	resp, err := conn.EntityAPIInstance.ListEntities(page, limit, filter, orderBy, selectParam)
	if err != nil {
		return diag.Errorf("error while listing entities: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("entities", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of entities.",
		}}
	}

	entities := resp.Data.GetValue().([]iamConfig.Entity)

	if err := d.Set("entities", flattenEntities(entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenEntities(entities []iamConfig.Entity) []map[string]interface{} {
	if len(entities) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(entities))
	for _, e := range entities {
		m := map[string]interface{}{
			"tenant_id":         utils.StringValue(e.TenantId),
			"ext_id":            utils.StringValue(e.ExtId),
			"links":             flattenEntityLinks(e.Links),
			"name":              utils.StringValue(e.Name),
			"description":       utils.StringValue(e.Description),
			"display_name":      utils.StringValue(e.DisplayName),
			"client_name":       utils.StringValue(e.ClientName),
			"search_url":        utils.StringValue(e.SearchURL),
			"created_time":      utils.TimeStringValue(e.CreatedTime),
			"last_updated_time": utils.TimeStringValue(e.LastUpdatedTime),
			"created_by":        utils.StringValue(e.CreatedBy),
			"attribute_list":    flattenAttributeList(e.AttributeList),
			"is_logical_and_supported_for_attributes": utils.BoolValue(e.IsLogicalAndSupportedForAttributes),
		}
		result = append(result, m)
	}
	return result
}
