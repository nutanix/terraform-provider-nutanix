package lcmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Optional: true,
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
	var clusterID *string
	if clusterExtID != "" {
		clusterID = utils.StringPtr(clusterExtID)
	} else {
		clusterID = nil
	}
	readResp, err := conn.LcmConfigAPIInstance.GetConfig(clusterID)
	if err != nil {
		return diag.Errorf("error while fetching the Lcm config : %v", err)
	}

	// Extract E-Tag Header
	etagValue := conn.LcmConfigAPIInstance.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := readResp.Data.GetValue().(lcmconfigimport1.Config)
	var connectivityTypeMap = map[string]lcmconfigimport1.ConnectivityType{
		"$UNKNOWN":               lcmconfigimport1.CONNECTIVITYTYPE_UNKNOWN,
		"$REDACTED":              lcmconfigimport1.CONNECTIVITYTYPE_REDACTED,
		"CONNECTED_SITE":         lcmconfigimport1.CONNECTIVITYTYPE_CONNECTED_SITE,
		"DARKSITE_WEB_SERVER":    lcmconfigimport1.CONNECTIVITYTYPE_DARKSITE_WEB_SERVER,
		"DARKSITE_DIRECT_UPLOAD": lcmconfigimport1.CONNECTIVITYTYPE_DARKSITE_DIRECT_UPLOAD,
	}

	if url, ok := d.GetOk("url"); ok {
		body.Url = utils.StringPtr(url.(string))
	}
	if IsExplicitlySet(d, "is_auto_inventory_enabled") {
		v := d.Get("is_auto_inventory_enabled").(bool)
		body.IsAutoInventoryEnabled = utils.BoolPtr(v)
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
	if IsExplicitlySet(d, "is_https_enabled") {
		v := d.Get("is_https_enabled").(bool)
		body.IsHttpsEnabled = utils.BoolPtr(v)
	}
	if IsExplicitlySet(d, "has_module_auto_upgrade_enabled") {
		v := d.Get("has_module_auto_upgrade_enabled").(bool)
		body.HasModuleAutoUpgradeEnabled = utils.BoolPtr(v)
	}
	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] LCM Update Config Request Spec: %s", string(aJSON))

	resp, err := conn.LcmConfigAPIInstance.UpdateConfig(&body, clusterID, args)
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

func IsExplicitlySet(d *schema.ResourceData, key string) bool {
	rawConfig := d.GetRawConfig() // Get raw Terraform config as cty.Value
	log.Printf("[DEBUG] Raw Config: %s", rawConfig)
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return false // If rawConfig is null/unknown, key wasn't explicitly set
	}

	// Convert rawConfig to map and check if key exists
	configMap := rawConfig.AsValueMap()
	if val, exists := configMap[key]; exists {
		log.Printf("[DEBUG1] Key: %s, Value: %s", key, val)
		log.Printf("[DEBUG2] values %t", val.IsNull())
		return !val.IsNull() // Ensure key exists and isn't explicitly null
	}
	return false
}
