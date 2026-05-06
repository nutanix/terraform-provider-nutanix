package volumesv2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "The start time for the stats query in RFC3339 format (e.g. 2023-01-01T00:00:00Z).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"end_time": {
				Description: "The end time for the stats query in RFC3339 format (e.g. 2023-01-02T00:00:00Z).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"sampling_interval": {
				Description: "The sampling interval in seconds.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"tenant_id": {
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"volume_group_ext_id": {
				Description: "Uuid of the Volume Group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"controller_avg_io_latency_usecs":       schemaForTimeValuePair("Controller average I/O latency measured in microseconds."),
			"controller_avg_read_io_latency_usecs":  schemaForTimeValuePair("Controller average read I/O latency measured in microseconds."),
			"controller_avg_write_io_latency_usecs": schemaForTimeValuePair("Controller average write I/O latency measured in microseconds."),
			"controller_io_bandwidth_k_bps":         schemaForTimeValuePair("Controller I/O bandwidth measured in Kbps."),
			"controller_num_iops":                   schemaForTimeValuePair("Controller I/O rate measured in iops."),
			"controller_num_read_iops":              schemaForTimeValuePair("Controller read I/O measured in iops."),
			"controller_num_write_iops":             schemaForTimeValuePair("Controller write I/O measured in iops."),
			"controller_read_io_bandwidth_k_bps":    schemaForTimeValuePair("Controller read I/O bandwidth measured in Kbps."),
			"controller_user_bytes":                 schemaForTimeValuePair("Controller user bytes."),
			"controller_write_io_bandwidth_k_bps":   schemaForTimeValuePair("Controller write I/O bandwidth measured in Kbps."),
			"hydration_remaining_bytes":             schemaForTimeValuePair("Number of bytes that are left to hydrate the Volume Group."),
		},
	}
}

func schemaForTimeValuePair(description string) *schema.Schema {
	return &schema.Schema{
		Description: description,
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
			},
		},
	}
}

func DatasourceNutanixVolumeGroupStatsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	extID := d.Get("ext_id").(string)

	startTimeStr := d.Get("start_time").(string)
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return diag.Errorf("error parsing start_time: %v", err)
	}

	var endTimePtr *time.Time
	if v, ok := d.GetOk("end_time"); ok {
		endTime, parseErr := time.Parse(time.RFC3339, v.(string))
		if parseErr != nil {
			return diag.Errorf("error parsing end_time: %v", parseErr)
		}
		endTimePtr = &endTime
	}

	var samplingIntervalPtr *int
	if v, ok := d.GetOk("sampling_interval"); ok {
		si := v.(int)
		samplingIntervalPtr = &si
	}

	resp, err := conn.VolumeAPIInstance.GetVolumeGroupStats(utils.StringPtr(extID), &startTime, endTimePtr, samplingIntervalPtr, nil, nil)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group Stats: %v", err)
	}

	if resp.Data == nil || resp.Data.GetValue() == nil {
		d.SetId(extID)
		return nil
	}

	getResp := resp.Data.GetValue().(volumesStats.VolumeGroupStats)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
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
	if err := d.Set("controller_io_bandwidth_k_bps", flattenTimeValuePairs(getResp.ControllerIOBandwidthKBps)); err != nil {
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
	if err := d.Set("controller_read_io_bandwidth_k_bps", flattenTimeValuePairs(getResp.ControllerReadIOBandwidthKBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_user_bytes", flattenTimeValuePairs(getResp.ControllerUserBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_write_io_bandwidth_k_bps", flattenTimeValuePairs(getResp.ControllerWriteIOBandwidthKBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hydration_remaining_bytes", flattenTimeValuePairs(getResp.HydrationRemainingBytes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extID)
	return nil
}

func flattenTimeValuePairs(pairs []volumesStats.TimeValuePair) []map[string]interface{} {
	if len(pairs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(pairs))
	for i, pair := range pairs {
		entry := map[string]interface{}{}
		if pair.Timestamp != nil {
			entry["timestamp"] = pair.Timestamp.String()
		}
		if pair.Value != nil {
			entry["value"] = int(*pair.Value)
		}
		result[i] = entry
	}
	return result
}
