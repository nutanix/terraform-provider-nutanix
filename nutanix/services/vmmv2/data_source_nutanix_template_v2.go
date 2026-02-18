package vmmv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/iam/v4/authn"
	import6 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	import5 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixTemplateV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixTemplateV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"template_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_version_spec": {
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
						"links": schemaForLinks(),
						"version_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_spec": schemaForTemplateVMSpec(),
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForTemplateUser(),
						},
						"is_active_version": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_gc_override_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"guest_update_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployed_vm_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForTemplateUser(),
			},
			"updated_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForTemplateUser(),
			},
			"category_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DatasourceNutanixTemplateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")

	resp, err := conn.TemplatesAPIInstance.GetTemplateById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching template : %v", err)
	}

	getResp := resp.Data.GetValue().(import5.Template)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenAPILink(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("template_name", getResp.TemplateName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("template_description", getResp.TemplateDescription); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("template_version_spec", flattenTemplateVersionSpec(getResp.TemplateVersionSpec)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guest_update_status", flattenGuestUpdateStatus(getResp.GuestUpdateStatus)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreateTime != nil {
		t := getResp.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.UpdateTime != nil {
		t := getResp.UpdateTime
		if err := d.Set("update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("created_by", flattenTemplateUser(getResp.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_by", flattenTemplateUser(getResp.UpdatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_ext_ids", getResp.CategoryExtIds); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenGuestUpdateStatus(pr *import5.GuestUpdateStatus) []map[string]interface{} {
	if pr != nil {
		gStatus := make([]map[string]interface{}, 0)
		g := make(map[string]interface{})

		g["deployed_vm_reference"] = pr.DeployedVmReference

		gStatus = append(gStatus, g)
		return gStatus
	}
	return nil
}

func flattenTemplateUser(pr *import5.TemplateUser) []map[string]interface{} {
	if pr != nil {
		tmps := make([]map[string]interface{}, 0)

		tmp := make(map[string]interface{})

		tmp["ext_id"] = pr.ExtId
		tmp["username"] = pr.Username
		if pr.UserType != nil {
			tmp["user_type"] = flattenUserType(pr.UserType)
		}
		if pr.IdpId != nil {
			tmp["idp_id"] = pr.IdpId
		}
		if pr.DisplayName != nil {
			tmp["display_name"] = pr.DisplayName
		}
		if pr.FirstName != nil {
			tmp["first_name"] = pr.FirstName
		}
		if pr.MiddleInitial != nil {
			tmp["middle_initial"] = pr.MiddleInitial
		}
		if pr.LastName != nil {
			tmp["last_name"] = pr.LastName
		}
		if pr.EmailId != nil {
			tmp["email_id"] = pr.EmailId
		}
		if pr.Locale != nil {
			tmp["locale"] = pr.Locale
		}
		if pr.Region != nil {
			tmp["region"] = pr.Region
		}
		if pr.IsForceResetPasswordEnabled != nil {
			tmp["is_force_reset_password_enabled"] = pr.IsForceResetPasswordEnabled
		}
		if pr.AdditionalAttributes != nil {
			tmp["additional_attributes"] = flattenCustomKVPair(pr.AdditionalAttributes)
		}
		if pr.Status != nil {
			tmp["status"] = flattenUserStatusType(pr.Status)
		}
		if pr.BucketsAccessKeys != nil {
			tmp["buckets_access_keys"] = flattenBucketsAccessKey(pr.BucketsAccessKeys)
		}
		if pr.LastLoginTime != nil {
			tmp["last_login_time"] = pr.LastLoginTime
		}
		if pr.CreatedTime != nil {
			t := pr.CreatedTime
			tmp["created_time"] = t.String()
		}
		if pr.LastUpdatedTime != nil {
			t := pr.LastUpdatedTime
			tmp["last_updated_time"] = t.String()
		}
		if pr.CreatedBy != nil {
			tmp["created_by"] = pr.CreatedBy
		}
		if pr.LastUpdatedBy != nil {
			tmp["updated_by"] = pr.LastUpdatedTime
		}

		tmps = append(tmps, tmp)
		return tmps
	}
	return nil
}

func flattenUserType(pr *authn.UserType) string {
	const two, three, four, five = 2, 3, 4, 5
	if pr != nil {
		if *pr == authn.UserType(two) {
			return "LOCAL"
		}
		if *pr == authn.UserType(three) {
			return "SAML"
		}
		if *pr == authn.UserType(four) {
			return "LDAP"
		}
		if *pr == authn.UserType(five) {
			return "EXTERNAL"
		}
	}
	return "UNKNOWN"
}

func flattenCustomKVPair(kvPairs []config.KVPair) []interface{} {
	if len(kvPairs) > 0 {
		kvps := make([]interface{}, len(kvPairs))

		for k, v := range kvPairs {
			kvp := make(map[string]interface{})

			if v.Name != nil {
				kvp["name"] = v.Name
			}
			if v.Value != nil {
				kvp["value"] = flattenKVValue(v.Value.GetValue())
			}
			kvps[k] = kvp
		}

		return kvps
	}
	return nil
}

func flattenKVValue(value interface{}) []interface{} {
	valueMap := make(map[string]interface{})
	switch v := value.(type) {
	case string:
		valueMap["string"] = v
	case int:
		valueMap["integer"] = v
	case bool:
		valueMap["boolean"] = v
	case []string:
		valueMap["string_list"] = v
	case []int:
		valueMap["integer_list"] = v
	case map[string]string:
		valueMap["object"] = v

	case []config.MapOfStringWrapper:
		mapOfStrings := make([]interface{}, len(v))
		for i, m := range v {
			mapOfStrings[i] = m
		}

		valueMap["map_of_strings"] = mapOfStrings
	default:
		log.Printf("[WARN] Unknown type %T", v)
		return nil
	}
	return []interface{}{valueMap}
}

func flattenUserStatusType(pr *authn.UserStatusType) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == authn.UserStatusType(two) {
			return "ACTIVE"
		}
		if *pr == authn.UserStatusType(three) {
			return "INACTIVE"
		}
	}
	return "UNKNOWN"
}

func flattenBucketsAccessKey(pr []authn.BucketsAccessKey) []interface{} {
	if len(pr) > 0 {
		bckts := make([]interface{}, len(pr))

		for k, v := range pr {
			bkt := make(map[string]interface{})

			if v.ExtId != nil {
				bkt["ext_id"] = v.ExtId
			}
			if v.AccessKeyName != nil {
				bkt["access_key_name"] = v.AccessKeyName
			}
			if v.SecretAccessKey != nil {
				bkt["secret_access_key"] = v.SecretAccessKey
			}
			if v.UserId != nil {
				bkt["user_id"] = v.UserId
			}
			if v.CreatedTime != nil {
				t := v.CreatedTime
				bkt["created_time"] = t.String()
			}
			bckts[k] = bkt
		}

		return bckts
	}
	return nil
}

func flattenTemplateVersionSpec(pr *import5.TemplateVersionSpec) []map[string]interface{} {
	if pr != nil {
		tmps := make([]map[string]interface{}, 0)

		tmp := make(map[string]interface{})

		if pr.TenantId != nil {
			tmp["tenant_id"] = pr.TenantId
		}
		if pr.ExtId != nil {
			tmp["ext_id"] = pr.ExtId
		}
		if pr.Links != nil {
			tmp["links"] = flattenAPILink(pr.Links)
		}
		if pr.VersionName != nil {
			tmp["version_name"] = pr.VersionName
		}
		if pr.VersionDescription != nil {
			tmp["version_description"] = pr.VersionDescription
		}
		if pr.VmSpec != nil {
			tmp["vm_spec"] = flattenVM(pr.VmSpec)
		}
		if pr.CreateTime != nil {
			t := pr.CreateTime
			tmp["create_time"] = t.String()
		}
		if pr.CreatedBy != nil {
			tmp["created_by"] = flattenTemplateUser(pr.CreatedBy)
		}
		// if pr.VersionSource != nil {
		//	tmp["version_source"] = flattenTemplateVersionSource(pr.VersionSource)
		//}
		if pr.IsActiveVersion != nil {
			tmp["is_active_version"] = pr.IsActiveVersion
		}
		if pr.IsGcOverrideEnabled != nil {
			tmp["is_gc_override_enabled"] = pr.IsGcOverrideEnabled
		}

		tmps = append(tmps, tmp)
		return tmps
	}
	return nil
}

func flattenTemplateVersionSource(versionSource *import5.OneOfTemplateVersionSpecVersionSource) []map[string]interface{} {
	if versionSource != nil {
		tmps := make([]map[string]interface{}, 0)

		tmp := make(map[string]interface{})

		if *versionSource.ObjectType_ == "vmm.v4.content.TemplateVmReference" {
			vmReferenceMap := make(map[string]interface{})
			vmReference := versionSource.GetValue().(import5.TemplateVmReference)

			vmReferenceMap["ext_id"] = vmReference.ExtId
			vmReferenceMap["guest_customization"] = flattenGuestCustomizationParams(vmReference.GuestCustomization)

			tmp["template_vm_reference"] = []map[string]interface{}{vmReferenceMap}
		}
		if *versionSource.ObjectType_ == "vmm.v4.content.TemplateVersionReference" {
			tempVersionReferenceMap := make(map[string]interface{})
			versionReference := versionSource.GetValue().(import5.TemplateVersionReference)

			tempVersionReferenceMap["version_id"] = versionReference.VersionId
			tempVersionReferenceMap["override_vm_config"] = flattenTemplateVMRefOverrideVMConfig(versionReference.OverrideVmConfig)

			tmp["template_version_reference"] = []map[string]interface{}{tempVersionReferenceMap}
		}

		tmps = append(tmps, tmp)
		return tmps
	}
	return nil
}

// func flattenTemplateGuestCustomization(guestCustomization *import6.GuestCustomizationParams) []map[string]interface{} {
//	if guestCustomization != nil {
//		guestCustomizationMap := make(map[string]interface{})
//		if guestCustomization.Config != nil {
//			guestCustomizationMap["domain"] = flattenGuestCustomizationConfig(guestCustomization.Config)
//		}
//		return []map[string]interface{}{guestCustomizationMap}
//	}
//	return nil
//}

// func flattenGuestCustomizationConfig(customizationParamsConfig *import6.OneOfGuestCustomizationParamsConfig) []map[string]interface{} {
//	if customizationParamsConfig != nil {
//		customizationParamsConfigMap := make(map[string]interface{})
//		if *customizationParamsConfig.ObjectType_ == "vmm.v4.ahv.config.SysprepConfig" {
//			sysprepConfigMap := make(map[string]interface{})
//			sysprepConfig := customizationParamsConfig.GetValue().(import6.Sysprep)
//
//			sysprepConfigMap["sysprep"] = flattenSysprepConfig(sysprepConfig)
//			customizationParamsConfigMap["config"] = []map[string]interface{}{sysprepConfigMap}
//		}
//		if *customizationParamsConfig.ObjectType_ == "vmm.v4.ahv.config.CloudInit" {
//			cloudInitConfigMap := make(map[string]interface{})
//			cloudInitConfig := customizationParamsConfig.GetValue().(import6.CloudInit)
//
//			cloudInitConfigMap["cloud_init"] = flattenCloudInitConfig(&cloudInitConfig)
//			customizationParamsConfigMap["config"] = []map[string]interface{}{cloudInitConfigMap}
//		}
//		return []map[string]interface{}{customizationParamsConfigMap}
//	}
//	return nil
//}
//
//func flattenSysprepConfig(sysprepConfig import6.Sysprep) interface{} {
//
//}

// func flattenCloudInitConfig(cloudInitConfig *import6.CloudInit) []map[string]interface{} {
//	if cloudInitConfig != nil {
//		cloudInitConfigMap := make(map[string]interface{})
//
//		if cloudInitConfig.DatasourceType != nil {
//			datasourceType := cloudInitConfig.DatasourceType
//			const CONFIG_DRIVE_V2 = 2
//			switch *datasourceType {
//			case CONFIG_DRIVE_V2:
//				cloudInitConfigMap["datasource_type"] = "CONFIG_DRIVE_V2"
//				break
//			default:
//				cloudInitConfigMap["datasource_type"] = "UNKNOWN"
//			}
//			if cloudInitConfig.Metadata != nil {
//				cloudInitConfigMap["metadata"] = cloudInitConfig.Metadata
//			}
//			if cloudInitConfig.CloudInitScript != nil {
//				cloudInitScriptMap := make(map[string]interface{})
//
//				if cloudInitConfig.CloudInitScript.GetValue() == "vmm.v4.ahv.config.Userdata" {
//					userDataMap := make(map[string]interface{})
//					userData := cloudInitConfig.CloudInitScript.GetValue().(import6.Userdata)
//
//					userDataMap["value"] = userData.Value
//					cloudInitScriptMap["cloud_init_script"] = []map[string]interface{}{userDataMap}
//				}
//				if cloudInitConfig.CloudInitScript.GetValue() == "vmm.v4.ahv.config.CustomKeyValues" {
//					customKVMap := make(map[string]interface{})
//					kvValues := cloudInitConfig.CloudInitScript.GetValue().(import6.CustomKeyValues)
//
//					customKVMap["custom_key_values"] = flattenCustomKVPair(kvValues.KeyValuePairs)
//
//					cloudInitScriptMap["cloud_init_script"] = []map[string]interface{}{customKVMap}
//				}
//			}
//
//		}
//		return []map[string]interface{}{cloudInitConfigMap}
//	}
//	return nil
//}

func flattenTemplateVMRefOverrideVMConfig(vmConfig *import5.VmConfigOverride) []map[string]interface{} {
	if vmConfig != nil {
		vmConfigMap := make(map[string]interface{})
		if vmConfig.Name != nil {
			vmConfigMap["name"] = vmConfig.Name
		}
		if vmConfig.NumSockets != nil {
			vmConfigMap["num_sockets"] = vmConfig.NumSockets
		}
		if vmConfig.NumCoresPerSocket != nil {
			vmConfigMap["num_cores_per_socket"] = vmConfig.NumCoresPerSocket
		}
		if vmConfig.NumThreadsPerCore != nil {
			vmConfigMap["num_threads_per_core"] = vmConfig.NumThreadsPerCore
		}
		if vmConfig.MemorySizeBytes != nil {
			vmConfigMap["memory_size_bytes"] = vmConfig.MemorySizeBytes
		}
		if vmConfig.Nics != nil {
			vmConfigMap["nics"] = flattenNic(vmConfig.Nics)
		}
		if vmConfig.GuestCustomization != nil {
			vmConfigMap["guest_customization"] = flattenGuestCustomizationParams(vmConfig.GuestCustomization)
		}

		return []map[string]interface{}{vmConfigMap}
	}
	return nil
}

func SchemaForCreateByAndUpdateByUser() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ext_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"username": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"user_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"idp_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"display_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"first_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"middle_initial": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"last_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"email_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"locale": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"region": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"is_force_reset_password_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"additional_attributes": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"value": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"buckets_access_keys": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"access_key_name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"secret_access_key": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"user_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"created_time": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"last_login_time": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"created_time": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"last_updated_time": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"created_by": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     schemaForTemplateUser(),
				},
				"updated_by": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     schemaForTemplateUser(),
				},
			},
		},
	}
}

