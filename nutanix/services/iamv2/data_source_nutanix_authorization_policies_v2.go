package iamv2

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authz"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/authorizationpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAuthorizationPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixAuthorizationPoliciesV2Read,
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
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reserved": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"entities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reserved": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"role": {
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
						"is_system_defined": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"authorization_policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
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

func DatasourceNutanixAuthorizationPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	listAuthorizationPoliciesRequest := import2.ListAuthorizationPoliciesRequest{}

	if v, ok := d.GetOk("page"); ok {
		listAuthorizationPoliciesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listAuthorizationPoliciesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listAuthorizationPoliciesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listAuthorizationPoliciesRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listAuthorizationPoliciesRequest.Select_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("expand"); ok {
		listAuthorizationPoliciesRequest.Expand_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.AuthAPIInstance.ListAuthorizationPolicies(ctx, &listAuthorizationPoliciesRequest)
	if err != nil {
		fmt.Println(err)
		return diag.Errorf("error while fetching auth policies: %v", err)
	}

	// Log the value of resp.Data.GetValue() under DEBUG level
	log.Printf("[DEBUG] resp.Data.GetValue(): %+v\n", resp.Data.GetValue())

	getVal := resp.ObjectType_

	if *getVal == "iam.v4.authz.AuthorizationPolicyProjection" {
		fmt.Println("policyProjection")
	}

	if resp.Data == nil {
		if err := d.Set("auth_policies", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of authorization policies.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.AuthorizationPolicyProjection)

	if err := d.Set("auth_policies", flattenAuthorizationPolicyEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenAuthorizationPolicyEntities(pr []import1.AuthorizationPolicyProjection) []interface{} {
	if len(pr) > 0 {
		auths := make([]interface{}, len(pr))

		for k, v := range pr {
			log.Printf("[DEBUG] flattenAuthorizationPolicyEntities[%v].ExtId: %v", k, *v.ExtId)
			auth := make(map[string]interface{})

			auth["ext_id"] = v.ExtId
			auth["display_name"] = v.DisplayName
			auth["description"] = v.Description
			auth["client_name"] = v.ClientName
			auth["entities"] = flattenEntityFilters(v.Entities)

			auth["identities"] = flattenIdentityFilters(v.Identities)
			log.Printf("[DEBUG] flattenAuthorizationPolicyEntities[%v].Identities: %v", k, v.Identities)
			log.Printf("[DEBUG] flattenAuthorizationPolicyEntities[%v].Identities: %v", k, auth["identities"])
			auth["role"] = v.Role
			if v.CreatedTime != nil {
				t := v.CreatedTime
				auth["created_time"] = t.String()
			}
			if v.LastUpdatedTime != nil {
				t := v.LastUpdatedTime
				auth["last_updated_time"] = t.String()
			}
			auth["created_by"] = v.CreatedBy
			auth["is_system_defined"] = v.IsSystemDefined
			auth["authorization_policy_type"] = flattenAuthorizationPolicyType(v.AuthorizationPolicyType)
			auth["project_ext_id"] = v.ProjectExtId
			auth["share_with_all_projects"] = v.SharedWithAllProjects
			auths[k] = auth
		}
		log.Printf("[DEBUG] flattenAuthorizationPolicyEntities return: %+v", auths[0])
		return auths
	}
	log.Printf("[DEBUG] flattenAuthorizationPolicyEntities return nil")
	return nil
}
