package datapoliciesv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import3 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixStoragePoliciesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixStoragePoliciesV2Create,
		ReadContext:   ResourceNutanixStoragePoliciesV2Read,
		UpdateContext: ResourceNutanixStoragePoliciesV2Update,
		DeleteContext: ResourceNutanixStoragePoliciesV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category_ext_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"compression_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compression_state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DISABLED", "POSTPROCESS", "INLINE", "SYSTEM_DERIVED"}, false),
						},
					},
				},
			},
			"encryption_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption_state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SYSTEM_DERIVED", "ENABLED"}, false),
						},
					},
				},
			},
			"qos_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"throttled_iops": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"fault_tolerance_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replication_factor": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SYSTEM_DERIVED", "TWO", "THREE"}, false),
						},
					},
				},
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixStoragePoliciesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	body := &import1.StoragePolicy{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}

	if categoryExtIds, ok := d.GetOk("category_ext_ids"); ok {
		categoriesList := common.InterfaceToSlice(categoryExtIds)
		body.CategoryExtIds = common.ExpandListOfString(categoriesList)
	}

	if compressionSpec, ok := d.GetOk("compression_spec"); ok && len(compressionSpec.([]interface{})) > 0 {
		body.CompressionSpec = buildCompressionSpec(compressionSpec.([]interface{})[0].(map[string]interface{}))
	}

	if encryptionSpec, ok := d.GetOk("encryption_spec"); ok && len(encryptionSpec.([]interface{})) > 0 {
		body.EncryptionSpec = buildEncryptionSpec(encryptionSpec.([]interface{})[0].(map[string]interface{}))
	}

	if faultToleranceSpec, ok := d.GetOk("fault_tolerance_spec"); ok && len(faultToleranceSpec.([]interface{})) > 0 {
		body.FaultToleranceSpec = buildFaultToleranceSpec(faultToleranceSpec.([]interface{})[0].(map[string]interface{}))
	}

	if qosSpec, ok := d.GetOk("qos_spec"); ok && len(qosSpec.([]interface{})) > 0 {
		body.QosSpec = buildQosSpec(qosSpec.([]interface{})[0].(map[string]interface{}))
	}

	// Helper function to check if a spec is SYSTEM_DERIVED safely
	isSystemDerived := func(spec interface{ GetName() string }) bool {
		if spec == nil {
			return false
		}
		return spec.GetName() == "SYSTEM_DERIVED"
	}

	// Check each spec safely
	compressionDerived, encryptionDerived, replicationDerived := true, true, true
	if body.CompressionSpec != nil {
		compressionDerived = isSystemDerived(body.CompressionSpec.CompressionState)
	}
	if body.EncryptionSpec != nil {
		encryptionDerived = isSystemDerived(body.EncryptionSpec.EncryptionState)
	}
	if body.FaultToleranceSpec != nil {
		replicationDerived = isSystemDerived(body.FaultToleranceSpec.ReplicationFactor)
	}

	// Are all system-derived? Only check those that exist
	allSystemDerived := compressionDerived && encryptionDerived && replicationDerived

	// Validate qos_spec presence
	if allSystemDerived && body.QosSpec == nil {
		return diag.Errorf("qos_spec must be provided when compression_state, encryption_state, and replication_factor are all SYSTEM_DERIVED")
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create Storage Policy Payload: %s", string(aJSON))
	res, err := conn.StoragePolicies.CreateStoragePolicy(body)
	if err != nil {
		return diag.Errorf("error while creating Storage Policy: %v", err)
	}

	TaskRef := res.Data.GetValue().(import2.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Storage Policy to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Storage Policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		return diag.Errorf("error while fetching vm UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

	return ResourceNutanixStoragePoliciesV2Read(ctx, d, meta)
}

func ResourceNutanixStoragePoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	resp, err := conn.StoragePolicies.GetStoragePolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while reading Storage Policy: %v", err)
	}
	body := resp.Data.GetValue().(import1.StoragePolicy)
	metadata := resp.Metadata
	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Read Storage Policy Response: %s", string(aJSON))
	return commonReadStateStoragePolicy(d, body, metadata)
}

func ResourceNutanixStoragePoliciesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI
	resp, err := conn.StoragePolicies.GetStoragePolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Storage Policy: %v", err)
	}
	etagValue := conn.StoragePolicies.ApiClient.GetEtag(resp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	updateSpec := resp.Data.GetValue().(import1.StoragePolicy)
	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("category_ext_ids") {
		categoriesList := common.InterfaceToSlice(d.Get("category_ext_ids"))
		updateSpec.CategoryExtIds = common.ExpandListOfString(categoriesList)
	}
	if d.HasChange("compression_spec") {
		updateSpec.CompressionSpec = buildCompressionSpec(d.Get("compression_spec").([]interface{})[0].(map[string]interface{}))
	}
	if d.HasChange("encryption_spec") {
		if updateSpec.EncryptionSpec.EncryptionState.GetName() == "ENABLED" {
			return diag.Errorf("Encryption value cannot be changed once enabled because it is not supported.")
		}
		updateSpec.EncryptionSpec = buildEncryptionSpec(d.Get("encryption_spec").([]interface{})[0].(map[string]interface{}))
	}
	if d.HasChange("qos_spec") {
		updateSpec.QosSpec = buildQosSpec(d.Get("qos_spec").([]interface{})[0].(map[string]interface{}))
	}
	if d.HasChange("fault_tolerance_spec") {
		updateSpec.FaultToleranceSpec = buildFaultToleranceSpec(d.Get("fault_tolerance_spec").([]interface{})[0].(map[string]interface{}))
	}
	// Policy type is not updatable, so we need to set it to nil
	updateSpec.PolicyType = nil

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] Update Storage Policy Payload: %s", string(aJSON))
	res, err := conn.StoragePolicies.UpdateStoragePolicyById(utils.StringPtr(d.Id()), &updateSpec, headers)
	if err != nil {
		return diag.Errorf("error while updating Storage Policy: %v", err)
	}

	return waitForTaskCompletion(ctx, d, meta, res, "update")
}

func ResourceNutanixStoragePoliciesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	// Fetch the e-tag
	resp, err := conn.StoragePolicies.GetStoragePolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Storage Policy: %v", err)
	}
	etagValue := conn.StoragePolicies.ApiClient.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	res, err := conn.StoragePolicies.DeleteStoragePolicyById(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while deleting Storage Policy: %v", err)
	}
	return waitForTaskCompletion(ctx, d, meta, res, "delete")
}

func flattenCompressionSpec(compressionSpec *import1.CompressionSpec) []interface{} {
	if compressionSpec == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"compression_state": compressionSpec.CompressionState.GetName(),
	}
	return []interface{}{m}
}

func flattenEncryptionSpec(encryptionSpec *import1.EncryptionSpec) []interface{} {
	if encryptionSpec == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"encryption_state": encryptionSpec.EncryptionState.GetName(),
	}
	return []interface{}{m}
}

func flattenQosSpec(qosSpec *import1.QosSpec) []interface{} {
	if qosSpec == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"throttled_iops": qosSpec.ThrottledIops,
	}
	return []interface{}{m}
}

func flattenFaultToleranceSpec(faultToleranceSpec *import1.FaultToleranceSpec) []interface{} {
	if faultToleranceSpec == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"replication_factor": faultToleranceSpec.ReplicationFactor.GetName(),
	}
	return []interface{}{m}
}

// Common helper functions
func buildCompressionSpec(specMap map[string]interface{}) *import1.CompressionSpec {
	compressionSpec := &import1.CompressionSpec{}
	if compressionState, ok := specMap["compression_state"]; ok {
		var cs import1.CompressionState
		err := cs.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, compressionState.(string))))
		if err == nil {
			compressionSpec.CompressionState = cs.Ref()
		}
	}
	return compressionSpec
}

func buildEncryptionSpec(specMap map[string]interface{}) *import1.EncryptionSpec {
	encryptionSpec := &import1.EncryptionSpec{}
	if encryptionState, ok := specMap["encryption_state"]; ok {
		var es import1.EncryptionState
		err := es.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, encryptionState.(string))))
		if err == nil {
			encryptionSpec.EncryptionState = es.Ref()
		}
	}
	return encryptionSpec
}

func buildQosSpec(specMap map[string]interface{}) *import1.QosSpec {
	qosSpec := &import1.QosSpec{}
	if throttledIops, ok := specMap["throttled_iops"]; ok {
		throttledIopsInt := throttledIops.(int)
		qosSpec.ThrottledIops = &throttledIopsInt
	}
	return qosSpec
}

func buildFaultToleranceSpec(specMap map[string]interface{}) *import1.FaultToleranceSpec {
	faultToleranceSpec := &import1.FaultToleranceSpec{}
	if replicationFactor, ok := specMap["replication_factor"]; ok {
		var rf import1.ReplicationFactor
		err := rf.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, replicationFactor.(string))))
		if err == nil {
			faultToleranceSpec.ReplicationFactor = rf.Ref()
		}
	}
	return faultToleranceSpec
}

func waitForTaskCompletion(ctx context.Context, d *schema.ResourceData, meta interface{}, res interface{}, operation string) diag.Diagnostics {
	TaskRef := res.(interface{ GetData() interface{} }).GetData().(import2.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Storage Policy (%s) to %s: %s", utils.StringValue(taskUUID), operation, errWaitTask)
	}

	log.Printf("[DEBUG] Storage Policy (%s) %s is successful", d.Id(), operation)
	if operation == "delete" {
		return nil
	}

	return ResourceNutanixStoragePoliciesV2Read(ctx, d, meta)
}

func commonReadStateStoragePolicy(d *schema.ResourceData, res import1.StoragePolicy, metadata *import3.ApiResponseMetadata) diag.Diagnostics {
	if res.ExtId != nil {
		if err := d.Set("ext_id", *res.ExtId); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.Name != nil {
		if err := d.Set("name", res.Name); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.CategoryExtIds != nil {
		if err := d.Set("category_ext_ids", res.CategoryExtIds); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.CompressionSpec != nil {
		if err := d.Set("compression_spec", flattenCompressionSpec(res.CompressionSpec)); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.EncryptionSpec != nil {
		if err := d.Set("encryption_spec", flattenEncryptionSpec(res.EncryptionSpec)); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.QosSpec != nil {
		if err := d.Set("qos_spec", flattenQosSpec(res.QosSpec)); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.FaultToleranceSpec != nil {
		if err := d.Set("fault_tolerance_spec", flattenFaultToleranceSpec(res.FaultToleranceSpec)); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.PolicyType != nil {
		if err := d.Set("policy_type", res.PolicyType.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.TenantId != nil {
		if err := d.Set("tenant_id", res.TenantId); err != nil {
			return diag.FromErr(err)
		}
	}
	if metadata.Links != nil {
		if err := d.Set("links", flattenLinks(metadata.Links)); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
