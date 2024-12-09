package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixNDBMaintenanceWindows() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBMaintenanceWindowsRead,
		Schema: map[string]*schema.Schema{
			"maintenance_windows": {
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
				},
			},
		},
	}
}

func dataSourceNutanixNDBMaintenanceWindowsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.ListMaintenanceWindow(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if e := d.Set("maintenance_windows", flattenMaintenanceWindowsResponse(resp)); err != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era clusters: %+v", err)
	}
	d.SetId(uuid)

	return nil
}

func flattenMaintenanceWindowsResponse(pr *era.ListMaintenanceWindowResponse) []interface{} {
	if pr != nil {
		windowResp := make([]interface{}, 0)
		for _, v := range *pr {
			window := map[string]interface{}{}
			window["id"] = v.ID
			window["name"] = v.Name
			window["description"] = v.Description
			window["schedule"] = flattenMaintenanceSchedule(v.Schedule)
			window["owner_id"] = v.OwnerID
			window["date_created"] = v.DateCreated
			window["date_modified"] = v.DateModified
			window["access_level"] = v.AccessLevel
			window["tags"] = flattenDBTags(v.Tags)
			window["status"] = v.Status
			window["next_run_time"] = v.NextRunTime
			window["entity_task_assoc"] = flattenEntityTaskAssoc(v.EntityTaskAssoc)
			window["timezone"] = v.Timezone
			if v.Properties != nil {
				props := []interface{}{}
				for _, prop := range v.Properties {
					props = append(props, map[string]interface{}{
						"name":  prop.Name,
						"value": prop.Value,
					})
				}
				window["properties"] = props
			}

			windowResp = append(windowResp, window)
		}

		return windowResp
	}
	return nil
}
