package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixUdaPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixUdaPolicyV2Create,
		ReadContext:   ResourceNutanixUdaPolicyV2Read,
		UpdateContext: ResourceNutanixUdaPolicyV2Update,
		DeleteContext: ResourceNutanixUdaPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Title of the policy.",
			},
			"entity_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Entity type associated with the User-Defined Alert policy. Allowed values are VM, node and cluster.",
			},
			"trigger_conditions": schemaForTriggerConditionsInput(),
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of the policy.",
			},
			"filters": schemaForFiltersInput(),
			"impact_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Impact types for the associated resulting alert.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_auto_resolved": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the auto-resolve feature is enabled for this policy.",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable/Disable flag for the policy.",
			},
			"is_expected_to_error_on_conflict": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Error when conflicting alert policies are found.",
			},
			"trigger_wait_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Waiting duration in seconds before triggering the alert, when the specified condition is met. It is set to 600s by default.",
			},
			"policies_to_override": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of IDs of the alert policies that should be overridden.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
			"links": schemaForLinks(),
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
			"related_policies": schemaForRelatedPoliciesComputed(),
		},
	}
}

func ResourceNutanixUdaPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	body := serviceability.NewUserDefinedPolicy()

	if title, ok := d.GetOk("title"); ok {
		body.Title = utils.StringPtr(title.(string))
	}
	if entityType, ok := d.GetOk("entity_type"); ok {
		body.EntityType = utils.StringPtr(entityType.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if tc, ok := d.GetOk("trigger_conditions"); ok {
		body.TriggerConditions = expandTriggerConditions(tc.([]interface{}))
	}
	if f, ok := d.GetOk("filters"); ok {
		body.Filters = expandFilters(f.([]interface{}))
	}
	if it, ok := d.GetOk("impact_types"); ok {
		body.ImpactTypes = expandImpactTypes(it.([]interface{}))
	}
	if iar, ok := d.GetOk("is_auto_resolved"); ok {
		body.IsAutoResolved = utils.BoolPtr(iar.(bool))
	}
	if ie, ok := d.GetOk("is_enabled"); ok {
		body.IsEnabled = utils.BoolPtr(ie.(bool))
	}
	if ieoc, ok := d.GetOk("is_expected_to_error_on_conflict"); ok {
		body.IsExpectedToErrorOnConflict = utils.BoolPtr(ieoc.(bool))
	}
	if twp, ok := d.GetOk("trigger_wait_period"); ok {
		body.TriggerWaitPeriod = utils.Int64Ptr(int64(twp.(int)))
	}
	if pto, ok := d.GetOk("policies_to_override"); ok {
		ptoList := pto.([]interface{})
		ptoStrings := make([]string, len(ptoList))
		for i, v := range ptoList {
			ptoStrings[i] = v.(string)
		}
		body.PoliciesToOverride = ptoStrings
	}

	resp, err := conn.UserDefinedPolicies.CreateUdaPolicy(body)
	if err != nil {
		return diag.Errorf("error while creating User-Defined Alert policy: %v", err)
	}

	getResp := resp.Data.GetValue().(serviceability.UserDefinedPolicy)
	d.SetId(utils.StringValue(getResp.ExtId))

	return ResourceNutanixUdaPolicyV2Read(ctx, d, meta)
}

func ResourceNutanixUdaPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Id()

	resp, err := conn.UserDefinedPolicies.GetUdaPolicyById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching User-Defined Alert policy: %s", err)
	}

	getResp := resp.Data.GetValue().(serviceability.UserDefinedPolicy)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
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
	if err := d.Set("is_expected_to_error_on_conflict", getResp.IsExpectedToErrorOnConflict); err != nil {
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

	return nil
}

func ResourceNutanixUdaPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Id()

	readResp, err := conn.UserDefinedPolicies.GetUdaPolicyById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching User-Defined Alert policy for update: %s", err)
	}

	getResp := readResp.Data.GetValue().(serviceability.UserDefinedPolicy)
	body := &getResp

	if d.HasChange("title") {
		body.Title = utils.StringPtr(d.Get("title").(string))
	}
	if d.HasChange("entity_type") {
		body.EntityType = utils.StringPtr(d.Get("entity_type").(string))
	}
	if d.HasChange("description") {
		body.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("trigger_conditions") {
		body.TriggerConditions = expandTriggerConditions(d.Get("trigger_conditions").([]interface{}))
	}
	if d.HasChange("filters") {
		body.Filters = expandFilters(d.Get("filters").([]interface{}))
	}
	if d.HasChange("impact_types") {
		body.ImpactTypes = expandImpactTypes(d.Get("impact_types").([]interface{}))
	}
	if d.HasChange("is_auto_resolved") {
		body.IsAutoResolved = utils.BoolPtr(d.Get("is_auto_resolved").(bool))
	}
	if d.HasChange("is_enabled") {
		body.IsEnabled = utils.BoolPtr(d.Get("is_enabled").(bool))
	}
	if d.HasChange("is_expected_to_error_on_conflict") {
		body.IsExpectedToErrorOnConflict = utils.BoolPtr(d.Get("is_expected_to_error_on_conflict").(bool))
	}
	if d.HasChange("trigger_wait_period") {
		body.TriggerWaitPeriod = utils.Int64Ptr(int64(d.Get("trigger_wait_period").(int)))
	}
	if d.HasChange("policies_to_override") {
		ptoList := d.Get("policies_to_override").([]interface{})
		ptoStrings := make([]string, len(ptoList))
		for i, v := range ptoList {
			ptoStrings[i] = v.(string)
		}
		body.PoliciesToOverride = ptoStrings
	}

	_, err = conn.UserDefinedPolicies.UpdateUdaPolicyById(utils.StringPtr(extID), body)
	if err != nil {
		return diag.Errorf("error while updating User-Defined Alert policy: %v", err)
	}

	return ResourceNutanixUdaPolicyV2Read(ctx, d, meta)
}

func ResourceNutanixUdaPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Id()

	_, err := conn.UserDefinedPolicies.DeleteUdaPolicyById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while deleting User-Defined Alert policy: %v", err)
	}

	d.SetId("")
	return nil
}
