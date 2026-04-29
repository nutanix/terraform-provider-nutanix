package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	monitoringService "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixManageAlertV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixManageAlertV2Create,
		ReadContext:   ResourceNutanixManageAlertV2Read,
		DeleteContext: ResourceNutanixManageAlertV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier of an alert that can be resolved or acknowledged.",
			},
			"action_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The action to perform on the alert. Valid values are ACKNOWLEDGE and RESOLVE.",
				ValidateFunc: validation.StringInSlice([]string{
					"ACKNOWLEDGE",
					"RESOLVE",
				}, false),
			},
			"task_ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier for the task.",
			},
		},
	}
}

func ResourceNutanixManageAlertV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Get("ext_id").(string)
	actionTypeStr := d.Get("action_type").(string)

	getResp, err := conn.Alerts.GetAlertById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching alert for ETag: %s", err)
	}

	etagValue := conn.Alerts.ApiClient.GetEtag(getResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	var actionType monitoringService.ActionType
	switch actionTypeStr {
	case "ACKNOWLEDGE":
		actionType = monitoringService.ACTIONTYPE_ACKNOWLEDGE
	case "RESOLVE":
		actionType = monitoringService.ACTIONTYPE_RESOLVE
	}

	body := &monitoringService.AlertActionSpec{
		ActionType: &actionType,
	}

	resp, err := conn.ManageAlerts.ManageAlert(utils.StringPtr(extID), body, args)
	if err != nil {
		return diag.Errorf("error while managing alert: %s", err)
	}

	taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
	if taskRef.ExtId != nil {
		if err := d.Set("task_ext_id", utils.StringValue(taskRef.ExtId)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(extID)
	return nil
}

func ResourceNutanixManageAlertV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixManageAlertV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
