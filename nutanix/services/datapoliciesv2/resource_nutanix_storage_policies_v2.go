package datapoliciesv2

import (
	"context"


	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	 conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	 import1 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/config"
	 "github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixStoragePoliciesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixStoragePoliciesV2Create,
		ReadContext:   ResourceNutanixStoragePoliciesV2Read,
		UpdateContext: ResourceNutanixStoragePoliciesV2Update,
		DeleteContext: ResourceNutanixStoragePoliciesV2Delete,
		Schema: map[string]*schema.Schema{
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
							Type:     schema.TypeString,
							Required: true,
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
							Type:     schema.TypeString,
							Required: true,
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
		 compressionSpecMap := v.(map[string]interface{})
		 compressionSpec := &import1.CompressionSpec{}
		 if compressionState, ok := compressionSpecMap["compression_state"]; ok {
			 compressionSpec.CompressionState = utils.StringPtr(compressionState.(string))
		 }
		 body.CompressionSpec = compressionSpec
	 }
	 if v, ok := d.GetOk("encryption_spec"); ok {
		 encryptionSpecMap := v.(map[string]interface{})
		 encryptionSpec := &import1.EncryptionSpec{}
		 if encryptionState, ok := encryptionSpecMap["encryption_state"]; ok {
			 encryptionSpec.EncryptionState = encryptionState.(string)
		 }
		 body.EncryptionSpec = encryptionSpec
	 }
	 if v, ok := d.GetOk("qos_spec"); ok {
		 qosSpecMap := v.(map[string]interface{})
		 qosSpec := &import1.QosSpec{}
		 if throttledIops, ok := qosSpecMap["throttled_iops"]; ok {
			 throttledIopsInt := int64(throttledIops.(int))
			 qosSpec.ThrottledIops = &throttledIopsInt
		 }
		 body.QosSpec = qosSpec
	 }
	 if v, ok := d.GetOk("fault_tolerance_spec"); ok {
		 faultToleranceSpecMap := v.(map[string]interface{})
		 faultToleranceSpec := &import1.FaultToleranceSpec{}
		 if replicationFactor, ok := faultToleranceSpecMap["replication_factor"]; ok {
			 faultToleranceSpec.ReplicationFactor = int64(replicationFactor.(int))
		 }
		 body.FaultToleranceSpec = faultToleranceSpec
	 }

	 if policyType, ok := d.GetOk("policy_type"); ok {
		 body.PolicyType = policyType.(int)
	 }

	 res, err := conn.StoragePolicies.CreateStoragePolicy(body)
	 if err != nil {
		 return diag.Errorf("error while creating Storage Policy: %v", err)
	 }
}

func ResourceNutanixStoragePoliciesV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {}

func ResourceNutanixStoragePoliciesV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {}

func ResourceNutanixStoragePoliciesV2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {}