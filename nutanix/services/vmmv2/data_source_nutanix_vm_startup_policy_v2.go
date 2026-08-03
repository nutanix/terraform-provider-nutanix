package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vmstartuppolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVmStartupPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixVmStartupPolicyV2Read,
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
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"categories": {
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
			"start_conditions": {
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
			"updated_by": {
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
			"num_dependency_conflicts": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_start_condition_conflicts": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
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
	}
}

func datasourceNutanixVmStartupPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id").(string)

	getRequest := import2.GetVmStartupPolicyByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.GetVmStartupPolicyById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching VM Startup Policy: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.VmStartupPolicy)

	if err := flattenVmStartupPolicy(d, &getResp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extID)
	return nil
}
