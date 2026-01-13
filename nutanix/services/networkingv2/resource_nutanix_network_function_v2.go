package networkingv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNetworkFunctionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixNetworkFunctionV2Create,
		ReadContext:   ResourceNutanixNetworkFunctionV2Read,
		UpdateContext: ResourceNutanixNetworkFunctionV2Update,
		DeleteContext: ResourceNutanixNetworkFunctionV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Optional: true,
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {Type: schema.TypeString, Computed: true},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {Type: schema.TypeString, Computed: true},
						"rel":  {Type: schema.TypeString, Computed: true},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV2(),
				},
			},
			"name":        {Type: schema.TypeString, Required: true},
			"description": {Type: schema.TypeString, Optional: true},
			"failure_handling": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NO_ACTION", "FAIL_CLOSE", "FAIL_OPEN"}, false),
			},
			"high_availability_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE_PASSIVE"}, false),
			},
			"traffic_forwarding_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"INLINE", "VTAP"}, false),
			},
			"data_plane_health_check_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_threshold": {Type: schema.TypeInt, Optional: true},
						"interval_secs":     {Type: schema.TypeInt, Optional: true},
						"success_threshold": {Type: schema.TypeInt, Optional: true},
						"timeout_secs":      {Type: schema.TypeInt, Optional: true},
					},
				},
			},
			"nic_pairs": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ingress_nic_reference":    {Type: schema.TypeString, Required: true},
						"egress_nic_reference":     {Type: schema.TypeString, Optional: true},
						"is_enabled":               {Type: schema.TypeBool, Required: true},
						"vm_reference":             {Type: schema.TypeString, Optional: true, Computed: true},
						"data_plane_health_status": {Type: schema.TypeString, Computed: true},
						"high_availability_state":  {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func ResourceNutanixNetworkFunctionV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	inputSpec := import1.NetworkFunction{}

	nfName := d.Get("name").(string)
	inputSpec.Name = utils.StringPtr(nfName)

	if v, ok := d.GetOk("description"); ok {
		inputSpec.Description = utils.StringPtr(v.(string))
	}

	if v, ok := d.GetOk("failure_handling"); ok {
		inputSpec.FailureHandling = common.ExpandEnum(v.(string), networkFunctionFailureHandlingMap, "failure_handling")
	}

	ha := common.ExpandEnum(d.Get("high_availability_mode").(string), networkFunctionHighAvailabilityModeMap, "high_availability_mode")
	if ha == nil {
		return diag.Errorf("invalid high_availability_mode: %s", d.Get("high_availability_mode").(string))
	}
	inputSpec.HighAvailabilityMode = ha

	if v, ok := d.GetOk("traffic_forwarding_mode"); ok {
		inputSpec.TrafficForwardingMode = common.ExpandEnum(v.(string), networkFunctionTrafficForwardingModeMap, "traffic_forwarding_mode")
	}

	if v, ok := d.GetOk("data_plane_health_check_config"); ok {
		inputSpec.DataPlaneHealthCheckConfig = expandDataPlaneHealthCheckConfig(v)
	}

	inputSpec.NicPairs = expandNicPairs(d.Get("nic_pairs"))

	resp, err := conn.NetworkFunctionAPI.CreateNetworkFunction(&inputSpec)
	if err != nil {
		return diag.Errorf("error while creating network function : %v", err)
	}

	taskVal := resp.Data.GetValue()
	taskRef, ok := taskVal.(import4.TaskReference)
	if !ok {
		return diag.Errorf("unexpected create network function task type: %T", taskVal)
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	taskRaw, errWait := stateConf.WaitForStateContext(ctx)
	if errWait != nil {
		return diag.Errorf("error waiting for network function (%s) to create: %s", utils.StringValue(taskUUID), errWait)
	}

	if taskDetails, ok := taskRaw.(prismConfig.Task); ok {
		uuid, errUUID := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeNetworkFunction, "Network function")
		if errUUID == nil && uuid != nil && utils.StringValue(uuid) != "" {
			d.SetId(utils.StringValue(uuid))
			_ = d.Set("ext_id", utils.StringValue(uuid))
			return ResourceNutanixNetworkFunctionV2Read(ctx, d, meta)
		}
	}

	// Fallback: lookup created entity by name via List API.
	filter := fmt.Sprintf("name eq '%s'", nfName)
	listResp, errList := conn.NetworkFunctionAPI.ListNetworkFunctions(nil, nil, &filter, nil)
	if errList == nil && listResp != nil && listResp.Data != nil {
		raw := listResp.Data.GetValue()
		if items, ok := raw.([]import1.NetworkFunction); ok && len(items) > 0 && items[0].ExtId != nil {
			d.SetId(utils.StringValue(items[0].ExtId))
			_ = d.Set("ext_id", utils.StringValue(items[0].ExtId))
			return ResourceNutanixNetworkFunctionV2Read(ctx, d, meta)
		}
	}

	return diag.Errorf("network function created but ext_id could not be determined")
}

func ResourceNutanixNetworkFunctionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.NetworkFunctionAPI.GetNetworkFunctionById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching network function : %v", err)
	}

	raw := resp.Data.GetValue()
	var getResp import1.NetworkFunction
	switch v := raw.(type) {
	case import1.NetworkFunction:
		getResp = v
	case *import1.NetworkFunction:
		if v == nil {
			return diag.Errorf("network function response was nil")
		}
		getResp = *v
	default:
		return diag.Errorf("unexpected network function response type: %T", raw)
	}

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("failure_handling", common.FlattenPtrEnum(getResp.FailureHandling)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("high_availability_mode", common.FlattenPtrEnum(getResp.HighAvailabilityMode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("traffic_forwarding_mode", common.FlattenPtrEnum(getResp.TrafficForwardingMode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_plane_health_check_config", flattenDataPlaneHealthCheckConfig(getResp.DataPlaneHealthCheckConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_pairs", flattenNicPairs(getResp.NicPairs)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixNetworkFunctionV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.NetworkFunctionAPI.GetNetworkFunctionById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching network function : %v", err)
	}

	raw := resp.Data.GetValue()
	var current import1.NetworkFunction
	switch v := raw.(type) {
	case import1.NetworkFunction:
		current = v
	case *import1.NetworkFunction:
		if v == nil {
			return diag.Errorf("network function response was nil")
		}
		current = *v
	default:
		return diag.Errorf("unexpected network function response type: %T", raw)
	}

	updateSpec := current

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			updateSpec.Description = utils.StringPtr(v.(string))
		} else {
			updateSpec.Description = nil
		}
	}
	if d.HasChange("failure_handling") {
		if v, ok := d.GetOk("failure_handling"); ok {
			updateSpec.FailureHandling = common.ExpandEnum(v.(string), networkFunctionFailureHandlingMap, "failure_handling")
		} else {
			updateSpec.FailureHandling = nil
		}
	}
	if d.HasChange("high_availability_mode") {
		updateSpec.HighAvailabilityMode = common.ExpandEnum(d.Get("high_availability_mode").(string), networkFunctionHighAvailabilityModeMap, "high_availability_mode")
	}
	if d.HasChange("traffic_forwarding_mode") {
		if v, ok := d.GetOk("traffic_forwarding_mode"); ok {
			updateSpec.TrafficForwardingMode = common.ExpandEnum(v.(string), networkFunctionTrafficForwardingModeMap, "traffic_forwarding_mode")
		} else {
			updateSpec.TrafficForwardingMode = nil
		}
	}
	if d.HasChange("data_plane_health_check_config") {
		if v, ok := d.GetOk("data_plane_health_check_config"); ok {
			updateSpec.DataPlaneHealthCheckConfig = expandDataPlaneHealthCheckConfig(v)
		} else {
			updateSpec.DataPlaneHealthCheckConfig = nil
		}
	}
	if d.HasChange("nic_pairs") {
		updateSpec.NicPairs = expandNicPairs(d.Get("nic_pairs"))
	}

	updateResp, err := conn.NetworkFunctionAPI.UpdateNetworkFunctionById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating network function : %v", err)
	}

	taskVal := updateResp.Data.GetValue()
	taskRef, ok := taskVal.(import4.TaskReference)
	if !ok {
		return diag.Errorf("unexpected update network function task type: %T", taskVal)
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for network function (%s) to update: %s", utils.StringValue(taskUUID), errWait)
	}

	return ResourceNutanixNetworkFunctionV2Read(ctx, d, meta)
}

func ResourceNutanixNetworkFunctionV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.NetworkFunctionAPI.DeleteNetworkFunctionById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting network function : %v", err)
	}

	taskVal := resp.Data.GetValue()
	taskRef, ok := taskVal.(import4.TaskReference)
	if !ok {
		return diag.Errorf("unexpected delete network function task type: %T", taskVal)
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for network function (%s) to delete: %s", utils.StringValue(taskUUID), errWait)
	}

	return nil
}
