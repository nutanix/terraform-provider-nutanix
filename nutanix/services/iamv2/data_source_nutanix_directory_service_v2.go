package iamv2

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixDirectoryServiceV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixDirectoryServiceV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"open_ldap_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_object_class": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_search_base": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"username_attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"user_group_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_object_class": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_search_base": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_member_attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_member_attribute_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"group_search_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"white_listed_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func DatasourceNutanixDirectoryServiceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("ext_id")
	resp, err := conn.DirectoryServiceAPIInstance.GetDirectoryServiceById(utils.StringPtr(extID.(string)))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching address group : %v", errorMessage["message"])
	}

	getResp := resp.Data.GetValue().(import1.DirectoryService)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", getResp.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secondary_urls", getResp.SecondaryUrls); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_name", getResp.DomainName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("directory_type", flattenDirectoryType(getResp.DirectoryType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("service_account", flattenDsServiceAccount(getResp.ServiceAccount)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("open_ldap_configuration", flattenOpenLdapConfig(getResp.OpenLdapConfiguration)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_search_type", flattenGroupSearchType(getResp.GroupSearchType)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("white_listed_groups", getResp.WhiteListedGroups); err != nil {
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

	d.SetId(*getResp.ExtId)
	return nil
}

func flattenDirectoryType(pr *import1.DirectoryType) string {
	if pr != nil {
		const two, three = 2, 3

		if *pr == import1.DirectoryType(two) {
			return "ACTIVE_DIRECTORY"
		}
		if *pr == import1.DirectoryType(three) {
			return "OPEN_LDAP"
		}
	}
	return "UNKNOWN"
}

func flattenDsServiceAccount(pr *import1.DsServiceAccount) []map[string]interface{} {
	if pr != nil {
		accs := make([]map[string]interface{}, 0)
		acc := make(map[string]interface{})

		acc["username"] = pr.Username
		acc["password"] = pr.Password

		accs = append(accs, acc)
		return accs
	}
	return nil
}

func flattenOpenLdapConfig(pr *import1.OpenLdapConfig) []map[string]interface{} {
	if pr != nil {
		accs := make([]map[string]interface{}, 0)
		acc := make(map[string]interface{})

		if pr.UserConfiguration != nil {
			acc["user_configuration"] = flattenUserConfiguration(pr.UserConfiguration)
		}
		if pr.UserGroupConfiguration != nil {
			acc["user_group_configuration"] = flattenUserGroupConfiguration(pr.UserGroupConfiguration)
		}
		accs = append(accs, acc)
		return accs
	}
	return nil
}

func flattenUserConfiguration(pr *import1.UserConfiguration) []map[string]interface{} {
	if pr != nil {
		configs := make([]map[string]interface{}, 0)
		cfg := make(map[string]interface{})

		if pr.UsernameAttribute != nil {
			cfg["username_attribute"] = pr.UsernameAttribute
		}
		if pr.UserSearchBase != nil {
			cfg["user_search_base"] = pr.UserSearchBase
		}
		if pr.UserObjectClass != nil {
			cfg["user_object_class"] = pr.UserObjectClass
		}

		configs = append(configs, cfg)
		return configs
	}
	return nil
}

func flattenUserGroupConfiguration(pr *import1.UserGroupConfiguration) []map[string]interface{} {
	if pr != nil {
		groups := make([]map[string]interface{}, 0)
		grp := make(map[string]interface{})

		if pr.GroupMemberAttribute != nil {
			grp["group_object_class"] = pr.GroupMemberAttribute
		}
		if pr.GroupMemberAttributeValue != nil {
			grp["group_search_base"] = pr.GroupMemberAttributeValue
		}
		if pr.GroupObjectClass != nil {
			grp["group_member_attribute"] = pr.GroupObjectClass
		}
		if pr.GroupSearchBase != nil {
			grp["group_member_attribute_value"] = pr.GroupSearchBase
		}

		groups = append(groups, grp)
		return groups
	}
	return nil
}

func flattenGroupSearchType(pr *import1.GroupSearchType) string {
	const two, three = 2, 3
	if pr != nil {
		if *pr == import1.GroupSearchType(two) {
			return "NON_RECURSIVE"
		}
		if *pr == import1.GroupSearchType(three) {
			return "RECURSIVE"
		}
	}
	return "UNKNOWN"
}
