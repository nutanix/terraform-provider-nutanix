package licensingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/agreements"
	import2 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/error"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixAcceptEULAV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixAcceptEULACreateV2,
		ReadContext:   resourceNutanixAcceptEULAReadV2,
		DeleteContext: resourceNutanixAcceptEULADeleteV2,
		UpdateContext: resourceNutanixAcceptEULAUpdateV2,
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"company_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arguments_map": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixAcceptEULACreateV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	body := &import1.EndUser{}
	if v, ok := d.GetOk("user_name"); ok {
		body.UserName = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("login_id"); ok {
		body.LoginId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("job_title"); ok {
		body.JobTitle = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("company_name"); ok {
		body.CompanyName = utils.StringPtr(v.(string))
	}

	resp, err := conn.LicensingEULAAPIInstance.AddUser(body)
	if err != nil {
		return diag.FromErr(err)
	}

	appMessages, ok := resp.Data.GetValue().([]import2.AppMessage)
	if !ok || len(appMessages) == 0 {
		return diag.Errorf("unexpected response type or empty AppMessage slice")
	}
	getResp := appMessages[0]
	if getResp.Message != nil {
		d.Set("message", getResp.Message)
	}
	if getResp.Severity != nil {
		d.Set("severity", getResp.Severity.GetName())
	}
	if getResp.Code != nil {
		d.Set("code", getResp.Code)
	}
	if getResp.Locale != nil {
		d.Set("locale", getResp.Locale)
	}
	if getResp.ErrorGroup != nil {
		d.Set("error_group", getResp.ErrorGroup)
	}
	if getResp.ArgumentsMap != nil {
		d.Set("arguments_map", flattenArgumentsMap(getResp.ArgumentsMap))
	} else {
		d.Set("arguments_map", []map[string]interface{}{})
	}

	d.SetId(utils.GenUUID())
	return resourceNutanixAcceptEULAReadV2(ctx, d, meta)
}

func resourceNutanixAcceptEULAReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixAcceptEULADeleteV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixAcceptEULAUpdateV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func flattenArgumentsMap(argumentsMap map[string]string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(argumentsMap))
	for key, value := range argumentsMap {
		result = append(result, map[string]interface{}{
			"property_name": key,
			"value":         value,
		})
	}
	return result
}
