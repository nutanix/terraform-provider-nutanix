package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNDBTmsCapability() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBTmsCapabilityRead,
		Schema: map[string]*schema.Schema{
			"time_machine_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nx_cluster_association_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sla_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"overall_continuous_range_end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_continuous_snapshot_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_catchup_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"heal_with_reset_capability": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"database_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// check data schema later
			"log_time_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"capability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"to": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_unit_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"snapshots": {
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
									"properties": dataSourceEraDatabaseProperties(),
									"tags":       dataSourceEraDBInstanceTags(),
									"snapshot_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nx_cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protection_domain_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parent_snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"time_machine_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"database_node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"app_info_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"applicable_types": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"snapshot_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"metadata": {
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
												"info": dataSourceEraDatabaseInfo(),
												"deregister_info": {
													Type:     schema.TypeMap,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"from_timestamp": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"to_timestamp": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"replication_retry_count": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"last_replication_retyr_source_snapshot_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"async": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"stand_by": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"curation_retry_count": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"operations_using_snapshot": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"software_snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"software_database_snapshot": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"dbserver_storage_metadata_version": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"santized": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"santized_from_snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"processed": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"database_snapshot": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"from_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"to_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dbserver_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dbserver_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dbserver_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"replicated_snapshots": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"software_snapshot": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"santized_snapshots": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_family": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_timestamp_date": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"lcm_config": dataSourceEraLCMConfig(),
									"parent_snapshot": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"snapshot_size": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"continuous_region": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"to_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sub_range": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"unknown_time_ranges": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"processed_ranges": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"first": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"second": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"unprocessed_ranges": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"first": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"second": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"partial_ranges": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"time_range_and_databases": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshots": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_logs": {
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
												"era_log_drive_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"database_node_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"from_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"to_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"size": {
													Type:     schema.TypeInt,
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
															"unknown_time_range": {
																Type:     schema.TypeBool,
																Computed: true,
															},
														},
													},
												},
												"metadata": {
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
															"deregister_info": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"message": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"operations": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"curation_retry_count": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"created_directly": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"updated_directly": {
																Type:     schema.TypeBool,
																Computed: true,
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
												"owner_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"database_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"unprocessed": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"log_copy_operation_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"databases_continuous_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"capability_reset_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_db_log": {
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
						"era_log_drive_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"to_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"metadata": {
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
									"deregister_info": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operations": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"curation_retry_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"created_directly": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"updated_directly": {
										Type:     schema.TypeBool,
										Computed: true,
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
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unprocessed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"log_copy_operation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"last_continuous_snapshot": {
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
						"properties": dataSourceEraDatabaseProperties(),
						"tags":       dataSourceEraDBInstanceTags(),
						"snapshot_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nx_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_machine_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_info_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"applicable_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"snapshot_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
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
									"info": dataSourceEraDatabaseInfo(),
									"deregister_info": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"from_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"to_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"replication_retry_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"last_replication_retry_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_replication_retry_source_snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"async": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"stand_by": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"curation_retry_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"operations_using_snapshot": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"software_snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"software_database_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dbserver_storage_metadata_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"santized": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"santized_from_snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"processed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"database_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"from_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"to_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dbserver_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dbserver_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dbserver_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replicated_snapshots": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"software_snapshot": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"santized_snapshots": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_timestamp_date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lcm_config": dataSourceEraLCMConfig(),
						"parent_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"snapshot_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixNDBTmsCapabilityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	tmsID := d.Get("time_machine_id")
	resp, er := conn.Service.TimeMachineCapability(ctx, tmsID.(string))
	if er != nil {
		return diag.FromErr(er)
	}

	if err := d.Set("output_time_zone", resp.OutputTimeZone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nx_cluster_id", resp.NxClusterID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("source", resp.Source); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nx_cluster_association_type", resp.NxClusterAssociationType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("sla_id", resp.SLAID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("overall_continuous_range_end_time", resp.OverallContinuousRangeEndTime); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("last_continuous_snapshot_time", resp.LastContinuousSnapshotTime); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("log_catchup_start_time", resp.LogCatchupStartTime); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("heal_with_reset_capability", resp.HealWithResetCapability); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_ids", utils.StringValueSlice(resp.DatabaseIds)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("capability", flattenTmsCapability(resp.Capability)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("capability_reset_time", resp.CapabilityResetTime); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("last_db_log", flattenLastDBLog(resp.LastDBLog)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("last_continuous_snapshot", flattenLastContinousSnapshot(resp.LastContinuousSnapshot)); err != nil {
		return diag.FromErr(err)
	}
	uuid, e := uuid.GenerateUUID()

	if e != nil {
		return diag.Errorf("Error generating UUID for era tms capability: %+v", e)
	}
	d.SetId(uuid)
	return nil
}

func flattenTmsCapability(pr []*era.Capability) []map[string]interface{} {
	if len(pr) > 0 {
		tmsList := []map[string]interface{}{}

		for _, v := range pr {
			cap := map[string]interface{}{}

			cap["continuous_region"] = flattenContinousRegion(v.ContinuousRegion)
			cap["database_ids"] = utils.StringValueSlice(v.DatabaseIds)
			cap["databases_continuous_region"] = v.DatabasesContinuousRegion
			cap["from"] = v.From
			cap["mode"] = v.Mode
			cap["snapshots"] = flattenSnapshotsList(v.Snapshots)
			cap["time_unit"] = v.TimeUnit
			cap["time_unit_number"] = v.TimeUnitNumber
			cap["to"] = v.To

			tmsList = append(tmsList, cap)
		}
		return tmsList
	}
	return nil
}

func flattenContinousRegion(pr *era.ContinuousRegion) []map[string]interface{} {
	if pr != nil {
		continousRegion := make([]map[string]interface{}, 0)
		conReg := map[string]interface{}{}

		conReg["from_time"] = pr.FromTime
		conReg["to_time"] = pr.ToTime
		conReg["sub_range"] = pr.SubRange
		conReg["message"] = pr.Message
		conReg["snapshot_ids"] = utils.StringSlice(pr.SnapshotIds)
		conReg["unknown_time_ranges"] = pr.UnknownTimeRanges
		conReg["processed_ranges"] = flattenProcessedRanges(pr.ProcessedRanges)
		conReg["unprocessed_ranges"] = flattenProcessedRanges(pr.UnprocessedRanges)
		conReg["partial_ranges"] = pr.PartialRanges
		conReg["time_range_and_databases"] = pr.TimeRangeAndDatabases
		conReg["snapshots"] = pr.Snapshots
		conReg["db_logs"] = flattenDBLogs(pr.DBLogs)
		conReg["timezone"] = pr.TimeZone

		continousRegion = append(continousRegion, conReg)
		return continousRegion
	}
	return nil
}

func flattenDBLogs(pr []*era.DBLogs) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))

		for _, v := range pr {
			val := map[string]interface{}{}

			val["id"] = v.ID
			val["name"] = v.Name
			val["era_log_drive_id"] = v.EraLogDriveID
			val["database_node_id"] = v.DatabaseNodeID
			val["from_time"] = v.FromTime
			val["to_time"] = v.ToTime
			val["status"] = v.Status
			val["size"] = v.Size
			val["metadata"] = flattenDBLogMetadata(v.Metadata)
			val["date_created"] = v.DateCreated
			val["date_modified"] = v.DateModified
			val["owner_id"] = v.OwnerID
			val["database_id"] = v.DatabaseID
			val["message"] = v.Message
			val["unprocessed"] = v.Unprocessed
			val["log_copy_operation_id"] = v.LogCopyOperationID

			res = append(res, val)
		}
		return res
	}
	return nil
}

func flattenDBLogMetadata(pr *era.DBLogsMetadata) []map[string]interface{} {
	if pr != nil {
		logsMeta := make([]map[string]interface{}, 0)
		log := map[string]interface{}{}

		log["secure_info"] = pr.SecureInfo
		log["info"] = pr.Info
		log["deregister_info"] = flattenDeRegiserInfo(pr.DeregisterInfo)
		log["curation_retry_count"] = pr.CurationRetryCount
		log["created_directly"] = pr.CreatedDirectly
		log["updated_directly"] = pr.UpdatedDirectly

		logsMeta = append(logsMeta, log)
		return logsMeta
	}
	return nil
}

func flattenLastDBLog(pr *era.DBLogs) []map[string]interface{} {
	if pr != nil {
		res := make([]map[string]interface{}, 0)
		val := map[string]interface{}{}

		val["id"] = pr.ID
		val["name"] = pr.Name
		val["era_log_drive_id"] = pr.EraLogDriveID
		val["database_node_id"] = pr.DatabaseNodeID
		val["from_time"] = pr.FromTime
		val["to_time"] = pr.ToTime
		val["status"] = pr.Status
		val["size"] = pr.Size
		val["metadata"] = flattenDBLogMetadata(pr.Metadata)
		val["date_created"] = pr.DateCreated
		val["date_modified"] = pr.DateModified
		val["owner_id"] = pr.OwnerID
		val["database_id"] = pr.DatabaseID
		val["message"] = pr.Message
		val["unprocessed"] = pr.Unprocessed
		val["log_copy_operation_id"] = pr.LogCopyOperationID

		res = append(res, val)
		return res
	}
	return nil
}

func flattenLastContinousSnapshot(pr *era.LastContinuousSnapshot) []map[string]interface{} {
	if pr != nil {
		snpList := make([]map[string]interface{}, 0)
		snap := map[string]interface{}{}

		snap["id"] = pr.ID
		snap["name"] = pr.Name
		snap["description"] = pr.Description
		snap["owner_id"] = pr.OwnerID
		snap["date_created"] = pr.DateCreated
		snap["date_modified"] = pr.DateModified
		snap["properties"] = flattenDBInstanceProperties(pr.Properties)
		snap["tags"] = flattenDBTags(pr.Tags)
		snap["snapshot_uuid"] = pr.SnapshotUUID
		snap["nx_cluster_id"] = pr.NxClusterID
		snap["protection_domain_id"] = pr.ProtectionDomainID
		snap["parent_snapshot_id"] = pr.ParentSnapshotID
		snap["time_machine_id"] = pr.TimeMachineID
		snap["database_node_id"] = pr.DatabaseNodeID
		snap["app_info_version"] = pr.AppInfoVersion
		snap["status"] = pr.Status
		snap["type"] = pr.Type
		snap["applicable_types"] = pr.ApplicableTypes
		snap["snapshot_timestamp"] = pr.SnapshotTimeStamp
		snap["metadata"] = flattenLastContinousSnapshotMetadata(pr.Metadata)
		snap["software_snapshot_id"] = pr.SoftwareSnapshotID
		snap["software_database_snapshot"] = pr.SoftwareDatabaseSnapshot
		snap["santized_from_snapshot_id"] = pr.SanitizedFromSnapshotID
		snap["processed"] = pr.Processed
		snap["database_snapshot"] = pr.DatabaseSnapshot
		snap["from_timestamp"] = pr.FromTimeStamp
		snap["to_timestamp"] = pr.ToTimeStamp
		snap["dbserver_id"] = pr.DBserverID
		snap["dbserver_name"] = pr.DBserverName
		snap["dbserver_ip"] = pr.DBserverIP
		snap["replicated_snapshots"] = pr.ReplicatedSnapshots
		snap["software_snapshot"] = pr.SoftwareSnapshot
		snap["santized_snapshots"] = pr.SanitizedSnapshots
		snap["snapshot_family"] = pr.SnapshotFamily
		snap["snapshot_timestamp_date"] = pr.SnapshotTimeStampDate
		snap["lcm_config"] = flattenDBLcmConfig(pr.LcmConfig)
		snap["parent_snapshot"] = pr.ParentSnapshot
		snap["snapshot_size"] = pr.SnapshotSize

		snpList = append(snpList, snap)
		return snpList
	}
	return nil
}

func flattenLastContinousSnapshotMetadata(pr *era.LastContinuousSnapshotMetadata) []map[string]interface{} {
	if pr != nil {
		res := make([]map[string]interface{}, 0)

		meta := map[string]interface{}{}

		meta["secure_info"] = pr.SecureInfo
		meta["info"] = pr.Info
		meta["deregister_info"] = pr.DeregisterInfo
		meta["from_timestamp"] = pr.FromTimeStamp
		meta["to_timestamp"] = pr.ToTimeStamp
		meta["replication_retry_count"] = pr.ReplicationRetryCount
		meta["last_replication_retry_timestamp"] = pr.LastReplicationRetryTimestamp
		meta["last_replication_retry_source_snapshot_id"] = pr.LastReplicationRetrySourceSnapshotID
		meta["async"] = pr.Async
		meta["stand_by"] = pr.Standby
		meta["curation_retry_count"] = pr.CurationRetryCount
		meta["operations_using_snapshot"] = pr.OperationsUsingSnapshot

		res = append(res, meta)
		return res
	}
	return nil
}

func flattenProcessedRanges(pr []*era.ProcessedRanges) []interface{} {
	if len(pr) > 0 {
		res := make([]interface{}, len(pr))

		for _, v := range pr {
			proRanges := map[string]interface{}{}

			proRanges["first"] = v.First
			proRanges["second"] = v.Second

			res = append(res, proRanges)
		}
		return res
	}
	return nil
}
