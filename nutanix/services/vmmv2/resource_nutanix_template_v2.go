package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	vmmCommon "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/config"
	vmmAuthn "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/iam/v4/authn"
	vmmProsmConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	vmmConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	vmmContent "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixTemplatesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTemplatesV2Create,
		ReadContext:   ResourceNutanixTemplatesV2Read,
		UpdateContext: ResourceNutanixTemplatesV2Update,
		DeleteContext: ResourceNutanixTemplatesV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_version_spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": schemaForLinks(),
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version_description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vm_spec": schemaForTemplateVMSpec(),
						"created_by": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     schemaForTemplateUser(),
						},
						"version_source": schemaForVersionSource(),
						"version_source_discriminator": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_active_version": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"is_gc_override_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"guest_update_status": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployed_vm_reference": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"created_by": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schemaForTemplateUser(),
			},
			"updated_by": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schemaForTemplateUser(),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_ext_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
		},
	}
}

func ResourceNutanixTemplatesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	body := vmmContent.NewTemplate()

	if name, ok := d.GetOk("template_name"); ok {
		body.TemplateName = utils.StringPtr(name.(string))
	}
	if tempDesc, ok := d.GetOk("template_description"); ok {
		body.TemplateDescription = utils.StringPtr(tempDesc.(string))
	}
	if tempVersionSpec, ok := d.GetOk("template_version_spec"); ok {
		versionSpecData := tempVersionSpec.([]interface{})[0].(map[string]interface{})
		versionSourceData := versionSpecData["version_source"].([]interface{})[0].(map[string]interface{})
		log.Printf("[DEBUG] versionSource : %v", versionSourceData)
		if templateVMReference, ok := versionSourceData["template_vm_reference"]; ok {
			log.Printf("[DEBUG] templateVmReference : %v", templateVMReference)
			if len(templateVMReference.([]interface{})) == 0 {
				return diag.Errorf("template_vm_reference is required for template creation")
			}
			templateVMReferenceData := templateVMReference.([]interface{})[0].(map[string]interface{})
			log.Printf("[DEBUG] templateVmReferenceData : %v", templateVMReferenceData)
			vmExtID := templateVMReferenceData["ext_id"].(string)
			if vmExtID == "" {
				return diag.Errorf("ext_id is required for template_vm_reference")
			}
			templateVersionSourceObj := &vmmContent.OneOfTemplateVersionSpecVersionSource{}
			vmRefInput := vmmContent.NewTemplateVmReference()

			vmRefInput.ExtId = utils.StringPtr(vmExtID)
			if guest, ok := templateVMReferenceData["guest_customization"]; ok && len(guest.([]interface{})) > 0 {
				vmRefInput.GuestCustomization = expandTemplateGuestCustomizationParams(guest)
			}

			err := templateVersionSourceObj.SetValue(*vmRefInput)
			if err != nil {
				return diag.Errorf("error while setting version source : %v", err)
			}

			templateVersionSpecObj := &vmmContent.TemplateVersionSpec{}
			templateVersionSpecObj.VersionSource = templateVersionSourceObj

			aJSON, _ := json.Marshal(templateVersionSpecObj)
			log.Printf("[DEBUG] templateVersionSpecObj : %v", string(aJSON))
			body.TemplateVersionSpec = templateVersionSpecObj
		} else {
			return diag.Errorf("template_version_spec is required for template creation")
		}
	}
	if guestUpdateStatus, ok := d.GetOk("guest_update_status"); ok {
		body.GuestUpdateStatus = expandGuestUpdateStatus(guestUpdateStatus)
	}
	if createdBy, ok := d.GetOk("created_by"); ok && len(createdBy.([]interface{})) > 0 {
		body.CreatedBy = expandTemplateUser(createdBy)
	}
	if categoryExtIDs, ok := d.GetOk("category_ext_ids"); ok {
		body.CategoryExtIds = common.ExpandListOfString(categoryExtIDs.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Template create request body :\n %s", string(aJSON))
	resp, err := conn.TemplatesAPIInstance.CreateTemplate(body)
	if err != nil {
		return diag.Errorf("error while creating template : %v", err)
	}
	TaskRef := resp.Data.GetValue().(vmmProsmConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the template to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching template create task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Template Create Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeTemplates, "Template")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return ResourceNutanixTemplatesV2Read(ctx, d, meta)
}

func ResourceNutanixTemplatesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	tempVersionSpecData := d.Get("template_version_spec").([]interface{})
	log.Printf("[DEBUG] tempVersionSpecData: %v", tempVersionSpecData)
	resp, err := conn.TemplatesAPIInstance.GetTemplateById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching template : %v", err)
	}
	getResp := resp.Data.GetValue().(vmmContent.Template)
	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] Get Template call: %s", string(aJSON))
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
	// version source not returned in API response, so set the value from the config
	// this logic to ignore changes on terraform plan when there is no change in version source
	// if there is a change in version source, then terraform plan will show the change
	if len(tempVersionSpecData) > 0 {
		if versionSource, ok := tempVersionSpecData[0].(map[string]interface{})["version_source"]; ok {
			tempVersionSpecFlattened := d.Get("template_version_spec").([]interface{})[0].(map[string]interface{})
			tempVersionSpecFlattened["version_source"] = versionSource
			log.Printf("[DEBUG] template_version_spec flattened : %v", tempVersionSpecFlattened)
			if err := d.Set("template_version_spec", []map[string]interface{}{tempVersionSpecFlattened}); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	log.Printf("[DEBUG] template_version_spec : %v", d.Get("template_version_spec"))
	if err := d.Set("category_ext_ids", getResp.CategoryExtIds); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixTemplatesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	readResp, err := conn.TemplatesAPIInstance.GetTemplateById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching template : %v", err)
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	respTemplate := readResp.Data.GetValue().(vmmContent.Template)

	updateSpec := &vmmContent.Template{}
	updateSpec.ExtId = respTemplate.ExtId

	if d.HasChange("template_name") {
		updateSpec.TemplateName = utils.StringPtr(d.Get("template_name").(string))
	}
	if d.HasChange("template_description") {
		updateSpec.TemplateDescription = utils.StringPtr(d.Get("template_description").(string))
	}
	if d.HasChange("template_version_spec") {
		// version_name, version_description is required for update
		log.Printf("[DEBUG] HasChanged version_name: %v", d.HasChange("template_version_spec.0.version_name"))
		if !d.HasChange("template_version_spec.0.version_name") {
			return diag.Errorf("version_name is required for update operation, please provide/update version_name")
		}
		if !d.HasChange("template_version_spec.0.version_description") {
			return diag.Errorf("version_description is required for update operation, please provide/update version_description")
		}
		if vSpec, ok := d.GetOk("template_version_spec"); ok {
			vName, okName := vSpec.([]interface{})[0].(map[string]interface{})["version_name"]
			vDescription, okDsec := vSpec.([]interface{})[0].(map[string]interface{})["version_description"]
			log.Printf("[DEBUG] version_name : %v", vName)
			log.Printf("[DEBUG] ok : %v, !ok || vName== \"\": %v", ok, !ok || vName == "")
			if !okName || vName == "" {
				return diag.Errorf("version_name is required for update operation")
			}
			if !okDsec || vDescription == "" {
				return diag.Errorf("version_description is required for update operation")
			}
		}

		updateSpec.TemplateVersionSpec = expandTemplateVersionSpec(d.Get("template_version_spec"))
		updateSpec.TemplateVersionSpec.ExtId = nil
		updateSpec.TemplateVersionSpec.CreatedBy = nil
	}
	if d.HasChange("guest_update_status") {
		updateSpec.GuestUpdateStatus = expandGuestUpdateStatus(d.Get("guest_update_status"))
	}
	if d.HasChange("created_by") {
		updateSpec.CreatedBy = expandTemplateUser(d.Get("created_by"))
	}
	if d.HasChange("updated_by") {
		updateSpec.UpdatedBy = expandTemplateUser(d.Get("updated_by"))
	}

	if updateSpec.TemplateVersionSpec != nil &&
		updateSpec.TemplateVersionSpec.VersionSource != nil {
		log.Printf("[DEBUG] Check version id in tf configuration")
		templateVersionReference := updateSpec.TemplateVersionSpec.VersionSource.GetValue()
		//nolint:gocritic // Type switch not used intentionally for demonstration
		switch templateVersionReference.(type) {
		// we need only to set version id for TemplateVersionReference type
		case vmmContent.TemplateVersionReference:
			log.Printf("[DEBUG] Template version reference type")
			versionID := templateVersionReference.(vmmContent.TemplateVersionReference).VersionId

			if versionID != nil || utils.StringValue(versionID) != "" {
				log.Printf("[DEBUG] Template version Id provided in tf configuration")
			}
			log.Printf("[DEBUG] Template version Id not provided in tf configuration, will use the latest version as default")
			templateVersions, errTempVersion := conn.TemplatesAPIInstance.ListTemplateVersions(utils.StringPtr(d.Id()), nil, nil, nil, nil, nil)
			if errTempVersion != nil {
				return diag.Errorf("error while fetching template versions : %v", errTempVersion)
			}
			templateVersion := templateVersions.Data.GetValue().([]vmmContent.TemplateVersionSpec)
			tmplVersion := templateVersion[0]
			if len(templateVersion) == 0 {
				return diag.Errorf("No template versions found for template %s", d.Id())
			}
			for _, version := range templateVersion {
				if version.CreateTime.After(*tmplVersion.CreateTime) {
					tmplVersion = version
				}
			}
			versionSource := updateSpec.TemplateVersionSpec.VersionSource.GetValue().(vmmContent.TemplateVersionReference)
			versionSource.VersionId = tmplVersion.ExtId
			errVs := updateSpec.TemplateVersionSpec.VersionSource.SetValue(versionSource)
			if errVs != nil {
				return diag.Errorf("error while setting version source : %v", err)
			}
		case vmmContent.TemplateVmReference:
			log.Printf("[DEBUG] Template vm reference type, no need to set version id")
		default:
			log.Printf("[DEBUG] Template version reference type not found")
		}
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] Template update request body :\n %v", string(aJSON))

	respUpdate, err := conn.TemplatesAPIInstance.UpdateTemplateById(utils.StringPtr(d.Id()), updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating template : %v", err)
	}

	TaskRef := respUpdate.Data.GetValue().(vmmProsmConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the template to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixTemplatesV2Read(ctx, d, meta)
}

func ResourceNutanixTemplatesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.TemplatesAPIInstance.DeleteTemplateById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting template : %v", err)
	}
	TaskRef := resp.Data.GetValue().(vmmProsmConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the template to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

// Schema's functions
func schemaForLinks() *schema.Schema {
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

func schemaForTemplateVMSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"links": schemaForLinks(),
				"ext_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"create_time": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"update_time": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"source": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"entity_type": {
								Type:         schema.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: validation.StringInSlice([]string{"VM_RECOVERY_POINT", "VM"}, false),
							},
							"ext_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"num_sockets": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"num_cores_per_socket": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"num_threads_per_core": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"num_numa_nodes": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"memory_size_bytes": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"is_vcpu_hard_pinning_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"is_cpu_passthrough_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"enabled_cpu_features": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{"HARDWARE_VIRTUALIZATION"}, false),
					},
				},
				"is_memory_overcommit_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"is_gpu_console_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"is_cpu_hotplug_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"is_scsi_controller_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"generation_uuid": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"bios_uuid": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"categories": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
				"ownership_info": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"owner": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ext_id": {
											Type:     schema.TypeString,
											Optional: true,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
				"host": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
				"cluster": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
				// not visible in API reference
				"availability_zone": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
				"guest_customization": schemaForTemplateGuestCustomization(),
				"guest_tools": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"version": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"is_installed": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"is_iso_inserted": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"available_version": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"guest_os_version": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"is_reachable": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"is_vss_snapshot_capable": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"is_vm_mobility_drivers_installed": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"is_enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Computed: true,
							},
							"capabilities": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Schema{
									Type:         schema.TypeString,
									ValidateFunc: validation.StringInSlice([]string{"SELF_SERVICE_RESTORE", "VSS_SNAPSHOT"}, false),
								},
							},
						},
					},
				},
				"hardware_clock_timezone": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"is_branding_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"boot_config": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"legacy_boot": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"boot_device": schemaForBootDevice(),
										"boot_order":  schemaForBootOrder(),
									},
								},
							},
							"uefi_boot": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"is_secure_boot_enabled": {
											Type:     schema.TypeBool,
											Optional: true,
											Computed: true,
										},
										"nvram_device": {
											Type:     schema.TypeList,
											Optional: true,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"backing_storage_info": {
														Type:     schema.TypeList,
														Optional: true,
														Computed: true,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"disk_ext_id": {
																	Type:     schema.TypeString,
																	Computed: true,
																},
																"disk_size_bytes": {
																	Type:     schema.TypeInt,
																	Optional: true,
																	Computed: true,
																},
																"storage_container": schemaForStorageContainer(),
																"storage_config":    schemaForStorageConfig(),
																"data_source": {
																	Type:     schema.TypeList,
																	Optional: true,
																	Computed: true,
																	Elem: &schema.Resource{
																		Schema: map[string]*schema.Schema{
																			"reference": {
																				Type:     schema.TypeList,
																				Optional: true,
																				Computed: true,
																				Elem: &schema.Resource{
																					Schema: map[string]*schema.Schema{
																						"image_reference": {
																							Type:     schema.TypeList,
																							Optional: true,
																							Computed: true,
																							Elem: &schema.Resource{
																								Schema: map[string]*schema.Schema{
																									"image_ext_id": {
																										Type:     schema.TypeString,
																										Optional: true,
																										Computed: true,
																									},
																								},
																							},
																						},
																						"vm_disk_reference": {
																							Type:     schema.TypeList,
																							Optional: true,
																							Computed: true,
																							Elem: &schema.Resource{
																								Schema: map[string]*schema.Schema{
																									"disk_ext_id": {
																										Type:     schema.TypeString,
																										Optional: true,
																										Computed: true,
																									},
																									"disk_address": {
																										Type:     schema.TypeList,
																										Optional: true,
																										Computed: true,
																										Elem: &schema.Resource{
																											Schema: map[string]*schema.Schema{
																												"bus_type": {
																													Type:         schema.TypeString,
																													Optional:     true,
																													Computed:     true,
																													ValidateFunc: validation.StringInSlice([]string{"SCSI", "SPAPR", "PCI", "IDE", "SATA"}, false),
																												},
																												"index": {
																													Type:     schema.TypeInt,
																													Optional: true,
																													Computed: true,
																												},
																											},
																										},
																									},
																									"vm_reference": {
																										Type:     schema.TypeList,
																										Optional: true,
																										Computed: true,
																										Elem: &schema.Resource{
																											Schema: map[string]*schema.Schema{
																												"ext_id": {
																													Type:     schema.TypeString,
																													Optional: true,
																													Computed: true,
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
																"is_migration_in_progress": {
																	Type:     schema.TypeBool,
																	Computed: true,
																},
															},
														},
													},
												},
											},
										},
										"boot_device": schemaForBootDevice(),
										"boot_order":  schemaForBootOrder(),
									},
								},
							},
						},
					},
				},
				"is_vga_console_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"machine_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice([]string{"PSERIES", "Q35", "PC"}, false),
				},
				"power_state": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice([]string{"PAUSED", "UNDETERMINED", "OFF", "ON"}, false),
				},
				"vtpm_config": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"is_vtpm_enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Computed: true,
							},
							"version": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"is_agent_vm": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"apc_config": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"is_apc_enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Computed: true,
							},
							"cpu_model": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ext_id": {
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
						},
					},
				},
				"is_live_migrate_capable": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"is_cross_cluster_migration_in_progress": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"storage_config": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"is_flash_mode_enabled": {
								Type:     schema.TypeBool,
								Optional: true,
								Computed: true,
							},
							"qos_config": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"throttled_iops": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
				"disks": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tenant_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"links": schemaForLinks(),
							"ext_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"disk_address": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"bus_type": {
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											ValidateFunc: validation.StringInSlice([]string{"SCSI", "SPAPR", "PCI", "IDE", "SATA"}, false),
										},
										"index": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
									},
								},
							},
							"backing_info": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"vm_disk": {
											Type:     schema.TypeList,
											Optional: true,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"disk_ext_id": {
														Type:     schema.TypeString,
														Computed: true,
													},
													"disk_size_bytes": {
														Type:     schema.TypeInt,
														Optional: true,
														Computed: true,
													},
													"storage_container": schemaForStorageContainer(),
													"storage_config":    schemaForStorageConfig(),
													"data_source":       schemaForDataSource(),
													"is_migration_in_progress": {
														Type:     schema.TypeBool,
														Optional: true,
														Computed: true,
													},
												},
											},
										},
										"adfs_volume_group_reference": {
											Type:     schema.TypeList,
											Optional: true,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"volume_group_ext_id": {
														Type:     schema.TypeString,
														Optional: true,
														Computed: true,
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
				"cd_roms": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tenant_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"links": schemaForLinks(),
							"ext_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"disk_address": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"bus_type": {
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											ValidateFunc: validation.StringInSlice([]string{"IDE", "SATA"}, false),
										},
										"index": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
									},
								},
							},
							"backing_info": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"disk_ext_id": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"disk_size_bytes": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
										"storage_container": schemaForStorageContainer(),
										"storage_config":    schemaForStorageConfig(),
										"data_source":       schemaForDataSource(),
										"is_migration_in_progress": {
											Type:     schema.TypeBool,
											Computed: true,
										},
									},
								},
							},
							"iso_type": {
								Type:         schema.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: validation.StringInSlice([]string{"OTHER", "GUEST_TOOLS", "GUEST_CUSTOMIZATION"}, false),
							},
						},
					},
				},
				"nics": schemaForNics(),
				"gpus": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tenant_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"links": schemaForLinks(),
							"ext_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"mode": {
								Type:         schema.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: validation.StringInSlice([]string{"PASSTHROUGH_GRAPHICS", "PASSTHROUGH_COMPUTE", "VIRTUAL"}, false),
							},
							"device_id": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
							},
							"vendor": {
								Type:         schema.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: validation.StringInSlice([]string{"NVIDIA", "AMD", "INTEL"}, false),
							},
							// not visible in API reference
							"pci_address": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"segment": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
										"bus": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
										"device": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
										"func": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
									},
								},
							},
							"guest_driver_version": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"frame_buffer_size_bytes": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"num_virtual_display_heads": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"fraction": {
								Type:     schema.TypeInt,
								Computed: true,
							},
						},
					},
				},
				"serial_ports": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tenant_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"links": schemaForLinks(),
							"ext_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"is_connected": {
								Type:     schema.TypeBool,
								Optional: true,
								Computed: true,
							},
							"index": {
								Type:         schema.TypeInt,
								Optional:     true,
								Computed:     true,
								ValidateFunc: validation.IntBetween(0, 3),
							},
						},
					},
				},
				"protection_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice([]string{"PD_PROTECTED", "UNPROTECTED", "RULE_PROTECTED"}, false),
				},
				"protection_policy_state": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"policy": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ext_id": {
											Type:     schema.TypeString,
											Optional: true,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
				"pci_devices": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tenant_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"links": schemaForLinks(),
							"ext_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"assigned_device_info": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"device": {
											Type:     schema.TypeList,
											Optional: true,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"device_ext_id": {
														Type:     schema.TypeString,
														Optional: true,
														Computed: true,
													},
												},
											},
										},
									},
								},
							},
							"backing_info": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"pcie_device_reference": {
											Type:     schema.TypeList,
											Optional: true,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"device_ext_id": {
														Type:     schema.TypeString,
														Optional: true,
														Computed: true,
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

func schemaForVersionSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"template_vm_reference": {
					Type:         schema.TypeList,
					Optional:     true,
					Computed:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"template_version_spec.0.version_source.0.template_version_reference", "template_version_spec.0.version_source.0.template_vm_reference"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:     schema.TypeString,
								Required: true,
							},
							"guest_customization": schemaForTemplateGuestCustomization(),
						},
					},
				},
				"template_version_reference": {
					Type:         schema.TypeList,
					Optional:     true,
					Computed:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"template_version_spec.0.version_source.0.template_version_reference", "template_version_spec.0.version_source.0.template_vm_reference"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"version_id": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
							"override_vm_config": {
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Optional: true,
											Computed: true,
										},
										"num_sockets": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
										"num_cores_per_socket": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
										"num_threads_per_core": {
											Type:     schema.TypeInt,
											Optional: true,
											Computed: true,
										},
										"memory_size_bytes": {
											Type:     schema.TypeInt,
											Optional: true,
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
	}
}

