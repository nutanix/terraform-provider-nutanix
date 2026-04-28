package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAuditV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixAuditV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the generated audit.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
			},
			"links": schemaForLinks(),
			"audit_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name for a given audit type. For example, VMCloneAudit or VMDeleteAudit.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional message associated with the audit.",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service which raised the event or audit. For internal Nutanix services, this value is set to \"Nutanix\".",
			},
			"operation_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operation type of the audit.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the audit.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time in ISO 8601 format when the audit was created.",
			},
			"operation_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The audit operation start time in ISO 8601 format.",
			},
			"operation_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The audit operation end time in ISO 8601 format.",
			},
			"affected_entities": schemaForEntityReferences(),
			"cluster_reference": schemaForEntityReference(),
			"source_entity": schemaForAuditEntityReference(),
			"user_reference": schemaForUserReference(),
			"parameters": schemaForParameters(),
		},
	}
}

func datasourceNutanixAuditV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Get("ext_id")

	resp, err := conn.AuditsAPI.GetAuditById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching audit: %v", err)
	}

	getResp := resp.Data.GetValue().(serviceability.Audit)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("audit_type", getResp.AuditType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("message", getResp.Message); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("service_name", getResp.ServiceName); err != nil {
		return diag.FromErr(err)
	}
	if getResp.OperationType != nil {
		if err := d.Set("operation_type", getResp.OperationType.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.Status != nil {
		if err := d.Set("status", getResp.Status.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.CreationTime != nil {
		if err := d.Set("creation_time", getResp.CreationTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.OperationStartTime != nil {
		if err := d.Set("operation_start_time", getResp.OperationStartTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.OperationEndTime != nil {
		if err := d.Set("operation_end_time", getResp.OperationEndTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("affected_entities", flattenEntityReferences(getResp.AffectedEntities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_reference", flattenEntityReference(getResp.ClusterReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_entity", flattenAuditEntityReference(getResp.SourceEntity)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_reference", flattenUserReference(getResp.UserReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parameters", flattenParameters(getResp.Parameters)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
