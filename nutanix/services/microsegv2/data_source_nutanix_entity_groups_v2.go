package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEntityGroupsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixEntityGroupsV2Read,
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
			"entity_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForEntityGroupListElem(),
			},
		},
	}
}

func schemaForEntityGroupListElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
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
			"allowed_config": schemaForAllowedConfig(),
			"except_config":  schemaForExceptConfig(),
			"policy_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixEntityGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	var page, limit *int
	var filter, orderBy, selectVal *string

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
		selectVal = utils.StringPtr(v.(string))
	}

	resp, err := conn.EntityGroupsAPIInstance.ListEntityGroups(page, limit, filter, orderBy, selectVal)
	if err != nil {
		return diag.Errorf("error while listing Entity Groups: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("entity_groups", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	listVal, ok := resp.Data.GetValue().([]import2.EntityGroup)
	if !ok {
		if err := d.Set("entity_groups", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	if err := d.Set("entity_groups", flattenEntityGroups(listVal)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenEntityGroups(groups []import2.EntityGroup) []map[string]interface{} {
	if len(groups) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(groups))
	for _, g := range groups {
		m := map[string]interface{}{
			"ext_id":         utils.StringValue(g.ExtId),
			"name":           utils.StringValue(g.Name),
			"description":    utils.StringValue(g.Description),
			"allowed_config": flattenAllowedConfig(g.AllowedConfig),
			"except_config":  flattenExceptConfig(g.ExceptConfig),
			"policy_ext_ids": g.PolicyExtIds,
			"owner_ext_id":   utils.StringValue(g.OwnerExtId),
			"tenant_id":      utils.StringValue(g.TenantId),
			"links":          flattenLinksEntityGroup(g.Links),
		}
		if g.CreationTime != nil {
			m["creation_time"] = g.CreationTime.Format("2006-01-02T15:04:05.000Z")
		}
		if g.LastUpdateTime != nil {
			m["last_update_time"] = g.LastUpdateTime.Format("2006-01-02T15:04:05.000Z")
		}
		result = append(result, m)
	}
	return result
}
