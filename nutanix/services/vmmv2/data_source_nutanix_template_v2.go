package vmmv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/common/v1/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/iam/v4/authn"
	import6 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/config"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/content"
	import7 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/templates"
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
						"guest_customization_profile": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"version_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"template_vm_reference": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"guest_customization_profile": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ext_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"guest_customization": schemaForTemplateGuestCustomization(),
											},
										},
									},
									"template_version_reference": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"override_vm_config": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"num_sockets": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"num_cores_per_socket": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"num_threads_per_core": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"memory_size_bytes": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"nics":                schemaForNics(),
															"guest_customization": schemaForTemplateGuestCustomization(),
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
			"project_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixTemplateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")
	getTemplateByIdRequest := import7.GetTemplateByIdRequest{
		ExtId: utils.StringPtr(extID.(string)),
	}
	resp, err := conn.TemplatesAPIInstance.GetTemplateById(ctx, &getTemplateByIdRequest)
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
	if err := d.Set("project_ext_id", getResp.ProjectExtId); err != nil {
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
		if pr.VersionSource != nil {
			tmp["version_source"] = flattenTemplateVersionSource(pr.VersionSource)
		}
		if pr.IsActiveVersion != nil {
			tmp["is_active_version"] = pr.IsActiveVersion
		}
		if pr.IsGcOverrideEnabled != nil {
			tmp["is_gc_override_enabled"] = pr.IsGcOverrideEnabled
		}
		if pr.GuestCustomizationProfile != nil {
			tmp["guest_customization_profile"] = flattenVmGcProfileReference(pr.GuestCustomizationProfile)
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
			vmReferenceMap["guest_customization_profile"] = flattenVmGcProfileReference(vmReference.GuestCustomizationProfile)

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

func schemaForStringOrDiscardOverride() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"discard": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func schemaForVmGcProfileConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"profile": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"config_override_spec": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"sysprep_config": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"customization": {
											Type:     schema.TypeList,
											Optional: true,
											MaxItems: 1,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"answer_file": {
														Type:     schema.TypeList,
														Optional: true,
														MaxItems: 1,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"unattend_xml": {
																	Type:     schema.TypeString,
																	Required: true,
																},
															},
														},
													},
													"sysprep_params": {
														Type:     schema.TypeList,
														Optional: true,
														MaxItems: 1,
														Elem: &schema.Resource{
															Schema: schemaForSysprepParamsOverrideSpec(),
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
			},
		},
	}
}

