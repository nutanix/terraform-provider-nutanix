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

func DatasourceNutanixVolumeGroupCategoryAssociationsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the category details that are associated with the Volume Group identified by {volumeGroupExtId}. Deprecated: This API has been deprecated.",
		ReadContext: DatasourceNutanixVolumeGroupCategoryAssociationsV2Read,
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of a Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"category_associations": {
				Description: "List of category details associated with the Volume Group.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Description: "A globally unique identifier of an instance that is suitable for external consumption.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Name of the entity represented by this reference.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"entity_type": {
							Description: "The entity type.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"uris": {
							Description: "URI of entity represented by this reference.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVolumeGroupCategoryAssociationsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id").(string)

	resp, err := conn.VolumeAPIInstance.ListCategoryAssociationsByVolumeGroupId(utils.StringPtr(volumeGroupExtID), nil, nil)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group Category Associations : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("category_associations", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return nil
	}

	getResp := resp.Data.GetValue().([]volumesClient.CategoryDetails)

	if err := d.Set("category_associations", flattenCategoryDetails(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenCategoryDetails(categories []volumesClient.CategoryDetails) []map[string]interface{} {
	if len(categories) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(categories))
	for i, c := range categories {
		entry := map[string]interface{}{}
		if c.ExtId != nil {
			entry["ext_id"] = utils.StringValue(c.ExtId)
		}
		if c.Name != nil {
			entry["name"] = utils.StringValue(c.Name)
		}
		if c.EntityType != nil {
			entry["entity_type"] = flattenCategoryEntityType(c.EntityType)
		}
		if c.Uris != nil {
			entry["uris"] = c.Uris
		}
		result[i] = entry
	}
	return result
}

func flattenCategoryEntityType(entityType interface{}) string {
	if entityType == nil {
		return ""
	}
	return ""
}
