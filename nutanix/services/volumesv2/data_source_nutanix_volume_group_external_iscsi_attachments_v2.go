package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVolumeGroupExternalIscsiAttachmentsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the list of external iSCSI attachments for a Volume Group identified by {extId}. Deprecated: This API has been deprecated.",
		ReadContext: DatasourceNutanixVolumeGroupExternalIscsiAttachmentsV2Read,
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of a Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"external_iscsi_attachments": {
				Description: "List of external iSCSI attachments.",
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

func DatasourceNutanixVolumeGroupExternalIscsiAttachmentsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id").(string)

	resp, err := conn.VolumeAPIInstance.ListExternalIscsiAttachmentsByVolumeGroupId(utils.StringPtr(volumeGroupExtID), nil, nil, nil, nil, nil, nil)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group External iSCSI Attachments: %v", err)
	}

	var attachmentsList []map[string]interface{}

	if resp.Data != nil && resp.Data.GetValue() != nil {
		if attachments, ok := resp.Data.GetValue().([]volumesClient.IscsiClientAttachment); ok {
			attachmentsList = make([]map[string]interface{}, len(attachments))
			for i, att := range attachments {
				entry := map[string]interface{}{}
				if att.ExtId != nil {
					entry["ext_id"] = utils.StringValue(att.ExtId)
				}
				if att.ClusterReference != nil {
					entry["cluster_reference"] = utils.StringValue(att.ClusterReference)
				}
				attachmentsList[i] = entry
			}
		}
	}

	if err := d.Set("external_iscsi_attachments", attachmentsList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(volumeGroupExtID)
	return nil
}
