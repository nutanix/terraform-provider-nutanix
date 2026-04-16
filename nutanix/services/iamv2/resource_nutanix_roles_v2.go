package iamv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/common/v1/config"
	iamConfig "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRolesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRolesV4Create,
		ReadContext:   ResourceNutanixRolesV4Read,
		UpdateContext: ResourceNutanixRolesV4Update,
		DeleteContext: ResourceNutanixRolesV4Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "ExtId for the Role.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"display_name": {
				Description: "The display name for the Role.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the Role.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"client_name": {
				Description: "Client that created the entity.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"operations": {
				Description: "List of Operations for the Role.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Description: "A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Description: "The URL at which the entity described by the link can be accessed.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"rel": {
							Description: "A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of \"self\" identifies the URL for the object.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"accessible_clients": {
				Description: "List of Accessible Clients for the Role.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"accessible_entity_types": {
				Description: "List of Accessible Entity Types for the Role.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"assigned_users_count": {
				Description: "Number of Users assigned to given Role.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"assigned_users_groups_count": {
				Description: "Number of User Groups assigned to given Role.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_system_defined": {
				Description: "Flag identifying if the Role is system defined or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func ResourceNutanixRolesV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI
	body := &iamConfig.Role{}

	if extID, ok := d.GetOk("ext_id"); ok {
		body.ExtId = utils.StringPtr(extID.(string))
	}
	if displayName, ok := d.GetOk("display_name"); ok {
		body.DisplayName = utils.StringPtr(displayName.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(description.(string))
	}
	if clientName, ok := d.GetOk("client_name"); ok {
		body.ClientName = utils.StringPtr(clientName.(string))
	}
	if operations, ok := d.GetOk("operations"); ok {
		operationsList := operations.([]interface{})
		operationsListStr := make([]string, len(operationsList))
		for i, v := range operationsList {
			operationsListStr[i] = v.(string)
		}
		body.Operations = operationsListStr
	}

	resp, err := conn.RolesAPIInstance.CreateRole(body)
	if err != nil {
		return diag.Errorf("error while creating role: %v", err)
	}

	getResp := resp.Data.GetValue().(iamConfig.Role)
	d.SetId(utils.StringValue(getResp.ExtId))
	return ResourceNutanixRolesV4Read(ctx, d, meta)
}

func ResourceNutanixRolesV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.RolesAPIInstance.GetRoleById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while Reading role: %v", err)
	}

	getResp := resp.Data.GetValue().(iamConfig.Role)

	// after creating role, operations saved in remote in different order than local
	if len(getResp.Operations) > 0 {
		// read the remote operations and local operations list
		remoteOperations := getResp.Operations
		localOperations := d.Get("operations").([]interface{})

		// final result for checking if operations are different
		diff := false

		// convert local operations to string slice
		localOperationsStr := make([]string, len(localOperations))
		for i, v := range localOperations {
			localOperationsStr[i] = (v.(string))
		}

		log.Printf("[DEBUG] localOperationsStr: %v", localOperationsStr)

		// check if remote operations are different from local operations
		for _, operation := range remoteOperations {
			offset := indexOf(localOperationsStr, operation)

			if offset == -1 {
				log.Printf("[DEBUG] Operation %v not found in local operations", operation)
				diff = true
				break
			}
		}

		// if operations are different, update local operations
		if diff {
			log.Printf("[DEBUG] Operations are different. Updating local operations")
			if err := d.Set("operations", getResp.Operations); err != nil {
				return diag.FromErr(err)
			}
		} else {
			// if operations are same, do not update local operations
			log.Printf("[DEBUG] Operations are same. Not updating local operations")
		}
	}

	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_name", getResp.DisplayName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_name", getResp.ClientName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("accessible_clients", getResp.AccessibleClients); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("accessible_entity_types", getResp.AccessibleEntityTypes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("assigned_users_count", getResp.AssignedUsersCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("assigned_users_groups_count", getResp.AssignedUserGroupsCount); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		t := getResp.CreatedTime
		if err := d.Set("created_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdatedTime != nil {
		t := getResp.LastUpdatedTime
		if err := d.Set("last_updated_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_system_defined", getResp.IsSystemDefined); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixRolesV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := utils.StringPtr(d.Id())

	updatedSpec := iamConfig.Role{}

	readResp, err := conn.RolesAPIInstance.GetRoleById(extID)
	if err != nil {
		return diag.Errorf("error while fetching role: %v", err)
	}

	// get etag value from read response to pass in update request If-Match header, Required for update request
	etagValue := conn.RolesAPIInstance.ApiClient.GetEtag(readResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	updatedSpec = readResp.Data.GetValue().(iamConfig.Role)

	if d.HasChange("display_name") {
		updatedSpec.DisplayName = utils.StringPtr(d.Get("display_name").(string))
	}
	if d.HasChange("description") {
		updatedSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("client_name") {
		updatedSpec.ClientName = utils.StringPtr(d.Get("client_name").(string))
	}
	if d.HasChange("operations") {
		operations := d.Get("operations").([]interface{})
		operationsListStr := make([]string, len(operations))
		for i, v := range operations {
			operationsListStr[i] = v.(string)
		}
		updatedSpec.Operations = operationsListStr
	}

	updateResp, err := conn.RolesAPIInstance.UpdateRoleById(extID, &updatedSpec, headers)
	if err != nil {
		return diag.Errorf("error while updating role: %v", err)
	}
	log.Printf("[DEBUG] Role updated. Response: %v", *updateResp)

	updateTaskResp := updateResp.Data.GetValue().(config.Message)

	if updateTaskResp.Message != nil {
		log.Printf("[DEBUG] %v", *updateTaskResp.Message)
	}
	return ResourceNutanixRolesV4Read(ctx, d, meta)
}

func ResourceNutanixRolesV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	readResp, err := conn.RolesAPIInstance.GetRoleById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching role: %v", err)
	}

	etagValue := conn.RolesAPIInstance.ApiClient.GetEtag(readResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	resp, err := conn.RolesAPIInstance.DeleteRoleById(utils.StringPtr(d.Id()), headers)
	if err != nil {
		return diag.Errorf("error while Deleting role: %v", err)
	}

	if resp == nil {
		log.Println("[DEBUG] Role deleted successfully.")
	}
	return nil
}

func indexOf(slice []string, target string) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}
