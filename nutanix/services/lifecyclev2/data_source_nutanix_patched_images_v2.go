package lifecyclev2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixPatchedImagesV2 reads a list of patched images. Each element
// exposes the claim_token_ext_id cross-resource reference.
func DatasourceNutanixPatchedImagesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixPatchedImagesV2Read,
		Schema: map[string]*schema.Schema{
			"page":     {Type: schema.TypeInt, Optional: true},
			"limit":    {Type: schema.TypeInt, Optional: true},
			"filter":   {Type: schema.TypeString, Optional: true},
			"order_by": {Type: schema.TypeString, Optional: true},
			"select":   {Type: schema.TypeString, Optional: true},
			"patched_images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixPatchedImageV2(),
			},
		},
	}
}

func DatasourceNutanixPatchedImagesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	var page, limit *int
	var filter, orderBy, selectQ *string

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
		selectQ = utils.StringPtr(v.(string))
	}

	resp, err := conn.PatchedImagesAPIInstance.ListPatchedImages(page, limit, filter, orderBy, selectQ)
	if err != nil {
		return diag.Errorf("error while listing patched images : %v", err)
	}

	d.SetId(resource.UniqueId())

	if resp.Data == nil {
		if err := d.Set("patched_images", make([]map[string]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of patched images.",
		}}
	}

	images := resp.Data.GetValue().([]import1.PatchedImage)

	if err := d.Set("patched_images", flattenPatchedImages(images)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenPatchedImages(images []import1.PatchedImage) []map[string]interface{} {
	if len(images) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(images))
	for i := range images {
		img := images[i]
		m := map[string]interface{}{
			"ext_id":                      utils.StringValue(img.ExtId),
			"claim_token_ext_id":          utils.StringValue(img.ClaimTokenExtId),
			"name":                        utils.StringValue(img.Name),
			"version":                     utils.StringValue(img.Version),
			"image_details":               flattenImageDetails(img.ImageDetails),
			"node_configurations":         flattenNodeConfigurations(img.NodeConfigurations),
			"patched_iso_url":             utils.StringValue(img.PatchedIsoUrl),
			"patched_iso_sha256_checksum": utils.StringValue(img.PatchedIsoSha256Checksum),
			"owner_ext_id":                utils.StringValue(img.OwnerExtId),
			"tenant_id":                   utils.StringValue(img.TenantId),
			"links":                       flattenLifecycleLinks(img.Links),
		}
		if img.HostType != nil {
			m["host_type"] = img.HostType.GetName()
		}
		if img.CreatedTime != nil {
			m["created_time"] = img.CreatedTime.Format(time.RFC3339)
		}
		out = append(out, m)
	}
	return out
}
