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
			"profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixVmGuestCustomizationProfileV2(),
			},
		},
	}
}

func DatasourceNutanixVmGuestCustomizationProfilesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	var filter, orderBy, selectParam *string
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
	if sel, ok := d.GetOk("select"); ok {
		selectParam = utils.StringPtr(sel.(string))
	}

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.ListVmGuestCustomizationProfiles(page, limit, filter, orderBy, selectParam)
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profiles: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("profiles", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return nil
	}

	getResp := resp.Data.GetValue().([]config.VmGuestCustomizationProfile)

	if err := d.Set("profiles", flattenVmGuestCustomizationProfileEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenVmGuestCustomizationProfileEntities(profiles []config.VmGuestCustomizationProfile) []interface{} {
	if len(profiles) == 0 {
		return nil
	}
	result := make([]interface{}, len(profiles))
	for i, p := range profiles {
		result[i] = flattenVmGuestCustomizationProfileEntity(p)
	}
	return result
}
