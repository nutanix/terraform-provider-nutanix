package nutanix

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-nutanix/client/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixNDBProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBProfileCreate,
		ReadContext:   resourceNutanixNDBProfileRead,
		UpdateContext: resourceNutanixNDBProfileUpdate,
		DeleteContext: resourceNutanixNDBProfileDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"published": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"compute_profile": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"software_profile", "network_profile", "database_parameter_profile"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpus": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  1,
						},
						"core_per_cpu": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  1,
						},
						"memory_size": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  2,
						},
					},
				},
			},
			"software_profile": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"compute_profile", "network_profile", "database_parameter_profile"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topology": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"cluster", "single"}, false),
						},
						"postgres_database": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_dbserver_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"base_profile_version_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"base_profile_version_description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"os_notes": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"db_software_notes": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"available_cluster_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"network_profile": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"compute_profile", "software_profile", "database_parameter_profile"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topology": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"cluster", "single"}, false),
						},
						"postgres_database": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"single_instance": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vlan_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"enable_ip_address_selection": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"ha_instance": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vlan_name": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"cluster_name": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"cluster_id": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"num_of_clusters": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"version_cluster_association": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nx_cluster_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"database_parameter_profile": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"compute_profile", "software_profile", "network_profile"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"postgres_database": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_connections": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "100",
										Description: "Determines the maximum number of concurrent connections to the database server. The default is typically 100, but might be less if your kernel settings will not support it (as determined during initdb).",
									},
									"max_replication_slots": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "10",
										Description: "Specifies the maximum number of replication slots that the server can support. The default is zero. wal_level must be set to archive or higher to allow replication slots to be used. Setting it to a lower value than the number of currently existing replication slots will prevent the server from starting.",
									},
									"effective_io_concurrency": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "1",
										Description: "Sets the number of concurrent disk I/O operations that PostgreSQL expects can be executed simultaneously. Raising this value will increase the number of I/O operations that any individual PostgreSQL session attempts to initiate in parallel.",
									},
									"timezone": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "UTC",
										Description: "Sets the time zone for displaying and interpreting time stamps",
									},
									"max_prepared_transactions": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "0",
										Description: "Sets the maximum number of transactions that can be in the prepared state simultaneously. Setting this parameter to zero (which is the default) disables the prepared-transaction feature. If you are not planning to use prepared transactions, this parameter should be set to zero to prevent accidental creation of prepared transactions. If you are using prepared transactions, you will probably want max_prepared_transactions to be at least as large as max_connections, so that every session can have a prepared transaction pending.",
									},
									"max_locks_per_transaction": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "64",
										Description: "The shared lock table tracks locks on max_locks_per_transaction * (max_connections + max_prepared_transactions) objects (e.g., tables); hence, no more than this many distinct objects can be locked at any one time. This parameter controls the average number of object locks allocated for each transaction; individual transactions can lock more objects as long as the locks of all transactions fit in the lock table. This is not the number of rows that can be locked; that value is unlimited. The default, 64, has historically proven sufficient, but you might need to raise this value if you have clients that touch many different tables in a single transaction. Increasing this parameter might cause PostgreSQL to request more System V shared memory than your operating system's default configuration allows.",
									},
									"max_wal_senders": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "10",
										Description: "Specifies the maximum number of concurrent connections from standby servers or streaming base backup clients (i.e., the maximum number of simultaneously running WAL sender processes). The default is 10. The value 0 means replication is disabled. WAL sender processes count towards the total number of connections, so the parameter cannot be set higher than max_connections. Abrupt streaming client disconnection might cause an orphaned connection slot until a timeout is reached, so this parameter should be set slightly higher than the maximum number of expected clients so disconnected clients can immediately reconnect. wal_level must be set to replica or higher to allow connections from standby servers.",
									},
									"max_worker_processes": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "8",
										Description: "Sets the maximum number of background processes that the system can support. The default is 8. When running a standby server, you must set this parameter to the same or higher value than on the master server. Otherwise, queries will not be allowed in the standby server.",
									},
									"min_wal_size": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "80MB",
										Description: "As long as WAL disk usage stays below this setting, old WAL files are always recycled for future use at a checkpoint, rather than removed. This can be used to ensure that enough WAL space is reserved to handle spikes in WAL usage, for example when running large batch jobs. The default is 80 MB.",
									},
									"max_wal_size": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "1GB",
										Description: "Maximum size to let the WAL grow to between automatic WAL checkpoints. This is a soft limit; WAL size can exceed max_wal_size under special circumstances, like under heavy load, a failing archive_command, or a high wal_keep_segments setting. The default is 1 GB. Increasing this parameter can increase the amount of time needed for crash recovery.",
									},
									"checkpoint_timeout": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "5min",
										Description: "Sets the maximum time between automatic WAL checkpoints . High Value gives Good Performance, but takes More Recovery Time, Reboot time. can reduce the I/O load on your system, especially when using large values for shared_buffers.",
									},
									"autovacuum": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "on",
										Description: "Controls whether the server should run the autovacuum launcher daemon. This is on by default; however, track_counts must also be enabled for autovacuum to work.",
									},
									"checkpoint_completion_target": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "0.5",
										Description: "Specifies the target of checkpoint completion, as a fraction of total time between checkpoints. Time spent flushing dirty buffers during checkpoint, as fraction of checkpoint interval . Formula - (checkpoint_timeout - 2min) / checkpoint_timeout. The default is 0.5.",
									},
									"autovacuum_freeze_max_age": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "200000000",
										Description: "Age at which to autovacuum a table to prevent transaction ID wraparound",
									},
									"autovacuum_vacuum_threshold": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "50",
										Description: "Min number of row updates before vacuum. Minimum number of tuple updates or deletes prior to vacuum. Take value in KB",
									},
									"autovacuum_vacuum_scale_factor": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "0.2",
										Description: "Number of tuple updates or deletes prior to vacuum as a fraction of reltuples",
									},
									"autovacuum_work_mem": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "-1",
										Description: "Sets the maximum memory to be used by each autovacuum worker process. Unit is in KB",
									},
									"autovacuum_max_workers": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "3",
										Description: "Sets the maximum number of simultaneously running autovacuum worker processes",
									},
									"autovacuum_vacuum_cost_delay": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "2ms",
										Description: "Vacuum cost delay in milliseconds, for autovacuum. Specifies the cost delay value that will be used in automatic VACUUM operations",
									},
									"wal_buffers": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "-1",
										Description: "Sets the number of disk-page buffers in shared memory for WAL. The amount of shared memory used for WAL data that has not yet been written to disk. The default setting of -1 selects a size equal to 1/32nd (about 3%) of shared_buffers, but not less than 64kB nor more than the size of one WAL segment, typically 16MB",
									},
									"synchronous_commit": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "on",
										Description: "Sets the current transaction's synchronization level. Specifies whether transaction commit will wait for WAL records to be written to disk before the command returns a success indication to the client. https://www.postgresql.org/docs/12/runtime-config-wal.html#GUC-SYNCHRONOUS-COMMIT",
									},
									"random_page_cost": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "4",
										Description: "Sets the planner's estimate of the cost of a nonsequentially fetched disk page. Sets the planner's estimate of the cost of a non-sequentially-fetched disk page. The default is 4.0. This value can be overridden for tables and indexes in a particular tablespace by setting the tablespace",
									},
									"wal_keep_segments": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "700",
										Description: "Sets the number of WAL files held for standby servers, Specifies the minimum number of past log file segments kept in the pg_wal directory, in case a standby server needs to fetch them for streaming replication. Each segment is normally 16 megabytes.",
									},
								},
							},
						},
					},
				},
			},

			// computed arguments
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_profile": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"published": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deprecated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secure": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"properties_map": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"version_cluster_association": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nx_cluster_id": {
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
									"owner_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"profile_version_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"secure": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"optimized_for_provisioning": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assoc_databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_availability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nx_cluster_id": {
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
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixNDBProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	req := &era.ProfileRequest{}
	softwareProfile := false
	if name, ok := d.GetOk("name"); ok {
		req.Name = utils.StringPtr(name.(string))
	}

	if des, ok := d.GetOk("description"); ok {
		req.Description = utils.StringPtr(des.(string))
	}

	if engType, ok := d.GetOk("engine_type"); ok {
		req.EngineType = utils.StringPtr(engType.(string))
	}

	if cp, ok := d.GetOk("compute_profile"); ok {
		req.Properties = buildComputeProfileRequest(cp)
		// setting some defaults values which are generated at runtime
		req.Topology = utils.StringPtr("ALL")
		req.Type = utils.StringPtr("Compute")
		req.SystemProfile = false
		req.DBVersion = utils.StringPtr("ALL")
	}

	if np, ok := d.GetOk("network_profile"); ok {
		nps := np.([]interface{})

		for _, v := range nps {
			val := v.(map[string]interface{})

			if tp, ok := val["topology"]; ok {
				req.Topology = utils.StringPtr(tp.(string))

				// other details
				req.Type = utils.StringPtr("Network")
				req.SystemProfile = false
				req.DBVersion = utils.StringPtr("ALL")
			}

			if ps, ok := val["postgres_database"]; ok {
				req.Properties = expandNetworkProfileProperties(ps.([]interface{}), ctx, meta)
			}

			if cls, ok := val["version_cluster_association"]; ok {
				clster := cls.([]interface{})
				out := make([]*era.VersionClusterAssociation, len(clster))
				for _, v := range clster {
					val := v.(map[string]interface{})

					if p1, ok1 := val["nx_cluster_id"]; ok1 {
						out = append(out, &era.VersionClusterAssociation{
							NxClusterID: utils.StringPtr(p1.(string)),
						})
					}
				}
				req.VersionClusterAssociation = out
			}
		}
	}

	if db, ok := d.GetOk("database_parameter_profile"); ok {
		req.Properties = buildDatabaseProfileProperties(db.([]interface{}))

		// setting some defaults values which are generated at runtime
		req.Topology = utils.StringPtr("ALL")
		req.Type = utils.StringPtr("Database_Parameter")
		req.SystemProfile = false
		req.DBVersion = utils.StringPtr("ALL")
	}

	if sp, ok := d.GetOk("software_profile"); ok {
		softwareProfile = true
		splist := sp.([]interface{})

		for _, v := range splist {
			val := v.(map[string]interface{})

			if tp, ok := val["topology"]; ok {
				req.Topology = utils.StringPtr(tp.(string))

				// other details
				req.Type = utils.StringPtr("Software")
				req.SystemProfile = false
				req.DBVersion = utils.StringPtr("ALL")
			}

			if ps, ok := val["postgres_database"]; ok {
				req.Properties = expandSoftwareProfileProp(ps.([]interface{}))
			}

			if ac, ok1 := d.GetOk("available_cluster_ids"); ok1 {
				st := ac.([]interface{})
				sublist := make([]*string, len(st))

				for a := range st {
					sublist[a] = utils.StringPtr(st[a].(string))
				}
				req.AvailableClusterIds = sublist
			}
		}
	}

	if softwareProfile {
		resp, er := conn.Service.CreateSoftwareProfiles(ctx, req)
		if er != nil {
			return diag.FromErr(er)
		}

		// Get Operation ID from response of SoftwareProfile  and poll for the operation to get completed.
		opID := resp.OperationId
		if opID == utils.StringPtr("") {
			return diag.Errorf("error: operation ID is an empty string")
		}
		opReq := era.GetOperationRequest{
			OperationID: utils.StringValue(opID),
		}

		log.Printf("polling for operation with id: %s\n", *opID)

		// Poll for operation here - Operation GET Call
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING"},
			Target:  []string{"COMPLETED", "FAILED"},
			Refresh: eraRefresh(ctx, conn, opReq),
			Timeout: d.Timeout(schema.TimeoutCreate),
			Delay:   eraDelay,
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for software profile	 (%s) to create: %s", *resp.EntityId, errWaitTask)
		}
		d.SetId(*resp.EntityId)
		return resourceNutanixNDBProfileRead(ctx, d, meta)
	}

	resp, err := conn.Service.CreateProfiles(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	return resourceNutanixNDBProfileRead(ctx, d, meta)
}

func resourceNutanixNDBProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	resp, err := conn.Service.GetProfiles(ctx, "", "", d.Id(), "")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner", resp.Owner); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("latest_version", resp.Latestversion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("versions", flattenVersions(resp.Versions)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("latest_version_id", resp.Latestversionid); err != nil {
		return diag.FromErr(err)
	}

	if resp.Assocdbservers != nil {
		d.Set("assoc_db_servers", resp.Assocdbservers)
	} else {
		d.Set("assoc_db_servers", nil)
	}

	if resp.Assocdatabases != nil {
		d.Set("assoc_databases", resp.Assocdatabases)
	} else {
		d.Set("assoc_databases", nil)
	}

	if resp.Nxclusterid != nil {
		d.Set("nx_cluster_id", resp.Nxclusterid)
	} else {
		d.Set("nx_cluster_id", nil)
	}

	if resp.Clusteravailability != nil {
		d.Set("cluster_availability", flattenClusterAvailability(resp.Clusteravailability))
	} else {
		d.Set("cluster_availability", nil)
	}

	return nil
}

func resourceNutanixNDBProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	req := &era.ProfileRequest{}

	netReq := &era.UpdateProfileRequest{}

	res, err := conn.Service.GetProfiles(ctx, "", "", d.Id(), "")
	if err != nil {
		diag.FromErr(err)
	}

	if res != nil {
		netReq.Name = res.Name
		netReq.Description = res.Description
		req.Properties = res.Versions[0].Properties
	}

	if pub, ok := d.GetOk("published"); ok {
		req.Published = pub.(bool)
	}

	if d.HasChange("name") {
		netReq.Name = utils.StringPtr(d.Get("name").(string))
		// update version name as well
		versionName := d.Get("name").(string)
		updateVersionName := versionName + " " + " (1.0)"
		req.Name = utils.StringPtr(updateVersionName)
	}

	if d.HasChange("description") {
		netReq.Description = utils.StringPtr(d.Get("description").(string))
		req.Description = utils.StringPtr(d.Get("description").(string))
	}

	if d.HasChange("compute_profile") {
		req.Properties = buildComputeProfileRequest(d.Get("compute_profile"))
	}

	if d.HasChange("network_profile") {
		nps := d.Get("network_profile").([]interface{})

		for _, v := range nps {
			val := v.(map[string]interface{})

			if ps, ok := val["postgres_database"]; ok {
				req.Properties = expandNetworkProfileProperties(ps.([]interface{}), ctx, meta)
			}

			if cls, ok := val["version_cluster_association"]; ok {
				clster := cls.([]interface{})
				out := make([]*era.VersionClusterAssociation, len(clster))
				for _, v := range clster {
					val := v.(map[string]interface{})

					if p1, ok1 := val["nx_cluster_id"]; ok1 {
						out = append(out, &era.VersionClusterAssociation{
							NxClusterID: utils.StringPtr(p1.(string)),
						})
					}
				}
				req.VersionClusterAssociation = out
			}
		}
	}

	if d.HasChange("database_parameter_profile") {
		req.Properties = buildDatabaseProfileProperties(d.Get("database_parameter_profile").([]interface{}))
	}

	if d.HasChange("software_profile") {
		splist := d.Get("software_profile").([]interface{})

		for _, v := range splist {
			val := v.(map[string]interface{})

			if ps, ok := val["postgres_database"]; ok {
				req.Properties = expandSoftwareProfileProp(ps.([]interface{}))
			}

			if ac, ok1 := d.GetOk("available_cluster_ids"); ok1 {
				st := ac.([]interface{})
				sublist := make([]*string, len(st))

				for a := range st {
					sublist[a] = utils.StringPtr(st[a].(string))
				}
				req.AvailableClusterIds = sublist
			}
		}
	}

	version_id := res.Versions[0].ID

	_, eror := conn.Service.UpdateProfile(ctx, netReq, d.Id())
	if eror != nil {
		return diag.FromErr(eror)
	}

	_, er := conn.Service.UpdateProfileVersion(ctx, req, d.Id(), *version_id)
	if er != nil {
		return diag.FromErr(er)
	}

	return resourceNutanixNDBProfileRead(ctx, d, meta)
}

func resourceNutanixNDBProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	resp, err := conn.Service.DeleteProfile(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == utils.StringPtr("Profile Successfully Deleted.") {
		d.SetId("")
	}
	return nil
}

func buildComputeProfileRequest(p interface{}) []*era.ProfileProperties {
	if p != nil {
		computeProp := []*era.ProfileProperties{}
		pc := p.([]interface{})
		for _, v := range pc {
			val := v.(map[string]interface{})
			if cpu, ok := val["cpus"]; ok {
				computeProp = append(computeProp, &era.ProfileProperties{
					Name:   utils.StringPtr("CPUS"),
					Value:  utils.StringPtr(cpu.(string)),
					Secure: false,
				})
			}

			if core_cpu, ok := val["core_per_cpu"]; ok {
				computeProp = append(computeProp, &era.ProfileProperties{
					Name:   utils.StringPtr("CORE_PER_CPU"),
					Value:  utils.StringPtr(core_cpu.(string)),
					Secure: false,
				})
			}

			if mem, ok := val["memory_size"]; ok {
				computeProp = append(computeProp, &era.ProfileProperties{
					Name:   utils.StringPtr("MEMORY_SIZE"),
					Value:  utils.StringPtr(mem.(string)),
					Secure: false,
				})
			}
		}
		return computeProp
	}
	return nil
}

