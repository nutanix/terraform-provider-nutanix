package volumesv2

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// List all the VM attachments for a Volume Group.
func DataSourceNutanixVolumeGroupVmsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the list of VM attachments for a Volume Group identified by {extId}.",
		ReadContext: DataSourceNutanixVolumeGroupVmsV4Read,

		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "The external identifier of the volume group.",
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
				Description: "A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions. For example, filter '$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields: extId",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"orderby": {
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:  extId",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vms_attachments": {
				Description: "List of Volume Groups.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ext_id": {
							Description: "A globally unique identifier of an instance that is suitable for external consumption.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"links": {
							Description: "A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Description: "The URL at which the entity described by the link can be accessed.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"rel": {
										Description: "A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of \"self\" identifies the URL for the object.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// DataSourceNutanixVolumeGroupVmsV4Read
func DataSourceNutanixVolumeGroupVmsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("ext_id")

	var filter, orderBy *string
	var page, limit *int

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
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}

	// get the volume groups response
	resp, err := conn.VolumeAPIInstance.ListVmAttachmentsByVolumeGroupId(utils.StringPtr(volumeGroupExtID.(string)), page, limit, filter, orderBy)

	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching volumes : %v", errorMessage["message"])
	}

	vmsAttachmentsResp := resp.Data

	if vmsAttachmentsResp != nil {
		// set the volume groups data in the terraform resource
		if err := d.Set("vms_attachments", flattenVolumeGroupVmsEntities(vmsAttachmentsResp.GetValue().([]volumesClient.VmAttachment))); err != nil {
			return diag.FromErr(err)
		}

	} else {
		// set the volume groups data in the terraform resource
		d.Set("volumes", make([]volumesClient.VolumeGroup, 0))
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVolumeGroupVmsEntities(vms []volumesClient.VmAttachment) []interface{} {
	if len(vms) > 0 {
		vmAttachmentList := make([]interface{}, len(vms))

		for k, v := range vms {
			vmAttachment := make(map[string]interface{})

			if v.ExtId != nil {
				vmAttachment["ext_id"] = v.ExtId
			}
			vmAttachmentList[k] = vmAttachment

		}
		return vmAttachmentList
	}
	return nil
}
