package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vmstartuppolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func vmReferenceListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vm_startup_policy_ext_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"vms": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ext_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func flattenVmReferences(refs []import1.VmReference) []map[string]interface{} {
	if len(refs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(refs))
	for i, r := range refs {
		m := map[string]interface{}{}
		if r.ExtId != nil {
			m["ext_id"] = utils.StringValue(r.ExtId)
		}
		result[i] = m
	}
	return result
}

// Dependency Conflict Dependee VMs
func DatasourceNutanixVmStartupPolicyDependencyConflictDependeeVmsV2() *schema.Resource {
	s := vmReferenceListSchema()
	s["dependency_conflict_ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyDependencyConflictDependeeVmsV2Read,
		Schema:      s,
	}
}

func datasourceNutanixVmStartupPolicyDependencyConflictDependeeVmsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)
	conflictExtID := d.Get("dependency_conflict_ext_id").(string)

	listRequest := import2.ListVmStartupPolicyDependencyConflictDependeeVmsRequest{
		VmStartupPolicyExtId:    utils.StringPtr(policyExtID),
		DependencyConflictExtId: utils.StringPtr(conflictExtID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyDependencyConflictDependeeVms(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing Dependency Conflict Dependee VMs: %v", err)
	}

	vms := resp.Data.GetValue().([]import1.VmReference)

	if err := d.Set("vms", flattenVmReferences(vms)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// Dependency Conflict Dependent VMs
func DatasourceNutanixVmStartupPolicyDependencyConflictDependentVmsV2() *schema.Resource {
	s := vmReferenceListSchema()
	s["dependency_conflict_ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyDependencyConflictDependentVmsV2Read,
		Schema:      s,
	}
}

func datasourceNutanixVmStartupPolicyDependencyConflictDependentVmsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)
	conflictExtID := d.Get("dependency_conflict_ext_id").(string)

	listRequest := import2.ListVmStartupPolicyDependencyConflictDependentVmsRequest{
		VmStartupPolicyExtId:    utils.StringPtr(policyExtID),
		DependencyConflictExtId: utils.StringPtr(conflictExtID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyDependencyConflictDependentVms(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing Dependency Conflict Dependent VMs: %v", err)
	}

	vms := resp.Data.GetValue().([]import1.VmReference)

	if err := d.Set("vms", flattenVmReferences(vms)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// Start Condition Conflict Dependee VMs
func DatasourceNutanixVmStartupPolicyStartConditionConflictDependeeVmsV2() *schema.Resource {
	s := vmReferenceListSchema()
	s["start_condition_conflict_ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyStartConditionConflictDependeeVmsV2Read,
		Schema:      s,
	}
}

func datasourceNutanixVmStartupPolicyStartConditionConflictDependeeVmsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)
	conflictExtID := d.Get("start_condition_conflict_ext_id").(string)

	listRequest := import2.ListVmStartupPolicyStartConditionConflictDependeeVmsRequest{
		VmStartupPolicyExtId:        utils.StringPtr(policyExtID),
		StartConditionConflictExtId: utils.StringPtr(conflictExtID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyStartConditionConflictDependeeVms(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing Start Condition Conflict Dependee VMs: %v", err)
	}

	vms := resp.Data.GetValue().([]import1.VmReference)

	if err := d.Set("vms", flattenVmReferences(vms)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// Start Condition Conflict Dependent VMs
func DatasourceNutanixVmStartupPolicyStartConditionConflictDependentVmsV2() *schema.Resource {
	s := vmReferenceListSchema()
	s["start_condition_conflict_ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyStartConditionConflictDependentVmsV2Read,
		Schema:      s,
	}
}

func datasourceNutanixVmStartupPolicyStartConditionConflictDependentVmsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)
	conflictExtID := d.Get("start_condition_conflict_ext_id").(string)

	listRequest := import2.ListVmStartupPolicyStartConditionConflictDependentVmsRequest{
		VmStartupPolicyExtId:        utils.StringPtr(policyExtID),
		StartConditionConflictExtId: utils.StringPtr(conflictExtID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyStartConditionConflictDependentVms(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing Start Condition Conflict Dependent VMs: %v", err)
	}

	vms := resp.Data.GetValue().([]import1.VmReference)

	if err := d.Set("vms", flattenVmReferences(vms)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
