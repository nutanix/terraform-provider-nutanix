package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/error"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixUserRevokeKeyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixUserRevokeKeyV2Create,
		ReadContext:   resourceNutanixUserRevokeKeyV2Read,
		UpdateContext: resourceNutanixUserRevokeKeyV2Update,
		DeleteContext: resourceNutanixUserRevokeKeyV2Delete,
		Schema: map[string]*schema.Schema{
			"user_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
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

func resourceNutanixUserRevokeKeyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	var userExtID *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtID = utils.StringPtr(v.(string))
	}

	var ExtID *string
	if v, ok := d.GetOk("ext_id"); ok {
		ExtID = utils.StringPtr(v.(string))
	}

	resp, err := conn.UsersAPIInstance.RevokeUserKey(userExtID, ExtID)
	if err != nil {
		return diag.Errorf("error while revoking the user key: %v | ExtId: %s | userExtId: %s", err, *ExtID, *userExtID)
	}

	revokeConfig := resp.Data.GetValue().(import1.AppMessage)
	if revokeConfig.Message != nil {
		d.Set("message", revokeConfig.Message)
	}
	if revokeConfig.Severity != nil {
		d.Set("severity", revokeConfig.Severity)
	}
	if revokeConfig.Code != nil {
		d.Set("code", revokeConfig.Code)
	}
	if revokeConfig.Locale != nil {
		d.Set("locale", revokeConfig.Locale)
	}
	if revokeConfig.ErrorGroup != nil {
		d.Set("error_group", revokeConfig.ErrorGroup)
	}
	if revokeConfig.ArgumentsMap != nil {
		d.Set("arguments_map", flattenArgumentsMap(revokeConfig.ArgumentsMap))
	} else {
		d.Set("arguments_map", []map[string]interface{}{})
	}
	d.SetId(*ExtID)
	return nil
}

func resourceNutanixUserRevokeKeyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixUserRevokeKeyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixUserRevokeKeyV2Create(ctx, d, meta)
}

func resourceNutanixUserRevokeKeyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