func schemaForTemplateGuestCustomization() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"config": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"sysprep": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"install_type": {
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											ValidateFunc: validation.StringInSlice([]string{"PREPARED", "FRESH"}, false),
										},
										"sysprep_script": {
											Type:     schema.TypeList,
											Optional: true,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"unattend_xml": {
														Type:     schema.TypeList,
														Optional: true,
														Computed: true,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"value": {
																	Type:     schema.TypeString,
																	Optional: true,
																	Computed: true,
																},
															},
														},
													},
													"custom_key_values": schemaForCustomKeyValuePairs(),
												},
											},
										},
									},
								},
							},
							"cloud_init": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"datasource_type": {
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											ValidateFunc: validation.StringInSlice([]string{"CONFIG_DRIVE_V2"}, false),
										},
										"metadata": {
											Type:     schema.TypeString,
											Optional: true,
											Computed: true,
										},
										"cloud_init_script": {
											Type:     schema.TypeList,
											Optional: true,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"user_data": {
														Type:     schema.TypeList,
														Optional: true,
														Computed: true,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"value": {
																	Type:     schema.TypeString,
																	Optional: true,
																	Computed: true,
																},
															},
														},
													},
													"custom_key_values": schemaForCustomKeyValuePairs(),
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

func schemaForCustomKeyValuePairs() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key_value_pairs": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
							"value": schemaForValue(),
						},
					},
				},
			},
		},
	}
}

