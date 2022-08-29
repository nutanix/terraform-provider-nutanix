package nutanix

import (
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
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

func generateCustomFilters(PermissionList []string) []interface{} {
	customMap := make(map[string]v3.EntityFilterExpressionList)
	view_image := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("image"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["view_image"] = view_image

	view_app_icon := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("app_icon"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_App_Icon"] = view_app_icon

	view_name_category := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("category"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_Name_Category"] = view_name_category

	create_or_update_name_category := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("category"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["Create_Or_Update_Name_Category"] = create_or_update_name_category

	view_environment := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("environment"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("SELF_OWNED"),
		},
	}
	customMap["View_Environment"] = view_environment

	view_marketplace_item := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("marketplace_item"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("SELF_OWNED"),
		},
	}
	customMap["View_Marketplace_Item"] = view_marketplace_item

	view_user := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("user_group"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_User"] = view_user

	view_user_group := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("image"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_User_Group"] = view_user_group

	view_role := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("role"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_Role"] = view_role

	view_directory_service := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("directory_service"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_Directory_Service"] = view_directory_service

	search_directory_service := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("directory_service"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["Search_Directory_Service"] = search_directory_service

	view_identity_provider := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("identity_provider"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("ALL"),
		},
	}
	customMap["View_Identity_Provider"] = view_identity_provider

	view_app_task := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("app_task"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("SELF_OWNED"),
		},
	}
	customMap["View_App_Task"] = view_app_task

	view_app_variable := v3.EntityFilterExpressionList{
		Operator: "IN",
		LeftHandSide: v3.LeftHandSide{
			EntityType: utils.StringPtr("app_variable"),
		},
		RightHandSide: v3.RightHandSide{
			Collection: utils.StringPtr("SELF_OWNED"),
		},
	}
	customMap["View_App_Variable"] = view_app_variable

	filterLIST := make([]interface{}, 0)

	for _, v := range PermissionList {
		val := v
		if vals, ok := customMap[val]; ok {
			filterLIST = append(filterLIST, vals)
		}
	}

	return filterLIST
}

func getRolesPermission(roleId string, meta interface{}, projectUUID string, clusterUUID string, pcCollab bool) []*v3.ContextList {
	conn := meta.(*Client).API

	resp, _ := conn.V3.GetRole(roleId)

	roleName := utils.StringValue(resp.Status.Name)

	if roleName == "Developer" || roleName == "Consumer" || roleName == ("Operator") || roleName == ("Project Admin") {

		contextOut := make([]*v3.ContextList, 0)

		contextOut = append(contextOut, getFilterCollaboration(pcCollab, projectUUID)...)

		contextOut = append(contextOut, getSystemRoles(roleName, projectUUID, clusterUUID)...)

		contextOut = append(contextOut, getDefaultContext(projectUUID)...)

		return contextOut
	} else {
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
}
