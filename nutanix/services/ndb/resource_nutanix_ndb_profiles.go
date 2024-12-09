package ndb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBProfileCreate,
		ReadContext:   resourceNutanixNDBProfileRead,
		UpdateContext: resourceNutanixNDBProfileUpdate,
		DeleteContext: resourceNutanixNDBProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Computed: true,
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
							Default:  "1",
						},
						"core_per_cpu": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "1",
						},
						"memory_size": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "2",
						},
					},
				},
			},
			"software_profile": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
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
										Type:     schema.TypeString,
										Optional: true,
										Default:  "100",
									},
									"max_replication_slots": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "10",
									},
									"effective_io_concurrency": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "1",
									},
									"timezone": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "UTC",
									},
									"max_prepared_transactions": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "0",
									},
									"max_locks_per_transaction": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "64",
									},
									"max_wal_senders": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "10",
									},
									"max_worker_processes": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "8",
									},
									"min_wal_size": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "80MB",
									},
									"max_wal_size": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "1GB",
									},
									"checkpoint_timeout": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "5min",
									},
									"autovacuum": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "on",
									},
									"checkpoint_completion_target": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "0.5",
									},
									"autovacuum_freeze_max_age": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "200000000",
									},
									"autovacuum_vacuum_threshold": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "50",
									},
									"autovacuum_vacuum_scale_factor": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "0.2",
									},
									"autovacuum_work_mem": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "-1",
									},
									"autovacuum_max_workers": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "3",
									},
									"autovacuum_vacuum_cost_delay": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "2ms",
									},
									"wal_buffers": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "-1",
									},
									"synchronous_commit": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "on",
									},
									"random_page_cost": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "4",
									},
									"wal_keep_segments": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "700",
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
			"assoc_db_servers": {
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
	conn := meta.(*conns.Client).Era

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
				req.Properties = expandNetworkProfileProperties(ctx, meta, ps.([]interface{}))
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
		opID := resp.OperationID
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
			return diag.Errorf("error waiting for software profile	 (%s) to create: %s", *resp.EntityID, errWaitTask)
		}
		d.SetId(*resp.EntityID)
	} else {
		resp, err := conn.Service.CreateProfiles(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(*resp.ID)
	}

	// Now if published is present args

	if publish, ok := d.GetOk("published"); ok {
		req := &era.ProfileRequest{}
		netReq := &era.UpdateProfileRequest{}

		req.Published = publish.(bool)

		// profile filter spec
		profileFilter := &era.ProfileFilter{}
		profileFilter.ProfileID = d.Id()
		res, err := conn.Service.GetProfile(ctx, profileFilter)
		if err != nil {
			diag.FromErr(err)
		}

		if res != nil {
			netReq.Name = res.Name
			netReq.Description = res.Description
			req.Properties = res.Versions[0].Properties
		}
		versionID := res.Versions[0].ID

		_, eror := conn.Service.UpdateProfile(ctx, netReq, d.Id())
		if eror != nil {
			return diag.FromErr(eror)
		}

		_, er := conn.Service.UpdateProfileVersion(ctx, req, d.Id(), *versionID)
		if er != nil {
			return diag.FromErr(er)
		}
	}
	log.Printf("NDB Profile with %s id is created successfully", d.Id())
	return resourceNutanixNDBProfileRead(ctx, d, meta)
}

func resourceNutanixNDBProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// profile filter spec
	profileFilter := &era.ProfileFilter{}
	profileFilter.ProfileID = d.Id()

	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("id is required for read operation")
	}
	resp, err := conn.Service.GetProfile(ctx, profileFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("engine_type", resp.Enginetype); err != nil {
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
	conn := meta.(*conns.Client).Era

	req := &era.ProfileRequest{}

	netReq := &era.UpdateProfileRequest{}

	// profile filter spec
	profileFilter := &era.ProfileFilter{}
	profileFilter.ProfileID = d.Id()

	res, err := conn.Service.GetProfile(ctx, profileFilter)
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
				req.Properties = expandNetworkProfileProperties(ctx, meta, ps.([]interface{}))
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

	versionID := res.Versions[0].ID

	_, eror := conn.Service.UpdateProfile(ctx, netReq, d.Id())
	if eror != nil {
		return diag.FromErr(eror)
	}

	_, er := conn.Service.UpdateProfileVersion(ctx, req, d.Id(), *versionID)
	if er != nil {
		return diag.FromErr(er)
	}
	log.Printf("NDB Profile with %s id is updated successfully", d.Id())
	return resourceNutanixNDBProfileRead(ctx, d, meta)
}

func resourceNutanixNDBProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.DeleteProfile(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == utils.StringPtr("Profile Successfully Deleted.") {
		log.Printf("NDB Profile with %s id is deleted successfully", d.Id())
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

			if coreCPU, ok := val["core_per_cpu"]; ok {
				computeProp = append(computeProp, &era.ProfileProperties{
					Name:   utils.StringPtr("CORE_PER_CPU"),
					Value:  utils.StringPtr(coreCPU.(string)),
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

func expandNetworkProfileProperties(ctx context.Context, meta interface{}, ps []interface{}) []*era.ProfileProperties {
	prop := []*era.ProfileProperties{}
	if len(ps) > 0 {
		for _, v := range ps {
			inst := v.(map[string]interface{})

			if sIns, ok := inst["single_instance"]; ok && len(sIns.([]interface{})) > 0 {
				prop = expandNetworkSingleInstance(sIns.([]interface{}))
			}

			if hIns, ok := inst["ha_instance"]; ok && len(hIns.([]interface{})) > 0 {
				prop = expandNetworkHAInstance(ctx, meta, hIns.([]interface{}))
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

				postgresProp := brr[0].(map[string]interface{})
				for key, value := range postgresProp {
					prop = append(prop, &era.ProfileProperties{
						Name:   utils.StringPtr(key),
						Value:  utils.StringPtr(value.(string)),
						Secure: false,
					})
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

			if p1, ok1 := val["source_dbserver_id"]; ok1 && len(p1.(string)) > 0 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("SOURCE_DBSERVER_ID"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("ID of the database server that should be used as a reference to create the software profile"),
				})
			}
			if p1, ok1 := val["base_profile_version_name"]; ok1 && len(p1.(string)) > 0 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("BASE_PROFILE_VERSION_NAME"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("Name of the base profile version."),
				})
			}
			if p1, ok1 := val["base_profile_version_description"]; ok1 && len(p1.(string)) > 0 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("BASE_PROFILE_VERSION_DESCRIPTION"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("Description of the base profile version."),
				})
			}
			if p1, ok1 := val["os_notes"]; ok1 && len(p1.(string)) > 0 {
				prop = append(prop, &era.ProfileProperties{
					Name:        utils.StringPtr("OS_NOTES"),
					Value:       utils.StringPtr(p1.(string)),
					Secure:      false,
					Description: utils.StringPtr("Notes or description for the Operating System."),
				})
			}
			if p1, ok1 := val["db_software_notes"]; ok1 && len(p1.(string)) > 0 {
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
		return prop
	}
	return nil
}

func expandNetworkHAInstance(ctx context.Context, meta interface{}, ps []interface{}) []*era.ProfileProperties {
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
				conn := meta.(*conns.Client).Era
				resp, _ := conn.Service.GetCluster(ctx, "", vl.(string))

				prop = append(prop, &era.ProfileProperties{
					Name:  utils.StringPtr(fmt.Sprintf("CLUSTER_ID_%d", k)),
					Value: utils.StringPtr(*resp.ID),
				})
			}
		}

		if clsID, ok := val["cluster_id"]; ok && len(clsID.([]interface{})) > 0 {
			vlans := clsID.([]interface{})
			for k, vl := range vlans {
				prop = append(prop, &era.ProfileProperties{
					Name:  utils.StringPtr(fmt.Sprintf("CLUSTER_ID_%d", k)),
					Value: utils.StringPtr(vl.(string)),
				})

				conn := meta.(*conns.Client).Era
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
