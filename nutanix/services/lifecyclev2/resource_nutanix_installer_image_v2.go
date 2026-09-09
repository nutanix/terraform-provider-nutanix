package lifecyclev2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/response"
	config "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixInstallerImageV2 manages the lifecycle of an installer image registered
// with Foundation Central. An installer image bundles AOS and hypervisor software
// used for node provisioning. The v4 API is asynchronous (Create/Update/Delete
// return a TaskReference), so every handler polls the resulting task with the
// shared Prism task-group helper before reading the entity back.
func ResourceNutanixInstallerImageV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixInstallerImageV2Create,
		ReadContext:   ResourceNutanixInstallerImageV2Read,
		UpdateContext: ResourceNutanixInstallerImageV2Update,
		DeleteContext: ResourceNutanixInstallerImageV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the image. This field is required when creating an image.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AOS", "AHV", "ESX"}, false),
				Description:  "Type of the installer image. One of AOS, AHV or ESX.",
			},
			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LOCAL", "REMOTE_URL"}, false),
				Description:  "Source of the image. One of LOCAL or REMOTE_URL.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "URL from where the image can be downloaded. This can be provided only through the update operation and only for images with REMOTE_URL as the source type.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Version of the image.",
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Certificate chain for the image URL.",
			},
			"metadata_download_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "URL from where the image metadata can be downloaded for an AOS image.",
			},
			"checksum": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Checksum value of the image. This is required and applicable only for AHV images.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sha256": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "SHA-256 checksum of the image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SHA-256 checksum value in hexadecimal format (64 characters).",
									},
								},
							},
						},
						"md5": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "MD5 checksum of the image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "MD5 checksum value in hexadecimal format (32 characters).",
									},
								},
							},
						},
					},
				},
			},
			"file_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the image file on the cluster.",
			},
			"metadata_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the image metadata file on the cluster.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the image entity was created.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL at which the entity described by the link can be accessed.",
						},
						"rel": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A name that identifies the relationship of the link to the object that is returned by the URL.",
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixInstallerImageV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	inputSpec := config.Image{}

	if name, ok := d.GetOk("name"); ok {
		inputSpec.Name = utils.StringPtr(name.(string))
	}
	if imageType, ok := d.GetOk("type"); ok {
		inputSpec.Type = expandImageType(imageType.(string))
	}
	if source, ok := d.GetOk("source"); ok {
		inputSpec.Source = expandImageSource(source.(string))
	}
	// NOTE: `url` (ImageUrl) is intentionally NOT set on Create. The API only
	// accepts the image URL through the Update operation (and only for images
	// with REMOTE_URL as the source type). It is applied via an Update call
	// immediately after the create task succeeds (see below).
	if version, ok := d.GetOk("version"); ok {
		inputSpec.Version = utils.StringPtr(version.(string))
	}
	if certChain, ok := d.GetOk("certificate_chain"); ok {
		inputSpec.CertificateChain = utils.StringPtr(certChain.(string))
	}
	if metadataURL, ok := d.GetOk("metadata_download_url"); ok {
		inputSpec.MetadataDownloadUrl = utils.StringPtr(metadataURL.(string))
	}
	if checksum, ok := d.GetOk("checksum"); ok {
		inputSpec.Checksum = expandImageChecksum(checksum.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(inputSpec, "", " ")
	log.Printf("[DEBUG] Image create payload: %s", string(aJSON))

	resp, err := conn.InstallerImagesAPIInstance.CreateImage(&inputSpec)
	if err != nil {
		return diag.Errorf("error while creating image : %v", err)
	}

	taskReference := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching image task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Image Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeInstallerImage, "Image")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	// The image URL can only be set through the Update API. If the user supplied
	// a `url`, apply it now via a read-modify-write update so Create leaves the
	// resource in the fully-configured state.
	if url, ok := d.GetOk("url"); ok {
		if diags := applyInstallerImageURL(ctx, d, meta, url.(string)); diags != nil {
			return diags
		}
	}

	return ResourceNutanixInstallerImageV2Read(ctx, d, meta)
}

// applyInstallerImageURL performs a read-modify-write Update that sets the image
// URL, which the API only accepts through the Update operation. The resulting
// task is polled with the shared Prism task-group helper.
func applyInstallerImageURL(ctx context.Context, d *schema.ResourceData, meta interface{}, url string) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.InstallerImagesAPIInstance.GetImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching image : %v", err)
	}

	updateSpec := resp.Data.GetValue().(config.Image)
	updateSpec.Url = utils.StringPtr(url)

	updateResp, err := conn.InstallerImagesAPIInstance.UpdateImageById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while setting image url : %v", err)
	}

	taskReference := updateResp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image (%s) url update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func ResourceNutanixInstallerImageV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.InstallerImagesAPIInstance.GetImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching image : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Image)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if getResp.Type != nil {
		if err := d.Set("type", getResp.Type.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.Source != nil {
		if err := d.Set("source", getResp.Source.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("url", getResp.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("certificate_chain", getResp.CertificateChain); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata_download_url", getResp.MetadataDownloadUrl); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("checksum", flattenImageChecksum(getResp.Checksum)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.FileStatus != nil {
		if err := d.Set("file_status", getResp.FileStatus.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.MetadataStatus != nil {
		if err := d.Set("metadata_status", getResp.MetadataStatus.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenImageLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixInstallerImageV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.InstallerImagesAPIInstance.GetImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching image : %v", err)
	}

	updateSpec := resp.Data.GetValue().(config.Image)

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("type") {
		updateSpec.Type = expandImageType(d.Get("type").(string))
	}
	if d.HasChange("source") {
		updateSpec.Source = expandImageSource(d.Get("source").(string))
	}
	if d.HasChange("url") {
		updateSpec.Url = utils.StringPtr(d.Get("url").(string))
	}
	if d.HasChange("version") {
		updateSpec.Version = utils.StringPtr(d.Get("version").(string))
	}
	if d.HasChange("certificate_chain") {
		updateSpec.CertificateChain = utils.StringPtr(d.Get("certificate_chain").(string))
	}
	if d.HasChange("metadata_download_url") {
		updateSpec.MetadataDownloadUrl = utils.StringPtr(d.Get("metadata_download_url").(string))
	}
	if d.HasChange("checksum") {
		updateSpec.Checksum = expandImageChecksum(d.Get("checksum").([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", " ")
	log.Printf("[DEBUG] Image update payload: %s", string(aJSON))

	updateResp, err := conn.InstallerImagesAPIInstance.UpdateImageById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating image : %v", err)
	}

	taskReference := updateResp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixInstallerImageV2Read(ctx, d, meta)
}

func ResourceNutanixInstallerImageV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.InstallerImagesAPIInstance.DeleteImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting image : %v", err)
	}

	taskReference := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

// expandImageType converts the schema string into the SDK ImageType enum.
func expandImageType(imageType string) *config.ImageType {
	const two, three, four = 2, 3, 4
	subMap := map[string]int{
		"AOS": two,
		"AHV": three,
		"ESX": four,
	}
	if val, ok := subMap[imageType]; ok {
		p := config.ImageType(val)
		return &p
	}
	return nil
}

// expandImageSource converts the schema string into the SDK ImageSource enum.
func expandImageSource(source string) *config.ImageSource {
	const two, three = 2, 3
	subMap := map[string]int{
		"LOCAL":      two,
		"REMOTE_URL": three,
	}
	if val, ok := subMap[source]; ok {
		p := config.ImageSource(val)
		return &p
	}
	return nil
}

// expandImageChecksum builds the OneOfImageChecksum request value from the
// checksum schema block. Exactly one of sha256/md5 is expected.
func expandImageChecksum(checksums []interface{}) *config.OneOfImageChecksum {
	if len(checksums) == 0 || checksums[0] == nil {
		return nil
	}
	checksumMap := checksums[0].(map[string]interface{})
	oneOf := config.NewOneOfImageChecksum()

	if sha, ok := checksumMap["sha256"].([]interface{}); ok && len(sha) > 0 && sha[0] != nil {
		shaMap := sha[0].(map[string]interface{})
		sha256 := config.NewImageSha256Checksum()
		sha256.HexDigest = utils.StringPtr(shaMap["hex_digest"].(string))
		if err := oneOf.SetValue(*sha256); err != nil {
			log.Printf("[ERROR] error setting sha256 checksum: %v", err)
			return nil
		}
		return oneOf
	}

	if md5Raw, ok := checksumMap["md5"].([]interface{}); ok && len(md5Raw) > 0 && md5Raw[0] != nil {
		md5Map := md5Raw[0].(map[string]interface{})
		md5 := config.NewImageMd5Checksum()
		md5.HexDigest = utils.StringPtr(md5Map["hex_digest"].(string))
		if err := oneOf.SetValue(*md5); err != nil {
			log.Printf("[ERROR] error setting md5 checksum: %v", err)
			return nil
		}
		return oneOf
	}

	return nil
}

// flattenImageChecksum converts the OneOfImageChecksum response value back into
// the checksum schema block.
func flattenImageChecksum(checksum *config.OneOfImageChecksum) []map[string]interface{} {
	if checksum == nil || checksum.ObjectType_ == nil {
		return nil
	}

	value := checksum.GetValue()
	if value == nil {
		return nil
	}

	checksumMap := make(map[string]interface{})
	switch *checksum.ObjectType_ {
	case "lifecycle.v4.config.ImageSha256Checksum":
		if sha, ok := value.(config.ImageSha256Checksum); ok {
			checksumMap["sha256"] = []map[string]interface{}{
				{"hex_digest": utils.StringValue(sha.HexDigest)},
			}
		}
	case "lifecycle.v4.config.ImageMd5Checksum":
		if md5, ok := value.(config.ImageMd5Checksum); ok {
			checksumMap["md5"] = []map[string]interface{}{
				{"hex_digest": utils.StringValue(md5.HexDigest)},
			}
		}
	}

	return []map[string]interface{}{checksumMap}
}

// flattenImageLinks converts the SDK ApiLink slice into the links schema block.
func flattenImageLinks(links []import2.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	linkList := make([]map[string]interface{}, len(links))
	for i, link := range links {
		linkList[i] = map[string]interface{}{
			"href": utils.StringValue(link.Href),
			"rel":  utils.StringValue(link.Rel),
		}
	}
	return linkList
}
