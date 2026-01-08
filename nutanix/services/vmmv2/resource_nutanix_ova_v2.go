package vmmv2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import3 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/iam/v4/authn"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixOvaV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixOvaV2Create,
		ReadContext:   ResourceNutanixOvaV2Read,
		UpdateContext: ResourceNutanixOvaV2Update,
		DeleteContext: ResourceNutanixOvaV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"checksum": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1, //nolint:gomnd
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ova_sha1_checksum": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ExactlyOneOf: []string{ // Exactly one of the following fields must be set
								"checksum.0.ova_sha1_checksum",
								"checksum.0.ova_sha256_checksum",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"ova_sha256_checksum": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ExactlyOneOf: []string{ // Exactly one of the following fields must be set
								"checksum.0.ova_sha1_checksum",
								"checksum.0.ova_sha256_checksum",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ova_url_source": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"should_allow_insecure_url": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"basic_auth": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"username": {
													Type:     schema.TypeString,
													Required: true,
												},
												"password": {
													Type:      schema.TypeString,
													Required:  true,
													Sensitive: true,
												},
											},
										},
									},
								},
							},
						},
						"ova_vm_source": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_ext_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"disk_file_format": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"object_lite_source": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"created_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": {
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
						},
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
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
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
									"value": schemaForValue(),
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster_location_ext_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parent_vm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ResourceNutanixVirtualMachineV2(),
			},
			"disk_format": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"VMDK", "QCOW2"}, false),
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixOvaV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	body := &import1.Ova{}
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if checksum, ok := d.GetOk("checksum"); ok {
		body.Checksum = expandOneOfOvaChecksum(checksum)
	}
	if source, ok := d.GetOk("source"); ok {
		body.Source = expandOneOfOvaSource(source)
	}
	if clsExts, ok := d.GetOk("cluster_location_ext_ids"); ok {
		body.ClusterLocationExtIds = flattenStringValue(clsExts.([]interface{}))
	}

	var diskFormatMap = map[string]import1.OvaDiskFormat{
		"$UNKNOWN":  import1.OVADISKFORMAT_UNKNOWN,
		"$REDACTED": import1.OVADISKFORMAT_REDACTED,
		"QCOW2":     import1.OVADISKFORMAT_QCOW2,
		"VMDK":      import1.OVADISKFORMAT_VMDK,
	}
	if diskFormat, ok := d.GetOk("disk_format"); ok {
		if strValue, isString := diskFormat.(string); isString {
			if enumValue, exists := diskFormatMap[strValue]; exists {
				body.DiskFormat = &enumValue
			}
		}
	}
	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] OVA body: %s", string(aJSON))

	resp, err := conn.OvasAPIInstance.CreateOva(body)
	if err != nil {
		return diag.Errorf("error creating OVA: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import2.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the OVA to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for OVA (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching OVA create task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(import3.Task)

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] OVA Create Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeOVA, "OVA")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixOvaV2Read(ctx, d, meta)
}

