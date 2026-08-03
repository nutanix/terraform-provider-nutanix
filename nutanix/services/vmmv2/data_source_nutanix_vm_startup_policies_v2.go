package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/vmstartuppolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVmStartupPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPoliciesV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixVmStartupPolicyV2(),
			},
		},
	}
}

func datasourceNutanixVmStartupPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	listRequest := import2.ListVmStartupPoliciesRequest{}

	if pagef, ok := d.GetOk("page"); ok {
		listRequest.Page_ = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		listRequest.Limit_ = utils.IntPtr(limitf.(int))
	}
	if filterf, ok := d.GetOk("filter"); ok {
		listRequest.Filter_ = utils.StringPtr(filterf.(string))
	}
	if order, ok := d.GetOk("order_by"); ok {
		listRequest.Orderby_ = utils.StringPtr(order.(string))
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicies(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing VM Startup Policies: %v", err)
	}

	policies := resp.Data.GetValue().([]import1.VmStartupPolicy)

	policyList := make([]map[string]interface{}, len(policies))
	for i, p := range policies {
		policyList[i] = flattenVmStartupPolicyToMap(&p)
	}

	if err := d.Set("policies", policyList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVmStartupPolicyToMap(policy *import1.VmStartupPolicy) map[string]interface{} {
	result := map[string]interface{}{}
	if policy.ExtId != nil {
		result["ext_id"] = utils.StringValue(policy.ExtId)
	}
	if policy.Name != nil {
		result["name"] = utils.StringValue(policy.Name)
	}
	if policy.Description != nil {
		result["description"] = utils.StringValue(policy.Description)
	}
	if policy.CreateTime != nil {
		result["create_time"] = utils.TimeStringValue(policy.CreateTime)
	}
	if policy.UpdateTime != nil {
		result["update_time"] = utils.TimeStringValue(policy.UpdateTime)
	}
	if policy.ProjectExtId != nil {
		result["project_ext_id"] = utils.StringValue(policy.ProjectExtId)
	}
	if policy.TenantId != nil {
		result["tenant_id"] = utils.StringValue(policy.TenantId)
	}
	if policy.NumCompliantVms != nil {
		result["num_compliant_vms"] = int(*policy.NumCompliantVms)
	}
	if policy.NumNonCompliantVms != nil {
		result["num_non_compliant_vms"] = int(*policy.NumNonCompliantVms)
	}
	if policy.NumPendingVms != nil {
		result["num_pending_vms"] = int(*policy.NumPendingVms)
	}
	if policy.NumDependencyConflicts != nil {
		result["num_dependency_conflicts"] = int(*policy.NumDependencyConflicts)
	}
	if policy.NumStartConditionConflicts != nil {
		result["num_start_condition_conflicts"] = int(*policy.NumStartConditionConflicts)
	}
	result["created_by"] = flattenUserReference(policy.CreatedBy)
	result["updated_by"] = flattenUserReference(policy.UpdatedBy)
	result["groups"] = flattenVmStartupPolicyGroups(policy.Groups)
	result["start_conditions"] = flattenVmStartupPolicyStartConditions(policy.StartConditions)
	result["links"] = flattenVmStartupPolicyLinks(policy.Links)
	return result
}
