package nutanix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixNDBDatabaseSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBDatabaseSnapshotCreate,
		ReadContext:   resourceNutanixNDBDatabaseSnapshotRead,
		UpdateContext: resourceNutanixNDBDatabaseSnapshotUpdate,
		DeleteContext: resourceNutanixNDBDatabaseSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"time_machine_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remove_schedule_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNutanixNDBDatabaseSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	req := &era.DatabaseSnapshotRequest{}
	timeMachineID := ""

	if tm, ok := d.GetOk("time_machine_id"); ok {
		timeMachineID = tm.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		req.Name = utils.StringPtr(name.(string))
	}

	if rm, ok := d.GetOk("remove_schedule_in_days"); ok {

		lcmConfig := &era.LcmConfig{}
		expDetails := &era.DBExpiryDetails{}

		expDetails.ExpireInDays = utils.IntPtr(rm.(int))

		lcmConfig.ExpiryDetails = expDetails
		req.LcmConfig = lcmConfig
	}

	// call the snapshot API

	resp, err := conn.Service.DatabaseSnapshot(ctx, timeMachineID, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Entityid)

	// Get Operation ID from response of snapshot and poll for the operation to get completed.
	opID := resp.Operationid
	if opID == "" {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := era.GetOperationRequest{
		OperationID: opID,
	}

	log.Printf("polling for operation with id: %s\n", opID)

	// Poll for operation here - Operation GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: eraRefresh(ctx, conn, opReq),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for snapshot	 (%s) to create: %s", resp.Entityid, errWaitTask)
	}

	return nil
}

func resourceNutanixNDBDatabaseSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBDatabaseSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBDatabaseSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
