package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamResponse "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/common/v1/response"
	iamConfig "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// List Role(s)
func DatasourceNutanixRolesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRolesV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Description: "A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"limit": {
				Description: "A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
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
				},
			},
		},
	}
}

func DatasourceNutanixRolesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	// initialize query params
	var filter, orderBy, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	resp, err := conn.RolesAPIInstance.ListRoles(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching roles: %v", err)
	}

	rolesRaw := resp.Data.GetValue()
	rolesList, ok := rolesRaw.([]iamConfig.Role)
	if !ok || len(rolesList) == 0 {
		if err := d.Set("roles", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of roles.",
		}}
	}

	if err := d.Set("roles", flattenRolesEntities(rolesList)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRolesEntities(roles []iamConfig.Role) []interface{} {
	if len(roles) > 0 {
		rolesList := make([]interface{}, len(roles))

		for k, v := range roles {
			role := make(map[string]interface{})

			if v.ExtId != nil {
				role["ext_id"] = v.ExtId
			}
			if v.TenantId != nil {
				role["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				role["links"] = flattenLinks(v.Links)
			}
			if v.DisplayName != nil {
				role["display_name"] = v.DisplayName
			}
			if v.Description != nil {
				role["description"] = v.Description
			}
			if v.ClientName != nil {
				role["client_name"] = v.ClientName
			}
			if v.Operations != nil {
				role["operations"] = v.Operations
			}
			if v.AccessibleClients != nil {
				role["accessible_clients"] = v.AccessibleClients
			}
			if v.AccessibleEntityTypes != nil {
				role["accessible_entity_types"] = v.AccessibleEntityTypes
			}
			if v.AssignedUsersCount != nil {
				role["assigned_users_count"] = v.AssignedUsersCount
			}
			if v.AssignedUserGroupsCount != nil {
				role["assigned_users_groups_count"] = v.AssignedUserGroupsCount
			}
			if v.CreatedTime != nil {
				t := v.CreatedTime
				role["created_time"] = t.String()
			}
			if v.LastUpdatedTime != nil {
				t := v.LastUpdatedTime
				role["last_updated_time"] = t.String()
			}
			if v.CreatedBy != nil {
				role["created_by"] = v.CreatedBy
			}
			if v.IsSystemDefined != nil {
				role["is_system_defined"] = v.IsSystemDefined
			}

			rolesList[k] = role
		}
		return rolesList
	}
	return nil
}

func flattenLinks(apiLinks []iamResponse.ApiLink) []map[string]interface{} {
	if len(apiLinks) > 0 {
		apiLinkList := make([]map[string]interface{}, len(apiLinks))

		for k, v := range apiLinks {
			links := map[string]interface{}{}
			if v.Href != nil {
				links["href"] = v.Href
			}
			if v.Rel != nil {
				links["rel"] = v.Rel
			}

			apiLinkList[k] = links
		}
		return apiLinkList
	}
	return nil
}
