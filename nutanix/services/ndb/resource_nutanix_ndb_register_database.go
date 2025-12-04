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

func ResourceNutanixNDBRegisterDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBRegisterDatabaseCreate,
		ReadContext:   resourceNutanixNDBRegisterDatabaseRead,
		UpdateContext: resourceNutanixNDBRegisterDatabaseUpdate,
		DeleteContext: resourceNutanixNDBRegisterDatabaseDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraProvisionTimeout),
			Update: schema.DefaultTimeout(EraProvisionTimeout),
			Delete: schema.DefaultTimeout(EraProvisionTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"database_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clustered": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"forced_install": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
			},
			"vm_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vm_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"vm_sshkey": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"vm_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reset_description_in_nx_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auto_tune_staging_drive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"working_directory": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/tmp",
			},
			"time_machine_info": timeMachineInfoSchema(),
			"tags":              dataSourceEraDBInstanceTags(),
			"actionarguments":   actionArgumentsSchema(),
			"postgress_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db_user": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"switch_log": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"allow_multiple_databases": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"backup_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "prefer_secondary",
						},
						"vm_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"postgres_software_home": {
							Type:     schema.TypeString,
							Required: true,
						},
						"software_home": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"db_password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"db_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			// delete values
			"delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"remove": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"soft_remove": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"forced": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"delete_time_machine": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"delete_logical_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// computed values

			"name": {
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
			"time_machine": dataSourceEraTimeMachine(),
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone": {
				Type:     schema.TypeBool,
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
			"database_status": {
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
			"parent_time_machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"info": dataSourceEraDatabaseInfo(),
			"group_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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
			"lcm_config": dataSourceEraLCMConfig(),
			"dbserver_logical_cluster": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"database_nodes":   dataSourceEraDatabaseNodes(),
			"linked_databases": dataSourceEraLinkedDatabases(),
		},
	}
}

func resourceNutanixNDBRegisterDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	log.Println("Creating the request!!!")
	req, err := buildReisterDBRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, er := conn.Service.RegisterDatabase(ctx, req)
	if er != nil {
		return diag.FromErr(er)
	}
	d.SetId(resp.Entityid)

	// Get Operation ID from response of RegisterDatabaseResponse and poll for the operation to get completed.
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
		return diag.Errorf("error waiting for db register	 (%s) to create: %s", resp.Entityid, errWaitTask)
	}
	log.Printf("NDB register database with %s id is created successfully", d.Id())
	return resourceNutanixNDBRegisterDatabaseRead(ctx, d, meta)
}

func resourceNutanixNDBRegisterDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	databaseID := d.Id()
	ctx = NewContext(ctx, dbID(databaseID))
	return readDatabaseInstance(ctx, d, meta)
}

func resourceNutanixNDBRegisterDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return updateDatabaseInstance(ctx, d, meta)
}

func resourceNutanixNDBRegisterDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	if conn == nil {
		return diag.Errorf("era is nil")
	}

	dbID := d.Id()

	req := &era.DeleteDatabaseRequest{}

	if delete, ok := d.GetOk("delete"); ok {
		req.Delete = delete.(bool)
	}

	if remove, ok := d.GetOk("remove"); ok {
		req.Remove = remove.(bool)
	}

	if softremove, ok := d.GetOk("soft_remove"); ok {
		req.Softremove = softremove.(bool)
	}

	if forced, ok := d.GetOk("forced"); ok {
		req.Forced = forced.(bool)
	}

	if deltms, ok := d.GetOk("delete_time_machine"); ok {
		req.Deletetimemachine = deltms.(bool)
	}

	if dellogicalcls, ok := d.GetOk("delete_logical_cluster"); ok {
		req.Deletelogicalcluster = dellogicalcls.(bool)
	}

	res, err := conn.Service.DeleteDatabase(ctx, req, dbID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Operation to unregister instance with id %s has started, operation id: %s", dbID, res.Operationid)
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
		return diag.Errorf("error waiting for unregister db Instance (%s) to delete: %s", res.Entityid, errWaitTask)
	}
	log.Printf("NDB register database with %s id is deleted successfully", d.Id())
	return nil
}

