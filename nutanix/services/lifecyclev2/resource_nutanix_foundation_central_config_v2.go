package lifecyclev2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fccfgimport "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/mgmt"
	lcmTaskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// foundationCentralConfigID is the synthetic identifier used for the Foundation
// Central configuration singleton. The API exposes a single configuration object
// with no external identifier, so a stable constant is used as the Terraform ID.
const foundationCentralConfigID = "foundation_central_config"

// ResourceNutanixFoundationCentralConfigV2 manages the configuration settings of
// the Foundation Central service. The Foundation Central APIs are served by the
// Foundation Central (FCVM) service, not Prism Central.
func ResourceNutanixFoundationCentralConfigV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixFoundationCentralConfigV2Create,
		ReadContext:   ResourceNutanixFoundationCentralConfigV2Read,
		UpdateContext: ResourceNutanixFoundationCentralConfigV2Update,
		DeleteContext: ResourceNutanixFoundationCentralConfigV2Delete,
		Schema: map[string]*schema.Schema{
			"ahv_installation_timeout_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Timeout in minutes for AHV installation",
			},
			"aos_download_timeout_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Timeout in minutes for AOS download",
			},
			"commit_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Foundation Central Commit ID",
			},
			"tls_certificate_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA-256 fingerprint of the TLS certificate used by Foundation Central service",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Foundation Central Version",
			},
		},
	}
}

// ResourceNutanixFoundationCentralConfigV2Create applies the desired Foundation
// Central configuration. The API has no create verb — the singleton is mutated via
// Update — so create performs an update and then refreshes state via Read.
func ResourceNutanixFoundationCentralConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if diags := updateFoundationCentralConfig(ctx, d, meta, schema.TimeoutCreate); diags != nil {
		return diags
	}

	d.SetId(foundationCentralConfigID)
	return ResourceNutanixFoundationCentralConfigV2Read(ctx, d, meta)
}

// ResourceNutanixFoundationCentralConfigV2Read fetches the current Foundation
// Central configuration and refreshes state.
func ResourceNutanixFoundationCentralConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.FoundationCentralConfigAPIInstance.GetFoundationCentralConfig()
	if err != nil {
		return diag.Errorf("error while fetching Foundation Central config: %v", err)
	}

	cfg := resp.Data.GetValue().(fccfgimport.FoundationCentralConfig)

	if err := d.Set("ahv_installation_timeout_minutes", utils.IntValue(cfg.AhvInstallationTimeoutMinutes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("aos_download_timeout_minutes", utils.IntValue(cfg.AosDownloadTimeoutMinutes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("commit_id", utils.StringValue(cfg.CommitId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tls_certificate_fingerprint", utils.StringValue(cfg.TlsCertificateFingerprint)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", utils.StringValue(cfg.Version)); err != nil {
		return diag.FromErr(err)
	}

	if d.Id() == "" {
		d.SetId(foundationCentralConfigID)
	}
	return nil
}

// ResourceNutanixFoundationCentralConfigV2Update applies configuration changes and
// refreshes state via Read.
func ResourceNutanixFoundationCentralConfigV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if diags := updateFoundationCentralConfig(ctx, d, meta, schema.TimeoutUpdate); diags != nil {
		return diags
	}

	return ResourceNutanixFoundationCentralConfigV2Read(ctx, d, meta)
}

// ResourceNutanixFoundationCentralConfigV2Delete is a no-op. The Foundation Central
// configuration is a singleton that cannot be deleted; removing the resource from
// state simply stops managing it.
func ResourceNutanixFoundationCentralConfigV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// updateFoundationCentralConfig performs a read-modify-write against the Foundation
// Central configuration: it fetches the current object, applies the changed fields,
// sends the update, and polls the resulting task to completion.
func updateFoundationCentralConfig(ctx context.Context, d *schema.ResourceData, meta interface{}, timeoutKey string) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	readResp, err := conn.FoundationCentralConfigAPIInstance.GetFoundationCentralConfig()
	if err != nil {
		return diag.Errorf("error while fetching Foundation Central config: %v", err)
	}

	// Extract E-Tag header for the conditional update.
	etagValue := conn.FoundationCentralConfigAPIInstance.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := readResp.Data.GetValue().(fccfgimport.FoundationCentralConfig)

	if v, ok := d.GetOk("ahv_installation_timeout_minutes"); ok {
		body.AhvInstallationTimeoutMinutes = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("aos_download_timeout_minutes"); ok {
		body.AosDownloadTimeoutMinutes = utils.IntPtr(v.(int))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Update Foundation Central Config Request Spec: %s", string(aJSON))

	resp, err := conn.FoundationCentralConfigAPIInstance.UpdateFoundationCentralConfig(&body, args)
	if err != nil {
		return diag.Errorf("error while updating Foundation Central config: %v", err)
	}

	taskRef := resp.Data.GetValue().(lcmTaskRef.TaskReference)
	taskUUID := taskRef.ExtId

	// Poll for completion of the update task using the shared prism task helper.
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(timeoutKey),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Foundation Central config (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Foundation Central config update task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Update Foundation Central Config Task Details: %s", string(aJSON))

	return nil
}
