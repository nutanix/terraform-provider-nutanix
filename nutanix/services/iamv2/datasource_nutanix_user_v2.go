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

func DatasourceNutanixUserV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixUserV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
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
		},
	}
}

func datasourceNutanixUserV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("ext_id")

	resp, err := conn.UsersAPIInstance.GetUserById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching user : %v", err)
	}

	getResp := resp.Data.GetValue().(iamConfig.User)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.Errorf("error setting ext_id: %v", err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.Errorf("error setting links: %v", err)
	}
	if err := d.Set("username", getResp.Username); err != nil {
		return diag.Errorf("error setting username: %v", err)
	}
	if err := d.Set("user_type", flattenUserType(getResp.UserType)); err != nil {
		return diag.Errorf("error setting user_type: %v", err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.Errorf("error setting description: %v", err)
	}
	if err := d.Set("idp_id", getResp.IdpId); err != nil {
		return diag.Errorf("error setting idp_id: %v", err)
	}
	if err := d.Set("display_name", getResp.DisplayName); err != nil {
		return diag.Errorf("error setting display_name: %v", err)
	}
	if err := d.Set("first_name", getResp.FirstName); err != nil {
		return diag.Errorf("error setting first_name: %v", err)
	}
	if err := d.Set("middle_initial", getResp.MiddleInitial); err != nil {
		return diag.Errorf("error setting middle_initial: %v", err)
	}
	if err := d.Set("last_name", getResp.LastName); err != nil {
		return diag.Errorf("error setting last_name: %v", err)
	}
	if err := d.Set("email_id", getResp.EmailId); err != nil {
		return diag.Errorf("error setting email_id: %v", err)
	}
	if err := d.Set("locale", getResp.Locale); err != nil {
		return diag.Errorf("error setting locale: %v", err)
	}
	if err := d.Set("region", getResp.Region); err != nil {
		return diag.Errorf("error setting region: %v", err)
	}
	if err := d.Set("password", getResp.Password); err != nil {
		return diag.Errorf("error setting password: %v", err)
	}
	if err := d.Set("is_force_reset_password", getResp.IsForceResetPasswordEnabled); err != nil {
		return diag.Errorf("error setting is_force_reset_password: %v", err)
	}
	if err := d.Set("additional_attributes", flattenAdditionalAttributes(getResp)); err != nil {
		return diag.Errorf("error setting additional_attributes for user %s: %s", d.Id(), err)
	}
	if err := d.Set("status", flattenUserStatusType(getResp.Status)); err != nil {
		return diag.Errorf("error setting status for user %s: %s", d.Id(), err)
	}
	if err := d.Set("buckets_access_keys", flattenBucketsAccessKeys(getResp)); err != nil {
		return diag.Errorf("error setting buckets_access_keys for user %s: %s", d.Id(), err)
	}
	if err := d.Set("last_login_time", getResp.LastLoginTime.Format("2006-01-02T15:04:05Z07:00")); err != nil {
		return diag.Errorf("error setting last_login_time for user %s: %s", d.Id(), err)
	}
	if err := d.Set("created_time", getResp.CreatedTime.Format("2006-01-02T15:04:05Z07:00")); err != nil {
		return diag.Errorf("error setting created_time for user %s: %s", d.Id(), err)
	}
	if err := d.Set("last_updated_time", getResp.LastUpdatedTime.Format("2006-01-02T15:04:05Z07:00")); err != nil {
		return diag.Errorf("error setting last_updated_time for user %s: %s", d.Id(), err)
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.Errorf("error setting created_by: %v", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
