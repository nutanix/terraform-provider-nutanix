package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixFindConflictingUdaPoliciesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixFindConflictingUdaPoliciesV2Create,
		ReadContext:   ResourceNutanixFindConflictingUdaPoliciesV2Read,
		DeleteContext: ResourceNutanixFindConflictingUdaPoliciesV2Delete,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Title of the policy.",
			},
			"entity_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Entity type associated with the User-Defined Alert policy. Allowed values are VM, node and cluster.",
			},
			"trigger_conditions": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Trigger conditions for the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},
									"threshold_value": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"int_value": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"double_value": {
													Type:     schema.TypeFloat,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"condition_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity_level": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the policy.",
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_filter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"group_filter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"impact_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_auto_resolved": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"trigger_wait_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"conflicting_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of conflicting policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique UUID associated with the User-Defined Alert policy, that conflicts with the given policy.",
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixFindConflictingUdaPoliciesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if twp, ok := d.GetOk("trigger_wait_period"); ok {
		body.TriggerWaitPeriod = utils.Int64Ptr(int64(twp.(int)))
	}

	resp, err := conn.UserDefinedPolicies.FindConflictingUdaPolicies(body)
	if err != nil {
		return diag.Errorf("error while finding conflicting User-Defined Alert policies: %v", err)
	}

	if resp.Data != nil {
		conflicts := resp.Data.GetValue().([]serviceability.ConflictingPolicy)
		conflictList := make([]map[string]interface{}, len(conflicts))
		for i, cp := range conflicts {
			conflictList[i] = map[string]interface{}{
				"ext_id": utils.StringValue(cp.ExtId),
			}
		}
		if err := d.Set("conflicting_policies", conflictList); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("conflicting_policies", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(utils.GenUUID())
	return nil
}

func ResourceNutanixFindConflictingUdaPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixFindConflictingUdaPoliciesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
