package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixNDBMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBMaintenanceWindowRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recurrence": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"day_of_week": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"week_of_month": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"threshold": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hour": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"minute": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
					},
				},
			},
			"tags": dataSourceEraDBInstanceTags(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_run_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_task_assoc": EntityTaskAssocSchema(),
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixNDBMaintenanceWindowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	maintainenanceWindowID := d.Get("id")

	resp, err := conn.Service.ReadMaintenanceWindow(ctx, maintainenanceWindowID.(string))
	if err != nil {
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

	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("next_run_time", resp.NextRunTime); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entity_task_assoc", flattenEntityTaskAssoc(resp.EntityTaskAssoc)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("timezone", resp.Timezone); err != nil {
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

	if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schedule", flattenMaintenanceSchedule(resp.Schedule)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(maintainenanceWindowID.(string))
	return nil
}

func flattenEntityTaskAssoc(pr []*era.MaintenanceTasksResponse) []interface{} {
	if len(pr) > 0 {
		tasks := make([]interface{}, 0)

		for _, v := range pr {
			entity := map[string]interface{}{}

			entity["access_level"] = v.AccessLevel
			entity["date_created"] = v.DateCreated
			entity["date_modified"] = v.DateModified
			entity["description"] = v.Description
			entity["entity"] = v.Entity
			entity["entity_id"] = v.EntityID
			entity["entity_type"] = v.EntityType
			entity["id"] = v.ID
			entity["maintenance_window_id"] = v.MaintenanceWindowID
			entity["maintenance_window_owner_id"] = v.MaintenanceWindowOwnerID
			entity["name"] = v.Name
			entity["owner_id"] = v.OwnerID
			entity["payload"] = flattenEntityTaskPayload(v.Payload)
			entity["status"] = v.Status
			entity["task_type"] = v.TaskType

			if v.Tags != nil {
				entity["tags"] = flattenDBTags(v.Tags)
			}

			if v.Properties != nil {
				props := []interface{}{}
				for _, prop := range v.Properties {
					props = append(props, map[string]interface{}{
						"name":  prop.Name,
						"value": prop.Value,
					})
				}
				entity["properties"] = props
			}

			tasks = append(tasks, entity)
		}
		return tasks
	}
	return nil
}

func flattenEntityTaskPayload(pr *era.Payload) []interface{} {
	if pr != nil {
		res := make([]interface{}, 0)

		payload := map[string]interface{}{}

		payload["pre_post_command"] = flattenPrePostCommand(pr.PrePostCommand)
		res = append(res, payload)

		return res
	}
	return nil
}

func flattenPrePostCommand(pr *era.PrePostCommand) []interface{} {
	if pr != nil {
		comms := make([]interface{}, 0)
		command := map[string]interface{}{}

		command["post_command"] = pr.PostCommand
		command["pre_command"] = pr.PreCommand

		comms = append(comms, command)
		return comms
	}
	return nil
}

func EntityTaskAssocSchema() *schema.Schema {
	return &schema.Schema{
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
						},
					},
				},
				"tags": dataSourceEraDBInstanceTags(),
				"maintenance_window_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"maintenance_window_owner_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"entity_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"entity_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"task_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"payload": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"pre_post_command": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"pre_command": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"post_command": {
											Type:     schema.TypeString,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
				"entity": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
