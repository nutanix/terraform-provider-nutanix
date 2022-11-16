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

func resourceNutanixNDBDatabaseRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBDatabaseRestoreCreate,
		ReadContext:   resourceNutanixNDBDatabaseRestoreRead,
		UpdateContext: resourceNutanixNDBDatabaseRestoreUpdate,
		DeleteContext: resourceNutanixNDBDatabaseRestoreDelete,
		Schema: map[string]*schema.Schema{
			"database_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"latest_snapshot": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_pitr_timestamp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNutanixNDBDatabaseRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era
	req := &era.DatabaseRestoreRequest{}

	databaseId := ""
	if dbId, ok := d.GetOk("database_id"); ok {
		databaseId = dbId.(string)
	}

	if snapId, ok := d.GetOk("snapshot_id"); ok {
		req.SnapshotId = utils.StringPtr(snapId.(string))
	}

	if latestsnap, ok := d.GetOk("latest_snapshot"); ok {
		req.LatestSnapshot = utils.StringPtr(latestsnap.(string))
	}

	if uptime, ok := d.GetOk("user_pitr_timestamp"); ok {
		req.UserPitrTimestamp = utils.StringPtr(uptime.(string))
	}

	if timezone, ok := d.GetOk("time_zone"); ok {
		req.TimeZone = utils.StringPtr(timezone.(string))
	}

	// getting action arguments

	actargs := []*era.Actionarguments{}

	actargs = append(actargs, &era.Actionarguments{
		Name:  "sameLocation",
		Value: "true",
	})

	req.ActionArguments = actargs

	// call the database restore API

	resp, er := conn.Service.DatabaseRestore(ctx, databaseId, req)
	if er != nil {
		return diag.FromErr(er)
	}

	d.SetId(resp.Entityid)

	// Get Operation ID from response of database restore and poll for the operation to get completed.
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
		return diag.Errorf("error waiting to perform db restore	 (%s) to create: %s", resp.Entityid, errWaitTask)
	}

	return nil
}

func resourceNutanixNDBDatabaseRestoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBDatabaseRestoreUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBDatabaseRestoreDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
