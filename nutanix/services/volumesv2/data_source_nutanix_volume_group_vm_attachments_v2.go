package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVolumeGroupVmAttachmentsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the list of VM attachments for a Volume Group identified by {extId}. Deprecated: This API has been deprecated.",
		ReadContext: DatasourceNutanixVolumeGroupVmAttachmentsV2Read,
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of a Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vm_attachments": {
				Description: "List of VM attachments.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Description: "The external identifier of the VM.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"index": {
							Description: "The index on the SCSI bus to attach the VM to the Volume Group. This is an optional field.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVolumeGroupVmAttachmentsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id").(string)

	resp, err := conn.VolumeAPIInstance.ListVmAttachmentsByVolumeGroupId(utils.StringPtr(volumeGroupExtID), nil, nil, nil, nil)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group VM Attachments: %v", err)
	}

	var vmAttachmentsList []map[string]interface{}

	if resp.Data != nil && resp.Data.GetValue() != nil {
		if vmAttachments, ok := resp.Data.GetValue().([]volumesClient.VmAttachment); ok {
			vmAttachmentsList = make([]map[string]interface{}, len(vmAttachments))
			for i, att := range vmAttachments {
				entry := map[string]interface{}{}
				if att.ExtId != nil {
					entry["ext_id"] = utils.StringValue(att.ExtId)
				}
				if att.Index != nil {
					entry["index"] = *att.Index
				}
				vmAttachmentsList[i] = entry
			}
		}
	}

	if err := d.Set("vm_attachments", vmAttachmentsList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(volumeGroupExtID)
	return nil
}
