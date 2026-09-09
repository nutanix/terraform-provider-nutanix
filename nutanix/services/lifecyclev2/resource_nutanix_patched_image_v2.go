package lifecyclev2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixPatchedImageV2 manages a patched hypervisor image created from a
// base image and customized for specific node deployments. A patched image
// references the claim token used during its creation (claim_token_ext_id),
// wiring the ClaimTokens cross-resource reference. Create and Delete are
// asynchronous task-based operations. The API does not support an update, so all
// input fields are ForceNew.
func ResourceNutanixPatchedImageV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixPatchedImageV2Create,
		ReadContext:   ResourceNutanixPatchedImageV2Read,
		DeleteContext: ResourceNutanixPatchedImageV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: patchedImageResourceSchema(),
	}
}

func patchedImageResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"claim_token_ext_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ExtId of the claim token used in the patched image",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Patched image name",
		},
		"host_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"AHV", "ESX"}, false),
			Description:  "Type of the host installed on the node",
		},
		"version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Version of the base hypervisor or host OS image used for patching",
		},
		"image_details": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "Details of the image used for patching",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"local_host_image_ext_id": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "External ID of the uploaded host image",
					},
				},
			},
		},
		"node_configurations": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			Description: "List of node configurations used for patching the image",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"node_ext_id": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "External ID of the node",
					},
					"host_configuration": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"hostname": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
									Description: "Hostname of the hypervisor or host OS",
								},
								"nameservers": {
									Type:        schema.TypeList,
									Optional:    true,
									ForceNew:    true,
									Description: "List of nameserver IP addresses to be used by the host",
									Elem:        patchedImageIPAddressResource(),
								},
								"ntpservers": {
									Type:        schema.TypeList,
									Optional:    true,
									ForceNew:    true,
									Description: "List of NTP server IP addresses to be used by the host",
									Elem:        patchedImageIPAddressOrFQDNResource(),
								},
								"network_details": {
									Type:     schema.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"management": {
												Type:     schema.TypeList,
												Optional: true,
												ForceNew: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"ip":      patchedImageIPAddressSchema(),
														"gateway": patchedImageIPAddressSchema(),
														"vlan_id": {
															Type:     schema.TypeInt,
															Optional: true,
															ForceNew: true,
														},
														"mtu_bytes": {
															Type:     schema.TypeInt,
															Optional: true,
															ForceNew: true,
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
		"patched_iso_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "URL for downloading the patched image",
		},
		"patched_iso_sha256_checksum": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "SHA-256 checksum of the patched image",
		},
		"created_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Timestamp when the patched image was created",
		},
		"owner_ext_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "External ID of the owner",
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "A globally unique identifier that represents the tenant that owns this entity.",
		},
		"links": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"href": {Type: schema.TypeString, Computed: true},
					"rel":  {Type: schema.TypeString, Computed: true},
				},
			},
		},
	}
}

func patchedImageIPAddressResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": patchedImageIPValuePrefix(),
			"ipv6": patchedImageIPValuePrefix(),
		},
	}
}

func patchedImageIPAddressOrFQDNResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": patchedImageIPValuePrefix(),
			"ipv6": patchedImageIPValuePrefix(),
			"fqdn": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {Type: schema.TypeString, Optional: true, ForceNew: true},
					},
				},
			},
		},
	}
}

func patchedImageIPAddressSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem:     patchedImageIPAddressResource(),
	}
}

func patchedImageIPValuePrefix() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value":         {Type: schema.TypeString, Optional: true, ForceNew: true},
				"prefix_length": {Type: schema.TypeInt, Optional: true, ForceNew: true},
			},
		},
	}
}

func ResourceNutanixPatchedImageV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	body := *import1.NewPatchedImage()
	if v, ok := d.GetOk("claim_token_ext_id"); ok {
		body.ClaimTokenExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("host_type"); ok {
		hostType := import1.HostType(hostTypeToEnum(v.(string)))
		body.HostType = &hostType
	}
	if v, ok := d.GetOk("version"); ok {
		body.Version = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("image_details"); ok {
		body.ImageDetails = expandImageDetails(v.([]interface{}))
	}
	if v, ok := d.GetOk("node_configurations"); ok {
		body.NodeConfigurations = expandNodeConfigurations(v.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Patched Image create payload : %s", string(aJSON))

	resp, err := conn.PatchedImagesAPIInstance.CreatePatchedImage(&body)
	if err != nil {
		return diag.Errorf("error while creating patched image : %v", err)
	}

	taskReference := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for patched image (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching patched image task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Patched Image Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypePatchedImage, "PatchedImage")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return ResourceNutanixPatchedImageV2Read(ctx, d, meta)
}

func ResourceNutanixPatchedImageV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.PatchedImagesAPIInstance.GetPatchedImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching patched image : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.PatchedImage)
	return flattenPatchedImageToState(d, &getResp)
}

func flattenPatchedImageToState(d *schema.ResourceData, getResp *import1.PatchedImage) diag.Diagnostics {
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("claim_token_ext_id", getResp.ClaimTokenExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if getResp.HostType != nil {
		if err := d.Set("host_type", getResp.HostType.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("image_details", flattenImageDetails(getResp.ImageDetails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_configurations", flattenNodeConfigurations(getResp.NodeConfigurations)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("patched_iso_url", getResp.PatchedIsoUrl); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("patched_iso_sha256_checksum", getResp.PatchedIsoSha256Checksum); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLifecycleLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixPatchedImageV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.PatchedImagesAPIInstance.DeletePatchedImageById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting patched image : %v", err)
	}

	taskReference := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for patched image (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

// hostTypeToEnum maps the string host type to the SDK HostType enum ordinal.
func hostTypeToEnum(v string) int {
	const ahv, esx = 2, 3
	switch v {
	case "AHV":
		return ahv
	case "ESX":
		return esx
	default:
		return 0
	}
}
