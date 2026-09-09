package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixInstallerImagesV2 returns a paginated list of all installer images
// registered in Foundation Central. The element schema is identical to the
// singular image datasource, so it is reused via Elem.
func DatasourceNutanixInstallerImagesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixInstallerImagesV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A URL query parameter that specifies the total number of records returned in the result set.",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to filter a collection of resources.",
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.",
			},
			"select": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.",
			},
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of installer images registered in Foundation Central.",
				Elem:        DatasourceNutanixInstallerImageV2(),
			},
		},
	}
}

func DatasourceNutanixInstallerImagesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	// initialize query params
	var filter, orderBy, selects *string
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
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	}

	resp, err := conn.InstallerImagesAPIInstance.ListImages(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching images : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("images", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of images.",
		}}
	}

	getResp := resp.Data.GetValue().([]config.Image)

	if err := d.Set("images", flattenImagesEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// flattenImagesEntities converts the SDK image list into the datasource schema shape.
func flattenImagesEntities(images []config.Image) []map[string]interface{} {
	if len(images) == 0 {
		return nil
	}

	imageList := make([]map[string]interface{}, len(images))
	for i, image := range images {
		img := make(map[string]interface{})
		img["ext_id"] = utils.StringValue(image.ExtId)
		img["name"] = utils.StringValue(image.Name)
		if image.Type != nil {
			img["type"] = image.Type.GetName()
		}
		if image.Source != nil {
			img["source"] = image.Source.GetName()
		}
		img["url"] = utils.StringValue(image.Url)
		img["version"] = utils.StringValue(image.Version)
		img["certificate_chain"] = utils.StringValue(image.CertificateChain)
		img["metadata_download_url"] = utils.StringValue(image.MetadataDownloadUrl)
		img["checksum"] = flattenImageChecksum(image.Checksum)
		if image.FileStatus != nil {
			img["file_status"] = image.FileStatus.GetName()
		}
		if image.MetadataStatus != nil {
			img["metadata_status"] = image.MetadataStatus.GetName()
		}
		if image.CreatedTime != nil {
			img["created_time"] = image.CreatedTime.String()
		}
		img["tenant_id"] = utils.StringValue(image.TenantId)
		img["links"] = flattenImageLinks(image.Links)
		imageList[i] = img
	}
	return imageList
}
