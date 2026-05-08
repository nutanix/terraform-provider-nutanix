package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVmGuestCustomizationProfilesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmGuestCustomizationProfilesV2Read,
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
			"vm_guest_customization_profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixVmGuestCustomizationProfileV2(),
			},
		},
	}
}

func DatasourceNutanixVmGuestCustomizationProfilesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	var filter, orderBy, selectQ *string
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
	if selectQf, ok := d.GetOk("select"); ok {
		selectQ = utils.StringPtr(selectQf.(string))
	}

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.ListVmGuestCustomizationProfiles(page, limit, filter, orderBy, selectQ)
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profiles: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("vm_guest_customization_profiles", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of VM Guest Customization Profiles.",
		}}
	}

	profiles := resp.Data.GetValue().([]config.VmGuestCustomizationProfile)

	if err := d.Set("vm_guest_customization_profiles", flattenVmGcProfileEntities(profiles)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVmGcProfileEntities(profiles []config.VmGuestCustomizationProfile) []map[string]interface{} {
	if len(profiles) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(profiles))
	for i, p := range profiles {
		profileMap := map[string]interface{}{
			"ext_id":      utils.StringValue(p.ExtId),
			"tenant_id":   utils.StringValue(p.TenantId),
			"links":       flattenVmGcProfileLinks(p.Links),
			"name":        utils.StringValue(p.Name),
			"description": utils.StringValue(p.Description),
			"config":      flattenVmGcProfileConfig(p.Config),
			"created_by":  flattenVmGcProfileUserReference(p.CreatedBy),
			"updated_by":  flattenVmGcProfileUserReference(p.UpdatedBy),
		}
		if p.CreateTime != nil {
			profileMap["create_time"] = p.CreateTime.String()
		}
		if p.UpdateTime != nil {
			profileMap["update_time"] = p.UpdateTime.String()
		}
		result[i] = profileMap
	}
	return result
}
