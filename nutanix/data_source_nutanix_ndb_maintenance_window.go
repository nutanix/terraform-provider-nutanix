package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNutanixNDBMaintenanceWindow() *schema.Resource {
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
			"entity_task_assoc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixNDBMaintenanceWindowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

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

	if err := d.Set("entity_task_assoc", resp.EntityTaskAssoc); err != nil {
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
