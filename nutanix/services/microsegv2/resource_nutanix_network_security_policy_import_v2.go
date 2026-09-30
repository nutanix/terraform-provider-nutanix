package microsegv2

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/request/networksecuritypolicies"
	prismMicroseg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNetworkSecurityPolicyImportV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNetworkSecurityPolicyImportV2Create,
		ReadContext:   resourceNutanixNetworkSecurityPolicyImportV2Read,
		// No UpdateContext: every field is ForceNew or Computed, so Terraform always
		// destroys and recreates the resource instead of updating it in place.
		DeleteContext: resourceNutanixNetworkSecurityPolicyImportV2Delete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The file path of the data file to import network security policies. Changing this forces the imported policies to be destroyed and re-imported.",
			},
			"dryrun": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The header parameter specifies if the import should be a dry run.",
			},
			"ntnx_purge_policies": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The header parameter specifies if the existing policies need to be deleted or retained upon network security policy import.",
			},
			"ntnx_project_ext_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The header parameter specifies the project external identifier for policies being imported.",
			},
			"task_ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier for the task created by the import operation.",
			},
			"imported_policy_ext_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of network security policy external identifiers (UUIDs) that were created by this import.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNutanixNetworkSecurityPolicyImportV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	path := d.Get("path").(string)
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		return diag.Errorf("The provided path '%s' is not a valid file", path)
	}

	req := import1.ApplyNetworkSecurityPolicyImportRequest{
		Path: utils.StringPtr(path),
	}
	if v := common.IsExplicitlySet(d, "ntnx_purge_policies"); v {
		req.NTNXPurgePolicies = utils.BoolPtr(d.Get("ntnx_purge_policies").(bool))
	}
	if v, ok := d.GetOk("ntnx_project_ext_id"); ok {
		req.NTNXProjectExtId = utils.StringPtr(v.(string))
	}
	if v := common.IsExplicitlySet(d, "dryrun"); v {
		req.Dryrun_ = utils.BoolPtr(d.Get("dryrun").(bool))
	}

	aJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("[DEBUG] ApplyNetworkSecurityPolicyImport Request: %s", string(aJSON))

	resp, err := conn.NetworkingSecurityInstance.ApplyNetworkSecurityPolicyImport(ctx, &req)
	if err != nil {
		return diag.Errorf("error applying network security policy import: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("unexpected response type from ApplyNetworkSecurityPolicyImport")
	}
	taskUUID := taskRef.ExtId

	taskConn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskConn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for network security policy import task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	getTaskByIDRequest := import2.GetTaskByIdRequest{
		ExtId: taskUUID,
	}
	taskResp, err := taskConn.TaskRefAPI.GetTaskById(ctx, &getTaskByIDRequest)
	if err != nil {
		return diag.Errorf("error fetching network security policy import task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] ApplyNetworkSecurityPolicyImport Task Details: %s", string(aJSON))

	// Parse the task (and its subtasks) to discover the policies created by the import.
	// A dry run only validates the import and creates nothing, so there are no imported
	// policies to track.
	importedExtIDs := make([]string, 0)
	if !d.Get("dryrun").(bool) {
		importedExtIDs = collectImportedPolicyExtIDs(ctx, taskConn, taskDetails)
	}

	d.SetId(utils.StringValue(taskDetails.ExtId))
	if err := d.Set("task_ext_id", utils.StringValue(taskDetails.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("imported_policy_ext_ids", importedExtIDs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// collectImportedPolicyExtIDs returns the ext_ids of the network security policies created by
// an import task. The cluster reports the created policies either on the parent task's
// entitiesAffected or only on its subtasks, so both are inspected and de-duplicated.
func collectImportedPolicyExtIDs(ctx context.Context, taskConn *prism.Client, task prismConfig.Task) []string {
	ids := make([]string, 0)
	seen := make(map[string]struct{})

	addEntities := func(entities []prismConfig.EntityReference) {
		for _, entity := range entities {
			if !isNetworkSecurityPolicyEntity(entity.Rel) || entity.ExtId == nil {
				continue
			}
			if _, ok := seen[*entity.ExtId]; ok {
				continue
			}
			seen[*entity.ExtId] = struct{}{}
			ids = append(ids, *entity.ExtId)
		}
	}

	addEntities(task.EntitiesAffected)

	for _, sub := range task.SubTasks {
		if sub.ExtId == nil {
			continue
		}
		getReq := import2.GetTaskByIdRequest{ExtId: sub.ExtId}
		subResp, err := taskConn.TaskRefAPI.GetTaskById(ctx, &getReq)
		if err != nil {
			log.Printf("[WARN] error fetching import subtask (%s): %v", *sub.ExtId, err)
			continue
		}
		subTask, ok := subResp.Data.GetValue().(prismConfig.Task)
		if !ok {
			continue
		}
		addEntities(subTask.EntitiesAffected)
	}

	return ids
}

func resourceNutanixNetworkSecurityPolicyImportV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	importedExtIDs := make([]string, 0)
	for _, v := range d.Get("imported_policy_ext_ids").([]interface{}) {
		if v == nil {
			continue
		}
		importedExtIDs = append(importedExtIDs, v.(string))
	}

	if len(importedExtIDs) == 0 {
		return nil
	}

	// Reconcile state with reality: keep only the policies that still exist on the cluster.
	stillExisting := make([]string, 0, len(importedExtIDs))
	for _, extID := range importedExtIDs {
		getReq := import1.GetNetworkSecurityPolicyByIdRequest{
			ExtId: utils.StringPtr(extID),
		}
		if _, err := conn.NetworkingSecurityInstance.GetNetworkSecurityPolicyById(ctx, &getReq); err != nil {
			log.Printf("[DEBUG] imported network security policy (%s) not found: %v", extID, err)
			continue
		}
		stillExisting = append(stillExisting, extID)
	}

	// If every imported policy was deleted out-of-band, drop the resource so it is recreated.
	if len(stillExisting) == 0 {
		log.Printf("[WARN] all imported network security policies for import (%s) were deleted; removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("imported_policy_ext_ids", stillExisting); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixNetworkSecurityPolicyImportV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI
	taskConn := meta.(*conns.Client).PrismAPI

	for _, v := range d.Get("imported_policy_ext_ids").([]interface{}) {
		if v == nil {
			continue
		}
		extID := v.(string)

		deleteReq := import1.DeleteNetworkSecurityPolicyByIdRequest{
			ExtId: utils.StringPtr(extID),
		}
		resp, err := conn.NetworkingSecurityInstance.DeleteNetworkSecurityPolicyById(ctx, &deleteReq)
		if err != nil {
			return diag.Errorf("error deleting imported network security policy (%s): %v", extID, err)
		}

		taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
		if !ok {
			return diag.Errorf("unexpected response type from DeleteNetworkSecurityPolicyById for policy (%s)", extID)
		}
		taskUUID := taskRef.ExtId

		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskConn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutDelete),
		}
		if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
			return diag.Errorf("error waiting for imported network security policy (%s) to delete: %s", extID, errWait)
		}
	}

	return nil
}
