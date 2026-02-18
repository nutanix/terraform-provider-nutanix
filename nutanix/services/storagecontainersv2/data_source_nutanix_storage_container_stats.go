package storagecontainersv2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	clustermgmtStats "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/stats"
	clsstats "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/stats"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixStorageStatsInfoV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixStorageStatsInfoV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sampling_interval": {
				Type:     schema.TypeInt,
				Default:  1,
				Optional: true,
			},
			"stat_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"AVG", "MIN", "MAX", "LAST", "SUM", "COUNT"}, false),
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"controller_num_iops":                            SchemaForValueTimestamp(),
			"controller_io_bandwidth_kbps":                   SchemaForValueTimestamp(),
			"controller_avg_io_latencyu_secs":                SchemaForValueTimestamp(),
			"controller_num_read_iops":                       SchemaForValueTimestamp(),
			"controller_num_write_iops":                      SchemaForValueTimestamp(),
			"controller_read_io_bandwidth_kbps":              SchemaForValueTimestamp(),
			"controller_write_io_bandwidth_kbps":             SchemaForValueTimestamp(),
			"controller_avg_read_io_latencyu_secs":           SchemaForValueTimestamp(),
			"controller_avg_write_io_latencyu_secs":          SchemaForValueTimestamp(),
			"storage_reserved_capacity_bytes":                SchemaForValueTimestamp(),
			"storage_actual_physical_usage_bytes":            SchemaForValueTimestamp(),
			"data_reduction_saving_ratio_ppm":                SchemaForValueTimestamp(),
			"data_reduction_total_saving_ratio_ppm":          SchemaForValueTimestamp(),
			"storage_free_bytes":                             SchemaForValueTimestamp(),
			"storage_capacity_bytes":                         SchemaForValueTimestamp(),
			"data_reduction_saved_bytes":                     SchemaForValueTimestamp(),
			"data_reduction_overall_pre_reduction_bytes":     SchemaForValueTimestamp(),
			"data_reduction_overall_post_reduction_bytes":    SchemaForValueTimestamp(),
			"data_reduction_compression_saving_ratio_ppm":    SchemaForValueTimestamp(),
			"data_reduction_dedup_saving_ratio_ppm":          SchemaForValueTimestamp(),
			"data_reduction_erasure_coding_saving_ratio_ppm": SchemaForValueTimestamp(),
			"data_reduction_thin_provision_saving_ratio_ppm": SchemaForValueTimestamp(),
			"data_reduction_clone_saving_ratio_ppm":          SchemaForValueTimestamp(),
			"data_reduction_snapshot_saving_ratio_ppm":       SchemaForValueTimestamp(),
			"data_reduction_zero_write_savings_bytes":        SchemaForValueTimestamp(),
			"controller_read_io_ratio_ppm":                   SchemaForValueTimestamp(),
			"controller_write_io_ratio_ppm":                  SchemaForValueTimestamp(),
			"storage_replication_factor":                     SchemaForValueTimestamp(),
			"storage_usage_bytes":                            SchemaForValueTimestamp(),
			"storage_tier_das_sata_usage_bytes":              SchemaForValueTimestamp(),
			"storage_tier_ssd_usage_bytes":                   SchemaForValueTimestamp(),
			"health":                                         SchemaForValueTimestamp(),
			"container_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func SchemaForValueTimestamp() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"timestamp": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func DatasourceNutanixStorageStatsInfoV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	extID := d.Get("ext_id")
	startTime := d.Get("start_time")
	endTime := d.Get("end_time")
	samplingInterval := d.Get("sampling_interval")

	if samplingInterval.(int) <= 0 {
		return diag.Errorf("sampling_interval should be greater than 0")
	}

	const two, three, four, five, six, seven = 2, 3, 4, 5, 6, 7
	statType := clsstats.DownSamplingOperator(seven) // Default value is LAST, Aggregation containing only the last recorded value.

	subMap := map[string]interface{}{
		"SUM":   two,
		"MIN":   three,
		"MAX":   four,
		"AVG":   five,
		"COUNT": six,
		"LAST":  seven,
	}
	pVal := subMap[d.Get("stat_type").(string)]
	if pVal != nil {
		statType = clsstats.DownSamplingOperator(pVal.(int))
	}
	resp, err := conn.StorageContainersAPI.GetStorageContainerById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching Storage Container : %v", err)
	}

	// Extract E-Tag Header
	etagValue := conn.ClusterEntityAPI.ApiClient.GetEtag(resp)

	args := make(map[string]interface{})
	args["If-Match"] = etagValue

	startTimeVal, err := time.Parse(time.RFC3339, startTime.(string))
	if err != nil {
		return diag.Errorf("error while parsing start_time : %v", err)
	}
	endTimeVal, err := time.Parse(time.RFC3339, endTime.(string))
	if err != nil {
		return diag.Errorf("error while parsing end_time : %v", err)
	}

	statsResp, err := conn.StorageContainersAPI.GetStorageContainerStats(utils.StringPtr(extID.(string)), &startTimeVal, &endTimeVal, utils.IntPtr(samplingInterval.(int)), &statType, args)
	if err != nil {
		return diag.Errorf("error while fetching Storage Container : %v", err)
	}

	getStatsResp := statsResp.Data.GetValue().(clustermgmtStats.StorageContainerStats)

	if err := d.Set("ext_id", getStatsResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getStatsResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getStatsResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("container_ext_id", getStatsResp.ContainerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_iops", flattenValueTimestamp(getStatsResp.ControllerNumIops)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_io_bandwidth_kbps", flattenValueTimestamp(getStatsResp.ControllerIoBandwidthkBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_io_latencyu_secs", flattenValueTimestamp(getStatsResp.ControllerAvgIoLatencyuSecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_read_iops", flattenValueTimestamp(getStatsResp.ControllerNumReadIops)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_num_write_iops", flattenValueTimestamp(getStatsResp.ControllerNumWriteIops)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_read_io_bandwidth_kbps", flattenValueTimestamp(getStatsResp.ControllerReadIoBandwidthkBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_write_io_bandwidth_kbps", flattenValueTimestamp(getStatsResp.ControllerWriteIoBandwidthkBps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_read_io_latencyu_secs", flattenValueTimestamp(getStatsResp.ControllerAvgReadIoLatencyuSecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_avg_write_io_latencyu_secs", flattenValueTimestamp(getStatsResp.ControllerAvgWriteIoLatencyuSecs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_reserved_capacity_bytes", flattenValueTimestamp(getStatsResp.StorageReservedCapacityBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_actual_physical_usage_bytes", flattenValueTimestamp(getStatsResp.StorageActualPhysicalUsageBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_saving_ratio_ppm", flattenValueTimestamp(getStatsResp.DataReductionSavingRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_total_saving_ratio_ppm", flattenValueTimestamp(getStatsResp.DataReductionTotalSavingRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_free_bytes", flattenValueTimestamp(getStatsResp.StorageFreeBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_capacity_bytes", flattenValueTimestamp(getStatsResp.StorageCapacityBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_saved_bytes", flattenValueTimestamp(getStatsResp.DataReductionSavedBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_overall_pre_reduction_bytes", flattenValueTimestamp(getStatsResp.DataReductionOverallPreReductionBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_overall_post_reduction_bytes", flattenValueTimestamp(getStatsResp.DataReductionOverallPostReductionBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_compression_saving_ratio_ppm", flattenValueTimestamp(getStatsResp.DataReductionCompressionSavingRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_dedup_saving_ratio_ppm", flattenValueTimestamp(getStatsResp.DataReductionDedupSavingRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_erasure_coding_saving_ratio_ppm", flattenValueTimestamp(getStatsResp.DataReductionErasureCodingSavingRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_thin_provision_saving_ratio_ppm", flattenValueTimestamp(getStatsResp.DataReductionThinProvisionSavingRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_clone_saving_ratio_ppm", flattenValueTimestamp(getStatsResp.DataReductionCloneSavingRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_snapshot_saving_ratio_ppm", flattenValueTimestamp(getStatsResp.DataReductionSnapshotSavingRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_reduction_zero_write_savings_bytes", flattenValueTimestamp(getStatsResp.DataReductionZeroWriteSavingsBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_read_io_ratio_ppm", flattenValueTimestamp(getStatsResp.ControllerReadIoRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_write_io_ratio_ppm", flattenValueTimestamp(getStatsResp.ControllerWriteIoRatioPpm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_replication_factor", flattenValueTimestamp(getStatsResp.StorageReplicationFactor)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_usage_bytes", flattenValueTimestamp(getStatsResp.StorageUsageBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_tier_das_sata_usage_bytes", flattenValueTimestamp(getStatsResp.StorageTierDasSataUsageBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_tier_ssd_usage_bytes", flattenValueTimestamp(getStatsResp.StorageTierSsdUsageBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("health", flattenValueTimestamp(getStatsResp.Health)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getStatsResp.ContainerExtId))
	return nil
}

func flattenValueTimestamp(timeIntValuePairs []clsstats.TimeIntValuePair) []map[string]interface{} {
	if len(timeIntValuePairs) > 0 {
		timeIntValueList := make([]map[string]interface{}, len(timeIntValuePairs))

		for k, v := range timeIntValuePairs {
			timeValuePair := map[string]interface{}{}
			if v.Value != nil {
				timeValuePair["value"] = v.Value
			}
			if v.Timestamp != nil {
				timeValuePair["timestamp"] = v.Timestamp.Format("2006-01-02T15:04:05Z07:00")
			}

			timeIntValueList[k] = timeValuePair
		}
		return timeIntValueList
	}
	return nil
}
