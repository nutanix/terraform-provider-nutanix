package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBMaintenanceWindowCreate,
		ReadContext:   resourceNutanixNDBMaintenanceWindowRead,
		UpdateContext: resourceNutanixNDBMaintenanceWindowUpdate,
		DeleteContext: resourceNutanixNDBMaintenanceWindowDelete,
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
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Asia/Calcutta",
			},

			"recurrence": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MONTHLY", "WEEKLY"}, false),
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "2",
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"day_of_week": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY",
					"FRIDAY", "SATURDAY", "SUNDAY",
				}, false),
			},
			"week_of_month": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4}),
			},

			// compute

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
		},
	}
}

func resourceNutanixNDBMaintenanceWindowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.MaintenanceWindowInput{}
	schedule := &era.MaintenaceSchedule{}

	if name, ok := d.GetOk("name"); ok {
		req.Name = utils.StringPtr(name.(string))
	}

	if desc, ok := d.GetOk("description"); ok {
		req.Description = utils.StringPtr(desc.(string))
	}

	if timezone, ok := d.GetOk("timezone"); ok {
		req.Timezone = utils.StringPtr(timezone.(string))
	}

	if recurrence, ok := d.GetOk("recurrence"); ok {
		schedule.Recurrence = utils.StringPtr(recurrence.(string))
	}

	if duration, ok := d.GetOk("duration"); ok {
		schedule.Duration = utils.IntPtr(duration.(int))
	}

	if startTime, ok := d.GetOk("start_time"); ok {
		schedule.StartTime = utils.StringPtr(startTime.(string))
	}

	if dayOfWeek, ok := d.GetOk("day_of_week"); ok && len(dayOfWeek.(string)) > 0 {
		schedule.DayOfWeek = utils.StringPtr(dayOfWeek.(string))
	}

	if weekOfMonth, ok := d.GetOk("week_of_month"); ok {
		schedule.WeekOfMonth = utils.IntPtr(weekOfMonth.(int))
	}

	req.Schedule = schedule

	resp, err := conn.Service.CreateMaintenanceWindow(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	log.Printf("NDB Maintenance Window with %s id is created successfully", d.Id())
	return resourceNutanixNDBMaintenanceWindowRead(ctx, d, meta)
}

func resourceNutanixNDBMaintenanceWindowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("id is required for read operation")
	}
	resp, err := conn.Service.ReadMaintenanceWindow(ctx, d.Id())
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
	return nil
}

func resourceNutanixNDBMaintenanceWindowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.MaintenanceWindowInput{}
	sch := &era.MaintenaceSchedule{}

	resp, err := conn.Service.ReadMaintenanceWindow(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {
		req.Name = resp.Name
		req.Description = resp.Description
		req.Timezone = resp.Timezone

		// read schedule info

		if resp.Schedule != nil {
			sch.DayOfWeek = resp.Schedule.DayOfWeek
			sch.Duration = resp.Schedule.Duration
			sch.StartTime = resp.Schedule.StartTime
			sch.Recurrence = resp.Schedule.Recurrence
			sch.WeekOfMonth = resp.Schedule.WeekOfMonth
		}
	}
	if d.HasChange("name") {
		req.Name = utils.StringPtr(d.Get("name").(string))
		req.ResetName = utils.BoolPtr(true)
	}

	if d.HasChange("description") {
		req.Description = utils.StringPtr(d.Get("description").(string))
		req.ResetDescription = utils.BoolPtr(true)
	}

	if d.HasChange("timezone") {
		req.Timezone = utils.StringPtr(d.Get("timezone").(string))
	}

	if d.HasChange("recurrence") {
		sch.Recurrence = utils.StringPtr(d.Get("recurrence").(string))
	}

	if d.HasChange("duration") {
		sch.Duration = utils.IntPtr(d.Get("duration").(int))
	}

	if d.HasChange("start_time") {
		sch.StartTime = utils.StringPtr(d.Get("start_time").(string))
	}

	if d.HasChange("day_of_week") {
		sch.DayOfWeek = utils.StringPtr(d.Get("day_of_week").(string))
	}

	if d.HasChange("week_of_month") {
		sch.WeekOfMonth = utils.IntPtr(d.Get("week_of_month").(int))
	}

	req.Schedule = sch
	req.ResetSchedule = utils.BoolPtr(true)

	respUpdate, err := conn.Service.UpdateMaintenaceWindow(ctx, req, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("NDB Maintenance Window with %s id is updated successfully", *respUpdate.ID)
	return resourceNutanixNDBMaintenanceWindowRead(ctx, d, meta)
}

func resourceNutanixNDBMaintenanceWindowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.DeleteMaintenanceWindow(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Status == utils.StringPtr("success") {
		log.Printf("NDB Maintenance Window with %s id is deleted successfully", d.Id())
		d.SetId("")
	}
	return nil
}

func flattenMaintenanceSchedule(pr *era.MaintenaceSchedule) []map[string]interface{} {
	if pr != nil {
		res := make([]map[string]interface{}, 0)

		schedule := map[string]interface{}{}

		schedule["recurrence"] = pr.Recurrence
		schedule["duration"] = pr.Duration
		schedule["start_time"] = pr.StartTime
		schedule["day_of_week"] = pr.DayOfWeek
		schedule["week_of_month"] = pr.WeekOfMonth
		schedule["threshold"] = pr.Threshold
		schedule["hour"] = pr.Hour
		schedule["minute"] = pr.Minute
		schedule["timezone"] = pr.TimeZone

		res = append(res, schedule)
		return res
	}
	return nil
}

func expandMaintenanceSchdeule(pr []interface{}) *era.MaintenaceSchedule {
	if len(pr) > 0 {
		sch := &era.MaintenaceSchedule{}

		for _, v := range pr {
			val := v.(map[string]interface{})

			if recurrence, ok := val["recurrence"]; ok {
				sch.Recurrence = utils.StringPtr(recurrence.(string))
			}

			if duration, ok := val["duration"]; ok {
				sch.Duration = utils.IntPtr(duration.(int))
			}

			if startTime, ok := val["start_time"]; ok {
				sch.StartTime = utils.StringPtr(startTime.(string))
			}

			if dayOfWeek, ok := val["day_of_week"]; ok {
				sch.DayOfWeek = utils.StringPtr(dayOfWeek.(string))
			}

			if weekOfMonth, ok := val["week_of_month"]; ok && len(weekOfMonth.(string)) > 0 {
				sch.WeekOfMonth = utils.IntPtr(weekOfMonth.(int))
			}
		}
		return sch
	}
	return nil
}
