package nutanix

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixUserGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixUserGroupsCreate,
		ReadContext:   resourceNutanixUserGroupsRead,
		UpdateContext: resourceNutanixUserGroupsUpdate,
		DeleteContext: resourceNutanixUserGroupsDelete,
		Schema: map[string]*schema.Schema{
			"directory_service_user_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distinguished_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"saml_user_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idp_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"directory_service_ou": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distinguished_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"metadata": {
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
	conn := meta.(*Client).API

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
		res.DirectoryServiceUserGroup = expandDirectoryUserGroup(ds.([]interface{}))
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
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      VpcDelayTime,
		MinTimeout: VpcMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for user group (%s) to create: %s", uuid, errWaitTask)
	}

	d.SetId(uuid)
	return nil
}

func resourceNutanixUserGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	resp, err := conn.V3.GetUserGroup(d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error reading user group %s: %s", d.Id(), err)
	}

	m, _ := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return diag.Errorf("error setting metadata for VPC %s: %s", d.Id(), err)
	}

	if err = d.Set("directory_service_user_group", flattenDirectoryServiceUserGroup(resp.Spec.Resources.DirectoryServiceUserGroup)); err != nil {
		return diag.Errorf("error setting directory_service_user_group for user group %s: %s", d.Id(), err)
	}

	// if err = d.Set("directory_service_ou", flattenDirectoryServiceUserGroup(resp.Spec.Resources.DirectoryServiceUserGroup)); err != nil {
	// 	return diag.Errorf("error setting directory_service_ou for user group %s: %s", d.Id(), err)
	// }

	return nil
}

func resourceNutanixUserGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixUserGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API
	log.Printf("[DEBUG] Deleting User Group: %s", d.Id())
	resp, err := conn.V3.DeleteUserGroup(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error while deleting user group UUID(%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETE_PENDING", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, resp.Status.ExecutionContext.TaskUUID.(string)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      VpcDelayTime,
		MinTimeout: VpcMinTimeout,
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
