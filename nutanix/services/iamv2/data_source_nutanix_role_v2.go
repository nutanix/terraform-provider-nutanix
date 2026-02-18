package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamConfig "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// List Role(s)
func DatasourceNutanixRoleV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRoleV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "ExtId for the Role.",
				Type:        schema.TypeString,
				Required:    true,
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
			"display_name": {
				Description: "The display name for the Role.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the Role.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"client_name": {
				Description: "Client that created the entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"operations": {
				Description: "List of Operations for the Role.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
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
				Description: "The creation time of the Role.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated_time": {
				Description: "The time when the Role was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_by": {
				Description: "User or Service Name that created the Role.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_system_defined": {
				Description: "Flag identifying if the Role is system defined or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func DatasourceNutanixRoleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	roleExtID := d.Get("ext_id").(string)

	resp, err := conn.RolesAPIInstance.GetRoleById(&roleExtID)
	if err != nil {
		return diag.Errorf("error while fetching role: %v", err)
	}

	getResp := resp.Data.GetValue().(iamConfig.Role)

	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
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
	if err := d.Set("operations", getResp.Operations); err != nil {
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

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
