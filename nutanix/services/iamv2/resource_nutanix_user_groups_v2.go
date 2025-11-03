package iamv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixUserGroupsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixUserGroupsV4Create,
		ReadContext:   ResourceNutanixUserGroupsV4Read,
		UpdateContext: ResourceNutanixUserGroupsV4Update,
		DeleteContext: ResourceNutanixUserGroupsV4Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SAML", "LDAP"}, false),
			},
			"idp_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"distinguished_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
		},
	}
}

func ResourceNutanixUserGroupsV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI
	input := &import1.UserGroup{}

	if gType, ok := d.GetOk("group_type"); ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"SAML": two,
			"LDAP": three,
		}
		pInt := subMap[gType.(string)]
		p := import1.GroupType(pInt.(int))
		input.GroupType = &p
	}

	if idp, ok := d.GetOk("idp_id"); ok {
		input.IdpId = utils.StringPtr(idp.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		input.Name = utils.StringPtr(name.(string))
	}
	if dName, ok := d.GetOk("distinguished_name"); ok {
		input.DistinguishedName = utils.StringPtr(dName.(string))
	}

	resp, err := conn.UserGroupsAPIInstance.CreateUserGroup(input)
	if err != nil {
		return diag.Errorf("error while creating user groups: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.UserGroup)
	d.SetId(*getResp.ExtId)
	return ResourceNutanixUserGroupsV4Read(ctx, d, meta)
}

func ResourceNutanixUserGroupsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.UserGroupsAPIInstance.GetUserGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching user groups: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.UserGroup)
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_type", flattenGroupType(getResp.GroupType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("idp_id", getResp.IdpId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("distinguished_name", getResp.DistinguishedName); err != nil {
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
	return nil
}

func ResourceNutanixUserGroupsV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixUserGroupsV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	readResp, err := conn.UserGroupsAPIInstance.GetUserGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching role: %v", err)
	}

	etagValue := conn.UserGroupsAPIInstance.ApiClient.GetEtag(readResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	resp, err := conn.UserGroupsAPIInstance.DeleteUserGroupById(utils.StringPtr(d.Id()), headers)
	if err != nil {
		return diag.Errorf("error while deleting user group : %v", err)
	}

	if resp == nil {
		log.Println("[DEBUG] User group deleted successfully.")
	}
	return nil
}
