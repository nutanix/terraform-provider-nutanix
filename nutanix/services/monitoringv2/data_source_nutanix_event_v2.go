package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/request/events"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEventV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixEventV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the generated event.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
			"links": schemaForLinks(),
			"affected_entities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all the entities that are affected by the event or audit.",
				Elem: &schema.Resource{
					Schema: schemaForEntityReference(),
				},
			},
			"classifications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Various categories into which this event type can be classified. For example, hardware, storage, or license.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the cluster associated with the entity.",
			},
			"cluster_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster UUID associated with the cluster where the event was first raised.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time in ISO 8601 format when the event was created.",
			},
			"event_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A preconfigured or dynamically generated unique value for each event type.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional message associated with the event.",
			},
			"metric_details": schemaForMetricDetails(),
			"operation_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operation type associated with the audit.",
			},
			"parameters": schemaForParameters(),
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service which raised the event or audit. For internal Nutanix services, this value is set to \"Nutanix\".",
			},
			"source_cluster_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster UUID associated with the source entity of the event.",
			},
			"source_entity": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The source entity associated with the event.",
				Elem: &schema.Resource{
					Schema: schemaForEntityReference(),
				},
			},
		},
	}
}

func datasourceNutanixEventV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Get("ext_id").(string)

	req := &events.GetEventByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.EventsAPI.GetEventById(ctx, req)
	if err != nil {
		return diag.Errorf("error while fetching event: %v", err)
	}

	getResp := resp.Data.GetValue().(serviceability.Event)

	if err := d.Set("tenant_id", utils.StringValue(getResp.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("affected_entities", flattenEntityReferences(getResp.AffectedEntities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("classifications", getResp.Classifications); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_name", utils.StringValue(getResp.ClusterName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_uuid", utils.StringValue(getResp.ClusterUUID)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreationTime != nil {
		if err := d.Set("creation_time", getResp.CreationTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("event_type", utils.StringValue(getResp.EventType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("message", utils.StringValue(getResp.Message)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metric_details", flattenMetricDetails(getResp.MetricDetails)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.OperationType != nil {
		if err := d.Set("operation_type", getResp.OperationType.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("parameters", flattenParameters(getResp.Parameters)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("service_name", utils.StringValue(getResp.ServiceName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_cluster_uuid", utils.StringValue(getResp.SourceClusterUUID)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_entity", flattenEventEntityReference(getResp.SourceEntity)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