func schemaForValue() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"string": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"integer": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"boolean": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"string_list": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"object": {
					Type:     schema.TypeMap,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"map_of_strings": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"map": {
								Type:     schema.TypeMap,
								Optional: true,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"integer_list": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
		},
	}
}

func schemaForBootDevice() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"boot_device_disk": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"disk_address": schemaForDiskAddress(),
						},
					},
				},
				"boot_device_nic": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"mac_address": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func schemaForDiskAddress() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bus_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice([]string{"SCSI", "SPAPR", "PCI", "IDE", "SATA"}, false),
				},
				"index": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func schemaForBootOrder() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"CDROM", "DISK", "NETWORK"}, false),
		},
	}
}

func schemaForStorageContainer() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ext_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func schemaForStorageConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_flash_mode_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func schemaForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"reference": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"image_reference": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"image_ext_id": {
											Type:     schema.TypeString,
											Optional: true,
											Computed: true,
										},
									},
								},
							},
							"vm_disk_reference": {
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"disk_ext_id": {
											Type:     schema.TypeString,
											Optional: true,
											Computed: true,
										},
										"disk_address": schemaForDiskAddress(),
										"vm_reference": {
											Type:     schema.TypeList,
											Optional: true,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"ext_id": {
														Type:     schema.TypeString,
														Optional: true,
														Computed: true,
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

func schemaForNics() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem:     nicsElemSchemaV2WithTenantLinks(),
	}
}

func schemaForTemplateUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SERVICE_ACCOUNT", "LDAP", "EXTERNAL", "LOCAL", "SAML"}, false),
			},
			"idp_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"middle_initial": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_force_reset_password_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"additional_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": schemaForValue(),
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"creation_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PREDEFINED", "SERVICEDEFINED", "USERDEFINED"}, false),
			},
		},
	}
}

