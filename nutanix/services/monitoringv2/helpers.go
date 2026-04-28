package monitoringv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/response"
	monitoringCommon "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/common"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The URL at which the entity described by the link can be accessed.",
				},
				"rel": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of \"self\" identifies the URL for the object.",
				},
			},
		},
	}
}

func schemaForEntityReferences() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of all the entities that are affected by the event or audit.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ext_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "UUID of the entity.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The name of the entity.",
				},
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The type of entity. For example, VM, node, or cluster.",
				},
			},
		},
	}
}

func schemaForEntityReference() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ext_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "UUID of the entity.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The name of the entity.",
				},
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The type of entity. For example, VM, node, or cluster.",
				},
			},
		},
	}
}

func schemaForAuditEntityReference() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ext_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "UUID of the entity.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The name of the entity.",
				},
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The type of entity. For example, VM, node, or cluster.",
				},
			},
		},
	}
}

func schemaForUserReference() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ext_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unique UUID of the user who initiated the operation.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The name of the user who initiated the operation.",
				},
				"ip_address": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The IP address from where the operation was triggered.",
				},
			},
		},
	}
}

func schemaForParameters() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Additional parameters associated with the audit. These parameters can be used to indicate custom key-value pairs for a given audit instance. For example, a service down audit in Prism Central can have the service name as a parameter.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"param_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name or key of additional parameter for an instance.",
				},
				"param_value": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Value of additional parameter for an instance.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"string_value": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Denotes a value of type string.",
							},
							"bool_value": {
								Type:        schema.TypeBool,
								Computed:    true,
								Description: "Denotes a value of type boolean.",
							},
							"int_value": {
								Type:        schema.TypeInt,
								Computed:    true,
								Description: "Denotes a value of type integer.",
							},
						},
					},
				},
			},
		},
	}
}

func flattenLinks(links []response.ApiLink) []map[string]interface{} {
	if len(links) > 0 {
		linkList := make([]map[string]interface{}, 0)
		for _, link := range links {
			linkMap := make(map[string]interface{})
			if link.Rel != nil {
				linkMap["rel"] = utils.StringValue(link.Rel)
			}
			if link.Href != nil {
				linkMap["href"] = utils.StringValue(link.Href)
			}
			linkList = append(linkList, linkMap)
		}
		return linkList
	}
	return []map[string]interface{}{}
}

func flattenEntityReferences(refs []monitoringCommon.EntityReference) []map[string]interface{} {
	if len(refs) > 0 {
		refList := make([]map[string]interface{}, len(refs))
		for i, ref := range refs {
			refMap := make(map[string]interface{})
			refMap["ext_id"] = utils.StringValue(ref.ExtId)
			refMap["name"] = utils.StringValue(ref.Name)
			refMap["type"] = utils.StringValue(ref.Type)
			refList[i] = refMap
		}
		return refList
	}
	return []map[string]interface{}{}
}

func flattenEntityReference(ref *monitoringCommon.EntityReference) []map[string]interface{} {
	if ref != nil {
		refMap := make(map[string]interface{})
		refMap["ext_id"] = utils.StringValue(ref.ExtId)
		refMap["name"] = utils.StringValue(ref.Name)
		refMap["type"] = utils.StringValue(ref.Type)
		return []map[string]interface{}{refMap}
	}
	return []map[string]interface{}{}
}

func flattenAuditEntityReference(ref *serviceability.AuditEntityReference) []map[string]interface{} {
	if ref != nil {
		refMap := make(map[string]interface{})
		refMap["ext_id"] = utils.StringValue(ref.ExtId)
		refMap["name"] = utils.StringValue(ref.Name)
		refMap["type"] = utils.StringValue(ref.Type)
		return []map[string]interface{}{refMap}
	}
	return []map[string]interface{}{}
}

func flattenUserReference(ref *serviceability.UserReference) []map[string]interface{} {
	if ref != nil {
		refMap := make(map[string]interface{})
		refMap["ext_id"] = utils.StringValue(ref.ExtId)
		refMap["name"] = utils.StringValue(ref.Name)
		refMap["ip_address"] = utils.StringValue(ref.IpAddress)
		return []map[string]interface{}{refMap}
	}
	return []map[string]interface{}{}
}

func flattenParameters(params []monitoringCommon.Parameter) []map[string]interface{} {
	if len(params) > 0 {
		paramList := make([]map[string]interface{}, len(params))
		for i, param := range params {
			paramMap := make(map[string]interface{})
			paramMap["param_name"] = utils.StringValue(param.ParamName)
			paramMap["param_value"] = flattenParamValue(param.ParamValue)
			paramList[i] = paramMap
		}
		return paramList
	}
	return []map[string]interface{}{}
}

func flattenParamValue(paramValue *monitoringCommon.OneOfParameterParamValue) []map[string]interface{} {
	if paramValue != nil && paramValue.ObjectType_ != nil {
		valueMap := make(map[string]interface{})
		value := paramValue.GetValue()
		if value != nil {
			switch *paramValue.ObjectType_ {
			case "monitoring.v4.common.StringValue":
				if strVal, ok := value.(monitoringCommon.StringValue); ok && strVal.StringValue != nil {
					valueMap["string_value"] = utils.StringValue(strVal.StringValue)
				}
			case "monitoring.v4.common.BoolValue":
				if boolVal, ok := value.(monitoringCommon.BoolValue); ok && boolVal.BoolValue != nil {
					valueMap["bool_value"] = utils.BoolValue(boolVal.BoolValue)
				}
			case "monitoring.v4.common.IntValue":
				if intVal, ok := value.(monitoringCommon.IntValue); ok && intVal.IntValue != nil {
					valueMap["int_value"] = int(utils.Int64Value(intVal.IntValue))
				}
			}
		}
		return []map[string]interface{}{valueMap}
	}
	return []map[string]interface{}{}
}
