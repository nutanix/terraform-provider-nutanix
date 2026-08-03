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

func DatasourceNutanixVmStartupPolicyDependencyConflictsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyDependencyConflictsV2Read,
		Schema: map[string]*schema.Schema{
			"vm_startup_policy_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dependency_conflicts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dependencyConflictSchema(),
				},
			},
		},
	}
}

func datasourceNutanixVmStartupPolicyDependencyConflictsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	policyExtID := d.Get("vm_startup_policy_ext_id").(string)

	listRequest := import2.ListVmStartupPolicyDependencyConflictsRequest{
		VmStartupPolicyExtId: utils.StringPtr(policyExtID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.ListVmStartupPolicyDependencyConflicts(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing VM Startup Policy Dependency Conflicts: %v", err)
	}

	conflicts := resp.Data.GetValue().([]import1.DependencyConflict)

	conflictList := make([]map[string]interface{}, len(conflicts))
	for i, c := range conflicts {
		conflictList[i] = flattenDependencyConflictToMap(&c)
	}

	if err := d.Set("dependency_conflicts", conflictList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
