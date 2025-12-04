package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBLinkedDB() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBLinkedDBCreate,
		ReadContext:   resourceNutanixNDBLinkedDBRead,
		UpdateContext: resourceNutanixNDBLinkedDBUpdate,
		DeleteContext: resourceNutanixNDBLinkedDBDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraProvisionTimeout),
			Delete: schema.DefaultTimeout(EraProvisionTimeout),
		},
		Schema: map[string]*schema.Schema{
			"database_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_database_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_linked_database_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secure_info": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"created_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"metric": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixNDBLinkedDBCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.CreateLinkedDatabasesRequest{}

	databaseID := ""
	databaseName := ""
	SetID := ""
	if dbID, dok := d.GetOk("database_id"); dok && len(dbID.(string)) > 0 {
		databaseID = dbID.(string)
	} else {
		return diag.Errorf("database_id is a required field")
	}

	dbNames := []*era.LinkedDatabases{}

	if dbName, ok := d.GetOk("database_name"); ok {
		dbNames = append(dbNames, &era.LinkedDatabases{
			DatabaseName: utils.StringPtr(dbName.(string)),
		})
		databaseName = dbName.(string)
	}

	req.Databases = dbNames

	// call the Linked Databases API

	resp, err := conn.Service.CreateLinkedDatabase(ctx, databaseID, req)
	if err != nil {
		return diag.FromErr(err)
	}

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
		return diag.Errorf("error waiting for databases	 (%s) to add: %s", resp.Entityid, errWaitTask)
	}

	// call the databases API

	response, er := conn.Service.GetDatabaseInstance(ctx, resp.Entityid)
	if er != nil {
		return diag.FromErr(er)
	}

	linkDbs := response.Linkeddatabases

	for _, v := range linkDbs {
		if v.DatabaseName == databaseName {
			SetID = v.ID
			break
		}
	}

	d.SetId(SetID)
	log.Printf("NDB linked database with %s id is created successfully", d.Id())

	return resourceNutanixNDBLinkedDBRead(ctx, d, meta)
}

func resourceNutanixNDBLinkedDBRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	databaseID := d.Get("database_id")

	// check if database id is nil
	if databaseID == "" {
		return diag.Errorf("database id is required for read operation")
	}
	response, er := conn.Service.GetDatabaseInstance(ctx, databaseID.(string))
	if er != nil {
		return diag.FromErr(er)
	}

	linkDbs := response.Linkeddatabases
	currentLinkedDB := &era.Linkeddatabases{}

	for _, v := range linkDbs {
		if v.ID == d.Id() {
			*currentLinkedDB = v
			break
		}
	}

	if err := d.Set("database_name", currentLinkedDB.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database_status", currentLinkedDB.Databasestatus); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("date_created", currentLinkedDB.Datecreated); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("date_modified", currentLinkedDB.Datemodified); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", currentLinkedDB.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("info", flattenLinkedDBInfo(currentLinkedDB.Info)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metric", currentLinkedDB.Metric); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", currentLinkedDB.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parent_database_id", currentLinkedDB.ParentDatabaseID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parent_linked_database_id", currentLinkedDB.ParentLinkedDatabaseID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("snapshot_id", currentLinkedDB.SnapshotID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", currentLinkedDB.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("timezone", currentLinkedDB.TimeZone); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixNDBLinkedDBUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBLinkedDBDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	dbID := d.Get("database_id")

	req := &era.DeleteLinkedDatabaseRequest{
		Delete: true,
		Forced: true,
	}

	// API to delete linked databases

	res, err := conn.Service.DeleteLinkedDatabase(ctx, dbID.(string), d.Id(), req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Operation to delete linked databases with id %s has started, operation id: %s", d.Id(), res.Operationid)
	opID := res.Operationid
	if opID == "" {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := era.GetOperationRequest{
		OperationID: opID,
	}

	log.Printf("polling for operation with id: %s\n", opID)

	// Poll for operation here - Cluster GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: eraRefresh(ctx, conn, opReq),
		Timeout: d.Timeout(schema.TimeoutDelete),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for linked db (%s) to delete: %s", d.Id(), errWaitTask)
	}
	log.Printf("NDB linked database with %s id is deleted successfully", d.Id())
	return nil
}

func flattenLinkedDBInfo(pr era.Info) []interface{} {
	res := make([]interface{}, 0)
	info := make(map[string]interface{})

	if pr.Secureinfo != nil {
		info["secure_info"] = pr.Secureinfo
	}

	if pr.Info != nil {
		inf := make([]interface{}, 0)
		infval := make(map[string]interface{})

		if pr.Info.CreatedBy != nil {
			infval["created_by"] = pr.Info.CreatedBy
		}

		inf = append(inf, infval)
		info["info"] = inf
	}

	res = append(res, info)
	return res
}
