package ndb

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBDatabaseRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBDatabaseRestoreCreate,
		ReadContext:   resourceNutanixNDBDatabaseRestoreRead,
		UpdateContext: resourceNutanixNDBDatabaseRestoreUpdate,
		DeleteContext: resourceNutanixNDBDatabaseRestoreDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraProvisionTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"database_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_pitr_timestamp"},
			},
			"latest_snapshot": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_pitr_timestamp": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"snapshot_id"},
				RequiredWith:  []string{"time_zone_pitr"},
			},
			"time_zone_pitr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restore_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// computed Values

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"properties": {
				Type:        schema.TypeList,
				Description: "List of all the properties",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},

						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": dataSourceEraDBInstanceTags(),
			"clone": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbserver_logical_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"info": dataSourceEraDatabaseInfo(),
			"metric": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parent_database_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lcm_config":   dataSourceEraLCMConfig(),
			"time_machine": dataSourceEraTimeMachine(),
			"dbserver_logical_cluster": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"database_nodes":   dataSourceEraDatabaseNodes(),
			"linked_databases": dataSourceEraLinkedDatabases(),
			"database_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixNDBDatabaseRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	req := &era.DatabaseRestoreRequest{}

	databaseID := ""
	if dbID, ok := d.GetOk("database_id"); ok && len(dbID.(string)) > 0 {
		databaseID = dbID.(string)
	} else {
		return diag.Errorf("database_id is a required field to perform restore")
	}

	if snapID, ok := d.GetOk("snapshot_id"); ok {
		req.SnapshotID = utils.StringPtr(snapID.(string))
	}

	if latestsnap, ok := d.GetOk("latest_snapshot"); ok {
		req.LatestSnapshot = utils.StringPtr(latestsnap.(string))
	}

	if uptime, ok := d.GetOk("user_pitr_timestamp"); ok {
		req.UserPitrTimestamp = utils.StringPtr(uptime.(string))
	}

	if timezone, ok := d.GetOk("time_zone_pitr"); ok {
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

	resp, er := conn.Service.DatabaseRestore(ctx, databaseID, req)
	if er != nil {
		return diag.FromErr(er)
	}

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

	setID := databaseID + "/" + resp.Operationid
	d.SetId(setID)
	log.Printf("NDB database restore  with %s id is performed successfully", databaseID)
	return resourceNutanixNDBDatabaseRestoreRead(ctx, d, meta)
}

func resourceNutanixNDBDatabaseRestoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	splitID := strings.Split(d.Id(), "/")
	dbUUID := splitID[0]

	if databaseID, ok := d.GetOk("database_id"); ok {
		ctx = NewContext(ctx, dbID(databaseID.(string)))
	} else {
		ctx = NewContext(ctx, dbID(dbUUID))
	}
	return readDatabaseInstance(ctx, d, meta)
}

func resourceNutanixNDBDatabaseRestoreUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixNDBDatabaseRestoreCreate(ctx, d, meta)
}

func resourceNutanixNDBDatabaseRestoreDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
