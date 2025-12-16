package networkingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	config "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const minItems = 2

func ResourceNutanixNetworkSecurityPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixNetworkSecurityPolicyV2Create,
		ReadContext:   ResourceNutanixNetworkSecurityPolicyV2Read,
		UpdateContext: ResourceNutanixNetworkSecurityPolicyV2Update,
		DeleteContext: ResourceNutanixNetworkSecurityPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"QUARANTINE", "ISOLATION", "APPLICATION", "SHAREDSERVICE"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"SAVE", "MONITOR", "ENFORCE"}, false),
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"QUARANTINE", "TWO_ENV_ISOLATION", "APPLICATION", "INTRA_GROUP", "MULTI_ENV_ISOLATION", "SHARED_SERVICE"}, false),
						},
						"spec": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"two_env_isolation_rule_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"first_isolation_group": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"second_isolation_group": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"application_rule_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"secured_group_category_associated_entity_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													Default:      "VM",
													ValidateFunc: validation.StringInSlice([]string{"SUBNET", "VM", "VPC"}, false),
												},
												"secured_group_category_references": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"secured_group_entity_group_reference": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"src_allow_spec": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													ValidateFunc: validation.StringInSlice([]string{"ALL", "NONE"}, false),
												},
												"dest_allow_spec": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													ValidateFunc: validation.StringInSlice([]string{"ALL", "NONE"}, false),
												},
												"src_category_associated_entity_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													Default:      "VM",
													ValidateFunc: validation.StringInSlice([]string{"SUBNET", "VM", "VPC"}, false),
												},
												"src_category_references": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"src_entity_group_reference": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"dest_category_associated_entity_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													Default:      "VM",
													ValidateFunc: validation.StringInSlice([]string{"SUBNET", "VM", "VPC"}, false),
												},
												"dest_category_references": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"dest_entity_group_reference": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"src_subnet": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"dest_subnet": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"src_address_group_references": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"dest_address_group_references": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"service_group_references": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"is_all_protocol_allowed": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"tcp_services": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_port": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"end_port": {
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"udp_services": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_port": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"end_port": {
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"icmp_services": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_all_allowed": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
															"type": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
															"code": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"network_function_chain_reference": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"network_function_reference": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"intra_entity_group_rule_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"secured_group_category_associated_entity_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													Default:      "VM",
													ValidateFunc: validation.StringInSlice([]string{"SUBNET", "VM", "VPC"}, false),
												},
												"secured_group_category_references": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"secured_group_entity_group_reference": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"secured_group_action": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
												},
												"secured_group_category_references": {
													Type:     schema.TypeList,
													Computed: true,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"secured_group_entity_group_reference": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"secured_group_service_references": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"tcp_services": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_port": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"end_port": {
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"udp_services": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_port": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"end_port": {
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"icmp_services": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_all_allowed": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
															"type": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
															"code": {
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
									"multi_env_isolation_rule_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"spec": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"all_to_all_isolation_group": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"isolation_group": {
																			Type:     schema.TypeList,
																			Required: true,
																			MinItems: minItems,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"group_category_associated_entity_type": {
																						Type:         schema.TypeString,
																						Optional:     true,
																						Computed:     true,
																						Default:      "VM",
																						ValidateFunc: validation.StringInSlice([]string{"SUBNET", "VM", "VPC"}, false),
																					},
																					"group_category_references": {
																						Type:     schema.TypeList,
																						Required: true,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"group_entity_group_reference": {
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
														},
													},
												},
											},
										},
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
					},
				},
			},
			"is_ipv6_traffic_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_hitlog_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_VLAN", "GLOBAL", "ALL_VPC", "VPC_LIST", "VPC_AS_CATEGORY"}, false),
			},
			"vpc_reference": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"secured_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
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

func ResourceNutanixNetworkSecurityPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	spec := import1.NetworkSecurityPolicy{}

	if name, ok := d.GetOk("name"); ok {
		spec.Name = utils.StringPtr(name.(string))
	}
	if types, ok := d.GetOk("type"); ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"QUARANTINE":  two,
			"ISOLATION":   three,
			"APPLICATION": four,
		}
		pInt := subMap[types.(string)]
		p := import1.SecurityPolicyType(pInt.(int))
		spec.Type = &p
	}
	if desc, ok := d.GetOk("description"); ok {
		spec.Description = utils.StringPtr(desc.(string))
	}
	if state, ok := d.GetOk("state"); ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"SAVE":    two,
			"MONITOR": three,
			"ENFORCE": four,
		}
		pInt := subMap[state.(string)]
		p := import1.SecurityPolicyState(pInt.(int))
		spec.State = &p
	}
	if rules, ok := d.GetOk("rules"); ok {
		spec.Rules = expandNetworkSecurityPolicyRule(rules.([]interface{}))
	}
	if isipv6, ok := d.GetOk("is_ipv6_traffic_allowed"); ok {
		spec.IsIpv6TrafficAllowed = utils.BoolPtr(isipv6.(bool))
	}
	if ishitlog, ok := d.GetOk("is_hitlog_enabled"); ok {
		spec.IsHitlogEnabled = utils.BoolPtr(ishitlog.(bool))
	}
	if scope, ok := d.GetOk("scope"); ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"ALL_VLAN": two,
			"ALL_VPC":  three,
			"VPC_LIST": four,
		}
		pInt := subMap[scope.(string)]
		p := import1.SecurityPolicyScope(pInt.(int))
		spec.Scope = &p
	}
	if vpcRef, ok := d.GetOk("vpc_reference"); ok {
		spec.VpcReferences = common.ExpandListOfString(vpcRef.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(spec, "", "  ")
	log.Printf("[DEBUG] Create Network Security Policy Payload: %s", string(aJSON))

	resp, err := conn.NetworkingSecurityInstance.CreateNetworkSecurityPolicy(&spec)
	if err != nil {
		return diag.Errorf("error while creating network security policy: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Network security  policy to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for network security policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching vpc UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

	return ResourceNutanixNetworkSecurityPolicyV2Read(ctx, d, meta)
}

func ResourceNutanixNetworkSecurityPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.NetworkingSecurityInstance.GetNetworkSecurityPolicyById(utils.StringPtr((d.Id())))
	if err != nil {
		return diag.Errorf("error while fetching network security policy: %v", err)
	}
	getResp := resp.Data.GetValue().(import1.NetworkSecurityPolicy)

	if err := d.Set("name", utils.StringValue(getResp.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", flattenSecurityPolicyType(getResp.Type)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", utils.StringValue(getResp.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", flattenPolicyState(getResp.State)); err != nil {
		return diag.FromErr(err)
	}

	// after creating role, operations saved in remote in different order than local
	if len(getResp.Rules) > 0 {
		// read the remote operations and local operations list
		remoteOperations := flattenNetworkSecurityPolicyRule(getResp.Rules)
		localOperations := expandNetworkSecurityPolicyRule(d.Get("rules").([]interface{}))

		// final result for checking if operations are different
		diff := false

		// convert local operations to string slice
		localOperationsStr := make([]string, len(localOperations))
		for i, v := range localOperations {
			localOperationsStr[i] = (flattenRuleType(v.Type))
		}

		log.Printf("[DEBUG] localOperationsStr: %v", localOperationsStr)

		// check if remote operations are different from local operations
		for _, operation := range remoteOperations {
			opsType := operation.(map[string]interface{})["type"]
			offset := indexOf(localOperationsStr, opsType.(string))

			if offset == -1 {
				log.Printf("[DEBUG] Rules %v not found in local rules", operation)
				diff = true
				break
			}
		}

		// if operations are different, update local operations
		if diff {
			log.Printf("[DEBUG] Rules are different. Updating local rules")
			if err := d.Set("rules", flattenNetworkSecurityPolicyRule(getResp.Rules)); err != nil {
				return diag.FromErr(err)
			}
		} else {
			// if operations are same, do not update local operations
			log.Printf("[DEBUG] Rules are same. Not updating local rules")
		}
	}

	if err := d.Set("is_ipv6_traffic_allowed", utils.BoolValue(getResp.IsIpv6TrafficAllowed)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_hitlog_enabled", utils.BoolValue(getResp.IsHitlogEnabled)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scope", flattenSecurityPolicyScope(getResp.Scope)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_reference", utils.StringSlice(getResp.VpcReferences)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secured_groups", utils.StringSlice(getResp.SecuredGroups)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreationTime != nil {
		t := getResp.CreationTime
		if err := d.Set("creation_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdateTime != nil {
		t := getResp.LastUpdateTime
		if err := d.Set("last_update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("is_system_defined", utils.BoolValue(getResp.IsSystemDefined)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", utils.StringValue(getResp.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", utils.StringValue(getResp.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(getResp.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksMicroSeg(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixNetworkSecurityPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	updatedSpec := import1.NetworkSecurityPolicy{}

	resp, err := conn.NetworkingSecurityInstance.GetNetworkSecurityPolicyById(utils.StringPtr((d.Id())))
	if err != nil {
		return diag.Errorf("error while fetching network security : %v", err)
	}

	updatedSpec = resp.Data.GetValue().(import1.NetworkSecurityPolicy)

	if d.HasChange("name") {
		updatedSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("type") {
		state := d.Get("type")
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"QUARANTINE":  two,
			"ISOLATION":   three,
			"APPLICATION": four,
		}
		pInt := subMap[state.(string)]
		p := import1.SecurityPolicyType(pInt.(int))
		updatedSpec.Type = &p
	}
	if d.HasChange("description") {
		updatedSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("rules") {
		updatedSpec.Rules = expandNetworkSecurityPolicyRule(d.Get("rules").([]interface{}))
	}
	if d.HasChange("state") {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"SAVE":    two,
			"MONITOR": three,
			"ENFORCE": four,
		}
		pInt := subMap[d.Get("state").(string)]
		p := import1.SecurityPolicyState(pInt.(int))
		updatedSpec.State = &p
	}
	if d.HasChange("is_ipv6_traffic_allowed") {
		updatedSpec.IsIpv6TrafficAllowed = utils.BoolPtr(d.Get("is_ipv6_traffic_allowed").(bool))
	}
	if d.HasChange("is_hitlog_enabled") {
		updatedSpec.IsHitlogEnabled = utils.BoolPtr(d.Get("is_hitlog_enabled").(bool))
	}
	if d.HasChange("scope") {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"ALL_VLAN": two,
			"ALL_VPC":  three,
			"VPC_LIST": four,
		}
		pInt := subMap[d.Get("scope").(string)]
		p := import1.SecurityPolicyScope(pInt.(int))
		updatedSpec.Scope = &p
	}
	if d.HasChange("vpc_reference") {
		updatedSpec.VpcReferences = common.ExpandListOfString(d.Get("vpc_reference").([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(updatedSpec, "", "  ")
	log.Printf("[DEBUG] Update Network Security Policy Payload: %s", string(aJSON))

	updatedResp, err := conn.NetworkingSecurityInstance.UpdateNetworkSecurityPolicyById(utils.StringPtr(d.Id()), &updatedSpec)
	if err != nil {
		return diag.Errorf("error while updating network security: %v", err)
	}

	TaskRef := updatedResp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Service Group to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for network security (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixNetworkSecurityPolicyV2Read(ctx, d, meta)
}

func ResourceNutanixNetworkSecurityPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.NetworkingSecurityInstance.DeleteNetworkSecurityPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting network security: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Service Group to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for network security (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandNetworkSecurityPolicyRule(pr []interface{}) []import1.NetworkSecurityPolicyRule {
	if len(pr) > 0 {
		nets := make([]import1.NetworkSecurityPolicyRule, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			net := import1.NetworkSecurityPolicyRule{}

			if desc, ok := val["description"]; ok {
				net.Description = utils.StringPtr(desc.(string))
			}
			if ty, ok := val["type"]; ok {
				const two, three, four, five, six = 2, 3, 4, 5, 6
				subMap := map[string]interface{}{
					"QUARANTINE":          two,
					"TWO_ENV_ISOLATION":   three,
					"APPLICATION":         four,
					"INTRA_GROUP":         five,
					"MULTI_ENV_ISOLATION": six,
				}
				pInt := subMap[ty.(string)]
				p := import1.RuleType(pInt.(int))
				net.Type = &p
			}
			if spec, ok := val["spec"]; ok {
				net.Spec = expandOneOfNetworkSecurityPolicyRuleSpec(spec)
			}
			nets[k] = net
		}
		return nets
	}
	return nil
}

func expandOneOfNetworkSecurityPolicyRuleSpec(pr interface{}) *import1.OneOfNetworkSecurityPolicyRuleSpec {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		policyRules := &import1.OneOfNetworkSecurityPolicyRuleSpec{}

		if isolation, ok := val["two_env_isolation_rule_spec"]; ok && len(isolation.([]interface{})) > 0 {
			iso := import1.NewTwoEnvIsolationRuleSpec()

			isoI := isolation.([]interface{})
			isoVal := isoI[0].(map[string]interface{})

			if firstIso, ok := isoVal["first_isolation_group"]; ok && len(firstIso.([]interface{})) > 0 {
				iso.FirstIsolationGroup = common.ExpandListOfString(firstIso.([]interface{}))
			}
			if secIso, ok := isoVal["second_isolation_group"]; ok && len(secIso.([]interface{})) > 0 {
				iso.SecondIsolationGroup = common.ExpandListOfString(secIso.([]interface{}))
			}
			policyRules.SetValue(*iso)
		}

		if appRule, ok := val["application_rule_spec"]; ok && len(appRule.([]interface{})) > 0 {
			app := import1.NewApplicationRuleSpec()

			appI := appRule.([]interface{})
			appVal := appI[0].(map[string]interface{})

			if secGroupCatAssocEntityType, ok := appVal["secured_group_category_associated_entity_type"]; ok && len(secGroupCatAssocEntityType.(string)) > 0 {
				const one, two, three = 1, 2, 3
				subMap := map[string]interface{}{
					"SUBNET": one,
					"VM":     two,
					"VPC":    three,
				}
				pInt := subMap[secGroupCatAssocEntityType.(string)]
				p := config.EntityType(pInt.(int))
				app.SecuredGroupCategoryAssociatedEntityType = &p
			}
			if secGroup, ok := appVal["secured_group_category_references"]; ok && len(secGroup.([]interface{})) > 0 {
				app.SecuredGroupCategoryReferences = common.ExpandListOfString(secGroup.([]interface{}))
			}
			if secGroupEntityGrpRef, ok := appVal["secured_group_entity_group_reference"]; ok && len(secGroupEntityGrpRef.(string)) > 0 {
				app.SecuredGroupEntityGroupReference = utils.StringPtr(secGroupEntityGrpRef.(string))
			}
			if srcAllow, ok := appVal["src_allow_spec"]; ok && len(srcAllow.(string)) > 0 {
				const two, three = 2, 3
				subMap := map[string]interface{}{
					"ALL":  two,
					"NONE": three,
				}
				pInt := subMap[srcAllow.(string)]
				p := import1.AllowType(pInt.(int))
				app.SrcAllowSpec = &p
			}
			if denyAllow, ok := appVal["dest_allow_spec"]; ok && len(denyAllow.(string)) > 0 {
				const two, three = 2, 3
				subMap := map[string]interface{}{
					"ALL":  two,
					"NONE": three,
				}
				pInt := subMap[denyAllow.(string)]
				p := import1.AllowType(pInt.(int))
				app.DestAllowSpec = &p
			}
			if srcCatAssocEntityType, ok := appVal["src_category_associated_entity_type"]; ok && len(srcCatAssocEntityType.(string)) > 0 {
				const one, two, three = 1, 2, 3
				subMap := map[string]interface{}{
					"SUBNET": one,
					"VM":     two,
					"VPC":    three,
				}
				pInt := subMap[srcCatAssocEntityType.(string)]
				p := config.EntityType(pInt.(int))
				app.SrcCategoryAssociatedEntityType = &p
			}
			if srcCatRef, ok := appVal["src_category_references"]; ok && len(srcCatRef.([]interface{})) > 0 {
				app.SrcCategoryReferences = common.ExpandListOfString(srcCatRef.([]interface{}))
			}
			if srcEntityGrpRef, ok := appVal["src_entity_group_reference"]; ok && len(srcEntityGrpRef.(string)) > 0 {
				app.SrcEntityGroupReference = utils.StringPtr(srcEntityGrpRef.(string))
			}
			if destCatAssocEntityType, ok := appVal["dest_category_associated_entity_type"]; ok && len(destCatAssocEntityType.(string)) > 0 {
				const one, two, three = 1, 2, 3
				subMap := map[string]interface{}{
					"SUBNET": one,
					"VM":     two,
					"VPC":    three,
				}
				pInt := subMap[destCatAssocEntityType.(string)]
				p := config.EntityType(pInt.(int))
				app.DestCategoryAssociatedEntityType = &p
			}
			if destCatRef, ok := appVal["dest_category_references"]; ok && len(destCatRef.([]interface{})) > 0 {
				app.DestCategoryReferences = common.ExpandListOfString(destCatRef.([]interface{}))
			}
			if destEntityGrpRef, ok := appVal["dest_entity_group_reference"]; ok && len(destEntityGrpRef.(string)) > 0 {
				app.DestEntityGroupReference = utils.StringPtr(destEntityGrpRef.(string))
			}
			if srcSubnet, ok := appVal["src_subnet"]; ok && len(srcSubnet.([]interface{})) > 0 {
				app.SrcSubnet = expandIPv4AddressMicroseg(srcSubnet)
			}
			if destSubnet, ok := appVal["dest_subnet"]; ok && len(destSubnet.([]interface{})) > 0 {
				app.DestSubnet = expandIPv4AddressMicroseg(destSubnet)
			}
			if srcAddGrpRef, ok := appVal["src_address_group_references"]; ok && len(srcAddGrpRef.([]interface{})) > 0 {
				app.SrcAddressGroupReferences = common.ExpandListOfString(srcAddGrpRef.([]interface{}))
			}
			if destAddGrpRef, ok := appVal["dest_address_group_references"]; ok && len(destAddGrpRef.([]interface{})) > 0 {
				app.DestAddressGroupReferences = common.ExpandListOfString(destAddGrpRef.([]interface{}))
			}
			if serviceGrpRef, ok := appVal["service_group_references"]; ok && len(serviceGrpRef.([]interface{})) > 0 {
				app.ServiceGroupReferences = common.ExpandListOfString(serviceGrpRef.([]interface{}))
			}
			if allProto, ok := appVal["is_all_protocol_allowed"]; ok {
				app.IsAllProtocolAllowed = utils.BoolPtr(allProto.(bool))
			}

			if tcp, ok := appVal["tcp_services"]; ok && len(tcp.([]interface{})) > 0 {
				app.TcpServices = expandTCPPortRangeSpec(tcp.([]interface{}))
			}
			if udp, ok := appVal["udp_services"]; ok && len(udp.([]interface{})) > 0 {
				app.UdpServices = expandUDPPortRangeSpec(udp.([]interface{}))
			}
			if icmp, ok := appVal["icmp_services"]; ok && len(icmp.([]interface{})) > 0 {
				app.IcmpServices = expandIcmpTypeCodeSpec(icmp.([]interface{}))
			}
			if netFuncChain, ok := appVal["network_function_chain_reference"]; ok && len(netFuncChain.(string)) > 0 {
				app.NetworkFunctionChainReference = utils.StringPtr(netFuncChain.(string))
			}
			if netFuncRef, ok := appVal["network_function_reference"]; ok && len(netFuncRef.(string)) > 0 {
				app.NetworkFunctionReference = utils.StringPtr(netFuncRef.(string))
			}
			policyRules.SetValue(*app)
		}

		if intraGroup, ok := val["intra_entity_group_rule_spec"]; ok && len(intraGroup.([]interface{})) > 0 {
			intra := import1.NewIntraEntityGroupRuleSpec()

			intraI := intraGroup.([]interface{})
			intraVal := intraI[0].(map[string]interface{})

			if secGroupCatAssocEntityType, ok := intraVal["secured_group_category_associated_entity_type"]; ok && len(secGroupCatAssocEntityType.(string)) > 0 {
				const one, two, three = 1, 2, 3
				subMap := map[string]interface{}{
					"SUBNET": one,
					"VM":     two,
					"VPC":    three,
				}
				pInt := subMap[secGroupCatAssocEntityType.(string)]
				p := config.EntityType(pInt.(int))
				intra.SecuredGroupCategoryAssociatedEntityType = &p
			}
			if secGroup, ok := intraVal["secured_group_category_references"]; ok && len(secGroup.([]interface{})) > 0 {
				intra.SecuredGroupCategoryReferences = common.ExpandListOfString(secGroup.([]interface{}))
			}
			if secGroupEntityGrpRef, ok := intraVal["secured_group_entity_group_reference"]; ok && len(secGroupEntityGrpRef.(string)) > 0 {
				intra.SecuredGroupEntityGroupReference = utils.StringPtr(secGroupEntityGrpRef.(string))
			}
			if secGroupAction, ok := intraVal["secured_group_action"]; ok && len(secGroupAction.(string)) > 0 {
				const two, three = 2, 3
				subMap := map[string]interface{}{
					"ALLOW": two,
					"DENY":  three,
				}
				pInt := subMap[secGroupAction.(string)]
				p := import1.IntraEntityGroupRuleAction(pInt.(int))
				intra.SecuredGroupAction = &p
			}
			if secGroupServiceRef, ok := intraVal["secured_group_service_references"]; ok && len(secGroupServiceRef.([]interface{})) > 0 {
				intra.SecuredGroupServiceReferences = common.ExpandListOfString(secGroupServiceRef.([]interface{}))
			}
			if tcp, ok := intraVal["tcp_services"]; ok && len(tcp.([]interface{})) > 0 {
				intra.TcpServices = expandTCPPortRangeSpec(tcp.([]interface{}))
			}
			if udp, ok := intraVal["udp_services"]; ok && len(udp.([]interface{})) > 0 {
				intra.UdpServices = expandUDPPortRangeSpec(udp.([]interface{}))
			}
			if icmp, ok := intraVal["icmp_services"]; ok && len(icmp.([]interface{})) > 0 {
				intra.IcmpServices = expandIcmpTypeCodeSpec(icmp.([]interface{}))
			}
			policyRules.SetValue(*intra)
		}

		if multiEnv, ok := val["multi_env_isolation_rule_spec"]; ok && len(multiEnv.([]interface{})) > 0 {
			multi := import1.NewMultiEnvIsolationRuleSpec()

			multiI := multiEnv.([]interface{})
			multiVal := multiI[0].(map[string]interface{})

			if spec, ok := multiVal["spec"]; ok && len(spec.([]interface{})) > 0 {
				oneOfMultiEnvIsolationRuleSpecSpec := expandOneOfMultiEnvIsolationRuleSpecSpec(spec.([]interface{}))

				aJSON, _ := json.Marshal(oneOfMultiEnvIsolationRuleSpecSpec)
				log.Printf("[DEBUG] OneOfMultiEnvIsolationRuleSpecSpec: %s", aJSON)

				multi.Spec = oneOfMultiEnvIsolationRuleSpecSpec

				aJSON, _ = json.Marshal(multi)
				log.Printf("[DEBUG] MultiEnvIsolationRuleSpec: %s", aJSON)
			}

			err := policyRules.SetValue(*multi)
			if err != nil {
				log.Printf("[ERROR] Error while setting value for multi env isolation rule: %v", err)
				return nil
			}
		}
		aJSON, _ := json.Marshal(policyRules)
		log.Printf("[DEBUG] PolicyRules: %s", aJSON)
		return policyRules
	}
	return nil
}

func expandOneOfMultiEnvIsolationRuleSpecSpec(spec []interface{}) *import1.OneOfMultiEnvIsolationRuleSpecSpec {
	if len(spec) > 0 {
		specVal := spec[0].(map[string]interface{})
		allToAllIso := import1.NewAllToAllIsolationGroup()

		oneOfMultiEnv := import1.NewOneOfMultiEnvIsolationRuleSpecSpec()

		if allToAll, ok := specVal["all_to_all_isolation_group"]; ok && len(allToAll.([]interface{})) > 0 {
			allToAllI := allToAll.([]interface{})
			allToAllVal := allToAllI[0].(map[string]interface{})

			if isoGroup, ok := allToAllVal["isolation_group"]; ok && len(isoGroup.([]interface{})) > 0 {
				allToAllIso.IsolationGroups = expandIsolationGroup(isoGroup.([]interface{}))
			}

			aJSON, _ := json.Marshal(allToAllIso)
			log.Printf("[DEBUG] AllToAllIsolationGroup: %s", aJSON)

			err := oneOfMultiEnv.SetValue(*allToAllIso)
			if err != nil {
				log.Printf("[ERROR] Error while setting value for MultiEnvIsolationRuleSpec.spec.allToAllIsolationGroup: %v", err)
				return nil
			}
			bJSON, _ := json.Marshal(oneOfMultiEnv)
			log.Printf("[DEBUG] OneOfMultiEnvIsolationRuleSpecSpec: %s", bJSON)
			return oneOfMultiEnv
		}
	}
	return nil
}

func expandIsolationGroup(isolationGroup []interface{}) []import1.IsolationGroup {
	if len(isolationGroup) > 0 {
		isolations := make([]import1.IsolationGroup, len(isolationGroup))

		for k, v := range isolationGroup {
			val := v.(map[string]interface{})
			iso := import1.IsolationGroup{}

			if groupCatAssocEntityType, ok := val["group_category_associated_entity_type"]; ok && len(groupCatAssocEntityType.(string)) > 0 {
				const one, two, three = 1, 2, 3
				subMap := map[string]interface{}{
					"SUBNET": one,
					"VM":     two,
					"VPC":    three,
				}
				pInt := subMap[groupCatAssocEntityType.(string)]
				p := config.EntityType(pInt.(int))
				iso.GroupCategoryAssociatedEntityType = &p
			}
			if groupCat, ok := val["group_category_references"]; ok && len(groupCat.([]interface{})) > 0 {
				iso.GroupCategoryReferences = common.ExpandListOfString(groupCat.([]interface{}))
			}
			if groupEntityGrpRef, ok := val["group_entity_group_reference"]; ok && len(groupEntityGrpRef.(string)) > 0 {
				iso.GroupEntityGroupReference = utils.StringPtr(groupEntityGrpRef.(string))
			}
			isolations[k] = iso
		}

		aJSON, _ := json.Marshal(isolations)
		log.Printf("[DEBUG] IsolationGroups: %s", aJSON)

		return isolations
	}
	return nil
}

func expandIPv4AddressMicroseg(pr interface{}) *config.IPv4Address {
	if pr == nil {
		return nil
	}

	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	ipv4 := &config.IPv4Address{}

	if v, ok := valMap["value"]; ok {
		if s, ok2 := v.(string); ok2 {
			ipv4.Value = utils.StringPtr(s)
		}
	}

	if p, ok := valMap["prefix_length"]; ok {
		if n, ok2 := p.(int); ok2 {
			ipv4.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv4
}

func indexOf(slice []string, target string) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}