func expandNetworkProfileProperties(ps []interface{}, ctx context.Context, meta interface{}) []*era.ProfileProperties {
	prop := []*era.ProfileProperties{}

	if len(ps) > 0 {
		for _, v := range ps {
			inst := v.(map[string]interface{})

			fmt.Println("Hellooooooooo")
			if sIns, ok := inst["single_instance"]; ok && len(sIns.([]interface{})) > 0 {
				fmt.Println("SINGLEEEEEEEEEE")
				prop = expandNetworkSingleInstance(sIns.([]interface{}))
			}

			if hIns, ok := inst["ha_instance"]; ok && len(hIns.([]interface{})) > 0 {
				prop = expandNetworkHAInstance(hIns.([]interface{}), ctx, meta)
			}
		}
	}
	return prop
}

func buildDatabaseProfileProperties(ps []interface{}) []*era.ProfileProperties {
	prop := []*era.ProfileProperties{}

	if len(ps) > 0 {

		for _, v := range ps {
			val := v.(map[string]interface{})
			if psdb, ok := val["postgres_database"]; ok {
				brr := psdb.([]interface{})

				for _, v := range brr {
					val := v.(map[string]interface{})

					if p1, ok1 := val["max_connections"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("max_connections"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Determines the maximum number of concurrent connections to the database server. The default is typically 100, but might be less if your kernel settings will not support it (as determined during initdb)."),
						})
					}

					if p1, ok1 := val["max_replication_slots"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("max_replication_slots"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Specifies the maximum number of replication slots that the server can support. The default is zero. wal_level must be set to archive or higher to allow replication slots to be used. Setting it to a lower value than the number of currently existing replication slots will prevent the server from starting."),
						})
					}
					if p1, ok1 := val["effective_io_concurrency"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("effective_io_concurrency"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the number of concurrent disk I/O operations that PostgreSQL expects can be executed simultaneously. Raising this value will increase the number of I/O operations that any individual PostgreSQL session attempts to initiate in parallel."),
						})
					}
					if p1, ok1 := val["timezone"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("timezone"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the time zone for displaying and interpreting time stamps"),
						})
					}
					if p1, ok1 := val["max_prepared_transactions"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("max_prepared_transactions"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the maximum number of transactions that can be in the prepared state simultaneously. Setting this parameter to zero (which is the default) disables the prepared-transaction feature. If you are not planning to use prepared transactions, this parameter should be set to zero to prevent accidental creation of prepared transactions. If you are using prepared transactions, you will probably want max_prepared_transactions to be at least as large as max_connections, so that every session can have a prepared transaction pending."),
						})
					}
					if p1, ok1 := val["max_locks_per_transaction"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("max_locks_per_transaction"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr(" The shared lock table tracks locks on max_locks_per_transaction * (max_connections + max_prepared_transactions) objects (e.g., tables); hence, no more than this many distinct objects can be locked at any one time. This parameter controls the average number of object locks allocated for each transaction; individual transactions can lock more objects as long as the locks of all transactions fit in the lock table. This is not the number of rows that can be locked; that value is unlimited. The default, 64, has historically proven sufficient, but you might need to raise this value if you have clients that touch many different tables in a single transaction. Increasing this parameter might cause PostgreSQL to request more System V shared memory than your operating system's default configuration allows."),
						})
					}
					if p1, ok1 := val["max_wal_senders"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("max_wal_senders"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Specifies the maximum number of concurrent connections from standby servers or streaming base backup clients (i.e., the maximum number of simultaneously running WAL sender processes). The default is 10. The value 0 means replication is disabled. WAL sender processes count towards the total number of connections, so the parameter cannot be set higher than max_connections. Abrupt streaming client disconnection might cause an orphaned connection slot until a timeout is reached, so this parameter should be set slightly higher than the maximum number of expected clients so disconnected clients can immediately reconnect. wal_level must be set to replica or higher to allow connections from standby servers."),
						})
					}
					if p1, ok1 := val["max_worker_processes"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("max_worker_processes"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the maximum number of background processes that the system can support. The default is 8. When running a standby server, you must set this parameter to the same or higher value than on the master server. Otherwise, queries will not be allowed in the standby server."),
						})
					}
					if p1, ok1 := val["min_wal_size"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("min_wal_size"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("As long as WAL disk usage stays below this setting, old WAL files are always recycled for future use at a checkpoint, rather than removed. This can be used to ensure that enough WAL space is reserved to handle spikes in WAL usage, for example when running large batch jobs. The default is 80 MB."),
						})
					}
					if p1, ok1 := val["max_wal_size"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("max_wal_size"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Maximum size to let the WAL grow to between automatic WAL checkpoints. This is a soft limit; WAL size can exceed max_wal_size under special circumstances, like under heavy load, a failing archive_command, or a high wal_keep_segments setting. The default is 1 GB. Increasing this parameter can increase the amount of time needed for crash recovery."),
						})
					}
					if p1, ok1 := val["checkpoint_timeout"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("checkpoint_timeout"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the maximum time between automatic WAL checkpoints . High Value gives Good Performance, but takes More Recovery Time, Reboot time. can reduce the I/O load on your system, especially when using large values for shared_buffers."),
						})
					}
					if p1, ok1 := val["autovacuum"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("autovacuum"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Controls whether the server should run the autovacuum launcher daemon. This is on by default; however, track_counts must also be enabled for autovacuum to work."),
						})
					}
					if p1, ok1 := val["checkpoint_completion_target"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("checkpoint_completion_target"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Specifies the target of checkpoint completion, as a fraction of total time between checkpoints. Time spent flushing dirty buffers during checkpoint, as fraction of checkpoint interval . Formula - (checkpoint_timeout - 2min) / checkpoint_timeout. The default is 0.5."),
						})
					}
					if p1, ok1 := val["autovacuum_freeze_max_age"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("autovacuum_freeze_max_age"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Age at which to autovacuum a table to prevent transaction ID wraparound"),
						})
					}
					if p1, ok1 := val["autovacuum_vacuum_threshold"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("autovacuum_vacuum_threshold"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Min number of row updates before vacuum. Minimum number of tuple updates or deletes prior to vacuum. Take value in KB"),
						})
					}
					if p1, ok1 := val["autovacuum_vacuum_scale_factor"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("autovacuum_vacuum_scale_factor"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Number of tuple updates or deletes prior to vacuum as a fraction of reltuples"),
						})
					}
					if p1, ok1 := val["autovacuum_work_mem"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("autovacuum_work_mem"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the maximum memory to be used by each autovacuum worker process. Unit is in KB"),
						})
					}
					if p1, ok1 := val["autovacuum_max_workers"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("autovacuum_max_workers"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the maximum number of simultaneously running autovacuum worker processes"),
						})
					}
					if p1, ok1 := val["autovacuum_vacuum_cost_delay"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("autovacuum_vacuum_cost_delay"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Vacuum cost delay in milliseconds, for autovacuum. Specifies the cost delay value that will be used in automatic VACUUM operations"),
						})
					}
					if p1, ok1 := val["wal_buffers"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("wal_buffers"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the number of disk-page buffers in shared memory for WAL. The amount of shared memory used for WAL data that has not yet been written to disk. The default setting of -1 selects a size equal to 1/32nd (about 3%) of shared_buffers, but not less than 64kB nor more than the size of one WAL segment, typically 16MB"),
						})
					}

					if p1, ok1 := val["synchronous_commit"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("synchronous_commit"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the current transaction's synchronization level. Specifies whether transaction commit will wait for WAL records to be written to disk before the command returns a success indication to the client. https://www.postgresql.org/docs/12/runtime-config-wal.html#GUC-SYNCHRONOUS-COMMIT"),
						})
					}
					if p1, ok1 := val["random_page_cost"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("random_page_cost"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the planner's estimate of the cost of a nonsequentially fetched disk page. Sets the planner's estimate of the cost of a non-sequentially-fetched disk page. The default is 4.0. This value can be overridden for tables and indexes in a particular tablespace by setting the tablespace"),
						})
					}
					if p1, ok1 := val["wal_keep_segments"]; ok1 {
						prop = append(prop, &era.ProfileProperties{
							Name:        utils.StringPtr("wal_keep_segments"),
							Value:       utils.StringPtr(p1.(string)),
							Secure:      false,
							Description: utils.StringPtr("Sets the number of WAL files held for standby servers, Specifies the minimum number of past log file segments kept in the pg_wal directory, in case a standby server needs to fetch them for streaming replication. Each segment is normally 16 megabytes."),
						})
					}
				}
			}
		}

	}

	return prop
}

