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

// DatasourceNutanixVolumeGroupVmsV2 lists the VM attachments for a Volume Group.
//
// Deprecated: This API has been deprecated.
func DatasourceNutanixVolumeGroupVmsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the list of VM attachments for a Volume Group identified by {extId}. Deprecated: This API has been deprecated.",
		ReadContext: DatasourceNutanixVolumeGroupVmsV2Read,
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of a Volume Group.",
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
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"orderby": {
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"attachments": {
				Description: "List of VM attachments for a Volume Group identified by {extId}.",
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

func DatasourceNutanixVolumeGroupVmsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	var filter, orderBy *string
	var page, limit *int

	volumeGroupExtID := d.Get("volume_group_ext_id")

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
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("orderby"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}

	// get the VM attachments response
	resp, err := conn.VolumeAPIInstance.ListVmAttachmentsByVolumeGroupId(utils.StringPtr(volumeGroupExtID.(string)), page, limit, filter, orderBy)
	if err != nil {
		return diag.Errorf("error while fetching VM attachments for the volume group : %v", err)
	}

	// extract the VM attachments data from the response
	if resp.Data == nil {
		if err := d.Set("attachments", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of VM attachments.",
		}}
	}

	getResp := resp.Data.GetValue().([]volumesClient.VmAttachment)
	// set the VM attachments data in the terraform resource
	if err := d.Set("attachments", flattenVMAttachmentEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVMAttachmentEntities(vmAttachments []volumesClient.VmAttachment) []interface{} {
	if len(vmAttachments) > 0 {
		vmAttachmentList := make([]interface{}, len(vmAttachments))

		for k, v := range vmAttachments {
			vmAttachment := make(map[string]interface{})

			if v.ExtId != nil {
				vmAttachment["ext_id"] = v.ExtId
			}
			if v.Index != nil {
				vmAttachment["index"] = v.Index
			}

			vmAttachmentList[k] = vmAttachment
		}
		return vmAttachmentList
	}
	return nil
}
