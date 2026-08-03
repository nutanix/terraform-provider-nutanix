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

func DatasourceNutanixVmStartupPolicyVmComplianceStatesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyVmComplianceStatesV2Read,
		Schema: map[string]*schema.Schema{
			"vm_startup_policy_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_compliance_states": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associated_categories": {
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
						"cluster": {
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
						"compliance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"non_compliance_reason": {
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
					},
				},
			},
		},
	}
}

func datasourceNutanixVmStartupPolicyVmComplianceStatesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)

	listRequest := import2.ListVmStartupPolicyVmComplianceStatesRequest{
		VmStartupPolicyExtId: utils.StringPtr(policyExtID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyVmComplianceStates(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing VM Startup Policy VM Compliance States: %v", err)
	}

	states := resp.Data.GetValue().([]import1.VmStartupPolicyVmComplianceState)

	stateList := make([]map[string]interface{}, len(states))
	for i, s := range states {
		stateList[i] = flattenVmComplianceState(&s)
	}

	if err := d.Set("vm_compliance_states", stateList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

type vmComplianceInfo struct {
	Status string
	Reason string
}

var complianceStatusMap = map[string]vmComplianceInfo{
	"vmm.v4.ahv.policies.VmStartupPolicyCompliantVm":    {Status: "COMPLIANT"},
	"vmm.v4.ahv.policies.VmStartupPolicyNonCompliantVm": {Status: "NON_COMPLIANT"},
	"vmm.v4.ahv.policies.VmStartupPolicyPendingVm":      {Status: "PENDING"},
}

var nonComplianceReasonMap = map[string]string{
	"vmm.v4.ahv.policies.VmStartupPolicyNgtNotEnabled":       "NGT_NOT_ENABLED",
	"vmm.v4.ahv.policies.VmStartupPolicyHaNotSupported":      "HA_NOT_SUPPORTED",
	"vmm.v4.ahv.policies.VmStartupPolicyClusterNotSupported": "CLUSTER_NOT_SUPPORTED",
}

func flattenVmComplianceState(state *import1.VmStartupPolicyVmComplianceState) map[string]interface{} {
	result := map[string]interface{}{}
	if state.ExtId != nil {
		result["ext_id"] = utils.StringValue(state.ExtId)
	}
	if state.TenantId != nil {
		result["tenant_id"] = utils.StringValue(state.TenantId)
	}
	result["associated_categories"] = flattenCategoryRefList(state.AssociatedCategories)
	if state.Cluster != nil {
		clusterMap := map[string]interface{}{}
		if state.Cluster.ExtId != nil {
			clusterMap["ext_id"] = utils.StringValue(state.Cluster.ExtId)
		}
		result["cluster"] = []map[string]interface{}{clusterMap}
	}

	if state.ComplianceStatus != nil && state.ComplianceStatus.ObjectType_ != nil {
		if info, ok := complianceStatusMap[*state.ComplianceStatus.ObjectType_]; ok {
			result["compliance_status"] = info.Status

			if info.Status == "NON_COMPLIANT" {
				nonCompliantVm := state.ComplianceStatus.GetValue().(import1.VmStartupPolicyNonCompliantVm)
				if nonCompliantVm.NonComplianceReason != nil && nonCompliantVm.NonComplianceReason.ObjectType_ != nil {
					if reason, ok := nonComplianceReasonMap[*nonCompliantVm.NonComplianceReason.ObjectType_]; ok {
						result["non_compliance_reason"] = reason
					}
				}
			}
		}
	}

	result["links"] = flattenVmStartupPolicyLinks(state.Links)
	return result
}
