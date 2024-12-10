package iamv2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixUserV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixUserV2Create,
		ReadContext:   resourceNutanixUserV2Read,
		UpdateContext: resourceNutanixUserV2Update,
		DeleteContext: resourceNutanixUserV2Delete,
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
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LOCAL", "SAML", "LDAP", "EXTERNAL"}, false),
			},
			"idp_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"middle_initial": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_reset_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"additional_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type: schema.TypeString,
				// Optional: true,
				Computed: true,
			},
			"created_time": {
				Type: schema.TypeString,
				// Optional: true,
				Computed: true,
			},
			"last_updated_time": {
				Type: schema.TypeString,
				// Optional: true,
				Computed: true,
			},
			"created_by": {
				Type: schema.TypeString,
				// Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNutanixUserV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	spec := &import1.User{}

	if un, ok := d.GetOk("username"); ok {
		spec.Username = utils.StringPtr(un.(string))
	}
	if ut, ok := d.GetOk("user_type"); ok {
		const two, three, four, five = 2, 3, 4, 5
		usertypeMap := map[string]interface{}{
			"LOCAL":    two,
			"SAML":     three,
			"LDAP":     four,
			"EXTERNAL": five,
		}
		pInt := usertypeMap[ut.(string)]
		p := import1.UserType(pInt.(int))
		spec.UserType = &p
	}

	if idp, ok := d.GetOk("idp_id"); ok {
		spec.IdpId = utils.StringPtr(idp.(string))
	}
	if displayName, ok := d.GetOk("display_name"); ok {
		spec.DisplayName = utils.StringPtr(displayName.(string))
	}
	if fName, ok := d.GetOk("first_name"); ok {
		spec.FirstName = utils.StringPtr(fName.(string))
	}
	if middle, ok := d.GetOk("middle_initial"); ok {
		spec.MiddleInitial = utils.StringPtr(middle.(string))
	}

	if lName, ok := d.GetOk("last_name"); ok {
		spec.LastName = utils.StringPtr(lName.(string))
	}
	if email, ok := d.GetOk("email_id"); ok {
		spec.EmailId = utils.StringPtr(email.(string))
	}
	if locale, ok := d.GetOk("locale"); ok {
		spec.Locale = utils.StringPtr(locale.(string))
	}
	if region, ok := d.GetOk("region"); ok {
		spec.Region = utils.StringPtr(region.(string))
	}
	if pass, ok := d.GetOk("password"); ok {
		spec.Password = utils.StringPtr(pass.(string))
	}
	if frp, ok := d.GetOk("force_reset_password"); ok {
		spec.IsForceResetPasswordEnabled = utils.BoolPtr(frp.(bool))
	}
	if status, ok := d.GetOk("status"); ok {
		const two, three = 2, 3
		statusMap := map[string]interface{}{
			"ACTIVE":   two,
			"INACTIVE": three,
		}
		pInt := statusMap[status.(string)]
		p := import1.UserStatusType(pInt.(int))
		spec.Status = &p
	}
	if lastLogin, ok := d.GetOk("last_login_time"); ok {
		t := lastLogin.(time.Time)
		spec.LastLoginTime = &t
	}
	if cTime, ok := d.GetOk("created_time"); ok {
		t := cTime.(time.Time)
		spec.CreatedTime = &t
	}
	if lastUpdate, ok := d.GetOk("last_updated_time"); ok {
		t := lastUpdate.(time.Time)
		spec.LastUpdatedTime = &t
	}
	if cBy, ok := d.GetOk("created_by"); ok {
		spec.CreatedBy = utils.StringPtr(cBy.(string))
	}
	if addAttr, ok := d.GetOk("additional_attributes"); ok {
		spec.AdditionalAttributes = expandKVPair(addAttr.([]interface{}))
	}

	resp, err := conn.UsersAPIInstance.CreateUser(spec)
	if err != nil {
		return diag.Errorf("error while creating User : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.User)

	d.SetId(*getResp.ExtId)
	return resourceNutanixUserV2Read(ctx, d, meta)
}

func resourceNutanixUserV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.UsersAPIInstance.GetUserById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching user : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.User)

	if err = d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.Errorf("error setting ext_id for user %s: %s", d.Id(), err)
	}
	if err = d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("username", getResp.Username); err != nil {
		return diag.Errorf("error setting username for user %s: %s", d.Id(), err)
	}
	if err = d.Set("user_type", flattenUserType(getResp.UserType)); err != nil {
		return diag.Errorf("error setting user_type for user %s: %s", d.Id(), err)
	}
	if err = d.Set("idp_id", getResp.IdpId); err != nil {
		return diag.Errorf("error setting idp_id for user %s: %s", d.Id(), err)
	}
	if err = d.Set("display_name", getResp.DisplayName); err != nil {
		return diag.Errorf("error setting display_name for user %s: %s", d.Id(), err)
	}
	if err = d.Set("first_name", getResp.FirstName); err != nil {
		return diag.Errorf("error setting first_name for user %s: %s", d.Id(), err)
	}
	if err = d.Set("middle_initial", getResp.MiddleInitial); err != nil {
		return diag.Errorf("error setting middle_initial for user %s: %s", d.Id(), err)
	}
	if err = d.Set("last_name", getResp.LastName); err != nil {
		return diag.Errorf("error setting last_name for user %s: %s", d.Id(), err)
	}
	if err = d.Set("email_id", getResp.EmailId); err != nil {
		return diag.Errorf("error setting email_id for user %s: %s", d.Id(), err)
	}
	if err = d.Set("locale", getResp.Locale); err != nil {
		return diag.Errorf("error setting username for user %s: %s", d.Id(), err)
	}
	if err = d.Set("region", getResp.Region); err != nil {
		return diag.Errorf("error setting region for user %s: %s", d.Id(), err)
	}
	if err = d.Set("force_reset_password", getResp.IsForceResetPasswordEnabled); err != nil {
		return diag.Errorf("error setting force_reset_password for user %s: %s", d.Id(), err)
	}
	if err = d.Set("additional_attributes", flattenAdditionalAttributes(getResp)); err != nil {
		return diag.Errorf("error setting additional_attributes for user %s: %s", d.Id(), err)
	}
	if err = d.Set("status", flattenUserStatusType(getResp.Status)); err != nil {
		return diag.Errorf("error setting status for user %s: %s", d.Id(), err)
	}
	if err = d.Set("buckets_access_keys", flattenBucketsAccessKeys(getResp)); err != nil {
		return diag.Errorf("error setting buckets_access_keys for user %s: %s", d.Id(), err)
	}

	if err = d.Set("last_login_time", getResp.LastLoginTime.Format("2006-01-02T15:04:05Z07:00")); err != nil {
		return diag.Errorf("error setting last_login_time for user %s: %s", d.Id(), err)
	}
	if err = d.Set("created_time", getResp.CreatedTime.Format("2006-01-02T15:04:05Z07:00")); err != nil {
		return diag.Errorf("error setting created_time for user %s: %s", d.Id(), err)
	}
	if err = d.Set("last_updated_time", getResp.LastUpdatedTime.Format("2006-01-02T15:04:05Z07:00")); err != nil {
		return diag.Errorf("error setting last_updated_time for user %s: %s", d.Id(), err)
	}
	if err = d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.Errorf("error setting created_by for user %s: %s", d.Id(), err)
	}
	return nil
}

