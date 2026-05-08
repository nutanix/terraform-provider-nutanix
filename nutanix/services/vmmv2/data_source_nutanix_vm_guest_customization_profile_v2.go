package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVmGuestCustomizationProfileV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmGuestCustomizationProfileV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": schemaForVmGcProfileConfig(true),
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"updated_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVmGuestCustomizationProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.GetVmGuestCustomizationProfileById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile: %v", err)
	}

	profile := resp.Data.GetValue().(config.VmGuestCustomizationProfile)

	if err := d.Set("ext_id", utils.StringValue(profile.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(profile.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenVmGcProfileLinks(profile.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", utils.StringValue(profile.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", utils.StringValue(profile.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config", flattenVmGcProfileConfig(profile.Config)); err != nil {
		return diag.FromErr(err)
	}
	if profile.CreateTime != nil {
		if err := d.Set("create_time", profile.CreateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if profile.UpdateTime != nil {
		if err := d.Set("update_time", profile.UpdateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("created_by", flattenVmGcProfileUserReference(profile.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_by", flattenVmGcProfileUserReference(profile.UpdatedBy)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(profile.ExtId))
	return nil
}
