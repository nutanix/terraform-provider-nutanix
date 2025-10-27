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
	import1 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
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
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"compression_spec": {
				Type:     schema.TypeMap,
				Optional: true,
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
				Type:     schema.TypeMap,
				Optional: true,
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
				Type:     schema.TypeMap,
				Optional: true,
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
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replication_factor": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "USER",
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
		body.CategoryExtIds = []string{}
		for _, id := range categoryExtIds.([]interface{}) {
			body.CategoryExtIds = append(body.CategoryExtIds, id.(string))
		}
	}

	if v, ok := d.GetOk("compression_spec"); ok {
		body.CompressionSpec = buildCompressionSpec(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("encryption_spec"); ok {
		body.EncryptionSpec = buildEncryptionSpec(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("qos_spec"); ok {
		body.QosSpec = buildQosSpec(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("fault_tolerance_spec"); ok {
		body.FaultToleranceSpec = buildFaultToleranceSpec(v.(map[string]interface{}))
	}

	if policyType, ok := d.GetOk("policy_type"); ok {
		body.PolicyType = buildPolicyType(policyType.(string))
	}

	res, err := conn.StoragePolicies.CreateStoragePolicy(body)
	if err != nil {
		return diag.Errorf("error while creating Storage Policy: %v", err)
	}

	TaskRef := res.Data.GetValue().(import2.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
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

 	return commonReadStateStoragePolicy(ctx, d, meta, body)


}

func ResourceNutanixStoragePoliciesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	body := &import1.StoragePolicy{}
	if d.HasChange("name") {
		body.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("category_ext_ids") {
		body.CategoryExtIds = []string{}
		for _, id := range d.Get("category_ext_ids").([]interface{}) {
			body.CategoryExtIds = append(body.CategoryExtIds, id.(string))
		}
	}
	if d.HasChange("compression_spec") {
		body.CompressionSpec = buildCompressionSpec(d.Get("compression_spec").(map[string]interface{}))
	}
	if d.HasChange("encryption_spec") {
		body.EncryptionSpec = buildEncryptionSpec(d.Get("encryption_spec").(map[string]interface{}))
	}
	if d.HasChange("qos_spec") {
		body.QosSpec = buildQosSpec(d.Get("qos_spec").(map[string]interface{}))
	}
	if d.HasChange("fault_tolerance_spec") {
		body.FaultToleranceSpec = buildFaultToleranceSpec(d.Get("fault_tolerance_spec").(map[string]interface{}))
	}
	if d.HasChange("policy_type") {
		body.PolicyType = buildPolicyType(d.Get("policy_type").(string))
	}

	resp, err := conn.StoragePolicies.GetStoragePolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Storage Policy: %v", err)
	}
	etagValue := conn.StoragePolicies.ApiClient.GetEtag(resp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	res, err := conn.StoragePolicies.UpdateStoragePolicyById(utils.StringPtr(d.Id()), body, headers)
	if err != nil {
		return diag.Errorf("error while updating Storage Policy: %v", err)
	}

	return waitForTaskCompletion(ctx, d, meta, res, "update")
}

func ResourceNutanixStoragePoliciesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	res, err := conn.StoragePolicies.DeleteStoragePolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting Storage Policy: %v", err)
	}
	return waitForTaskCompletion(ctx, d, meta, res, "delete")
}

func flattenCompressionSpec(compressionSpec *import1.CompressionSpec) map[string]interface{} {
	return map[string]interface{}{
		"compression_state": compressionSpec.CompressionState.GetName(),
	}
}

func flattenEncryptionSpec(encryptionSpec *import1.EncryptionSpec) map[string]interface{} {
	return map[string]interface{}{
		"encryption_state": encryptionSpec.EncryptionState.GetName(),
	}
}

func flattenQosSpec(qosSpec *import1.QosSpec) map[string]interface{} {
	return map[string]interface{}{
		"throttled_iops": qosSpec.ThrottledIops,
	}
}

func flattenFaultToleranceSpec(faultToleranceSpec *import1.FaultToleranceSpec) map[string]interface{} {
	return map[string]interface{}{
		"replication_factor": faultToleranceSpec.ReplicationFactor.GetName(),
	}
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
		throttledIopsInt := int(throttledIops.(int))
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

func buildPolicyType(policyTypeStr string) *import1.PolicyType {
	var pt import1.PolicyType
	err := pt.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, policyTypeStr)))
	if err == nil {
		return pt.Ref()
	}
	return nil
}

func waitForTaskCompletion(ctx context.Context, d *schema.ResourceData, meta interface{}, res interface{}, operation string) diag.Diagnostics {
	TaskRef := res.(interface{ GetData() interface{} }).GetData().(import2.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Storage Policy (%s) to %s: %s", utils.StringValue(taskUUID), operation, errWaitTask)
	}

	log.Printf("[DEBUG] Storage Policy (%s) %s successfully", utils.StringValue(taskUUID), operation)
	if operation == "delete" {
		return nil
	}

	return ResourceNutanixStoragePoliciesV2Read(ctx, d, meta)
}

func commonReadStateStoragePolicy(ctx context.Context, d *schema.ResourceData, meta interface{}, res import1.StoragePolicy) diag.Diagnostics {
	if res.ExtId != nil {
		d.Set("ext_id", *res.ExtId)
	}
	if res.Name != nil {
		d.Set("name", res.Name)
	}
	if res.CategoryExtIds != nil {
		d.Set("category_ext_ids", res.CategoryExtIds)
	}
	if res.CompressionSpec != nil {
		d.Set("compression_spec", flattenCompressionSpec(res.CompressionSpec))
	}
	if res.EncryptionSpec != nil {
		d.Set("encryption_spec", flattenEncryptionSpec(res.EncryptionSpec))
	}
	if res.QosSpec != nil {
		d.Set("qos_spec", flattenQosSpec(res.QosSpec))
	}
	if res.FaultToleranceSpec != nil {
		d.Set("fault_tolerance_spec", flattenFaultToleranceSpec(res.FaultToleranceSpec))
	}
	if res.PolicyType != nil {
		d.Set("policy_type", res.PolicyType.GetName())
	}
	if res.Links != nil {
		d.Set("links", flattenLinks(res.Links))
	}
	if res.TenantId != nil {
		d.Set("tenant_id", res.TenantId)
	}
	return nil
}