package monitoringv2

import (
	"fmt"

	import1 "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/common"
	import2 "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/response"
	import3 "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func flattenAuditData(data interface{}) (map[string]interface{}, error) {
	audit, ok := data.(*import3.Audit)
	if !ok {
		return nil, fmt.Errorf("failed to cast audit data")
	}

	auditMap := make(map[string]interface{})

	if audit.ExtId != nil {
		auditMap["ext_id"] = utils.StringValue(audit.ExtId)
	}

	if audit.AffectedEntities != nil && len(audit.AffectedEntities) > 0 {
		auditMap["affected_entities"] = flattenEntityReferences(audit.AffectedEntities)
	}

	if audit.AuditType != nil {
		auditMap["audit_type"] = utils.StringValue(audit.AuditType)
	}

	if audit.ClusterReference != nil {
		auditMap["cluster_reference"] = flattenEntityReference(audit.ClusterReference)
	}

	if audit.CreationTime != nil {
		auditMap["creation_time"] = audit.CreationTime.String()
	}

	if audit.Links != nil && len(audit.Links) > 0 {
		auditMap["links"] = flattenApiLinks(audit.Links)
	}

	if audit.Message != nil {
		auditMap["message"] = utils.StringValue(audit.Message)
	}

	if audit.OperationEndTime != nil {
		auditMap["operation_end_time"] = audit.OperationEndTime.String()
	}

	if audit.OperationStartTime != nil {
		auditMap["operation_start_time"] = audit.OperationStartTime.String()
	}

	if audit.OperationType != nil {
		auditMap["operation_type"] = flattenOperationType(audit.OperationType)
	}

	if audit.Parameters != nil && len(audit.Parameters) > 0 {
		auditMap["parameters"] = flattenParameters(audit.Parameters)
	}

	if audit.ServiceName != nil {
		auditMap["service_name"] = utils.StringValue(audit.ServiceName)
	}

	if audit.SourceEntity != nil {
		auditMap["source_entity"] = flattenAuditEntityReference(audit.SourceEntity)
	}

	if audit.Status != nil {
		auditMap["status"] = flattenStatus(audit.Status)
	}

	if audit.TenantId != nil {
		auditMap["tenant_id"] = utils.StringValue(audit.TenantId)
	}

	if audit.UserReference != nil {
		auditMap["user_reference"] = flattenUserReference(audit.UserReference)
	}

	return auditMap, nil
}

func flattenAuditsList(data interface{}) ([]interface{}, error) {
	audits, ok := data.([]import3.Audit)
	if !ok {
		return nil, fmt.Errorf("failed to cast audits list data")
	}

	if len(audits) == 0 {
		return make([]interface{}, 0), nil
	}

	auditsList := make([]interface{}, len(audits))
	for i, audit := range audits {
		auditMap := make(map[string]interface{})

		if audit.ExtId != nil {
			auditMap["ext_id"] = utils.StringValue(audit.ExtId)
		}

		if audit.AffectedEntities != nil && len(audit.AffectedEntities) > 0 {
			auditMap["affected_entities"] = flattenEntityReferences(audit.AffectedEntities)
		}

		if audit.AuditType != nil {
			auditMap["audit_type"] = utils.StringValue(audit.AuditType)
		}

		if audit.ClusterReference != nil {
			auditMap["cluster_reference"] = flattenEntityReference(audit.ClusterReference)
		}

		if audit.CreationTime != nil {
			auditMap["creation_time"] = audit.CreationTime.String()
		}

		if audit.Links != nil && len(audit.Links) > 0 {
			auditMap["links"] = flattenApiLinks(audit.Links)
		}

		if audit.Message != nil {
			auditMap["message"] = utils.StringValue(audit.Message)
		}

		if audit.OperationEndTime != nil {
			auditMap["operation_end_time"] = audit.OperationEndTime.String()
		}

		if audit.OperationStartTime != nil {
			auditMap["operation_start_time"] = audit.OperationStartTime.String()
		}

		if audit.OperationType != nil {
			auditMap["operation_type"] = flattenOperationType(audit.OperationType)
		}

		if audit.Parameters != nil && len(audit.Parameters) > 0 {
			auditMap["parameters"] = flattenParameters(audit.Parameters)
		}

		if audit.ServiceName != nil {
			auditMap["service_name"] = utils.StringValue(audit.ServiceName)
		}

		if audit.SourceEntity != nil {
			auditMap["source_entity"] = flattenAuditEntityReference(audit.SourceEntity)
		}

		if audit.Status != nil {
			auditMap["status"] = flattenStatus(audit.Status)
		}

		if audit.TenantId != nil {
			auditMap["tenant_id"] = utils.StringValue(audit.TenantId)
		}

		if audit.UserReference != nil {
			auditMap["user_reference"] = flattenUserReference(audit.UserReference)
		}

		auditsList[i] = auditMap
	}

	return auditsList, nil
}

