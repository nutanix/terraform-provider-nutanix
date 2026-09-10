package vmmv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vmstartuppolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	vmm "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dependencyConflictSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dependee_category": {
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
		"dependent_category": {
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
		"dependee_vms_associated_categories": {
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
		"dependent_vms_associated_categories": {
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
		"category_dependency_chain": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"dependee_category": {
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
					"dependent_category": {
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
					"policy": {
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
				},
			},
		},
		"dependee_vms": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of dependee VMs involved in this conflict.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ext_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"dependent_vms": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of dependent VMs involved in this conflict.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ext_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
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
	}
}

func DatasourceNutanixVmStartupPolicyDependencyConflictV2() *schema.Resource {
	s := dependencyConflictSchema()
	s["vm_startup_policy_ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	s["ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyDependencyConflictV2Read,
		Schema:      s,
	}
}

func datasourceNutanixVmStartupPolicyDependencyConflictV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)
	extID := d.Get("ext_id").(string)

	getRequest := import2.GetVmStartupPolicyDependencyConflictByIdRequest{
		VmStartupPolicyExtId: utils.StringPtr(policyExtID),
		ExtId:                utils.StringPtr(extID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.GetVmStartupPolicyDependencyConflictById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching VM Startup Policy Dependency Conflict: %v", err)
	}

	conflict := resp.Data.GetValue().(import1.DependencyConflict)

	if err := flattenDependencyConflict(d, &conflict); err != nil {
		return diag.FromErr(err)
	}

	dependeeVms, dependentVms := fetchDependencyConflictVms(ctx, conn, policyExtID, extID)
	if err := d.Set("dependee_vms", flattenVmReferences(dependeeVms)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dependent_vms", flattenVmReferences(dependentVms)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extID)
	return nil
}

func fetchDependencyConflictVms(ctx context.Context, conn *vmm.Client, policyExtID, conflictExtID string) ([]import1.VmReference, []import1.VmReference) {
	var dependeeVms, dependentVms []import1.VmReference

	dependeeReq := import2.ListVmStartupPolicyDependencyConflictDependeeVmsRequest{
		VmStartupPolicyExtId:    utils.StringPtr(policyExtID),
		DependencyConflictExtId: utils.StringPtr(conflictExtID),
	}
	if resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyDependencyConflictDependeeVms(ctx, &dependeeReq); err == nil {
		dependeeVms = resp.Data.GetValue().([]import1.VmReference)
	} else {
		log.Printf("[WARN] Failed to fetch dependency conflict dependee VMs: %v", err)
	}

	dependentReq := import2.ListVmStartupPolicyDependencyConflictDependentVmsRequest{
		VmStartupPolicyExtId:    utils.StringPtr(policyExtID),
		DependencyConflictExtId: utils.StringPtr(conflictExtID),
	}
	if resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyDependencyConflictDependentVms(ctx, &dependentReq); err == nil {
		dependentVms = resp.Data.GetValue().([]import1.VmReference)
	} else {
		log.Printf("[WARN] Failed to fetch dependency conflict dependent VMs: %v", err)
	}

	return dependeeVms, dependentVms
}

func flattenDependencyConflict(d *schema.ResourceData, conflict *import1.DependencyConflict) error {
	if conflict.ExtId != nil {
		if err := d.Set("ext_id", utils.StringValue(conflict.ExtId)); err != nil {
			return err
		}
	}
	if conflict.TenantId != nil {
		if err := d.Set("tenant_id", utils.StringValue(conflict.TenantId)); err != nil {
			return err
		}
	}
	if err := d.Set("dependee_category", flattenCategoryRef(conflict.DependeeCategory)); err != nil {
		return err
	}
	if err := d.Set("dependent_category", flattenCategoryRef(conflict.DependentCategory)); err != nil {
		return err
	}
	if err := d.Set("dependee_vms_associated_categories", flattenCategoryRefList(conflict.DependeeVmsAssociatedCategories)); err != nil {
		return err
	}
	if err := d.Set("dependent_vms_associated_categories", flattenCategoryRefList(conflict.DependentVmsAssociatedCategories)); err != nil {
		return err
	}
	if err := d.Set("category_dependency_chain", flattenCategoryDependencyChain(conflict.CategoryDependencyChain)); err != nil {
		return err
	}
	if err := d.Set("links", flattenVmStartupPolicyLinks(conflict.Links)); err != nil {
		return err
	}
	return nil
}

func flattenCategoryRef(ref *import1.CategoryReference) []map[string]interface{} {
	if ref == nil {
		return nil
	}
	result := map[string]interface{}{}
	if ref.ExtId != nil {
		result["ext_id"] = utils.StringValue(ref.ExtId)
	}
	return []map[string]interface{}{result}
}

func flattenCategoryRefList(refs []import1.CategoryReference) []map[string]interface{} {
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

func flattenCategoryDependencyChain(chain []import1.CategoryDependency) []map[string]interface{} {
	if len(chain) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(chain))
	for i, cd := range chain {
		m := map[string]interface{}{}
		m["dependee_category"] = flattenCategoryRef(cd.DependeeCategory)
		m["dependent_category"] = flattenCategoryRef(cd.DependentCategory)
		if cd.Policy != nil {
			policyMap := map[string]interface{}{}
			if cd.Policy.ExtId != nil {
				policyMap["ext_id"] = utils.StringValue(cd.Policy.ExtId)
			}
			m["policy"] = []map[string]interface{}{policyMap}
		}
		result[i] = m
	}
	return result
}

func flattenDependencyConflictToMap(conflict *import1.DependencyConflict) map[string]interface{} {
	result := map[string]interface{}{}
	if conflict.ExtId != nil {
		result["ext_id"] = utils.StringValue(conflict.ExtId)
	}
	if conflict.TenantId != nil {
		result["tenant_id"] = utils.StringValue(conflict.TenantId)
	}
	result["dependee_category"] = flattenCategoryRef(conflict.DependeeCategory)
	result["dependent_category"] = flattenCategoryRef(conflict.DependentCategory)
	result["dependee_vms_associated_categories"] = flattenCategoryRefList(conflict.DependeeVmsAssociatedCategories)
	result["dependent_vms_associated_categories"] = flattenCategoryRefList(conflict.DependentVmsAssociatedCategories)
	result["category_dependency_chain"] = flattenCategoryDependencyChain(conflict.CategoryDependencyChain)
	result["links"] = flattenVmStartupPolicyLinks(conflict.Links)
	return result
}
