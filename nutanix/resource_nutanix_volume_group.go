package nutanix

import (
	"fmt"
	"log"
	"strconv"
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
	// Get client connection
	conn := meta.(*Client).API

	// Prepare request
	request := &v3.VolumeGroupInput{}
	spec := &v3.VolumeGroup{}
	metadata := &v3.Metadata{}
	res := &v3.VolumeGroupResources{}

	// Read Arguments and set request values
	n, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")

	if v, ok := d.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}
	if !nok {
		return fmt.Errorf("Please provide the required name attribute")
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

	// Make request to the API
	resp, err := conn.V3.CreateVolumeGroup(request)

	if err != nil {
		return err
	}

	uuid := *resp.Metadata.UUID

	// Set terraform state id
	d.SetId(uuid)

	// Wait for the VM to be available
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
			"Error waiting for volume_group (%s) to create: %s", d.Id(), err)
	}

	return resourceNutanixVolumeGroupRead(d, meta)
}

func resourceNutanixVolumeGroupRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetVolumeGroup(d.Id())
	if err != nil {
		return err
	}

	log.Printf("Reading Volume Group values %s", d.Id())
	fmt.Printf("Reading Volume Group values %s", d.Id())

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = resp.Metadata.LastUpdateTime.String()
	metadata["kind"] = utils.StringValue(resp.Metadata.Kind)
	metadata["uuid"] = utils.StringValue(resp.Metadata.UUID)
	metadata["creation_time"] = resp.Metadata.CreationTime.String()
	metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(resp.Metadata.SpecVersion)))
	metadata["spec_hash"] = utils.StringValue(resp.Metadata.SpecHash)
	metadata["name"] = utils.StringValue(resp.Metadata.Name)
	if err := d.Set("metadata", metadata); err != nil {
		return err
	}

	if resp.Metadata.Categories != nil {
		categories := resp.Metadata.Categories
		var catList []map[string]interface{}

		for name, values := range categories {
			catItem := make(map[string]interface{})
			catItem["name"] = name
			catItem["value"] = values
			catList = append(catList, catItem)
		}
		if err := d.Set("categories", catList); err != nil {
			return err
		}
	}

	pr := make(map[string]interface{})
	if resp.Metadata.ProjectReference != nil {
		pr["kind"] = utils.StringValue(resp.Metadata.ProjectReference.Kind)
		pr["name"] = utils.StringValue(resp.Metadata.ProjectReference.Name)
		pr["uuid"] = utils.StringValue(resp.Metadata.ProjectReference.UUID)

	}
	if err := d.Set("project_reference", pr); err != nil {
		return err
	}
	or := make(map[string]interface{})
	if resp.Metadata.OwnerReference != nil {
		or["kind"] = utils.StringValue(resp.Metadata.OwnerReference.Kind)
		or["name"] = utils.StringValue(resp.Metadata.OwnerReference.Name)
		or["uuid"] = utils.StringValue(resp.Metadata.OwnerReference.UUID)

	}
	if err := d.Set("owner_reference", or); err != nil {
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

			// set vm_reference value
			vmRef := make(map[string]interface{})
			if v.VMReference != nil {
				vmRef["kind"] = utils.StringValue(v.VMReference.Kind)
				vmRef["uuid"] = utils.StringValue(v.VMReference.UUID)
			}
			attach["vm_reference"] = vmRef

			// set iscsi_initiator_name
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

			// set vm_reference value
			dsRef := make(map[string]interface{})
			if v.DataSourceReference != nil {
				dsRef["kind"] = utils.StringValue(v.DataSourceReference.Kind)
				dsRef["uuid"] = utils.StringValue(v.DataSourceReference.UUID)
			}
			vgDisk["vm_reference"] = dsRef

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

func resourceNutanixVolumeGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	log.Printf("Updating Volume Group values %s", d.Id())
	fmt.Printf("Updating Volume Group values %s", d.Id())

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

	// get state
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
		r := &v3.Reference{
			Kind: utils.String(or["kind"].(string)),
			UUID: utils.String(or["uuid"].(string)),
			Name: utils.String(or["name"].(string)),
		}
		metadata.OwnerReference = r
	}
	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pr["kind"].(string)),
			UUID: utils.String(pr["uuid"].(string)),
			Name: utils.String(pr["name"].(string)),
		}
		metadata.ProjectReference = r
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
						v := value.(map[string]interface{})
						ref := &v3.Reference{}
						if j, ok1 := v["kind"]; ok1 {
							ref.Kind = utils.String(j.(string))
						}
						if j, ok1 := v["uuid"]; ok1 {
							ref.UUID = utils.String(j.(string))
						}
						attachment.VMReference = ref
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
			"Error waiting for volume group (%s) to update: %s", d.Id(), err)
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
		return fmt.Errorf(
			"Error waiting for volume group (%s) to delete: %s", d.Id(), err)
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
					v := value.(map[string]interface{})
					ref := &v3.Reference{}
					if j, ok1 := v["kind"]; ok1 {
						ref.Kind = utils.String(j.(string))
					}
					if j, ok1 := v["uuid"]; ok1 {
						ref.UUID = utils.String(j.(string))
					}
					attachment.VMReference = ref
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
				d := &v3.VGDisk{}

				if value, ok := val["vmdisk_uuid"]; ok && value.(string) != "" {
					d.VmdiskUUID = utils.String(value.(string))
				}

				if value, ok := val["index"]; ok && value.(int) >= 0 {
					d.Index = utils.Int64(int64(value.(int)))
				}

				if value, ok := val["data_source_reference"]; ok && len(value.(map[string]interface{})) != 0 {
					v := value.(map[string]interface{})
					ref := &v3.Reference{}
					if j, ok1 := v["kind"]; ok1 {
						ref.Kind = utils.String(j.(string))
					}
					if j, ok1 := v["uuid"]; ok1 {
						ref.UUID = utils.String(j.(string))
					}
					d.DataSourceReference = ref
				}

				if value, ok := val["disk_size_mib"]; ok && value.(int) >= 0 {
					d.DiskSizeMib = utils.Int64(int64(value.(int)))
				}

				if value, ok := val["storage_container_uuid"]; ok && value.(string) != "" {
					d.StorageContainerUUID = utils.String(value.(string))
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
