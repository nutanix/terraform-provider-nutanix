package nutanix

import (
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

	metadata := &v3.DSMetadata{}

	// Get the metadata request
	metadata, err := readListMetadata(d, "volume_group")
	if err != nil {
		return err
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
		m, c := setRSEntityMetadata(v.Metadata)

		entity["metadata"] = m
		entity["project_reference"] = getReferenceValues(v.Metadata.ProjectReference)
		entity["owner_reference"] = getReferenceValues(v.Metadata.OwnerReference)
		entity["categories"] = c
		entity["api_version"] = utils.StringValue(v.APIVersion)
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
				attach["vm_reference"] = getClusterReferenceValues(v.VMReference)
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
				vgDisk["vm_reference"] = getClusterReferenceValues(v.DataSourceReference)
				diskList[k] = vgDisk
			}

		}

		entity["attachment_list"] = attachList
		entity["disk_list"] = diskList
		entity["iscsi_target_prefix"] = v.Status.Resources.IscsiTargetPrefix

		entities[k] = entity
	}

	d.SetId(resource.UniqueId())

	return d.Set("entities", entities)
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