func buildReisterDBRequest(d *schema.ResourceData) (*era.RegisterDBInputRequest, error) {
	res := &era.RegisterDBInputRequest{}

	if dbType, ok := d.GetOk("database_type"); ok && len(dbType.(string)) > 0 {
		res.DatabaseType = utils.StringPtr(dbType.(string))
	}

	if dbName, ok := d.GetOk("database_name"); ok && len(dbName.(string)) > 0 {
		res.DatabaseName = utils.StringPtr(dbName.(string))
	}

	if desc, ok := d.GetOk("description"); ok && len(desc.(string)) > 0 {
		res.Description = utils.StringPtr(desc.(string))
	}

	if cls, ok := d.GetOk("clustered"); ok {
		res.Clustered = cls.(bool)
	}

	if forcedInstall, ok := d.GetOk("forced_install"); ok {
		res.ForcedInstall = forcedInstall.(bool)
	}

	if category, ok := d.GetOk("category"); ok && len(category.(string)) > 0 {
		res.Category = utils.StringPtr(category.(string))
	}

	if vmIP, ok := d.GetOk("vm_ip"); ok && len(vmIP.(string)) > 0 {
		res.VMIP = utils.StringPtr(vmIP.(string))
	}

	if vmUsername, ok := d.GetOk("vm_username"); ok && len(vmUsername.(string)) > 0 {
		res.VMUsername = utils.StringPtr(vmUsername.(string))
	}

	if vmPass, ok := d.GetOk("vm_password"); ok && len(vmPass.(string)) > 0 {
		res.VMPassword = utils.StringPtr(vmPass.(string))
	}

	if vmSshkey, ok := d.GetOk("vm_sshkey"); ok && len(vmSshkey.(string)) > 0 {
		res.VMSshkey = utils.StringPtr(vmSshkey.(string))
	}

	if forcedInstall, ok := d.GetOk("vm_description"); ok && len(forcedInstall.(string)) > 0 {
		res.ForcedInstall = forcedInstall.(bool)
	}

	if nxCls, ok := d.GetOk("nx_cluster_id"); ok && len(nxCls.(string)) > 0 {
		res.NxClusterID = utils.StringPtr(nxCls.(string))
	}

	if resetDesc, ok := d.GetOk("reset_description_in_nx_cluster"); ok {
		res.ResetDescriptionInNxCluster = resetDesc.(bool)
	}

	if autoTune, ok := d.GetOk("auto_tune_staging_drive"); ok {
		res.AutoTuneStagingDrive = (autoTune.(bool))
	}

	if wrk, ok := d.GetOk("working_directory"); ok && len(wrk.(string)) > 0 {
		res.WorkingDirectory = utils.StringPtr(wrk.(string))
	}

	if tms, ok := d.GetOk("time_machine_info"); ok && len(tms.(*schema.Set).List()) > 0 {
		res.TimeMachineInfo = buildTimeMachineFromResourceData(tms.(*schema.Set))
	}

	if tags, ok := d.GetOk("tags"); ok && len(tags.([]interface{})) > 0 {
		res.Tags = expandTags(tags.([]interface{}))
	}

	if maintenance, ok := d.GetOk("maintenance_tasks"); ok && len(maintenance.([]interface{})) > 0 {
		res.MaintenanceTasks = expandMaintenanceTasks(maintenance.([]interface{}))
	}

	res.Actionarguments = expandRegisterDBActionArguments(d)
	return res, nil
}

func expandRegisterDBActionArguments(d *schema.ResourceData) []*era.Actionarguments {
	args := []*era.Actionarguments{}
	if post, ok := d.GetOk("postgress_info"); ok {
		brr := post.([]interface{})

		for _, arg := range brr {
			val := arg.(map[string]interface{})
			var values interface{}
			if plist, pok := val["listener_port"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "listener_port",
					Value: values,
				})
			}
			if plist, pok := val["db_user"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "db_user",
					Value: values,
				})
			}
			if plist, pok := val["switch_log"]; pok && plist.(bool) {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "switch_log",
					Value: values,
				})
			}
			if plist, pok := val["allow_multiple_databases"]; pok && plist.(bool) {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "allow_multiple_databases",
					Value: values,
				})
			}
			if plist, pok := val["backup_policy"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "backup_policy",
					Value: values,
				})
			}
			if plist, pok := val["vm_ip"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "vmIp",
					Value: values,
				})
			}
			if plist, pok := val["postgres_software_home"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "postgres_software_home",
					Value: values,
				})
			}
			if plist, pok := val["software_home"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "software_home",
					Value: values,
				})
			}
			if plist, pok := val["db_password"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "db_password",
					Value: values,
				})
			}
			if plist, pok := val["db_name"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "db_name",
					Value: values,
				})
			}
		}
	}

	resp := buildActionArgumentsFromResourceData(d.Get("actionarguments").(*schema.Set), args)
	return resp
}
