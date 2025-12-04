package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesClientResponse "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/response"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// List all the Volume Groups.
func DatasourceNutanixVolumeGroupsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the list of Volume Groups.",
		ReadContext: DatasourceNutanixVolumeGroupsV2Read,
		Schema: map[string]*schema.Schema{
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
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields: clusterReference, extId, name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields: clusterReference, extId, name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"volumes": {
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
						"name": {
							Description: "Volume Group name. This is an optional field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Volume Group description. This is an optional field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"should_load_balance_vm_attachments": {
							Description: "Indicates whether to enable Volume Group load balancing for VM attachments. This cannot be enabled if there are iSCSI client attachments already associated with the Volume Group, and vice-versa. This is an optional field.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"sharing_status": {
							Description: "Indicates whether the Volume Group can be shared across multiple iSCSI initiators. The mode cannot be changed from SHARED to NOT_SHARED on a Volume Group with multiple attachments. Similarly, a Volume Group cannot be associated with more than one attachment as long as it is in exclusive mode. This is an optional field",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"target_name": {
							Description: "Name of the external client target that will be visible and accessible to the client. This is an optional field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"enabled_authentications": {
							Description: "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"iscsi_features": {
							Description: "iSCSI specific settings for the Volume Group. This is an optional field.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled_authentications": {
										Description: "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"created_by": {
							Description: "Service/user who created this Volume Group. This is an optional field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cluster_reference": {
							Description: "The UUID of the cluster that will host the Volume Group. This is a mandatory field for creating a Volume Group on Prism Central.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"storage_features": {
							Description: "Storage optimization features which must be enabled on the Volume Group. This is an optional field.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flash_mode": {
										Description: "Once configured, this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.",
										Type:        schema.TypeList,
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Description: "Indicates whether the flash mode is enabled for the Volume Group.",
													Type:        schema.TypeBool,
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"usage_type": {
							Description: "Expected usage type for the Volume Group. This is an indicative hint on how the caller will consume the Volume Group. This is an optional field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"is_hidden": {
							Description: "Indicates whether the Volume Group is meant to be hidden or not. This is an optional field. If omitted, the VG will not be hidden.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVolumeGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	var filter, orderBy, expand, selects *string
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

	// get the volume groups response
	resp, err := conn.VolumeAPIInstance.ListVolumeGroups(page, limit, filter, orderBy, expand, selects)
	if err != nil {
		return diag.Errorf("error while fetching volumes : %v", err)
	}

	volumesResp := resp.Data

	if volumesResp != nil {
		// set the volume groups data in the terraform resource
		if err := d.Set("volumes", flattenVolumesEntities(volumesResp.GetValue().([]volumesClient.VolumeGroup))); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return nil
	}

	// set the volume groups data to empty list
	d.Set("volumes", make([]volumesClient.VolumeGroup, 0))

	d.SetId(utils.GenUUID())

	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "🫙 No data found.",
		Detail:   "The API returned an empty list of volume groups.",
	}}
}

func flattenVolumesEntities(volumeGroups []volumesClient.VolumeGroup) []interface{} {
	if len(volumeGroups) > 0 {
		volumeGroupList := make([]interface{}, len(volumeGroups))

		for k, v := range volumeGroups {
			volumeGroup := make(map[string]interface{})

			if v.TenantId != nil {
				volumeGroup["tenant_id"] = v.TenantId
			}
			if v.ExtId != nil {
				volumeGroup["ext_id"] = v.ExtId
			}
			if v.Links != nil {
				volumeGroup["links"] = flattenLinks(v.Links)
			}
			if v.Name != nil {
				volumeGroup["name"] = v.Name
			}
			if v.Description != nil {
				volumeGroup["description"] = v.Description
			}
			if v.ShouldLoadBalanceVmAttachments != nil {
				volumeGroup["should_load_balance_vm_attachments"] = v.ShouldLoadBalanceVmAttachments
			}
			if v.SharingStatus != nil {
				volumeGroup["sharing_status"] = flattenSharingStatus(v.SharingStatus)
			}
			if v.TargetName != nil {
				volumeGroup["target_name"] = v.TargetName
			}
			if v.EnabledAuthentications != nil {
				volumeGroup["enabled_authentications"] = flattenEnabledAuthentications(v.EnabledAuthentications)
			}
			if v.IscsiFeatures != nil {
				volumeGroup["iscsi_features"] = flattenIscsiFeatures(v.IscsiFeatures)
			}
			if v.CreatedBy != nil {
				volumeGroup["created_by"] = v.CreatedBy
			}
			if v.ClusterReference != nil {
				volumeGroup["cluster_reference"] = v.ClusterReference
			}
			if v.StorageFeatures != nil {
				volumeGroup["storage_features"] = flattenStorageFeatures(v.StorageFeatures)
			}
			if v.UsageType != nil {
				volumeGroup["usage_type"] = flattenUsageType(v.UsageType)
			}
			if v.IsHidden != nil {
				volumeGroup["is_hidden"] = v.IsHidden
			}

			volumeGroupList[k] = volumeGroup
		}
		return volumeGroupList
	}
	return nil
}

func flattenLinks(apiLinks []volumesClientResponse.ApiLink) []map[string]interface{} {
	if len(apiLinks) > 0 {
		apiLinkList := make([]map[string]interface{}, len(apiLinks))

		for k, v := range apiLinks {
			links := map[string]interface{}{}
			if v.Href != nil {
				links["href"] = v.Href
			}
			if v.Rel != nil {
				links["rel"] = v.Rel
			}

			apiLinkList[k] = links
		}
		return apiLinkList
	}
	return nil
}

func flattenSharingStatus(sharingStatus *volumesClient.SharingStatus) string {
	var sharingStatusStr string
	if sharingStatus != nil {
		const two, three = 2, 3
		if *sharingStatus == volumesClient.SharingStatus(two) {
			sharingStatusStr = "SHARED"
		}
		if *sharingStatus == volumesClient.SharingStatus(three) {
			sharingStatusStr = "NOT_SHARED"
		}
	}
	return sharingStatusStr
}

func flattenEnabledAuthentications(authenticationType *volumesClient.AuthenticationType) string {
	var enabledAuthentications string
	if authenticationType != nil {
		const two, three = 2, 3
		if *authenticationType == volumesClient.AuthenticationType(two) {
			enabledAuthentications = "CHAP"
		}
		if *authenticationType == volumesClient.AuthenticationType(three) {
			enabledAuthentications = "NONE"
		}
	}
	return enabledAuthentications
}

func flattenIscsiFeatures(iscsiFeatures *volumesClient.IscsiFeatures) []map[string]interface{} {
	if iscsiFeatures != nil {
		enabledAuthentications := make(map[string]interface{})
		enabledAuthentications["enabled_authentications"] = flattenEnabledAuthentications(iscsiFeatures.EnabledAuthentications)
		return []map[string]interface{}{enabledAuthentications}
	}
	return nil
}

func flattenFlashMode(flashMode *volumesClient.FlashMode) []map[string]interface{} {
	if flashMode != nil {
		flashModeList := make([]map[string]interface{}, 0)
		isEnabled := make(map[string]interface{})
		isEnabled["is_enabled"] = flashMode.IsEnabled
		flashModeList = append(flashModeList, isEnabled)
		return flashModeList
	}
	return nil
}

func flattenStorageFeatures(storageFeatures *volumesClient.StorageFeatures) []map[string]interface{} {
	if storageFeatures != nil {
		storageFeaturesList := make([]map[string]interface{}, 0)
		flashMode := make(map[string]interface{})
		flashMode["flash_mode"] = flattenFlashMode(storageFeatures.FlashMode)
		storageFeaturesList = append(storageFeaturesList, flashMode)
		return storageFeaturesList
	}
	return nil
}

func flattenUsageType(usageType *volumesClient.UsageType) string {
	var usageTypeStr string
	if usageType != nil {
		const two, three, four, five = 2, 3, 4, 5
		if *usageType == volumesClient.UsageType(two) {
			usageTypeStr = "USER"
		}
		if *usageType == volumesClient.UsageType(three) {
			usageTypeStr = "INTERNAL"
		}
		if *usageType == volumesClient.UsageType(four) {
			usageTypeStr = "TEMPORARY"
		}
		if *usageType == volumesClient.UsageType(five) {
			usageTypeStr = "BACKUP_TARGET"
		}
	}
	return usageTypeStr
}
