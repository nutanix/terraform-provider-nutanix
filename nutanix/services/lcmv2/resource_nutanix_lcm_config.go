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

func ResourceLcmConfigV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceLcmConfigV2Create,
		ReadContext:   ResourceLcmConfigV2Read,
		UpdateContext: ResourceLcmConfigV2Update,
		DeleteContext: ResourceLcmConfigV2Delete,
		Schema: map[string]*schema.Schema{
			"ntnx_request_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"if_match": {
				Type:     schema.TypeString,
				Required: true,
			},
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

func ResourceLcmConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterId := d.Get("x_cluster_id").(string)
	ntnxRequestId, ok := d.Get("ntnx_request_id").(string)
	if !ok || ntnxRequestId == "" {
		return diag.Errorf("ntnx_request_id is required and cannot be null or empty")
	}
	if_match, ok := d.Get("if_match").(string)
	if !ok || if_match == "" {
		return diag.Errorf("if_match is required and cannot be null or empty")
	}

	args := make(map[string]interface{})
	args["X-Cluster-Id"] = clusterId

	body := lcmconfigimport1.Config{}

	resp, err := conn.LcmConfigAPIInstance.UpdateConfig(&body, &clusterId, args)
	if err != nil {
		return diag.Errorf("error while updating the LCM config: %v", err)
	}

	getResp := resp.Data.GetValue().(lcmconfigimport1.UpdateConfigApiResponse)
	TaskRef := getResp.Data.GetValue().(taskRef.TaskReference)
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
	return ResourceLcmConfigV2Read(ctx, d, meta)
}

func ResourceLcmConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterId := d.Get("x_cluster_id").(string)

	args := make(map[string]interface{})
	args["X-Cluster-Id"] = clusterId

	resp, err := conn.LcmConfigAPIInstance.GetConfig(&clusterId, args)
	if err != nil {
		return diag.Errorf("error while fetching the Lcm config : %v", err)
	}

	lcmConfig := resp.Data.GetValue().(lcmconfigimport1.Config)
	if err := d.Set("tenant_id", lcmConfig.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(lcmConfig.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_auto_inventory_enabled", lcmConfig.IsAutoInventoryEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auto_inventory_schedule", lcmConfig.AutoInventorySchedule); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", lcmConfig.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_version", lcmConfig.DisplayVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connectivity_type", lcmConfig.ConnectivityType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_https_enabled", lcmConfig.IsHttpsEnabled); err != nil {
		return diag.FromErr(err)
	}
	// if err := d.Set("supported_software_entities", flattenSoftwareEntities(lcmConfig.SupportedSoftwareEntities)); err != nil {
	// 	return diag.FromErr(err)
	// }
	// if err := d.Set("deprecated_software_entities", flattenSoftwareEntities(lcmConfig.DeprecatedSoftwareEntities)); err != nil {
	// 	return diag.FromErr(err)
	// }
	if err := d.Set("is_framework_bundle_uploaded", lcmConfig.IsFrameworkBundleUploaded); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_module_auto_upgrade_enabled", lcmConfig.HasModuleAutoUpgradeEnabled); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*lcmConfig.ExtId)
	return nil
}

func ResourceLcmConfigV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceLcmConfigV2Create(ctx, d, meta)
}

func ResourceLcmConfigV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
