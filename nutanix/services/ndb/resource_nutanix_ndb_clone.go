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

func ResourceNutanixNDBClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBCloneCreate,
		ReadContext:   resourceNutanixNDBCloneRead,
		UpdateContext: resourceNutanixNDBCloneUpdate,
		DeleteContext: resourceNutanixNDBCloneDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraProvisionTimeout),
			Delete: schema.DefaultTimeout(EraProvisionTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"time_machine_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"time_machine_name"},
			},
			"time_machine_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"time_machine_id"},
			},
			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_pitr_timestamp"},
			},
			"user_pitr_timestamp": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"snapshot_id"},
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"nodes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"compute_profile_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_profile_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"new_db_server_time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nx_cluster_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"properties": {
							Type:        schema.TypeList,
							Description: "List of all the properties",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"dbserver_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"lcm_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_lcm_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expiry_details": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"expire_in_days": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"expiry_date_timezone": {
													Type:     schema.TypeString,
													Required: true,
												},
												"delete_database": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"refresh_details": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"refresh_in_days": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"refresh_time": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"refresh_date_timezone": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_public_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"compute_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_parameter_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vm_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"create_dbserver": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"clustered": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dbserver_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbserver_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbserver_logical_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"latest_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"postgresql_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dbserver_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"db_password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"pre_clone_cmd": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"post_clone_cmd": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"actionarguments": actionArgumentsSchema(),
			// delete arguments for clone resource.
			"delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"remove": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"remove_logical_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Computed values

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
		},
	}
}

func resourceNutanixNDBCloneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	req := &era.CloneRequest{}

	tmsID, tok := d.GetOk("time_machine_id")
	tmsName, tnOk := d.GetOk("time_machine_name")

	if !tok && !tnOk {
		return diag.Errorf("Atleast one of time_machine_id or time_machine_name is required to perform clone")
	}

	if len(tmsName.(string)) > 0 {
		// call time machine API with value-type name
		res, err := conn.Service.GetTimeMachine(ctx, "", tmsName.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		tmsID = *res.ID
	}

	req.TimeMachineID = utils.StringPtr(tmsID.(string))

	// build request for clone
	if err := buildCloneRequest(d, req); err != nil {
		return diag.FromErr(err)
	}

	// call clone API

	resp, err := conn.Service.CreateClone(ctx, tmsID.(string), req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.Entityid)

	// Get Operation ID from response of Clone and poll for the operation to get completed.
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
		return diag.Errorf("error waiting for time machine clone (%s) to create: %s", resp.Entityid, errWaitTask)
	}

	log.Printf("NDB clone with %s id is created successfully", d.Id())
	return resourceNutanixNDBCloneRead(ctx, d, meta)
}

func resourceNutanixNDBCloneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	filterParams := &era.FilterParams{}
	filterParams.Detailed = "true"
	filterParams.AnyStatus = "false"
	filterParams.LoadDBServerCluster = "false"
	filterParams.TimeZone = "UTC"

	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("id is required for read operation")
	}

	resp, err := conn.Service.GetClone(ctx, d.Id(), "", filterParams)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp != nil {
		if err = d.Set("description", resp.Description); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("name", resp.Name); err != nil {
			return diag.FromErr(err)
		}

		props := []interface{}{}
		for _, prop := range resp.Properties {
			props = append(props, map[string]interface{}{
				"name":  prop.Name,
				"value": prop.Value,
			})
		}
		if err := d.Set("properties", props); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("date_created", resp.Datecreated); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("date_modified", resp.Datemodified); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("clone", resp.Clone); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("database_name", resp.Databasename); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("type", resp.Type); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("database_cluster_type", resp.Databaseclustertype); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("status", resp.Status); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dbserver_logical_cluster_id", resp.Dbserverlogicalclusterid); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("time_zone", resp.Timezone); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("info", flattenDBInfo(resp.Info)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("metric", resp.Metric); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("parent_database_id", resp.ParentDatabaseID); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("time_machine", flattenDBTimeMachine(resp.TimeMachine)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dbserver_logical_cluster", resp.Dbserverlogicalcluster); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("database_nodes", flattenDBNodes(resp.Databasenodes)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("linked_databases", flattenDBLinkedDbs(resp.Linkeddatabases)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceNutanixNDBCloneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	dbID := d.Id()

	name := ""
	description := ""

	if d.HasChange("name") {
		name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		description = d.Get("description").(string)
	}

	tags := make([]*era.Tags, 0)
	if d.HasChange("tags") {
		tags = expandTags(d.Get("tags").([]interface{}))
	}

	updateReq := era.UpdateDatabaseRequest{
		Name:             name,
		Description:      description,
		Tags:             tags,
		Resetname:        true,
		Resetdescription: true,
		Resettags:        true,
	}

	res, err := conn.Service.UpdateCloneDatabase(ctx, dbID, &updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if res != nil {
		log.Printf("NDB clone with %s id is updated successfully", d.Id())
	}

	return resourceNutanixNDBCloneRead(ctx, d, meta)
}

func resourceNutanixNDBCloneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if remlogicalcls, ok := d.GetOk("remove_logical_cluster"); ok {
		req.Deletelogicalcluster = remlogicalcls.(bool)
	}

	res, err := conn.Service.DeleteClone(ctx, dbID, req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Operation to unregister clone instance with id %s has started, operation id: %s", dbID, res.Operationid)
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
		return diag.Errorf("error waiting for clone Instance (%s) to unregister: %s", res.Entityid, errWaitTask)
	}
	log.Printf("NDB clone with %s id is deleted successfully", d.Id())
	return nil
}

func buildCloneRequest(d *schema.ResourceData, res *era.CloneRequest) error {
	if name, ok := d.GetOk("name"); ok {
		res.Name = utils.StringPtr(name.(string))
	}

	if des, ok := d.GetOk("description"); ok {
		res.Description = utils.StringPtr(des.(string))
	}

	if nxcls, ok := d.GetOk("nx_cluster_id"); ok {
		res.NxClusterID = utils.StringPtr(nxcls.(string))
	}

	if ssh, ok := d.GetOk("ssh_public_key"); ok {
		res.SSHPublicKey = utils.StringPtr(ssh.(string))
	}
	if userPitrTimestamp, ok := d.GetOk("user_pitr_timestamp"); ok {
		res.UserPitrTimestamp = utils.StringPtr(userPitrTimestamp.(string))
	}
	if timeZone, ok := d.GetOk("time_zone"); ok && len(timeZone.(string)) > 0 {
		res.TimeZone = utils.StringPtr(timeZone.(string))
	}
	if computeProfileID, ok := d.GetOk("compute_profile_id"); ok {
		res.ComputeProfileID = utils.StringPtr(computeProfileID.(string))
	}
	if networkProfileID, ok := d.GetOk("network_profile_id"); ok {
		res.NetworkProfileID = utils.StringPtr(networkProfileID.(string))
	}
	if databaseParameterProfileID, ok := d.GetOk("database_parameter_profile_id"); ok {
		res.DatabaseParameterProfileID = utils.StringPtr(databaseParameterProfileID.(string))
	}
	if snapshotID, ok := d.GetOk("snapshot_id"); ok {
		res.SnapshotID = utils.StringPtr(snapshotID.(string))
	}

	if dbserverID, ok := d.GetOk("dbserver_id"); ok {
		res.DbserverID = utils.StringPtr(dbserverID.(string))
	}
	if dbserverClusterID, ok := d.GetOk("dbserver_cluster_id"); ok {
		res.DbserverClusterID = utils.StringPtr(dbserverClusterID.(string))
	}
	if dbserverLogicalClusterID, ok := d.GetOk("dbserver_logical_cluster_id"); ok {
		res.DbserverLogicalClusterID = utils.StringPtr(dbserverLogicalClusterID.(string))
	}
	if createDbserver, ok := d.GetOk("create_dbserver"); ok {
		res.CreateDbserver = createDbserver.(bool)
	}
	if clustered, ok := d.GetOk("clustered"); ok {
		res.Clustered = clustered.(bool)
	}
	if nodeCount, ok := d.GetOk("node_count"); ok {
		res.NodeCount = utils.IntPtr(nodeCount.(int))
	}

	if nodes, ok := d.GetOk("nodes"); ok {
		res.Nodes = expandClonesNodes(nodes.([]interface{}))
	}

	if lcmConfig, ok := d.GetOk("lcm_config"); ok {
		res.LcmConfig = expandLCMConfig(lcmConfig.([]interface{}))
	}

	if postgres, ok := d.GetOk("postgresql_info"); ok && len(postgres.([]interface{})) > 0 {
		res.ActionArguments = expandPostgreSQLCloneActionArgs(d, postgres.([]interface{}))
	}

	if tags, ok := d.GetOk("tags"); ok && len(tags.([]interface{})) > 0 {
		res.Tags = expandTags(tags.([]interface{}))
	}
	return nil
}

func expandClonesNodes(pr []interface{}) []*era.Nodes {
	nodes := make([]*era.Nodes, len(pr))
	if len(pr) > 0 {
		for k, v := range pr {
			val := v.(map[string]interface{})
			node := &era.Nodes{}

			if v1, ok1 := val["network_profile_id"]; ok1 && len(v1.(string)) > 0 {
				node.Networkprofileid = utils.StringPtr(v1.(string))
			}

			if v1, ok1 := val["compute_profile_id"]; ok1 && len(v1.(string)) > 0 {
				node.ComputeProfileID = utils.StringPtr(v1.(string))
			}

			if v1, ok1 := val["vm_name"]; ok1 && len(v1.(string)) > 0 {
				node.Vmname = utils.StringPtr(v1.(string))
			}

			if v1, ok1 := val["nx_cluster_id"]; ok1 && len(v1.(string)) > 0 {
				node.NxClusterID = utils.StringPtr(v1.(string))
			}

			if v1, ok1 := val["new_db_server_time_zone"]; ok1 && len(v1.(string)) > 0 {
				node.NewDBServerTimeZone = utils.StringPtr(v1.(string))
			}
			if v1, ok1 := val["properties"]; ok1 && len(v1.([]interface{})) > 0 {
				node.Properties = v1.([]*era.NodesProperties)
			}

			if v1, ok1 := val["dbserver_id"]; ok1 && len(v1.(string)) > 0 {
				node.DatabaseServerID = utils.StringPtr(v1.(string))
			}
			nodes[k] = node
		}
		return nodes
	}
	return nil
}

func expandPostgreSQLCloneActionArgs(d *schema.ResourceData, pr []interface{}) []*era.Actionarguments {
	if len(pr) > 0 {
		args := []*era.Actionarguments{}

		postgresProp := pr[0].(map[string]interface{})
		for key, value := range postgresProp {
			args = append(args, &era.Actionarguments{
				Name:  key,
				Value: utils.StringPtr(value.(string)),
			})
		}
		resp := buildActionArgumentsFromResourceData(d.Get("actionarguments").(*schema.Set), args)
		return resp
	}
	return nil
}

func expandLCMConfig(pr []interface{}) *era.CloneLCMConfig {
	if len(pr) > 0 {
		cloneLcm := &era.CloneLCMConfig{}
		for _, v := range pr {
			val := v.(map[string]interface{})

			if v1, ok1 := val["database_lcm_config"]; ok1 && len(v1.([]interface{})) > 0 {
				dbLcm := v1.([]interface{})
				dbLcmConfig := &era.DatabaseLCMConfig{}
				for _, v := range dbLcm {
					val := v.(map[string]interface{})

					if exp, ok1 := val["expiry_details"]; ok1 {
						dbLcmConfig.ExpiryDetails = expandDBExpiryDetails(exp.([]interface{}))
					}

					if ref, ok1 := val["refresh_details"]; ok1 {
						dbLcmConfig.RefreshDetails = expandDBRefreshDetails(ref.([]interface{}))
					}
				}
				cloneLcm.DatabaseLCMConfig = dbLcmConfig
			}
		}
		return cloneLcm
	}
	return nil
}

func expandDBExpiryDetails(pr []interface{}) *era.DBExpiryDetails {
	if len(pr) > 0 {
		expDetails := &era.DBExpiryDetails{}

		for _, v := range pr {
			val := v.(map[string]interface{})

			if v1, ok1 := val["expire_in_days"]; ok1 {
				expDetails.ExpireInDays = utils.IntPtr(v1.(int))
			}
			if v1, ok1 := val["expiry_date_timezone"]; ok1 && len(v1.(string)) > 0 {
				expDetails.ExpiryDateTimezone = utils.StringPtr(v1.(string))
			}
			if v1, ok1 := val["delete_database"]; ok1 {
				expDetails.DeleteDatabase = v1.(bool)
			}
		}
		return expDetails
	}
	return nil
}

func expandDBRefreshDetails(pr []interface{}) *era.DBRefreshDetails {
	if len(pr) > 0 {
		refDetails := &era.DBRefreshDetails{}

		for _, v := range pr {
			val := v.(map[string]interface{})

			if v1, ok1 := val["refresh_in_days"]; ok1 {
				refDetails.RefreshInDays = v1.(int)
			}
			if v1, ok1 := val["refresh_time"]; ok1 && len(v1.(string)) > 0 {
				refDetails.RefreshTime = v1.(string)
			}
			if v1, ok1 := val["refresh_date_timezone"]; ok1 && len(v1.(string)) > 0 {
				refDetails.RefreshDateTimezone = v1.(string)
			}
		}
		return refDetails
	}
	return nil
}
