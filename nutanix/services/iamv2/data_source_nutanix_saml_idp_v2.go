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

func DatasourceNutanixSamlIDPV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSamlIDPV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"idp_metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logout_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_id_policy_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups_delim": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"entity_issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_signed_authn_req_enabled": {
				Type:     schema.TypeBool,
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

func DatasourceNutanixSamlIDPV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("ext_id")

	resp, err := conn.SamlIdentityAPIInstance.GetSamlIdentityProviderById(utils.StringPtr(extID.(string)))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching saml identity providers: %v", errorMessage["message"])
	}

	getResp := resp.Data.GetValue().(import1.SamlIdentityProvider)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("idp_metadata", flattenIdpMetadata(getResp.IdpMetadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("username_attribute", getResp.UsernameAttribute); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_attribute", getResp.EmailAttribute); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("groups_attribute", getResp.GroupsAttribute); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("groups_delim", getResp.GroupsDelim); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_attributes", getResp.CustomAttributes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_issuer", getResp.EntityIssuer); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_signed_authn_req_enabled", getResp.IsSignedAuthnReqEnabled); err != nil {
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

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenIdpMetadata(pr *import1.IdpMetadata) []map[string]interface{} {
	if pr != nil {
		idps := make([]map[string]interface{}, 0)
		idp := make(map[string]interface{})

		idp["entity_id"] = pr.EntityId
		idp["login_url"] = pr.LoginUrl
		idp["logout_url"] = pr.LogoutUrl
		idp["error_url"] = pr.ErrorUrl
		idp["certificate"] = pr.Certificate
		if pr.NameIdPolicyFormat != nil {
			idp["name_id_policy_format"] = flattenNameIDPolicyFormat(pr.NameIdPolicyFormat)
		}

		idps = append(idps, idp)
		return idps
	}
	return nil
}

func flattenNameIDPolicyFormat(pr *import1.NameIdPolicyFormat) string {
	if pr != nil {
		const two, three, four, five, six, seven, eight, nine, ten = 2, 3, 4, 5, 6, 7, 8, 9, 10

		if *pr == import1.NameIdPolicyFormat(two) {
			return "emailAddress"
		}
		if *pr == import1.NameIdPolicyFormat(three) {
			return "unspecified"
		}
		if *pr == import1.NameIdPolicyFormat(four) {
			return "X509SubjectName"
		}
		if *pr == import1.NameIdPolicyFormat(five) {
			return "WindowsDomainQualifiedName"
		}
		if *pr == import1.NameIdPolicyFormat(six) {
			return "encrypted"
		}
		if *pr == import1.NameIdPolicyFormat(seven) {
			return "entity"
		}
		if *pr == import1.NameIdPolicyFormat(eight) {
			return "kerberos"
		}
		if *pr == import1.NameIdPolicyFormat(nine) {
			return "persistent"
		}
		if *pr == import1.NameIdPolicyFormat(ten) {
			return "transient"
		}
	}
	return "UNKNOWN"
}
