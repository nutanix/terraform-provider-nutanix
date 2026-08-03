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

func DatasourceNutanixVmStartupPolicyStartConditionConflictsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyStartConditionConflictsV2Read,
		Schema: map[string]*schema.Schema{
			"vm_startup_policy_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_condition_conflicts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: startConditionConflictSchema(),
				},
			},
		},
	}
}

func datasourceNutanixVmStartupPolicyStartConditionConflictsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)

	listRequest := import2.ListVmStartupPolicyStartConditionConflictsRequest{
		VmStartupPolicyExtId: utils.StringPtr(policyExtID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyStartConditionConflicts(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing VM Startup Policy Start Condition Conflicts: %v", err)
	}

	conflicts := resp.Data.GetValue().([]import1.StartConditionConflict)

	conflictList := make([]map[string]interface{}, len(conflicts))
	for i, c := range conflicts {
		conflictList[i] = flattenStartConditionConflictToMap(&c)
	}

	if err := d.Set("start_condition_conflicts", conflictList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
