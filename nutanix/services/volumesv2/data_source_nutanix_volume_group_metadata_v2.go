package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesConfig "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVolumeGroupMetadataV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query for metadata information which is associated with the Volume Group identified by {extId}. Deprecated: This API has been deprecated.",
		ReadContext: DatasourceNutanixVolumeGroupMetadataV2Read,
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of a Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"category_ids": {
				Description: "A list of globally unique identifiers that represent all the categories the resource is associated with.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_reference_id": {
				Description: "A globally unique identifier that represents the owner of this resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner_user_name": {
				Description: "The userName of the owner of this resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"project_name": {
				Description: "The name of the project this resource belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"project_reference_id": {
				Description: "A globally unique identifier that represents the project this resource belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func DatasourceNutanixVolumeGroupMetadataV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id").(string)

	resp, err := conn.VolumeAPIInstance.GetVolumeGroupMetadataById(utils.StringPtr(volumeGroupExtID))
	if err != nil {
		return diag.Errorf("error while fetching Volume Group Metadata: %v", err)
	}

	if resp.Data != nil && resp.Data.GetValue() != nil {
		getResp := resp.Data.GetValue().(volumesConfig.Metadata)

		if err := d.Set("category_ids", getResp.CategoryIds); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("owner_reference_id", getResp.OwnerReferenceId); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("owner_user_name", getResp.OwnerUserName); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("project_name", getResp.ProjectName); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("project_reference_id", getResp.ProjectReferenceId); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(volumeGroupExtID)
	return nil
}
