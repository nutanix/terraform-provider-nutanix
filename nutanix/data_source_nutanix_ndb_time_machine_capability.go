package nutanix

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixNDBTmsCapability() *schema.Resource {
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
									"santised": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"santised_from_snapshot_id": {
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
									"santised_snapshots": {
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
							Type:     schema.TypeString,
							Computed: true,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_continuous_snapshot": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixNDBTmsCapabilityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	tmsID := d.Get("time_machine_id")
	resp, err := conn.Service.TimeMachineCapability(ctx, tmsID.(string))
	if err != nil {
		return diag.FromErr(err)
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

	if err := d.Set("last_db_log", resp.LastDbLog); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("last_continuous_snapshot", resp.LastContinuousSnapshot); err != nil {
		return diag.FromErr(err)
	}
	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era tms capability: %+v", err)
	}
	d.SetId(uuid)
	return nil
}

func flattenTmsCapability(pr []*era.Capability) []map[string]interface{} {
	if len(pr) > 0 {
		tmsList := []map[string]interface{}{}

		for _, v := range pr {
			cap := map[string]interface{}{}

			cap["continuous_region"] = v.ContinuousRegion
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