// expanders
func expandTemplateVersionSpec(pr interface{}) *vmmContent.TemplateVersionSpec {
	if pr.([]interface{}) != nil {
		cfg := &vmmContent.TemplateVersionSpec{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if versionName, ok := val["version_name"]; ok && len(versionName.(string)) > 0 {
			cfg.VersionName = utils.StringPtr(versionName.(string))
		}
		if versionDescription, ok := val["version_description"]; ok && len(versionDescription.(string)) > 0 {
			cfg.VersionDescription = utils.StringPtr(versionDescription.(string))
		}
		if vmSpec, ok := val["vm_spec"]; ok {
			cfg.VmSpec = expandTemplateVMSpec(vmSpec)
		}
		if createdBy, ok := val["created_by"]; ok && len(createdBy.([]interface{})) > 0 {
			cfg.CreatedBy = expandTemplateUser(createdBy)
		}
		if versionSource, ok := val["version_source"]; ok {
			cfg.VersionSource = expandTemplateVersionSpecVersionSource(versionSource)
		}
		if versionSourceDiscriminator, ok := val["version_source_discriminator"]; ok {
			cfg.VersionSourceDiscriminator = utils.StringPtr(versionSourceDiscriminator.(string))
		}
		if isActive, ok := val["is_active_version"]; ok {
			cfg.IsActiveVersion = utils.BoolPtr(isActive.(bool))
		}
		if isGcOverride, ok := val["is_gc_override_enabled"]; ok {
			cfg.IsGcOverrideEnabled = utils.BoolPtr(isGcOverride.(bool))
		}

		aJSON, _ := json.Marshal(cfg)
		log.Printf("[DEBUG] expandTemplateVersionSpec: %s", string(aJSON))
		return cfg
	}
	return nil
}

func expandTemplateVMSpec(vmSpec interface{}) *vmmConfig.Vm {
	if len(vmSpec.([]interface{})) > 0 {
		vm := &vmmConfig.Vm{}
		vmI := vmSpec.([]interface{})
		vmVal := vmI[0].(map[string]interface{})

		if name, ok := vmVal["name"]; ok {
			vm.Name = utils.StringPtr(name.(string))
		}
		if description, ok := vmVal["description"]; ok {
			vm.Description = utils.StringPtr(description.(string))
		}
		if source, ok := vmVal["source"]; ok {
			vm.Source = expandVMSourceReference(source)
		}
		if numSockets, ok := vmVal["num_sockets"]; ok {
			vm.NumSockets = utils.IntPtr(numSockets.(int))
		}
		if numCoresPerSocket, ok := vmVal["num_cores_per_socket"]; ok {
			vm.NumCoresPerSocket = utils.IntPtr(numCoresPerSocket.(int))
		}
		if numThreadsPerCore, ok := vmVal["num_threads_per_core"]; ok {
			vm.NumThreadsPerCore = utils.IntPtr(numThreadsPerCore.(int))
		}
		if memorySizeBytes, ok := vmVal["memory_size_bytes"]; ok {
			vm.MemorySizeBytes = utils.Int64Ptr(int64(memorySizeBytes.(int)))
		}
		if isCPUPassThroughEnabled, ok := vmVal["is_vcpu_hard_pinning_enabled"]; ok {
			vm.IsVcpuHardPinningEnabled = utils.BoolPtr(isCPUPassThroughEnabled.(bool))
		}
		if isCPUPassThroughEnabled, ok := vmVal["is_cpu_passthrough_enabled"]; ok {
			vm.IsCpuPassthroughEnabled = utils.BoolPtr(isCPUPassThroughEnabled.(bool))
		}
		if enableCPUFeatures, ok := vmVal["enabled_cpu_features"]; ok && len(enableCPUFeatures.([]interface{})) > 0 {
			enCPUFeaturesList := enableCPUFeatures.([]interface{})
			const hardwareVirt = 2
			subMap := map[string]interface{}{
				"HARDWARE_VIRTUALIZATION": hardwareVirt,
			}
			var cpuFeatures []vmmConfig.CpuFeature
			for _, feature := range enCPUFeaturesList {
				pVal := subMap[feature.(string)]
				if pVal != nil {
					cpuFeatures = append(cpuFeatures, vmmConfig.CpuFeature(pVal.(int)))
				}
			}
			vm.EnabledCpuFeatures = cpuFeatures
		}
		if isMemoryOvercommitEnabled, ok := vmVal["is_memory_overcommit_enabled"]; ok {
			vm.IsMemoryOvercommitEnabled = utils.BoolPtr(isMemoryOvercommitEnabled.(bool))
		}
		if isGpuConsoleEnabled, ok := vmVal["is_gpu_console_enabled"]; ok {
			vm.IsGpuConsoleEnabled = utils.BoolPtr(isGpuConsoleEnabled.(bool))
		}
		if isCPUHotplugEnabled, ok := vmVal["is_cpu_hotplug_enabled"]; ok {
			vm.IsCpuHotplugEnabled = utils.BoolPtr(isCPUHotplugEnabled.(bool))
		}
		if isScsiControllerEnabled, ok := vmVal["is_scsi_controller_enabled"]; ok {
			vm.IsScsiControllerEnabled = utils.BoolPtr(isScsiControllerEnabled.(bool))
		}
		if generationUUID, ok := vmVal["generation_uuid"]; ok && generationUUID != "" {
			vm.GenerationUuid = utils.StringPtr(generationUUID.(string))
		}
		if biosUUID, ok := vmVal["bios_uuid"]; ok && biosUUID != "" {
			vm.BiosUuid = utils.StringPtr(biosUUID.(string))
		}
		if categories, ok := vmVal["categories"]; ok {
			vm.Categories = expandCategoryReference(categories.([]interface{}))
		}
		if ownershipInfo, ok := vmVal["ownership_info"]; ok {
			vm.OwnershipInfo = expandOwnershipInfo(ownershipInfo)
		}
		if host, ok := vmVal["host"]; ok {
			vm.Host = expandHostReference(host)
		}
		if cluster, ok := vmVal["cluster"]; ok {
			vm.Cluster = expandClusterReference(cluster)
		}
		if availabilityZone, ok := vmVal["availability_zone"]; ok {
			vm.AvailabilityZone = expandAvailabilityZone(availabilityZone)
		}
		if guestCustomization, ok := vmVal["guest_customization"]; ok {
			vm.GuestCustomization = expandTemplateGuestCustomizationParams(guestCustomization)
		}
		if guestTools, ok := vmVal["guest_tools"]; ok {
			vm.GuestTools = expandGuestTools(guestTools)
		}
		if hardwareClockTimezone, ok := vmVal["hardware_clock_timezone"]; ok {
			vm.HardwareClockTimezone = utils.StringPtr(hardwareClockTimezone.(string))
		}
		if isBrandingEnabled, ok := vmVal["is_branding_enabled"]; ok {
			vm.IsBrandingEnabled = utils.BoolPtr(isBrandingEnabled.(bool))
		}
		if bootConfig, ok := vmVal["boot_config"]; ok {
			vm.BootConfig = expandOneOfVMBootConfig(bootConfig)
		}
		if isVgaConsoleEnabled, ok := vmVal["is_vga_console_enabled"]; ok {
			vm.IsVgaConsoleEnabled = utils.BoolPtr(isVgaConsoleEnabled.(bool))
		}
		if machineType, ok := vmVal["machine_type"]; ok && machineType != "" {
			const two, three, four = 2, 3, 4
			subMap := map[string]interface{}{
				"PC":      two,
				"PSERIES": three,
				"Q35":     four,
			}
			pVal := subMap[machineType.(string)]
			p := vmmConfig.MachineType(pVal.(int))
			vm.MachineType = &p
		}
		if powerState, ok := vmVal["power_state"]; ok && powerState != "" {
			const two, three, four, five = 2, 3, 4, 5
			subMap := map[string]interface{}{
				"ON":           two,
				"OFF":          three,
				"PAUSED":       four,
				"UNDETERMINED": five,
			}
			pVal := subMap[powerState.(string)]
			p := vmmConfig.PowerState(pVal.(int))
			vm.PowerState = &p
		}
		if vtpmConfig, ok := vmVal["vtpm_config"]; ok {
			vm.VtpmConfig = expandVtpmConfig(vtpmConfig)
		}
		if isAgentVM, ok := vmVal["is_agent_vm"]; ok {
			vm.IsAgentVm = utils.BoolPtr(isAgentVM.(bool))
		}
		if apcConfig, ok := vmVal["apc_config"]; ok {
			vm.ApcConfig = expandApcConfig(apcConfig)
		}
		if storageConfig, ok := vmVal["storage_config"]; ok {
			vm.StorageConfig = expandADSFVmStorageConfig(storageConfig)
		}
		if disks, ok := vmVal["disks"]; ok {
			vm.Disks = expandDisk(disks.([]interface{}))
		}
		if cdRoms, ok := vmVal["cd_roms"]; ok {
			vm.CdRoms = expandCdRom(cdRoms.([]interface{}))
		}
		if nics, ok := vmVal["nics"]; ok {
			vm.Nics = expandNic(nics.([]interface{}))
		}
		if gpus, ok := vmVal["gpus"]; ok {
			vm.Gpus = expandGpu(gpus.([]interface{}))
		}
		if serialPorts, ok := vmVal["serial_ports"]; ok {
			vm.SerialPorts = expandSerialPort(serialPorts.([]interface{}))
		}
		if protectionType, ok := vmVal["protection_type"]; ok && protectionType != "" {
			const two, three, four = 2, 3, 4
			subMap := map[string]interface{}{
				"UNPROTECTED":    two,
				"PD_PROTECTED":   three,
				"RULE_PROTECTED": four,
			}
			pVal := subMap[protectionType.(string)]
			p := vmmConfig.ProtectionType(pVal.(int))
			vm.ProtectionType = &p
		}
		if protectionPolicyState, ok := vmVal["protection_policy_state"]; ok {
			vm.ProtectionPolicyState = expandProtectionPolicyState(protectionPolicyState)
		}
		if pcieDevices, ok := vmVal["pcie_devices"]; ok {
			vm.PcieDevices = expandPcieDevices(pcieDevices)
		}
		return vm
	}
	return nil
}

func expandAvailabilityZone(availabilityZone interface{}) *vmmConfig.AvailabilityZoneReference {
	if availabilityZone != nil && len(availabilityZone.([]interface{})) > 0 {
		availabilityZoneObj := &vmmConfig.AvailabilityZoneReference{}
		availabilityZoneData := availabilityZone.([]interface{})

		if extID := availabilityZoneData[0].(map[string]interface{})["ext_id"]; extID != nil {
			availabilityZoneObj.ExtId = utils.StringPtr(extID.(string))
		}
		return availabilityZoneObj
	}
	return nil
}

func expandPcieDevices(pcieDevices interface{}) []vmmConfig.PcieDevice {
	if len(pcieDevices.([]interface{})) > 0 {
		var pcieDevicesList []vmmConfig.PcieDevice

		for _, pcieDevice := range pcieDevices.([]interface{}) {
			pcieDeviceObj := vmmConfig.PcieDevice{}
			pcieDeviceData := pcieDevice.(map[string]interface{})

			if assignedDeviceInfo, ok := pcieDeviceData["assigned_device_info"]; ok {
				pcieDeviceObj.AssignedDeviceInfo = expandTemplateAssignedDeviceInfo(assignedDeviceInfo)
			}
			if backingInfo, ok := pcieDeviceData["backing_info"]; ok {
				pcieDeviceObj.BackingInfo = expandTemplatePcieDeviceBackingInfo(backingInfo)
			}
			pcieDevicesList = append(pcieDevicesList, pcieDeviceObj)
		}
		return pcieDevicesList
	}
	return nil
}

func expandTemplateAssignedDeviceInfo(assignedDeviceInfo interface{}) *vmmConfig.PcieDeviceInfo {
	if assignedDeviceInfo != nil {
		assignedDeviceInfoObj := &vmmConfig.PcieDeviceInfo{}
		assignedDeviceInfoData := assignedDeviceInfo.(map[string]interface{})

		if device, ok := assignedDeviceInfoData["device"]; ok {
			deviceObj := &vmmConfig.PcieDeviceReference{}
			deviceData := device.(map[string]interface{})

			if deviceExtID := deviceData["device_ext_id"]; deviceExtID != nil {
				deviceObj.DeviceExtId = utils.StringPtr(deviceExtID.(string))
			}
			assignedDeviceInfoObj.Device = deviceObj
		}
		return assignedDeviceInfoObj
	}
	return nil
}

func expandTemplatePcieDeviceBackingInfo(backingInfo interface{}) *vmmConfig.OneOfPcieDeviceBackingInfo {
	if backingInfo != nil {
		backingInfoObj := &vmmConfig.OneOfPcieDeviceBackingInfo{}
		backingInfoData := backingInfo.(map[string]interface{})

		if pcieDeviceReference, ok := backingInfoData["pcie_device_reference"]; ok {
			pcieDeviceReferenceObj := &vmmConfig.PcieDeviceReference{}
			pcieDeviceReferenceData := pcieDeviceReference.(map[string]interface{})

			if deviceExtID := pcieDeviceReferenceData["device_ext_id"]; deviceExtID != nil {
				pcieDeviceReferenceObj.DeviceExtId = utils.StringPtr(deviceExtID.(string))
			}
			err := backingInfoObj.SetValue(pcieDeviceReferenceObj)
			if err != nil {
				log.Printf("[ERROR] Error setting value for pcie_device_reference: %v", err)
				diag.Errorf("Error setting value for pcie_device_reference: %v", err)
				return nil
			}
		}
	}
	return nil
}

func expandTemplateVersionSpecVersionSource(versionSource interface{}) *vmmContent.OneOfTemplateVersionSpecVersionSource {
	if len(versionSource.([]interface{})) > 0 {
		templateVersionSpecVersionSource := &vmmContent.OneOfTemplateVersionSpecVersionSource{}
		versionSourceData := versionSource.([]interface{})[0].(map[string]interface{})

		if templateVMReference, ok := versionSourceData["template_vm_reference"]; ok && len(templateVMReference.([]interface{})) > 0 {
			vmRefInput := vmmContent.NewTemplateVmReference()
			prI := templateVMReference.([]interface{})
			val := prI[0].(map[string]interface{})

			if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
				vmRefInput.ExtId = utils.StringPtr(extID.(string))
			}
			if guest, ok := val["guest_customization"]; ok && len(guest.([]interface{})) > 0 {
				vmRefInput.GuestCustomization = expandTemplateGuestCustomizationParams(guest)
			}
			aJSON, _ := json.Marshal(vmRefInput)
			log.Printf("[DEBUG] templateVMReference: %v", string(aJSON))
			err := templateVersionSpecVersionSource.SetValue(*vmRefInput)
			if err != nil {
				log.Printf("[ERROR] templateVMReference: Error setting value for templateVMReference: %v", err)
				return nil
			}
		}
		if templateVersionReference, ok := versionSourceData["template_version_reference"]; ok && len(templateVersionReference.([]interface{})) > 0 {
			versionReference := vmmContent.NewTemplateVersionReference()
			prI := templateVersionReference.([]interface{})
			val := prI[0].(map[string]interface{})

			if versionID, ok := val["version_id"]; ok && len(versionID.(string)) > 0 {
				versionReference.VersionId = utils.StringPtr(versionID.(string))
			}
			if overrideVMConfig, ok := val["override_vm_config"]; ok && len(overrideVMConfig.([]interface{})) > 0 {
				versionReference.OverrideVmConfig = expandVMConfigOverrideTemplate(overrideVMConfig)
			}
			aJSON, _ := json.Marshal(versionReference)
			log.Printf("[DEBUG] templateVersionReference: %v", string(aJSON))
			err := templateVersionSpecVersionSource.SetValue(*versionReference)
			if err != nil {
				log.Printf("[ERROR] templateVersionReference: Error setting value for templateVersionReference: %v", err)
				return nil
			}
		}

		return templateVersionSpecVersionSource
	}
	return nil
}

