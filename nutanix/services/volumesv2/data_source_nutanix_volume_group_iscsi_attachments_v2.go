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

func DatasourceNutanixVolumeGroupIscsiAttachmentsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the list of external iSCSI attachments for a Volume Group identified by {extId}. Deprecated: This API has been deprecated.",
		ReadContext: DatasourceNutanixVolumeGroupIscsiAttachmentsV2Read,
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of a Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"iscsi_attachments": {
				Description: "List of external iSCSI attachments for the Volume Group.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Description: "The external identifier of an iSCSI client.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cluster_reference": {
							Description: "The UUID of the cluster that will host the iSCSI client. This field is read-only.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVolumeGroupIscsiAttachmentsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id").(string)

	resp, err := conn.VolumeAPIInstance.ListExternalIscsiAttachmentsByVolumeGroupId(utils.StringPtr(volumeGroupExtID), nil, nil, nil, nil, nil, nil)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group iSCSI Attachments : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("iscsi_attachments", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return nil
	}

	getResp := resp.Data.GetValue().([]volumesClient.IscsiClientAttachment)

	if err := d.Set("iscsi_attachments", flattenIscsiClientAttachments(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenIscsiClientAttachments(attachments []volumesClient.IscsiClientAttachment) []map[string]interface{} {
	if len(attachments) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(attachments))
	for i, a := range attachments {
		entry := map[string]interface{}{}
		if a.ExtId != nil {
			entry["ext_id"] = utils.StringValue(a.ExtId)
		}
		if a.ClusterReference != nil {
			entry["cluster_reference"] = utils.StringValue(a.ClusterReference)
		}
		result[i] = entry
	}
	return result
}
