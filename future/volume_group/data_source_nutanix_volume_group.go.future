package nutanix

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixVolumeGroup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNutanixVolumeGroupRead,
		Schema: getDataSourceVolumeGroupSchema(),
	}
}

func dataSourceNutanixVolumeGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Network Security Rule: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	volumeGroupID, ok := d.GetOk("volume_group_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute volume_group_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetVolumeGroup(volumeGroupID.(string))

	if err != nil {
		return err
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", getReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", getReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("api_version", utils.StringValue(resp.APIVersion)); err != nil {
		return err
	}
	if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
		return err
	}
	if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
		return err
	}

	// set state value
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}
	// set flash_mode
	if err := d.Set("flash_mode", utils.StringValue(resp.Status.Resources.FlashMode)); err != nil {
		return err
	}

	if err := d.Set("file_system_type", utils.StringValue(resp.Status.Resources.FileSystemType)); err != nil {
		return err
	}

	if err := d.Set("sharing_status", utils.StringValue(resp.Status.Resources.SharingStatus)); err != nil {
		return err
	}

	// set attachment value
	al := resp.Status.Resources.AttachmentList
	attachList := make([]map[string]interface{}, 0)
	if al != nil {
		attachList = make([]map[string]interface{}, len(al))
		for k, v := range al {
			attach := make(map[string]interface{})
			attach["vm_reference"] = getClusterReferenceValues(v.VMReference)
			attach["iscsi_initiator_name"] = utils.StringValue(v.IscsiInitiatorName)
			attachList[k] = attach
		}
	}
	if err := d.Set("attachment_list", attachList); err != nil {
		return err
	}

	// set disk_list value
	dl := resp.Status.Resources.DiskList
	diskList := make([]map[string]interface{}, 0)
	if dl != nil {
		diskList = make([]map[string]interface{}, len(dl))
		for k, v := range dl {
			vgDisk := make(map[string]interface{})

			// simple first
			vgDisk["vmdisk_uuid"] = utils.StringValue(v.VmdiskUUID)
			vgDisk["index"] = utils.Int64Value(v.Index)
			vgDisk["disk_size_mib"] = utils.Int64Value(v.DiskSizeMib)
			vgDisk["storage_container_uuid"] = utils.StringValue(v.StorageContainerUUID)
			vgDisk["vm_reference"] = getClusterReferenceValues(v.DataSourceReference)
			diskList[k] = vgDisk
		}
	}
	if err := d.Set("disk_list", diskList); err != nil {
		return err
	}

	// set iscsi_target_prefix value
	if err := d.Set("iscsi_target_prefix", resp.Status.Resources.IscsiTargetPrefix); err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func getDataSourceVolumeGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"volume_group_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"metadata": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_hash": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"categories": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"project_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"owner_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"api_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"flash_mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"file_system_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sharing_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"attachment_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"vm_reference": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Required: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"iscsi_initiator_name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"disk_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"vmdisk_uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"index": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"data_source_reference": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Required: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"disk_size_mib": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"storage_container_uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"iscsi_target_prefix": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
