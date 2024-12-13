package iamv2

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAuthorizationPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixAuthorizationPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func DatasourceNutanixAuthorizationPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DatasourceNutanixAuthorizationPolicyV2Read \n")
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("ext_id")
	resp, err := conn.AuthAPIInstance.GetAuthorizationPolicyById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching authorization polices: %v", err)
	}
	getResp := resp.Data.GetValue().(import1.AuthorizationPolicy)

	if err := d.Set("display_name", getResp.DisplayName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_name", getResp.ClientName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identities", flattenIdentityFilters(getResp.Identities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entities", flattenEntityFilters(getResp.Entities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role", getResp.Role); err != nil {
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
	if err := d.Set("is_system_defined", getResp.IsSystemDefined); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("authorization_policy_type", flattenAuthorizationPolicyType(getResp.AuthorizationPolicyType)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}

func flattenAuthorizationPolicyType(pr *import1.AuthorizationPolicyType) string {
	if pr != nil {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		if *pr == import1.AuthorizationPolicyType(two) {
			return "USER_DEFINED"
		}
		if *pr == import1.AuthorizationPolicyType(three) {
			return "SERVICE_DEFINED"
		}
		if *pr == import1.AuthorizationPolicyType(four) {
			return "PREDEFINED_READ_ONLY"
		}
		if *pr == import1.AuthorizationPolicyType(five) {
			return "PREDEFINED_UPDATE_IDENTITY_ONLY"
		}
		if *pr == import1.AuthorizationPolicyType(six) {
			return "SERVICE_DEFINED_READ_ONLY"
		}
	}
	return "UNKNOWN"
}

func flattenIdentityFilters(identityFilters []import1.IdentityFilter) []interface{} {
	if len(identityFilters) > 0 {
		identities := make([]interface{}, len(identityFilters))
		log.Printf("[DEBUG] flattenIdentityFilters \n")
		for k, v := range identityFilters {
			identity := make(map[string]interface{})
			log.Printf("[DEBUG] flattenIdentityFilters  %v:%v\n", k, v)
			log.Printf("[DEBUG] flattenIdentityFilters val type : %v\n", reflect.TypeOf(v))

			reservedMap, err := json.Marshal(v.Reserved_)
			if err != nil {
				log.Printf("[DEBUG] flattenIdentityFiltersError [%v]:%v err : %v\n", k, v, err)
			}
			log.Printf("[DEBUG] flattenIdentityFilters reserved : %v\n", string(reservedMap))
			identity["reserved"] = string(reservedMap)

			identities[k] = identity
		}
		return identities
	}
	return nil
}

func flattenEntityFilters(entityFilters []import1.EntityFilter) []interface{} {
	if len(entityFilters) > 0 {
		entities := make([]interface{}, len(entityFilters))

		for k, v := range entityFilters {
			entity := make(map[string]interface{})
			log.Printf("[DEBUG] flattenIdentityFilters  %v:%v\n", k, v)
			log.Printf("[DEBUG] flattenIdentityFilters val type : %v\n", reflect.TypeOf(v))
			reservedMap, err := json.Marshal(v.Reserved_)
			if err != nil {
				log.Printf("[DEBUG] flattenIdentityFiltersError [%v]:%v err : %v\n", k, v, err)
			}
			log.Printf("[DEBUG] flattenIdentityFilters reserved : %v\n", string(reservedMap))
			entity["reserved"] = string(reservedMap)
			entities[k] = entity
		}
		return entities
	}
	return nil
}
