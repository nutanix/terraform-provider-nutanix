package volumesv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Description: "List of category associations.",
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
		return diag.Errorf("error while fetching Volume Group Category Associations: %v", err)
	}

	var categoryList []map[string]interface{}

	if resp.Data != nil && resp.Data.GetValue() != nil {
		if categories, ok := resp.Data.GetValue().([]volumesClient.CategoryDetails); ok {
			categoryList = make([]map[string]interface{}, len(categories))
			for i, cat := range categories {
				entry := map[string]interface{}{}
				if cat.ExtId != nil {
					entry["ext_id"] = utils.StringValue(cat.ExtId)
				}
				if cat.Name != nil {
					entry["name"] = utils.StringValue(cat.Name)
				}
				if cat.EntityType != nil {
					entry["entity_type"] = flattenCategoryEntityType(cat.EntityType)
				}
				entry["uris"] = cat.Uris
				categoryList[i] = entry
			}
		}
	}

	if err := d.Set("category_associations", categoryList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(volumeGroupExtID)
	return nil
}

func flattenCategoryEntityType(entityType interface{}) string {
	if entityType == nil {
		return ""
	}
	return fmt.Sprintf("%v", entityType)
}
