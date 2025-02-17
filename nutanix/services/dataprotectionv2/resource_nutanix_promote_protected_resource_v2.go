package dataprotectionv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dataprtotectionPrismConfig "github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixPromoteProtectedResourceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixPromoteProtectedResourceV2Create,
		ReadContext:   ResourceNutanixPromoteProtectedResourceV2Read,
		UpdateContext: ResourceNutanixPromoteProtectedResourceV2Update,
		DeleteContext: ResourceNutanixPromoteProtectedResourceV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"promoted_vm_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// ResourceNutanixPromoteProtectedResourceV2Create to Promote Protected Resource
func ResourceNutanixPromoteProtectedResourceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataProtectionAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.ProtectedResource.PromoteProtectedResource(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("Error while promoting protected resource: %s", err)
	}

	TaskRef := resp.Data.GetValue().(dataprtotectionPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for task to complete: %s", err)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("Error while getting task by ID: %s", err)
	}

	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ := json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Promote Protected Resource Task Details: %s", aJSON)

	// // extract promoted item UUID
	// promotedItemName := rUUID.CompletionDetails[0].Name
	// promotedItemUUID := rUUID.CompletionDetails[0].Value.GetValue().(string)

	// if utils.StringValue(promotedItemName) == "promotedVmExtId" {
	// 	err = d.Set("promoted_vm_ext_id", promotedItemUUID)
	// 	if err != nil {
	// 		return diag.Errorf("Error while setting promoted Vm Ext ID: %s", err)
	// 	}
	// }

	// d.SetId(promotedItemUUID)
	d.SetId(utils.GenUUID())

	return ResourceNutanixPromoteProtectedResourceV2Read(ctx, d, meta)
}

// ResourceNutanixPromoteProtectedResourceV2Read to Promote Protected Resource
func ResourceNutanixPromoteProtectedResourceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixPromoteProtectedResourceV2Update to Promote Protected Resource
func ResourceNutanixPromoteProtectedResourceV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixPromoteProtectedResourceV2Delete to Promote Protected Resource
func ResourceNutanixPromoteProtectedResourceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
