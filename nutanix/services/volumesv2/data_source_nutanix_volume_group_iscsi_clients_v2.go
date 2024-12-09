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

// List all the iSCSI attachments associated with the given Volume Group.
func DatasourceNutanixVolumeGroupIscsiClientsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the list of external iSCSI attachments for a Volume Group identified by {extId}.",
		ReadContext: DatasourceNutanixVolumeGroupIscsiClientsV2Read,
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
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"orderby": {
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields: clusterReference, extId",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields: clusterReference, extId",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"iscsi_clients": {
				Description: "List of the iSCSI attachments associated with the given Volume Group.",
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

func DatasourceNutanixVolumeGroupIscsiClientsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	var filter, orderBy, expand, selects *string
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
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	// get the volume group iscsi clients
	resp, err := conn.VolumeAPIInstance.ListExternalIscsiAttachmentsByVolumeGroupId(utils.StringPtr(volumeGroupExtID.(string)), page, limit, filter, orderBy, expand, selects)

	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching External Iscsi Attachments : %v", errorMessage["message"])
	}

	diskResp := resp.Data

	// extract the volume groups data from the response
	if diskResp != nil {
		// set the volume groups iscsi clients  data in the terraform resource
		if err := d.Set("iscsi_clients", flattenVolumeIscsiClientsEntities(diskResp.GetValue().([]volumesClient.IscsiClientAttachment))); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(resource.UniqueId())
	return nil

}

func flattenVolumeIscsiClientsEntities(iscsiClientAttachments []volumesClient.IscsiClientAttachment) []interface{} {
	if len(iscsiClientAttachments) > 0 {
		iscsiClientList := make([]interface{}, len(iscsiClientAttachments))

		for k, v := range iscsiClientAttachments {
			iscsiClient := make(map[string]interface{})

			if v.ExtId != nil {
				iscsiClient["ext_id"] = v.ExtId
			}

			if v.ClusterReference != nil {
				iscsiClient["cluster_reference"] = v.ClusterReference
			}

			iscsiClientList[k] = iscsiClient
		}
		return iscsiClientList
	}
	return nil
}
