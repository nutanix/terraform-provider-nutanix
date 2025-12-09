package prismv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixUnregisterClusterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixUnregisterClusterV2Create,
		ReadContext:   ResourceNutanixUnregisterClusterV2Read,
		UpdateContext: ResourceNutanixUnregisterClusterV2Update,
		DeleteContext: ResourceNutanixUnregisterClusterV2Delete,
		Schema: map[string]*schema.Schema{
			"pc_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func ResourceNutanixUnregisterClusterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clsConn := meta.(*conns.Client).ClusterAPI
	conn := meta.(*conns.Client).PrismAPI
	pcExtID := d.Get("pc_ext_id")

	readClsResp, readErr := conn.DomainManagerAPIInstance.GetDomainManagerById(utils.StringPtr(pcExtID.(string)))
	if readErr != nil {
		return diag.Errorf("error while fetching PC: %v", readErr)
	}

	args := make(map[string]interface{})
	eTag := clsConn.ClusterEntityAPI.ApiClient.GetEtag(readClsResp)
	args["If-Match"] = utils.StringPtr(eTag)

	extID := d.Get("ext_id")
	body := management.ClusterReference{
		ExtId: utils.StringPtr(extID.(string)),
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Unregister Cluster Request payload: %s", string(aJSON))

	// pass nil for the new dyRun flag
	resp, err := conn.DomainManagerAPIInstance.Unregister(utils.StringPtr(pcExtID.(string)), &body, nil, args)

	if err != nil {
		return diag.Errorf("error while unregistering cluster : %v", err)
	}

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster unregistration to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for cluster unregistration (%s) to complete: %s", utils.StringValue(taskUUID), err)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster unregistration task (%s): %v", utils.StringValue(taskUUID), err)
	}

	taskDetails := taskResp.Data.GetValue().(config.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Unregister Cluster task details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeDomainManager,
		"Unregistered Domain Manager")
	if err != nil {
		return diag.Errorf("error while extracting domain manager UUID from task response: %s", err)
	}
	d.SetId(utils.StringValue(uuid))

	return nil
}

func ResourceNutanixUnregisterClusterV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixUnregisterClusterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixUnregisterClusterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
