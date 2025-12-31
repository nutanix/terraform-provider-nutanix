package iamv2

import (
	"context"
	"encoding/json"
	"log"
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				ValidateFunc: validation.StringInSlice([]string{"LOCAL", "SAML", "LDAP", "EXTERNAL", "SERVICE_ACCOUNT"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
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

func resourceNutanixUserV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	spec := &import1.User{}

	if un, ok := d.GetOk("username"); ok {
		spec.Username = utils.StringPtr(un.(string))
	}
	if ut, ok := d.GetOk("user_type"); ok {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		usertypeMap := map[string]interface{}{
			"LOCAL":           two,
			"SAML":            three,
			"LDAP":            four,
			"EXTERNAL":        five,
			"SERVICE_ACCOUNT": six,
		}
		pInt := usertypeMap[ut.(string)]
		p := import1.UserType(pInt.(int))
		spec.UserType = &p
	}
	if description, ok := d.GetOk("description"); ok {
		spec.Description = utils.StringPtr(description.(string))
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

	aJSON, _ := json.MarshalIndent(spec, "", "  ")
	log.Printf("[DEBUG] create user spec: %s", aJSON)

	resp, err := conn.UsersAPIInstance.CreateUser(spec)
	if err != nil {
		return diag.Errorf("error while creating User : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.User)

	d.SetId(utils.StringValue(getResp.ExtId))
	return resourceNutanixUserV2Read(ctx, d, meta)
}

func resourceNutanixUserV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.UsersAPIInstance.GetUserById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching user : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.User)

	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] resourceNutanixUserV2Read: get user response: %s", aJSON)

	if err = d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.Errorf("error setting ext_id for user %s: %s", d.Id(), err)
	}
	if err = d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("username", getResp.Username); err != nil {
		return diag.Errorf("error setting username for user/service account %s: %s", d.Id(), err)
	}
	if err = d.Set("user_type", flattenUserType(getResp.UserType)); err != nil {
		return diag.Errorf("error setting user_type for user %s: %s", d.Id(), err)
	}
	if err = d.Set("description", getResp.Description); err != nil {
		return diag.Errorf("error setting description for user/service account %s: %s", d.Id(), err)
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
		return diag.Errorf("error setting email_id for user/service account %s: %s", d.Id(), err)
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

	// get Resp
	getResp, er := conn.UsersAPIInstance.GetUserById(utils.StringPtr(d.Id()))
	if er != nil {
		return diag.FromErr(er)
	}

	getUserResp := getResp.Data.GetValue().(import1.User)

	updateSpec := &getUserResp

	// validation on update spec
	// Note: user read response has "" as default value for  middleInitial, emailId, displayName.
	if updateSpec.MiddleInitial != nil && utils.StringValue(updateSpec.MiddleInitial) == "" {
		updateSpec.MiddleInitial = nil
	}
	if updateSpec.EmailId != nil && utils.StringValue(updateSpec.EmailId) == "" {
		updateSpec.EmailId = nil
	}
	if updateSpec.DisplayName != nil && utils.StringValue(updateSpec.DisplayName) == "" {
		updateSpec.DisplayName = nil
	}

	// checking if attribute is updated or not

	if d.HasChange("user_type") {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		usertypeMap := map[string]interface{}{
			"LOCAL":           two,
			"SAML":            three,
			"LDAP":            four,
			"EXTERNAL":        five,
			"SERVICE_ACCOUNT": six,
		}
		pInt := usertypeMap[d.Get("user_type").(string)]
		p := import1.UserType(pInt.(int))
		updateSpec.UserType = &p
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("idp_id") {
		updateSpec.IdpId = utils.StringPtr(d.Get("idp_id").(string))
	}
	if d.HasChange("display_name") {
		updateSpec.DisplayName = utils.StringPtr(d.Get("display_name").(string))
	} else if utils.StringValue(updateSpec.DisplayName) == "" {
		// If display_name is empty value (""), we should not send it in the update request
		updateSpec.DisplayName = nil
	}
	if d.HasChange("first_name") {
		updateSpec.FirstName = utils.StringPtr(d.Get("first_name").(string))
	} else if utils.StringValue(updateSpec.FirstName) == "" {
		// If first_name is empty value (""), we should not send it in the update request
		updateSpec.FirstName = nil
	}
	if d.HasChange("middle_initial") {
		updateSpec.MiddleInitial = utils.StringPtr(d.Get("middle_initial").(string))
	} else if utils.StringValue(updateSpec.MiddleInitial) == "" {
		// If middle_initial is empty value (""), we should not send it in the
		// update request
		updateSpec.MiddleInitial = nil
	}
	if d.HasChange("last_name") {
		updateSpec.LastName = utils.StringPtr(d.Get("last_name").(string))
	} else if utils.StringValue(updateSpec.LastName) == "" {
		// If last_name is empty value (""), we should not send it in the
		// update request
		updateSpec.LastName = nil
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

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] update user spec: %s", aJSON)

	updateresp, err := conn.UsersAPIInstance.UpdateUserById(utils.StringPtr(d.Id()), updateSpec, args)

	if err != nil {
		return diag.FromErr(err)
	}
	updateResp := updateresp.Data.GetValue().(import1.User)

	if d.Id() != utils.StringValue(updateResp.ExtId) {
		return diag.Errorf("ext_id is different in update user")
	}
	return resourceNutanixUserV2Read(ctx, d, meta)
}

func resourceNutanixUserV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixUserV2Delete : Delete not supported yet")
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  "Delete operation not supported",
			Detail:   "Deleting users via Terraform is not supported yet. Please delete the user manually from the Prism Central UI if required, or use v3 resource for now",
		},
	}
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
		const two, three, four, five, six = 2, 3, 4, 5, 6
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
		if *pr == import1.UserType(six) {
			return "SERVICE_ACCOUNT"
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
