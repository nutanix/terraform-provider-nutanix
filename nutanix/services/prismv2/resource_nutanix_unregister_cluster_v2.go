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
		return diag.Errorf("error while fetching Domain Manager: %v", readErr)
	}

	args := make(map[string]interface{})
	eTag := clsConn.ClusterEntityAPI.ApiClient.GetEtag(readClsResp)
	args["If-Match"] = utils.StringPtr(eTag)

	extID := d.Get("ext_id")
	body := management.ClusterReference{
		ExtId: utils.StringPtr(extID.(string)),
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Unregister Cluster Request Body: %s", string(aJSON))

	resp, err := conn.DomainManagerAPIInstance.Unregister(utils.StringPtr(pcExtID.(string)), &body, args)

	if err != nil {
		return diag.Errorf("error while unregistering cluster : %v", err)
	}

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for Cluster Unregister Task to complete: %s", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching task details : %v", err)
	}

	rUUID := resourceUUID.Data.GetValue().(config.Task)
	aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Unregister Cluster Task Details: %s", string(aJSON))

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

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