func flattenEntityReferences(entities []import1.EntityReference) []map[string]interface{} {
	if len(entities) == 0 {
		return nil
	}

	entityList := make([]map[string]interface{}, len(entities))
	for i, entity := range entities {
		entityMap := make(map[string]interface{})

		if entity.ExtId != nil {
			entityMap["ext_id"] = utils.StringValue(entity.ExtId)
		}

		if entity.Name != nil {
			entityMap["name"] = utils.StringValue(entity.Name)
		}

		if entity.Type != nil {
			entityMap["type"] = utils.StringValue(entity.Type)
		}

		entityList[i] = entityMap
	}

	return entityList
}

func flattenEntityReference(entity *import1.EntityReference) []map[string]interface{} {
	if entity == nil {
		return nil
	}

	entityMap := make(map[string]interface{})

	if entity.ExtId != nil {
		entityMap["ext_id"] = utils.StringValue(entity.ExtId)
	}

	if entity.Name != nil {
		entityMap["name"] = utils.StringValue(entity.Name)
	}

	if entity.Type != nil {
		entityMap["type"] = utils.StringValue(entity.Type)
	}

	return []map[string]interface{}{entityMap}
}

func flattenAuditEntityReference(entity *import3.AuditEntityReference) []map[string]interface{} {
	if entity == nil {
		return nil
	}

	entityMap := make(map[string]interface{})

	if entity.ExtId != nil {
		entityMap["ext_id"] = utils.StringValue(entity.ExtId)
	}

	if entity.Name != nil {
		entityMap["name"] = utils.StringValue(entity.Name)
	}

	if entity.Type != nil {
		entityMap["type"] = utils.StringValue(entity.Type)
	}

	return []map[string]interface{}{entityMap}
}

func flattenApiLinks(links []import2.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}

	linkList := make([]map[string]interface{}, len(links))
	for i, link := range links {
		linkMap := make(map[string]interface{})

		if link.Href != nil {
			linkMap["href"] = utils.StringValue(link.Href)
		}

		if link.Rel != nil {
			linkMap["rel"] = utils.StringValue(link.Rel)
		}

		linkList[i] = linkMap
	}

	return linkList
}

func flattenOperationType(opType *import1.OperationType) string {
	if opType == nil {
		return ""
	}
	return fmt.Sprintf("%v", *opType)
}

func flattenStatus(status *import3.Status) string {
	if status == nil {
		return ""
	}
	return fmt.Sprintf("%v", *status)
}

func flattenParameters(params []import1.Parameter) []map[string]interface{} {
	if len(params) == 0 {
		return nil
	}

	paramList := make([]map[string]interface{}, len(params))
	for i, param := range params {
		paramMap := make(map[string]interface{})

		if param.ParamName != nil {
			paramMap["param_name"] = utils.StringValue(param.ParamName)
		}

		if param.ParamValue != nil {
			paramMap["param_value"] = flattenParameterValue(param.ParamValue)
		}

		paramList[i] = paramMap
	}

	return paramList
}

func flattenParameterValue(paramValue *import1.OneOfParameterParamValue) []map[string]interface{} {
	if paramValue == nil || paramValue.ObjectType_ == nil {
		return nil
	}

	valueMap := make(map[string]interface{})
	value := paramValue.GetValue()

	if value != nil {
		switch *paramValue.ObjectType_ {
		case "monitoring.v4.common.StringValue":
			if strVal, ok := value.(import1.StringValue); ok && strVal.StringValue != nil {
				valueMap["string_value"] = utils.StringValue(strVal.StringValue)
			}
		case "monitoring.v4.common.BoolValue":
			if boolVal, ok := value.(import1.BoolValue); ok && boolVal.BoolValue != nil {
				valueMap["bool_value"] = utils.BoolValue(boolVal.BoolValue)
			}
		case "monitoring.v4.common.IntValue":
			if intVal, ok := value.(import1.IntValue); ok && intVal.IntValue != nil {
				valueMap["int_value"] = int(utils.Int64Value(intVal.IntValue))
			}
		}
	}

	return []map[string]interface{}{valueMap}
}

func flattenUserReference(userRef *import3.UserReference) []map[string]interface{} {
	if userRef == nil {
		return nil
	}

	userMap := make(map[string]interface{})

	if userRef.ExtId != nil {
		userMap["ext_id"] = utils.StringValue(userRef.ExtId)
	}

	if userRef.IpAddress != nil {
		userMap["ip_address"] = utils.StringValue(userRef.IpAddress)
	}

	if userRef.Name != nil {
		userMap["name"] = utils.StringValue(userRef.Name)
	}

	return []map[string]interface{}{userMap}
}
