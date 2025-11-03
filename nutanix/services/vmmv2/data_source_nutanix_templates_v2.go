package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import5 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixTemplatesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixTemplatesV2Read,
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
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixTemplateV2(),
			},
		},
	}
}

func DatasourceNutanixTemplatesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

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
	if selectQf, ok := d.GetOk("select"); ok {
		selectQ = utils.StringPtr(selectQf.(string))
	} else {
		selectQ = nil
	}
	resp, err := conn.TemplatesAPIInstance.ListTemplates(page, limit, filter, orderBy, selectQ)
	if err != nil {
		return diag.Errorf("error while fetching templates : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("templates", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of templates.",
		}}
	}
	getResp := resp.Data.GetValue().([]import5.Template)

	if err := d.Set("templates", flattenTemplatesEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenTemplatesEntities(pr []import5.Template) []interface{} {
	if len(pr) > 0 {
		temps := make([]interface{}, len(pr))

		for k, v := range pr {
			temp := make(map[string]interface{})

			temp["tenant_id"] = v.TenantId
			temp["links"] = flattenAPILink(v.Links)
			temp["ext_id"] = v.ExtId
			temp["template_name"] = v.TemplateName
			temp["template_description"] = v.TemplateDescription
			temp["template_version_spec"] = flattenTemplateVersionSpec(v.TemplateVersionSpec)
			temp["guest_update_status"] = flattenGuestUpdateStatus(v.GuestUpdateStatus)
			if v.CreateTime != nil {
				t := v.CreateTime
				temp["create_time"] = t.String()
			}
			if v.UpdateTime != nil {
				t := v.UpdateTime
				temp["update_time"] = t.String()
			}
			temp["created_by"] = flattenTemplateUser(v.CreatedBy)
			temp["updated_by"] = flattenTemplateUser(v.UpdatedBy)
			temp["category_ext_ids"] = v.CategoryExtIds

			temps[k] = temp
		}
		return temps
	}
	return nil
}