func resourceNutanixUserV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	updateSpec := &import1.User{}

	// get Resp
	getResp, er := conn.UsersAPIInstance.GetUserById(utils.StringPtr(d.Id()))
	if er != nil {
		return diag.FromErr(er)
	}

	getUserResp := getResp.Data.GetValue().(import1.User)

	updateSpec = &getUserResp

	// checking if attribute is updated or not

	if d.HasChange("user_type") {
		const two, three, four, five = 2, 3, 4, 5
		usertypeMap := map[string]interface{}{
			"LOCAL":    two,
			"SAML":     three,
			"LDAP":     four,
			"EXTERNAL": five,
		}
		pInt := usertypeMap[d.Get("user_type").(string)]
		p := import1.UserType(pInt.(int))
		updateSpec.UserType = &p
	}
	if d.HasChange("idp_id") {
		updateSpec.IdpId = utils.StringPtr(d.Get("idp_id").(string))
	}
	if d.HasChange("display_name") {
		updateSpec.DisplayName = utils.StringPtr(d.Get("display_name").(string))
	}
	if d.HasChange("first_name") {
		updateSpec.FirstName = utils.StringPtr(d.Get("first_name").(string))
	}
	if d.HasChange("middle_initial") {
		updateSpec.MiddleInitial = utils.StringPtr(d.Get("middle_initial").(string))
	}
	if d.HasChange("last_name") {
		updateSpec.LastName = utils.StringPtr(d.Get("last_name").(string))
	}
	if d.HasChange("email_id") {
		updateSpec.EmailId = utils.StringPtr(d.Get("email_id").(string))
	}
	if d.HasChange("locale") {
		updateSpec.Locale = utils.StringPtr(d.Get("locale").(string))
	}
	if d.HasChange("region") {
		updateSpec.Region = utils.StringPtr(d.Get("region").(string))
	}
	if d.HasChange("password") {
		updateSpec.Password = utils.StringPtr(d.Get("password").(string))
	}
	if d.HasChange("force_reset_password") {
		updateSpec.IsForceResetPasswordEnabled = utils.BoolPtr(d.Get("force_reset_password").(bool))
	}
	if d.HasChange("status") {
		const two, three = 2, 3
		statusMap := map[string]interface{}{
			"ACTIVE":   two,
			"INACTIVE": three,
		}
		pInt := statusMap[d.Get("status").(string)]
		p := import1.UserStatusType(pInt.(int))
		updateSpec.Status = &p
	}
	if d.HasChange("additional_attributes") {
		updateSpec.AdditionalAttributes = expandKVPair(d.Get("additional_attributes").([]interface{}))
	}

	// Extract E-Tag Header
	etagValue := conn.APIClientInstance.GetEtag(getResp)

	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	updateresp, err := conn.UsersAPIInstance.UpdateUserById(utils.StringPtr(d.Id()), updateSpec, args)
	if err != nil {
		return diag.FromErr(err)
	}
	updateResp := updateresp.Data.GetValue().(import1.User)

	if d.Id() != *updateResp.ExtId {
		return diag.Errorf("ext_id is different in update user")
	}
	return resourceNutanixUserV2Read(ctx, d, meta)
}

