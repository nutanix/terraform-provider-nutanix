package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	Era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixEraDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraDatabaseRead,
		Schema: map[string]*schema.Schema{
			"database_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
			"owner_id": {
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
			"properties": {
				Type:        schema.TypeList,
				Description: "List of all the properties",
				Computed:    true,
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
						"ref_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secure": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clustered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"clone": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"era_created": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"internal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"placeholder": {
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
						"bpg_configs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"data_disks": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"count": {
																Type:     schema.TypeFloat,
																Computed: true,
															},
														},
													},
												},
												"log_disks": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"count": {
																Type:     schema.TypeFloat,
																Computed: true,
															},
															"size": {
																Type:     schema.TypeFloat,
																Computed: true,
															},
														},
													},
												},
												"archive_storage": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"size": {
																Type:     schema.TypeFloat,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"vm_properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"nr_hugepages": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"overcommit_memory": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"dirty_ratio": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"dirty_background_ratio": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"dirty_expire_centisecs": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"dirty_writeback_centisecs": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"swappiness": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
											},
										},
									},
									"bpg_db_param": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"shared_buffers": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"maintenance_work_mem": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"work_mem": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"effective_cache_size": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"max_worker_processes": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"max_parallel_workers_per_gather": {
													Type:     schema.TypeString,
													Computed: true,
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
			"group_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metadata": {
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
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_database_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_source_database_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lcm_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiry_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"remind_before_in_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"effective_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expiry_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expiry_date_timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_created": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"expire_in_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"delete_database": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"delete_time_machine": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"delete_vm": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"refresh_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"refresh_in_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"refresh_in_hours": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"refresh_in_months": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"last_refresh_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"next_refresh_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"refresh_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"refresh_date_timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"pre_delete_command": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"post_delete_command": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"time_machine": {
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
						"owner_id": {
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
						"access_level": {
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
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ref_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secure": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"clustered": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"clone": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"internal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"database_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ea_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sla_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clones": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_nx_clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sla_update_in_progress": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"metric": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sla_update_metadata": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sla": {
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
									"unique_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"owner_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"system_sla": {
										Type:     schema.TypeBool,
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

									"continuous_retention": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"daily_retention": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"weekly_retention": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"monthly_retention": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"quarterly_retention": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"yearly_retention": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"reference_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"pitr_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"current_active_frequency": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"schedule": {
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
									"unique_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"owner_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"system_policy": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"global_policy": {
										Type:     schema.TypeBool,
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
									"snapshot_time_of_day": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"hours": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"minutes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"extra": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"continuous_schedule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"log_backup_interval": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"snapshots_per_day": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"weekly_schedule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"day_of_week": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"day_of_week_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"monthly_schedule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"day_of_month": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"yearly_schedule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"month": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"month_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"day_of_month": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"quartely_schedule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_month": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"start_month_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"day_of_month": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"daily_schedule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"reference_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"dbserver_logical_cluster": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"database_nodes": {
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
						"owner_id": {
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
						"access_level": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"database_id": {
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
						"primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dbserver_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"software_installation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
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
									"secure_info": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"info": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"dbserver": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"protection_domain": {
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
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_id": {
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
									"primary_host": {
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
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ref_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"secure": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"description": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"era_created": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"assoc_entities": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"linked_databases": {
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
						"database_name": {
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
						"owner_id": {
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
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"metadata": {
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
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"databases": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"database_group_state_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNutanixEraDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era
	dUUID, ok := d.GetOk("database_id")
	if !ok {
		return diag.Errorf("please provide `database_id`")
	}

	resp, err := conn.Service.GetDatabaseInstance(ctx, dUUID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner_id", resp.Ownerid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_created", resp.Datecreated); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_modified", resp.Datemodified); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("properties", flattenDbInstanceProperties(resp.Properties)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", flattenDbTags(resp.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clone", resp.Clone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clustered", resp.Clustered); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("internal", resp.Internal); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("placeholder", resp.Placeholder); err != nil {
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

	if err := d.Set("database_status", resp.Databasestatus); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_logical_cluster_id", resp.Dbserverlogicalclusterid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_machine_id", resp.Timemachineid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_time_machine_id", resp.Parenttimemachineid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_zone", resp.Timezone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("info", flattenDbInfo(resp.Info)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("group_info", resp.GroupInfo); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("metadata", flattenDbInstanceMetadata(resp.Metadata)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("metric", resp.Metric); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("category", resp.Category); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_database_id", resp.ParentDatabaseId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_source_database_id", resp.ParentSourceDatabaseId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("lcm_config", flattenDbLcmConfig(resp.Lcmconfig)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_machine", flattenDbTimeMachine(resp.TimeMachine)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_logical_cluster", resp.Dbserverlogicalcluster); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_nodes", flattenDbNodes(resp.Databasenodes)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("linked_databases", flattenDbLinkedDbs(resp.Linkeddatabases)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("databases", resp.Databases); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database_group_state_info", resp.DatabaseGroupStateInfo); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	return nil
}

func flattenDbInstanceProperties(pr []Era.DBInstanceProperties) []map[string]interface{} {
	if len(pr) > 0 {
		res := []map[string]interface{}{}
		for _, v := range pr {
			prop := map[string]interface{}{}

			prop["description"] = v.Description
			prop["name"] = v.Name
			prop["ref_id"] = v.RefID
			prop["secure"] = v.Secure
			prop["value"] = v.Value

			res = append(res, prop)
		}
		return res
	}
	return nil
}

func flattenDbInstanceMetadata(pr *Era.DBInstanceMetadata) map[string]interface{} {
	if pr != nil {
		pmeta := make(map[string]interface{})

		pmeta["secure_info"] = pr.Secureinfo
		pmeta["info"] = pr.Info
		pmeta["deregister_info"] = pr.Deregisterinfo
		pmeta["tm_activate_operation_id"] = pr.Tmactivateoperationid
		pmeta["created_dbservers"] = pr.Createddbservers
		pmeta["registered_dbservers"] = pr.Registereddbservers
		pmeta["last_refresh_timestamp"] = pr.Lastrefreshtimestamp
		pmeta["last_requested_refresh_timestamp"] = pr.Lastrequestedrefreshtimestamp
		pmeta["capability_reset_time"] = pr.CapabilityResetTime
		pmeta["state_before_refresh"] = pr.Statebeforerefresh
		pmeta["state_before_restore"] = pr.Statebeforerestore
		pmeta["state_before_scaling"] = pr.Statebeforescaling
		pmeta["log_catch_up_for_restore_dispatched"] = pr.Logcatchupforrestoredispatched
		pmeta["last_log_catch_up_for_restore_operation_id"] = pr.Lastlogcatchupforrestoreoperationid
		pmeta["base_size_computed"] = pr.BaseSizeComputed
		pmeta["original_database_name"] = pr.Originaldatabasename
		pmeta["provision_operation_id"] = pr.ProvisionOperationId
		pmeta["source_snapshot_id"] = pr.SourceSnapshotId
		pmeta["pitr_based"] = pr.PitrBased
		pmeta["sanitised"] = pr.Sanitised
		pmeta["refresh_blocker_info"] = pr.RefreshBlockerInfo
		pmeta["deregistered_with_delete_time_machine"] = pr.DeregisteredWithDeleteTimeMachine
	}
	return nil
}

func flattenDbNodes(pr []Era.Databasenodes) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			db := map[string]interface{}{}

			db["access_level"] = v.AccessLevel
			db["database_id"] = v.Databaseid
			db["database_status"] = v.Databasestatus
			db["date_created"] = v.Datecreated
			db["date_modified"] = v.Datemodified
			db["dbserver_id"] = v.Dbserverid
			db["description"] = v.Description
			db["id"] = v.ID
			db["metadata"] = v.Metadata
			db["name"] = v.Name
			db["owner_id"] = v.Ownerid
			db["primary"] = v.Primary
			db["properties"] = v.Properties
			db["protection_domain"] = flattenDbProtectionDomain(v.Protectiondomain)
			db["protection_domain_id"] = v.Protectiondomainid
			db["software_installation_id"] = v.Softwareinstallationid
			db["status"] = v.Status
			db["tags"] = flattenDbTags(v.Tags)

			res[k] = db
		}
		return res
	}
	return nil
}

func flattenDbLinkedDbs(pr []Era.Linkeddatabases) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ld := map[string]interface{}{}

			ld["database_name"] = v.DatabaseName
			ld["database_status"] = v.Databasestatus
			ld["date_created"] = v.Datecreated
			ld["date_modified"] = v.Datemodified
			ld["description"] = v.Description
			ld["id"] = v.ID
			ld["metadata"] = v.Metadata
			ld["metric"] = v.Metric
			ld["name"] = v.Name
			ld["owner_id"] = v.Ownerid
			ld["parent_database_id"] = v.ParentDatabaseId
			ld["parent_linked_database_id"] = v.ParentLinkedDatabaseId
			ld["snapshot_id"] = v.SnapshotId
			ld["status"] = v.Status
			ld["timezone"] = v.TimeZone

			res[k] = ld
		}
		return res
	}
	return nil
}

func flattenDbProtectionDomain(pr *Era.Protectiondomain) []map[string]interface{} {
	pDList := make([]map[string]interface{}, 0)
	if pr != nil {
		pmeta := make(map[string]interface{}, 0)

		pmeta["cloud_id"] = pr.Cloudid
		pmeta["date_created"] = pr.Datecreated
		pmeta["date_modified"] = pr.Datemodified
		pmeta["description"] = pr.Description
		pmeta["era_created"] = pr.Eracreated
		pmeta["id"] = pr.ID
		pmeta["name"] = pr.Name
		pmeta["owner_id"] = pr.Ownerid
		pmeta["primary_host"] = pr.PrimaryHost
		pmeta["properties"] = flattenDbInstanceProperties(pr.Properties)
		pmeta["status"] = pr.Status
		if pr.Tags != nil {
			pmeta["tags"] = flattenDbTags(pr.Tags)
		}
		pmeta["type"] = pr.Type

		pDList = append(pDList, pmeta)
		return pDList
	}
	return nil
}

func flattenDbTags(pr []Era.EraTags) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			tag := map[string]interface{}{}

			tag["entity_id"] = v.EntityId
			tag["entity_name"] = v.EntityType
			tag["id"] = v.TagId
			tag["tag_name"] = v.TagName
			tag["value"] = v.Value

			res[k] = tag
		}
		return res
	}
	return nil
}

func flattenDbInfo(pr *Era.Info) []map[string]interface{} {
	infoList := make([]map[string]interface{}, 0)
	if pr != nil {
		info := make(map[string]interface{}, 0)

		if pr.Secureinfo != nil {
			info["secure_info"] = pr.Secureinfo
		}
		if pr.Info != nil {
			info["bpg_configs"] = flattenBpgConfig(pr.Info.BpgConfigs)
		}
		infoList = append(infoList, info)
		return infoList
	}
	return nil
}

func flattenBpgConfig(pr *Era.BpgConfigs) []map[string]interface{} {
	bpgList := make([]map[string]interface{}, 0)
	if pr != nil {
		bpg := make(map[string]interface{}, 0)

		var bgdbParams []map[string]interface{}
		if pr.BpgDbParam != nil {
			bg := make(map[string]interface{})
			bg["maintenance_work_mem"] = utils.StringValue(&pr.BpgDbParam.MaintenanceWorkMem)
			bg["effective_cache_size"] = utils.StringValue(&pr.BpgDbParam.EffectiveCacheSize)
			bg["max_parallel_workers_per_gather"] = utils.StringValue(&pr.BpgDbParam.MaxParallelWorkersPerGather)
			bg["max_worker_processes"] = utils.StringValue(&pr.BpgDbParam.MaxWorkerProcesses)
			bg["shared_buffers"] = utils.StringValue(&pr.BpgDbParam.SharedBuffers)
			bg["work_mem"] = utils.StringValue(&pr.BpgDbParam.WorkMem)
			bgdbParams = append(bgdbParams, bg)
		}
		bpg["bpg_db_param"] = bgdbParams

		var storg []map[string]interface{}
		if pr.Storage != nil {
			str := make(map[string]interface{})

			var storgArch []map[string]interface{}
			if pr.Storage.ArchiveStorage != nil {
				arc := make(map[string]interface{})

				arc["size"] = pr.Storage.ArchiveStorage.Size
				storgArch = append(storgArch, arc)
			}
			str["archive_storage"] = storgArch

			var stdisk []map[string]interface{}
			if pr.Storage.DataDisks != nil {
				arc := make(map[string]interface{})

				arc["count"] = pr.Storage.DataDisks.Count
				stdisk = append(stdisk, arc)
			}
			str["data_disks"] = stdisk

			var stgLog []map[string]interface{}
			if pr.Storage.LogDisks != nil {
				arc := make(map[string]interface{})

				arc["size"] = pr.Storage.LogDisks.Size
				arc["count"] = pr.Storage.LogDisks.Count
				stgLog = append(stgLog, arc)
			}
			str["log_disks"] = stgLog

			storg = append(storg, str)
		}
		bpg["storage"] = storg

		var vmProp []map[string]interface{}
		if pr.VMProperties != nil {
			vmp := make(map[string]interface{})
			vmp["dirty_background_ratio"] = pr.VMProperties.DirtyBackgroundRatio
			vmp["dirty_expire_centisecs"] = pr.VMProperties.DirtyExpireCentisecs
			vmp["dirty_ratio"] = pr.VMProperties.DirtyRatio
			vmp["dirty_writeback_centisecs"] = pr.VMProperties.DirtyWritebackCentisecs
			vmp["nr_hugepages"] = pr.VMProperties.NrHugepages
			vmp["overcommit_memory"] = pr.VMProperties.OvercommitMemory
			vmp["swappiness"] = pr.VMProperties.Swappiness

			vmProp = append(vmProp, vmp)
		}

		bpg["vm_properties"] = vmProp

		bpgList = append(bpgList, bpg)
		return bpgList
	}
	return nil
}

func flattenDbLcmConfig(pr *Era.LcmConfig) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		lcm := map[string]interface{}{}

		lcm["expiryDetails"] = flattenEraExpiryDetails(pr.ExpiryDetails)
		lcm["refresh_details"] = flattenEraRefreshDetails(pr.RefreshDetails)

		var preLcmComm []map[string]interface{}
		if pr.PreDeleteCommand != nil {
			pre := map[string]interface{}{}

			pre["command"] = pr.PreDeleteCommand.Command

			preLcmComm = append(preLcmComm, pre)

		}
		lcm["pre_delete_command"] = preLcmComm

		var postLcmComm []map[string]interface{}
		if pr.PreDeleteCommand != nil {
			pre := map[string]interface{}{}

			pre["command"] = pr.PostDeleteCommand.Command

			postLcmComm = append(postLcmComm, pre)

		}
		lcm["post_delete_command"] = postLcmComm

		res = append(res, lcm)
		return res
	}
	return nil
}

func flattenEraExpiryDetails(pr *Era.EraDbExpiryDetails) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		expiry := map[string]interface{}{}

		expiry["delete_database"] = pr.DeleteDatabase
		expiry["delete_time_machine"] = pr.DeleteTimeMachine
		expiry["delete_vm"] = pr.DeleteVM
		expiry["effective_timestamp"] = pr.EffectiveTimestamp
		expiry["expire_in_days"] = pr.ExpireInDays
		expiry["expiry_date_timezone"] = pr.ExpiryDateTimezone
		expiry["expiry_timestamp"] = pr.ExpiryTimestamp
		expiry["remind_before_in_days"] = pr.RemindBeforeInDays
		expiry["user_created"] = pr.UserCreated

		res = append(res, expiry)
		return res
	}
	return nil
}

func flattenEraRefreshDetails(pr *Era.EraDbRefreshDetails) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		refresh := map[string]interface{}{}

		refresh["last_refresh_date"] = pr.LastRefreshDate
		refresh["next_refresh_date"] = pr.NextRefreshDate
		refresh["refresh_date_timezone"] = pr.RefreshDateTimezone
		refresh["refresh_in_days"] = pr.RefreshInDays
		refresh["refresh_in_hours"] = pr.RefreshInHours
		refresh["refresh_in_months"] = pr.RefreshInMonths
		refresh["refresh_time"] = pr.RefreshTime

		res = append(res, refresh)
		return res
	}
	return nil
}

func flattenDbTimeMachine(pr *Era.TimeMachine) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		tmac := map[string]interface{}{}

		tmac["id"] = pr.ID
		tmac["name"] = pr.Name
		tmac["description"] = pr.Description
		tmac["owner_id"] = pr.OwnerID
		tmac["date_created"] = pr.DateCreated
		tmac["date_modified"] = pr.DateModified
		tmac["access_level"] = pr.AccessLevel
		tmac["properties"] = flattenDbInstanceProperties(pr.Properties)
		tmac["tags"] = flattenDbTags(pr.Tags)
		tmac["clustered"] = pr.Clustered
		tmac["clone"] = pr.Clone
		tmac["internal"] = pr.Internal
		tmac["database_id"] = pr.DatabaseID
		tmac["type"] = pr.Type
		tmac["category"] = pr.Category
		tmac["status"] = pr.Status
		tmac["ea_status"] = pr.EaStatus
		tmac["scope"] = pr.Scope
		tmac["sla_id"] = pr.SLAID
		tmac["schedule_id"] = pr.ScheduleID
		tmac["metric"] = pr.Metric
		tmac["sla_update_metadata"] = pr.SLAUpdateMetadata
		tmac["database"] = pr.Database
		tmac["clones"] = pr.Clones
		tmac["source_nx_clusters"] = pr.SourceNxClusters
		tmac["sla_update_in_progress"] = pr.SLAUpdateInProgress
		tmac["sla"] = flattenDbSla(pr.SLA)
		tmac["schedule"] = flattenSchedule(pr.Schedule)

		res = append(res, tmac)
		return res

	}
	return nil
}

func flattenDbSla(pr *Era.ListSLAResponse) []map[string]interface{} {
	res := []map[string]interface{}{}
	if pr != nil {
		sla := map[string]interface{}{}

		sla["id"] = pr.ID
		sla["name"] = pr.Name
		sla["continuous_retention"] = pr.Continuousretention
		sla["daily_retention"] = pr.Dailyretention
		sla["date_modified"] = pr.Datemodified
		sla["date_created"] = pr.Datecreated
		sla["description"] = pr.Description
		sla["monthly_retention"] = pr.Monthlyretention
		sla["owner_id"] = pr.Ownerid
		sla["quarterly_retention"] = pr.Quarterlyretention
		sla["reference_count"] = pr.Referencecount
		sla["system_sla"] = pr.Systemsla
		sla["unique_name"] = pr.Uniquename
		sla["weekly_retention"] = pr.Weeklyretention
		sla["yearly_retention"] = pr.Yearlyretention

		res = append(res, sla)
		return res
	}
	return nil
}

func flattenSchedule(pr *Era.Schedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		sch := map[string]interface{}{}

		sch["continuous_schedule"] = flattenContinousSch(pr.Continuousschedule)
		sch["date_created"] = pr.Datecreated
		sch["date_modified"] = pr.Datemodified
		sch["description"] = pr.Description
		sch["global_policy"] = pr.GlobalPolicy
		sch["id"] = pr.ID
		sch["monthly_schedule"] = flattenMonthlySchedule(pr.Monthlyschedule)
		sch["name"] = pr.Name
		sch["owner_id"] = pr.OwnerID
		sch["quartely_schedule"] = flattenQuartelySchedule(pr.Quartelyschedule)
		sch["reference_count"] = pr.ReferenceCount
		sch["snapshot_time_of_day"] = flattenSnapshotTimeOfDay(pr.Snapshottimeofday)
		sch["start_time"] = pr.StartTime
		sch["system_policy"] = pr.SystemPolicy
		sch["time_zone"] = pr.TimeZone
		sch["unique_name"] = pr.UniqueName
		sch["weekly_schedule"] = flattenWeeklySchedule(pr.Weeklyschedule)
		sch["yearly_schedule"] = flattenYearlylySchedule(pr.Yearlyschedule)
		sch["daily_schedule"] = flattenDailySchedule(pr.Dailyschedule)

		res = append(res, sch)
		return res
	}
	return nil
}

func flattenContinousSch(pr *Era.Continuousschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["log_backup_interval"] = pr.Logbackupinterval
		cr["snapshots_per_day"] = pr.Snapshotsperday

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenMonthlySchedule(pr *Era.Monthlyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["day_of_month"] = pr.Dayofmonth

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenQuartelySchedule(pr *Era.Quartelyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["day_of_month"] = pr.Dayofmonth
		cr["start_month"] = pr.Startmonth

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenSnapshotTimeOfDay(pr *Era.Snapshottimeofday) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["hours"] = pr.Hours
		cr["minutes"] = pr.Minutes
		cr["seconds"] = pr.Seconds

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenWeeklySchedule(pr *Era.Weeklyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["day_of_week"] = pr.Dayofweek

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenYearlylySchedule(pr *Era.Yearlyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["day_of_month"] = pr.Dayofmonth
		cr["month"] = pr.Month

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenDailySchedule(pr *Era.Dailyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		res = append(res, cr)
		return res
	}
	return nil
}