func expandSoftwareProfileProp(ps []interface{}) []*era.ProfileProperties {
	prop := []*era.ProfileProperties{}

	if len(ps) > 0 {
		for _, v := range ps {
			val := v.(map[string]interface{})

			if p1, ok1 := val["source_dbserver_id"]; ok1 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("SOURCE_DBSERVER_ID"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("ID of the database server that should be used as a reference to create the software profile"),
				})
			}
			if p1, ok1 := val["base_profile_version_name"]; ok1 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("BASE_PROFILE_VERSION_NAME"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("Name of the base profile version."),
				})
			}
			if p1, ok1 := val["base_profile_version_description"]; ok1 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("BASE_PROFILE_VERSION_DESCRIPTION"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("Description of the base profile version."),
				})
			}
			if p1, ok1 := val["os_notes"]; ok1 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("OS_NOTES"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("Notes or description for the Operating System."),
				})
			}
			if p1, ok1 := val["db_software_notes"]; ok1 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("DB_SOFTWARE_NOTES"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("Description of the Postgres database software."),
				})
			}
		}
		return prop
	}
	return nil
}

func expandNetworkSingleInstance(ps []interface{}) []*era.ProfileProperties {
	if len(ps) > 0 {
		prop := []*era.ProfileProperties{}

		for _, v := range ps {
			val := v.(map[string]interface{})

			if p1, ok1 := val["vlan_name"]; ok1 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("VLAN_NAME"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("Name of the vLAN"),
				})
			}

			if p1, ok1 := val["enable_ip_address_selection"]; ok1 {
				prop = append(prop, &era.ProfileProperties{
					Name:  utils.StringPtr("ENABLE_IP_ADDRESS_SELECTION"),
					Value: utils.StringPtr(p1.(string)),
				})
			}
		}
		fmt.Println("INSIDE NETWORKINGGGGGGG")
		return prop
	}
	return nil
}

