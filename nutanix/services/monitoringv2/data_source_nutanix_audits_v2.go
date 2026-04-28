package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAuditsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixAuditsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A URL query parameter that specifies the page number of the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A URL query parameter that specifies the total number of records returned in the result set.",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to filter a collection of resources.",
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.",
			},
			"select": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity.",
			},
			"audits": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of audits.",
				Elem:        DatasourceNutanixAuditV2(),
			},
		},
	}
}

func datasourceNutanixAuditsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	var filter, orderBy, selects *string
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

	resp, err := conn.AuditsAPI.ListAudits(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching audits: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("audits", []map[string]interface{}{}); err != nil {
			return diag.Errorf("error setting audits: %v", err)
		}
		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of audits.",
		}}
	}

	getResp := resp.Data.GetValue().([]serviceability.Audit)

	if err := d.Set("audits", flattenAudits(getResp)); err != nil {
		return diag.Errorf("error setting audits: %v", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenAudits(audits []serviceability.Audit) []map[string]interface{} {
	if len(audits) == 0 {
		return []map[string]interface{}{}
	}

	auditsList := make([]map[string]interface{}, len(audits))

	for i, audit := range audits {
		auditMap := make(map[string]interface{})

		auditMap["ext_id"] = utils.StringValue(audit.ExtId)
		auditMap["tenant_id"] = utils.StringValue(audit.TenantId)
		auditMap["links"] = flattenLinks(audit.Links)
		auditMap["audit_type"] = utils.StringValue(audit.AuditType)
		auditMap["message"] = utils.StringValue(audit.Message)
		auditMap["service_name"] = utils.StringValue(audit.ServiceName)

		if audit.OperationType != nil {
			auditMap["operation_type"] = audit.OperationType.GetName()
		}
		if audit.Status != nil {
			auditMap["status"] = audit.Status.GetName()
		}
		if audit.CreationTime != nil {
			auditMap["creation_time"] = audit.CreationTime.String()
		}
		if audit.OperationStartTime != nil {
			auditMap["operation_start_time"] = audit.OperationStartTime.String()
		}
		if audit.OperationEndTime != nil {
			auditMap["operation_end_time"] = audit.OperationEndTime.String()
		}

		auditMap["affected_entities"] = flattenEntityReferences(audit.AffectedEntities)
		auditMap["cluster_reference"] = flattenEntityReference(audit.ClusterReference)
		auditMap["source_entity"] = flattenAuditEntityReference(audit.SourceEntity)
		auditMap["user_reference"] = flattenUserReference(audit.UserReference)
		auditMap["parameters"] = flattenParameters(audit.Parameters)

		auditsList[i] = auditMap
	}

	return auditsList
}
