package iam

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	userGroupDelayTime  = 2 * time.Second
	userGroupMinTimeout = 5 * time.Second
)

func ResourceNutanixUserGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixUserGroupsCreate,
		ReadContext:   resourceNutanixUserGroupsRead,
		UpdateContext: resourceNutanixUserGroupsUpdate,
		DeleteContext: resourceNutanixUserGroupsDelete,
		Schema: map[string]*schema.Schema{
			"directory_service_user_group": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"saml_user_group", "directory_service_ou"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distinguished_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"saml_user_group": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"directory_service_user_group", "directory_service_ou"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idp_uuid": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"directory_service_ou": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"directory_service_user_group", "saml_user_group"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distinguished_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNutanixUserGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.UserGroupIntentInput{}
	spec := &v3.UserGroupSpec{}
	res := &v3.UserGroupResources{}
	metadata := &v3.Metadata{}

	if err := getMetadataAttributes(d, metadata, "user_group"); err != nil {
		return diag.FromErr(err)
	}

	if ds, ok := d.GetOk("directory_service_user_group"); ok {
		res.DirectoryServiceUserGroup = expandDirectoryUserGroup(ds.([]interface{}))
	}

	if ds, ok := d.GetOk("directory_service_ou"); ok {
		res.DirectoryServiceOU = expandDirectoryUserGroup(ds.([]interface{}))
	}

	if su, ok := d.GetOk("saml_user_group"); ok {
		res.SamlUserGroup = expandSamlUserGroup(su.([]interface{}))
	}

	spec.Resources = res
	request.Spec = spec
	request.Metadata = metadata

	// Create User Group API

	resp, err := conn.V3.CreateUserGroup(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}
	uuid := *resp.Metadata.UUID
	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the UserGroup to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      userGroupDelayTime,
		MinTimeout: userGroupMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for user group (%s) to create: %s", uuid, errWaitTask)
	}

	d.SetId(uuid)
	return resourceNutanixUserGroupsRead(ctx, d, meta)
}

func resourceNutanixUserGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	resp, err := conn.V3.GetUserGroup(d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return diag.FromErr(err)
		}
		return diag.Errorf("error reading user group %s: %s", d.Id(), err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return diag.Errorf("error setting metadata for User Groups %s: %s", d.Id(), err)
	}

	if err = d.Set("categories", c); err != nil {
		return diag.Errorf("error setting categories for user UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return diag.Errorf("error setting owner_reference for user UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("directory_service_user_group", flattenDirectoryServiceUserGroup(resp.Spec.Resources.DirectoryServiceUserGroup)); err != nil {
		return diag.Errorf("error setting directory_service_user_group for user group %s: %s", d.Id(), err)
	}

	if err = d.Set("directory_service_ou", flattenDirectoryServiceUserGroup(resp.Spec.Resources.DirectoryServiceOU)); err != nil {
		return diag.Errorf("error setting directory_service_ou for user group %s: %s", d.Id(), err)
	}

	if err = d.Set("saml_user_group", flattenSamlUserGroup(resp.Spec.Resources.SamlUserGroup)); err != nil {
		return diag.Errorf("error setting saml_user_group for user group %s: %s", d.Id(), err)
	}

	return nil
}

func resourceNutanixUserGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.UserGroupIntentInput{}
	spec := &v3.UserGroupSpec{}
	metadata := &v3.Metadata{}
	res := &v3.UserGroupResources{}

	response, err := conn.V3.GetUserGroup(d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return diag.FromErr(err)
		}
		return diag.Errorf("error reading User Group %s: %s", d.Id(), err)
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

	if d.HasChange("directory_service_user_group") {
		res.DirectoryServiceUserGroup = expandDirectoryUserGroup(d.Get("directory_service_user_group").([]interface{}))
	}

	if d.HasChange("directory_service_ou") {
		res.DirectoryServiceUserGroup = expandDirectoryUserGroup(d.Get("directory_service_ou").([]interface{}))
	}

	if d.HasChange("saml_user_group") {
		res.SamlUserGroup = expandSamlUserGroup(d.Get("saml_user_group").([]interface{}))
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

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	// Make request to the API
	resp, err := conn.V3.UpdateUserGroup(ctx, d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the User Group to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      userGroupDelayTime,
		MinTimeout: userGroupMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for user group (%s) to create: %s", d.Id(), errWaitTask)
	}

	return resourceNutanixUserGroupsRead(ctx, d, meta)
}

func resourceNutanixUserGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API
	log.Printf("[DEBUG] Deleting User Group: %s", d.Id())
	resp, err := conn.V3.DeleteUserGroup(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return diag.FromErr(err)
		}
		return diag.Errorf("error while deleting user group UUID(%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETE_PENDING", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, resp.Status.ExecutionContext.TaskUUID.(string)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      userGroupDelayTime,
		MinTimeout: userGroupMinTimeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf(
			"error waiting for user group (%s) to delete: %s", d.Id(), err)
	}
	d.SetId("")

	return nil
}

func expandDirectoryUserGroup(pr []interface{}) *v3.DirectoryServiceUserGroup {
	if len(pr) > 0 {
		res := &v3.DirectoryServiceUserGroup{}
		ent := pr[0].(map[string]interface{})

		if pnum, pk := ent["distinguished_name"]; pk && len(pnum.(string)) > 0 {
			res.DistinguishedName = utils.StringPtr(pnum.(string))
		}
		return res
	}
	return nil
}

func expandSamlUserGroup(pr []interface{}) *v3.SamlUserGroup {
	if len(pr) > 0 {
		res := &v3.SamlUserGroup{}
		ent := pr[0].(map[string]interface{})

		if idp, iok := ent["idp_uuid"]; iok {
			res.IdpUUID = utils.StringPtr(idp.(string))
		}

		if name, nok := ent["name"]; nok {
			res.Name = utils.StringPtr(name.(string))
		}

		return res
	}
	return nil
}

func flattenSamlUserGroup(su *v3.SamlUserGroup) []interface{} {
	if su != nil {
		res := make([]interface{}, 0)
		sug := make(map[string]interface{})

		sug["idp_uuid"] = su.IdpUUID
		sug["name"] = su.Name

		res = append(res, sug)
		return res
	}
	return nil
}