func expandTemplateGuestCustomizationParams(guestCustomization interface{}) *vmmConfig.GuestCustomizationParams {
	if len(guestCustomization.([]interface{})) > 0 {
		guestCustomizationParams := &vmmConfig.GuestCustomizationParams{}
		guestCustomizationData := guestCustomization.([]interface{})[0].(map[string]interface{})

		if config, ok := guestCustomizationData["config"]; ok && len(config.([]interface{})) > 0 {
			log.Printf("[DEBUG] guestCustomizationParams.Config: %v", config)
			guestCustomizationParams.Config = expandTemplateGuestCustomizationConfig(config)
		}
		aJSON, _ := json.Marshal(guestCustomizationParams)
		log.Printf("[DEBUG] guestCustomizationParams: %v", string(aJSON))

		return guestCustomizationParams
	}
	return nil
}

func expandTemplateGuestCustomizationConfig(config interface{}) *vmmConfig.OneOfGuestCustomizationParamsConfig {
	if len(config.([]interface{})) > 0 {
		guestCustomizationConfig := vmmConfig.NewOneOfGuestCustomizationParamsConfig()
		configData := config.([]interface{})[0].(map[string]interface{})

		if sysprep, ok := configData["sysprep"]; ok && len(sysprep.([]interface{})) > 0 {
			sysprepObj := vmmConfig.NewSysprep()
			sysprepData := sysprep.([]interface{})[0].(map[string]interface{})

			if installType, ok := sysprepData["install_type"]; ok {
				if installType != nil && installType != "" {
					const two, three = 2, 3
					subMap := map[string]interface{}{
						"FRESH":    two,
						"PREPARED": three,
					}
					pVal := subMap[installType.(string)]
					if pVal == nil {
						sysprepObj.InstallType = nil
					}
					p := vmmConfig.InstallType(pVal.(int))
					sysprepObj.InstallType = &p
				}
			}
			if sysprepScript, ok := sysprepData["sysprep_script"]; ok && len(sysprepScript.([]interface{})) > 0 {
				sysprepObj.SysprepScript = expandSysprepScript(sysprepScript)
			}
			aJSON, _ := json.Marshal(sysprepObj)
			log.Printf("[DEBUG] sysprep.sysprep_script expanded: %v", string(aJSON))
			err := guestCustomizationConfig.SetValue(*sysprepObj)
			if err != nil {
				log.Printf("[ERROR] Error setting value for sysprep: %v", err)
				return nil
			}
		}

		if cloudInit, ok := configData["cloud_init"]; ok && len(cloudInit.([]interface{})) > 0 {
			cloudInitObj := vmmConfig.NewCloudInit()
			cloudInitData := cloudInit.([]interface{})[0].(map[string]interface{})

			if datasourceType, ok := cloudInitData["datasource_type"]; ok && len(datasourceType.(string)) > 0 {
				if datasourceType != nil && datasourceType != "" {
					const two = 2
					subMap := map[string]interface{}{
						"CONFIG_DRIVE_V2": two,
					}
					pVal := subMap[datasourceType.(string)]
					if pVal == nil {
						cloudInitObj.DatasourceType = nil
					}
					p := vmmConfig.CloudInitDataSourceType(pVal.(int))
					cloudInitObj.DatasourceType = &p
				}
			}
			if metadata, ok := cloudInitData["metadata"]; ok && len(metadata.(string)) > 0 {
				cloudInitObj.Metadata = utils.StringPtr(metadata.(string))
			}
			if cloudInitScript, ok := cloudInitData["cloud_init_script"]; ok && len(cloudInitScript.([]interface{})) > 0 {
				cloudInitScriptObj := vmmConfig.NewOneOfCloudInitCloudInitScript()
				cloudInitScriptData := cloudInitScript.([]interface{})[0].(map[string]interface{})

				if userdata := cloudInitScriptData["user_data"]; userdata != nil && len(userdata.([]interface{})) > 0 {
					user := vmmConfig.NewUserdata()
					userVal := userdata.([]interface{})[0].(map[string]interface{})

					if value, ok := userVal["value"]; ok {
						user.Value = utils.StringPtr(value.(string))
					}

					err := cloudInitScriptObj.SetValue(*user)
					if err != nil {
						log.Printf("[ERROR] cloudInitScript : Error setting value for userdata: %v", err)
						return nil
					}
				}
				if customKeyValues, ok := cloudInitScriptData["custom_key_values"]; ok && len(customKeyValues.([]interface{})) > 0 {
					log.Printf("[DEBUG] cloud_init.cloud_init_script.customKeyValues: %v", customKeyValues)
					customKeyValuesObj := expandTemplateCustomKeyValuesPairs(customKeyValues)
					aJSON, _ := json.Marshal(customKeyValuesObj)
					log.Printf("[DEBUG] cloud_init.cloud_init_script.customKeyValues expanded: %v", string(aJSON))
					err := cloudInitScriptObj.SetValue(*customKeyValuesObj)
					if err != nil {
						log.Printf("[ERROR] cloudInitScript: Error setting value for custom key values: %v", err)
						return nil
					}
				}
				cloudInitObj.CloudInitScript = cloudInitScriptObj
			}

			aJSON, _ := json.Marshal(cloudInitObj)
			log.Printf("[DEBUG] cloudInitObj expanded: %v", string(aJSON))

			err := guestCustomizationConfig.SetValue(*cloudInitObj)
			if err != nil {
				log.Printf("[ERROR] Error setting value for cloud init: %v", err)
				return nil
			}
		}

		aJSON, _ := json.Marshal(guestCustomizationConfig)
		log.Printf("[DEBUG] guestCustomizationConfig expanded: %v", string(aJSON))
		return guestCustomizationConfig
	}
	return nil
}