func ResourceNutanixOvaV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	getResp, err := conn.OvasAPIInstance.GetOvaById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error reading OVA (%s): %v", d.Id(), err)
	}
	ova := getResp.Data.GetValue().(import1.Ova)
	aJSON, _ := json.MarshalIndent(ova, "", "  ")
	log.Printf("[DEBUG] Get Network call: %s", string(aJSON))
	if err := d.Set("ext_id", utils.StringValue(ova.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(ova.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenAPILink(ova.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", utils.StringValue(ova.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("checksum", flattenOneOfOvaChecksum(ova.Checksum)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("size_bytes", int(*ova.SizeBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", flattenCreatedBy(ova.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_location_ext_ids", ova.ClusterLocationExtIds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parent_vm", utils.StringValue(ova.ParentVm)); err != nil {
		return diag.FromErr(err)
	}

	// Set the VM config
	fields, diags := extractVMConfigFields(*ova.VmConfig)
	if diags.HasError() {
		return diags
	}
	if err := d.Set("vm_config", []interface{}{fields}); err != nil {
		return diag.FromErr(fmt.Errorf("failed setting vm_config: %w", err))
	}

	if err := d.Set("disk_format", flattenOvaDiskFormat(ova.DiskFormat)); err != nil {
		return diag.FromErr(err)
	}

	if ova.CreateTime != nil {
		t := ova.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if ova.LastUpdateTime != nil {
		t := ova.LastUpdateTime
		if err := d.Set("last_update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] OVA (%s) read successfully", d.Id())
	return nil
}

func ResourceNutanixOvaV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	getResp, err := conn.OvasAPIInstance.GetOvaById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error reading OVA (%s): %v", d.Id(), err)
	}

	respOvas := getResp.Data.GetValue().(import1.Ova)
	updateSpec := respOvas

	// Only update of name is allowed
	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			updateSpec.Name = utils.StringPtr(v.(string))
		}
	}

	// remove created by from the request as it cause the request to fail
	// with error -> [Path '/createdBy'] Object has missing required properties ([\"userType\"])
	updateSpec.CreatedBy = nil

	// Prepare the update request
	updateResp, err := conn.OvasAPIInstance.UpdateOvaById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error updating OVA (%s): %v", d.Id(), err)
	}

	TaskRef := updateResp.Data.GetValue().(import2.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the OVA to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for OVA (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixOvaV2Read(ctx, d, meta)
}

func ResourceNutanixOvaV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	// Delete the OVA resource
	if d.Id() == "" {
		return diag.FromErr(errors.New("resource id is empty, cannot delete ova"))
	}

	deleteResp, err := conn.OvasAPIInstance.DeleteOvaById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error deleting OVA (%s): %v", d.Id(), err)
	}

	TaskRef := deleteResp.Data.GetValue().(import2.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the OVA to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for OVA (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId("")
	return nil
}

func expandOneOfOvaChecksum(pr interface{}) *import1.OneOfOvaChecksum {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		chksum := &import1.OneOfOvaChecksum{}

		if val["ova_sha1_checksum"] != nil && len(val["ova_sha1_checksum"].([]interface{})) > 0 {
			hexDigestI := val["ova_sha1_checksum"].([]interface{})
			hexDigestVal := hexDigestI[0].(map[string]interface{})
			sha1 := import1.NewOvaSha1Checksum()

			sha1.HexDigest = utils.StringPtr(hexDigestVal["hex_digest"].(string))
			chksum.SetValue(*sha1)
		} else if val["ova_sha256_checksum"] != nil && len(val["ova_sha256_checksum"].([]interface{})) > 0 {
			hexDigestI := val["ova_sha256_checksum"].([]interface{})
			hexDigestVal := hexDigestI[0].(map[string]interface{})

			sha256 := import1.NewOvaSha256Checksum()
			sha256.HexDigest = utils.StringPtr(hexDigestVal["hex_digest"].(string))
			chksum.SetValue(*sha256)
		}
		return chksum
	}
	return nil
}

func flattenOneOfOvaChecksum(checksum *import1.OneOfOvaChecksum) []map[string]interface{} {
	if checksum != nil {
		resList := make([]map[string]interface{}, 0)

		sha := make(map[string]interface{})

		getVal := checksum.ObjectType_

		if utils.StringValue(getVal) == "vmm.v4.content.OvaSha1Checksum" {
			sha1 := checksum.GetValue().(import1.OvaSha1Checksum)

			sha1List := make([]map[string]interface{}, 0)
			sha1Map := make(map[string]interface{})
			sha1Map["hex_digest"] = sha1.HexDigest
			sha1List = append(sha1List, sha1Map)
			sha["ova_sha1_checksum"] = sha1List
		} else {
			sha256 := checksum.GetValue().(import1.OvaSha256Checksum)

			sha256List := make([]map[string]interface{}, 0)
			sha256Map := make(map[string]interface{})
			sha256Map["hex_digest"] = sha256.HexDigest
			sha256List = append(sha256List, sha256Map)
			sha["ova_sha256_checksum"] = sha256List
		}
		resList = append(resList, sha)
		return resList
	}
	return nil
}

func expandOneOfOvaSource(pr interface{}) *import1.OneOfOvaSource {
	if pr != nil && len(pr.([]interface{})) > 0 {
		imgSrc := &import1.OneOfOvaSource{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		if urlSrc, ok := val["ova_url_source"]; ok && len(urlSrc.([]interface{})) > 0 {
			urlSrcInput := import1.OvaUrlSource{}

			urlIn := urlSrc.([]interface{})
			urlMap := urlIn[0].(map[string]interface{})

			if url, ok := urlMap["url"]; ok && len(url.(string)) > 0 {
				urlSrcInput.Url = utils.StringPtr(url.(string))
			}
			if basicAuth, ok := urlMap["basic_auth"]; ok && len(basicAuth.([]interface{})) > 0 {
				urlSrcInput.BasicAuth = expandOvaURLBasicAuth(basicAuth)
			}
			urlSrcInput.ObjectType_ = utils.StringPtr("vmm.v4.content.OvaUrlSource")
			err := imgSrc.SetValue(urlSrcInput)
			if err != nil {
				log.Fatalf("SetValue failed: %v", err)
			}
		}

		if vmDisk, ok := val["ova_vm_source"]; ok && len(vmDisk.([]interface{})) > 0 {
			OvavmDiskSrc := import1.OvaVmSource{}

			vmDiskIn := vmDisk.([]interface{})
			vmDiskMap := vmDiskIn[0].(map[string]interface{})

			OvavmDiskSrc.VmExtId = utils.StringPtr(vmDiskMap["vm_ext_id"].(string))
			OvavmDiskSrc.ObjectType_ = utils.StringPtr("vmm.v4.content.OvaVmSource")
			if diskFormat, ok := vmDiskMap["disk_file_format"]; ok && len(diskFormat.(string)) > 0 {
				var diskFormatMap = map[string]import1.OvaDiskFormat{
					"$UNKNOWN":  import1.OVADISKFORMAT_UNKNOWN,
					"$REDACTED": import1.OVADISKFORMAT_REDACTED,
					"QCOW2":     import1.OVADISKFORMAT_QCOW2,
					"VMDK":      import1.OVADISKFORMAT_VMDK,
				}
				if enumValue, exists := diskFormatMap[diskFormat.(string)]; exists {
					OvavmDiskSrc.DiskFileFormat = &enumValue
				}
			}
			err := imgSrc.SetValue(OvavmDiskSrc)
			if err != nil {
				log.Fatalf("SetValue failed: %v", err)
			}
		}

		if objLite, ok := val["object_lite_source"]; ok && len(objLite.([]interface{})) > 0 {
			objLiteIn := objLite.([]interface{})
			objLiteMap := objLiteIn[0].(map[string]interface{})
			objLiteSrc := import1.ObjectsLiteSource{}
			objLiteSrc.ObjectType_ = utils.StringPtr("vmm.v4.content.ObjectsLiteSource")
			objLiteSrc.Key = utils.StringPtr(objLiteMap["key"].(string))
			err := imgSrc.SetValue(objLiteSrc)
			if err != nil {
				log.Fatalf("SetValue failed: %v", err)
			}
		}
		return imgSrc
	}
	return nil
}

func flattenOneOfOvaSource(source *import1.OneOfOvaSource) []map[string]interface{} {
	if source != nil {
		resList := make([]map[string]interface{}, 0)

		urlSrcMap := make(map[string]interface{})
		urlSrcList := make([]map[string]interface{}, 0)

		vmDiskSrcMap := make(map[string]interface{})
		vmDiskSrcList := make([]map[string]interface{}, 0)

		objLiteSrc := make(map[string]interface{})
		objLiteSrcList := make([]map[string]interface{}, 0)

		if utils.StringValue(source.ObjectType_) == "vmm.v4.content.OvaUrlSource" {
			urlSrc := source.GetValue().(import1.UrlSource)

			urlSrcObj := make(map[string]interface{})
			urlSrcObjList := make([]map[string]interface{}, 0)

			if urlSrc.Url != nil {
				urlSrcObj["url"] = urlSrc.Url
			}
			if urlSrc.BasicAuth != nil {
				urlSrcObj["basic_auth"] = flattenURLBasicAuth(urlSrc.BasicAuth)
			}

			urlSrcObjList = append(urlSrcObjList, urlSrcObj)

			urlSrcMap["ova_url_source"] = urlSrcObjList
			urlSrcList = append(urlSrcList, urlSrcMap)

			return urlSrcList
		}

		if utils.StringValue(source.ObjectType_) == "vmm.v4.content.OvaVmSource" {
			vmDiskSrc := source.GetValue().(import1.VmDiskSource)

			vmDiskObj := make(map[string]interface{})
			vmDiskObjList := make([]map[string]interface{}, 0)

			if vmDiskSrc.ExtId != nil {
				vmDiskObj["ext_id"] = vmDiskSrc.ExtId
			}

			vmDiskObjList = append(vmDiskObjList, vmDiskObj)

			vmDiskSrcMap["ova_vm_source"] = vmDiskObjList

			vmDiskSrcList = append(vmDiskSrcList, vmDiskSrcMap)
			return vmDiskSrcList
		}

		if utils.StringValue(source.ObjectType_) == "vmm.v4.content.ObjectsLiteSource" {
			objLite := source.GetValue().(import1.ObjectsLiteSource)

			objLiteMap := make(map[string]interface{})
			objLiteMapList := make([]map[string]interface{}, 0)
			if objLite.Key != nil {
				objLiteMap["key"] = objLite.Key
			}
			objLiteMapList = append(objLiteMapList, objLiteMap)
			objLiteSrc["object_lite_source"] = objLiteMapList
			objLiteSrcList = append(objLiteSrcList, objLiteSrc)
			return objLiteSrcList
		}
		resList = append(resList, urlSrcMap)
		resList = append(resList, vmDiskSrcMap)
		resList = append(resList, objLiteSrc)
		return resList
	}
	return nil
}

func expandOvaURLBasicAuth(pr interface{}) *import1.UrlBasicAuth {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		basicAuths := &import1.UrlBasicAuth{}

		if user, ok := val["username"]; ok {
			basicAuths.Username = utils.StringPtr(user.(string))
		}
		if pass, ok := val["password"]; ok {
			basicAuths.Password = utils.StringPtr(pass.(string))
		}
		return basicAuths
	}
	return nil
}

func flattenOvaDiskFormat(diskFormat *import1.OvaDiskFormat) string {
	if diskFormat != nil {
		switch *diskFormat {
		case import1.OVADISKFORMAT_QCOW2:
			return "QCOW2"
		case import1.OVADISKFORMAT_VMDK:
			return "VMDK"
		default:
			return ""
		}
	}
	return ""
}

func flattenCreatedBy(createdBy *import4.User) []map[string]interface{} {
	if createdBy != nil {
		resList := make([]map[string]interface{}, 0)
		createdByMap := make(map[string]interface{})
		if v := utils.StringValue(createdBy.TenantId); v != "" {
			createdByMap["tenant_id"] = v
		}
		if v := createdBy.Links; len(v) > 0 {
			createdByMap["links"] = flattenAPILink(v)
		}
		if v := utils.StringValue(createdBy.ExtId); v != "" {
			createdByMap["ext_id"] = v
		}
		if v := utils.StringValue(createdBy.Username); v != "" {
			createdByMap["username"] = v
		}
		if v := flattenUserType(createdBy.UserType); v != "" {
			createdByMap["user_type"] = v
		}
		if v := utils.StringValue(createdBy.IdpId); v != "" {
			createdByMap["idp_id"] = v
		}
		if v := utils.StringValue(createdBy.DisplayName); v != "" {
			createdByMap["display_name"] = v
		}
		if v := utils.StringValue(createdBy.FirstName); v != "" {
			createdByMap["first_name"] = v
		}
		if v := utils.StringValue(createdBy.MiddleInitial); v != "" {
			createdByMap["middle_initial"] = v
		}
		if v := utils.StringValue(createdBy.LastName); v != "" {
			createdByMap["last_name"] = v
		}
		if v := utils.StringValue(createdBy.EmailId); v != "" {
			createdByMap["email_id"] = v
		}
		if v := utils.StringValue(createdBy.Locale); v != "" {
			createdByMap["locale"] = v
		}
		if v := utils.StringValue(createdBy.Region); v != "" {
			createdByMap["region"] = v
		}
		if v := utils.StringValue(createdBy.Password); v != "" {
			createdByMap["password"] = v
		}
		if createdBy.IsForceResetPasswordEnabled != nil {
			createdByMap["is_force_reset_password_enabled"] = *createdBy.IsForceResetPasswordEnabled
		}
		if len(createdBy.AdditionalAttributes) > 0 {
			createdByMap["additional_attributes"] = flattenCustomKVPair(createdBy.AdditionalAttributes)
		}
		if v := flattenUserStatusType(createdBy.Status); v != "" {
			createdByMap["status"] = v
		}
		if len(createdBy.BucketsAccessKeys) > 0 {
			createdByMap["buckets_access_keys"] = flattenBucketsAccessKey(createdBy.BucketsAccessKeys)
		}
		if createdBy.LastLoginTime != nil {
			createdByMap["last_login_time"] = createdBy.LastLoginTime.String()
		}
		if createdBy.CreatedTime != nil {
			createdByMap["created_time"] = createdBy.CreatedTime.String()
		}
		if createdBy.LastUpdatedTime != nil {
			createdByMap["last_updated_time"] = createdBy.LastUpdatedTime.String()
		}
		if v := utils.StringValue(createdBy.CreatedBy); v != "" {
			createdByMap["created_by"] = v
		}
		if v := utils.StringValue(createdBy.Description); v != "" {
			createdByMap["description"] = v
		}
		if v := flattenCreationType(createdBy.CreationType); v != "" {
			createdByMap["creation_type"] = v
		}
		resList = append(resList, createdByMap)
		return resList
	}
	return nil
}

func flattenCreationType(pr *import4.CreationType) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == import4.CreationType(two) {
			return "PREDEFINED"
		}
		if *pr == import4.CreationType(three) {
			return "USERDEFINED"
		}
		if *pr == import4.CreationType(four) {
			return "SERVICEDEFINED"
		}
	}
	return "UNKNOWN"
}
