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

func resourceNutanixNDBScaleDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBScaleDatabaseCreate,
		ReadContext:   resourceNutanixNDBScaleDatabaseRead,
		UpdateContext: resourceNutanixNDBScaleDatabaseUpdate,
		DeleteContext: resourceNutanixNDBScaleDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"database_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_storage_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"pre_script_cmd": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"post_script_cmd": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNutanixNDBScaleDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	req := &era.DatabaseScale{}
	dbUUID := ""
	if db, ok := d.GetOk("database_uuid"); ok {
		dbUUID = db.(string)
	}

	if app, ok := d.GetOk("application_type"); ok {
		req.ApplicationType = utils.StringPtr(app.(string))
	}

	// action arguments

	args := []*era.Actionarguments{}

	if dataSize, ok := d.GetOk("data_storage_size"); ok {
		args = append(args, &era.Actionarguments{
			Name:  "data_storage_size",
			Value: utils.IntPtr(dataSize.(int)),
		})
	}

	if pre, ok := d.GetOk("pre_script_cmd"); ok {
		args = append(args, &era.Actionarguments{
			Name:  "pre_script_cmd",
			Value: utils.StringPtr(pre.(string)),
		})
	}

	if post, ok := d.GetOk("post_script_cmd"); ok {
		args = append(args, &era.Actionarguments{
			Name:  "post_script_cmd",
			Value: utils.StringPtr(post.(string)),
		})
	}

	req.Actionarguments = args

	// call API

	resp, err := conn.Service.DatabaseScale(ctx, dbUUID, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Entityid)

	// Get Operation ID from response of ProvisionDatabaseResponse and poll for the operation to get completed.
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
		return diag.Errorf("error waiting for db Instance    (%s) to scale: %s", resp.Entityid, errWaitTask)
	}

	return nil
}

func resourceNutanixNDBScaleDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBScaleDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBScaleDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
