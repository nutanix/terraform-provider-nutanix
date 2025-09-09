package vmmv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	import5 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixImageV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixImageV4Create,
		ReadContext:   ResourceNutanixImageV4Read,
		UpdateContext: ResourceNutanixImageV4Update,
		DeleteContext: ResourceNutanixImageV4Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			sources := []string{"source.0.url_source", "source.0.vm_disk_source", "source.0.object_lite_source"}
			count := 0
			for _, s := range sources {
				if _, ok := d.GetOk(s); ok {
					count++
				}
			}
			if count > 1 {
				return fmt.Errorf("only one of url_source, vm_disk_source, or object_lite_source can be specified in source")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISK_IMAGE", "ISO_IMAGE"}, false),
			},
			"checksum": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hex_digest": {
							Type:     schema.TypeString,
							Required: true,
						},
						"object_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"sha1", "sha256"}, false),
						},
					},
				},
			},
			"source": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url_source": {
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
										Default:  false,
									},
									"basic_auth": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"username": {
													Type:     schema.TypeString,
													Required: true,
												},
												"password": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"vm_disk_source": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
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
			"category_ext_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"size_bytes": {
				Type:     schema.TypeInt,
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
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"placement_policy_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"placement_policy_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"compliance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enforcement_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_cluster_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enforced_cluster_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"conflicting_policy_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixImageV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	body := &import5.Image{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if types, ok := d.GetOk("type"); ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"DISK_IMAGE": two,
			"ISO_IMAGE":  three,
		}
		pVal := subMap[types.(string)]
		p := import5.ImageType(pVal.(int))
		body.Type = &p
	}
	if checksum, ok := d.GetOk("checksum"); ok {
		body.Checksum = expandOneOfImageChecksum(checksum)
	}
	if src, ok := d.GetOk("source"); ok {
		body.Source = expandOneOfImageSource(src)
	}
	if ctgExts, ok := d.GetOk("category_ext_ids"); ok {
		body.CategoryExtIds = flattenStringValue(ctgExts.([]interface{}))
	}
	if clsExts, ok := d.GetOk("cluster_location_ext_ids"); ok {
		body.ClusterLocationExtIds = flattenStringValue(clsExts.([]interface{}))
	}

	resp, err := conn.ImagesAPIInstance.CreateImage(body)
	if err != nil {
		return diag.Errorf("error while creating Image : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Image UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)
	return ResourceNutanixImageV4Read(ctx, d, meta)
}

func ResourceNutanixImageV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.ImagesAPIInstance.GetImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching images : %v", err)
	}

	getResp := resp.Data.GetValue().(import5.Image)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenAPILink(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", flattenImageType(getResp.Type)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("checksum", flattenOneOfImageChecksum(getResp.Checksum)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("size_bytes", getResp.SizeBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source", flattenOneOfImageSource(getResp.Source)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_ext_ids", getResp.CategoryExtIds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_location_ext_ids", getResp.ClusterLocationExtIds); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreateTime != nil {
		t := getResp.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdateTime != nil {
		t := getResp.LastUpdateTime
		if err := d.Set("last_update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("placement_policy_status", flattenImagePlacementStatus(getResp.PlacementPolicyStatus)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixImageV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.ImagesAPIInstance.GetImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching images : %v", err)
	}

	respImages := resp.Data.GetValue().(import5.Image)
	updateSpec := respImages

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("type") {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"DISK_IMAGE": two,
			"ISO_IMAGE":  three,
		}
		pVal := subMap[d.Get("type").(string)]
		p := import5.ImageType(pVal.(int))
		updateSpec.Type = &p
	}
	if d.HasChange("checksum") {
		updateSpec.Checksum = expandOneOfImageChecksum(d.Get("checksum"))
	}
	if d.HasChange("source") {
		updateSpec.Source = expandOneOfImageSource(d.Get("source"))
	}
	if d.HasChange("category_ext_ids") {
		updateSpec.CategoryExtIds = flattenStringValue(d.Get("category_ext_ids").([]interface{}))
	}
	if d.HasChange("cluster_location_ext_ids") {
		updateSpec.ClusterLocationExtIds = flattenStringValue(d.Get("cluster_location_ext_ids").([]interface{}))
	}

	updateResp, er := conn.ImagesAPIInstance.UpdateImageById(utils.StringPtr(d.Id()), &updateSpec)
	if er != nil {
		return diag.Errorf("error while updating images : %v", err)
	}
	TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return ResourceNutanixImageV4Read(ctx, d, meta)
}

