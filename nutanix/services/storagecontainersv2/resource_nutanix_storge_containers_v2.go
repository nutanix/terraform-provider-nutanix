package storagecontainersv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	clustermgmtConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clsCommonConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	clsPrismConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	timePeriod = 1 * time.Minute
	timeSleep  = 2 * time.Minute
)

func ResourceNutanixStorageContainersV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixStorageContainersV2Create,
		ReadContext:   ResourceNutanixStorageContainersV2Read,
		UpdateContext: ResourceNutanixStorageContainersV2Update,
		DeleteContext: ResourceNutanixStorageContainersV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
			"container_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_pool_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_marked_for_removal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"max_capacity_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"logical_explicit_reserved_capacity_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"logical_implicit_reserved_capacity_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"logical_advertised_capacity_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"replication_factor": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"nfs_whitelist_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": resourceSchemaForValuePrefixLength(),
						"ipv6": resourceSchemaForValuePrefixLength(),
						"fqdn": resourceSchemaForFqdnValue(),
					},
				},
			},
			"erasure_code": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "OFF", "ON"}, false),
			},
			"is_inline_ec_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"has_higher_ec_fault_domain_preference": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"erasure_code_delay_secs": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"cache_deduplication": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "OFF", "ON"}, false),
			},
			"on_disk_dedup": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "OFF", "POST_PROCESS"}, false),
			},
			"is_compression_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"compression_delay_secs": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"is_internal": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_software_encryption_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"affinity_host_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ignore_small_files": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func ResourceNutanixStorageContainersV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	body := &clustermgmtConfig.StorageContainer{}

	clusterExtID := d.Get("cluster_ext_id")

	if extID, ok := d.GetOk("ext_id"); ok {
		body.ExtId = utils.StringPtr(extID.(string))
	}
	if containerExtID, ok := d.GetOk("container_ext_id"); ok {
		body.ContainerExtId = utils.StringPtr(containerExtID.(string))
	}
	if ownerExtID, ok := d.GetOk("owner_ext_id"); ok {
		body.OwnerExtId = utils.StringPtr(ownerExtID.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if logicalExplicitReservedCapacityBytes, ok := d.GetOk("logical_explicit_reserved_capacity_bytes"); ok {
		body.LogicalExplicitReservedCapacityBytes = utils.Int64Ptr(int64(logicalExplicitReservedCapacityBytes.(int)))
		log.Printf("[DEBUG] logicalAdvertisedCapacityBytes: %v", utils.Int64Ptr(int64(logicalExplicitReservedCapacityBytes.(int))))
	}
	if logicalAdvertisedCapacityBytes, ok := d.GetOk("logical_advertised_capacity_bytes"); ok {
		body.LogicalAdvertisedCapacityBytes = utils.Int64Ptr(int64(logicalAdvertisedCapacityBytes.(int)))
		log.Printf("[DEBUG] logical_explicit_reserved_capacity_bytes: %v", utils.Int64Ptr(int64(logicalAdvertisedCapacityBytes.(int))))
	}
	if replicationFactor, ok := d.GetOk("replication_factor"); ok {
		body.ReplicationFactor = utils.IntPtr(replicationFactor.(int))
		log.Printf("[DEBUG] replicationFactor: %v", utils.IntPtr(replicationFactor.(int)))
	}
	if nfsWhitelistAddresses, ok := d.GetOk("nfs_whitelist_addresses"); ok {
		body.NfsWhitelistAddress = expandNfsWhitelistAddresses(nfsWhitelistAddresses)
	}
	if erasureCode, ok := d.GetOk("erasure_code"); ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"NONE": two,
			"OFF":  three,
			"ON":   four,
		}
		pVal := subMap[erasureCode.(string)]
		p := clustermgmtConfig.ErasureCodeStatus(pVal.(int))
		body.ErasureCode = &p
	}
	if isInlineEcEnabled, ok := d.GetOk("is_inline_ec_enabled"); ok {
		body.IsInlineEcEnabled = utils.BoolPtr(isInlineEcEnabled.(bool))
	}
	if hasHigherEcFaultDomainPreference, ok := d.GetOk("has_higher_ec_fault_domain_preference"); ok {
		body.HasHigherEcFaultDomainPreference = utils.BoolPtr(hasHigherEcFaultDomainPreference.(bool))
	}
	if erasureCodeDelaySecs, ok := d.GetOk("erasure_code_delay_secs"); ok {
		body.ErasureCodeDelaySecs = utils.IntPtr(erasureCodeDelaySecs.(int))
	}
	if cacheDeduplication, ok := d.GetOk("cache_deduplication"); ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"NONE": two,
			"OFF":  three,
			"ON":   four,
		}
		pVal := subMap[cacheDeduplication.(string)]
		p := clustermgmtConfig.CacheDeduplication(pVal.(int))
		body.CacheDeduplication = &p
	}
	if onDiskDedup, ok := d.GetOk("on_disk_dedup"); ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"NONE":         two,
			"OFF":          three,
			"POST_PROCESS": four,
		}
		pVal := subMap[onDiskDedup.(string)]
		p := clustermgmtConfig.OnDiskDedup(pVal.(int))
		body.OnDiskDedup = &p
	}
	if isCompressionEnabled, ok := d.GetOk("is_compression_enabled"); ok {
		body.IsCompressionEnabled = utils.BoolPtr(isCompressionEnabled.(bool))
	}
	if compressionDelaySecs, ok := d.GetOk("compression_delay_secs"); ok {
		body.CompressionDelaySecs = utils.IntPtr(compressionDelaySecs.(int))
	}
	if isInternal, ok := d.GetOk("is_internal"); ok {
		body.IsInternal = utils.BoolPtr(isInternal.(bool))
	}
	if isSoftwareEncryptionEnabled, ok := d.GetOk("is_software_encryption_enabled"); ok {
		body.IsSoftwareEncryptionEnabled = utils.BoolPtr(isSoftwareEncryptionEnabled.(bool))
	}
	if affinityHostExtID, ok := d.GetOk("affinity_host_ext_id"); ok {
		body.AffinityHostExtId = utils.StringPtr(affinityHostExtID.(string))
	}

	jsonBody, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] create storage container body: %s", string(jsonBody))
	resp, err := conn.StorageContainersAPI.CreateStorageContainer(body, utils.StringPtr(clusterExtID.(string)))
	if err != nil {
		return diag.Errorf("error while creating storage containers : %v", err)
	}

	TaskRef := resp.Data.GetValue().(clsPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the storage container to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for storage container (%s) to be created: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching storage container create task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Storage Container Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeStorageContainer, "Storage container")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	// Delay/sleep for 2 Minute, replication factor is not updated immediately
	time.Sleep(timeSleep)

	return ResourceNutanixStorageContainersV2Read(ctx, d, meta)
}

