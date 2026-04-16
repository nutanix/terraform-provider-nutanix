package lcmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/response"
	lcmstatusimport1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/resources"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLcmStatusV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLcmStatusV2Create,
		Schema: map[string]*schema.Schema{
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"framework_version": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_update_needed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"in_progress_operation": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"is_cancel_intent_set": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"upload_task_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixLcmStatusV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI

	var clusterID *string
	if id := d.Get("x_cluster_id").(string); id != "" {
		clusterID = &id
	}

	resp, err := conn.LcmStatusAPIInstance.GetStatus(clusterID)
	if err != nil {
		return diag.Errorf("error while fetching the Lcm status : %v", err)
	}

	lcmStatusResp := resp.Data.GetValue().(lcmstatusimport1.StatusInfo)

	if err := d.Set("tenant_id", lcmStatusResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(lcmStatusResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_cancel_intent_set", lcmStatusResp.IsCancelIntentSet); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("framework_version", flattenFrameworkVersion(lcmStatusResp.FrameworkVersion)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("in_progress_operation", flattenInProgressOperation(lcmStatusResp.InProgressOperation)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("upload_task_uuid", lcmStatusResp.UploadTaskUuid); err != nil {
		return diag.FromErr(err)
	}

	// setting the data source id to the random uuid
	d.SetId(utils.GenUUID())
	return nil
}

func flattenFrameworkVersion(pr *lcmstatusimport1.FrameworkVersionInfo) []map[string]interface{} {
	if pr != nil {
		frameworkVersionRef := make([]map[string]interface{}, 0)
		frameworkVersion := make(map[string]interface{})

		frameworkVersion["current_version"] = pr.CurrentVersion
		frameworkVersion["available_version"] = pr.AvailableVersion
		frameworkVersion["is_update_needed"] = pr.IsUpdateNeeded

		frameworkVersionRef = append(frameworkVersionRef, frameworkVersion)
		return frameworkVersionRef
	}
	return nil
}

func flattenInProgressOperation(pr *lcmstatusimport1.InProgressOpInfo) []map[string]interface{} {
	if pr != nil {
		OperationRef := make([]map[string]interface{}, 0)
		Operation := make(map[string]interface{})

		operationType := ""
		if pr.OperationType != nil {
			operationType = pr.OperationType.GetName()
		}
		Operation["operation_type"] = operationType
		Operation["operation_id"] = pr.OperationId

		OperationRef = append(OperationRef, Operation)
		return OperationRef
	}
	return nil
}

func flattenLinks(pr []import2.ApiLink) []map[string]interface{} {
	if len(pr) > 0 {
		linkList := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			links := map[string]interface{}{}
			if v.Href != nil {
				links["href"] = v.Href
			}
			if v.Rel != nil {
				links["rel"] = v.Rel
			}

			linkList[k] = links
		}
		return linkList
	}
	return nil
}
