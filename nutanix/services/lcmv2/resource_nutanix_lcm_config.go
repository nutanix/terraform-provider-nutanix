package lcmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	lcmconfigimport1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/resources"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixLcmConfigV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixLcmConfigV2Create,
		ReadContext:   ResourceNutanixLcmConfigV2Read,
		UpdateContext: ResourceNutanixLcmConfigV2Update,
		DeleteContext: ResourceNutanixLcmConfigV2Delete,
		Schema: map[string]*schema.Schema{
			"if_match": {
				Type:     schema.TypeString,
				Required: true,
			},
			"x_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_auto_inventory_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_inventory_schedule": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connectivity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_https_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"supported_software_entities": {
				Type:     schema.TypeList,
				Computed: true,
			},
			"deprecated_software_entities": {
				Type:     schema.TypeList,
				Computed: true,
			},
			"is_framework_bundle_uploaded": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_module_auto_upgrade_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixLcmConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterId := d.Get("x_cluster_id").(string)
	if_match, ok := d.Get("if_match").(string)
	if !ok || if_match == "" {
		return diag.Errorf("if_match is required and cannot be null or empty")
	}

	body := lcmconfigimport1.Config{}

	resp, err := conn.LcmConfigAPIInstance.UpdateConfig(&body, utils.StringPtr(clusterId))
	if err != nil {
		return diag.Errorf("error while updating the LCM config: %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI

	// Wait for the Config Update to be successful
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroup(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Config Update task failed: %s", errWaitTask)
	}
	d.SetId(*taskUUID)
	return nil
}

func ResourceNutanixLcmConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixLcmConfigV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixLcmConfigV2Create(ctx, d, meta)
}

func ResourceNutanixLcmConfigV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
