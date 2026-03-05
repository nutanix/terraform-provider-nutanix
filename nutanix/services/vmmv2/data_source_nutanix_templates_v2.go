package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/content"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/templates"
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

	listTemplatesRequest := import2.ListTemplatesRequest{}

	if v, ok := d.GetOk("page"); ok {
		listTemplatesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listTemplatesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listTemplatesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listTemplatesRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listTemplatesRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.TemplatesAPIInstance.ListTemplates(ctx, &listTemplatesRequest)
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
