package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVMAntiAffinityPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVMAntiAffinityPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"updated_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"num_compliant_vms": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_non_compliant_vms": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_pending_vms": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixVMAntiAffinityPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")

	resp, err := conn.VMAntiAffinityPolicyAPIInstance.GetVmAntiAffinityPolicyById(utils.StringPtr(extID.(string)))

	if err != nil {
		return diag.Errorf("error while fetching policy : %v", err)
	}

	getResp := resp.Data.GetValue().(policies.VmAntiAffinityPolicy)

	flattenedPolicy := flattenVMAntiAffinityPolicyEntity(getResp)

	for k, v := range flattenedPolicy {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(flattenedPolicy["ext_id"].(string))

	return nil
}

func flattenVMAntiAffinityPolicyEntity(policy policies.VmAntiAffinityPolicy) map[string]interface{} {
	result := make(map[string]interface{})
	if policy.ExtId != nil {
		result["ext_id"] = *policy.ExtId
	}
	if policy.Name != nil {
		result["name"] = *policy.Name
	}
	if policy.Description != nil {
		result["description"] = *policy.Description
	}
	if policy.CreateTime != nil {
		result["create_time"] = utils.TimeStringValue(policy.CreateTime)
	}
	if policy.UpdateTime != nil {
		result["update_time"] = utils.TimeStringValue(policy.UpdateTime)
	}
	if policy.CreatedBy != nil && policy.CreatedBy.ExtId != nil {
		result["created_by"] = map[string]string{"ext_id": *policy.CreatedBy.ExtId}
	}
	if policy.UpdatedBy != nil && policy.UpdatedBy.ExtId != nil {
		result["updated_by"] = map[string]string{"ext_id": *policy.UpdatedBy.ExtId}
	}
	if policy.Categories != nil {
		result["categories"] = flattenPolicyCategoryReference(policy.Categories)
	}
	if policy.NumCompliantVms != nil {
		result["num_compliant_vms"] = *policy.NumCompliantVms
	}
	if policy.NumNonCompliantVms != nil {
		result["num_non_compliant_vms"] = *policy.NumNonCompliantVms
	}
	if policy.NumPendingVms != nil {
		result["num_pending_vms"] = *policy.NumPendingVms
	}
	return result
}
