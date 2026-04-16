package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func DataSourceNutanixNDBTimeMachine() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBTimeMachineRead,
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
			"properties": dataSourceEraDatabaseProperties(),
			"tags":       dataSourceEraDBInstanceTags(),
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
	}
}

func dataSourceNutanixNDBTimeMachineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	tmsID, tok := d.GetOk("time_machine_id")
	tmsName, tnOk := d.GetOk("time_machine_name")

	if !tok && !tnOk {
		return diag.Errorf("Atleast one of time_machine_id or time_machine_name is required to perform clone")
	}

	// call time Machine API

	resp, err := conn.Service.GetTimeMachine(ctx, tmsID.(string), tmsName.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner_id", resp.OwnerID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("date_created", resp.DateCreated); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("date_modified", resp.DateModified); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("access_level", resp.AccessLevel); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("properties", flattenDBInstanceProperties(resp.Properties)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clustered", resp.Clustered); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clone", resp.Clone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("internal", resp.Internal); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_id", resp.DatabaseID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("category", resp.Category); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ea_status", resp.EaStatus); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clustered", resp.Clustered); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clone", resp.Clone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("scope", resp.Scope); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("sla_id", resp.SLAID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schedule_id", resp.ScheduleID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("metric", resp.Metric); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database", resp.Database); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clones", resp.Clones); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_nx_clusters", resp.SourceNxClusters); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("sla_update_in_progress", resp.SLAUpdateInProgress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sla", flattenDBSLA(resp.SLA)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schedule", flattenSchedule(resp.Schedule)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	return nil
}
