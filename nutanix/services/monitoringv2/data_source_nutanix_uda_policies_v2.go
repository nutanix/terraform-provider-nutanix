package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixUdaPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixUdaPoliciesV2Read,
		Schema: map[string]*schema.Schema{
			"uda_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A globally unique identifier of an instance that is suitable for external consumption.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A globally unique identifier that represents the tenant that owns this entity.",
						},
						"links": schemaForLinks(),
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Title of the policy.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the policy.",
						},
						"entity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity type associated with the User-Defined Alert policy. Allowed values are VM, node and cluster.",
						},
						"trigger_conditions": schemaForTriggerConditionsComputed(),
						"filters": schemaForFiltersComputed(),
						"impact_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Impact types for the associated resulting alert.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_auto_resolved": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the auto-resolve feature is enabled for this policy.",
						},
						"is_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable/Disable flag for the policy.",
						},
						"trigger_wait_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Waiting duration in seconds before triggering the alert, when the specified condition is met.",
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Username of the user who created the policy.",
						},
						"last_updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last updated time of the policy in ISO 8601 format.",
						},
						"policies_to_override": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of IDs of the alert policies that should be overridden.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"related_policies": schemaForRelatedPoliciesComputed(),
						"is_expected_to_error_on_conflict": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Error when conflicting alert policies are found.",
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixUdaPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	resp, err := conn.UserDefinedPolicies.ListUdaPolicies(nil, nil, nil, nil, nil)
	if err != nil {
		return diag.Errorf("error while listing User-Defined Alert policies: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("uda_policies", []map[string]interface{}{}); err != nil {
			return diag.Errorf("error setting User-Defined Alert policies: %s", err)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	getResp := resp.Data.GetValue().([]serviceability.UserDefinedPolicy)

	if err := d.Set("uda_policies", flattenUdaPolicies(getResp)); err != nil {
		return diag.Errorf("error setting User-Defined Alert policies: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}
