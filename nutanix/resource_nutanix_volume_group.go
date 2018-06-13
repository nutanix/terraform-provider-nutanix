package nutanix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixVolumeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixVolumeGroupCreate,
		Read:   resourceNutanixVolumeGroupRead,
		Update: resourceNutanixVolumeGroupUpdate,
		Delete: resourceNutanixVolumeGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getVGSchema(),
	}
}

func resourceNutanixVolumeGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.VolumeGroupInput{}
	spec := &v3.VolumeGroup{}
	metadata := &v3.Metadata{}
	res := &v3.VolumeGroupResources{}

	n, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")

	if v, ok := d.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}
	if !nok {
		return fmt.Errorf("please provide the required name attribute")
	}
	if err := getMetadataAttributes(d, metadata, "volume_group"); err != nil {
		return err
	}
	if descok {
		spec.Description = utils.String(desc.(string))
	}

	if err := getVolumeGroupResources(d, res); err != nil {
		return err
	}

	spec.Name = utils.String(n.(string))
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	resp, err := conn.V3.CreateVolumeGroup(request)

	if err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    volumeGroupStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for volume_group (%s) to create: %s", d.Id(), err)
	}

	return resourceNutanixVolumeGroupRead(d, meta)
}

func resourceNutanixVolumeGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	resp, err := conn.V3.GetVolumeGroup(d.Id())
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

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))
	d.Set("state", utils.StringValue(resp.Status.State))
	d.Set("flash_mode", utils.StringValue(resp.Status.Resources.FlashMode))
	d.Set("file_system_type", utils.StringValue(resp.Status.Resources.FileSystemType))
	d.Set("sharing_status", utils.StringValue(resp.Status.Resources.SharingStatus))

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

	dl := resp.Status.Resources.DiskList
	diskList := make([]map[string]interface{}, 0)
	if dl != nil {
		diskList = make([]map[string]interface{}, len(dl))
		for k, v := range dl {
			vgDisk := make(map[string]interface{})
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

	d.Set("iscsi_target_prefix", utils.StringValue(resp.Status.Resources.IscsiTargetPrefix))
	d.SetId(*resp.Metadata.UUID)

	return nil
}

func resourceNutanixVolumeGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.VolumeGroupInput{}
	metadata := &v3.Metadata{}
	res := &v3.VolumeGroupResources{}
	spec := &v3.VolumeGroup{}

	response, err := conn.V3.GetVolumeGroup(d.Id())

	if err != nil {
		return err
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			res = response.Spec.Resources
		}
	}

	if d.HasChange("categories") {
		catl := d.Get("categories").([]interface{})

		if len(catl) > 0 {
			cl := make(map[string]string)
			for _, v := range catl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["value"]; kok && k.(string) != "" {
						cl[i.(string)] = k.(string)
					}
				}
			}
			metadata.Categories = cl
		} else {
			metadata.Categories = nil
		}
	}
	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").(map[string]interface{})
		metadata.OwnerReference = validateRef(or)
	}
	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		metadata.ProjectReference = validateRef(pr)
	}
	if d.HasChange("name") {
		spec.Name = utils.String(d.Get("name").(string))
	}
	if d.HasChange("description") {
		spec.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("flash_mode") {
		res.FlashMode = utils.String(d.Get("flash_mode").(string))
	}

	if d.HasChange("file_system_type") {
		res.FileSystemType = utils.String(d.Get("file_system_type").(string))
	}

	if d.HasChange("sharing_status") {
		res.SharingStatus = utils.String(d.Get("sharing_status").(string))
	}

	if d.HasChange("attachment_list") {
		if v, ok := d.GetOk("attachment_list"); ok {
			n := v.([]interface{})
			if len(n) > 0 {
				attachments := make([]*v3.VMAttachment, len(n))

				for k, nc := range n {
					val := nc.(map[string]interface{})
					attachment := &v3.VMAttachment{}

					if value, ok := val["vm_reference"]; ok && len(value.(map[string]interface{})) != 0 {
						attachment.VMReference = validateShortRef(value.(map[string]interface{}))
					}

					if value, ok := val["iscsi_initiator_name"]; ok && value.(string) != "" {
						attachment.IscsiInitiatorName = utils.String(value.(string))
					}
					attachments[k] = attachment
				}
				res.AttachmentList = attachments
			}
		}
	}

	if d.HasChange("iscsi_target_prefix") {
		res.IscsiTargetPrefix = utils.String(d.Get("iscsi_target_prefix").(string))
	}
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	log.Printf("[DEBUG] Updating Volume Group: %s, %s", d.Get("name").(string), d.Id())
	fmt.Printf("[DEBUG] Updating Volume Group: %s, %s", d.Get("name").(string), d.Id())

	_, errUpdate := conn.V3.UpdateVolumeGroup(d.Id(), request)
	if errUpdate != nil {
		return errUpdate
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    volumeGroupStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for volume group (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixVolumeGroupRead(d, meta)
}

func resourceNutanixVolumeGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting Volume Group: %s", d.Get("name").(string))

	conn := meta.(*Client).API
	UUID := d.Id()

	if err := conn.V3.DeleteVolumeGroup(UUID); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "DELETE_IN_PROGRESS", "COMPLETE"},
		Target:     []string{"DELETED"},
		Refresh:    volumeGroupStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for volume group (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getVolumeGroupResources(d *schema.ResourceData, vg *v3.VolumeGroupResources) error {
	if v, ok := d.GetOk("flash_mode"); ok {
		vg.FlashMode = utils.String(v.(string))
	}

	if v, ok := d.GetOk("file_system_type"); ok {
		vg.FileSystemType = utils.String(v.(string))
	}

	if v, ok := d.GetOk("sharing_status"); ok {
		vg.SharingStatus = utils.String(v.(string))
	}

	if v, ok := d.GetOk("attachment_list"); ok {
		n := v.([]interface{})
		if len(n) > 0 {
			attachments := make([]*v3.VMAttachment, len(n))

			for k, nc := range n {
				val := nc.(map[string]interface{})
				attachment := &v3.VMAttachment{}

				if value, ok := val["vm_reference"]; ok && len(value.(map[string]interface{})) != 0 {
					attachment.VMReference = validateShortRef(value.(map[string]interface{}))
				}

				if value, ok := val["iscsi_initiator_name"]; ok && value.(string) != "" {
					attachment.IscsiInitiatorName = utils.String(value.(string))
				}
				attachments[k] = attachment
			}
			vg.AttachmentList = attachments
		}
	}

	if v, ok := d.GetOk("disk_list"); ok {
		n := v.([]interface{})
		if len(n) > 0 {
			dl := make([]*v3.VGDisk, len(n))

			for k, nc := range n {
				val := nc.(map[string]interface{})
				disk := &v3.VGDisk{}

				if value, ok := val["vmdisk_uuid"]; ok && value.(string) != "" {
					disk.VmdiskUUID = utils.String(value.(string))
				}

				if value, ok := val["index"]; ok && value.(int) >= 0 {
					disk.Index = utils.Int64(int64(value.(int)))
				}

				if value, ok := val["data_source_reference"]; ok && len(value.(map[string]interface{})) != 0 {
					disk.DataSourceReference = validateShortRef(value.(map[string]interface{}))
				}

				if value, ok := val["disk_size_mib"]; ok && value.(int) >= 0 {
					disk.DiskSizeMib = utils.Int64(int64(value.(int)))
				}

				if value, ok := val["storage_container_uuid"]; ok && value.(string) != "" {
					disk.StorageContainerUUID = utils.String(value.(string))
				}

				dl[k] = d
			}
			vg.DiskList = dl
		}
	}

	if v, ok := d.GetOk("iscsi_target_prefix"); ok {
		vg.IscsiTargetPrefix = utils.String(v.(string))
	}

	return nil
}

func getVGSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Optional: true,
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
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"flash_mode": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"file_system_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"sharing_status": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"attachment_list": {
			Type:     schema.TypeList,
			Optional: true,
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
			Optional: true,
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
			Optional: true,
			Computed: true,
		},
	}
}

func volumeGroupStateRefreshFunc(client *v3.Client, uuid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.V3.GetVolumeGroup(uuid)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return v, "DELETED", nil
			}
			log.Printf("ERROR %s", err)
			return nil, "", err
		}

		return v, *v.Status.State, nil
	}
}
