package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func DataSourceNutanixEraSLA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraSLARead,
		Schema: map[string]*schema.Schema{
			"sla_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"sla_name"},
			},
			"sla_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"sla_id"},
			},
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
			"quartely_retention": {
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
	}
}

func dataSourceNutanixEraSLARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	slaID, iok := d.GetOk("sla_id")
	slaName, nok := d.GetOk("sla_name")

	if !iok && !nok {
		return diag.Errorf("please provide one of sla_id or sla_name attributes")
	}

	resp, err := conn.Service.GetSLA(ctx, slaID.(string), slaName.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("unique_name", resp.Uniquename); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner_id", resp.Ownerid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("system_sla", resp.Systemsla); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("date_created", resp.Datecreated); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("date_modified", resp.Datemodified); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("continuous_retention", resp.Continuousretention); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("daily_retention", resp.Dailyretention); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("weekly_retention", resp.Weeklyretention); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("monthly_retention", resp.Monthlyretention); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("quartely_retention", resp.Quarterlyretention); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("yearly_retention", resp.Yearlyretention); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("reference_count", resp.Referencecount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pitr_enabled", resp.PitrEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("current_active_frequency", resp.CurrentActiveFrequency); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	return nil
}