func flattenVM(v *import6.Vm) []map[string]interface{} {
	if v != nil {
		vmList := make([]map[string]interface{}, 0)
		vm := make(map[string]interface{})

		if v.TenantId != nil {
			vm["tenant_id"] = v.TenantId
		}
		if v.Links != nil {
			vm["links"] = flattenAPILink(v.Links)
		}
		if v.ExtId != nil {
			vm["ext_id"] = v.ExtId
		}
		if v.Name != nil {
			vm["name"] = v.Name
		}
		if v.Description != nil {
			vm["description"] = v.Description
		}
		if v.CreateTime != nil {
			t := v.CreateTime
			vm["create_time"] = t.String()
		}
		if v.UpdateTime != nil {
			t := v.UpdateTime
			vm["update_time"] = t.String()
		}
		if v.Source != nil {
			vm["source"] = flattenVMSourceReference(v.Source)
		}
		if v.NumSockets != nil {
			vm["num_sockets"] = v.NumSockets
		}
		if v.NumCoresPerSocket != nil {
			vm["num_cores_per_socket"] = v.NumCoresPerSocket
		}
		if v.NumThreadsPerCore != nil {
			vm["num_threads_per_core"] = v.NumThreadsPerCore
		}
		if v.NumNumaNodes != nil {
			vm["num_numa_nodes"] = v.NumNumaNodes
		}
		if v.MemorySizeBytes != nil {
			vm["memory_size_bytes"] = v.MemorySizeBytes
		}
		if v.IsVcpuHardPinningEnabled != nil {
			vm["is_vcpu_hard_pinning_enabled"] = v.IsVcpuHardPinningEnabled
		}
		if v.IsCpuPassthroughEnabled != nil {
			vm["is_cpu_passthrough_enabled"] = v.IsCpuPassthroughEnabled
		}
		if v.EnabledCpuFeatures != nil {
			vm["enabled_cpu_features"] = flattenCPUFeature(v.EnabledCpuFeatures)
		}
		if v.IsMemoryOvercommitEnabled != nil {
			vm["is_memory_overcommit_enabled"] = v.IsMemoryOvercommitEnabled
		}
		if v.IsGpuConsoleEnabled != nil {
			vm["is_gpu_console_enabled"] = v.IsGpuConsoleEnabled
		}
		if v.GenerationUuid != nil {
			vm["generation_uuid"] = v.GenerationUuid
		}
		if v.BiosUuid != nil {
			vm["bios_uuid"] = v.BiosUuid
		}
		if v.Categories != nil {
			vm["categories"] = flattenCategoryReference(v.Categories)
		}
		if v.OwnershipInfo != nil {
			vm["ownership_info"] = flattenOwnershipInfo(v.OwnershipInfo)
		}
		if v.Host != nil {
			vm["host"] = flattenHostReference(v.Host)
		}
		if v.Cluster != nil {
			vm["cluster"] = flattenClusterReference(v.Cluster)
		}
		if v.GuestCustomization != nil {
			vm["guest_customization"] = flattenGuestCustomizationParams(v.GuestCustomization)
		}
		if v.GuestTools != nil {
			vm["guest_tools"] = flattenGuestTools(v.GuestTools)
		}
		if v.HardwareClockTimezone != nil {
			vm["hardware_clock_timezone"] = v.HardwareClockTimezone
		}
		if v.IsBrandingEnabled != nil {
			vm["is_branding_enabled"] = v.IsBrandingEnabled
		}
		if v.BootConfig != nil {
			vm["boot_config"] = flattenOneOfVMBootConfig(v.BootConfig)
		}
		if v.IsVgaConsoleEnabled != nil {
			vm["is_vga_console_enabled"] = v.IsVgaConsoleEnabled
		}
		if v.MachineType != nil {
			vm["machine_type"] = flattenMachineType(v.MachineType)
		}
		if v.PowerState != nil {
			vm["power_state"] = flattenPowerState(v.PowerState)
		}
		if v.VtpmConfig != nil {
			vm["vtpm_config"] = flattenVtpmConfig(v.VtpmConfig)
		}
		if v.IsAgentVm != nil {
			vm["is_agent_vm"] = v.IsAgentVm
		}
		if v.ApcConfig != nil {
			vm["apc_config"] = flattenApcConfig(v.ApcConfig)
		}
		if v.IsLiveMigrateCapable != nil {
			vm["is_live_migrate_capable"] = v.IsLiveMigrateCapable
		}
		if v.IsCrossClusterMigrationInProgress != nil {
			vm["is_cross_cluster_migration_in_progress"] = v.IsCrossClusterMigrationInProgress
		}
		if v.StorageConfig != nil {
			vm["storage_config"] = flattenADSFVmStorageConfig(v.StorageConfig)
		}
		if v.Disks != nil {
			vm["disks"] = flattenDisk(v.Disks)
		}
		if v.CdRoms != nil {
			vm["cd_roms"] = flattenCdRom(v.CdRoms)
		}
		if v.Nics != nil {
			vm["nics"] = flattenNic(v.Nics)
		}
		if v.Gpus != nil {
			vm["gpus"] = flattenGpu(v.Gpus)
		}
		if v.SerialPorts != nil {
			vm["serial_ports"] = flattenSerialPort(v.SerialPorts)
		}
		if v.ProtectionType != nil {
			vm["protection_type"] = flattenProtectionType(v.ProtectionType)
		}
		if v.ProtectionPolicyState != nil {
			vm["protection_policy_state"] = flattenProtectionPolicyState(v.ProtectionPolicyState)
		}

		vmList = append(vmList, vm)
		return vmList
	}
	return nil
}