func expandNetworkHAInstance(ps []interface{}, ctx context.Context, meta interface{}) []*era.ProfileProperties {
	prop := []*era.ProfileProperties{}

	for _, v := range ps {
		val := v.(map[string]interface{})

		if numCls, ok := val["num_of_clusters"]; ok {
			prop = append(prop, &era.ProfileProperties{
				Name:  utils.StringPtr("NUM_CLUSTERS"),
				Value: utils.StringPtr(numCls.(string)),
			})
		}

		if p1, ok1 := val["enable_ip_address_selection"]; ok1 {
			prop = append(prop, &era.ProfileProperties{
				Name:  utils.StringPtr("ENABLE_IP_ADDRESS_SELECTION"),
				Value: utils.StringPtr(p1.(string)),
			})
		}

		if numVlan, ok := val["vlan_name"]; ok {

			vlans := numVlan.([]interface{})

			for k, vl := range vlans {
				prop = append(prop, &era.ProfileProperties{
					Name:  utils.StringPtr(fmt.Sprintf("VLAN_NAME_%d", k)),
					Value: utils.StringPtr(vl.(string)),
				})
			}
		}

		if clsName, ok := val["cluster_name"]; ok && len(clsName.([]interface{})) > 0 {

			vlans := clsName.([]interface{})

			for k, vl := range vlans {
				prop = append(prop, &era.ProfileProperties{
					Name:  utils.StringPtr(fmt.Sprintf("CLUSTER_NAME_%d", k)),
					Value: utils.StringPtr(vl.(string)),
				})

				// call the cluster API to fetch cluster id
				conn := meta.(*Client).Era
				resp, _ := conn.Service.GetCluster(ctx, "", vl.(string))

				prop = append(prop, &era.ProfileProperties{
					Name:  utils.StringPtr(fmt.Sprintf("CLUSTER_ID_%d", k)),
					Value: utils.StringPtr(*resp.ID),
				})
			}
		}

		if clsId, ok := val["cluster_id"]; ok && len(clsId.([]interface{})) > 0 {

			vlans := clsId.([]interface{})

			for k, vl := range vlans {
				prop = append(prop, &era.ProfileProperties{
					Name:  utils.StringPtr(fmt.Sprintf("CLUSTER_ID_%d", k)),
					Value: utils.StringPtr(vl.(string)),
				})

				conn := meta.(*Client).Era
				resp, _ := conn.Service.GetCluster(ctx, vl.(string), "")

				prop = append(prop, &era.ProfileProperties{
					Name:  utils.StringPtr(fmt.Sprintf("CLUSTER_NAME_%d", k)),
					Value: utils.StringPtr(*resp.Uniquename),
				})
			}
		}
	}
	return prop
}
