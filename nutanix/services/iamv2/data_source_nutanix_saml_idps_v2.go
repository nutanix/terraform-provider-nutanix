package iamv2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authn"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/samlidentityproviders"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixSamlIDPsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSamlIDPsV2Read,
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
			"identity_providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
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
						"project_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shared_with_projects": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"share_with_all_projects": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixSamlIDPsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	listSamlIdentityProvidersRequest := import2.ListSamlIdentityProvidersRequest{}
	if v, ok := d.GetOk("page"); ok {
		listSamlIdentityProvidersRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listSamlIdentityProvidersRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listSamlIdentityProvidersRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listSamlIdentityProvidersRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listSamlIdentityProvidersRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.SamlIdentityAPIInstance.ListSamlIdentityProviders(ctx, &listSamlIdentityProvidersRequest)
	if err != nil {
		fmt.Println(err)
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching identity providers: %v", errorMessage["message"])
	}

	if resp.Data == nil {
		if err := d.Set("identity_providers", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of identity providers.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.SamlIdentityProvider)
	if err := d.Set("identity_providers", flattenIdentityProvidersEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenIdentityProvidersEntities(pr []import1.SamlIdentityProvider) []interface{} {
	if len(pr) > 0 {
		idps := make([]interface{}, len(pr))

		for k, v := range pr {
			idp := make(map[string]interface{})

			idp["ext_id"] = v.ExtId

			if v.Name != nil {
				idp["name"] = v.Name
			}
			if v.IdpMetadata != nil {
				idp["idp_metadata"] = flattenIdpMetadata(v.IdpMetadata)
			}
			if v.UsernameAttribute != nil {
				idp["username_attribute"] = v.UsernameAttribute
			}
			if v.EmailAttribute != nil {
				idp["email_attribute"] = v.EmailAttribute
			}
			if v.GroupsAttribute != nil {
				idp["groups_attribute"] = v.GroupsAttribute
			}
			if v.GroupsDelim != nil {
				idp["groups_delim"] = v.GroupsDelim
			}
			if v.CustomAttributes != nil {
				idp["custom_attributes"] = v.CustomAttributes
			}
			if v.EntityIssuer != nil {
				idp["entity_issuer"] = v.EntityIssuer
			}
			if v.IsSignedAuthnReqEnabled != nil {
				idp["is_signed_authn_req_enabled"] = v.IsSignedAuthnReqEnabled
			}
			if v.CreatedTime != nil {
				t := v.CreatedTime
				idp["created_time"] = t.String()
			}
			if v.LastUpdatedTime != nil {
				t := v.LastUpdatedTime
				idp["last_updated_time"] = t.String()
			}
			if v.CreatedBy != nil {
				idp["created_by"] = v.CreatedBy
			}
			idp["project_ext_id"] = v.ProjectExtId
			idp["shared_with_projects"] = v.SharedWithProjects
			idp["share_with_all_projects"] = v.SharedWithAllProjects
			idps[k] = idp
		}
		return idps
	}
	return nil
}