func ResourceNutanixImageV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.ImagesAPIInstance.DeleteImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting images : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func taskStateRefreshPrismTaskGroupFunc(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)
		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(import2.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

func getTaskStatus(pr *import2.TaskStatus) string {
	if pr != nil {
		const two, three, five, six, seven = 2, 3, 5, 6, 7
		if *pr == import2.TaskStatus(six) {
			return "FAILED"
		}
		if *pr == import2.TaskStatus(seven) {
			return "CANCELED"
		}
		if *pr == import2.TaskStatus(two) {
			return "QUEUED"
		}
		if *pr == import2.TaskStatus(three) {
			return "RUNNING"
		}
		if *pr == import2.TaskStatus(five) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}

func expandOneOfImageChecksum(pr interface{}) *import5.OneOfImageChecksum {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		chksum := &import5.OneOfImageChecksum{}

		if val["object_type"] == "sha1" {
			sha1 := chksum.GetValue().(import5.ImageSha1Checksum)

			sha1.HexDigest = utils.StringPtr(val["hex_digest"].(string))
			chksum.SetValue(sha1)
		} else {
			sha256 := chksum.GetValue().(import5.ImageSha256Checksum)
			sha256.HexDigest = utils.StringPtr(val["hex_digest"].(string))
			chksum.SetValue(sha256)
		}
		return chksum
	}
	return nil
}

func expandOneOfImageSource(pr interface{}) *import5.OneOfImageSource {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		imgSrc := &import5.OneOfImageSource{}

		if urlSrc, ok := val["url_source"]; ok && len(urlSrc.([]interface{})) > 0 {
			urlSrcInput := import5.NewUrlSource()

			urlIn := urlSrc.([]interface{})
			urlMap := urlIn[0].(map[string]interface{})

			if url, ok := urlMap["url"]; ok && len(url.(string)) > 0 {
				urlSrcInput.Url = utils.StringPtr(url.(string))
			}
			if shouldAllow, ok := urlMap["should_allow_insecure_url"]; ok {
				urlSrcInput.ShouldAllowInsecureUrl = utils.BoolPtr(shouldAllow.(bool))
			}
			if basicAuth, ok := urlMap["basic_auth"]; ok && len(basicAuth.([]interface{})) > 0 {
				urlSrcInput.BasicAuth = expandURLBasicAuth(basicAuth)
			}
			imgSrc.SetValue(*urlSrcInput)
		}

		if vmDisk, ok := val["vm_disk_source"]; ok && len(vmDisk.([]interface{})) > 0 {
			vmDiskSrc := import5.NewVmDiskSource()

			vmDiskIn := vmDisk.([]interface{})
			vmDiskMap := vmDiskIn[0].(map[string]interface{})

			vmDiskSrc.ExtId = utils.StringPtr(vmDiskMap["ext_id"].(string))
			imgSrc.SetValue(*vmDiskSrc)
		}

		if objLite, ok := val["object_lite_source"]; ok && len(objLite.([]interface{})) > 0 {
			objLiteIn := objLite.([]interface{})
			objLiteMap := objLiteIn[0].(map[string]interface{})

			objLiteSrc := import5.NewObjectsLiteSource()

			objLiteSrc.Key = utils.StringPtr(objLiteMap["key"].(string))
			imgSrc.SetValue(*objLiteSrc)
		}
		return imgSrc
	}
	return nil
}

func expandURLBasicAuth(pr interface{}) *import5.UrlBasicAuth {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		basicAuths := &import5.UrlBasicAuth{}

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

func flattenStringValue(pr []interface{}) []string {
	if len(pr) == 0 {
		return []string{} // return empty slice, not nil
	}

	res := make([]string, 0, len(pr))
	for _, v := range pr {
		str, ok := v.(string)
		if !ok {
			// handle the error gracefully — maybe skip or log?
			continue
		}
		res = append(res, str)
	}
	return res
}
