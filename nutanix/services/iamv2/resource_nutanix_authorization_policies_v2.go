package iamv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixAuthPoliciesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixAuthPoliciesV2Create,
		ReadContext:   ResourceNutanixAuthPoliciesV2Read,
		UpdateContext: ResourceNutanixAuthPoliciesV2Update,
		DeleteContext: ResourceNutanixAuthPoliciesV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, rd *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				log.Printf("[DEBUG] Importing Authorization Policy")
				return []*schema.ResourceData{rd}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"client_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identities": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reserved": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: SuppressEquivalentAuthPolicyDiffs,
							StateFunc: func(v interface{}) string {
								log.Printf("[DEBUG] StateFunc value: %v\n", v)
								json, err := structure.NormalizeJsonString(v)
								if err != nil {
									log.Printf("[DEBUG] StateFunc err : %v\n", err)
								}
								return json
							},
						},
					},
				},
			},

			"entities": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reserved": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: SuppressEquivalentAuthPolicyDiffs,
							StateFunc: func(v interface{}) string {
								log.Printf("[DEBUG] StateFunc v : %v\n", v)
								json, err := structure.NormalizeJsonString(v)
								if err != nil {
									log.Printf("[DEBUG] StateFunc err : %v\n", err)
								}
								return json
							},
						},
					},
				},
			},

			"role": {
				Type:     schema.TypeString,
				Required: true,
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
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PREDEFINED_READ_ONLY", "SERVICE_DEFINED_READ_ONLY",
					"PREDEFINED_UPDATE_IDENTITY_ONLY", "SERVICE_DEFINED", "USER_DEFINED",
				}, false),
			},
		},
	}
}

func ResourceNutanixAuthPoliciesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating Authorization Policy")
	conn := meta.(*conns.Client).IamAPI
	input := &import1.AuthorizationPolicy{}
	log.Printf("[DEBUG] Creating Authorization Policy")

	if display, ok := d.GetOk("display_name"); ok {
		input.DisplayName = utils.StringPtr(display.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		input.Description = utils.StringPtr(desc.(string))
	}
	if identities, ok := d.GetOk("identities"); ok {
		identities, err := expandIdentityFilter(identities.([]interface{}))
		if err != nil {
			return diag.Errorf("error while creating  Authorization Policy in identities err: %v", err)
		}

		input.Identities = identities
	}
	if entities, ok := d.GetOk("entities"); ok {
		entities, err := expandEntityFilter(entities.([]interface{}))
		if err != nil {
			return diag.Errorf("error while creating  Authorization Policy in entities err: %v", err)
		}

		input.Entities = entities
	}
	if role, ok := d.GetOk("role"); ok {
		input.Role = utils.StringPtr(role.(string))
	}
	if authPolicyType, ok := d.GetOk("authorization_policy_type"); ok {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		subMap := map[string]interface{}{
			"USER_DEFINED":                    two,
			"SERVICE_DEFINED":                 three,
			"PREDEFINED_READ_ONLY":            four,
			"PREDEFINED_UPDATE_IDENTITY_ONLY": five,
			"SERVICE_DEFINED_READ_ONLY":       six,
		}
		pInt := subMap[authPolicyType.(string)]
		p := import1.AuthorizationPolicyType(pInt.(int))
		input.AuthorizationPolicyType = &p
	}

	resp, err := conn.AuthAPIInstance.CreateAuthorizationPolicy(input)
	if err != nil {
		return diag.Errorf("error while creating authorization policies : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.AuthorizationPolicy)

	log.Printf("[DEBUG] Creating Authorization Policy Return")

	d.Set("ext_id", *getResp.ExtId)
	d.SetId(*getResp.ExtId)
	return ResourceNutanixAuthPoliciesV2Read(ctx, d, meta)
}

func ResourceNutanixAuthPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Authorization Policy")
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.AuthAPIInstance.GetAuthorizationPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching authorization polices: %v", err)
	}
	getResp := resp.Data.GetValue().(import1.AuthorizationPolicy)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
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
	return nil
}

func ResourceNutanixAuthPoliciesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Authorization Policy")
	conn := meta.(*conns.Client).IamAPI
	updatedSpec := import1.AuthorizationPolicy{}

	resp, err := conn.AuthAPIInstance.GetAuthorizationPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching  Authorization Policy: %v", err)
	}

	etagValue := conn.AuthAPIInstance.ApiClient.GetEtag(resp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	updatedSpec = resp.Data.GetValue().(import1.AuthorizationPolicy)

	if d.HasChange("display_name") {
		updatedSpec.DisplayName = utils.StringPtr(d.Get("display_name").(string))
	}
	if d.HasChange("description") {
		updatedSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("identities") {
		identities, errID := expandIdentityFilter(d.Get("identities").([]interface{}))
		if errID != nil {
			return diag.Errorf("error while updating  Authorization Policy in identities err: %v", errID)
		}
		updatedSpec.Identities = identities
	}
	if d.HasChange("entities") {
		entities, errEn := expandEntityFilter(d.Get("entities").([]interface{}))
		if errEn != nil {
			return diag.Errorf("error while updating  Authorization Policy in entities err: %v", errEn)
		}
		updatedSpec.Entities = entities
	}
	if d.HasChange("role") {
		updatedSpec.Role = utils.StringPtr(d.Get("role").(string))
	}
	if d.HasChange("authorization_policy_type") {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		subMap := map[string]interface{}{
			"USER_DEFINED":                    two,
			"SERVICE_DEFINED":                 three,
			"PREDEFINED_READ_ONLY":            four,
			"PREDEFINED_UPDATE_IDENTITY_ONLY": five,
			"SERVICE_DEFINED_READ_ONLY":       six,
		}
		pInt := subMap[d.Get("authorization_policy_type").(string)]
		p := import1.AuthorizationPolicyType(pInt.(int))
		updatedSpec.AuthorizationPolicyType = &p
	}

	updatedResp, err := conn.AuthAPIInstance.UpdateAuthorizationPolicyById(utils.StringPtr(d.Id()), &updatedSpec, headers)
	if err != nil {
		return diag.Errorf("error while updating  Authorization Policy: %v", err)
	}

	updatedResponse := updatedResp.Data.GetValue().(config.Message)
	log.Printf("[DEBUG] updatedResponse : %v\n", updatedResponse)

	if updatedResponse.Message != nil {
		log.Println("[DEBUG] updated the Authorization Policy")
	}
	return ResourceNutanixAuthPoliciesV2Read(ctx, d, meta)
}

func ResourceNutanixAuthPoliciesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting Authorization Policy")
	conn := meta.(*conns.Client).IamAPI

	readResp, err := conn.AuthAPIInstance.GetAuthorizationPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching auth policy: %v", err)
	}

	etagValue := conn.AuthAPIInstance.ApiClient.GetEtag(readResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)
	resp, err := conn.AuthAPIInstance.DeleteAuthorizationPolicyById(utils.StringPtr(d.Id()), headers)
	if err != nil {
		return diag.Errorf("error while deleting auth policy : %v", err)
	}

	if resp == nil {
		log.Println("[DEBUG] auth policy deleted successfully.")
	}
	return nil
}

func expandIdentityFilter(identities []interface{}) ([]import1.IdentityFilter, error) {
	if len(identities) > 0 {
		filters := make([]import1.IdentityFilter, len(identities))

		log.Printf("[DEBUG] expandIdentityFilter ")
		for key, value := range identities {
			log.Printf("[DEBUG] expandIdentityFilter key:%v\n", key)
			item, ok := value.(map[string]interface{})
			if !ok {
				// Handle error or continue based on requirements
				log.Printf("[DEBUG] expandIdentityFilter continue\n")
				continue
			}
			log.Printf("[DEBUG] expandIdentityFilter item : %v\n", item)
			log.Printf("[DEBUG] expandIdentityFilter item type : %v\n", reflect.TypeOf(item))
			filter := import1.IdentityFilter{}

			if val, exists := item["reserved"]; exists {
				// Assuming the field is of type string, adjust the type assertion accordingly
				log.Printf("[DEBUG] expandIdentityFilter val : %v\n", val.(string))
				log.Printf("[DEBUG] expandIdentityFilter val type : %v\n", reflect.TypeOf(val))
				reserved, err := deserializeJSONStringToMap(val.(string))
				if err != nil {
					return nil, fmt.Errorf("%s", err.Error())
				}
				log.Printf("[DEBUG] expandIdentityFilter reserved : %v\n", reserved)
				filter.Reserved_ = reserved
			}
			// Repeat for other fields as necessary
			log.Printf("[DEBUG] expandIdentityFilter key:%v  filter.Reserved_ : %v\n", key, filter.Reserved_)
			filters[key] = filter
		}
		return filters, nil
	}
	return nil, nil
}

