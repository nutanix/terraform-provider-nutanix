package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func ResourceNutanixNDBLogCatchUps() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBLogCatchUpsCreate,
		ReadContext:   resourceNutanixNDBLogCatchUpsRead,
		UpdateContext: resourceNutanixNDBLogCatchUpsUpdate,
		DeleteContext: resourceNutanixNDBLogCatchUpsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraProvisionTimeout),
		},
		Schema: map[string]*schema.Schema{
			"time_machine_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"database_id"},
				ForceNew:      true,
			},
			"database_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"time_machine_id"},
				ForceNew:      true,
			},
			"for_restore": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"log_catchup_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNutanixNDBLogCatchUpsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	req := &era.LogCatchUpRequest{}

	tmsID := ""

	tm, tmOk := d.GetOk("time_machine_id")

	dbID, dbOk := d.GetOk("database_id")

	if !tmOk && !dbOk {
		return diag.Errorf("please provide the required `time_machine_id` or `database_id`  attribute")
	}

	if tmOk {
		tmsID = tm.(string)
	}

	if dbOk {
		// get the time machine id by getting database details

		dbResp, er := conn.Service.GetDatabaseInstance(ctx, dbID.(string))
		if er != nil {
			return diag.FromErr(er)
		}

		tmsID = dbResp.Timemachineid
	}

	// call log-catchup API

	actargs := []*era.Actionarguments{}
	//nolint:staticcheck
	if restore, rok := d.GetOkExists("for_restore"); rok && restore.(bool) {
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
	resp, err := conn.Service.LogCatchUp(ctx, tmsID, req)
	if err != nil {
		return diag.FromErr(err)
	}

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
	d.SetId(resp.Operationid)
	log.Printf("NDB log catchup with %s id is performed successfully", d.Id())
	return nil
}

func resourceNutanixNDBLogCatchUpsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBLogCatchUpsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixNDBLogCatchUpsCreate(ctx, d, meta)
}

func resourceNutanixNDBLogCatchUpsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