func expandSysprepScript(sysprepScript interface{}) *vmmConfig.OneOfSysprepSysprepScript {
	if len(sysprepScript.([]interface{})) > 0 {
		sysprepScriptObj := vmmConfig.NewOneOfSysprepSysprepScript()
		sysprepScriptData := sysprepScript.([]interface{})[0].(map[string]interface{})
		aJSON, _ := json.Marshal(sysprepScriptData)
		log.Printf("[DEBUG] sysprep.sysprep_script.sysprepScriptData: %s", string(aJSON))
		if unattendXML, ok := sysprepScriptData["unattend_xml"]; ok && len(unattendXML.([]interface{})) > 0 {
			unattendXMLObj := expandTemplateUnattendXML(unattendXML)
			aJSON, _ = json.Marshal(unattendXMLObj)
			log.Printf("[DEBUG] sysprep.sysprep_script.unattend_xml expanded: %v", string(aJSON))
			err := sysprepScriptObj.SetValue(*unattendXMLObj)
			if err != nil {
				log.Printf("[ERROR] SysprepScript: Error setting value for unattend Xml: %v", err)
				return nil
			}
		}
		if customKeyValues, ok := sysprepScriptData["custom_key_values"]; ok && len(customKeyValues.([]interface{})) > 0 {
			// customKeyValuesObj := vmmConfig.NewCustomKeyValues()
			customKeyValuesObj := expandTemplateCustomKeyValuesPairs(customKeyValues)
			aJSON, _ = json.Marshal(customKeyValuesObj)
			log.Printf("[DEBUG] sysprep.sysprep_script.customKeyValues expanded: %v", string(aJSON))
			err := sysprepScriptObj.SetValue(*customKeyValuesObj)
			if err != nil {
				log.Printf("[ERROR] SysprepScript: Error setting value for custom key values: %v", err)
				return nil
			}
		}

		return sysprepScriptObj
	}
	return nil
}

