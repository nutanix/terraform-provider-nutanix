package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixNDBTimeMachines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBTimeMachinesRead,
		Schema: map[string]*schema.Schema{
			"time_machines": dataSourceEraTimeMachine(),
		},
	}
}

func dataSourceNutanixNDBTimeMachinesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// call tms API
	resp, err := conn.Service.ListTimeMachines(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("time_machines", flattenTimeMachines(resp)); e != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era time machines: %+v", err)
	}
	d.SetId(uuid)
	return nil
}

func flattenTimeMachines(tms *era.ListTimeMachines) []map[string]interface{} {
	if tms != nil {
		lst := []map[string]interface{}{}

		for _, pr := range *tms {
			tmac := map[string]interface{}{}

			tmac["id"] = pr.ID
			tmac["name"] = pr.Name
			tmac["description"] = pr.Description
			tmac["date_created"] = pr.DateCreated
			tmac["date_modified"] = pr.DateModified
			tmac["access_level"] = pr.AccessLevel
			tmac["properties"] = flattenDBInstanceProperties(pr.Properties)
			tmac["tags"] = flattenDBTags(pr.Tags)
			tmac["clustered"] = pr.Clustered
			tmac["clone"] = pr.Clone
			tmac["database_id"] = pr.DatabaseID
			tmac["type"] = pr.Type
			tmac["status"] = pr.Status
			tmac["ea_status"] = pr.EaStatus
			tmac["scope"] = pr.Scope
			tmac["sla_id"] = pr.SLAID
			tmac["schedule_id"] = pr.ScheduleID
			tmac["metric"] = pr.Metric
			tmac["database"] = pr.Database
			tmac["clones"] = pr.Clones
			tmac["source_nx_clusters"] = pr.SourceNxClusters
			tmac["sla_update_in_progress"] = pr.SLAUpdateInProgress
			tmac["sla"] = flattenDBSLA(pr.SLA)
			tmac["schedule"] = flattenSchedule(pr.Schedule)

			lst = append(lst, tmac)
		}
		return lst
	}
	return nil
}
