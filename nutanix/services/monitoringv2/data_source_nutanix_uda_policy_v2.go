package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixUdaPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixUdaPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
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
	}
}

func DatasourceNutanixUdaPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.UserDefinedPolicies.GetUdaPolicyById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching User-Defined Alert policy: %s", err)
	}

	getResp := resp.Data.GetValue().(serviceability.UserDefinedPolicy)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("title", getResp.Title); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_type", getResp.EntityType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("trigger_conditions", flattenTriggerConditions(getResp.TriggerConditions)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("filters", flattenFilters(getResp.Filters)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("impact_types", flattenImpactTypes(getResp.ImpactTypes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_auto_resolved", getResp.IsAutoResolved); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", getResp.IsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("trigger_wait_period", getResp.TriggerWaitPeriod); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if getResp.LastUpdatedTime != nil {
		if err := d.Set("last_updated_time", getResp.LastUpdatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("policies_to_override", getResp.PoliciesToOverride); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("related_policies", flattenRelatedPolicies(getResp.RelatedPolicies)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_expected_to_error_on_conflict", getResp.IsExpectedToErrorOnConflict); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extID)
	return nil
}
