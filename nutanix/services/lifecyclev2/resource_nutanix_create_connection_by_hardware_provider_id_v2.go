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

// ResourceNutanixCreateConnectionByHardwareProviderIDV2 is an action resource that
// creates a connection to an endpoint for a specific hardware provider. This is the
// standalone create-action surface for the connection create API; the full lifecycle
// (read/update/delete) is managed by nutanix_connection_v2.
func ResourceNutanixCreateConnectionByHardwareProviderIDV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixCreateConnectionByHardwareProviderIDV2Create,
		ReadContext:   ResourceNutanixCreateConnectionByHardwareProviderIDV2Read,
		DeleteContext: ResourceNutanixCreateConnectionByHardwareProviderIDV2Delete,
		Schema: map[string]*schema.Schema{
			"hardware_provider_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "External ID of the hardware provider",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the connection",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Region for the connection",
			},
			"access_details": func() *schema.Schema {
				s := resourceConnectionAccessDetailsSchema()
				s.ForceNew = true
				return s
			}(),
			"ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier of the created connection",
			},
		},
	}
}

func ResourceNutanixCreateConnectionByHardwareProviderIDV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	log.Printf("[DEBUG] Create connection by hardware provider payload : %s", string(aJSON))

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

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeHardwareProviderConn, "Connection")
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", utils.StringValue(uuid)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return ResourceNutanixCreateConnectionByHardwareProviderIDV2Read(ctx, d, meta)
}

func ResourceNutanixCreateConnectionByHardwareProviderIDV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return nil
}

func ResourceNutanixCreateConnectionByHardwareProviderIDV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Action resource: the connection lifecycle (delete) is owned by
	// nutanix_connection_v2. The framework clears the ID after this returns.
	return nil
}
