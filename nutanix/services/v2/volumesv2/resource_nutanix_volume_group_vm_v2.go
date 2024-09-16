package volumesv2

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	taskPoll "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v16/models/prism/v4/config"
	volumesPrism "github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v16/models/prism/v4/config"
	volumesClient "github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v16/models/volumes/v4/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Attach an AHV VM to the given Volume Group.
func ResourceNutanixVolumeAttachVmToVolumeGroupV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Attaches VM to a Volume Group identified by {extId}.",
		CreateContext: ResourceNutanixVolumeAttachVmToVolumeGroupV2Create,
		ReadContext:   ResourceNutanixVolumeAttachVmToVolumeGroupV2Read,
		UpdateContext: ResourceNutanixVolumeAttachVmToVolumeGroupV2Update,
		DeleteContext: ResourceNutanixVolumeAttachVmToVolumeGroupV2Delete,

		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of the volume group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vm_ext_id": {
				Description: "A globally unique identifier of an instance that is suitable for external consumption. This Field is Required.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"index": {
				Description: "The index on the SCSI bus to attach the VM to the Volume Group. This is an optional field.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"ext_id": {
				Description: "A globally unique identifier of a task.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func ResourceNutanixVolumeAttachVmToVolumeGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtId := d.Get("volume_group_ext_id")

	body := volumesClient.VmAttachment{}

	if vmExtId, ok := d.GetOk("vm_ext_id"); ok {
		body.ExtId = utils.StringPtr(vmExtId.(string))
	}

	if index, ok := d.GetOk("index"); ok {
		body.Index = utils.IntPtr(index.(int))
	}

	resp, err := conn.VolumeAPIInstance.AttachVm(utils.StringPtr(volumeGroupExtId.(string)), &body)

	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while Attaching Vm to Volume Group : %v", errorMessage["message"])
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template (%s) to Attach Vm to Volume Group: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while Attaching Vm to Volume Group: %v", errorMessage["message"])
	}
	rUUID := resourceUUID.Data.GetValue().(taskPoll.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId

	d.SetId(*uuid)
	d.Set("ext_id", *uuid)

	return nil
}

func ResourceNutanixVolumeAttachVmToVolumeGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVolumeAttachVmToVolumeGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVolumeAttachVmToVolumeGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtId := d.Get("volume_group_ext_id")

	body := volumesClient.VmAttachment{}

	if vmExtId, ok := d.GetOk("vm_ext_id"); ok {
		body.ExtId = utils.StringPtr(vmExtId.(string))
	}

	if index, ok := d.GetOk("index"); ok {
		body.Index = utils.IntPtr(index.(int))
	}

	resp, err := conn.VolumeAPIInstance.DetachVm(utils.StringPtr(volumeGroupExtId.(string)), &body)

	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while Detaching Vm to Volume Group : %v", errorMessage["message"])
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template (%s) to Detach Vm to Volume Group: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while Detaching Vm to Volume Group: %v", errorMessage["message"])
	}
	rUUID := resourceUUID.Data.GetValue().(taskPoll.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId

	d.SetId(*uuid)
	d.Set("ext_id", *uuid)

	return nil

}
