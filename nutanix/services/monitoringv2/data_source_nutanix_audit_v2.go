package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import4 "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/request/audits"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixAuditV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixAuditV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the generated audit.",
			},
			"affected_entities": {
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
			},
			"audit_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name for a given audit type. For example, VMCloneAudit or VMDeleteAudit.",
			},
			"cluster_reference": {
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
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time in ISO 8601 format when the audit was created.",
			},
			"links": {
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
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional message associated with the audit.",
			},
			"operation_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The audit operation end time in ISO 8601 format.",
			},
			"operation_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The audit operation start time in ISO 8601 format.",
			},
			"operation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
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
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service which raised the event or audit. For internal Nutanix services, this value is set to \"Nutanix\".",
			},
			"source_entity": {
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
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
			},
			"user_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique UUID of the user who initiated the operation.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address from where the operation was triggered.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the user who initiated the operation.",
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixAuditV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Get("ext_id").(string)

	req := &import4.GetAuditByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.AuditsAPI.GetAuditById(ctx, req)
	if err != nil {
		return diag.Errorf("error while fetching audit: %v", err)
	}

	if resp.Data == nil {
		return diag.Errorf("audit data is nil for ext_id: %s", extID)
	}

	getResp := resp.Data.GetValue()
	if getResp == nil {
		return diag.Errorf("error getting audit value for ext_id: %s", extID)
	}

	audit, err := flattenAuditData(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range audit {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(extID)
	return nil
}
