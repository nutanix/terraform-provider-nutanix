package prism

import (
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func getSystemRoles(kind string, projectUUID string, clusterUUID string) []*v3.ContextList {
	if kind == "Project Admin" {
		return []*v3.ContextList{
			{
				EntityFilterExpressionList: []v3.EntityFilterExpressionList{
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("image"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("marketplace_item"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("directory_service"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("role"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("project"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							UUIDList: []string{projectUUID},
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("user"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("user_group"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("environment"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_icon"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("category"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_task"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("cluster"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							UUIDList: []string{clusterUUID},
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_variable"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("identity_provider"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("vm_recovery_point"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						Operator: "IN",
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("virtual_network"),
						},
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
				},
			},
		}
	}
	if kind == "Developer" {
		return []*v3.ContextList{
			{
				EntityFilterExpressionList: []v3.EntityFilterExpressionList{
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_icon"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_task"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_variable"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("category"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("cluster"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							UUIDList: []string{clusterUUID},
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("image"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("marketplace_item"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("vm_recovery_point"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						Operator: "IN",
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("virtual_network"),
						},
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
				},
			},
		}
	}
	if kind == "Consumer" {
		return []*v3.ContextList{
			{
				EntityFilterExpressionList: []v3.EntityFilterExpressionList{
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("image"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("marketplace_item"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_icon"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("category"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("cluster"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							UUIDList: []string{clusterUUID},
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_task"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_variable"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("SELF_OWNED"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("vm_recovery_point"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						Operator: "IN",
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("virtual_network"),
						},
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
				},
			},
		}
	}
	if kind == "Operator" {
		return []*v3.ContextList{
			{
				EntityFilterExpressionList: []v3.EntityFilterExpressionList{
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("app_icon"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("category"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("vm_recovery_point"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
					{
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("cluster"),
						},
						Operator: "IN",
						RightHandSide: v3.RightHandSide{
							UUIDList: []string{clusterUUID},
						},
					},
					{
						Operator: "IN",
						LeftHandSide: v3.LeftHandSide{
							EntityType: utils.StringPtr("virtual_network"),
						},
						RightHandSide: v3.RightHandSide{
							Collection: utils.StringPtr("ALL"),
						},
					},
				},
			},
		}
	}
	return nil
}

func getDefaultContext(projectUUID string) []*v3.ContextList {
	return []*v3.ContextList{
		{
			EntityFilterExpressionList: []v3.EntityFilterExpressionList{
				{
					LeftHandSide: v3.LeftHandSide{
						EntityType: utils.StringPtr("blueprint"),
					},
					Operator: "IN",
					RightHandSide: v3.RightHandSide{
						Collection: utils.StringPtr("ALL"),
					},
				},
				{
					LeftHandSide: v3.LeftHandSide{
						EntityType: utils.StringPtr("environment"),
					},
					Operator: "IN",
					RightHandSide: v3.RightHandSide{
						Collection: utils.StringPtr("ALL"),
					},
				},
				{
					LeftHandSide: v3.LeftHandSide{
						EntityType: utils.StringPtr("marketplace_item"),
					},
					Operator: "IN",
					RightHandSide: v3.RightHandSide{
						Collection: utils.StringPtr("ALL"),
					},
				},
			},
			ScopeFilterExpressionList: []*v3.ScopeFilterExpressionList{
				{
					LeftHandSide: "PROJECT",
					Operator:     "IN",
					RightHandSide: v3.RightHandSide{
						UUIDList: []string{projectUUID},
					},
				},
			},
		},
	}
}

func getFilterCollaboration(collab bool, projectUUID string) []*v3.ContextList {
	var FilterScope string
	if collab {
		FilterScope = "ALL"
	} else {
		FilterScope = "SELF_OWNED"
	}
	return []*v3.ContextList{
		{
			ScopeFilterExpressionList: []*v3.ScopeFilterExpressionList{
				{
					LeftHandSide: "PROJECT",
					Operator:     "IN",
					RightHandSide: v3.RightHandSide{
						UUIDList: []string{projectUUID},
					},
				},
			},
			EntityFilterExpressionList: []v3.EntityFilterExpressionList{
				{
					LeftHandSide: v3.LeftHandSide{
						EntityType: utils.StringPtr("ALL"),
					},
					Operator: "IN",
					RightHandSide: v3.RightHandSide{
						Collection: utils.StringPtr(FilterScope),
					},
				},
			},
		},
	}
}

// Custom Roles

func generateCustomFilters(permissionList []string) []interface{} {
	customMap := make(map[string]v3.EntityFilterExpressionList)
	viewImage := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("image"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["view_image"] = viewImage

	viewAppIcon := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("app_icon"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_App_Icon"] = viewAppIcon

	viewNameCategory := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("category"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_Name_Category"] = viewNameCategory

	createOrUpdateNameCategory := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("category"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["Create_Or_Update_Name_Category"] = createOrUpdateNameCategory

	viewEnvironment := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("environment"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("SELF_OWNED"),
		},
	}
	customMap["View_Environment"] = viewEnvironment

	viewMarketplaceItem := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("marketplace_item"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("SELF_OWNED"),
		},
	}
	customMap["View_Marketplace_Item"] = viewMarketplaceItem

	viewUser := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("user_group"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_User"] = viewUser

	viewUserGroup := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("image"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_User_Group"] = viewUserGroup

	viewRole := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("role"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_Role"] = viewRole

	viewDirectoryService := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("directory_service"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_Directory_Service"] = viewDirectoryService

	searchDirectoryService := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("directory_service"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["Search_Directory_Service"] = searchDirectoryService

	viewIdentityProvider := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("identity_provider"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_Identity_Provider"] = viewIdentityProvider

	viewAppTask := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("app_task"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("SELF_OWNED"),
		},
	}
	customMap["View_App_Task"] = viewAppTask

	viewAppVariable := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("app_variable"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("SELF_OWNED"),
		},
	}
	customMap["View_App_Variable"] = viewAppVariable

	filterLIST := make([]interface{}, 0)

	for _, v := range permissionList {
		val := v
		if vals, ok := customMap[val]; ok {
			filterLIST = append(filterLIST, vals)
		}
	}

	return filterLIST
}

func getRolesPermission(roleID string, meta interface{}, projectUUID string, clusterUUID string, pcCollab bool) []*v3.ContextList {
	conn := meta.(*conns.Client).API

	resp, _ := conn.V3.GetRole(roleID)

	roleName := utils.StringValue(resp.Status.Name)

	if roleName == "Developer" || roleName == "Consumer" || roleName == ("Operator") || roleName == ("Project Admin") {
		contextOut := make([]*v3.ContextList, 0)

		contextOut = append(contextOut, getFilterCollaboration(pcCollab, projectUUID)...)

		contextOut = append(contextOut, getSystemRoles(roleName, projectUUID, clusterUUID)...)

		contextOut = append(contextOut, getDefaultContext(projectUUID)...)

		return contextOut
	}
	permList := resp.Status.Resources.PermissionReferenceList

	permS := []string{}

	for _, v := range permList {
		permS = append(permS, utils.StringValue(v.Name))
	}

	customPerms := generateCustomFilters(permS)

	ans := make([]v3.EntityFilterExpressionList, 0)
	for _, v := range customPerms {
		val := v.(v3.EntityFilterExpressionList)
		ans = append(ans, val)
	}

	if clusterUUID != "" {
		ans = append(ans, v3.EntityFilterExpressionList{
			LeftHandSide: v3.LeftHandSide{
				EntityType: utils.StringPtr("cluster"),
			},
			Operator: "IN",
			RightHandSide: v3.RightHandSide{
				UUIDList: []string{clusterUUID},
			},
		})
	}

	out := &v3.ContextList{}
	out.EntityFilterExpressionList = ans

	contextOut := make([]*v3.ContextList, 0)

	contextOut = append(contextOut, getFilterCollaboration(pcCollab, projectUUID)...)

	contextOut = append(contextOut, out)

	contextOut = append(contextOut, getDefaultContext(projectUUID)...)

	return contextOut
}
