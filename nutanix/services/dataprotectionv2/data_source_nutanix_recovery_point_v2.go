package dataprotectionv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/common"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

const (
	ApplicationConsistentPropertiesVss1 = "dataprotection.v4.common.VssProperties"
	ApplicationConsistentPropertiesVss2 = "dataprotection.v4.r0.b1.common.VssProperties"
)

func DatasourceNutanixRecoveryPointV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRecoveryPointV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": SchemaForLinks(),
			"location_agnostic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recovery_point_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_references": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vm_recovery_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": SchemaForLinks(),
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recovery_point_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consistency_group_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_agnostic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_recovery_points": SchemaForDiskRecoveryPoints(),
						"vm_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"application_consistent_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backup_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"should_include_writers": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"writers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"should_store_vss_metadata": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"object_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"volume_group_recovery_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": SchemaForLinks(),
						"consistency_group_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_agnostic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_group_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_group_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"disk_recovery_points": SchemaForDiskRecoveryPoints(),
					},
				},
			},
		},
	}
}

func SchemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func SchemaForDiskRecoveryPoints() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"disk_recovery_point_ext_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"disk_ext_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func DatasourceNutanixRecoveryPointV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DatasourceNutanixRecoveryPointV2Read \n")

	conn := meta.(*conns.Client).DataProtectionAPI

	recoveryPointExtID := d.Get("ext_id").(string)

	resp, err := conn.RecoveryPoint.GetRecoveryPointById(&recoveryPointExtID)
	if err != nil {
		return diag.Errorf("error while fetching recovery point: %v", err)
	}

	getResp := resp.Data.GetValue().(config.RecoveryPoint)

	aJSON, _ := json.Marshal(getResp)
	log.Printf("[DEBUG] DatasourceNutanixRecoveryPointV2Read response: \n%v\n", string(aJSON))

	log.Printf("[DEBUG] RecoveryPoint.Name: %v\n", getResp.Name)
	log.Printf("[DEBUG] RecoveryPoint.ExtId: %v\n", getResp.ExtId)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location_agnostic_id", getResp.LocationAgnosticId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_time", flattenTime(getResp.CreationTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("expiration_time", flattenTime(getResp.ExpirationTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", flattenStatus(getResp.Status)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("recovery_point_type", flattenRecoveryPointType(getResp.RecoveryPointType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location_references", flattenLocationReferences(getResp.LocationReferences)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_recovery_points", flattenVMRecoveryPoints(getResp.VmRecoveryPoints)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("volume_group_recovery_points", flattenVolumeGroupRecoveryPoints(getResp.VolumeGroupRecoveryPoints)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}

func flattenTime(inTime *time.Time) string {
	if inTime != nil {
		return inTime.UTC().Format(time.RFC3339)
	}
	return ""
}

func flattenVolumeGroupRecoveryPoints(volumeGroupRecoveryPoints []config.VolumeGroupRecoveryPoint) []map[string]interface{} {
	if len(volumeGroupRecoveryPoints) > 0 {
		volumeGroupRecoveryPointList := make([]map[string]interface{}, len(volumeGroupRecoveryPoints))

		for k, v := range volumeGroupRecoveryPoints {
			volumeGroupRecoveryPoint := map[string]interface{}{}
			if v.TenantId != nil {
				volumeGroupRecoveryPoint["tenant_id"] = v.TenantId
			}
			if v.ExtId != nil {
				volumeGroupRecoveryPoint["ext_id"] = v.ExtId
			}
			if v.Links != nil {
				volumeGroupRecoveryPoint["links"] = flattenLinks(v.Links)
			}
			if v.ConsistencyGroupExtId != nil {
				volumeGroupRecoveryPoint["consistency_group_ext_id"] = v.ConsistencyGroupExtId
			}
			if v.LocationAgnosticId != nil {
				volumeGroupRecoveryPoint["location_agnostic_id"] = v.LocationAgnosticId
			}
			if v.VolumeGroupExtId != nil {
				volumeGroupRecoveryPoint["volume_group_ext_id"] = v.VolumeGroupExtId
			}
			if v.VolumeGroupCategories != nil {
				volumeGroupRecoveryPoint["volume_group_categories"] = v.VolumeGroupCategories
			}
			if v.DiskRecoveryPoints != nil {
				volumeGroupRecoveryPoint["disk_recovery_points"] = flattenDiskRecoveryPoints(v.DiskRecoveryPoints)
			}

			volumeGroupRecoveryPointList[k] = volumeGroupRecoveryPoint
		}
		return volumeGroupRecoveryPointList
	}
	return nil
}

func flattenVMRecoveryPoints(vmRecoveryPoints []config.VmRecoveryPoint) []map[string]interface{} {
	if len(vmRecoveryPoints) > 0 {
		vmRecoveryPointList := make([]map[string]interface{}, len(vmRecoveryPoints))

		for k, v := range vmRecoveryPoints {
			vmRecoveryPoint := map[string]interface{}{}
			if v.TenantId != nil {
				vmRecoveryPoint["tenant_id"] = v.TenantId
			}
			if v.ExtId != nil {
				vmRecoveryPoint["ext_id"] = v.ExtId
			}
			if v.Links != nil {
				vmRecoveryPoint["links"] = flattenLinks(v.Links)
			}
			if v.ConsistencyGroupExtId != nil {
				vmRecoveryPoint["consistency_group_ext_id"] = v.ConsistencyGroupExtId
			}
			if v.LocationAgnosticId != nil {
				vmRecoveryPoint["location_agnostic_id"] = v.LocationAgnosticId
			}
			if v.DiskRecoveryPoints != nil {
				vmRecoveryPoint["disk_recovery_points"] = flattenDiskRecoveryPoints(v.DiskRecoveryPoints)
			}
			if v.VmExtId != nil {
				vmRecoveryPoint["vm_ext_id"] = v.VmExtId
			}
			if v.VmCategories != nil {
				vmRecoveryPoint["vm_categories"] = v.VmCategories
			}
			if v.Name != nil {
				vmRecoveryPoint["name"] = v.Name
			}
			if v.ExpirationTime != nil {
				vmRecoveryPoint["expiration_time"] = flattenTime(v.ExpirationTime)
			}
			if v.Status != nil {
				vmRecoveryPoint["status"] = flattenStatus(v.Status)
			}
			if v.RecoveryPointType != nil {
				vmRecoveryPoint["recovery_point_type"] = flattenRecoveryPointType(v.RecoveryPointType)
			}
			if v.CreationTime != nil {
				vmRecoveryPoint["creation_time"] = flattenTime(v.CreationTime)
			}
			if v.ApplicationConsistentProperties != nil {
				vmRecoveryPoint["application_consistent_properties"] = flattenApplicationConsistentProperties(v.ApplicationConsistentProperties)
			}

			aJSON, _ := json.MarshalIndent(v, "", "  ")
			log.Printf("[DEBUG] VM Recovery Point v: %v\n", string(aJSON))

			vmRecoveryPointList[k] = vmRecoveryPoint
		}

		aJSON, _ := json.MarshalIndent(vmRecoveryPointList, "", "  ")
		log.Printf("[DEBUG] VM Recovery Points Flattened: %v\n", string(aJSON))
		return vmRecoveryPointList
	}
	return nil
}

func flattenApplicationConsistentProperties(vmRecoveryProperties *config.OneOfVmRecoveryPointApplicationConsistentProperties) []map[string]interface{} {
	if vmRecoveryProperties != nil {
		vmRecProps := make(map[string]interface{})
		if *vmRecoveryProperties.ObjectType_ == ApplicationConsistentPropertiesVss1 ||
			*vmRecoveryProperties.ObjectType_ == ApplicationConsistentPropertiesVss2 {
			properties := vmRecoveryProperties.GetValue().(common.VssProperties)
			vmRecProps["backup_type"] = flattenBackupType(properties.BackupType)
			vmRecProps["should_include_writers"] = properties.ShouldIncludeWriters
			vmRecProps["writers"] = properties.Writers
			vmRecProps["should_store_vss_metadata"] = properties.ShouldStoreVssMetadata
			vmRecProps["object_type"] = properties.ObjectType_
		}
		return []map[string]interface{}{vmRecProps}
	}
	return nil
}

func flattenBackupType(backupType *common.BackupType) string {
	if backupType != nil {
		const two, three = 2, 3
		if *backupType == common.BackupType(two) {
			return "FULL_BACKUP"
		}
		if *backupType == common.BackupType(three) {
			return "COPY_BACKUP"
		}
	}
	return "UNKNOWN"
}

func flattenDiskRecoveryPoints(diskRecoveryPoints []common.DiskRecoveryPoint) []map[string]interface{} {
	if len(diskRecoveryPoints) > 0 {
		diskRecoveryPointList := make([]map[string]interface{}, len(diskRecoveryPoints))

		for k, v := range diskRecoveryPoints {
			diskRecoveryPoint := map[string]interface{}{}
			if v.DiskRecoveryPointExtId != nil {
				diskRecoveryPoint["disk_recovery_point_ext_id"] = v.DiskRecoveryPointExtId
			}
			if v.DiskExtId != nil {
				diskRecoveryPoint["disk_ext_id"] = v.DiskExtId
			}

			diskRecoveryPointList[k] = diskRecoveryPoint
		}
		return diskRecoveryPointList
	}
	return nil
}

func flattenLocationReferences(references []config.LocationReference) []map[string]interface{} {
	if len(references) > 0 {
		locationReferences := make([]map[string]interface{}, len(references))

		for k, v := range references {
			reference := map[string]interface{}{}
			if v.LocationExtId != nil {
				reference["location_ext_id"] = v.LocationExtId
			}

			locationReferences[k] = reference
		}
		return locationReferences
	}
	return nil
}

func flattenRecoveryPointType(recoveryPointType *common.RecoveryPointType) string {
	if recoveryPointType != nil {
		const two, three = 2, 3
		if *recoveryPointType == common.RecoveryPointType(two) {
			return "CRASH_CONSISTENT"
		}
		if *recoveryPointType == common.RecoveryPointType(three) {
			return "APPLICATION_CONSISTENT"
		}
	}
	return "UNKNOWN"
}

func flattenStatus(status *common.RecoveryPointStatus) string {
	if status != nil {
		const two = 2
		if *status == common.RecoveryPointStatus(two) {
			return "COMPLETE"
		}
	}
	return "UNKNOWN"
}

func flattenLinks(apiLinks []response.ApiLink) []map[string]interface{} {
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
