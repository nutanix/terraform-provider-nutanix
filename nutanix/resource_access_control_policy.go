package nutanix

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/spf13/cast"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
	"strings"
)

func resourceNutanixAccessControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixAccessControlPolicyCreate,
		Read:   resourceNutanixAccessControlPolicyRead,
		Update: resourceNutanixAccessControlPolicyUpdate,
		Delete: resourceNutanixAccessControlPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_reference_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"user_group_reference_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"role_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"role"}, false),
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"filter_list": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"context_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scope_filter_expression_list": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"left_hand_side": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"CATEGORY", "PROJECT"}, false),
												},
												"operator": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"IN", "IN_ALL", "NOT_IN"}, false),
												},
												"right_hand_side": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"collection": {
																Type:         schema.TypeString,
																Optional:     true,
																Computed:     true,
																ValidateFunc: validation.StringInSlice([]string{"ALL"}, false),
															},
															"categories": categoriesSchema(),
															"uuid_list": {
																Type:     schema.TypeSet,
																Optional: true,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
									"entity_filter_expression_list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"left_hand_side_entity_type": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"operator": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"IN", "NOT_IN"}, false),
												},
												"right_hand_side": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"collection": {
																Type:         schema.TypeString,
																Optional:     true,
																Computed:     true,
																ValidateFunc: validation.StringInSlice([]string{"ALL", "SELF_OWNED"}, false),
															},
															"categories": categoriesSchema(),
															"uuid_list": {
																Type:     schema.TypeSet,
																Optional: true,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
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
	}
}

func resourceNutanixAccessControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.AccessControlPolicy{}
	spec := &v3.AccessControlPolicySpec{}
	metadata := &v3.Metadata{}
	access := &v3.AccessControlPolicyResources{}

	rf, rfOk := d.GetOk("role_reference")

	if !rfOk {
		return fmt.Errorf("please provide the required `role_reference` attribute")
	}

	if err := getMetadataAttributes(d, metadata, "access_control_policy"); err != nil {
		return err
	}

	access.RoleReference = validateRefList(rf.([]interface{}))

	expandAccessControlPolicyResources(d, access)

	if description, ok := d.GetOk("description"); ok {
		spec.Description = utils.StringPtr(description.(string))
	}

	if name, ok := d.GetOk("name"); ok {
		spec.Name = utils.StringPtr(name.(string))
	}
	spec.Resources = access
	request.Metadata = metadata
	request.Spec = spec

	resp, err := conn.V3.CreateAccessControlPolicy(request)
	if err != nil {
		return fmt.Errorf("error creating Nutanix AccessControlPolicy %s: %+v", utils.StringValue(spec.Name), err)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the AccessControlPolicy to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING", "PENDING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		id := d.Id()
		d.SetId("")
		return fmt.Errorf("error waiting for access control policy id (%s) to create: %+v", id, err)
	}

	// Setting Description because in Get request is not present.
	d.Set("description", utils.StringValue(resp.Spec.Description))

	d.SetId(utils.StringValue(resp.Metadata.UUID))

	return resourceNutanixAccessControlPolicyRead(d, meta)
}

func resourceNutanixAccessControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	id := d.Id()
	resp, err := conn.V3.GetAccessControlPolicy(id)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		errDel := resourceNutanixSubnetDelete(d, meta)
		if errDel != nil {
			return fmt.Errorf("error deleting access control policy (%s) after read error: %+v", id, errDel)
		}
		d.SetId("")
		return fmt.Errorf("error reading access control policy id (%s): %+v", id, err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", flattenReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	d.Set("api_version", resp.APIVersion)

	if status := resp.Status; status != nil {
		if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
			return err
		}
		if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
			return err
		}
		if err := d.Set("state", utils.StringValue(resp.Status.State)); err != nil {
			return err
		}

		if res := status.Resources; res != nil {
			if err := d.Set("user_reference_list", flattenArrayReferenceValues(status.Resources.UserReferenceList)); err != nil {
				return err
			}
			if err := d.Set("user_group_reference_list", flattenArrayReferenceValues(status.Resources.UserGroupReferenceList)); err != nil {
				return err
			}
			if err := d.Set("role_reference", flattenReferenceValuesList(status.Resources.RoleReference)); err != nil {
				return err
			}
			if status.Resources.FilterList.ContextList != nil {
				if err := d.Set("filter_list", flattenFilterList(status.Resources.FilterList)); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func resourceNutanixAccessControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	request := &v3.AccessControlPolicy{}
	metadata := &v3.Metadata{}
	res := &v3.AccessControlPolicyResources{}
	spec := &v3.AccessControlPolicySpec{}

	id := d.Id()
	response, err := conn.V3.GetAccessControlPolicy(id)

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return fmt.Errorf("error retrieving for access control policy id (%s) :%+v", id, err)
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			res = response.Spec.Resources
		}
	}

	if d.HasChange("categories") {
		metadata.Categories = expandCategories(d.Get("categories"))
	}
	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").(map[string]interface{})
		metadata.OwnerReference = validateRef(or)
	}
	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		metadata.ProjectReference = validateRef(pr)
	}
	if d.HasChange("name") {
		spec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		spec.Description = utils.StringPtr(d.Get("description").(string))
	}

	if d.HasChange("user_reference_list") {
		res.UserGroupReferenceList = validateArrayRef(d.Get("user_reference_list").([]interface{}), utils.StringPtr("user"))
	}

	if d.HasChange("user_group_reference_list") {
		res.UserGroupReferenceList = validateArrayRef(d.Get("user_group_reference_list").([]interface{}), utils.StringPtr("user_group"))
	}

	if d.HasChange("role_reference") {
		res.RoleReference = validateRefList(d.Get("role_reference").([]interface{}))
	}

	if d.HasChange("filter_list") {
		res.FilterList.ContextList = expandContextFilterList(d)
	}

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	resp, errUpdate := conn.V3.UpdateAccessControlPolicy(d.Id(), request)
	if errUpdate != nil {
		return fmt.Errorf("error updating access control policy id %s): %s", d.Id(), errUpdate)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for access control policy (%s) to update: %s", d.Id(), err)
	}
	// Setting Description because in Get request is not present.
	d.Set("description", utils.StringValue(resp.Spec.Description))

	return resourceNutanixAccessControlPolicyRead(d, meta)
}

func resourceNutanixAccessControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	resp, err := conn.V3.DeleteAccessControlPolicy(d.Id())
	if err != nil {
		return fmt.Errorf("error deleting access control policy id %s): %s", d.Id(), err)
	}

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING", "DELETED_PENDING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, cast.ToString(resp.Status.ExecutionContext.TaskUUID)),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for access control policy (%s) to update: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixAccessControlPolicyExists(conn *v3.Client, name string) (*string, error) {
	var accessUUID *string

	filter := fmt.Sprintf("name==%s", name)
	accessList, err := conn.V3.ListAllAccessControlPolicy(filter)

	if err != nil {
		return nil, err
	}

	for _, access := range accessList.Entities {
		if utils.StringValue(access.Status.Name) == name {
			accessUUID = access.Metadata.UUID
		}
	}
	return accessUUID, nil
}

func expandAccessControlPolicyResources(d *schema.ResourceData, access *v3.AccessControlPolicyResources) {
	var filterList v3.FilterList

	if v, ok := d.GetOk("user_reference_list"); ok {
		access.UserReferenceList = validateArrayRef(v.(*schema.Set), utils.StringPtr("user"))
	}

	if v, ok := d.GetOk("user_group_reference_list"); ok {
		access.UserGroupReferenceList = validateArrayRef(v.(*schema.Set), utils.StringPtr("user_group"))
	}

	if v, ok := d.GetOk("role_reference"); ok {
		access.RoleReference = validateRefList(v.([]interface{}))
	}

	filterList.ContextList = expandContextFilterList(d)

	if filterList.ContextList != nil {
		access.FilterList = &filterList
	}
}

func expandContextFilterList(d *schema.ResourceData) []*v3.ContextList {
	if v1, ok := d.GetOk("filter_list"); ok {
		fl := v1.([]interface{})
		for _, filter := range fl {
			v := filter.(map[string]interface{})
			if v2, ok := v["context_list"]; ok {
				contextList := make([]*v3.ContextList, 0)
				for _, a1 := range v2.([]interface{}) {
					var context v3.ContextList
					con := a1.(map[string]interface{})

					context.ScopeFilterExpressionList = expandScopeExpressionList(con)
					context.EntityFilterExpressionList = expandEntityExpressionList(con)

					contextList = append(contextList, &context)
				}
				return contextList
			}
		}
	}
	return nil
}

