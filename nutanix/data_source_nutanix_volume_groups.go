package nutanix

import (
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixVolumeGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixVolumeGroupsRead,

		Schema: getDataSourceVolumeGroupsSchema(),
	}
}

func dataSourceNutanixVolumeGroupsRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	metadata := &v3.ListMetadata{}

	if v, ok := d.GetOk("metadata"); ok {
		m := v.(map[string]interface{})
		metadata.Kind = utils.String("volume_group")
		if mv, mok := m["sort_attribute"]; mok {
			metadata.SortAttribute = utils.String(mv.(string))
		}
		if mv, mok := m["filter"]; mok {
			metadata.Filter = utils.String(mv.(string))
		}
		if mv, mok := m["length"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Length = utils.Int64(int64(i))
		}
		if mv, mok := m["sort_order"]; mok {
			metadata.SortOrder = utils.String(mv.(string))
		}
		if mv, mok := m["offset"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Offset = utils.Int64(int64(i))
		}
	}

	// Make request to the API
	resp, err := conn.V3.ListVolumeGroup(metadata)
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})
		// set metadata values
		metadata := make(map[string]interface{})
		metadata["last_update_time"] = utils.TimeValue(v.Metadata.LastUpdateTime).String()
		metadata["kind"] = utils.StringValue(v.Metadata.Kind)
		metadata["uuid"] = utils.StringValue(v.Metadata.UUID)
		metadata["creation_time"] = utils.TimeValue(v.Metadata.CreationTime).String()
		metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.Metadata.SpecVersion)))
		metadata["spec_hash"] = utils.StringValue(v.Metadata.SpecHash)
		metadata["name"] = utils.StringValue(v.Metadata.Name)
		entity["metadata"] = metadata

		if v.Metadata.Categories != nil {
			categories := v.Metadata.Categories
			var catList []map[string]interface{}

			for name, values := range categories {
				catItem := make(map[string]interface{})
				catItem["name"] = name
				catItem["value"] = values
				catList = append(catList, catItem)
			}
			entity["categories"] = catList
		}

		entity["api_version"] = utils.StringValue(v.APIVersion)

		pr := make(map[string]interface{})
		if v.Metadata.ProjectReference != nil {
			pr["kind"] = utils.StringValue(v.Metadata.ProjectReference.Kind)
			pr["name"] = utils.StringValue(v.Metadata.ProjectReference.Name)
			pr["uuid"] = utils.StringValue(v.Metadata.ProjectReference.UUID)
		}
		entity["project_reference"] = pr

		or := make(map[string]interface{})
		if v.Metadata.OwnerReference != nil {
			or["kind"] = utils.StringValue(v.Metadata.OwnerReference.Kind)
			or["name"] = utils.StringValue(v.Metadata.OwnerReference.Name)
			or["uuid"] = utils.StringValue(v.Metadata.OwnerReference.UUID)
		}
		entity["owner_reference"] = or

		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)
		entity["state"] = utils.StringValue(v.Status.State)
		entity["flash_mode"] = utils.StringValue(v.Status.Resources.FlashMode)
		entity["file_system_type"] = utils.StringValue(v.Status.Resources.FileSystemType)
		entity["sharing_status"] = utils.StringValue(v.Status.Resources.SharingStatus)

		// set attachment value
		al := v.Status.Resources.AttachmentList
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

		// set disk_list value
		dl := v.Status.Resources.DiskList
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

		entity["attachment_list"] = attachList
		entity["disk_list"] = diskList
		entity["iscsi_target_prefix"] = v.Status.Resources.IscsiTargetPrefix

		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}

func getDataSourceVolumeGroupsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"sort_attribute": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"filter": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"length": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"sort_order": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"offset": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"api_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"entities": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}
