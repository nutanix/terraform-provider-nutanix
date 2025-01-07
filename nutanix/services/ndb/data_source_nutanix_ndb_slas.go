package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	Era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixEraSLAs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraSLAsRead,
		Schema: map[string]*schema.Schema{
			"slas": {
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
				},
			},
		},
	}
}

func dataSourceNutanixEraSLAsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.ListSLA(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("slas", flattenSLAsResponse(resp)); err != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era clusters: %+v", err)
	}
	d.SetId(uuid)
	return nil
}

func flattenSLAsResponse(sla *Era.SLAResponse) []map[string]interface{} {
	if sla != nil {
		lst := []map[string]interface{}{}
		for _, data := range *sla {
			d := map[string]interface{}{}
			d["id"] = data.ID
			d["name"] = data.Name
			d["unique_name"] = data.Uniquename
			d["description"] = data.Description
			d["owner_id"] = data.Ownerid
			d["system_sla"] = data.Systemsla
			d["date_created"] = data.Datecreated
			d["date_modified"] = data.Datemodified
			d["continuous_retention"] = data.Continuousretention
			d["daily_retention"] = data.Dailyretention
			d["weekly_retention"] = data.Weeklyretention
			d["monthly_retention"] = data.Monthlyretention
			d["quartely_retention"] = data.Quarterlyretention
			d["yearly_retention"] = data.Yearlyretention
			d["reference_count"] = data.Referencecount
			d["pitr_enabled"] = data.PitrEnabled
			d["current_active_frequency"] = data.CurrentActiveFrequency
			lst = append(lst, d)
		}
		return lst
	}
	return nil
}
