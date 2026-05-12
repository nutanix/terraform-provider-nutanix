package volumesv2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesResponse "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/response"
	volumesStats "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/stats"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVolumeGroupStatsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the Volume Group stats identified by {extId}.",
		ReadContext: DatasourceNutanixVolumeGroupStatsV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "The external identifier of a Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"start_time": {
				Description: "The start time for the stats query in RFC3339 format (e.g. 2024-01-01T00:00:00Z).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"end_time": {
				Description: "The end time for the stats query in RFC3339 format (e.g. 2024-01-02T00:00:00Z). If not provided, defaults to current time.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"sampling_interval": {
				Description: "The sampling interval in seconds.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"stat_type": {
				Description: "The down sampling operator for the stats query.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"select": {
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tenant_id": {
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"links": {
				Description: "A HATEOAS style link for the response.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Description: "The URL at which the entity described by the link can be accessed.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"rel": {
							Description: "A name that identifies the relationship of the link to the object that is returned by the URL.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"volume_group_ext_id": {
				Description: "Uuid of the Volume Group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"controller_avg_io_latency_usecs": {
				Description: "Controller average I/O latency measured in microseconds.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_avg_read_io_latency_usecs": {
				Description: "Controller average read I/O latency measured in microseconds.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_avg_write_io_latency_usecs": {
				Description: "Controller average write I/O latency measured in microseconds.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_io_bandwidth_kbps": {
				Description: "Controller I/O bandwidth measured in Kbps.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_num_iops": {
				Description: "Controller I/O rate measured in iops.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_num_read_iops": {
				Description: "Controller read I/O measured in iops.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_num_write_iops": {
				Description: "Controller write I/O measured in iops.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_read_io_bandwidth_kbps": {
				Description: "Controller read I/O bandwidth measured in Kbps.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_user_bytes": {
				Description: "Controller user bytes.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"controller_write_io_bandwidth_kbps": {
				Description: "Controller write I/O bandwidth measured in Kbps.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
			"hydration_remaining_bytes": {
				Description: "Number of bytes that are left to hydrate the Volume Group.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: timeValuePairSchema(),
				},
			},
		},
	}
}

func timeValuePairSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"timestamp": {
			Description: "Timestamp is returned in Epoch format.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"value": {
			Description: "Value of the stat at the corresponding timestamp value represented in Int64 format.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
	}
}

func DatasourceNutanixVolumeGroupStatsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	extID := d.Get("ext_id").(string)

	startTimeStr := d.Get("start_time").(string)
	startTime, parseErr := time.Parse(time.RFC3339, startTimeStr)
	if parseErr != nil {
		return diag.Errorf("error parsing start_time: %v", parseErr)
	}

	var endTimePtr *time.Time
	if v, ok := d.GetOk("end_time"); ok {
		endTime, endParseErr := time.Parse(time.RFC3339, v.(string))
		if endParseErr != nil {
			return diag.Errorf("error parsing end_time: %v", endParseErr)
		}
		endTimePtr = &endTime
	}

	var samplingIntervalPtr *int
	if v, ok := d.GetOk("sampling_interval"); ok {
		val := v.(int)
		samplingIntervalPtr = &val
	}

	var selectPtr *string
	if v, ok := d.GetOk("select"); ok {
		val := v.(string)
		selectPtr = &val
	}

	resp, err := conn.VolumeAPIInstance.GetVolumeGroupStats(utils.StringPtr(extID), &startTime, endTimePtr, samplingIntervalPtr, nil, selectPtr)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group Stats : %v", err)
	}

	getResp := resp.Data.GetValue().(volumesStats.VolumeGroupStats)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenStatsLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("volume_group_ext_id", getResp.VolumeGroupExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_io_latency_usecs", flattenTimeValuePairs(getResp.ControllerAvgIOLatencyUsecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_read_io_latency_usecs", flattenTimeValuePairs(getResp.ControllerAvgReadIOLatencyUsecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_write_io_latency_usecs", flattenTimeValuePairs(getResp.ControllerAvgWriteIOLatencyUsecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_io_bandwidth_kbps", flattenTimeValuePairs(getResp.ControllerIOBandwidthKBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_iops", flattenTimeValuePairs(getResp.ControllerNumIOPS)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_read_iops", flattenTimeValuePairs(getResp.ControllerNumReadIOPS)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_write_iops", flattenTimeValuePairs(getResp.ControllerNumWriteIOPS)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_read_io_bandwidth_kbps", flattenTimeValuePairs(getResp.ControllerReadIOBandwidthKBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_user_bytes", flattenTimeValuePairs(getResp.ControllerUserBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_write_io_bandwidth_kbps", flattenTimeValuePairs(getResp.ControllerWriteIOBandwidthKBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hydration_remaining_bytes", flattenTimeValuePairs(getResp.HydrationRemainingBytes)); err != nil {
		return diag.FromErr(err)
	}

	if getResp.ExtId != nil {
		d.SetId(utils.StringValue(getResp.ExtId))
	} else {
		d.SetId(extID)
	}
	return nil
}

func flattenTimeValuePairs(pairs []volumesStats.TimeValuePair) []map[string]interface{} {
	if len(pairs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(pairs))
	for i, p := range pairs {
		entry := map[string]interface{}{}
		if p.Timestamp != nil {
			entry["timestamp"] = p.Timestamp.String()
		}
		if p.Value != nil {
			entry["value"] = int(*p.Value)
		}
		result[i] = entry
	}
	return result
}

func flattenStatsLinks(apiLinks []volumesResponse.ApiLink) []map[string]interface{} {
	if len(apiLinks) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(apiLinks))
	for i, v := range apiLinks {
		link := map[string]interface{}{}
		if v.Href != nil {
			link["href"] = v.Href
		}
		if v.Rel != nil {
			link["rel"] = v.Rel
		}
		result[i] = link
	}
	return result
}
