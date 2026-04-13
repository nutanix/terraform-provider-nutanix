package iamv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamCommonConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/common/v1/config"
	iamConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authz"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Description: "A HATEOAS style link for the response.",
		Type:        schema.TypeList,
		Computed:    true,
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
	}
}

func schemaForKeyValuePairs() *schema.Schema {
	return &schema.Schema{
		Description: "Key-value pairs for the role membership.",
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func schemaForScopeTemplateNameValues() *schema.Schema {
	return &schema.Schema{
		Description: "Name value pairs to substitute in the scope template variables referenced by the role membership.",
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func flattenRmIdentityType(identityType *iamConfig.RmIdentityType) string {
	if identityType == nil {
		return ""
	}
	switch *identityType {
	case iamConfig.RMIDENTITYTYPE_USER:
		return "USER"
	case iamConfig.RMIDENTITYTYPE_GROUP:
		return "GROUP"
	default:
		return "UNKNOWN"
	}
}

func expandRmIdentityType(identityType string) *iamConfig.RmIdentityType {
	var val iamConfig.RmIdentityType
	switch identityType {
	case "USER":
		val = iamConfig.RMIDENTITYTYPE_USER
	case "GROUP":
		val = iamConfig.RMIDENTITYTYPE_GROUP
	default:
		val = iamConfig.RMIDENTITYTYPE_UNKNOWN
	}
	return &val
}

func flattenScopeTemplateNameValues(kvPairs []iamCommonConfig.KVPair) []map[string]interface{} {
	if len(kvPairs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(kvPairs))
	for i, kv := range kvPairs {
		kvMap := map[string]interface{}{}
		if kv.Name != nil {
			kvMap["name"] = utils.StringValue(kv.Name)
		}
		if kv.Value != nil {
			val := kv.Value.GetValue()
			if val != nil {
				if strVal, ok := val.(*string); ok && strVal != nil {
					kvMap["value"] = utils.StringValue(strVal)
				} else if strVal, ok := val.(string); ok {
					kvMap["value"] = strVal
				}
			}
		}
		result[i] = kvMap
	}
	return result
}

func expandScopeTemplateNameValues(kvPairs []interface{}) []iamCommonConfig.KVPair {
	if len(kvPairs) == 0 {
		return nil
	}
	result := make([]iamCommonConfig.KVPair, len(kvPairs))
	for i, kv := range kvPairs {
		kvMap := kv.(map[string]interface{})
		kvPair := iamCommonConfig.KVPair{}
			if v, ok := kvMap["name"]; ok && v.(string) != "" {
			kvPair.Name = utils.StringPtr(v.(string))
		}
		if v, ok := kvMap["value"]; ok && v.(string) != "" {
			oneOfValue := iamCommonConfig.NewOneOfKVPairValue()
			oneOfValue.SetValue(v.(string))
			kvPair.Value = oneOfValue
		}
		result[i] = kvPair
	}
	return result
}