func schemaForSysprepParamsOverrideSpec() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"first_logon_commands": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"general_settings": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"administrator_password": schemaForStringOrDiscardOverride(),
					"auto_logon_settings": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"logon_count": {
									Type:     schema.TypeInt,
									Optional: true,
								},
							},
						},
					},
					"computer_name": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"use_vm_name": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"discard": {
									Type:     schema.TypeBool,
									Optional: true,
								},
							},
						},
					},
					"registered_organization": schemaForStringOrDiscardOverride(),
					"registered_owner":        schemaForStringOrDiscardOverride(),
					"timezone":                schemaForStringOrDiscardOverride(),
					"windows_product_key":     schemaForStringOrDiscardOverride(),
				},
			},
		},
		"locale_settings": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"system_locale": schemaForStringOrDiscardOverride(),
					"ui_language":   schemaForStringOrDiscardOverride(),
					"user_locale":   schemaForStringOrDiscardOverride(),
				},
			},
		},
		"network_settings": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nic_config_list": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dns_config": {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"preferred_dns_server_address": {
												Type:     schema.TypeString,
												Optional: true,
											},
											"alternate_dns_server_addresses": {
												Type:     schema.TypeList,
												Optional: true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"ipv4_config": {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"use_dhcp": {
												Type:     schema.TypeBool,
												Optional: true,
											},
											"ip_address": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"prefix_length": {
															Type:     schema.TypeInt,
															Optional: true,
														},
													},
												},
											},
											"default_gateways": {
												Type:     schema.TypeList,
												Optional: true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
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
		"workgroup_or_domain_info": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"workgroup": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"domain_settings": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"credentials": {
									Type:     schema.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"domain_name": {
												Type:     schema.TypeString,
												Optional: true,
											},
											"password": {
												Type:      schema.TypeString,
												Optional:  true,
												Sensitive: true,
											},
											"username": {
												Type:     schema.TypeString,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"discard": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}
}

func expandVmGcProfileConfigOverride(pr interface{}) *import6.VmGcProfileConfig {
	if pr == nil {
		return nil
	}
	list, ok := pr.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	val := list[0].(map[string]interface{})
	cfg := import6.NewVmGcProfileConfig()
	if profile, ok := val["profile"]; ok {
		cfg.Profile = expandVmGcProfileReference(profile)
	}
	if overrideSpec, ok := val["config_override_spec"]; ok {
		cfg.ConfigOverrideSpec = expandConfigOverrideSpec(overrideSpec)
	}
	return cfg
}

func expandConfigOverrideSpec(pr interface{}) *import6.OneOfVmGcProfileConfigConfigOverrideSpec {
	if pr == nil {
		return nil
	}
	list, ok := pr.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	val := list[0].(map[string]interface{})
	oneOf := import6.NewOneOfVmGcProfileConfigConfigOverrideSpec()

	if sysprepConfig, ok := val["sysprep_config"]; ok {
		sysList, ok := sysprepConfig.([]interface{})
		if ok && len(sysList) > 0 && sysList[0] != nil {
			sysVal := sysList[0].(map[string]interface{})
			spec := import6.NewVmGcProfileSysprepConfigOverrideSpec()

			if custRaw, ok := sysVal["customization"]; ok {
				custList, ok := custRaw.([]interface{})
				if ok && len(custList) > 0 && custList[0] != nil {
					custVal := custList[0].(map[string]interface{})
					cust := import6.NewOneOfVmGcProfileSysprepConfigOverrideSpecCustomization()
					valueSet := false

					if answerFile, ok := custVal["answer_file"]; ok {
						afList, ok := answerFile.([]interface{})
						if ok && len(afList) > 0 && afList[0] != nil {
							afVal := afList[0].(map[string]interface{})
							af := import6.NewVmGcProfileAnswerFileOverrideSpec()
							if xml, ok := afVal["unattend_xml"].(string); ok && xml != "" {
								af.UnattendXml = utils.StringPtr(xml)
							}
							cust.SetValue(*af)
							valueSet = true
						}
					}
					if !valueSet {
						if sysprepParams, ok := custVal["sysprep_params"]; ok {
							spList, ok := sysprepParams.([]interface{})
							if ok && len(spList) > 0 && spList[0] != nil {
								spVal := spList[0].(map[string]interface{})
								params := expandSysprepParamsOverrideSpec(spVal)
								cust.SetValue(*params)
								valueSet = true
							}
						}
					}
					if valueSet {
						spec.Customization = cust
					}
				}
			}
			oneOf.SetValue(*spec)
		}
	}
	return oneOf
}

func expandSysprepParamsOverrideSpec(spVal map[string]interface{}) *import6.VmGcProfileSysprepParamsOverrideSpec {
	params := import6.NewVmGcProfileSysprepParamsOverrideSpec()

	if cmds, ok := spVal["first_logon_commands"].([]interface{}); ok && len(cmds) > 0 {
		commands := make([]string, len(cmds))
		for i, c := range cmds {
			commands[i] = c.(string)
		}
		params.FirstLogonCommands = commands
	}

	if gsList, ok := spVal["general_settings"].([]interface{}); ok && len(gsList) > 0 && gsList[0] != nil {
		gsVal := gsList[0].(map[string]interface{})
		gs := import6.NewVmGcProfileGeneralSettingsOverrideSpec()

		if adminPwd, ok := gsVal["administrator_password"].([]interface{}); ok && len(adminPwd) > 0 && adminPwd[0] != nil {
			pwdVal := adminPwd[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileGeneralSettingsOverrideSpecAdministratorPassword()
			if discard, ok := pwdVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if v, ok := pwdVal["value"].(string); ok && v != "" {
				pwd := import6.NewVmGcProfileAdministratorPassword()
				pwd.Value = utils.StringPtr(v)
				oneOf.SetValue(*pwd)
			}
			gs.AdministratorPassword = oneOf
		}

		if alsList, ok := gsVal["auto_logon_settings"].([]interface{}); ok && len(alsList) > 0 && alsList[0] != nil {
			alsVal := alsList[0].(map[string]interface{})
			als := import6.NewVmGcProfileAutoLogonSettingsOverrideSpec()
			if lc, ok := alsVal["logon_count"].(int); ok {
				als.LogonCount = utils.IntPtr(lc)
			}
			gs.AutoLogonSettings = als
		}

		if cnList, ok := gsVal["computer_name"].([]interface{}); ok && len(cnList) > 0 && cnList[0] != nil {
			cnVal := cnList[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileGeneralSettingsOverrideSpecComputerName()
			if discard, ok := cnVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if useVM, ok := cnVal["use_vm_name"].(bool); ok && useVM {
				oneOf.SetValue(*import6.NewVmGcProfileUseVmNameOverrideSpec())
			} else if v, ok := cnVal["value"].(string); ok && v != "" {
				cn := import6.NewVmGcProfileComputerName()
				cn.Value = utils.StringPtr(v)
				oneOf.SetValue(*cn)
			}
			gs.ComputerName = oneOf
		}

		if roList, ok := gsVal["registered_organization"].([]interface{}); ok && len(roList) > 0 && roList[0] != nil {
			roVal := roList[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileGeneralSettingsOverrideSpecRegisteredOrganization()
			if discard, ok := roVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if v, ok := roVal["value"].(string); ok && v != "" {
				ro := import6.NewVmGcProfileRegisteredOrganization()
				ro.Value = utils.StringPtr(v)
				oneOf.SetValue(*ro)
			}
			gs.RegisteredOrganization = oneOf
		}

		if roList, ok := gsVal["registered_owner"].([]interface{}); ok && len(roList) > 0 && roList[0] != nil {
			roVal := roList[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileGeneralSettingsOverrideSpecRegisteredOwner()
			if discard, ok := roVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if v, ok := roVal["value"].(string); ok && v != "" {
				ro := import6.NewVmGcProfileRegisteredOwner()
				ro.Value = utils.StringPtr(v)
				oneOf.SetValue(*ro)
			}
			gs.RegisteredOwner = oneOf
		}

		if tzList, ok := gsVal["timezone"].([]interface{}); ok && len(tzList) > 0 && tzList[0] != nil {
			tzVal := tzList[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileGeneralSettingsOverrideSpecTimezone()
			if discard, ok := tzVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if v, ok := tzVal["value"].(string); ok && v != "" {
				tz := import6.NewVmGcProfileTimezone()
				tz.Value = utils.StringPtr(v)
				oneOf.SetValue(*tz)
			}
			gs.Timezone = oneOf
		}

		if wpkList, ok := gsVal["windows_product_key"].([]interface{}); ok && len(wpkList) > 0 && wpkList[0] != nil {
			wpkVal := wpkList[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileGeneralSettingsOverrideSpecWindowsProductKey()
			if discard, ok := wpkVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if v, ok := wpkVal["value"].(string); ok && v != "" {
				wpk := import6.NewVmGcProfileWindowsProductKey()
				wpk.Value = utils.StringPtr(v)
				oneOf.SetValue(*wpk)
			}
			gs.WindowsProductKey = oneOf
		}

		params.GeneralSettings = gs
	}

	if lsList, ok := spVal["locale_settings"].([]interface{}); ok && len(lsList) > 0 && lsList[0] != nil {
		lsVal := lsList[0].(map[string]interface{})
		ls := import6.NewVmGcProfileLocaleSettingsOverrideSpec()

		if slList, ok := lsVal["system_locale"].([]interface{}); ok && len(slList) > 0 && slList[0] != nil {
			slVal := slList[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileLocaleSettingsOverrideSpecSystemLocale()
			if discard, ok := slVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if v, ok := slVal["value"].(string); ok && v != "" {
				loc := import6.NewVmGcProfileLocaleSettingOverride()
				loc.Value = utils.StringPtr(v)
				oneOf.SetValue(*loc)
			}
			ls.SystemLocale = oneOf
		}

		if ulList, ok := lsVal["ui_language"].([]interface{}); ok && len(ulList) > 0 && ulList[0] != nil {
			ulVal := ulList[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileLocaleSettingsOverrideSpecUiLanguage()
			if discard, ok := ulVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if v, ok := ulVal["value"].(string); ok && v != "" {
				loc := import6.NewVmGcProfileLocaleSettingOverride()
				loc.Value = utils.StringPtr(v)
				oneOf.SetValue(*loc)
			}
			ls.UiLanguage = oneOf
		}

		if ulcList, ok := lsVal["user_locale"].([]interface{}); ok && len(ulcList) > 0 && ulcList[0] != nil {
			ulcVal := ulcList[0].(map[string]interface{})
			oneOf := import6.NewOneOfVmGcProfileLocaleSettingsOverrideSpecUserLocale()
			if discard, ok := ulcVal["discard"].(bool); ok && discard {
				oneOf.SetValue(*import6.NewVmGcProfileDiscardSettings())
			} else if v, ok := ulcVal["value"].(string); ok && v != "" {
				loc := import6.NewVmGcProfileLocaleSettingOverride()
				loc.Value = utils.StringPtr(v)
				oneOf.SetValue(*loc)
			}
			ls.UserLocale = oneOf
		}

		params.LocaleSettings = ls
	}

	if nsList, ok := spVal["network_settings"].([]interface{}); ok && len(nsList) > 0 && nsList[0] != nil {
		nsVal := nsList[0].(map[string]interface{})
		ns := import6.NewVmGcProfileNetworkSettingsOverrideSpec()
		if nicList, ok := nsVal["nic_config_list"].([]interface{}); ok && len(nicList) > 0 {
			nics := make([]import6.VmGcProfileNicConfigOverrideSpec, len(nicList))
			for i, nicRaw := range nicList {
				nicMap := nicRaw.(map[string]interface{})
				nic := *import6.NewVmGcProfileNicConfigOverrideSpec()
				if dnsList, ok := nicMap["dns_config"].([]interface{}); ok && len(dnsList) > 0 && dnsList[0] != nil {
					dnsMap := dnsList[0].(map[string]interface{})
					dns := import6.NewVmGcProfileDnsConfigOverrideSpec()
					if v, ok := dnsMap["preferred_dns_server_address"].(string); ok && v != "" {
						dns.PreferredDnsServerAddress = utils.StringPtr(v)
					}
					if altDns, ok := dnsMap["alternate_dns_server_addresses"].([]interface{}); ok && len(altDns) > 0 {
						addrs := make([]string, len(altDns))
						for j, a := range altDns {
							addrs[j] = a.(string)
						}
						dns.AlternateDnsServerAddresses = addrs
					}
					nic.DnsConfig = dns
				}
				if ipv4List, ok := nicMap["ipv4_config"].([]interface{}); ok && len(ipv4List) > 0 && ipv4List[0] != nil {
					ipv4Map := ipv4List[0].(map[string]interface{})
					oneOfIPv4 := import6.NewOneOfVmGcProfileNicConfigOverrideSpecIpv4Config()
					if useDhcp, ok := ipv4Map["use_dhcp"].(bool); ok && useDhcp {
						oneOfIPv4.SetValue(*import6.NewVmGcProfileUseDhcpOverrideSpec())
					} else {
						ipv4Override := import6.NewVmGcProfileNicIpv4ConfigOverrideSpec()
						if ipAddrList, ok := ipv4Map["ip_address"].([]interface{}); ok && len(ipAddrList) > 0 && ipAddrList[0] != nil {
							ipAddrMap := ipAddrList[0].(map[string]interface{})
							ipAddr := &config.IPv4Address{}
							if v, ok := ipAddrMap["value"].(string); ok && v != "" {
								ipAddr.Value = utils.StringPtr(v)
							}
							if pl, ok := ipAddrMap["prefix_length"].(int); ok {
								ipAddr.PrefixLength = utils.IntPtr(pl)
							}
							ipv4Override.IpAddress = ipAddr
						}
						if gws, ok := ipv4Map["default_gateways"].([]interface{}); ok && len(gws) > 0 {
							gwStrs := make([]string, len(gws))
							for j, g := range gws {
								gwStrs[j] = g.(string)
							}
							ipv4Override.DefaultGateways = gwStrs
						}
						oneOfIPv4.SetValue(*ipv4Override)
					}
					nic.Ipv4Config = oneOfIPv4
				}
				nics[i] = nic
			}
			ns.NicConfigList = nics
		}
		params.NetworkSettings = ns
	}

	if wdList, ok := spVal["workgroup_or_domain_info"].([]interface{}); ok && len(wdList) > 0 && wdList[0] != nil {
		wdVal := wdList[0].(map[string]interface{})
		oneOfWD := import6.NewOneOfVmGcProfileSysprepParamsOverrideSpecWorkgroupOrDomainInfo()
		if discard, ok := wdVal["discard"].(bool); ok && discard {
			oneOfWD.SetValue(*import6.NewVmGcProfileDiscardSettings())
		} else if wgList, ok := wdVal["workgroup"].([]interface{}); ok && len(wgList) > 0 && wgList[0] != nil {
			wgMap := wgList[0].(map[string]interface{})
			wg := import6.NewVmGcProfileWorkgroupOverrideSpec()
			if v, ok := wgMap["name"].(string); ok && v != "" {
				wg.Name = utils.StringPtr(v)
			}
			oneOfWD.SetValue(*wg)
		} else if dsList, ok := wdVal["domain_settings"].([]interface{}); ok && len(dsList) > 0 && dsList[0] != nil {
			dsMap := dsList[0].(map[string]interface{})
			ds := import6.NewVmGcProfileDomainSettingsOverrideSpec()
			if credsList, ok := dsMap["credentials"].([]interface{}); ok && len(credsList) > 0 && credsList[0] != nil {
				credsMap := credsList[0].(map[string]interface{})
				creds := import6.NewVmGcProfileDomainCredentialsOverrideSpec()
				if v, ok := credsMap["domain_name"].(string); ok && v != "" {
					creds.DomainName = utils.StringPtr(v)
				}
				if v, ok := credsMap["password"].(string); ok && v != "" {
					creds.Password = utils.StringPtr(v)
				}
				if v, ok := credsMap["username"].(string); ok && v != "" {
					creds.Username = utils.StringPtr(v)
				}
				ds.Credentials = creds
			}
			oneOfWD.SetValue(*ds)
		}
		params.WorkgroupOrDomainInfo = oneOfWD
	}

	return params
}

func flattenVmGcProfileConfigOverride(pr *import6.VmGcProfileConfig) []map[string]interface{} {
	if pr == nil {
		return nil
	}
	result := make(map[string]interface{})
	if pr.Profile != nil {
		result["profile"] = flattenVmGcProfileReference(pr.Profile)
	}
	if pr.ConfigOverrideSpec != nil {
		result["config_override_spec"] = flattenConfigOverrideSpec(pr.ConfigOverrideSpec)
	}
	return []map[string]interface{}{result}
}

func flattenConfigOverrideSpec(oneOf *import6.OneOfVmGcProfileConfigConfigOverrideSpec) []map[string]interface{} {
	if oneOf == nil {
		return nil
	}
	val := oneOf.GetValue()
	if val == nil {
		return nil
	}
	result := make(map[string]interface{})
	switch v := val.(type) {
	case import6.VmGcProfileSysprepConfigOverrideSpec:
		result["sysprep_config"] = flattenSysprepConfigOverrideSpec(&v)
	}
	return []map[string]interface{}{result}
}

func flattenSysprepConfigOverrideSpec(spec *import6.VmGcProfileSysprepConfigOverrideSpec) []map[string]interface{} {
	if spec == nil {
		return nil
	}
	result := make(map[string]interface{})
	if spec.Customization != nil {
		custMap := make(map[string]interface{})
		custVal := spec.Customization.GetValue()
		switch v := custVal.(type) {
		case import6.VmGcProfileAnswerFileOverrideSpec:
			custMap["answer_file"] = []map[string]interface{}{
				{"unattend_xml": v.UnattendXml},
			}
		case import6.VmGcProfileSysprepParamsOverrideSpec:
			custMap["sysprep_params"] = flattenSysprepParamsOverrideSpec(&v)
		}
		result["customization"] = []map[string]interface{}{custMap}
	}
	return []map[string]interface{}{result}
}

func flattenSysprepParamsOverrideSpec(params *import6.VmGcProfileSysprepParamsOverrideSpec) []map[string]interface{} {
	if params == nil {
		return nil
	}
	result := make(map[string]interface{})
	if len(params.FirstLogonCommands) > 0 {
		result["first_logon_commands"] = params.FirstLogonCommands
	}
	if params.GeneralSettings != nil {
		result["general_settings"] = flattenGeneralSettingsOverrideSpec(params.GeneralSettings)
	}
	if params.LocaleSettings != nil {
		result["locale_settings"] = flattenLocaleSettingsOverrideSpec(params.LocaleSettings)
	}
	if params.NetworkSettings != nil {
		result["network_settings"] = flattenNetworkSettingsOverrideSpec(params.NetworkSettings)
	}
	return []map[string]interface{}{result}
}

func flattenGeneralSettingsOverrideSpec(gs *import6.VmGcProfileGeneralSettingsOverrideSpec) []map[string]interface{} {
	if gs == nil {
		return nil
	}
	result := make(map[string]interface{})
	if gs.AutoLogonSettings != nil {
		result["auto_logon_settings"] = []map[string]interface{}{
			{"logon_count": gs.AutoLogonSettings.LogonCount},
		}
	}
	return []map[string]interface{}{result}
}

func flattenLocaleSettingsOverrideSpec(ls *import6.VmGcProfileLocaleSettingsOverrideSpec) []map[string]interface{} {
	if ls == nil {
		return nil
	}
	return []map[string]interface{}{make(map[string]interface{})}
}

func flattenNetworkSettingsOverrideSpec(ns *import6.VmGcProfileNetworkSettingsOverrideSpec) []map[string]interface{} {
	if ns == nil {
		return nil
	}
	result := make(map[string]interface{})
	if len(ns.NicConfigList) > 0 {
		nicList := make([]map[string]interface{}, len(ns.NicConfigList))
		for i, nic := range ns.NicConfigList {
			nicMap := make(map[string]interface{})
			if nic.DnsConfig != nil {
				dnsMap := map[string]interface{}{
					"preferred_dns_server_address":   nic.DnsConfig.PreferredDnsServerAddress,
					"alternate_dns_server_addresses": nic.DnsConfig.AlternateDnsServerAddresses,
				}
				nicMap["dns_config"] = []map[string]interface{}{dnsMap}
			}
			nicList[i] = nicMap
		}
		result["nic_config_list"] = nicList
	}
	return []map[string]interface{}{result}
}

func expandVmGcProfileReference(pr interface{}) *import6.VmGcProfileReference {
	if pr == nil {
		return nil
	}
	list, ok := pr.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	val := list[0].(map[string]interface{})
	ref := import6.NewVmGcProfileReference()
	if extID, ok := val["ext_id"]; ok && extID.(string) != "" {
		ref.ExtId = utils.StringPtr(extID.(string))
	}
	return ref
}

func flattenVmGcProfileReference(pr *import6.VmGcProfileReference) []map[string]interface{} {
	if pr == nil {
		return nil
	}
	return []map[string]interface{}{
		{"ext_id": pr.ExtId},
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
