package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamConfig "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixUsersV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixUsersV2Read,
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
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rel": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idp_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"middle_initial": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locale": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_force_reset_password": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"additional_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"buckets_access_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"links": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"rel": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"access_key_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secret_access_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"created_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"last_login_time": {
							Type:     schema.TypeString,
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
						"last_updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceNutanixUsersV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	resp, err := conn.UsersAPIInstance.ListUsers(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching users : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("users", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of users.",
		}}
	}

	getResp := resp.Data.GetValue().([]iamConfig.User)

	if err := d.Set("users", flattenUsersEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenUsersEntities(usersResp []iamConfig.User) []interface{} {
	if len(usersResp) > 0 {
		users := make([]interface{}, len(usersResp))
		for k, v := range usersResp {
			user := make(map[string]interface{})

			if v.TenantId != nil {
				user["tenant_id"] = v.TenantId
			}
			if v.ExtId != nil {
				user["ext_id"] = v.ExtId
			}
			if v.Links != nil {
				user["links"] = flattenLinks(v.Links)
			}
			if v.Username != nil {
				user["username"] = v.Username
			}
			if v.UserType != nil {
				user["user_type"] = flattenUserType(v.UserType)
			}
			if v.Description != nil {
				user["description"] = v.Description
			}
			if v.IdpId != nil {
				user["idp_id"] = v.IdpId
			}
			if v.DisplayName != nil {
				user["display_name"] = v.DisplayName
			}
			if v.FirstName != nil {
				user["first_name"] = v.FirstName
			}
			if v.MiddleInitial != nil {
				user["middle_initial"] = v.MiddleInitial
			}
			if v.LastName != nil {
				user["last_name"] = v.LastName
			}
			if v.EmailId != nil {
				user["email_id"] = v.EmailId
			}
			if v.Locale != nil {
				user["locale"] = v.Locale
			}
			if v.Region != nil {
				user["region"] = v.Region
			}
			if v.IsForceResetPasswordEnabled != nil {
				user["is_force_reset_password"] = *v.IsForceResetPasswordEnabled
			}
			if v.AdditionalAttributes != nil {
				user["additional_attributes"] = flattenAdditionalAttributes(v)
			}
			if v.Status != nil {
				user["status"] = flattenUserStatusType(v.Status)
			}
			if v.BucketsAccessKeys != nil {
				user["buckets_access_keys"] = flattenBucketsAccessKeys(v)
			}
			if v.LastLoginTime != nil {
				user["last_login_time"] = v.LastLoginTime.Format("2006-01-02T15:04:05Z07:00")
			}
			if v.CreatedTime != nil {
				user["created_time"] = v.CreatedTime.Format("2006-01-02T15:04:05Z07:00")
			}
			if v.LastUpdatedTime != nil {
				user["last_updated_time"] = v.LastUpdatedTime.Format("2006-01-02T15:04:05Z07:00")
			}
			if v.CreatedBy != nil {
				user["created_by"] = v.CreatedBy
			}
			if v.LastUpdatedBy != nil {
				user["last_updated_by"] = v.LastUpdatedBy
			}

			users[k] = user
		}
		return users
	}
	return nil
}

func flattenBucketsAccessKeys(user iamConfig.User) []map[string]interface{} {
	if len(user.BucketsAccessKeys) > 0 {
		bucketsAccessKeysList := make([]map[string]interface{}, len(user.BucketsAccessKeys))

		for k, v := range user.BucketsAccessKeys {
			bucketsAccessKeys := map[string]interface{}{}
			if v.ExtId != nil {
				bucketsAccessKeys["ext_id"] = *v.ExtId
			}
			if v.Links != nil {
				bucketsAccessKeys["links"] = flattenLinks(v.Links)
			}
			if v.AccessKeyName != nil {
				bucketsAccessKeys["access_key_name"] = *v.AccessKeyName
			}
			if v.SecretAccessKey != nil {
				bucketsAccessKeys["secret_access_key"] = *v.SecretAccessKey
			}
			if v.UserId != nil {
				bucketsAccessKeys["user_id"] = *v.UserId
			}
			if v.CreatedTime != nil {
				createdTime := v.CreatedTime
				bucketsAccessKeys["created_time"] = createdTime.String()
			}

			bucketsAccessKeysList[k] = bucketsAccessKeys
		}
		return bucketsAccessKeysList
	}
	return nil
}

func flattenAdditionalAttributes(user iamConfig.User) []interface{} {
	if len(user.AdditionalAttributes) > 0 {
		additionalAttributes := make([]interface{}, len(user.AdditionalAttributes))
		for i, attr := range user.AdditionalAttributes {
			additionalAttributes[i] = map[string]interface{}{
				"name":  *attr.Name,
				"value": *attr.Value,
			}
		}
		return additionalAttributes
	}
	return nil
}
