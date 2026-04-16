package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Get the details of a Volume Disk.
func DatasourceNutanixVolumeDiskV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the Volume Disk identified by {extId} in the Volume Group identified by {volumeGroupExtID}.",
		ReadContext: DatasourceNutanixVolumeDiskV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "The external identifier of the Volume Disk.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"volume_group_ext_id": {
				Description: "The external identifier of the Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
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
							Description: "Once configured, this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.",
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
	}
}

func DatasourceNutanixVolumeDiskV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id")
	volumeDiskExtID := d.Get("ext_id")

	resp, err := conn.VolumeAPIInstance.GetVolumeDiskById(utils.StringPtr(volumeGroupExtID.(string)), utils.StringPtr(volumeDiskExtID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching volume Disk : %v", err)
	}
	getResp := resp.Data.GetValue().(volumesClient.VolumeDisk)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("index", getResp.Index); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_size_bytes", getResp.DiskSizeBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_container_id", getResp.StorageContainerId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_data_source_reference", flattenDiskDataSourceReference(getResp.DiskDataSourceReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_storage_features", flattenDiskStorageFeatures(getResp.DiskStorageFeatures)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
