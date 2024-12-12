package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// List all the category details that are associated with the Volume Group.
func DatasourceNutanixVolumeCategoryDetailsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVolumeCategoryDetailsV2Read,

		Description: "Query the category details that are associated with the Volume Group identified by {volumeGroupExtID}.",
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "The external identifier of the Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"page": {
				Description: "A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"limit": {
				Description: "A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"category_details": {
				Description: "List of all category details that are associated with the Volume Group.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Description: "The external identifier of the category detail",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of the category detail.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"uris": {
							Description: "The uri list of the category detail.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Description: "",
								Type:        schema.TypeList,
							},
						},
						"entity_type": {
							Description: "The Entity Type of the category detail.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVolumeCategoryDetailsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	var page, limit *int

	volumeGroupExtID := d.Get("ext_id")

	// initialize the query parameters
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

	// get the volume groups response
	resp, err := conn.VolumeAPIInstance.ListCategoryAssociationsByVolumeGroupId(utils.StringPtr(volumeGroupExtID.(string)), page, limit)
	if err != nil {
		return diag.Errorf("error while fetching volumes : %v", err)
	}

	// extract the volume groups data from the response
	getResp := resp.Data.GetValue().([]volumesClient.CategoryDetails)

	// set the volume groups data in the terraform resource
	if err := d.Set("category_details", flattenCategoryDetails(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenCategoryDetails(categories []volumesClient.CategoryDetails) []interface{} {
	if len(categories) > 0 {
		categoriesList := make([]interface{}, len(categories))

		for k, v := range categories {
			category := make(map[string]interface{})

			category["ext_id"] = v.ExtId
			category["name"] = v.Name
			category["uris"] = v.Uris
			category["entity_type"] = v.EntityType.GetName()

			categoriesList[k] = category
		}
		return categoriesList
	}
	return nil
}