func ResourceNutanixStorageContainersV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading storage container with ID: %s", d.Id())
	conn := meta.(*conns.Client).ClusterAPI

	resp, err := conn.StorageContainersAPI.GetStorageContainerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching storage container: %v", err)
	}

	getResp := resp.Data.GetValue().(clustermgmtConfig.StorageContainer)

	jsonBody, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] read storage container body: %s", string(jsonBody))

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("container_ext_id", getResp.ContainerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_ext_id", getResp.ClusterExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_pool_ext_id", getResp.StoragePoolExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_marked_for_removal", getResp.IsMarkedForRemoval); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("max_capacity_bytes", getResp.MaxCapacityBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("logical_explicit_reserved_capacity_bytes", getResp.LogicalExplicitReservedCapacityBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("logical_implicit_reserved_capacity_bytes", getResp.LogicalImplicitReservedCapacityBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("logical_advertised_capacity_bytes", getResp.LogicalAdvertisedCapacityBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("replication_factor", getResp.ReplicationFactor); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nfs_whitelist_addresses", flattenNfsWhitelistAddresses(getResp.NfsWhitelistAddress)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("erasure_code", flattenErasureCodeStatus(getResp.ErasureCode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_inline_ec_enabled", getResp.IsInlineEcEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_higher_ec_fault_domain_preference", getResp.HasHigherEcFaultDomainPreference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("erasure_code_delay_secs", getResp.ErasureCodeDelaySecs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cache_deduplication", flattenCacheDeduplication(getResp.CacheDeduplication)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("on_disk_dedup", flattenOnDiskDedup(getResp.OnDiskDedup)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_compression_enabled", getResp.IsCompressionEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("compression_delay_secs", getResp.CompressionDelaySecs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_internal", getResp.IsInternal); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_software_encryption_enabled", getResp.IsSoftwareEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_encrypted", getResp.IsEncrypted); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("affinity_host_ext_id", getResp.AffinityHostExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_name", getResp.ClusterName); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixStorageContainersV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update Storage Container")
	conn := meta.(*conns.Client).ClusterAPI

	resp, err := conn.StorageContainersAPI.GetStorageContainerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching storage container : %v", err)
	}

	// Extract E-Tag Header
	etagValue := conn.ClusterEntityAPI.ApiClient.GetEtag(resp)

	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	respStorageContainer := resp.Data.GetValue().(clustermgmtConfig.StorageContainer)
	updateSpec := respStorageContainer

	if d.HasChange("ext_id") {
		updateSpec.ExtId = utils.StringPtr(d.Get("ext_id").(string))
	}
	if d.HasChange("container_ext_id") {
		updateSpec.ContainerExtId = utils.StringPtr(d.Get("container_ext_id").(string))
	}
	if d.HasChange("owner_ext_id") {
		updateSpec.OwnerExtId = utils.StringPtr(d.Get("owner_ext_id").(string))
	}
	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("logical_explicit_reserved_capacity_bytes") {
		updateSpec.LogicalExplicitReservedCapacityBytes = utils.Int64Ptr(int64(d.Get("logical_explicit_reserved_capacity_bytes").(int)))
	}
	if d.HasChange("logical_advertised_capacity_bytes") {
		updateSpec.LogicalAdvertisedCapacityBytes = utils.Int64Ptr(int64(d.Get("logical_advertised_capacity_bytes").(int)))
	}
	if d.HasChange("replication_factor") {
		updateSpec.ReplicationFactor = utils.IntPtr(d.Get("replication_factor").(int))
	}
	if d.HasChange("nfs_whitelist_addresses") {
		log.Printf("[DEBUG] nfs_whitelist_addresses: %v", d.Get("nfs_whitelist_addresses"))
		updateSpec.NfsWhitelistAddress = expandNfsWhitelistAddresses(d.Get("nfs_whitelist_addresses"))
	}
	if d.HasChange("erasure_code") {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"NONE": two,
			"OFF":  three,
			"ON":   four,
		}
		pVal := subMap[d.Get("erasure_code").(string)]
		p := clustermgmtConfig.ErasureCodeStatus(pVal.(int))
		updateSpec.ErasureCode = &p
	}
	if d.HasChange("is_inline_ec_enabled") {
		updateSpec.IsInlineEcEnabled = utils.BoolPtr(d.Get("is_inline_ec_enabled").(bool))
	}
	if d.HasChange("has_higher_ec_fault_domain_preference") {
		updateSpec.HasHigherEcFaultDomainPreference = utils.BoolPtr(d.Get("has_higher_ec_fault_domain_preference").(bool))
	}
	if d.HasChange("erasure_code_delay_secs") {
		updateSpec.ErasureCodeDelaySecs = utils.IntPtr(d.Get("erasure_code_delay_secs").(int))
	}
	if d.HasChange("cache_deduplication") {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"NONE": two,
			"OFF":  three,
			"ON":   four,
		}
		pVal := subMap[d.Get("cache_deduplication").(string)]
		p := clustermgmtConfig.CacheDeduplication(pVal.(int))
		updateSpec.CacheDeduplication = &p
	}
	if d.HasChange("on_disk_dedup") {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"NONE":         two,
			"OFF":          three,
			"POST_PROCESS": four,
		}
		pVal := subMap[d.Get("on_disk_dedup").(string)]
		p := clustermgmtConfig.OnDiskDedup(pVal.(int))
		updateSpec.OnDiskDedup = &p
	}
	if d.HasChange("is_compression_enabled") {
		updateSpec.IsCompressionEnabled = utils.BoolPtr(d.Get("is_compression_enabled").(bool))
	}
	if d.HasChange("compression_delay_secs") {
		updateSpec.CompressionDelaySecs = utils.IntPtr(d.Get("compression_delay_secs").(int))
	}
	if d.HasChange("is_internal") {
		updateSpec.IsInternal = utils.BoolPtr(d.Get("is_internal").(bool))
	}
	if d.HasChange("is_software_encryption_enabled") {
		updateSpec.IsSoftwareEncryptionEnabled = utils.BoolPtr(d.Get("is_software_encryption_enabled").(bool))
	}
	if d.HasChange("affinity_host_ext_id") {
		updateSpec.AffinityHostExtId = utils.StringPtr(d.Get("affinity_host_ext_id").(string))
	}

	updateResp, err := conn.StorageContainersAPI.UpdateStorageContainerById(utils.StringPtr(d.Id()), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating storage container : %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(clsPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the storage container to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for storage container (%s) to be updated: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching storage container update task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Storage Container Update Task Details: %s", string(aJSON))

	// delay/sleep for 1 Minute, replication factor is not updated immediately
	time.Sleep(timePeriod)
	return ResourceNutanixStorageContainersV2Read(ctx, d, meta)
}

func ResourceNutanixStorageContainersV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	// default value for ignoreSmallFiles is true
	ignoreSmallFiles := true
	if ignoreSmallFile, ok := d.GetOk("ignore_small_files"); ok {
		ignoreSmallFiles = *utils.BoolPtr(ignoreSmallFile.(bool))
	}

	resp, err := conn.StorageContainersAPI.DeleteStorageContainerById(utils.StringPtr(d.Id()), utils.BoolPtr(ignoreSmallFiles))
	if err != nil {
		return diag.Errorf("error while deleting storage container: %v", err)
	}

	TaskRef := resp.Data.GetValue().(clsPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the storage container to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for storage container (%s) to be deleted: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details for logging
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching storage container delete task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete Storage Container Task Details: %s", string(aJSON))

	return nil
}

func expandNfsWhitelistAddresses(nfsWhitelistAddresses interface{}) []clsCommonConfig.IPAddressOrFQDN {
	if nfsWhitelistAddresses != nil {
		nfsWhitelistAddressesList := nfsWhitelistAddresses.([]interface{})
		ips := make([]clsCommonConfig.IPAddressOrFQDN, len(nfsWhitelistAddressesList))

		ip := &clsCommonConfig.IPAddressOrFQDN{}
		prI := nfsWhitelistAddresses.([]interface{})
		val := prI[0].(map[string]interface{})

		if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
			ip.Ipv4 = expandIPv4Address(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
			log.Printf("[DEBUG] ipv6: %v", ipv6)

			ip.Ipv6 = expandIPv6Address(ipv6)
		}
		if fqdn, ok := val["fqdn"]; ok && len(fqdn.([]interface{})) > 0 {
			ip.Fqdn = expandFQDN(fqdn.([]interface{}))
		}
		ips[0] = *ip
		return ips
	}
	return nil
}

func resourceSchemaForValuePrefixLength() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"prefix_length": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func resourceSchemaForFqdnValue() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func expandIPv4Address(pr interface{}) *clsCommonConfig.IPv4Address {
	if pr != nil {
		ipv4 := &clsCommonConfig.IPv4Address{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv4.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv4
	}
	return nil
}

func expandIPv6Address(pr interface{}) *clsCommonConfig.IPv6Address {
	if pr != nil {
		ipv6 := &clsCommonConfig.IPv6Address{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv6.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv6.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv6
	}
	return nil
}

func expandFQDN(pr []interface{}) *clsCommonConfig.FQDN {
	if len(pr) > 0 {
		fqdn := clsCommonConfig.FQDN{}
		val := pr[0].(map[string]interface{})
		if value, ok := val["value"]; ok {
			fqdn.Value = utils.StringPtr(value.(string))
		}

		return &fqdn
	}
	return nil
}
