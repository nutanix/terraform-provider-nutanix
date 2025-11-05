package volumesv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// List all the Volume Disks attached to the Volume Group.
func DatasourceNutanixVolumeDisksV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the list of disks corresponding to a Volume Group identified by {volumeGroupExtId}.",
		ReadContext: DatasourceNutanixVolumeDisksV2Read,
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
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
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields: diskSizeBytes",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"select": {
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields: extId, storageContainerId",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"disks": {
				Description: "List of disks corresponding to a Volume Group identified by {volumeGroupExtId}.",
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
						"index": {
							Description: "Index of the disk in a Volume Group. This field is optional and immutable.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"disk_size_bytes": {
							Description: "Size of the disk in bytes. This field is mandatory during Volume Group creation if a new disk is being created on the storage container.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"storage_container_id": {
							Description: "Storage container on which the disk must be created. This is a read-only field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Volume Disk description. This is an optional field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"disk_data_source_reference": {
							Description: "Disk Data Source Reference.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Description: "The external identifier of the Data Source Reference.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"name": {
										Description: "The name of the Data Source Reference.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"uris": {
										Description: "The uri list of the Data Source Reference.",
										Type:        schema.TypeList,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeList,
										},
									},
									"entity_type": {
										Description: "The Entity Type of the Data Source Reference.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"disk_storage_features": {
							Description: "Storage optimization features which must be enabled on the Volume Disks. This is an optional field. If omitted, the disks will honor the Volume Group specific storage features setting.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flash_mode": {
										Description: "The flash mode of the disk.",
										Type:        schema.TypeList,
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Description: "The flash mode is enabled or not.",
													Type:        schema.TypeBool,
													Computed:    true,
												},
											},
										},
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

func DatasourceNutanixVolumeDisksV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	var filter, orderBy, selects *string
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
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	// get the volume disks response
	resp, err := conn.VolumeAPIInstance.ListVolumeDisksByVolumeGroupId(utils.StringPtr(volumeGroupExtID.(string)), page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching Disks attached to the volume group : %v", err)
	}
	// extract the volume disks data from the response
	if resp.Data == nil {
		if err := d.Set("disks", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of volume disks.",
		}}
	}

	getResp := resp.Data.GetValue().([]volumesClient.VolumeDisk)
	// set the volume groups data in the terraform resource
	if err := d.Set("disks", flattenDisksEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenDisksEntities(volumeDisks []volumesClient.VolumeDisk) []interface{} {
	if len(volumeDisks) > 0 {
		volumeDiskList := make([]interface{}, len(volumeDisks))

		for k, v := range volumeDisks {
			volumeDisk := make(map[string]interface{})

			if v.TenantId != nil {
				volumeDisk["tenant_id"] = v.TenantId
			}
			if v.ExtId != nil {
				volumeDisk["ext_id"] = v.ExtId
			}
			if v.Links != nil {
				volumeDisk["links"] = flattenLinks(v.Links)
			}
			if v.Index != nil {
				volumeDisk["index"] = v.Index
			}
			if v.DiskSizeBytes != nil {
				volumeDisk["disk_size_bytes"] = v.DiskSizeBytes
			}
			if v.StorageContainerId != nil {
				volumeDisk["storage_container_id"] = v.StorageContainerId
			}
			if v.Description != nil {
				volumeDisk["description"] = v.Description
			}
			if v.DiskDataSourceReference != nil {
				volumeDisk["disk_data_source_reference"] = flattenDiskDataSourceReference(v.DiskDataSourceReference)
			}
			if v.DiskStorageFeatures != nil {
				volumeDisk["disk_storage_features"] = flattenDiskStorageFeatures(v.DiskStorageFeatures)
			}
			log.Printf("[DEBUG] Disk : %v", volumeDisk)
			volumeDiskList[k] = volumeDisk
		}
		return volumeDiskList
	}
	return nil
}

func flattenDiskDataSourceReference(entityReference *config.EntityReference) []map[string]interface{} {
	if entityReference != nil {
		diskDataSourceReferenceList := make([]map[string]interface{}, 0)
		diskDataSourceReference := make(map[string]interface{})
		diskDataSourceReference["ext_id"] = entityReference.ExtId
		diskDataSourceReference["name"] = entityReference.Name
		diskDataSourceReference["uris"] = entityReference.Uris
		diskDataSourceReference["entity_type"] = entityReference.EntityType

		log.Printf("[DEBUG] Disks Data Source Reference ext_id: %v", diskDataSourceReference["ext_id"])
		log.Printf("[DEBUG] Disks Data Source Reference name: %v", diskDataSourceReference["name"])
		log.Printf("[DEBUG] Disks Data Source Reference uris: %v", diskDataSourceReference["uris"])
		log.Printf("[DEBUG] Disks Data Source Reference entity_type: %v", diskDataSourceReference["entity_type"])

		diskDataSourceReferenceList = append(diskDataSourceReferenceList, diskDataSourceReference)

		return diskDataSourceReferenceList
	}
	return nil
}

func flattenDiskStorageFeatures(diskStorageFeatures *volumesClient.DiskStorageFeatures) []map[string]interface{} {
	if diskStorageFeatures != nil {
		diskStorageFeaturesList := make([]map[string]interface{}, 0)
		flashMode := make(map[string]interface{})
		flashMode["flash_mode"] = flattenFlashMode(diskStorageFeatures.FlashMode)
		diskStorageFeaturesList = append(diskStorageFeaturesList, flashMode)
		return diskStorageFeaturesList
	}
	return nil
}
