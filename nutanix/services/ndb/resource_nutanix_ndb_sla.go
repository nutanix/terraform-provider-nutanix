package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBSla() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBSlaCreate,
		ReadContext:   resourceNutanixNDBSlaRead,
		UpdateContext: resourceNutanixNDBSlaUpdate,
		DeleteContext: resourceNutanixNDBSlaDelete,
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
			"continuous_retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "30",
			},
			"daily_retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "7",
			},
			"weekly_retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "2",
			},
			"monthly_retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "2",
			},
			"quarterly_retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "1",
			},
			"yearly_retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// computed
			"unique_name": {
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
	}
}

func resourceNutanixNDBSlaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.SLAIntentInput{}

	if name, ok1 := d.GetOk("name"); ok1 {
		req.Name = utils.StringPtr(name.(string))
	}

	if desc, ok1 := d.GetOk("description"); ok1 {
		req.Description = utils.StringPtr(desc.(string))
	}

	if conRen, ok1 := d.GetOk("continuous_retention"); ok1 {
		req.ContinuousRetention = utils.IntPtr(conRen.(int))
	}

	if dailyRen, ok1 := d.GetOk("daily_retention"); ok1 {
		req.DailyRetention = utils.IntPtr(dailyRen.(int))
	}
	if weeklyRen, ok1 := d.GetOk("weekly_retention"); ok1 {
		req.WeeklyRetention = utils.IntPtr(weeklyRen.(int))
	}

	if monthRen, ok1 := d.GetOk("monthly_retention"); ok1 {
		req.MonthlyRetention = utils.IntPtr(monthRen.(int))
	}
	if quartRen, ok1 := d.GetOk("quarterly_retention"); ok1 {
		req.QuarterlyRetention = utils.IntPtr(quartRen.(int))
	}

	resp, err := conn.Service.CreateSLA(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	log.Printf("NDB SLA with %s id is created successfully", d.Id())
	return resourceNutanixNDBSlaRead(ctx, d, meta)
}

func resourceNutanixNDBSlaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// get the sla

	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("id is required for read operation")
	}
	resp, err := conn.Service.GetSLA(ctx, d.Id(), "")
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", resp.Name); err != nil {
		return diag.Errorf("error setting name for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("description", resp.Description); err != nil {
		return diag.Errorf("error setting description for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("continuous_retention", resp.Continuousretention); err != nil {
		return diag.Errorf("error setting continuous_retention for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("daily_retention", resp.Dailyretention); err != nil {
		return diag.Errorf("error setting daily_retention for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("weekly_retention", resp.Weeklyretention); err != nil {
		return diag.Errorf("error setting weekly_retention for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("monthly_retention", resp.Monthlyretention); err != nil {
		return diag.Errorf("error setting monthly_retention for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("quarterly_retention", resp.Quarterlyretention); err != nil {
		return diag.Errorf("error setting quarterly_retention for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("unique_name", resp.Uniquename); err != nil {
		return diag.Errorf("error setting unique_name for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("owner_id", resp.Ownerid); err != nil {
		return diag.Errorf("error setting owner_id for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("system_sla", resp.Systemsla); err != nil {
		return diag.Errorf("error setting system_sla for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("date_created", resp.Datecreated); err != nil {
		return diag.Errorf("error setting date_created for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("date_modified", resp.Datemodified); err != nil {
		return diag.Errorf("error setting date_modified for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("yearly_retention", resp.Yearlyretention); err != nil {
		return diag.Errorf("error setting yearly_retention for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("reference_count", resp.Referencecount); err != nil {
		return diag.Errorf("error setting reference_count for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("pitr_enabled", resp.PitrEnabled); err != nil {
		return diag.Errorf("error setting pitr_enabled for sla %s: %s", d.Id(), err)
	}

	if err = d.Set("current_active_frequency", resp.CurrentActiveFrequency); err != nil {
		return diag.Errorf("error setting current_active_frequency for sla %s: %s", d.Id(), err)
	}
	return nil
}

func resourceNutanixNDBSlaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	req := &era.SLAIntentInput{}
	// get the current sla

	resp, er := conn.Service.GetSLA(ctx, d.Id(), "")
	if er != nil {
		return diag.FromErr(er)
	}

	if resp != nil {
		req.ContinuousRetention = &resp.Continuousretention
		req.DailyRetention = &resp.Dailyretention
		req.Description = resp.Description
		req.MonthlyRetention = &resp.Monthlyretention
		req.Name = resp.Name
		req.QuarterlyRetention = &resp.Quarterlyretention
		req.WeeklyRetention = &resp.Weeklyretention
	}

	if d.HasChange("name") {
		req.Name = utils.StringPtr(d.Get("name").(string))
	}

	if d.HasChange("description") {
		req.Description = utils.StringPtr(d.Get("description").(string))
	}

	if d.HasChange("continuous_retention") {
		req.ContinuousRetention = utils.IntPtr(d.Get("continuous_retention").(int))
	}

	if d.HasChange("daily_retention") {
		req.DailyRetention = utils.IntPtr(d.Get("daily_retention").(int))
	}
	if d.HasChange("weekly_retention") {
		req.WeeklyRetention = utils.IntPtr(d.Get("weekly_retention").(int))
	}

	if d.HasChange("monthly_retention") {
		req.MonthlyRetention = utils.IntPtr(d.Get("monthly_retention").(int))
	}
	if d.HasChange("quarterly_retention") {
		req.QuarterlyRetention = utils.IntPtr(d.Get("quarterly_retention").(int))
	}

	// Adding id in payload for update going to be implemented in future
	req.ID = utils.StringPtr(d.Id())

	_, err := conn.Service.UpdateSLA(ctx, req, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("NDB SLA with %s id is updated successfully", d.Id())
	return resourceNutanixNDBSlaRead(ctx, d, meta)
}

func resourceNutanixNDBSlaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.DeleteSLA(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Status == utils.StringPtr("success") {
		log.Printf("NDB SLA with %s id is deleted successfully", d.Id())
		d.SetId("")
	}
	return nil
}