func expandTemplateCustomKeyValuesPairs(customKeyValues interface{}) *vmmConfig.CustomKeyValues {
	if customKeyValues != nil {
		customKeyValuesObj := vmmConfig.NewCustomKeyValues()
		customKeyValuesData := customKeyValues.([]interface{})
		if len(customKeyValuesData) > 0 {
			if keyValues := customKeyValuesData[0].(map[string]interface{})["key_value_pairs"]; keyValues != nil {
				log.Printf("[DEBUG] key_value_pairs: %v", keyValues)
				customKeyValuesObj.KeyValuePairs = expandTemplateKVPairs(keyValues)
			}
		}
		return customKeyValuesObj
	}
	return nil
}

func expandTemplateUnattendXML(unattendXML interface{}) *vmmConfig.Unattendxml {
	if unattendXML != nil {
		unattendXMLObj := vmmConfig.NewUnattendxml()
		unattendXMLData := unattendXML.([]interface{})

		if len(unattendXMLData) > 0 {
			if value, ok := unattendXMLData[0].(map[string]interface{})["value"]; ok {
				unattendXMLObj.Value = utils.StringPtr(value.(string))
			}
		}
		return unattendXMLObj
	}
	return nil
}

func expandVMConfigOverrideTemplate(pr interface{}) *vmmContent.VmConfigOverride {
	if len(pr.([]interface{})) > 0 {
		cfg := &vmmContent.VmConfigOverride{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if name, ok := val["name"]; ok && len(name.(string)) > 0 {
			cfg.Name = utils.StringPtr(name.(string))
		}
		if numSockets, ok := val["num_sockets"]; ok {
			cfg.NumSockets = utils.IntPtr(numSockets.(int))
		}
		if numCoresPerSocket, ok := val["num_cores_per_socket"]; ok {
			cfg.NumCoresPerSocket = utils.IntPtr(numCoresPerSocket.(int))
		}
		if numThreadsPerCore, ok := val["num_threads_per_core"]; ok {
			cfg.NumThreadsPerCore = utils.IntPtr(numThreadsPerCore.(int))
		}
		if memorySizeBytes, ok := val["memory_size_bytes"]; ok {
			cfg.MemorySizeBytes = utils.Int64Ptr(int64(memorySizeBytes.(int)))
		}
		if nics, ok := val["nics"]; ok && len(nics.([]interface{})) > 0 {
			cfg.Nics = expandNic(nics.([]interface{}))
		}
		if guest, ok := val["guest_customization"]; ok && len(guest.([]interface{})) > 0 {
			cfg.GuestCustomization = expandTemplateGuestCustomizationParams(guest)
		}
		return cfg
	}
	return nil
}

func expandGuestUpdateStatus(status interface{}) *vmmContent.GuestUpdateStatus {
	guestUpdateStatus := status.([]interface{})[0].(map[string]interface{})
	guestUpdateStatusObj := &vmmContent.GuestUpdateStatus{}
	if deployedVMReference, ok := guestUpdateStatus["deployed_vm_reference"]; ok {
		guestUpdateStatusObj.DeployedVmReference = utils.StringPtr(deployedVMReference.(string))
	}
	return guestUpdateStatusObj
}

func expandTemplateUser(user interface{}) *vmmContent.TemplateUser {
	userObj := &vmmContent.TemplateUser{}
	userData := user.([]interface{})[0].(map[string]interface{})

	if username, ok := userData["username"]; ok && username != "" {
		userObj.Username = utils.StringPtr(username.(string))
	}
	if userType, ok := userData["user_type"]; ok && userType != "" {
		userObj.UserType = expandUserType(userType)
	}
	if idpID, ok := userData["idp_id"]; ok && idpID != "" {
		userObj.IdpId = utils.StringPtr(idpID.(string))
	}
	if displayName, ok := userData["display_name"]; ok && displayName != "" {
		userObj.DisplayName = utils.StringPtr(displayName.(string))
	}
	if firstName, ok := userData["first_name"]; ok && firstName != "" {
		userObj.FirstName = utils.StringPtr(firstName.(string))
	}
	if middleInitial, ok := userData["middle_initial"]; ok && middleInitial != "" {
		userObj.MiddleInitial = utils.StringPtr(middleInitial.(string))
	}
	if lastName, ok := userData["last_name"]; ok && lastName != "" {
		userObj.LastName = utils.StringPtr(lastName.(string))
	}
	if emailID, ok := userData["email_id"]; ok && emailID != "" {
		userObj.EmailId = utils.StringPtr(emailID.(string))
	}
	if locale, ok := userData["locale"]; ok && locale != "" {
		userObj.Locale = utils.StringPtr(locale.(string))
	}
	if region, ok := userData["region"]; ok && region != "" {
		userObj.Region = utils.StringPtr(region.(string))
	}
	if password, ok := userData["password"]; ok && password != "" {
		userObj.Password = utils.StringPtr(password.(string))
	}
	if isForceResetPasswordEnabled, ok := userData["is_force_reset_password_enabled"]; ok {
		userObj.IsForceResetPasswordEnabled = utils.BoolPtr(isForceResetPasswordEnabled.(bool))
	}
	if additionalAttributes, ok := userData["additional_attributes"]; ok {
		userObj.AdditionalAttributes = expandTemplateKVPairs(additionalAttributes)
		aJSON, _ := json.Marshal(userObj.AdditionalAttributes)
		log.Printf("[DEBUG] expanede additionalAttributes: %v", string(aJSON))
		ad := flattenCustomKVPair(userObj.AdditionalAttributes)
		aJSON, _ = json.Marshal(ad)
		log.Printf("[DEBUG] Flattened additionalAttributes: %v", string(aJSON))
	}
	if status, ok := userData["status"]; ok {
		if status != nil && status != "" {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"ACTIVE":   two,
				"INACTIVE": three,
			}
			pVal := subMap[status.(string)]
			if pVal == nil {
				userObj.Status = nil
			}
			p := vmmAuthn.UserStatusType(pVal.(int))
			userObj.Status = &p
		}
	}
	if description, ok := userData["description"]; ok && description != "" {
		userObj.Description = utils.StringPtr(description.(string))
	}
	if creationType, ok := userData["creation_type"]; ok {
		if creationType != nil && creationType != "" {
			const two, three, four = 2, 3, 4
			subMap := map[string]interface{}{
				"PREDEFINED":     two,
				"USERDEFINED":    three,
				"SERVICEDEFINED": four,
			}
			pVal := subMap[creationType.(string)]
			if pVal == nil {
				userObj.CreationType = nil
			}
			p := vmmAuthn.CreationType(pVal.(int))
			userObj.CreationType = &p
		}
	}
	return userObj
}

