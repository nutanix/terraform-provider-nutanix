package lcmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/resources"
	lcmconfigimport1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/resources"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
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
			"x_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_auto_inventory_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_inventory_schedule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connectivity_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_https_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"has_module_auto_upgrade_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func ResourceNutanixLcmConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterExtID := d.Get("x_cluster_id").(string)

	readResp, err := conn.LcmConfigAPIInstance.GetConfig(utils.StringPtr(clusterExtID))
	if err != nil {
		return diag.Errorf("error while fetching the Lcm config : %v", err)
	}

	// Extract E-Tag Header
	etagValue := conn.LcmConfigAPIInstance.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := &lcmconfigimport1.Config{}
	var connectivityTypeMap = map[string]resources.ConnectivityType{
		"$UNKNOWN":               resources.CONNECTIVITYTYPE_UNKNOWN,
		"$REDACTED":              resources.CONNECTIVITYTYPE_REDACTED,
		"CONNECTED_SITE":         resources.CONNECTIVITYTYPE_CONNECTED_SITE,
		"DARKSITE_WEB_SERVER":    resources.CONNECTIVITYTYPE_DARKSITE_WEB_SERVER,
		"DARKSITE_DIRECT_UPLOAD": resources.CONNECTIVITYTYPE_DARKSITE_DIRECT_UPLOAD,
	}

	if url, ok := d.GetOk("url"); ok {
		body.Url = utils.StringPtr(url.(string))
	}

	if isAutoInventoryEnabled, ok := d.GetOk("is_auto_inventory_enabled"); ok {
		body.IsAutoInventoryEnabled = utils.BoolPtr(isAutoInventoryEnabled.(bool))
	}
	if autoInventorySchedule, ok := d.GetOk("auto_inventory_schedule"); ok {
		body.AutoInventorySchedule = utils.StringPtr(autoInventorySchedule.(string))
	}
	if connectivityType, ok := d.GetOk("connectivity_type"); ok {
		if strValue, isString := connectivityType.(string); isString {
			if enumValue, exists := connectivityTypeMap[strValue]; exists {
				body.ConnectivityType = &enumValue
			}
		}
	}
	if isHttpsEnabled, ok := d.GetOk("is_https_enabled"); ok {
		body.IsHttpsEnabled = utils.BoolPtr(isHttpsEnabled.(bool))
	}
	if hasModuleAutoUpgradeEnabled, ok := d.GetOk("has_module_auto_upgrade_enabled"); ok {
		body.HasModuleAutoUpgradeEnabled = utils.BoolPtr(hasModuleAutoUpgradeEnabled.(bool))
	}
	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] LCM Update Config Request Spec: %s", string(aJSON))

	resp, err := conn.LcmConfigAPIInstance.UpdateConfig(body, utils.StringPtr(clusterExtID), args)
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

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching the Lcm upgrade task : %v", err)
	}

	task := resourceUUID.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(task, "", "  ")
	log.Printf("[DEBUG] LCM Config update Task Details: %s", string(aJSON))

	// randomly generating the id
	d.SetId(utils.GenUUID())
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
