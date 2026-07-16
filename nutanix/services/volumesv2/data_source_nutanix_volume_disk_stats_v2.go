package volumesv2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	volumesCommonStats "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/stats"
	volumesStats "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/stats"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixVolumeDiskStatsV2 queries the Volume Disk stats identified by {extId}
// in the Volume Group identified by {volumeGroupExtId}.
func DatasourceNutanixVolumeDiskStatsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the Volume Disk stats identified by {diskExtId}.",
		ReadContext: DatasourceNutanixVolumeDiskStatsV2Read,
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of a Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ext_id": {
				Description: "The external identifier of a Volume Disk.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"start_time": {
				Description: "The start time in RFC-3339 format from which the stats should be reported.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"end_time": {
				Description: "The end time in RFC-3339 format until which the stats should be reported.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"sampling_interval": {
				Description: "The sampling interval in seconds at which the stats should be reported.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"stat_type": {
				Description:  "The operator to use while performing down-sampling on stats data. Allowed values are SUM, MIN, MAX, AVG, COUNT and LAST.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"SUM", "MIN", "MAX", "AVG", "COUNT", "LAST"}, false),
			},
			"select": {
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tenant_id": {
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"links": {
				Description: "A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.",
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
							Description: "A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of \"self\" identifies the URL for the object.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"volume_disk_ext_id": {
				Description: "Uuid of the Volume Disk.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"controller_avg_io_latency_usecs":       schemaForVolumeDiskStatsTimeValuePair("Controller average I/O latency measured in microseconds."),
			"controller_avg_read_io_latency_usecs":  schemaForVolumeDiskStatsTimeValuePair("Controller average read I/O latency measured in microseconds."),
			"controller_avg_write_io_latency_usecs": schemaForVolumeDiskStatsTimeValuePair("Controller average write I/O latency measured in microseconds."),
			"controller_io_bandwidth_kbps":          schemaForVolumeDiskStatsTimeValuePair("Controller I/O bandwidth measured in Kbps."),
			"controller_num_iops":                   schemaForVolumeDiskStatsTimeValuePair("Controller I/O rate measured in iops."),
			"controller_num_read_iops":              schemaForVolumeDiskStatsTimeValuePair("Controller read I/O measured in iops."),
			"controller_num_write_iops":             schemaForVolumeDiskStatsTimeValuePair("Controller write I/O measured in iops."),
			"controller_read_io_bandwidth_kbps":     schemaForVolumeDiskStatsTimeValuePair("Controller read I/O bandwidth measured in Kbps."),
			"controller_user_bytes":                 schemaForVolumeDiskStatsTimeValuePair("Controller user bytes."),
			"controller_write_io_bandwidth_kbps":    schemaForVolumeDiskStatsTimeValuePair("Controller write I/O bandwidth measured in Kbps."),
		},
	}
}

// schemaForVolumeDiskStatsTimeValuePair returns the shared schema for a list of
// timestamp/value stat pairs used across all Volume Disk stat attributes.
func schemaForVolumeDiskStatsTimeValuePair(description string) *schema.Schema {
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

func DatasourceNutanixVolumeDiskStatsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id").(string)
	volumeDiskExtID := d.Get("ext_id").(string)

	startTimeVal, err := time.Parse(time.RFC3339, d.Get("start_time").(string))
	if err != nil {
		return diag.Errorf("error while parsing start_time : %v", err)
	}
	endTimeVal, err := time.Parse(time.RFC3339, d.Get("end_time").(string))
	if err != nil {
		return diag.Errorf("error while parsing end_time : %v", err)
	}

	var samplingInterval *int
	if si, ok := d.GetOk("sampling_interval"); ok {
		if si.(int) <= 0 {
			return diag.Errorf("sampling_interval should be greater than 0")
		}
		samplingInterval = utils.IntPtr(si.(int))
	}

	// Default value is LAST, aggregation containing only the last recorded value.
	statType := volumesCommonStats.DOWNSAMPLINGOPERATOR_LAST
	statTypeMap := map[string]volumesCommonStats.DownSamplingOperator{
		"SUM":   volumesCommonStats.DOWNSAMPLINGOPERATOR_SUM,
		"MIN":   volumesCommonStats.DOWNSAMPLINGOPERATOR_MIN,
		"MAX":   volumesCommonStats.DOWNSAMPLINGOPERATOR_MAX,
		"AVG":   volumesCommonStats.DOWNSAMPLINGOPERATOR_AVG,
		"COUNT": volumesCommonStats.DOWNSAMPLINGOPERATOR_COUNT,
		"LAST":  volumesCommonStats.DOWNSAMPLINGOPERATOR_LAST,
	}
	if st, ok := d.GetOk("stat_type"); ok {
		statType = statTypeMap[st.(string)]
	}

	var selectQuery *string
	if sel, ok := d.GetOk("select"); ok {
		selectQuery = utils.StringPtr(sel.(string))
	}

	resp, err := conn.VolumeAPIInstance.GetVolumeDiskStats(
		utils.StringPtr(volumeGroupExtID),
		utils.StringPtr(volumeDiskExtID),
		&startTimeVal,
		&endTimeVal,
		samplingInterval,
		&statType,
		selectQuery,
	)
	if err != nil {
		return diag.Errorf("error while fetching Volume Disk stats : %v", err)
	}

	getResp := resp.Data.GetValue().(volumesStats.VolumeDiskStats)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("volume_disk_ext_id", getResp.VolumeDiskExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_io_latency_usecs", flattenVolumeStatsTimeValuePair(getResp.ControllerAvgIOLatencyUsecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_read_io_latency_usecs", flattenVolumeStatsTimeValuePair(getResp.ControllerAvgReadIOLatencyUsecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_write_io_latency_usecs", flattenVolumeStatsTimeValuePair(getResp.ControllerAvgWriteIOLatencyUsecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_io_bandwidth_kbps", flattenVolumeStatsTimeValuePair(getResp.ControllerIOBandwidthKBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_iops", flattenVolumeStatsTimeValuePair(getResp.ControllerNumIOPS)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_read_iops", flattenVolumeStatsTimeValuePair(getResp.ControllerNumReadIOPS)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_write_iops", flattenVolumeStatsTimeValuePair(getResp.ControllerNumWriteIOPS)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_read_io_bandwidth_kbps", flattenVolumeStatsTimeValuePair(getResp.ControllerReadIOBandwidthKBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_user_bytes", flattenVolumeStatsTimeValuePair(getResp.ControllerUserBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_write_io_bandwidth_kbps", flattenVolumeStatsTimeValuePair(getResp.ControllerWriteIOBandwidthKBps)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

// flattenVolumeStatsTimeValuePair converts a slice of Volume stats TimeValuePair
// into the Terraform state representation.
func flattenVolumeStatsTimeValuePair(timeValuePairs []volumesStats.TimeValuePair) []map[string]interface{} {
	if len(timeValuePairs) == 0 {
		return nil
	}
	timeValueList := make([]map[string]interface{}, len(timeValuePairs))
	for k, v := range timeValuePairs {
		pair := map[string]interface{}{}
		if v.Timestamp != nil {
			pair["timestamp"] = v.Timestamp.Format(time.RFC3339)
		}
		if v.Value != nil {
			pair["value"] = int(utils.Int64Value(v.Value))
		}
		timeValueList[k] = pair
	}
	return timeValueList
}
