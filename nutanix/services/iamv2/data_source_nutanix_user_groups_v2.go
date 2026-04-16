package iamv2

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixUserGroupsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixUserGroupsV4Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
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
			"user_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixUserGroupV2(),
			},
		},
	}
}

func DatasourceNutanixUserGroupsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	resp, err := conn.UserGroupsAPIInstance.ListUserGroups(page, limit, filter, orderBy, selects)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching user groups: %v", errorMessage["message"])
	}

	if resp.Data == nil {
		if err := d.Set("user_groups", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of user groups.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.UserGroup)
	if err := d.Set("user_groups", flattenUserGroupEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenUserGroupEntities(userGroups []import1.UserGroup) []interface{} {
	if len(userGroups) > 0 {
		ugs := make([]interface{}, len(userGroups))

		for k, userGroup := range userGroups {
			ug := make(map[string]interface{})

			if userGroup.TenantId != nil {
				ug["tenant_id"] = userGroup.TenantId
			}
			if userGroup.ExtId != nil {
				ug["ext_id"] = userGroup.ExtId
			}
			if userGroup.Links != nil {
				ug["links"] = flattenLinks(userGroup.Links)
			}
			if userGroup.GroupType != nil {
				ug["group_type"] = flattenGroupType(userGroup.GroupType)
			}
			if userGroup.IdpId != nil {
				ug["idp_id"] = userGroup.IdpId
			}
			if userGroup.Name != nil {
				ug["name"] = userGroup.Name
			}
			if userGroup.DistinguishedName != nil {
				ug["distinguished_name"] = userGroup.DistinguishedName
			}
			if userGroup.CreatedBy != nil {
				ug["created_by"] = userGroup.CreatedBy
			}
			if userGroup.CreatedTime != nil {
				t := userGroup.CreatedTime
				ug["created_time"] = t.String()
			}
			if userGroup.LastUpdatedTime != nil {
				t := userGroup.LastUpdatedTime
				ug["last_updated_time"] = t.String()
			}

			ugs[k] = ug
		}
		return ugs
	}
	return nil
}
