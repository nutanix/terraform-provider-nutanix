package lifecyclev2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/mgmt"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixFoundationCentralConfigV2 manages the configuration settings of the
// Foundation Central service. Foundation Central exposes a single, cluster-wide
// configuration object (no entity identifier), so this resource behaves as a
// singleton: Create applies the desired configuration via an Update call, Read
// fetches the current configuration, Update re-applies changes, and Delete is a
// no-op (the configuration object cannot be removed).
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

func ResourceNutanixFoundationCentralConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return updateFoundationCentralConfig(ctx, d, meta, schema.TimeoutCreate)
}

func ResourceNutanixFoundationCentralConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.FoundationCentralConfigAPIInstance.GetFoundationCentralConfig()
	if err != nil {
		return diag.Errorf("error while fetching foundation central config : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.FoundationCentralConfig)

	if err := d.Set("ahv_installation_timeout_minutes", getResp.AhvInstallationTimeoutMinutes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("aos_download_timeout_minutes", getResp.AosDownloadTimeoutMinutes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("commit_id", getResp.CommitId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tls_certificate_fingerprint", getResp.TlsCertificateFingerprint); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixFoundationCentralConfigV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return updateFoundationCentralConfig(ctx, d, meta, schema.TimeoutUpdate)
}

func ResourceNutanixFoundationCentralConfigV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Foundation Central exposes a singleton configuration object that cannot be
	// deleted. The framework clears the ID after this returns.
	return nil
}

// updateFoundationCentralConfig performs the read-modify-write required to apply the
// desired Foundation Central configuration. It fetches the current config, overlays
// the changed fields from the schema, sends the Update, and polls the resulting task.
func updateFoundationCentralConfig(ctx context.Context, d *schema.ResourceData, meta interface{}, timeout string) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	readResp, err := conn.FoundationCentralConfigAPIInstance.GetFoundationCentralConfig()
	if err != nil {
		return diag.Errorf("error while fetching foundation central config : %v", err)
	}

	body := readResp.Data.GetValue().(import1.FoundationCentralConfig)

	if v, ok := d.GetOk("ahv_installation_timeout_minutes"); ok {
		body.AhvInstallationTimeoutMinutes = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("aos_download_timeout_minutes"); ok {
		body.AosDownloadTimeoutMinutes = utils.IntPtr(v.(int))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Foundation Central Config Update Request Spec: %s", string(aJSON))

	resp, err := conn.FoundationCentralConfigAPIInstance.UpdateFoundationCentralConfig(&body)
	if err != nil {
		return diag.Errorf("error while updating foundation central config : %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(timeout),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for foundation central config (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching foundation central config update task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Foundation Central Config Update Task Details: %s", string(aJSON))

	// Singleton configuration object with no entity identifier — use a synthetic ID.
	d.SetId(resource.UniqueId())
	return ResourceNutanixFoundationCentralConfigV2Read(ctx, d, meta)
}