func expandTemplateKVPairs(attributes interface{}) []vmmCommon.KVPair {
	var attributesList []vmmCommon.KVPair

	for _, attribute := range attributes.([]interface{}) {
		attributeData := attribute.(map[string]interface{})
		kvPair := vmmCommon.KVPair{}
		log.Printf("[DEBUG] attributeData: %v", attributeData)
		if attributeData["name"] != nil && attributeData["value"] != nil {
			kvPair.Name = utils.StringPtr(attributeData["name"].(string))
			kvPair.Value = expandValue(attributeData["value"])
			attributesList = append(attributesList, kvPair)
		}
	}
	aJSON, _ := json.Marshal(attributesList)
	log.Printf("[DEBUG] attributesList: %v", string(aJSON))
	return attributesList
}

func expandValue(kvPairValue interface{}) *vmmCommon.OneOfKVPairValue {
	valueObj := vmmCommon.NewOneOfKVPairValue()
	if kvPairValue != nil {
		valueData := kvPairValue.([]interface{})[0].(map[string]interface{})
		log.Printf("[DEBUG] kvPair valueData: %v", valueData)
		//nolint:gocritic // Keeping if-else for clarity in this specific case
		if valueData["string_list"] != nil && len(valueData["string_list"].([]interface{})) > 0 {
			log.Printf("[DEBUG] valueData of type string_list")
			stringList := valueData["string_list"].([]interface{})
			stringsListStr := make([]string, len(stringList))
			for i, v := range stringList {
				stringsListStr[i] = v.(string)
			}
			log.Printf("[DEBUG] stringsListStr: %v", stringsListStr)
			err := valueObj.SetValue(stringsListStr)
			if err != nil {
				log.Printf("[ERROR] Error setting value for string_list: %s", err)
				diag.Errorf("Error setting value for string_list: %s", err)
				return nil
			}
		} else if valueData["integer_list"] != nil && len(valueData["integer_list"].([]interface{})) > 0 {
			log.Printf("[DEBUG] valueData of type integer_list")
			integerList := valueData["integer_list"].([]interface{})
			integersListInt := make([]int, len(integerList))
			for i, v := range integerList {
				integersListInt[i] = v.(int)
			}
			err := valueObj.SetValue(integersListInt)
			if err != nil {
				log.Printf("[ERROR] Error setting value for integer_list: %s", err)
				diag.Errorf("Error setting value for integer_list: %s", err)
				return nil
			}
		} else if valueData["map_of_strings"] != nil && len(valueData["map_of_strings"].([]interface{})) > 0 {
			log.Printf("[DEBUG] valueData of type map_of_strings")
			mapOfStrings := make([]vmmCommon.MapOfStringWrapper, len(valueData["map_of_strings"].([]interface{})))

			for index, mapOfStringsData := range valueData["map_of_strings"].([]interface{}) {
				mapOfStringsDataMap := mapOfStringsData.(map[string]interface{})
				mapOfStringsObj := vmmCommon.MapOfStringWrapper{}
				mapOfStringsObj.Map = make(map[string]string)
				for k, v := range mapOfStringsDataMap["map"].(map[string]interface{}) {
					mapOfStringsObj.Map[k] = v.(string)
				}
				mapOfStrings[index] = mapOfStringsObj
			}
			aJSON, _ := json.Marshal(mapOfStrings)
			log.Printf("[DEBUG] mapOfStrings: %v", string(aJSON))
			err := valueObj.SetValue(mapOfStrings)
			if err != nil {
				log.Printf("[ERROR] Error setting value for map: %s", err)
				diag.Errorf("Error setting value for map: %s", err)
				return nil
			}
		} else if valueData["string"] != nil && valueData["string"] != "" {
			log.Printf("[DEBUG] valueData of type string")
			err := valueObj.SetValue(valueData["string"].(string))
			if err != nil {
				log.Printf("[ERROR] Error setting value for string: %s", err)
				diag.Errorf("Error setting value for string: %s", err)
				return nil
			}
		} else if valueData["object"] != nil && len(valueData["object"].(map[string]interface{})) > 0 {
			log.Printf("[DEBUG] valueData of type object")
			object := make(map[string]string)
			for k, v := range valueData["object"].(map[string]interface{}) {
				object[k] = v.(string)
			}
			err := valueObj.SetValue(object)
			if err != nil {
				log.Printf("[ERROR] Error setting value for object: %s", err)
				diag.Errorf("Error setting value for object: %s", err)
				return nil
			}
		} else if valueData["integer"] != nil && valueData["integer"] != 0 {
			log.Printf("[DEBUG] valueData of type integer")
			err := valueObj.SetValue(valueData["integer"].(int))
			if err != nil {
				log.Printf("[ERROR] Error setting value for integer: %s", err)
				diag.Errorf("Error setting value for integer: %s", err)
				return nil
			}
		} else if valueData["boolean"] != nil {
			log.Printf("[DEBUG] valueData of type boolean")
			err := valueObj.SetValue(valueData["boolean"].(bool))
			if err != nil {
				log.Printf("[ERROR] Error setting value for boolean: %s", err)
				diag.Errorf("Error setting value for boolean: %s", err)
				return nil
			}
		} else {
			log.Printf("[ERROR] invalid value type")
			return nil
		}
	}
	return valueObj
}

func expandUserType(userType interface{}) *vmmAuthn.UserType {
	if userType != nil && userType != "" {
		const zero, two, three, four, five, six = 0, 2, 3, 4, 5, 6
		subMap := map[string]interface{}{
			"LOCAL":           two,
			"SAML":            three,
			"LDAP":            four,
			"EXTERNAL":        five,
			"SERVICE_ACCOUNT": six,
		}
		pVal := subMap[userType.(string)]
		if pVal == nil {
			p := vmmAuthn.UserType(zero)
			return &p
		}
		p := vmmAuthn.UserType(pVal.(int))
		return &p
	}
	return nil
}
