package microsegv2

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	microsegConfig "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/request/networksecuritypolicies"
	prismMicroseg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/microseg"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixNetworkSecurityPolicyExportV2 exports a point-in-time snapshot of the
// supplied network security policies.
func ResourceNutanixNetworkSecurityPolicyExportV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNetworkSecurityPolicyExportV2Create,
		ReadContext:   resourceNutanixNetworkSecurityPolicyExportV2Read,
		// No UpdateContext: every field is ForceNew or Computed, so Terraform always
		// destroys and recreates the resource instead of updating it in place.
		DeleteContext: resourceNutanixNetworkSecurityPolicyExportV2Delete,
		Schema: map[string]*schema.Schema{
			"policy_ext_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The list of network security policy external identifiers (UUIDs) to export. If omitted, all network security policies are exported. May be combined with project_ext_id to export only the listed policies within a project context.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_ext_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The external identifier of the project associated with the policies being exported.",
			},
			"file_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Local file path where the raw exported payload is written. The resulting file can be fed directly to the `path` argument of nutanix_network_security_policy_import_v2.",
			},
			"exported_payload": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The exported network security policy payload returned by the cluster, base64-encoded (the raw payload is a binary octet-stream).",
			},
			"task_ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier for the task created by the export operation.",
			},
		},
	}
}

func resourceNutanixNetworkSecurityPolicyExportV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	spec := microsegConfig.NewNetworkSecurityPolicyExportSpec()
	// When policy_ext_ids is omitted, PolicyReferences is left empty so the cluster
	// exports every network security policy.
	if v, ok := d.GetOk("policy_ext_ids"); ok {
		spec.PolicyReferences = common.ExpandListOfString(common.InterfaceToSlice(v))
	}
	if v, ok := d.GetOk("project_ext_id"); ok {
		spec.ProjectExtId = utils.StringPtr(v.(string))
	}

	req := import1.ExportNetworkSecurityPolicyRequest{
		Body: spec,
	}

	aJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("[DEBUG] ExportNetworkSecurityPolicy Request: %s", string(aJSON))

	taskConn := meta.(*conns.Client).PrismAPI

	// The microseg backend serializes policy import/export/create/delete/update
	// operations: only one may run at a time cluster-wide. A create/delete task can
	// report SUCCEEDED while the backend briefly still holds that lock, so an export
	// issued immediately afterwards is rejected with a "concurrent operations are not
	// permitted" error. Retry the prepare-export (steps 1-2) until the lock clears.
	//
	// requestID is captured from the successful attempt because a single NTNX-Request-Id
	// correlates the asynchronous prepare-export with the subsequent synchronous download;
	// the cluster keys the staged export on this id and the SDK auto-generates a fresh id
	// per call, so it must be supplied explicitly to both calls.
	var requestID string
	var taskDetails prismConfig.Task
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error generating request id for network security policy export: %v", err))
		}
		requestID = id
		requestHeader := map[string]interface{}{
			"NTNX-Request-Id": utils.StringPtr(requestID),
		}

		// Step 1: trigger the asynchronous prepare-export task.
		resp, err := conn.NetworkingSecurityInstance.ExportNetworkSecurityPolicy(ctx, &req, requestHeader)
		if err != nil {
			if isConcurrentOperationErr(err) {
				log.Printf("[DEBUG] network security policy export blocked by a concurrent operation, retrying: %v", err)
				return resource.RetryableError(fmt.Errorf("error exporting network security policies: %v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("error exporting network security policies: %v", err))
		}

		taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("unexpected response type from ExportNetworkSecurityPolicy"))
		}
		taskUUID := taskRef.ExtId

		// Step 2: poll the task until it succeeds.
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskConn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutCreate),
		}
		if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
			if isConcurrentOperationErr(errWait) {
				log.Printf("[DEBUG] network security policy export task blocked by a concurrent operation, retrying: %v", errWait)
				return resource.RetryableError(fmt.Errorf("error waiting for network security policy export task (%s): %s", utils.StringValue(taskUUID), errWait))
			}
			return resource.NonRetryableError(fmt.Errorf("error waiting for network security policy export task (%s): %s", utils.StringValue(taskUUID), errWait))
		}

		getTaskByIDRequest := import2.GetTaskByIdRequest{
			ExtId: taskUUID,
		}
		taskResp, err := taskConn.TaskRefAPI.GetTaskById(ctx, &getTaskByIDRequest)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error fetching network security policy export task (%s): %v", utils.StringValue(taskUUID), err))
		}
		taskDetails = taskResp.Data.GetValue().(prismConfig.Task)
		return nil
	})
	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] ExportNetworkSecurityPolicy Task Details: %s", string(aJSON))

	// Step 3: download the serialized payload that the prepare-export task staged, reusing
	// the same request id so the cluster returns this export rather than another one.
	// The payload is a binary octet-stream, so it is kept as raw bytes here.
	payload, err := downloadExportedPayload(ctx, conn, requestID)
	if err != nil {
		return diag.Errorf("error retrieving exported network security policy payload: %v", err)
	}

	// Persist the raw payload to disk so it can be consumed directly by the import
	// resource (whose `path` argument expects this exact octet-stream format).
	filePath := d.Get("file_path").(string)
	if err := os.WriteFile(filePath, payload, 0o600); err != nil {
		return diag.Errorf("error writing exported payload to file (%s): %v", filePath, err)
	}

	// Step 4: persist results. The binary payload is base64-encoded so it is a valid,
	// round-trippable Terraform string value.
	if err := d.Set("exported_payload", base64.StdEncoding.EncodeToString(payload)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("task_ext_id", utils.StringValue(taskDetails.ExtId)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(taskDetails.ExtId))

	return nil
}