func expandScopeExpressionList(con map[string]interface{}) []*v3.ScopeFilterExpressionList {
	if v2, ok := con["scope_filter_expression_list"]; ok {
		scopes := make([]*v3.ScopeFilterExpressionList, 0)
		for _, a2 := range v2.([]interface{}) {
			var scope v3.ScopeFilterExpressionList
			sco := a2.(map[string]interface{})

			if v3, ok := sco["left_hand_side"]; ok {
				scope.LeftHandSide = v3.(string)
			}
			if v3, ok := sco["operator"]; ok {
				scope.Operator = v3.(string)
			}

			scope.RightHandSide = expandRightHandSide(sco)

			scopes = append(scopes, &scope)
		}
		return scopes

	}
	return nil
}

func expandEntityExpressionList(con map[string]interface{}) []v3.EntityFilterExpressionList {
	if v2, ok := con["entity_filter_expression_list"]; ok {
		entities := make([]v3.EntityFilterExpressionList, 0)
		for _, a2 := range v2.([]interface{}) {
			var entity v3.EntityFilterExpressionList
			sco := a2.(map[string]interface{})

			if v4, ok := sco["left_hand_side_entity_type"]; ok {
				var left v3.LeftHandSide

				left.EntityType = utils.StringPtr(v4.(string))

				entity.LeftHandSide = left
			}
			if v3, ok := sco["operator"]; ok {
				entity.Operator = v3.(string)
			}

			entity.RightHandSide = expandRightHandSide(sco)

			entities = append(entities, entity)
		}
		return entities
	}
	return nil
}

func expandRightHandSide(side map[string]interface{}) v3.RightHandSide {
	var right v3.RightHandSide
	if v4, ok := side["right_hand_side"]; ok {
		vrhs := v4.([]interface{})
		for _, vrh := range vrhs {
			rhd := vrh.(map[string]interface{})

			if v5, ok := rhd["collection"]; ok {
				if v5.(string) != "" {
					right.Collection = utils.StringPtr(v5.(string))
				}
			}
			if v5, ok := rhd["categories"]; ok {
				right.Categories = expandCategories(v5)
			}
			if v5, ok := rhd["uuid_list"]; ok {
				right.UUIDList = cast.ToStringSlice(v5.(*schema.Set).List())
			}
		}

	}
	return right
}

func flattenFilterList(filterList *v3.FilterList) []interface{} {
	filters := make([]interface{}, 0)

	filter := make(map[string]interface{})
	filter["context_list"] = flattenContextList(filterList.ContextList)
	filters = append(filters, filter)

	return filters
}

func flattenContextList(contextList []*v3.ContextList) []interface{} {
	contexts := make([]interface{}, 0)
	for _, con := range contextList {
		if con != nil {
			scope := make(map[string]interface{})
			scope["scope_filter_expression_list"] = flattenScopeExpressionList(con.ScopeFilterExpressionList)
			scope["entity_filter_expression_list"] = flattenEntityExpressionList(con.EntityFilterExpressionList)

			contexts = append(contexts, scope)
		}
	}

	return contexts

}

func flattenScopeExpressionList(scopeList []*v3.ScopeFilterExpressionList) []interface{} {
	scopes := make([]interface{}, 0)

	for _, sco := range scopeList {
		scope := make(map[string]interface{})
		scope["left_hand_side"] = sco.LeftHandSide
		scope["operator"] = sco.Operator
		scope["right_hand_side"] = flattenRightHandSide(sco.RightHandSide)

		scopes = append(scopes, scope)
	}

	return scopes
}

func flattenEntityExpressionList(entities []v3.EntityFilterExpressionList) []interface{} {
	scopes := make([]interface{}, 0)

	for _, ent := range entities {
		scope := make(map[string]interface{})
		scope["left_hand_side_entity_type"] = utils.StringValue(ent.LeftHandSide.EntityType)
		scope["operator"] = ent.Operator
		scope["right_hand_side"] = flattenRightHandSide(ent.RightHandSide)

		scopes = append(scopes, scope)
	}

	return scopes
}

func flattenRightHandSide(right v3.RightHandSide) []interface{} {
	rightHand := make([]interface{}, 0)

	r := make(map[string]interface{})
	r["collection"] = utils.StringValue(right.Collection)
	r["uuid_list"] = right.UUIDList
	r["categories"] = flattenCategories(right.Categories)

	rightHand = append(rightHand, r)

	return rightHand
}
