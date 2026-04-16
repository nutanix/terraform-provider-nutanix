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

func ResourceNutanixNDBScaleDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBScaleDatabaseCreate,
		ReadContext:   resourceNutanixNDBScaleDatabaseRead,
		UpdateContext: resourceNutanixNDBScaleDatabaseUpdate,
		DeleteContext: resourceNutanixNDBScaleDatabaseDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraProvisionTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
			"scale_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"databasetype": {
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
			"dbserver_logical_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
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
			"parent_source_database_id": {
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

func resourceNutanixNDBScaleDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.DatabaseScale{}
	dbUUID := ""
	if db, dok := d.GetOk("database_uuid"); dok && len(db.(string)) > 0 {
		dbUUID = db.(string)
	} else {
		return diag.Errorf("database_id is a required field to perform scale")
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

	// adding working dir

	args = append(args, &era.Actionarguments{
		Name:  "working_dir",
		Value: "/tmp",
	})

	req.Actionarguments = args

	// call API

	resp, err := conn.Service.DatabaseScale(ctx, dbUUID, req)
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
		return diag.Errorf("error waiting for db Instance    (%s) to scale: %s", resp.Entityid, errWaitTask)
	}

	setID := dbUUID + "/" + resp.Operationid
	d.SetId(setID)
	log.Printf("NDB database with %s id is scaled successfully", dbUUID)
	return resourceNutanixNDBScaleDatabaseRead(ctx, d, meta)
}

func resourceNutanixNDBScaleDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	splitID := strings.Split(d.Id(), "/")
	dbUUID := splitID[0]

	if databaseID, ok := d.GetOk("database_uuid"); ok {
		ctx = NewContext(ctx, dbID(databaseID.(string)))
	} else {
		ctx = NewContext(ctx, dbID(dbUUID))
	}
	return readDatabaseInstance(ctx, d, meta)
}

func resourceNutanixNDBScaleDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixNDBScaleDatabaseCreate(ctx, d, meta)
}

func resourceNutanixNDBScaleDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
