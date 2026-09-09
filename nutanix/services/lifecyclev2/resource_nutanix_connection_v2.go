package lifecyclev2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixConnectionV2 manages the lifecycle of a hardware provider connection.
// A connection stores the endpoint and authentication details used to reach an
// external hardware provider. Create/Update/Delete are asynchronous task-based APIs.
func ResourceNutanixConnectionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixConnectionV2Create,
		ReadContext:   ResourceNutanixConnectionV2Read,
		UpdateContext: ResourceNutanixConnectionV2Update,
		DeleteContext: ResourceNutanixConnectionV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"hardware_provider_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "External ID of the hardware provider",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier of the connection",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the connection",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region for the connection",
			},
			"access_details": resourceConnectionAccessDetailsSchema(),
			"deployment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of deployment for the hardware provider connection",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the connection",
			},
			"links":     common.LinksSchema(),
			"tenant_id": {Type: schema.TypeString, Computed: true},
		},
	}
}

func ResourceNutanixConnectionV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)

	body := import1.Connection{}
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("region"); ok {
		body.Region = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("access_details"); ok {
		body.AccessDetails = expandConnectionAccessDetails(v.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Connection create payload : %s", string(aJSON))

	resp, err := conn.HardwareProvidersAPIInstance.CreateConnectionByHardwareProviderId(utils.StringPtr(hpExtID), &body)
	if err != nil {
		return diag.Errorf("error while creating connection : %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for connection (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching connection create task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Connection Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeHardwareProviderConn, "Connection")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return ResourceNutanixConnectionV2Read(ctx, d, meta)
}

func ResourceNutanixConnectionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)

	resp, err := conn.HardwareProvidersAPIInstance.GetConnectionById(utils.StringPtr(hpExtID), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching connection : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Connection)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("region", getResp.Region); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("access_details", flattenConnectionAccessDetails(getResp.AccessDetails)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.DeploymentType != nil {
		if err := d.Set("deployment_type", getResp.DeploymentType.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func ResourceNutanixConnectionV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)

	readResp, err := conn.HardwareProvidersAPIInstance.GetConnectionById(utils.StringPtr(hpExtID), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching connection : %v", err)
	}

	updateSpec := readResp.Data.GetValue().(import1.Connection)

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("region") {
		updateSpec.Region = utils.StringPtr(d.Get("region").(string))
	}
	if d.HasChange("access_details") {
		updateSpec.AccessDetails = expandConnectionAccessDetails(d.Get("access_details").([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", " ")
	log.Printf("[DEBUG] Connection update payload : %s", string(aJSON))

	updateResp, err := conn.HardwareProvidersAPIInstance.UpdateConnectionById(utils.StringPtr(hpExtID), utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating connection : %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for connection (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixConnectionV2Read(ctx, d, meta)
}

func ResourceNutanixConnectionV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)

	resp, err := conn.HardwareProvidersAPIInstance.DeleteConnectionById(utils.StringPtr(hpExtID), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting connection : %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for connection (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