func resourceNutanixNetworkSecurityPolicyExportV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	policyExtIDs := common.ExpandListOfString(common.InterfaceToSlice(d.Get("policy_ext_ids")))
	if len(policyExtIDs) == 0 {
		return nil
	}

	// Verify the source policies still exist on the cluster. If every policy referenced by
	// this export has been deleted out-of-band, the snapshot no longer maps to anything, so
	// we drop it from state to allow a clean re-export.
	existing := 0
	for _, extID := range policyExtIDs {
		getReq := import1.GetNetworkSecurityPolicyByIdRequest{
			ExtId: utils.StringPtr(extID),
		}
		if _, err := conn.NetworkingSecurityInstance.GetNetworkSecurityPolicyById(ctx, &getReq); err != nil {
			log.Printf("[DEBUG] network security policy (%s) referenced by export not found: %v", extID, err)
			continue
		}
		existing++
	}

	if existing == 0 {
		log.Printf("[WARN] all source policies for export (%s) were deleted; removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

func resourceNutanixNetworkSecurityPolicyExportV2Delete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No-Op: destroying an export only removes the snapshot from Terraform state.
	return nil
}

// downloadExportedPayload performs the synchronous GET that flushes the prepared export from
// the cluster. The SDK downloads the octet-stream body to a temporary file and returns its
// path via FileDetail; we read that file back to capture the payload contents.
func downloadExportedPayload(ctx context.Context, conn *microseg.Client, requestID string) ([]byte, error) {
	listReq := import1.ListNetworkSecurityPoliciesRequest{}

	// Force the octet-stream content negotiation so the cluster returns the staged export
	// rather than a JSON policy list, and reuse the prepare-export request id so the right
	// export is flushed.
	headers := map[string]interface{}{
		"Accept":          utils.StringPtr("application/octet-stream"),
		"NTNX-Request-Id": utils.StringPtr(requestID),
	}

	listResp, err := conn.NetworkingSecurityInstance.ListNetworkSecurityPolicies(ctx, &listReq, headers)
	if err != nil {
		return nil, err
	}
	if listResp == nil || listResp.Data == nil {
		return nil, fmt.Errorf("empty response while downloading exported payload")
	}

	fileDetail, ok := listResp.Data.GetValue().(microsegConfig.FileDetail)
	if !ok || fileDetail.Path == nil {
		return nil, fmt.Errorf("unexpected response type while downloading exported payload")
	}

	// The SDK downloads the octet-stream into a temporary file (in the working directory);
	// read it and remove it so it does not linger after the export completes.
	defer func() {
		if rmErr := os.Remove(*fileDetail.Path); rmErr != nil {
			log.Printf("[WARN] failed to remove temporary export download (%s): %v", *fileDetail.Path, rmErr)
		}
	}()

	contents, err := os.ReadFile(*fileDetail.Path)
	if err != nil {
		return nil, fmt.Errorf("failed reading downloaded export file (%s): %w", *fileDetail.Path, err)
	}

	return contents, nil
}

// isConcurrentOperationErr reports whether an error is the microseg backend rejecting an
// operation because another policy import/export/create/delete/update is still in flight.
// The backend serializes these operations cluster-wide, so this is a transient condition
// that clears once the in-flight operation releases its lock.
func isConcurrentOperationErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "concurrent operations are not permitted")
}

// isNetworkSecurityPolicyEntity reports whether a task-affected entity reference corresponds to
// a network security policy, based on its `rel` (namespace:module[:submodule]:entityType).
func isNetworkSecurityPolicyEntity(rel *string) bool {
	if rel == nil {
		return false
	}
	r := strings.ToLower(*rel)
	return strings.Contains(r, utils.RelEntityTypeSecurityPolicy)
}