func resourceNutanixUserV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// DeleteUserById is not implemented in the GA release
	// conn := meta.(*conns.Client).IamAPI

	// readResp, err := conn.UsersAPIInstance.GetUserById(utils.StringPtr(d.Id()))
	// if err != nil {
	// 	return diag.Errorf("error while fetching user: %v", err)
	// }

	// etagValue := conn.UserGroupsAPIInstance.ApiClient.GetEtag(readResp)
	// headers := make(map[string]interface{})
	// headers["If-Match"] = utils.StringPtr(etagValue)

	// resp, err := conn.UsersAPIInstance.DeleteUserById(utils.StringPtr(d.Id()), headers)
	// if err != nil {
	// 	return diag.Errorf("error while deleting user  : %v", err)
	// }

	// if resp == nil {
	// 	log.Println("[DEBUG] User deleted successfully.")
	// }
	return nil
}

func expandKVPair(pr []interface{}) []config.KVPair {
	if len(pr) > 0 {
		kvPairs := make([]config.KVPair, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			pair := &config.KVPair{}

			if name, ok := val["name"]; ok {
				pair.Name = utils.StringPtr(name.(string))
			}
			if value, ok := val["value"]; ok {
				pair.Value = value.(*config.OneOfKVPairValue)
			}
			kvPairs[k] = *pair
		}
		return kvPairs
	}
	return nil
}

func flattenUserType(pr *import1.UserType) string {
	if pr != nil {
		const two, three, four, five = 2, 3, 4, 5
		if *pr == import1.UserType(two) {
			return "LOCAL"
		}
		if *pr == import1.UserType(three) {
			return "SAML"
		}
		if *pr == import1.UserType(four) {
			return "LDAP"
		}
		if *pr == import1.UserType(five) {
			return "EXTERNAL"
		}
	}
	return "UNKNOWN"
}

func flattenUserStatusType(pr *import1.UserStatusType) string {
	if pr != nil {
		const two, three = 2, 3

		if *pr == import1.UserStatusType(two) {
			return "ACTIVE"
		}
		if *pr == import1.UserStatusType(three) {
			return "INACTIVE"
		}
	}
	return "UNKNOWN"
}
