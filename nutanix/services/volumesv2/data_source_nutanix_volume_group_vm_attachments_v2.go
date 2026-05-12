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
				Description: "List of VM attachments for the Volume Group.",
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
		return diag.Errorf("error while fetching Volume Group VM Attachments : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("vm_attachments", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return nil
	}

	getResp := resp.Data.GetValue().([]volumesClient.VmAttachment)

	if err := d.Set("vm_attachments", flattenVmAttachments(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVmAttachments(attachments []volumesClient.VmAttachment) []map[string]interface{} {
	if len(attachments) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(attachments))
	for i, a := range attachments {
		entry := map[string]interface{}{}
		if a.ExtId != nil {
			entry["ext_id"] = utils.StringValue(a.ExtId)
		}
		if a.Index != nil {
			entry["index"] = *a.Index
		}
		result[i] = entry
	}
	return result
}
