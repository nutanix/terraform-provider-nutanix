package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEventsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixEventsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"events": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixEventV2(),
			},
		},
	}
}

func datasourceNutanixEventsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	var filter, orderBy, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	}

	resp, err := conn.EventsAPI.ListEvents(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching events: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("events", []map[string]interface{}{}); err != nil {
			return diag.Errorf("error setting events: %s", err)
		}
		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of events.",
		}}
	}

	getResp := resp.Data.GetValue().([]serviceability.Event)

	if err := d.Set("events", flattenEvents(getResp)); err != nil {
		return diag.Errorf("error setting events: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenEvents(events []serviceability.Event) []map[string]interface{} {
	if len(events) == 0 {
		return []map[string]interface{}{}
	}

	eventsList := make([]map[string]interface{}, 0, len(events))

	for _, event := range events {
		eventMap := make(map[string]interface{})

		eventMap["ext_id"] = utils.StringValue(event.ExtId)
		eventMap["tenant_id"] = utils.StringValue(event.TenantId)
		eventMap["links"] = flattenLinks(event.Links)
		eventMap["event_type"] = utils.StringValue(event.EventType)
		eventMap["message"] = utils.StringValue(event.Message)
		eventMap["creation_time"] = flattenTime(event.CreationTime)
		eventMap["cluster_name"] = utils.StringValue(event.ClusterName)
		eventMap["cluster_uuid"] = utils.StringValue(event.ClusterUUID)
		eventMap["service_name"] = utils.StringValue(event.ServiceName)
		eventMap["source_cluster_uuid"] = utils.StringValue(event.SourceClusterUUID)
		eventMap["operation_type"] = flattenOperationType(event.OperationType)
		eventMap["classifications"] = event.Classifications
		eventMap["source_entity"] = flattenEventEntityReference(event.SourceEntity)
		eventMap["affected_entities"] = flattenEntityReferences(event.AffectedEntities)
		eventMap["metric_details"] = flattenMetricDetails(event.MetricDetails)
		eventMap["parameters"] = flattenParameters(event.Parameters)

		eventsList = append(eventsList, eventMap)
	}

	return eventsList
}
