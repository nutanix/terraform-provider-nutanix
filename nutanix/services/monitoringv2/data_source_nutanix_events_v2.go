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

func DatasourceNutanixEventsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixEventsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A URL query parameter that specifies the page number of the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100.",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to filter a collection of resources.",
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.",
			},
			"select": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.",
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

	req := &events.ListEventsRequest{}

	if pagef, ok := d.GetOk("page"); ok {
		req.Page_ = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		req.Limit_ = utils.IntPtr(limitf.(int))
	}
	if filterf, ok := d.GetOk("filter"); ok {
		req.Filter_ = utils.StringPtr(filterf.(string))
	}
	if order, ok := d.GetOk("order_by"); ok {
		req.Orderby_ = utils.StringPtr(order.(string))
	}
	if selectf, ok := d.GetOk("select"); ok {
		req.Select_ = utils.StringPtr(selectf.(string))
	}

	resp, err := conn.EventsAPI.ListEvents(ctx, req)
	if err != nil {
		return diag.Errorf("error while fetching events: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("events", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenEvents(eventsList []serviceability.Event) []map[string]interface{} {
	if len(eventsList) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, len(eventsList))

	for i, event := range eventsList {
		eventMap := make(map[string]interface{})

		eventMap["ext_id"] = utils.StringValue(event.ExtId)
		eventMap["tenant_id"] = utils.StringValue(event.TenantId)
		eventMap["links"] = flattenLinks(event.Links)
		eventMap["affected_entities"] = flattenEntityReferences(event.AffectedEntities)
		eventMap["classifications"] = event.Classifications
		eventMap["cluster_name"] = utils.StringValue(event.ClusterName)
		eventMap["cluster_uuid"] = utils.StringValue(event.ClusterUUID)

		if event.CreationTime != nil {
			eventMap["creation_time"] = event.CreationTime.String()
		}

		eventMap["event_type"] = utils.StringValue(event.EventType)
		eventMap["message"] = utils.StringValue(event.Message)
		eventMap["metric_details"] = flattenMetricDetails(event.MetricDetails)

		if event.OperationType != nil {
			eventMap["operation_type"] = event.OperationType.GetName()
		}

		eventMap["parameters"] = flattenParameters(event.Parameters)
		eventMap["service_name"] = utils.StringValue(event.ServiceName)
		eventMap["source_cluster_uuid"] = utils.StringValue(event.SourceClusterUUID)
		eventMap["source_entity"] = flattenEventEntityReference(event.SourceEntity)

		result[i] = eventMap
	}

	return result
}
