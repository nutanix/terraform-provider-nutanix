package vmmv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/vmstartuppolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	vmm "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func startConditionConflictSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"conflicting_policy": {
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
		"conflicting_start_condition": startConditionSchemaComputed(),
		"start_condition":             startConditionSchemaComputed(),
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

func startConditionSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"delay_duration_secs": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"power_state_criteria": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"power_on": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{},
								},
							},
							"guest_bootup": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"timeout_duration_secs": {
											Type:     schema.TypeInt,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVmStartupPolicyStartConditionConflictV2() *schema.Resource {
	s := startConditionConflictSchema()
	s["vm_startup_policy_ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	s["ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyStartConditionConflictV2Read,
		Schema:      s,
	}
}

func datasourceNutanixVmStartupPolicyStartConditionConflictV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)
	extID := d.Get("ext_id").(string)

	getRequest := import2.GetVmStartupPolicyStartConditionConflictByIdRequest{
		VmStartupPolicyExtId: utils.StringPtr(policyExtID),
		ExtId:                utils.StringPtr(extID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.GetVmStartupPolicyStartConditionConflictById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching VM Startup Policy Start Condition Conflict: %v", err)
	}

	conflict := resp.Data.GetValue().(import1.StartConditionConflict)

	if err := flattenStartConditionConflict(d, &conflict); err != nil {
		return diag.FromErr(err)
	}

	dependeeVms, dependentVms := fetchStartConditionConflictVms(ctx, conn, policyExtID, extID)
	if err := d.Set("dependee_vms", flattenVmReferences(dependeeVms)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dependent_vms", flattenVmReferences(dependentVms)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extID)
	return nil
}

func fetchStartConditionConflictVms(ctx context.Context, conn *vmm.Client, policyExtID, conflictExtID string) ([]import1.VmReference, []import1.VmReference) {
	var dependeeVms, dependentVms []import1.VmReference

	dependeeReq := import2.ListVmStartupPolicyStartConditionConflictDependeeVmsRequest{
		VmStartupPolicyExtId:        utils.StringPtr(policyExtID),
		StartConditionConflictExtId: utils.StringPtr(conflictExtID),
	}
	if resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyStartConditionConflictDependeeVms(ctx, &dependeeReq); err == nil {
		dependeeVms = resp.Data.GetValue().([]import1.VmReference)
	} else {
		log.Printf("[WARN] Failed to fetch start condition conflict dependee VMs: %v", err)
	}

	dependentReq := import2.ListVmStartupPolicyStartConditionConflictDependentVmsRequest{
		VmStartupPolicyExtId:        utils.StringPtr(policyExtID),
		StartConditionConflictExtId: utils.StringPtr(conflictExtID),
	}
	if resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyStartConditionConflictDependentVms(ctx, &dependentReq); err == nil {
		dependentVms = resp.Data.GetValue().([]import1.VmReference)
	} else {
		log.Printf("[WARN] Failed to fetch start condition conflict dependent VMs: %v", err)
	}

	return dependeeVms, dependentVms
}

func flattenStartConditionConflict(d *schema.ResourceData, conflict *import1.StartConditionConflict) error {
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
	if conflict.ConflictingPolicy != nil {
		policyMap := map[string]interface{}{}
		if conflict.ConflictingPolicy.ExtId != nil {
			policyMap["ext_id"] = utils.StringValue(conflict.ConflictingPolicy.ExtId)
		}
		if err := d.Set("conflicting_policy", []map[string]interface{}{policyMap}); err != nil {
			return err
		}
	}
	if err := d.Set("conflicting_start_condition", flattenStartConditionSingle(conflict.ConflictingStartCondition)); err != nil {
		return err
	}
	if err := d.Set("start_condition", flattenStartConditionSingle(conflict.StartCondition)); err != nil {
		return err
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
	if err := d.Set("links", flattenVmStartupPolicyLinks(conflict.Links)); err != nil {
		return err
	}
	return nil
}

func flattenStartConditionSingle(sc *import1.StartCondition) []map[string]interface{} {
	if sc == nil {
		return nil
	}
	sMap := map[string]interface{}{}
	if sc.DelayDurationSecs != nil {
		sMap["delay_duration_secs"] = *sc.DelayDurationSecs
	}
	if sc.PowerStateCriteria != nil && sc.PowerStateCriteria.ObjectType_ != nil {
		pscMap := map[string]interface{}{}
		val := sc.PowerStateCriteria.GetValue()
		if val != nil {
			switch *sc.PowerStateCriteria.ObjectType_ {
			case "vmm.v4.ahv.policies.PowerStateCriteriaPowerOn":
				pscMap["power_on"] = []map[string]interface{}{{}}
				pscMap["guest_bootup"] = []map[string]interface{}{}
			case "vmm.v4.ahv.policies.PowerStateCriteriaGuestBootup":
				if gb, ok := val.(import1.PowerStateCriteriaGuestBootup); ok {
					gbMap := map[string]interface{}{}
					if gb.TimeoutDurationSecs != nil {
						gbMap["timeout_duration_secs"] = *gb.TimeoutDurationSecs
					}
					pscMap["guest_bootup"] = []map[string]interface{}{gbMap}
				}
				pscMap["power_on"] = []map[string]interface{}{}
			}
		}
		sMap["power_state_criteria"] = []map[string]interface{}{pscMap}
	}
	return []map[string]interface{}{sMap}
}

func flattenStartConditionConflictToMap(conflict *import1.StartConditionConflict) map[string]interface{} {
	result := map[string]interface{}{}
	if conflict.ExtId != nil {
		result["ext_id"] = utils.StringValue(conflict.ExtId)
	}
	if conflict.TenantId != nil {
		result["tenant_id"] = utils.StringValue(conflict.TenantId)
	}
	if conflict.ConflictingPolicy != nil {
		policyMap := map[string]interface{}{}
		if conflict.ConflictingPolicy.ExtId != nil {
			policyMap["ext_id"] = utils.StringValue(conflict.ConflictingPolicy.ExtId)
		}
		result["conflicting_policy"] = []map[string]interface{}{policyMap}
	}
	result["conflicting_start_condition"] = flattenStartConditionSingle(conflict.ConflictingStartCondition)
	result["start_condition"] = flattenStartConditionSingle(conflict.StartCondition)
	result["dependee_category"] = flattenCategoryRef(conflict.DependeeCategory)
	result["dependent_category"] = flattenCategoryRef(conflict.DependentCategory)
	result["dependee_vms_associated_categories"] = flattenCategoryRefList(conflict.DependeeVmsAssociatedCategories)
	result["dependent_vms_associated_categories"] = flattenCategoryRefList(conflict.DependentVmsAssociatedCategories)
	result["links"] = flattenVmStartupPolicyLinks(conflict.Links)
	return result
}
