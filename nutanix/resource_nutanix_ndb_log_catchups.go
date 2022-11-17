package nutanix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
)

func resourceNutanixNDBLogCatchUps() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBLogCatchUpsCreate,
		ReadContext:   resourceNutanixNDBLogCatchUpsRead,
		UpdateContext: resourceNutanixNDBLogCatchUpsUpdate,
		DeleteContext: resourceNutanixNDBLogCatchUpsDelete,
		Schema: map[string]*schema.Schema{
			"time_machine_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"database_id"},
			},
			"database_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"time_machine_id"},
			},
			"for_restore": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNutanixNDBLogCatchUpsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era
	req := &era.LogCatchUpRequest{}

	tmsId := ""

	tm, tmOk := d.GetOk("time_machine_id")

	dbId, dbOk := d.GetOk("database_id")

	if !tmOk && !dbOk {
		return diag.Errorf("please provide the required `time_machine_id` or `database_id`  attribute")
	}

	if tmOk {
		tmsId = tm.(string)
	}

	if dbOk {
		// get the time machine id by getting database details

		dbResp, er := conn.Service.GetDatabaseInstance(ctx, dbId.(string))
		if er != nil {
			return diag.FromErr(er)
		}

		tmsId = dbResp.Timemachineid
	}

	// call log-catchup API

	actargs := []*era.Actionarguments{}

	if restore, rok := d.GetOkExists("for_restore"); rok {
		forRestore := restore.(bool)

		req.ForRestore = forRestore

		actargs = append(actargs, &era.Actionarguments{
			Name:  "preRestoreLogCatchup",
			Value: forRestore,
		})
	}

	actargs = append(actargs, &era.Actionarguments{
		Name:  "switch_log",
		Value: "true",
	})

	req.Actionarguments = actargs
	resp, err := conn.Service.LogCatchUp(ctx, tmsId, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Entityid)

	// Get Operation ID from response of log-catchups and poll for the operation to get completed.
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
		return diag.Errorf("error waiting to perform log-catchups	 (%s) to create: %s", resp.Entityid, errWaitTask)
	}
	return nil
}

func resourceNutanixNDBLogCatchUpsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBLogCatchUpsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBLogCatchUpsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
