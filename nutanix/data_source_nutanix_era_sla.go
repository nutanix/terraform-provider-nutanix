package nutanix

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	Era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
)

func dataSourceNutanixEraSLA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraSLARead,
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

func dataSourceNutanixEraSLARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	resp, err := conn.Service.ListSLA(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("HELLLLLOOOOOO")
	aJSON, _ := json.Marshal(resp)
	fmt.Printf("JSON Print - \n%s\n", string(aJSON))

	if err := d.Set("slas", flattenSLAsResponse(resp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

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
