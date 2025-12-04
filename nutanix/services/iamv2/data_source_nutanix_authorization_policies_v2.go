package iamv2

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
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
					},
				},
			},
		},
	}
}

func DatasourceNutanixAuthorizationPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	// initialize query params
	var filter, orderBy, selects, expand *string
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
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}

	resp, err := conn.AuthAPIInstance.ListAuthorizationPolicies(page, limit, filter, orderBy, expand, selects)
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

			auths[k] = auth
		}
		log.Printf("[DEBUG] flattenAuthorizationPolicyEntities return: %+v", auths[0])
		return auths
	}
	log.Printf("[DEBUG] flattenAuthorizationPolicyEntities return nil")
	return nil
}