func expandEntityFilter(entities []interface{}) ([]import1.EntityFilter, error) {
	if len(entities) > 0 {
		filters := make([]import1.EntityFilter, len(entities))

		log.Printf("[DEBUG] expandEntityFilter ")
		for key, value := range entities {
			log.Printf("[DEBUG] expandEntityFilter key:%v\n", key)
			item, ok := value.(map[string]interface{})
			if !ok {
				// Handle error or continue based on requirements
				log.Printf("[DEBUG] expandEntityFilter continue\n")
				continue
			}
			log.Printf("[DEBUG] expandEntityFilter item : %v\n", item)
			log.Printf("[DEBUG] expandEntityFilter item type : %v\n", reflect.TypeOf(item))
			filter := import1.EntityFilter{}

			if val, exists := item["reserved"]; exists {
				// Assuming the field is of type string, adjust the type assertion accordingly
				log.Printf("[DEBUG] expandEntityFilter val : %v\n", val.(string))
				log.Printf("[DEBUG] expandEntityFilter val type : %v\n", reflect.TypeOf(val))
				reserved, err := deserializeJSONStringToMap(val.(string))
				if err != nil {
					return nil, fmt.Errorf("%s", err.Error())
				}
				log.Printf("[DEBUG] expandEntityFilter reserved : %v\n", reserved)
				filter.Reserved_ = reserved
			}
			// Repeat for other fields as necessary
			log.Printf("[DEBUG] expandEntityFilter key:%v  filter.Reserved_ : %v\n", key, filter.Reserved_)
			filters[key] = filter
		}
		return filters, nil
	}
	return nil, nil
}

func deserializeJSONStringToMap(jsonString string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		log.Printf("[DEBUG] deserializeJSONStringToMap err : %v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG] deserializeJSONStringToMap map : %v\n", m)
	return m, nil
}

func SuppressEquivalentAuthPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	return AuthPolicyStringsEquivalent(old, new)
}

func AuthPolicyStringsEquivalent(s1, s2 string) bool {
	if strings.TrimSpace(s1) == "" && strings.TrimSpace(s2) == "" {
		log.Printf("[DEBUG] AuthPolicyStringsEquivalent Both strings are empty")
		return true
	}

	if strings.TrimSpace(s1) == "{}" && strings.TrimSpace(s2) == "" {
		log.Printf("[DEBUG] AuthPolicyStringsEquivalent s1 is empty and s2 is {}")
		return true
	}

	if strings.TrimSpace(s1) == "" && strings.TrimSpace(s2) == "{}" {
		log.Printf("[DEBUG] AuthPolicyStringsEquivalent s1 is {} and s2 is empty")
		return true
	}

	if strings.TrimSpace(s1) == "{}" && strings.TrimSpace(s2) == "{}" {
		log.Printf("[DEBUG] AuthPolicyStringsEquivalent Both strings are {}")
		return true
	}
	log.Printf("[DEBUG] AuthPolicyStringsEquivalent s1: %s, s2: %s return false", s1, s2)

	return false
}

// SuppressEquivalentJSONDiffs returns a difference suppression function that compares
// two JSON strings and returns `true` if they are semantically equivalent.
func SuppressEquivalentJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	return JSONStringsEqual(old, new)
}

func JSONStringsEqual(s1, s2 string) bool {
	b1 := bytes.NewBufferString("")
	if err := json.Compact(b1, []byte(s1)); err != nil {
		return false
	}

	b2 := bytes.NewBufferString("")
	if err := json.Compact(b2, []byte(s2)); err != nil {
		return false
	}

	return JSONBytesEqual(b1.Bytes(), b2.Bytes())
}

func JSONBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
