package networkingv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/networkfunctions"
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
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV2(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nic_pairs": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ingress_nic_reference": {
							Type:     schema.TypeString,
							Required: true,
						},
						"egress_nic_reference": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"vm_reference": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						// these are computed attributes because they are set by the system
						// come from the API response to indicate the current state of the network function
						// indicates whether the pair is healthy or unhealthy
						"data_plane_health_status": {Type: schema.TypeString, Computed: true},
						// this is computed attribute because it is set by the system
						// come from the API response to indicate the current state of the network function
						// indicates whether the pair is active or passive
						"high_availability_state": {Type: schema.TypeString, Computed: true},
					},
				},
			},
			"high_availability_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(networkFunctionHighAvailabilityModeAllowed, false),
			},
			"failure_handling": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(networkFunctionFailureHandlingAllowed, false),
			},
			"traffic_forwarding_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(networkFunctionTrafficForwardingModeAllowed, false),
			},
			"data_plane_health_check_config": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_threshold": {Type: schema.TypeInt, Optional: true, Computed: true},
						"interval_secs":     {Type: schema.TypeInt, Optional: true, Computed: true},
						"success_threshold": {Type: schema.TypeInt, Optional: true, Computed: true},
						"timeout_secs":      {Type: schema.TypeInt, Optional: true, Computed: true},
					},
				},
			},

			"project_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed attributes
			"ext_id": {
				Optional: true,
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
		inputSpec.FailureHandling = common.ExpandEnum[import1.FailureHandling](v.(string))
	}

	ha := common.ExpandEnum[import1.HighAvailabilityMode](d.Get("high_availability_mode").(string))
	if ha == nil {
		return diag.Errorf("invalid high_availability_mode: %s", d.Get("high_availability_mode").(string))
	}
	inputSpec.HighAvailabilityMode = ha

	if v, ok := d.GetOk("traffic_forwarding_mode"); ok {
		inputSpec.TrafficForwardingMode = common.ExpandEnum[import1.TrafficForwardingMode](v.(string))
	}

	if v, ok := d.GetOk("data_plane_health_check_config"); ok {
		inputSpec.DataPlaneHealthCheckConfig = expandDataPlaneHealthCheckConfig(v)
	}
	if projectExtID, ok := d.GetOk("project_ext_id"); ok {
		inputSpec.ProjectExtId = utils.StringPtr(projectExtID.(string))
	}

	inputSpec.NicPairs = expandNicPairs(d.Get("nic_pairs"))

	createNetworkFunctionRequest := import2.CreateNetworkFunctionRequest{
		Body: &inputSpec,
	}
	resp, err := conn.NetworkFunctionAPI.CreateNetworkFunction(ctx, &createNetworkFunctionRequest)
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
	listNetworkFunctionsRequest := import2.ListNetworkFunctionsRequest{
		Filter_: &filter,
	}
	listResp, errList := conn.NetworkFunctionAPI.ListNetworkFunctions(ctx, &listNetworkFunctionsRequest)
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

	getNetworkFunctionByIDRequest := import2.GetNetworkFunctionByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.NetworkFunctionAPI.GetNetworkFunctionById(ctx, &getNetworkFunctionByIDRequest)
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
	if err := d.Set("project_ext_id", getResp.ProjectExtId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixNetworkFunctionV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("project_ext_id") {
		return diag.Errorf("error while updating project_ext_id: Update of project_ext_id is not supported")
	}
	conn := meta.(*conns.Client).NetworkingAPI

	getNetworkFunctionByIDRequest := import2.GetNetworkFunctionByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.NetworkFunctionAPI.GetNetworkFunctionById(ctx, &getNetworkFunctionByIDRequest)
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
			updateSpec.FailureHandling = common.ExpandEnum[import1.FailureHandling](v.(string))
		} else {
			updateSpec.FailureHandling = nil
		}
	}
	if d.HasChange("high_availability_mode") {
		updateSpec.HighAvailabilityMode = common.ExpandEnum[import1.HighAvailabilityMode](d.Get("high_availability_mode").(string))
	}
	if d.HasChange("traffic_forwarding_mode") {
		if v, ok := d.GetOk("traffic_forwarding_mode"); ok {
			updateSpec.TrafficForwardingMode = common.ExpandEnum[import1.TrafficForwardingMode](v.(string))
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

	updateNetworkFunctionByIDRequest := import2.UpdateNetworkFunctionByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  &updateSpec,
	}
	updateResp, err := conn.NetworkFunctionAPI.UpdateNetworkFunctionById(ctx, &updateNetworkFunctionByIDRequest)
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

	deleteNetworkFunctionByIDRequest := import2.DeleteNetworkFunctionByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.NetworkFunctionAPI.DeleteNetworkFunctionById(ctx, &deleteNetworkFunctionByIDRequest)
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
