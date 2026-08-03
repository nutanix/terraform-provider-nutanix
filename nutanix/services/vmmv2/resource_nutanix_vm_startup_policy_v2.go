package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/common/v1/response"
	vmmConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/prism/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/vmstartuppolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmStartupPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixVmStartupPolicyV2Create,
		ReadContext:   resourceNutanixVmStartupPolicyV2Read,
		UpdateContext: resourceNutanixVmStartupPolicyV2Update,
		DeleteContext: resourceNutanixVmStartupPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 2,
				MaxItems: 6,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"categories": {
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 5,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Optional: true,
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
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delay_duration_secs": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 600),
						},
						"power_state_criteria": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"power_on": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"guest_bootup": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timeout_duration_secs": {
													Type:     schema.TypeInt,
													Optional: true,
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
				Optional: true,
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

func resourceNutanixVmStartupPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	body := import1.VmStartupPolicy{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if groups, ok := d.GetOk("groups"); ok {
		body.Groups = expandVmStartupPolicyGroups(groups.([]interface{}))
	}
	if sc, ok := d.GetOk("start_conditions"); ok {
		body.StartConditions = expandVmStartupPolicyStartConditions(sc.([]interface{}))
	}
	if projectExtID, ok := d.GetOk("project_ext_id"); ok {
		body.ProjectExtId = utils.StringPtr(projectExtID.(string))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] VM Startup Policy Request Body: %s", string(aJSON))

	createRequest := import2.CreateVmStartupPolicyRequest{
		Body: &body,
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.CreateVmStartupPolicy(ctx, &createRequest)
	if err != nil {
		return diag.Errorf("error while creating VM Startup Policy: %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Startup Policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskRequest := import3.GetTaskByIdRequest{
		ExtId: taskUUID,
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskRequest)
	if err != nil {
		return diag.Errorf("error while fetching VM Startup Policy task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] VM Startup Policy Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeVMStartupPolicy, "VM Startup Policy")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return resourceNutanixVmStartupPolicyV2Read(ctx, d, meta)
}

func resourceNutanixVmStartupPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	getRequest := import2.GetVmStartupPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.GetVmStartupPolicyById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching VM Startup Policy: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.VmStartupPolicy)

	if err := flattenVmStartupPolicy(d, &getResp); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixVmStartupPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("project_ext_id") {
		return diag.Errorf("error while updating project_ext_id: Update of project_ext_id is not supported")
	}
	conn := meta.(*conns.Client).VmmAPI

	getRequest := import2.GetVmStartupPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.GetVmStartupPolicyById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching VM Startup Policy for update: %v", err)
	}

	respPolicy := resp.Data.GetValue().(import1.VmStartupPolicy)
	updateSpec := respPolicy

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("groups") {
		updateSpec.Groups = expandVmStartupPolicyGroups(d.Get("groups").([]interface{}))
	}
	if d.HasChange("start_conditions") {
		updateSpec.StartConditions = expandVmStartupPolicyStartConditions(d.Get("start_conditions").([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", " ")
	log.Printf("[DEBUG] VM Startup Policy Update Request Body: %s", string(aJSON))

	updateRequest := import2.UpdateVmStartupPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  &updateSpec,
	}

	args := make(map[string]interface{})
	etagValue := conn.VmStartupPoliciesAPIInstance.ApiClient.GetEtag(resp)
	args["If-Match"] = utils.StringPtr(etagValue)

	updateResp, err := conn.VmStartupPoliciesAPIInstance.UpdateVmStartupPolicyById(ctx, &updateRequest, args)
	if err != nil {
		return diag.Errorf("error while updating VM Startup Policy: %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Startup Policy (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return resourceNutanixVmStartupPolicyV2Read(ctx, d, meta)
}

func resourceNutanixVmStartupPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	getRequest := import2.GetVmStartupPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	readResp, err := conn.VmStartupPoliciesAPIInstance.GetVmStartupPolicyById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while reading VM Startup Policy for delete: %v", err)
	}

	args := make(map[string]interface{})
	etagValue := conn.VmStartupPoliciesAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etagValue)

	deleteRequest := import2.DeleteVmStartupPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	resp, err := conn.VmStartupPoliciesAPIInstance.DeleteVmStartupPolicyById(ctx, &deleteRequest, args)
	if err != nil {
		return diag.Errorf("error while deleting VM Startup Policy: %v", err)
	}
	TaskRef := resp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Startup Policy (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandVmStartupPolicyGroups(groups []interface{}) []import1.DependencyGroup {
	if len(groups) == 0 {
		return nil
	}
	result := make([]import1.DependencyGroup, len(groups))
	for i, g := range groups {
		gMap := g.(map[string]interface{})
		dg := import1.DependencyGroup{}
		if cats, ok := gMap["categories"]; ok {
			catList := cats.([]interface{})
			catRefs := make([]import1.CategoryReference, len(catList))
			for j, c := range catList {
				cMap := c.(map[string]interface{})
				if extID, ok := cMap["ext_id"]; ok && extID.(string) != "" {
					catRefs[j] = import1.CategoryReference{ExtId: utils.StringPtr(extID.(string))}
				}
			}
			dg.Categories = catRefs
		}
		result[i] = dg
	}
	return result
}

func expandVmStartupPolicyStartConditions(sc []interface{}) []import1.StartCondition {
	if len(sc) == 0 {
		return nil
	}
	result := make([]import1.StartCondition, len(sc))
	for i, s := range sc {
		sMap := s.(map[string]interface{})
		cond := import1.StartCondition{}
		if delay, ok := sMap["delay_duration_secs"]; ok {
			cond.DelayDurationSecs = utils.IntPtr(delay.(int))
		}
		if psc, ok := sMap["power_state_criteria"]; ok {
			pscList := psc.([]interface{})
			if len(pscList) > 0 && pscList[0] != nil {
				pscMap := pscList[0].(map[string]interface{})
				oneOf := import1.NewOneOfStartConditionPowerStateCriteria()
				set := false
				if gbRaw, gbOk := pscMap["guest_bootup"]; gbOk {
					gbList := gbRaw.([]interface{})
					if len(gbList) > 0 && gbList[0] != nil {
						gbMap := gbList[0].(map[string]interface{})
						gb := import1.NewPowerStateCriteriaGuestBootup()
						if timeout, ok := gbMap["timeout_duration_secs"]; ok {
							gb.TimeoutDurationSecs = utils.IntPtr(timeout.(int))
						}
						oneOf.SetValue(*gb)
						cond.PowerStateCriteria = oneOf
						set = true
					}
				}
				if !set {
					if poRaw, poOk := pscMap["power_on"]; poOk {
						poList := poRaw.([]interface{})
						if len(poList) > 0 {
							po := import1.NewPowerStateCriteriaPowerOn()
							oneOf.SetValue(*po)
							cond.PowerStateCriteria = oneOf
						}
					}
				}
			}
		}
		result[i] = cond
	}
	return result
}

func flattenVmStartupPolicy(d *schema.ResourceData, policy *import1.VmStartupPolicy) error {
	if policy.ExtId != nil {
		if err := d.Set("ext_id", utils.StringValue(policy.ExtId)); err != nil {
			return err
		}
	}
	if policy.Name != nil {
		if err := d.Set("name", utils.StringValue(policy.Name)); err != nil {
			return err
		}
	}
	if policy.Description != nil {
		if err := d.Set("description", utils.StringValue(policy.Description)); err != nil {
			return err
		}
	}
	if policy.CreateTime != nil {
		if err := d.Set("create_time", utils.TimeStringValue(policy.CreateTime)); err != nil {
			return err
		}
	}
	if policy.UpdateTime != nil {
		if err := d.Set("update_time", utils.TimeStringValue(policy.UpdateTime)); err != nil {
			return err
		}
	}
	if policy.ProjectExtId != nil {
		if err := d.Set("project_ext_id", utils.StringValue(policy.ProjectExtId)); err != nil {
			return err
		}
	}
	if policy.TenantId != nil {
		if err := d.Set("tenant_id", utils.StringValue(policy.TenantId)); err != nil {
			return err
		}
	}
	if policy.NumCompliantVms != nil {
		if err := d.Set("num_compliant_vms", int(*policy.NumCompliantVms)); err != nil {
			return err
		}
	}
	if policy.NumNonCompliantVms != nil {
		if err := d.Set("num_non_compliant_vms", int(*policy.NumNonCompliantVms)); err != nil {
			return err
		}
	}
	if policy.NumPendingVms != nil {
		if err := d.Set("num_pending_vms", int(*policy.NumPendingVms)); err != nil {
			return err
		}
	}
	if policy.NumDependencyConflicts != nil {
		if err := d.Set("num_dependency_conflicts", int(*policy.NumDependencyConflicts)); err != nil {
			return err
		}
	}
	if policy.NumStartConditionConflicts != nil {
		if err := d.Set("num_start_condition_conflicts", int(*policy.NumStartConditionConflicts)); err != nil {
			return err
		}
	}
	if err := d.Set("created_by", flattenUserReference(policy.CreatedBy)); err != nil {
		return err
	}
	if err := d.Set("updated_by", flattenUserReference(policy.UpdatedBy)); err != nil {
		return err
	}
	if err := d.Set("groups", flattenVmStartupPolicyGroups(policy.Groups)); err != nil {
		return err
	}
	if err := d.Set("start_conditions", flattenVmStartupPolicyStartConditions(policy.StartConditions)); err != nil {
		return err
	}
	if err := d.Set("links", flattenVmStartupPolicyLinks(policy.Links)); err != nil {
		return err
	}
	return nil
}

func flattenUserReference(ref *import1.UserReference) []map[string]interface{} {
	if ref == nil {
		return nil
	}
	result := map[string]interface{}{}
	if ref.ExtId != nil {
		result["ext_id"] = utils.StringValue(ref.ExtId)
	}
	return []map[string]interface{}{result}
}

func flattenVmStartupPolicyGroups(groups []import1.DependencyGroup) []map[string]interface{} {
	if len(groups) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(groups))
	for i, g := range groups {
		gMap := map[string]interface{}{}
		if g.Categories != nil {
			cats := make([]map[string]interface{}, len(g.Categories))
			for j, c := range g.Categories {
				catMap := map[string]interface{}{}
				if c.ExtId != nil {
					catMap["ext_id"] = utils.StringValue(c.ExtId)
				}
				cats[j] = catMap
			}
			gMap["categories"] = cats
		}
		result[i] = gMap
	}
	return result
}

func flattenVmStartupPolicyStartConditions(sc []import1.StartCondition) []map[string]interface{} {
	if len(sc) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(sc))
	for i, s := range sc {
		sMap := map[string]interface{}{}
		if s.DelayDurationSecs != nil {
			sMap["delay_duration_secs"] = *s.DelayDurationSecs
		}
		if s.PowerStateCriteria != nil && s.PowerStateCriteria.ObjectType_ != nil {
			pscMap := map[string]interface{}{}
			val := s.PowerStateCriteria.GetValue()
			if val != nil {
				switch *s.PowerStateCriteria.ObjectType_ {
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
		result[i] = sMap
	}
	return result
}

func flattenVmStartupPolicyLinks(links []response.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(links))
	for i, l := range links {
		lMap := map[string]interface{}{}
		if l.Href != nil {
			lMap["href"] = utils.StringValue(l.Href)
		}
		if l.Rel != nil {
			lMap["rel"] = utils.StringValue(l.Rel)
		}
		result[i] = lMap
	}
	return result
}
